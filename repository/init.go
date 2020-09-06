package repository

import (
	"database/sql"
)

type repository struct {
	db *sql.DB
}

// NewRepository is the method for construct a repository which will be registered in the context
func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}
