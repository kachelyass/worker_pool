package handlers

import (
	"net/http"
	"worker_pool/internal/database"
)

type TaskHandler struct {
	store *database.TaskStore
}

func NewTaskHandler(store *database.TaskStore) *TaskHandler {
	return &TaskHandler{store: store}
}

func (h *TaskHandler) GetALl(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.GetAll(r.Context())
	if err != nil {
		writerError(w, http.StatusInternalServerError, "Failed to get all tasks")
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request, id int) {
	task, err := h.store.GetByID(r.Context(), id)
	if err != nil {
		writerError(w, http.StatusInternalServerError, "Failed to get task")
	}
	writeJSON(w, http.StatusOK, task)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
}

func writerError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
