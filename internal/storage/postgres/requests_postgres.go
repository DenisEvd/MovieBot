package postgres

import (
	"MovieBot/internal/storage"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
)

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
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrNoRequest
		}
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
