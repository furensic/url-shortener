package repository

import (
	"context"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id, username"

	if err := a.db.QueryRow(ctx, query, &u.Username, &u.PasswordHash).Scan(&newUser.Id, &newUser.Username); err != nil {
		return nil, err
	}

	return &newUser, nil
}
