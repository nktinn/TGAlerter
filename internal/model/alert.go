package model

type Alert struct {
	ServiceId string `json:"service_id"`
	AlertType int    `json:"alert_type"`
	Message   string `json:"message"`
}
