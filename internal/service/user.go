package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

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

	newUser := &models.User{}
	newUser.Username = u.Username

	argonConfig := &Argon2Configuration{
		TimeCost:   2,
		MemoryCost: 64 * 1024,
		Threads:    4,
		KeyLength:  32,
	}

	argonConfig.Salt = make([]byte, 16)
	_, err := rand.Read(argonConfig.Salt)
	if err != nil {
		return nil, fmt.Errorf("salt generation failed: %w", err)
	}

	argonConfig.HashRaw = argon2.IDKey([]byte(u.Password), argonConfig.Salt, argonConfig.TimeCost, argonConfig.MemoryCost, argonConfig.Threads, argonConfig.KeyLength)

	// PHC Format: $argon2<variant>$v=<version>$m=<memory>,t=<iterations>,p=<parallelism>$<salt>$<hash>
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, argonConfig.MemoryCost, argonConfig.TimeCost, argonConfig.Threads, base64.RawStdEncoding.EncodeToString(argonConfig.Salt), base64.RawStdEncoding.EncodeToString(argonConfig.HashRaw))
	newUser.PasswordHash = encodedHash

	return s.storage.UserRepository.Create(*newUser)
}
