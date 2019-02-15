package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var db *sql.DB

func Conn() *sql.DB {
	if db != nil {
		return db // return existing DB
	}

	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("can't connect to database", err)
	}
	return db
}

func InsertCustomer(name, email, status string) *sql.Row {
	fmt.Println("-> InsertCustomer")
	return Conn().QueryRow("INSERT INTO customers (name, email, status) VALUES ($1, $2, $3) RETURNING id", name, email, status)
}
