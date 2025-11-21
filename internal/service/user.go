package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"codeberg.org/Kassiopeia/url-shortener/internal/models"
	"codeberg.org/Kassiopeia/url-shortener/internal/repository"
	"golang.org/x/crypto/argon2"
)

var ErrPasswordEmpty = errors.New("password cannot be empty.")
var ErrPasswordTooShort = errors.New("passwords needs to be 8 or more characters.")

type Argon2Configuration struct {
	HashRaw    []byte
	Salt       []byte
	TimeCost   uint32
	MemoryCost uint32
	Threads    uint8
	KeyLength  uint32
}

type UserService struct {
	storage repository.Repo
}

func NewUserService(storage repository.Repo) *UserService {
	return &UserService{storage: storage}
}

func (s *UserService) Create(u models.RegisterUserPayload) (*models.User, error) {
	if u.Password == "" {
		return nil, ErrPasswordEmpty
	}
	if len(u.Password) < 8 {
		return nil, ErrPasswordTooShort
	}
	if len(u.Username) >= 32 {
		return nil, errors.New("Username too long")
	}

	_, err := s.GetByUsername(u.Username)
	if err != repository.ErrUsernameNotFound { // maybe create a error for the service?
		return nil, errors.New("User already exists")
	}

	slog.Info("Test")

	newUser := &models.User{}
	newUser.Username = u.Username

	argonConfig := &Argon2Configuration{
		TimeCost:   2,
		MemoryCost: 64 * 1024,
		Threads:    4,
		KeyLength:  32,
	}

	argonConfig.Salt = make([]byte, 16)
	_, err = rand.Read(argonConfig.Salt)
	if err != nil {
		return nil, fmt.Errorf("salt generation failed: %w", err)
	}

	argonConfig.HashRaw = argon2.IDKey([]byte(u.Password), argonConfig.Salt, argonConfig.TimeCost, argonConfig.MemoryCost, argonConfig.Threads, argonConfig.KeyLength)

	// PHC Format: $argon2<variant>$v=<version>$m=<memory>,t=<iterations>,p=<parallelism>$<salt>$<hash>
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, argonConfig.MemoryCost, argonConfig.TimeCost, argonConfig.Threads, base64.RawStdEncoding.EncodeToString(argonConfig.Salt), base64.RawStdEncoding.EncodeToString(argonConfig.HashRaw))
	newUser.PasswordHash = encodedHash

	slog.Info(fmt.Sprintf("%s", encodedHash))

	return s.storage.UserRepository.Create(*newUser)
}

func (s *UserService) GetByUsername(username string) (*models.User, error) {
	if len(username) >= 32 {
		return nil, errors.New("Username too long")
	}
	userFound := &models.User{}

	userFound, err := s.storage.UserRepository.GetByUsername(username)
	if err != nil {
		if err == repository.ErrUsernameNotFound {
			return nil, repository.ErrUsernameNotFound
		}
		return nil, err
	}

	return userFound, nil
}

// returns either true if the credentials are correct or false if incorrect or not found
func (s *UserService) VerifyCredentials(payload models.LoginUserPayload) (bool, error) {
	user, err := s.storage.UserRepository.GetByUsername(payload.Username)
	if err != nil {
		return false, repository.ErrUsernameNotFound
	}

	// PHC Format: $argon2<variant>$v=<version>$m=<memory>,t=<iterations>,p=<parallelism>$<salt>$<hash>
	split := strings.Split(user.PasswordHash, "$")

	if split[0] == "argon2id" {
		slog.Debug("Correct argon2 variant")
	} else {
		slog.Debug("Incorrect argon2 variant")
	}

	return true, nil
}
