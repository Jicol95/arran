package dal

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
	Insert(resource ExampleResource, tx *sql.Tx) error
	FetchByID(exampleResourceID string, forUpdate bool, tx *sql.Tx) (ExampleResource, error)
	Update(resource ExampleResource, tx *sql.Tx) error
	DeleteByID(exampleResourceID string, tx *sql.Tx) error
}

func (r *exampleResourcePostgreSQLRepository) Insert(resource ExampleResource, tx *sql.Tx) error {
	query := "INSERT INTO example_resource (id, name, created_at) VALUES ($1, $2, $3);"

	var err error

	if tx != nil {
		_, err = tx.Exec(query, resource.ID, resource.Name, resource.CreatedAt.UTC())
	} else {
		_, err = r.db.Exec(query, resource.ID, resource.Name, resource.CreatedAt.UTC())
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *exampleResourcePostgreSQLRepository) Update(event ExampleResource, tx *sql.Tx) error {
	query := "UPDATE example_resource SET name = $1 WHERE id = $2;"

	var err error

	if tx != nil {
		_, err = tx.Exec(query, event.Name, event.ID)
	} else {
		_, err = r.db.Exec(query, event.Name, event.ID)
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

func NewExampleResourceRepository(db *sql.DB) ExampleResourceRepository {
	return &exampleResourcePostgreSQLRepository{
		db: db,
	}
}
