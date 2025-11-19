package main

import (
	"context"
	"errors"
	"log"
	"time"

	"codeberg.org/Kassiopeia/url-shortener/cmd/api/handlers"
	"codeberg.org/Kassiopeia/url-shortener/internal/repository"
	"codeberg.org/Kassiopeia/url-shortener/internal/service"
	"github.com/jackc/pgx/v5"
)

type application struct {
	config  config
	models  Models
	service *service.ShortenerService
}

type Models struct {
	ShortenedUri repository.ShortenedUriRepository
}

func NewModels(db *pgx.Conn) Models {
	return Models{
		ShortenedUri: repository.NewShortenedUriRepo(db),
	}
}

type Services struct {
	ShortenedUri service.ShortenerService
}

func NewPostgresDatabase(s string) (*pgx.Conn, error) {
	if s == "" {
		return nil, errors.New("Connection string can't be empty.")
	}

	conn, err := pgx.Connect(context.Background(), s) // change to svc:password@postgres:5432 when using the dockerized version
	if err != nil {
		log.Fatal("Error opening database connection: ", err.Error())
	}

	if err = conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return conn, nil
}

func main() {
	log.Print("Starting url shortener service")

	db, err := NewPostgresDatabase("postgres://svc:password@localhost:5432/url_shortener")
	if err != nil {
		log.Fatal(err.Error())
	}

	shortenedUriRepo := repository.NewShortenedUriRepo(db)

	shortenerService := service.NewShortenerService(shortenedUriRepo) // ?

	app_config := config{
		port:              8090,
		writeTimeout:      3 * time.Second,
		readTimeout:       3 * time.Second,
		readHeaderTimeout: 5 * time.Second,
		idleTimeout:       time.Minute,
	}

	app := &application{
		config:  app_config,
		models:  NewModels(db),
		service: shortenerService,
	}

	handler := &handlers.Handler{
		ShortenerService: shortenerService,
	}

	if err := app.serveHTTP(handler); err != nil {
		log.Fatal("Error listening and serving:", err.Error())
	}
}
