package main

import (
	"database/sql"
	"log"
	"net/http"

	"example.com/crud-example/internal/api"
	"example.com/crud-example/internal/data"
	_ "github.com/lib/pq"
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

	employeeRepository := data.NewSQLEmployeeRepository(db)
	employeeRepository.Init()

	apiHandler := api.NewAPIHandler(employeeRepository)

	server := &http.Server{
		Addr:    ":8080",
		Handler: apiHandler.Handler(),
	}

	log.Println("Starting API server on " + server.Addr)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
