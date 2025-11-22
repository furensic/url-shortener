package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"

	"codeberg.org/Kassiopeia/url-shortener/internal/repository"
	"codeberg.org/Kassiopeia/url-shortener/internal/service"
	"github.com/jackc/pgx/v5"
)

type application struct {
	config  config
	service Services
	logger  slog.Logger
}

type Services struct {
	ShortenerService service.ShortenerService
	UserService      service.UserService
}

func (app *application) NewPostgresDatabase(s string) (*pgx.Conn, error) {
	if s == "" {
		return nil, errors.New("Connection string can't be empty.")
	}

	conn, err := pgx.Connect(context.Background(), s) // change to svc:password@postgres:5432 when using the dockerized version
	if err != nil {
		log.Print("Error opening database connection: ", err.Error())
	}

	if err = conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return conn, nil
}

func CreateLogger(logLevel slog.Level) *slog.Logger {
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	}))
	slog.SetDefault(l)
	return l
}

func main() {
	logger := CreateLogger(slog.LevelDebug)

	app := &application{
		logger: *logger,
	}

	app_config := config{
		port:              8090,
		writeTimeout:      3 * time.Second,
		readTimeout:       3 * time.Second,
		readHeaderTimeout: 5 * time.Second,
		idleTimeout:       time.Minute,
	}

	db, err := app.NewPostgresDatabase("postgres://svc:password@localhost:5432/url_shortener")
	if err != nil {
		logger.Error(err.Error())
	}

	repositories := repository.NewRootRepository()
	repositories.ShortenedUriRepository = repository.NewShortenedUriPgxAdapter(db, repository.RepositoryConfiguration{Logger: logger.WithGroup("ShortenedUriRepository")})
	repositories.UserRepository = repository.NewUserPgxAdapter(db, repository.RepositoryConfiguration{Logger: logger.WithGroup("UserRepository")})

	shortenerService := service.NewShortenerService(*repositories)
	userService := service.NewUserService(*repositories)
	// i wonder how i could do it so that i wouldnt need to build seperate
	// repositories for each service e.g. ShortenerService wouldn't need User service?

	app.config = app_config
	app.service = Services{
		ShortenerService: *shortenerService,
		UserService:      *userService,
	}

	logger.Info("Launching HTTP server on :" + strconv.Itoa(app_config.port))
	if err := app.serveHTTP(); err != nil {
		logger.Error("Error listening and serving: " + err.Error())
	}
}
