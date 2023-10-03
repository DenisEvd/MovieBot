package storage

import (
	"errors"
	"github.com/jmoiron/sqlx"
)

type Requests interface {
	AddRequest(text string) (int, error)
	DeleteRequest(id int) (string, error)
}

type Movies interface {
	AddMovie(username string, movieID int, movieTitle string) error
	PickRandom(username string) (int, error)
	Remove(username string, movieID int) error
	IsExists(username string, movieID int) (bool, error)
}

type Storage struct {
	Requests
	Movies
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		Requests: NewRequestsPostgres(db),
		Movies:   NewMoviesPostgres(db),
	}
}

var ErrNoSavedMovies = errors.New("no saved movies")
