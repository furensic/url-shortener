package models

import "time"

type ShortenedUri struct {
	Id        int       `json:"id"`
	OriginUri string    `json:"origin_uri"`
	Timestamp time.Time `json:"-"`
}

type User struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}

type RegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
