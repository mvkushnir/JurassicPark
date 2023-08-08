package core

import (
	"database/sql"
	"errors"
	"fmt"
)

func queryDinosaurs(query string, db *sql.DB) ([]Dinosaur, error) {
	var dinosaurs []Dinosaur

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			name    string
			species Species
			cage    int
		)
		if err := rows.Scan(&name, &species, &cage); err != nil {
			return nil, err
		}
		dinosaurs = append(dinosaurs, Dinosaur{name, species, cage})
	}

	return dinosaurs, nil
}

func GetDinosaurs(species string, db *sql.DB) ([]Dinosaur, error) {
	query := "SELECT * FROM dinosaurs"
	if species != "" {
		query += fmt.Sprintf(" WHERE species = '%s'", species)
	}

	return queryDinosaurs(query, db)
}

func GetDinosaur(name string, db *sql.DB) (Dinosaur, error) {
	query := fmt.Sprintf("SELECT * FROM dinosaurs WHERE name = '%s'", name)
	row := db.QueryRow(query)
	var (
		n       string
		species Species
		cage    int
	)
	if err := row.Scan(&n, &species, &cage); err != nil {
		return Dinosaur{}, err
	}
	return Dinosaur{n, species, cage}, nil
}

func GetDinosaursInCage(number int, db *sql.DB) ([]Dinosaur, error) {
	query := fmt.Sprintf("SELECT * FROM dinosaurs WHERE cage = %d", number)
	return queryDinosaurs(query, db)
}

func CageIsAtCapacity(cage Cage, db *sql.DB) (bool, error) {
	dinosaurs, err := GetDinosaursInCage(cage.Number, db)
	if err != nil {
		return false, err
	}
	return len(dinosaurs) >= cage.Capacity, nil
}

func DestinationCageIsValid(dinosaur Dinosaur, cage Cage, db *sql.DB) error {
	// Check that cage is active
	if cage.PowerStatus != ACTIVE {
		return errors.New(fmt.Sprintf("Cage %d is not ACTIVE", cage.Number))
	}

	// Check that cage is not at capacity
	dinosaurs, err := GetDinosaursInCage(cage.Number, db)
	if err != nil {
		return err
	}
	if len(dinosaurs) >= cage.Capacity {
		return errors.New(fmt.Sprintf("Cage %d is at capacity", cage.Number))
	}

	// Check that dinosaurs match species conditions
	// (Carnivores are only with same species, herbivores are not in a cage with carnivores)
	if dinosaur.isCarnivore() {
		for _, other := range dinosaurs {
			if other.Species != dinosaur.Species {
				return errors.New("Carnivores cannot be placed in a cage with different species")
			}
		}
	} else {
		for _, other := range dinosaurs {
			if other.isCarnivore() {
				return errors.New("Herbivores cannot be placed in a cage with carnivores")
			}
		}
	}

	return nil
}

func GetCage(n int, db *sql.DB) (Cage, error) {
	query := fmt.Sprintf("SELECT * FROM cages WHERE number = %d", n)
	row := db.QueryRow(query)
	var (
		number      int
		powerStatus PowerStatus
		capacity    int
	)
	if err := row.Scan(&number, &powerStatus, &capacity); err != nil {
		return Cage{}, err
	}
	return Cage{number, powerStatus, capacity}, nil
}

func GetCages(powerStatus string, db *sql.DB) ([]Cage, error) {
	var cages []Cage

	query := "SELECT * FROM cages"
	if powerStatus != "" {
		query += fmt.Sprintf(" WHERE powerStatus = '%s'", powerStatus)
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			number      int
			powerStatus PowerStatus
			capacity    int
		)
		if err := rows.Scan(&number, &powerStatus, &capacity); err != nil {
			return nil, err
		}
		cages = append(cages, Cage{number, powerStatus, capacity})
	}

	return cages, nil
}

func CreateDinosaur(d Dinosaur, db *sql.DB) error {
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO dinosaurs(name, species, cage) VALUES('%s', '%s', %d)",
		d.Name, d.Species, d.Cage))
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}

func MoveDinosaur(name string, cage int, db *sql.DB) error {
	stmt, err := db.Prepare(fmt.Sprintf("UPDATE dinosaurs SET cage = %d WHERE name = '%s'", cage, name))
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}

func UpdateCage(number int, powerStatus string, db *sql.DB) error {
	stmt, err := db.Prepare(fmt.Sprintf("UPDATE cages SET powerStatus = '%s' WHERE number = %d", powerStatus, number))
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}

func CreateCage(c Cage, db *sql.DB) error {
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO cages(number, powerstatus, capacity) VALUES(%d, '%s', %d)",
		c.Number, c.PowerStatus, c.Capacity))
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}
