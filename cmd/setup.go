package main

import (
	"JurassicPark/core"
	"fmt"
)

func main() {
	db := core.Connect()
	defer db.Close()

	// Create cages table
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS cages (number INT NOT NULL PRIMARY KEY AUTO_INCREMENT, powerStatus VARCHAR(255) NOT NULL, capacity INT NOT NULL)")
	if err != nil {
		panic(err)
	}

	// Create dinosaurs table
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS dinosaurs (name VARCHAR(255) NOT NULL PRIMARY KEY, species VARCHAR(255) NOT NULL, cage INT NOT NULL, FOREIGN KEY (cage) REFERENCES cages(number))")
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully created tables")
}
