package db

import "github.com/jmoiron/sqlx"

var (
	// DBCon is the connection handle for the database
	DBCon *sqlx.DB
)
