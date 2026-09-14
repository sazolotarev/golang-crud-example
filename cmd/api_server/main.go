package main

import (
	"database/sql"
	"log"

	"example.com/crud-example/internal/api"
	"example.com/crud-example/internal/data"
)

func main() {
	// TODO: load DB connection settings from config/environment
	db, err := sql.Open("postgres", "host=postgres user=crud_user password=crud_password dbname=crud_example sslmode=disable connect_timeout=5")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	employeeRepository := data.NewEmployeeRepository(db)
	employeeRepository.Init()

	server := api.NewAPIServer(":8080", employeeRepository)
	server.Run()
}
