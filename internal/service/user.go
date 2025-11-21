package service

import (
	"log"

	"codeberg.org/Kassiopeia/url-shortener/internal/models"
	"codeberg.org/Kassiopeia/url-shortener/internal/repository"
)

type UserService struct {
	storage repository.Repo
}

func NewUserService(storage repository.Repo) *UserService {
	return &UserService{storage: storage}
}

func (s *UserService) Create(u models.RegisterUserRequest) (*models.User, error) {
	newUser := &models.User{}

	log.Print(u.Password, u.Username)

	return newUser, nil
}
