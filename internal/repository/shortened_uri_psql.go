package repository

import (
	"context"
	"sync"
	"time"

	"codeberg.org/Kassiopeia/url-shortener/internal/models"
	"github.com/jackc/pgx/v5"
)

type ShortenedUriPostgresAdapter struct {
	db     *pgx.Conn
	config RepositoryConfiguration
	dbLock sync.Mutex
}

func NewShortenedUriPgxAdapter(db *pgx.Conn, cfg RepositoryConfiguration) *ShortenedUriPostgresAdapter {
	return &ShortenedUriPostgresAdapter{db: db}
}

func (a *ShortenedUriPostgresAdapter) GetById(id int) (*models.ShortenedUri, error) {

	uri := models.ShortenedUri{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id, origin_uri FROM shortened_uri WHERE id=$1"
	a.dbLock.Lock()
	defer a.dbLock.Unlock()
	if err := a.db.QueryRow(ctx, query, id).Scan(&uri.Id, &uri.OriginUri); err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrShortenedUriNotFound
		}
		return nil, err
	}
	return &uri, nil
}

func (a *ShortenedUriPostgresAdapter) Create(u models.ShortenedUri) (*models.ShortenedUri, error) {
	uri := models.ShortenedUri{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO shortened_uri (origin_uri, timestamp) VALUES ($1, $2) returning id, origin_uri"
	a.dbLock.Lock()
	defer a.dbLock.Unlock()
	if err := a.db.QueryRow(ctx, query, &u.OriginUri, u.Timestamp.Unix()).Scan(&uri.Id, &uri.OriginUri); err != nil {
		return nil, err
	}

	return &uri, nil
}
