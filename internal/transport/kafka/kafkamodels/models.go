package kafkamodels

type CreateTaskMessage struct {
	ExternalID  string `json:"external_id"`
	Description string `json:"description"`
}
