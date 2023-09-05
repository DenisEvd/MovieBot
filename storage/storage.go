package storage

import (
	"MovieBot/entities"
	"errors"
)

type Storage interface {
	AddMovie(movie entities.Movie) error
	PickRandom(username string) (entities.Movie, error)
	Remove(chatID int, movie entities.Movie) error
	IsExists() (bool, error)
}

var ErrNoSavedMovies = errors.New("no saved movies")
