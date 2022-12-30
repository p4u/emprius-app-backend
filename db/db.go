package db

import (
	"github.com/genjidb/genji"
)

type Database struct {
	*genji.DB
}

func New(path string) (*Database, error) {
	db, err := genji.Open(path)
	if err != nil {
		return nil, err
	}
	return &Database{db}, nil
}

func (db *Database) Close() error {
	return db.DB.Close()
}

func (db *Database) CreateTables() error {
	if err := createUserTables(db); err != nil {
		return err
	}
	if err := createImageTable(db); err != nil {
		return err
	}
	if err := createTransportTables(db); err != nil {
		return err
	}
	if err := createToolTables(db); err != nil {
		return err
	}
	return nil
}
