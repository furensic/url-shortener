package database

import "github.com/jackc/pgx/v5"

type ShortenedUriModel struct {
	DB *pgx.Conn
}

type ShortenedUri struct {
	Id        int    `json:"id"`
	OriginUri string `json:"origin_uri"`
}

func (m *ShortenedUriModel) Create(s *ShortenedUri) error {
	return nil
}
