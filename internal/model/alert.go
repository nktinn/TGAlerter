package model

type Alert struct {
	ServiceID string `json:"service_id"`
	AlertType int    `json:"alert_type"`
	Message   string `json:"message"`
}
