package database

type Models struct {
	ShortenedUri ShortenedUriModel
}

func NewModels() Models {
	return Models{
		ShortenedUri: ShortenedUriModel{},
	}
}
