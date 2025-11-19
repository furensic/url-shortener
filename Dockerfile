FROM golang:1.24.10

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./tmp/url-shortener ./cmd/api

CMD ["./tmp/url-shortener"]