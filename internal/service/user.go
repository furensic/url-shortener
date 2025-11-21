package service

import (
	"errors"

	"codeberg.org/Kassiopeia/url-shortener/internal/models"
	"codeberg.org/Kassiopeia/url-shortener/internal/repository"
	"github.com/alexedwards/argon2id"
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

	hash, err := argon2id.CreateHash(u.Password, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}

	newUser := models.User{
		Username:     u.Username,
		PasswordHash: hash,
	}

	return s.storage.UserRepository.Create(newUser)
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

	match, err := argon2id.ComparePasswordAndHash(payload.Password, user.PasswordHash)
	if err != nil {
		return false, err
	}

	return match, nil
}
