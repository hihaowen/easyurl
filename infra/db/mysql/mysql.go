package mysql

import (
	"database/sql"
	"log"
)

import (
	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func Connect(dsn string) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(50)

	Db = db
}
