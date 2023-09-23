package storage

import (
	"MovieBot/internal/pkg/clients/kinopoisk"
	"errors"
)

type Storage interface {
	AddMovie(username string, movie *kinopoisk.MovieShortInfo)
	PickRandom(username string) (int, error)
	Remove(username string, movieID int) error
	IsExists(username string, movieID int) (bool, error)
}

var ErrNoSavedMovies = errors.New("no saved movies")
