package handlers_test

import (
	"testing"
	"time"

	"codeberg.org/Kassiopeia/url-shortener/cmd/api/handlers"
	"codeberg.org/Kassiopeia/url-shortener/internal/repository"
	"codeberg.org/Kassiopeia/url-shortener/internal/service"
)

func TestHealthHandler(t *testing.T) {
	// really have 0 clue what i need to do in order to import application and classes from my main package?
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

	}

	shortenedUriRepo := repository.NewPostgresAdapter(db)

	repositories := repository.Repo{
		ShortenedUriRepository: shortenedUriRepo,
	}

	shortenerService := service.NewShortenerService(repositories)
	// i wonder how i could do it so that i wouldnt need to build seperate
	// repositories for each service e.g. ShortenerService wouldn't need User service?

	handler := &handlers.Handler{
		ShortenerService: shortenerService,
	}

	app.config = app_config
	app.service = shortenerService
}
