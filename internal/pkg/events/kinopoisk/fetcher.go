package kinopoisk

import (
	"MovieBot/internal/pkg/clients/kinopoisk"
	"MovieBot/internal/pkg/events"
	"github.com/pkg/errors"
)

type KpFetcher struct {
	kp    kinopoisk.MovieAPI
	limit int
}

func NewKpFetcher(kp kinopoisk.MovieAPI, limit int) *KpFetcher {
	return &KpFetcher{
		kp:    kp,
		limit: limit,
	}
}

func (k *KpFetcher) FetchMovieById(id int) (events.Movie, error) {
	data, err := k.kp.GetMovieByID(id)
	if err != nil {
		return events.Movie{}, errors.Wrap(err, "getting movie by id")
	}

	movie := events.Movie{
		ID:          data.ID,
		Title:       data.Title,
		Year:        data.Year,
		Description: data.Description,
		Poster:      data.Poster.URL,
		Rating:      data.Rating.IMDB,
		Length:      data.MovieLength,
	}

	return movie, nil
}

func (k *KpFetcher) FetchMoviesByTitle(title string) ([]events.Movie, error) {
	data, err := k.kp.FindMovieByTitle(title, k.limit)
	if err != nil {
		return []events.Movie{}, errors.Wrap(err, "getting movie by title")
	}

	movies := make([]events.Movie, 0, len(data))
	for _, movie := range data {
		movies = append(movies, events.Movie{
			ID:          movie.ID,
			Title:       movie.Title,
			Year:        movie.Year,
			Description: movie.Description,
			Poster:      movie.Poster,
			Rating:      movie.Rating,
			Length:      movie.Length,
		})
	}

	return movies, nil
}
