package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type ShortenedUriModel struct {
	DB *pgx.Conn
}

type ShortenedUri struct {
	Id        int    `json:"id"`
	OriginUri string `json:"origin_uri"`
}

func (m *ShortenedUriModel) Create(s *ShortenedUri) (ShortenedUri, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO shortened_uri (origin_uri) VALUES ($1) RETURNING id"
	var id int
	if err := m.DB.QueryRow(ctx, query, s.OriginUri).Scan(&id); err != nil {
		return ShortenedUri{}, err
	}

	return ShortenedUri{
		Id:        id,
		OriginUri: s.OriginUri,
	}, nil
}

func (m *ShortenedUriModel) GetById(id int) (ShortenedUri, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT * FROM shortened_uri WHERE id=$1"
	var shortenedUri ShortenedUri
	if err := m.DB.QueryRow(ctx, query, id).Scan(&shortenedUri.Id, &shortenedUri.OriginUri); err != nil {
		return ShortenedUri{}, err
	}

	return shortenedUri, nil
}
