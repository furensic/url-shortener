package repository

import (
	"context"
	"time"

	"codeberg.org/Kassiopeia/url-shortener/internal/models"
	"github.com/jackc/pgx/v5"
)

type ShortenedUriPostgres struct {
	Db *pgx.Conn
}

var _ ShortenedUriRepository = (*ShortenedUriPostgres)(nil)

func NewShortenedUriRepo(db *pgx.Conn) *ShortenedUriPostgres {
	return &ShortenedUriPostgres{Db: db}
}

func (m *ShortenedUriPostgres) Create(s models.ShortenedUri) (models.ShortenedUri, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO shortened_uri (origin_uri) VALUES ($1) RETURNING id"
	var id int
	if err := m.Db.QueryRow(ctx, query, s.OriginUri).Scan(&id); err != nil {
		return models.ShortenedUri{}, err
	}

	return models.ShortenedUri{
		Id:        id,
		OriginUri: s.OriginUri, // maybe later use actual data returned by sql?
	}, nil
}

func (m *ShortenedUriPostgres) GetById(id int) (models.ShortenedUri, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT * FROM shortened_uri WHERE id=$1"
	var shortenedUri models.ShortenedUri
	if err := m.Db.QueryRow(ctx, query, id).Scan(&shortenedUri.Id, &shortenedUri.OriginUri); err != nil {
		return models.ShortenedUri{}, err
	}

	return shortenedUri, nil
}
