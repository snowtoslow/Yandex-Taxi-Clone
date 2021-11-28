package models

type HealthRequest struct{}

type HealthResponse struct {
	Status string `json:"status"`
}
