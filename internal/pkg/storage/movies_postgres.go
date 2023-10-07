package storage

import (
	"MovieBot/internal/pkg/events"
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

func (p *MoviesPostgres) PickRandom(username string) (events.Movie, error) {
	querySelectId := fmt.Sprintf("SELECT movie_id FROM %s r WHERE r.username=$1 ORDER BY random() LIMIT 1", recordsTable)

	var movieID int
	err := p.db.QueryRow(querySelectId, username).Scan(&movieID)

	if err == sql.ErrNoRows {
		return events.Movie{}, ErrNoSavedMovies
	}

	if err != nil {
		return events.Movie{}, errors.Wrap(err, "can't scan title from db")
	}

	var movie events.Movie
	querySelectMovie := fmt.Sprintf("SELECT * FROM %s WHERE movie_id=$1 LIMIT 1", moviesTable)
	err = p.db.QueryRow(querySelectMovie, movieID).Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Description, &movie.Poster, &movie.Rating, &movie.Length)
	if err != nil {
		return events.Movie{}, errors.Wrap(err, "select movie")
	}

	return movie, nil
}

func (p *MoviesPostgres) AddMovie(username string, movie *events.Movie) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return err
	}

	queryInsToRecords := fmt.Sprintf("INSERT INTO %s (username, movie_id, is_watched) VALUES ($1, $2, $3)", recordsTable)
	_, err = tx.Exec(queryInsToRecords, username, movie.ID, false)
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "adding record in db")
	}

	isExist, err := p.isExistMovie(movie.ID)
	if err != nil {
		return errors.Wrap(err, "check movie in table")
	}

	if !isExist {
		queryInsToMovies := fmt.Sprintf("INSERT INTO %s (movie_id, title, year, description, poster, rating, length) VALUES ($1, $2, $3, $4, $5, $6, $7)", moviesTable)
		_, err = tx.Exec(queryInsToMovies, movie.ID, movie.Title, movie.Year, movie.Description, movie.Poster, movie.Rating, movie.Length)
		if err != nil {
			_ = tx.Rollback()
			return errors.Wrap(err, "adding movie in db")
		}
	}

	return tx.Commit()
}

func (p *MoviesPostgres) Remove(username string, movieID int) error {
	query := fmt.Sprintf("DELETE FROM %s r WHERE r.username=$1 AND r.movie_id=$2", recordsTable)

	_, err := p.db.Exec(query, username, movieID)
	return err
}

func (p *MoviesPostgres) IsExistRecord(username string, movieID int) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s r WHERE r.username=$1 AND r.movie_id=$2 LIMIT 1", recordsTable)

	var count int
	err := p.db.Get(&count, query, username, movieID)
	if err != nil {
		return false, errors.Wrap(err, "can't check record in db:")
	}
	return count != 0, nil
}

func (p *MoviesPostgres) isExistMovie(movieID int) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE movie_id=$1", moviesTable)

	var count int
	err := p.db.Get(&count, query, movieID)
	if err != nil {
		return false, err
	}

	return count != 0, err
}

//func (p *MoviesPostgres) Watched() error {
//	return nil
//}
//
//func (p *MoviesPostgres) isWatched() (bool, error) {
//	return true, nil
//}
