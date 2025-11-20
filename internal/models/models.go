package models

import "time"

type ShortenedUri struct {
	Id        int       `json:"id"`
	OriginUri string    `json:"origin_uri"`
	Timestamp time.Time `json:"-"`
}
