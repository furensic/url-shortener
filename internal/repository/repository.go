package repository

import (
	"codeberg.org/Kassiopeia/url-shortener/internal/models"
)

type ShortenedUriRepository interface {
	Create(s *models.ShortenedUri) (models.ShortenedUri, error)
	GetById(id int) (models.ShortenedUri, error)
}

type PostgresDatabase struct {
	connectionString string
}
