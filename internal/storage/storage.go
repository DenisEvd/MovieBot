package storage

import (
	"MovieBot/internal/events"
	"errors"
	"github.com/jmoiron/sqlx"
)

type Requests interface {
	AddRequest(text string) (int, error)
	DeleteRequest(id int) (string, error)
}

type Movies interface {
	AddMovie(username string, movie *events.Movie) error
	GetAll(username string) ([]events.Movie, error)
	GetNMovie(username string, n int) (events.Movie, error)
	Watch(username string, movieID int) error
	IsWatched(username string, movieID int) (bool, error)
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
var ErrNoRequest = errors.New("no saved request")
