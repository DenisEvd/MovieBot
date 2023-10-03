package storage

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type MoviesPostgres struct {
	db *sqlx.DB
}

func NewMoviesPostgres(db *sqlx.DB) *MoviesPostgres {
	return &MoviesPostgres{db: db}
}

//func (p *MoviesPostgres) GetAllMovies(userID int) ([]entities.Movie, error) {
//	return []entities.Movie{}, nil
//}

func (p *MoviesPostgres) PickRandom(username string) (int, error) {
	query := fmt.Sprintf("SELECT movie_id FROM %s r WHERE r.username=$1 ORDER BY random() LIMIT 1", recordsTable)

	var movieID int
	err := p.db.QueryRow(query, username).Scan(&movieID)

	if err == sql.ErrNoRows {
		return 0, ErrNoSavedMovies
	}

	if err != nil {
		return 0, errors.Wrap(err, "can't scan title from db")
	}

	return movieID, nil
}

func (p *MoviesPostgres) AddMovie(username string, movieID int, movieTitle string) error {
	query := fmt.Sprintf("INSERT INTO %s (username, movie_title, movie_id) VALUES ($1, $2, $3)", recordsTable)
	_, err := p.db.Exec(query, username, movieTitle, movieID)
	if err != nil {
		return errors.Wrap(err, "adding record in db")
	}

	return nil
}

func (p *MoviesPostgres) Remove(username string, movieID int) error {
	query := fmt.Sprintf("DELETE FROM %s r WHERE r.username=$1 AND r.movie_id=$2", recordsTable)

	_, err := p.db.Exec(query, username, movieID)
	return err
}

func (p *MoviesPostgres) IsExists(username string, movieID int) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s r WHERE r.username=$1 AND r.movie_id=$2 LIMIT 1", recordsTable)

	var count int
	err := p.db.Get(&count, query, username, movieID)
	if err != nil {
		return false, errors.Wrap(err, "can't check record in db:")
	}
	return count != 0, nil
}

//func (p *MoviesPostgres) Watched() error {
//	return nil
//}
//
//func (p *MoviesPostgres) isWatched() (bool, error) {
//	return true, nil
//}
