package service

import (
	"log"

	"codeberg.org/Kassiopeia/url-shortener/internal/models"
	"codeberg.org/Kassiopeia/url-shortener/internal/repository"
)

type ShortenerService struct {
	storage repository.ShortenedUriRepository
}

func NewShortenerService(storage repository.ShortenedUriRepository) *ShortenerService {
	return &ShortenerService{storage: storage}
}

func (s *ShortenerService) Create(u models.ShortenedUri) (*models.ShortenedUri, error) {
	// logic to check uri in request
	return s.storage.Create(u)
}

func (s *ShortenerService) GetById(id int) (*models.ShortenedUri, error) {
	log.Print("before s.GetById(id)")

	return s.storage.GetById(id)
}
