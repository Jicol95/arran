package dao

import "database/sql"

type TransactionManager interface {
	BeginTx() (*sql.Tx, error)
}

type postgresTransactionManager struct {
	db *sql.DB
}

func (tm *postgresTransactionManager) BeginTx() (*sql.Tx, error) {
	return tm.db.Begin()
}

func NewTransactionManager(db *sql.DB) TransactionManager {
	return &postgresTransactionManager{db: db}
}
