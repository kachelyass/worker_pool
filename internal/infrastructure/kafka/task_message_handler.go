package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"worker_pool/internal/handlers/models"
)

type MessageStore interface {
	Create(ctx context.Context, task models.Task) (models.Task, error)
}

type MessageCreateHandler struct {
	store MessageStore
}

func NewMessageCreateHandler(store MessageStore) *MessageCreateHandler {
	return &MessageCreateHandler{store: store}
}

func (h *MessageCreateHandler) Handle(ctx context.Context, key []byte, value []byte) error {
	var msg models.CreateTaskMessage

	if err := json.Unmarshal(value, &msg); err != nil {
		return fmt.Errorf("unmarshal create task message: %w", err)
	}

	_, err := h.store.Create(ctx, models.Task{
		Description: msg.Description,
		Status:      "pending",
	})
	if err != nil {
		return fmt.Errorf("create task in db: %w", err)
	}

	return nil
}
