package db

import (
	"database/sql"
	"log"

	_ "github.com/marcboeker/go-duckdb"
)

var DuckDb *sql.DB

func init() {
	var err error
	DuckDb, err = sql.Open("duckdb", "")
	if err != nil {
		log.Fatal(err)
	}
	if err = DuckDb.Ping(); err != nil {
		log.Fatal(err)
	}
}
