package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func GetDBConnection() *sql.DB {
	if db != nil {
		return db
	}
	var err error
	db, err = sql.Open("sqlite3", "songbooks_of_praise.sqlite")
	if err != nil {
		panic(err)
	}
	return db
}

func Disconnect() {
	if db != nil {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}
}
