package core

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	// In a real environment I'd make these configurable, to allow for easily changing the values
	// between prod, dev, and local envs
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/jurassicpark")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Successfully connected to database")

	return db
}
