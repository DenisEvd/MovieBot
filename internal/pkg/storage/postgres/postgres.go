package postgres

import (
	"MovieBot/internal/pkg/storage"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

const (
	recordsTable  = "records"
	requestsTable = "requests"
)

type Postgres struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

//func (p *Postgres) GetAllMovies(userID int) ([]entities.Movie, error) {
//	return []entities.Movie{}, nil
//}

func (p *Postgres) AddRequest(text string) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (request) values ($1) RETURNING id", requestsTable)
	row := p.db.QueryRow(query, text)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (p *Postgres) DeleteRequest(id int) (string, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return "", err
	}

	var request string
	queryGet := fmt.Sprintf("SELECT r.request FROM %s r WHERE r.id=$1", requestsTable)
	row := tx.QueryRow(queryGet, id)
	if err := row.Scan(&request); err != nil {
		_ = tx.Rollback()
		return "", err
	}

	queryDelete := fmt.Sprintf("DELETE FROM %s r WHERE r.id=$1", requestsTable)
	_, err = p.db.Exec(queryDelete, id)
	if err != nil {
		_ = tx.Rollback()
		return "", err
	}

	return request, tx.Commit()
}

func (p *Postgres) PickRandom(username string) (int, error) {
	query := fmt.Sprintf("SELECT movie_id FROM %s r WHERE r.username=$1 ORDER BY random() LIMIT 1", recordsTable)

	var movieID int
	err := p.db.QueryRow(query, username).Scan(&movieID)

	if err == sql.ErrNoRows {
		return 0, storage.ErrNoSavedMovies
	}

	if err != nil {
		return 0, errors.Wrap(err, "can't scan title from db")
	}

	return movieID, nil
}

func (p *Postgres) AddMovie(username string, movieID int, movieTitle string) error {
	query := fmt.Sprintf("INSERT INTO %s (username, movie_title, movie_id) VALUES ($1, $2, $3)", recordsTable)
	_, err := p.db.Exec(query, username, movieTitle, movieID)
	if err != nil {
		return errors.Wrap(err, "adding record in db")
	}

	return nil
}

func (p *Postgres) Remove(username string, movieID int) error {
	query := fmt.Sprintf("DELETE FROM %s r WHERE r.username=$1 AND r.movie_id=$2", recordsTable)

	_, err := p.db.Exec(query, username, movieID)
	return err
}

func (p *Postgres) IsExists(username string, movieID int) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s r WHERE r.username=$1 AND r.movie_id=$2 LIMIT 1", recordsTable)

	var count int
	err := p.db.Get(&count, query, username, movieID)
	if err != nil {
		return false, errors.Wrap(err, "can't check record in db:")
	}
	return count != 0, nil
}

//func (p *Postgres) Watched() error {
//	return nil
//}
//
//func (p *Postgres) isWatched() (bool, error) {
//	return true, nil
//}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
