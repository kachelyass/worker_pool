package handlers

import (
	"encoding/json"
	"net/http"
	"worker_pool/internal/database"
	"worker_pool/internal/models"
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

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var in struct {
		Description string `json:"description"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&in); err != nil {
		writerError(w, http.StatusBadRequest, "Title is required")
		return
	}
	task := models.Task{
		Description: in.Description,
	}
	created, err := h.store.Create(r.Context(), task)
	if err != nil {
		writerError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)

}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writerError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
