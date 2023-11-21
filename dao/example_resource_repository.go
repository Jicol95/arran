package dao

import (
	"database/sql"
	"fmt"
	"time"
)

type ExampleResource struct {
	ID        string
	Name      string
	CreatedAt time.Time
}

type exampleResourcePostgreSQLRepository struct {
	db *sql.DB
}

type ExampleResourceRepository interface {
	Insert(event ExampleResource, tx *sql.Tx) error
	FetchByID(exampleResourceID string, forUpdate bool, tx *sql.Tx) (ExampleResource, error)
	DeleteByID(exampleResourceID string, tx *sql.Tx) error
	FetchAll(tx *sql.Tx) ([]ExampleResource, error)
}

func (r *exampleResourcePostgreSQLRepository) Insert(event ExampleResource, tx *sql.Tx) error {
	query := "INSERT INTO example_resource (id, name, created_at) VALUES ($1, $2, $3);"

	var err error

	if tx != nil {
		_, err = tx.Exec(query, event.ID, event.Name, event.CreatedAt.UTC())
	} else {
		_, err = r.db.Exec(query, event.ID, event.Name, event.CreatedAt.UTC())
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *exampleResourcePostgreSQLRepository) DeleteByID(exampleResourceID string, tx *sql.Tx) error {
	query := "DELETE FROME example_resource WHERE id = $1;"

	var err error

	if tx != nil {
		_, err = tx.Exec(query, exampleResourceID)
	} else {
		_, err = r.db.Exec(query, exampleResourceID)
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *exampleResourcePostgreSQLRepository) FetchByID(exampleResourceID string, forUpdate bool, tx *sql.Tx) (ExampleResource, error) {
	query := "SELECT id, name, created_at FROM example_resource WHERE id = $1 %s;"

	if forUpdate {
		query = fmt.Sprintf(query, "FOR UPDATE")
	} else {
		query = fmt.Sprintf(query, "")
	}

	var row *sql.Row

	if tx != nil {
		row = tx.QueryRow(query, exampleResourceID)
	} else {
		row = r.db.QueryRow(query, exampleResourceID)
	}

	var id string
	var name string
	var created_at time.Time

	if err := row.Scan(&id, &name, &created_at); err != nil {
		return ExampleResource{}, err
	}

	return ExampleResource{id, name, created_at}, nil
}

func (r *exampleResourcePostgreSQLRepository) FetchAll(tx *sql.Tx) ([]ExampleResource, error) {
	query := "SELECT id, name, created_at FROM example_resource;"

	var rows *sql.Rows
	var err error

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = r.db.Query(query)
	}

	if err != nil {
		return nil, err
	}

	var resultList []ExampleResource

	for rows.Next() {
		var id string
		var name string
		var createdAt time.Time

		// Scan the values from the current row into variables
		err := rows.Scan(&id, &name, &createdAt)
		if err != nil {
			return nil, err
		}

		resultList = append(resultList, ExampleResource{ID: id, Name: name, CreatedAt: createdAt})
	}

	return resultList, nil
}

func NewExampleResourceRepository(db *sql.DB) ExampleResourceRepository {
	return &exampleResourcePostgreSQLRepository{
		db: db,
	}
}
