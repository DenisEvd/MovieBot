package postgres

import (
	"MovieBot/internal/events"
	"MovieBot/internal/storage"
	"fmt"
	"github.com/pkg/errors"
)

func (p *Postgres) GetAll(username string) ([]events.Movie, error) {
	query := fmt.Sprintf("SELECT m.movie_id, m.title, m.year, m.description, m.poster, m.rating, m.length FROM %s r INNER JOIN %s m ON m.movie_id=r.movie_id WHERE r.username=$1 AND r.is_watched=false ORDER BY m.movie_id", recordsTable, moviesTable)

	var movies []events.Movie
	rows, err := p.db.Query(query, username)
	defer func() { _ = rows.Close() }()
	if err != nil {
		return []events.Movie{}, errors.Wrap(err, "error get movies from db")
	}

	for rows.Next() {
		var movie events.Movie
		err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Description, &movie.Poster, &movie.Rating, &movie.Length)
		if err != nil {
			return []events.Movie{}, errors.Wrap(err, "error scan rows")
		}

		movies = append(movies, movie)
	}

	if len(movies) == 0 {
		return []events.Movie{}, storage.ErrNoSavedMovies
	}

	return movies, nil
}

func (p *Postgres) GetNMovie(username string, n int) (events.Movie, error) {
	var count int
	querySelectCount := fmt.Sprintf("SELECT COUNT(*) FROM %s r WHERE r.username=$1 AND r.is_watched=false", recordsTable)
	if err := p.db.QueryRow(querySelectCount, username).Scan(&count); err != nil {
		return events.Movie{}, errors.Wrap(err, "error getting count")
	}

	if count == 0 {
		return events.Movie{}, storage.ErrNoSavedMovies
	}

	var movie events.Movie
	querySelectMovie := fmt.Sprintf("SELECT m.movie_id, m.title, m.year, m.description, m.poster, m.rating, m.length FROM %s r INNER JOIN %s m ON m.movie_id=r.movie_id WHERE r.username=$1 AND r.is_watched=false ORDER BY m.movie_id LIMIT 1 OFFSET $2", recordsTable, moviesTable)
	if err := p.db.QueryRow(querySelectMovie, username, (n-1)%count).Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Description, &movie.Poster, &movie.Rating, &movie.Length); err != nil {
		return events.Movie{}, errors.Wrap(err, "can't select n movie")
	}

	return movie, nil
}

func (p *Postgres) AddMovie(username string, movie events.Movie) error {
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

func (p *Postgres) Watch(username string, movieID int) error {
	query := fmt.Sprintf("UPDATE %s SET is_watched=true WHERE username=$1 AND movie_id=$2", recordsTable)

	_, err := p.db.Exec(query, username, movieID)
	if err != nil {
		return errors.Wrap(err, "update record")
	}

	return nil
}

func (p *Postgres) IsWatched(username string, movieID int) (bool, error) {
	query := fmt.Sprintf("SELECT is_watched FROM %s WHERE username=$1 AND movie_id=$2", recordsTable)

	var isWatched bool
	if err := p.db.QueryRow(query, username, movieID).Scan(&isWatched); err != nil {
		return false, errors.Wrap(err, "getting boolean")
	}

	return isWatched, nil
}

func (p *Postgres) Remove(username string, movieID int) error {
	query := fmt.Sprintf("DELETE FROM %s r WHERE r.username=$1 AND r.movie_id=$2", recordsTable)

	_, err := p.db.Exec(query, username, movieID)
	return err
}

func (p *Postgres) IsExistRecord(username string, movieID int) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s r WHERE r.username=$1 AND r.movie_id=$2", recordsTable)

	var count int
	err := p.db.Get(&count, query, username, movieID)
	if err != nil {
		return false, errors.Wrap(err, "can't check record in db:")
	}
	return count != 0, nil
}

func (p *Postgres) isExistMovie(movieID int) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE movie_id=$1", moviesTable)

	var count int
	err := p.db.Get(&count, query, movieID)
	if err != nil {
		return false, err
	}

	return count != 0, err
}
