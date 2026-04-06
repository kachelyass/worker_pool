package models

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type SetWorkersRequest struct {
	Count int `json:"count"`
}

type CreateTaskMessage struct {
	ExternalID  string `json:"external_id"`
	Description string `json:"description"`
}
