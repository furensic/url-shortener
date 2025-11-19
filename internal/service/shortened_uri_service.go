package service

import (
	"codeberg.org/Kassiopeia/url-shortener/internal/models"
	"codeberg.org/Kassiopeia/url-shortener/internal/repository"
)

type ShortenedUriService struct {
	Repo repository.ShortenedUriRepository
}

func (s *ShortenedUriService) Create(u *models.ShortenedUri) (models.ShortenedUri, error) {
	// logic to check uri in request
	return s.Repo.Create(u)
}
