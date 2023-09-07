package storage

import (
	"MovieBot/entities"
	"errors"
)

type Storage interface {
	AddMovie(username string, movie *entities.Movie)
	PickRandom(username string) (*entities.Movie, error)
	Remove(username string, movie *entities.Movie) error
	IsExists(username string, movie *entities.Movie) (bool, error)
}

var ErrNoSavedMovies = errors.New("no saved movies")
