package handlers

import (
	"encoding/json"
	"net/http"
	"worker_pool/internal/app/workerpool"
	"worker_pool/internal/handlers/models"
)

type Server struct {
	pool *workerpool.PoolManager
}

func NewServer(server *workerpool.PoolManager) *Server {
	return &Server{pool: server}
}

func (s *Server) GetWorkers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":  "ok",
		"workers": s.pool.Count(),
	})
}

func (s *Server) SetWorkers(w http.ResponseWriter, r *http.Request) {
	var req models.SetWorkersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	if req.Count < 1 || req.Count > 1000 {
		http.Error(w, "count must be between 1 and 1000", http.StatusBadRequest)
		return
	}

	current := s.pool.SetWorkers(req.Count)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"workers": current,
	})
}
