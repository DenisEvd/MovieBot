package movie_fetcher

import "MovieBot/internal/events"

type MovieFetcher interface {
	FetchMovieById(id int) (events.Movie, error)
	FetchMoviesByTitle(title string) ([]events.Movie, error)
}
