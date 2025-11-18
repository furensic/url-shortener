package database

type ShortenedUriModel struct {
}

type ShortenedUri struct {
	Id        int    `json:"id"`
	OriginUri string `json:"origin_uri"`
}

func (m *ShortenedUriModel) Create(s *ShortenedUri) error {
	return nil
}
