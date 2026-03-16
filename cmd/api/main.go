package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"worker_pool/internal/database/postgre"
	"worker_pool/internal/handlers"
)

func main() {
	ctx := context.Background()
	db, err := postgre.Connect(ctx, "postgresql://user:pass@localhost:5432/db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := postgre.NewTaskStore(db)
	taskHandler := handlers.NewTaskHandler(store)
	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.GetALl(w, r)
		case http.MethodPost:
			taskHandler.Create(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid task ID", http.StatusBadRequest)
				return
			}
			taskHandler.GetByID(w, r, id)
		case http.MethodPut:

			taskHandler.Create(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
