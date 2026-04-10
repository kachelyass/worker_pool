package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"worker_pool/internal/transport/kafka/kafkamodels"
	"worker_pool/internal/transport/rest/httpmodels"
)

type MessageStore interface {
	Create(ctx context.Context, task httpmodels.Task) (httpmodels.Task, error)
}

type MessageCreateHandler struct {
	store MessageStore
}

func NewMessageCreateHandler(store MessageStore) *MessageCreateHandler {
	return &MessageCreateHandler{store: store}
}

func (h *MessageCreateHandler) Handle(ctx context.Context, key []byte, value []byte) error {
	var msg kafkamodels.CreateTaskMessage

	if err := json.Unmarshal(value, &msg); err != nil {
		return fmt.Errorf("unmarshal create task message: %w", err)
	}

	_, err := h.store.Create(ctx, httpmodels.Task{
		Description: msg.Description,
		Status:      "pending",
	})
	if err != nil {
		return fmt.Errorf("create task in db: %w", err)
	}

	return nil
}
