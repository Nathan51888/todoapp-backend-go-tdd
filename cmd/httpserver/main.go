package main

import (
	"fmt"
	"log"
	"mytodoapp/adapters/httpserver"
	"mytodoapp/adapters/persistence/postgre"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbUrl := os.Getenv("DATABASE_URL")
	fmt.Printf("db url: %s\n", dbUrl)

	store := postgre.NewPostgreTodoStore(dbUrl)
	server := httpserver.NewTodoServer(store)
	log.Fatal(http.ListenAndServe(":8080", server))
}
