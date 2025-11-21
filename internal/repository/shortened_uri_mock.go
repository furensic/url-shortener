package repository

import (
	"codeberg.org/Kassiopeia/url-shortener/internal/models"
)

type ShortenedUriMockAdapter struct {
}

func NewMockAdapter() *ShortenedUriMockAdapter {
	return &ShortenedUriMockAdapter{}
}

func (a *ShortenedUriMockAdapter) GetById(id int) (*models.ShortenedUri, error) {
	uri := models.ShortenedUri{}

	return &uri, nil
}

func (a *ShortenedUriMockAdapter) Create(u models.ShortenedUri) (*models.ShortenedUri, error) {
	uri := models.ShortenedUri{}

	return &uri, nil
}
