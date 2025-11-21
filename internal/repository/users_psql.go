package repository

import (
	"codeberg.org/Kassiopeia/url-shortener/internal/models"
	"github.com/jackc/pgx/v5"
)

type UserPostgresAdapter struct {
	db *pgx.Conn
}

func NewUserPostgresAdapter(db *pgx.Conn) *UserPostgresAdapter {
	return &UserPostgresAdapter{
		db: db,
	}
}

func (a *UserPostgresAdapter) Create(u models.User) (*models.User, error) {
	newUser := models.User{}

	return &newUser, nil
}
