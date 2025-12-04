package main

import (
	"log"

	"github.com/fpmoles/go-microservices/internal/database"
	"github.com/fpmoles/go-microservices/internal/server"
)

func main() {
	db, err := database.NewDatabaseClient()
	if err != nil {
		log.Fatal("failed to init Database client: %s", err)
	}
	srv := server.NewEchoServer(db)
	if err := srv.Start(); err != nil {
		log.Fatal(err.Error())
	}
}