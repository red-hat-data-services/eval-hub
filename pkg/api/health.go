package api

import "time"

// HealthResponse represents health check response
type HealthResponse struct {
	Status            string                    `json:"status"`
	Version           string                    `json:"version"`
	Timestamp         *time.Time                `json:"timestamp"`
	Components        map[string]map[string]any `json:"components,omitempty"`
	Uptime            time.Duration             `json:"uptime"`
	ActiveEvaluations int                       `json:"active_evaluations,omitempty"`
}
