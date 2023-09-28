package storage

import (
	"errors"
)

type Storage interface {
	AddRequest(text string) (int, error)
	DeleteRequest(id int) (string, error)

	AddMovie(username string, movieID int, movieTitle string) error
	PickRandom(username string) (int, error)
	Remove(username string, movieID int) error
	IsExists(username string, movieID int) (bool, error)
}

var ErrNoSavedMovies = errors.New("no saved movies")
