package service

import (
	"time"

	"codeberg.org/Kassiopeia/url-shortener/internal/models"
	"codeberg.org/Kassiopeia/url-shortener/internal/repository"
)

type ShortenerService struct {
	storage repository.Repo
}

func NewShortenerService(storage repository.Repo) *ShortenerService {
	return &ShortenerService{storage: storage}
}

func (s *ShortenerService) Create(u models.ShortenedUri) (*models.ShortenedUri, error) {
	// add timestamp to u
	u.Timestamp = time.Now()
	// logic to check uri in request
	return s.storage.ShortenedUriRepository.Create(u)
}

func (s *ShortenerService) GetById(id int) (*models.ShortenedUri, error) {
	return s.storage.ShortenedUriRepository.GetById(id)
}
