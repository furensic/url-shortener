package service

import (
	"codeberg.org/Kassiopeia/url-shortener/internal/models"
	"codeberg.org/Kassiopeia/url-shortener/internal/repository"
)

type ShortenerService struct {
	storage *repository.ShortenedUriRepository
}

func NewShortenerService(storage repository.ShortenedUriRepository) *ShortenerService {
	return &ShortenerService{storage: &storage}
}

func (s *ShortenerService) Create(u models.ShortenedUri) (models.ShortenedUri, error) {
	// logic to check uri in request
	return models.ShortenedUri{}, nil
}

func (s *ShortenerService) GetById(id int) (models.ShortenedUri, error) {
	return s.GetById(id)
}
