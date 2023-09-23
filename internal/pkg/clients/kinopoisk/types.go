package kinopoisk

type MovieAPI interface {
	FindMovieByTitle(title string, limit int) ([]MovieShortInfo, error)
	GetMovieByID(movieID int) (Movie, error)
}

type MovieResponse struct {
	Docs []MovieShortInfo `json:"docs"`
}

type MovieShortInfo struct {
	ID     int     `json:"id"`
	Title  string  `json:"name"`
	Year   int     `json:"year"`
	Rating float32 `json:"rating"`
	Length int     `json:"movieLength"`
}

type Movie struct {
	Title       string  `json:"name"`
	Year        int     `json:"year"`
	Description string  `json:"shortDescription"`
	Poster      Poster  `json:"poster"`
	Rating      Ratings `json:"rating"`
	MovieLength int     `json:"movieLength"`
}

type Ratings struct {
	IMDB float32 `json:"imdb"`
}

type Poster struct {
	URL string `json:"url"`
}
