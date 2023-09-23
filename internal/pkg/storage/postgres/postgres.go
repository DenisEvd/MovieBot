package postgres

import (
	"MovieBot/internal/lib"
	"MovieBot/internal/pkg/clients/kinopoisk"
	"MovieBot/internal/pkg/storage"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
	recordsTable = "records"
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

func (p *Postgres) PickRandom(username string) (int, error) {
	query := fmt.Sprintf("SELECT movie_id FROM %s r WHERE r.username=$1 LIMIT 1", recordsTable)

	var movieID int
	err := p.db.QueryRow(query, username).Scan(&movieID)

	if err == sql.ErrNoRows {
		return 0, storage.ErrNoSavedMovies
	}

	if err != nil {
		return 0, lib.Wrap("can't scan title from db:", err)
	}

	return movieID, nil
}

func (p *Postgres) AddMovie(username string, movie *kinopoisk.MovieShortInfo) {
	query := fmt.Sprintf("INSERT INTO %s (username, movie_title, movie_id) VALUES ($1, $2, $3)", recordsTable)
	_, _ = p.db.Exec(query, username, movie.Title, movie.ID)
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
		return false, lib.Wrap("can't check record in db:", err)
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
