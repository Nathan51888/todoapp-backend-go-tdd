package main

import (
	"fmt"
	"log"
	"mytodoapp/adapters/httpserver"
	"mytodoapp/adapters/persistence/postgre"
	"net/http"
	"os"
)

func main() {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbDatabase := os.Getenv("POSTGRES_DB")

	dbConnString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s", dbHost, dbPort, dbDatabase, dbUser, dbPassword)
	log.Printf("db conn string: %s", dbConnString)

	store, err := postgre.NewPostgreTodoStore(dbConnString)
	if err != nil {
		log.Fatalf("Error creating postgre todo store: %v", err)
	}
	server := httpserver.NewTodoServer(store)
	log.Fatal(http.ListenAndServe(":8080", server))
}
