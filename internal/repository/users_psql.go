package repository

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"codeberg.org/Kassiopeia/url-shortener/internal/models"
	"github.com/jackc/pgx/v5"
)

type UserPostgresAdapter struct {
	db     *pgx.Conn
	config RepositoryConfiguration
}

func NewUserPgxAdapter(db *pgx.Conn, cfg RepositoryConfiguration) *UserPostgresAdapter {
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

var ErrUsernameNotFound = errors.New("Username not found")

func (a *UserPostgresAdapter) GetByUsername(username string) (*models.User, error) {
	userFound := models.User{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	slog.Info(username)

	query := "SELECT id, username FROM users WHERE username=$1"

	err := a.db.QueryRow(ctx, query, username).Scan(&userFound.Id, &userFound.Username)
	if err != nil {
		slog.Error(err.Error())
		if err == pgx.ErrNoRows {
			return nil, ErrUsernameNotFound
		}
		return nil, err
	}

	return &userFound, nil
}

func (a *UserPostgresAdapter) Verify(p models.LoginUserPayload) (*models.User, error) {
	userFound := models.User{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id, username, password_hash FROM users WHERE username=$1"

	err := a.db.QueryRow(ctx, query, p.Username).Scan(&userFound.Id, &userFound.Username, &userFound.PasswordHash)
	if err != nil {
		slog.Error(err.Error())
		if err == pgx.ErrNoRows {
			return nil, ErrUsernameNotFound
		}
		return nil, err
	}

	return &userFound, nil
}
