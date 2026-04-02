package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"worker_pool/internal/infrastructure/kafka"
	"worker_pool/internal/infrastructure/postgre"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := postgre.Connect(ctx, "postgresql://user:pass@localhost:5432/db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := postgre.NewTaskStore(db)
	handler := kafka.NewMessageCreateHandler(store)

	consumer, err := kafka.NewConsumer(
		[]string{"redpanda:19092"},
		"task-create-db-writer",
		"task-create-db-writer",
		[]string{"jobs"},
		handler,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	log.Println("kafka consumer started")

	if err := consumer.Consume(ctx); err != nil {
		log.Fatal(err)
	}
}
