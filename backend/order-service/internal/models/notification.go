package models

// EmailNotification — структура для email-уведомления
type EmailNotification struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Type    string `json:"type"` // order_created, order_status_changed, order_delivered
}