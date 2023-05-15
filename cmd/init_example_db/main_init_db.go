package main

import (
	"log"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func main() {

	dsn := "user=foo password=bar dbname=foobar host=localhost port=5432 sslmode=disable"

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	query := `CREATE TABLE users (
					id BIGSERIAL NOT NULL PRIMARY KEY,
					name TEXT NOT NULL UNIQUE,
					registered_at timestamp NOT NULL DEFAULT NOW()
				);`

	db.MustExec(query)

}
