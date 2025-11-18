migrate-up:
	migrate -database "postgres://svc:password@localhost:5432/url_shortener?sslmode=disable" -path cmd/db/migrations up