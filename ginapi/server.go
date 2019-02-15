package main

import (
	"github.com/Chatchaikan/finalexam/customer"
	"github.com/Chatchaikan/finalexam/database"

	_ "github.com/lib/pq"
)

func main() {

	database.Conn()
	customer.CreateTable()

	r := customer.NewRouter()
	r.Run(":2019")
}
