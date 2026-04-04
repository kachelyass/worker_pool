package postgre

import (
	"context"
	"errors"
	"log"
	"time"
	"worker_pool/internal/handlers/models"
	"worker_pool/pkg/metrics"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskStore struct {
	db *pgxpool.Pool
}

func NewTaskStore(db *pgxpool.Pool) *TaskStore {
	return &TaskStore{db: db}
}

func (s *TaskStore) GetAll(ctx context.Context) ([]models.Task, error) {
	start := time.Now()
	operation := "get_all_tasks"
	defer func() {
		metrics.DBQueryTotal.WithLabelValues(operation).Inc()
		metrics.DBQueryDuration.WithLabelValues(operation).Observe(time.Since(start).Seconds())
	}()
	rows, err := s.db.Query(ctx, `SELECT id, description FROM tasks`)
	if err != nil {
		metrics.DBQueryErrorsTotal.WithLabelValues(operation).Inc()
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Description); err != nil {
			log.Printf("GetAll error: %s", err)
			return nil, err
		}
		if err := rows.Err(); err != nil {
			metrics.DBQueryErrorsTotal.WithLabelValues(operation).Inc()
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *TaskStore) GetByID(ctx context.Context, id int) (models.Task, error) {
	var task models.Task
	start := time.Now()
	operation := "get_task_by_id"
	defer func() {
		metrics.DBQueryTotal.WithLabelValues(operation).Inc()
		metrics.DBQueryDuration.WithLabelValues(operation).Observe(time.Since(start).Seconds())
	}()
	err := s.db.QueryRow(ctx, `SELECT * FROM tasks WHERE id = $1 `, id).Scan(&task.ID, &task.Description, &task.Status)
	if errors.Is(err, pgx.ErrNoRows) {
		metrics.DBQueryErrorsTotal.WithLabelValues(operation).Inc()
		return task, nil
	}
	if err != nil {
		metrics.DBQueryErrorsTotal.WithLabelValues(operation).Inc()
		return task, err
	}
	return task, nil
}

func (s *TaskStore) Create(ctx context.Context, task models.Task) (models.Task, error) {
	start := time.Now()
	operation := "create_task"
	defer func() {
		metrics.DBQueryTotal.WithLabelValues(operation).Inc()
		metrics.DBQueryDuration.WithLabelValues(operation).Observe(time.Since(start).Seconds())
	}()
	err := s.db.QueryRow(ctx, `INSERT INTO tasks(description) VALUES ($1) RETURNING id, description, status`, task.Description).Scan(&task.ID, &task.Description, &task.Status)
	if err != nil {
		metrics.DBQueryErrorsTotal.WithLabelValues(operation).Inc()
		return task, err
	}
	return task, nil
}

func (s *TaskStore) ClaimNextIDs(ctx context.Context, limit int) ([]int, error) {
	start := time.Now()
	operation := "claim_next_ids"
	defer func() {
		metrics.DBQueryTotal.WithLabelValues(operation).Inc()
		metrics.DBQueryDuration.WithLabelValues(operation).Observe(time.Since(start).Seconds())
	}()
	rows, err := s.db.Query(ctx, `WITH picked AS (
SELECT id 
FROM tasks 
WHERE status = 'pending'
LIMIT $1
FOR UPDATE SKIP LOCKED
)
UPDATE tasks 
SET status = 'processing'
FROM picked 
WHERE tasks.id = picked.id
RETURNING tasks.id`, limit)
	if err != nil {
		metrics.DBQueryErrorsTotal.WithLabelValues(operation).Inc()
		return nil, err
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			metrics.DBQueryErrorsTotal.WithLabelValues(operation).Inc()
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		metrics.DBQueryErrorsTotal.WithLabelValues(operation).Inc()
		return nil, err
	}
	return ids, nil
}

func (s *TaskStore) MarkDone(ctx context.Context, id int) (models.Task, error) {
	start := time.Now()
	operation := "mark_done_task"
	defer func() {
		metrics.DBQueryTotal.WithLabelValues(operation).Inc()
		metrics.DBQueryDuration.WithLabelValues(operation).Observe(time.Since(start).Seconds())
	}()
	var task models.Task
	err := s.db.QueryRow(ctx, `UPDATE tasks SET status = 'done' WHERE id = $1 RETURNING id, description, status`, id).Scan(&task.ID, &task.Description, &task.Status)
	if err != nil {
		metrics.DBQueryErrorsTotal.WithLabelValues(operation).Inc()
		return task, err
	}
	return task, nil
}
