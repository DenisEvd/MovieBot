package kinopoisk

type MovieAPI interface {
	FindMovieByTitle(title string, limit int) ([]MovieByTitle, error)
	GetMovieByID(movieID int) (MovieByID, error)
}

type MovieResponse struct {
	Docs []MovieByTitle `json:"docs"`
}

type MovieByTitle struct {
	ID          int     `json:"id"`
	Title       string  `json:"name"`
	Year        int     `json:"year"`
	Description string  `json:"shortDescription"`
	Poster      string  `json:"poster"`
	Rating      float32 `json:"rating"`
	Length      int     `json:"movieLength"`
}

type MovieByID struct {
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
