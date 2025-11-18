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

func (m *ShortenedUriModel) Create(s *ShortenedUri) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO shortened_uri (origin_uri) VALUES ($1) RETURNING id"
	var id int
	if err := m.DB.QueryRow(ctx, query, s.OriginUri).Scan(&id); err != nil {
		return err
	}

	return nil
}
