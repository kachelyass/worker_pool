package main

import (
	"context"
	"log"
	"worker_pool/internal/database"
	"worker_pool/internal/handlers"
)

func main() {
	ctx := context.Background()
	db, err := database.Connect(ctx, "postgresql://user:pass@localhost:5432/db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	store := database.NewTaskStore(db)
	taskHandler := handlers.NewTaskHandler(store)

}
