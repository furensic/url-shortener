package database

import "github.com/jackc/pgx/v5"

type Models struct {
	ShortenedUri ShortenedUriModel
}

func NewModels(db *pgx.Conn) Models {
	return Models{
		ShortenedUri: ShortenedUriModel{DB: db},
	}
}
