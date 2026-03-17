package postgre

import (
	"context"
	"errors"
	"log"
	"worker_pool/internal/handlers/models"

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
	rows, err := s.db.Query(ctx, `SELECT id, description FROM tasks`)
	if err != nil {
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
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *TaskStore) GetByID(ctx context.Context, id int) (models.Task, error) {
	var task models.Task
	err := s.db.QueryRow(ctx, `SELECT * FROM tasks WHERE id = $1 `, id).Scan(&task.ID, &task.Description)
	if errors.Is(err, pgx.ErrNoRows) {
		return task, nil
	}
	if err != nil {
		return task, err
	}
	return task, nil
}

func (s *TaskStore) Create(ctx context.Context, task models.Task) (models.Task, error) {
	err := s.db.QueryRow(ctx, `INSERT INTO tasks(description) VALUES ($1) RETURNING id, description, status`, task.Description).Scan(&task.ID, &task.Description, &task.Status)
	if err != nil {
		return task, err
	}
	return task, nil
}

func (s *TaskStore) GetPendingIDs(ctx context.Context) ([]int, error) {
	rows, err := s.db.Query(ctx, `SELECT id FROM tasks WHERE status = 'pending'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

func (s *TaskStore) UpdateTaskStatus(ctx context.Context, id int) (models.Task, error) {
	var task models.Task
	err := s.db.QueryRow(ctx, `UPDATE tasks SET status = 'done' WHERE id = $1 RETURNING id, description, status`, id).Scan(&task.ID, &task.Description, &task.Status)
	if err != nil {
		return task, err
	}
	return task, nil
}
