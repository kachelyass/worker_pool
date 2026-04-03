package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"worker_pool/internal/handlers/models"
)

type TaskStore interface {
	GetAll(ctx context.Context) ([]models.Task, error)
	GetByID(ctx context.Context, id int) (models.Task, error)
	Create(ctx context.Context, task models.Task) (models.Task, error)
}

type TaskReader interface {
	GetAll(ctx context.Context) ([]models.Task, error)
	GetByID(ctx context.Context, id int) (models.Task, error)
}

type TaskPublisher interface {
	PublishCreateTask(description string) error
}

type TaskHandler struct {
	store     TaskStore
	publisher TaskPublisher
}

func NewTaskHandler(store TaskStore, publisher TaskPublisher) *TaskHandler {
	return &TaskHandler{store: store, publisher: publisher}
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
		return
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
		writerError(w, http.StatusBadRequest, "description is required")
		return
	}

	if in.Description == "" {
		writerError(w, http.StatusBadRequest, "description is required")
		return
	}

	if err := h.publisher.PublishCreateTask(in.Description); err != nil {
		writerError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]string{
		"status": "queued",
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writerError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
