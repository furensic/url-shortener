package handlers_test

import (
	"testing"

	"codeberg.org/Kassiopeia/url-shortener/cmd/api/handlers"
	"codeberg.org/Kassiopeia/url-shortener/internal/repository"
	"codeberg.org/Kassiopeia/url-shortener/internal/service"
	// "codeberg.org/Kassiopeia/url-shortener/internal/repository"
	// "codeberg.org/Kassiopeia/url-shortener/internal/service"
)

func TestHealthHandler(t *testing.T) {
	// really have 0 clue what i need to do in order to import application and classes from my main package?

	shortenedUriRepo := repository.NewMockAdapter()

	repositories := repository.Repo{
		ShortenedUriRepository: shortenedUriRepo,
	}

	shortenerService := service.NewShortenerService(repositories)

	_ = &handlers.Handler{
		ShortenerService: shortenerService,
	}

}
