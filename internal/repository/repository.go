package repository

import (
	"log/slog"

	"codeberg.org/Kassiopeia/url-shortener/internal/models"
)

type RepositoryConfiguration struct {
	Logger *slog.Logger
}

type Repo struct {
	ShortenedUriRepository interface {
		Create(s models.ShortenedUri) (*models.ShortenedUri, error)
		GetById(id int) (*models.ShortenedUri, error)
	}
	UserRepository interface {
		Create(u models.User) (*models.User, error)
		GetByUsername(username string) (*models.User, error)
		Verify(p models.LoginUserPayload) (*models.User, error)
	}
}

func NewRootRepository() *Repo {
	return &Repo{}
}
