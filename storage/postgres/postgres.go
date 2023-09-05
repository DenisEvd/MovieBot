package postgres

import (
	"MovieBot/entities"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type Postgres struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

func (p *Postgres) GetAllMovies(userID int) ([]entities.Movie, error) {
	return []entities.Movie{}, nil
}

func (p *Postgres) PickRandom(username string) (entities.Movie, error) {
	return entities.Movie{}, nil
}

func (p *Postgres) AddMovie(movie entities.Movie) error {
	return nil
}

func (p *Postgres) Remove(chatID int, movie entities.Movie) error {
	return nil
}

func (p *Postgres) IsExists() (bool, error) {
	return false, nil
}

func (p *Postgres) Watched() error {
	return nil
}

func (p *Postgres) isWatched() (bool, error) {
	return true, nil
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
