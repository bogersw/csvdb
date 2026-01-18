package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("sqlite", "")
	if err != nil {
		log.Fatal(err)
	}
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}
}
