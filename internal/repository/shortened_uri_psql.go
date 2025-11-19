package repository

import (
	"context"
	"time"

	"codeberg.org/Kassiopeia/url-shortener/internal/models"
	"github.com/jackc/pgx/v5"
)

type ShortenedUriRepoAdapter struct {
	db *pgx.Conn
}

func NewShortenedUriRepo(db *pgx.Conn) *ShortenedUriRepoAdapter {
	return &ShortenedUriRepoAdapter{db: db}
}

func (a *ShortenedUriRepoAdapter) GetById(id int) (*models.ShortenedUri, error) {
	uri := models.ShortenedUri{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id, origin_uri FROM shortened_uri WHERE id=$1"

	if err := a.db.QueryRow(ctx, query, id).Scan(&uri.Id, &uri.OriginUri); err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrShortenedUriNotFound
		}
		return nil, err
	}

	return &uri, nil
}

func (a *ShortenedUriRepoAdapter) Create(u models.ShortenedUri) (*models.ShortenedUri, error) {
	uri := models.ShortenedUri{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO shortened_uri (origin_uri) VALUES ($1) returning id, origin_uri"

	if err := a.db.QueryRow(ctx, query, &u.OriginUri).Scan(&uri.Id, &uri.OriginUri); err != nil {
		return nil, err
	}

	return &uri, nil
}
