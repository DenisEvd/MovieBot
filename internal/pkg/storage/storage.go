package storage

import (
	"MovieBot/internal/pkg/events"
	"errors"
	"github.com/jmoiron/sqlx"
)

type Requests interface {
	AddRequest(text string) (int, error)
	DeleteRequest(id int) (string, error)
}

type Movies interface {
	AddMovie(username string, movie *events.Movie) error
	PickRandom(username string) (events.Movie, error)
	Remove(username string, movieID int) error
	IsExistRecord(username string, movieID int) (bool, error)
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
