package main

import (
	"fmt"
	"log"
	"log/slog"
	"mytodoapp/adapters/httpserver"
	"mytodoapp/adapters/persistence/postgre"
	"net/http"
	"os"
)

func main() {
	port := ":8080"
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbDatabase := os.Getenv("POSTGRES_DB")
	dbConnString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s", dbHost, dbPort, dbDatabase, dbUser, dbPassword)
	// FIXME: don't expose db password
	logger.Info("DB conn string", slog.String("conn_string", dbConnString))

	todoStore, err := postgre.NewPostgreTodoStore(dbConnString)
	if err != nil {
		logger.Error("Error creating postgre todo store:", "err", err)
		log.Fatalf("Error creating postgre todo store: %v", err)
	}
	userStore, err := postgre.NewPostgreUserStore(dbConnString)
	if err != nil {
		logger.Error("Error creating postgre user store:", "err", err)
		log.Fatalf("Error creating postgre user store: %v", err)
	}
	server := httpserver.NewTodoServer(logger, todoStore, userStore)
	logger.Info("Server is starting", "address", port)
	log.Fatal(http.ListenAndServe(":8080", server))
}
