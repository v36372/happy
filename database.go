package happy

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattes/migrate/database/postgres"
)

// PGDB is wrapper of sqlx.DB
type PGDB struct {
	*sqlx.DB
}

// Tx is wrapper of Tx (transaction)
type Tx struct {
	*sqlx.Tx
}

// OpenDB returns a DB reference for a Postgres data source
func OpenDB(dataSourceString string) (*PGDB, error) {
	db, err := sqlx.Open("postgres", dataSourceString)
	if err != nil {
		return nil, err
	}
	return &PGDB{db}, nil
}

// MustOpenDB returns a DB reference for a Postgres data source
// and panics if there's any error
func MustOpenDB(dataSourceString string) *PGDB {
	return &PGDB{sqlx.MustOpen("postgres", dataSourceString)}
}

// Begin starts and returns a new transaction
func (db *PGDB) Begin() (*Tx, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}
