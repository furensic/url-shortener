package models

type ShortenedUri struct {
	Id        int    `json:"id"`
	OriginUri string `json:"origin_uri"`
}
