package service

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
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
	user, err := s.storage.UserRepository.Verify(payload)
	if err != nil {
		return false, repository.ErrUsernameNotFound
	}

	// PHC Format: $argon2<variant>$v=<version>$m=<memory>,t=<iterations>,p=<parallelism>$<salt>$<hash>
	split := strings.Split(user.PasswordHash, "$")

	slog.Info(fmt.Sprintf("%+v", user))
	slog.Info(split[1])
	if split[1] == "argon2id" {
		slog.Info("Correct argon2 variant")
	} else {
		return false, errors.New("Invalid argon2 variant")
	}

	argonConfig := struct {
		variant    string
		version    int
		memory     int
		iterations int
		threads    int
		salt       string
		hash       string
	}{
		variant: split[1],
	}

	for i, str := range split {
		slog.Info(fmt.Sprintf("i: %d, str: %v", i, str))
	}

	ver, ok := strings.CutPrefix(split[2], "v=")
	if !ok {
		slog.Info("Wrong version")
		return false, err
	}
	params := strings.Split(split[3], ",")
	slog.Info(fmt.Sprintf("params: %v", params))
	mem, ok := strings.CutPrefix(params[0], "m=")
	slog.Info(fmt.Sprintf("mem: %v", mem))
	if !ok {
		slog.Info("Wrong memory")
		return false, err
	}
	iter, ok := strings.CutPrefix(params[1], "t=")
	slog.Info(fmt.Sprintf("iter: %v", iter))
	if !ok {
		slog.Info("Wrong iterations")
		return false, err
	}
	thr, ok := strings.CutPrefix(params[2], "p=")
	slog.Info(fmt.Sprintf("thr: %v", thr))
	if !ok {
		slog.Info("Wrong thread")
		return false, err
	}

	argonConfig.version, err = strconv.Atoi(ver)
	argonConfig.memory, err = strconv.Atoi(mem)
	argonConfig.iterations, err = strconv.Atoi(iter)
	argonConfig.threads, err = strconv.Atoi(thr)
	argonConfig.salt = split[4]
	argonConfig.hash = split[5]

	slog.Info(fmt.Sprintf("salt: %s", argonConfig.salt))
	slog.Info(fmt.Sprintf("hash: %s", argonConfig.hash))
	inputPassword := argon2.IDKey([]byte(payload.Password), []byte(argonConfig.salt), uint32(argonConfig.iterations), uint32(argonConfig.memory), uint8(argonConfig.threads), 32)

	slog.Info(fmt.Sprintf("inputPassword: %s\ndbPassword: %s", base64.RawStdEncoding.EncodeToString(inputPassword), argonConfig.hash))
	if err != nil {
		slog.Info(err.Error())
	}
	if !bytes.Equal(inputPassword, []byte(argonConfig.hash)) {
		slog.Info("invalid password")
		return false, errors.New("Invalid password")
	}

	return true, nil
}
