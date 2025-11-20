package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"

	"codeberg.org/Kassiopeia/url-shortener/cmd/api/handlers"
	"codeberg.org/Kassiopeia/url-shortener/internal/repository"
	"codeberg.org/Kassiopeia/url-shortener/internal/service"
	"github.com/jackc/pgx/v5"
)

type application struct {
	config  config
	service *service.ShortenerService
	logger  slog.Logger
}

type Services struct {
	ShortenedUri service.ShortenerService
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

func main() {
	logger_lvl := new(slog.LevelVar)
	logger_lvl.Set(slog.LevelDebug)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logger_lvl,
	}))

	logger.Info("Starting url shortener service")

	logger.Debug("Initializing new application struct")
	app := &application{
		logger: *logger,
	}

	logger.Debug("Creating new app config")
	app_config := config{
		port:              8090,
		writeTimeout:      3 * time.Second,
		readTimeout:       3 * time.Second,
		readHeaderTimeout: 5 * time.Second,
		idleTimeout:       time.Minute,
	}
	logger.Debug("New app config: " + fmt.Sprintf("%+v", app_config))

	logger.Debug("Connecting to database")
	db, err := app.NewPostgresDatabase("postgres://svc:password@localhost:5432/url_shortener")
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Debug("Database connection established to " + db.Config().Host)

	logger.Debug("Creating new ShortenedUri Repository with postgres adapter")
	shortenedUriRepo := repository.NewPostgresAdapter(db)

	logger.Debug("Creating new Repository")
	repositories := repository.Repo{
		ShortenedUriRepository: shortenedUriRepo,
	}

	logger.Debug("Creating new ShortenerService using repository")
	shortenerService := service.NewShortenerService(repositories)
	// i wonder how i could do it so that i wouldnt need to build seperate
	// repositories for each service e.g. ShortenerService wouldn't need User service?

	logger.Debug("Creating new handlers with ShortenerService")
	handler := &handlers.Handler{
		ShortenerService: shortenerService,
	}

	logger.Debug("Updating app config")
	app.config = app_config
	logger.Debug("Updating app services")
	app.service = shortenerService

	logger.Info("Launching HTTP server on :" + strconv.Itoa(app_config.port))
	if err := app.serveHTTP(handler); err != nil {
		logger.Error("Error listening and serving: " + err.Error())
	}
}
