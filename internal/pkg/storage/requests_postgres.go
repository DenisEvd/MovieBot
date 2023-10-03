package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type RequestsPostgres struct {
	db *sqlx.DB
}

func NewRequestsPostgres(db *sqlx.DB) *RequestsPostgres {
	return &RequestsPostgres{db: db}
}

func (p *RequestsPostgres) AddRequest(text string) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (request) values ($1) RETURNING id", requestsTable)
	row := p.db.QueryRow(query, text)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (p *RequestsPostgres) DeleteRequest(id int) (string, error) {
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
