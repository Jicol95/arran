package dao

import (
	"database/sql"
	"fmt"
	"time"
)

type ProcessedEvent struct {
	ID          string    `json:"id"`
	Source      string    `json:"eventId"`
	ProcessedAt time.Time `json:"processedAt"`
}

type processedEventPostgreSQLRepository struct {
	db *sql.DB
}

type ProcessedEventRepository interface {
	Insert(event ProcessedEvent, tx *sql.Tx) error
	FetchByEventID(eventID string, forUpdate bool, tx *sql.Tx) (ProcessedEvent, error)
}

func (r *processedEventPostgreSQLRepository) Insert(event ProcessedEvent, tx *sql.Tx) error {
	query := "INSERT INTO processed_event (id, source, processed_at) VALUES ($1, $2, $3);"

	var err error

	if tx != nil {
		_, err = tx.Exec(query, event.ID, event.Source, event.ProcessedAt.UTC())
	} else {
		_, err = r.db.Exec(query, event.ID, event.Source, event.ProcessedAt.UTC())
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *processedEventPostgreSQLRepository) FetchByEventID(eventID string, forUpdate bool, tx *sql.Tx) (ProcessedEvent, error) {
	query := "SELECT id, source, processed_at FROM processed_event WHERE id = $1 %s"

	if forUpdate {
		query = fmt.Sprintf(query, "FOR UPDATE")
	} else {
		query = fmt.Sprintf(query, "")
	}

	var row *sql.Row

	if tx != nil {
		row = tx.QueryRow(query)
	} else {
		row = r.db.QueryRow(query)
	}

	var id string
	var source string
	var processedAt time.Time

	if err := row.Scan(&id, &source, &processedAt); err != nil {
		return ProcessedEvent{}, err
	}

	return ProcessedEvent{id, source, processedAt}, nil
}

func NewProcessedEventRepository() ProcessedEventRepository {
	return &processedEventPostgreSQLRepository{}
}
