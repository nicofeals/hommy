package pgsql

import (
	"github.com/go-pg/pg"
)

type Database struct {
	db *pg.DB
}

func NewDatabase(o *pg.Options) *Database {
	return &Database{
		db: pg.Connect(o),
	}
}
