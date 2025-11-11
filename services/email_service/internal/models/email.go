package models

import "time"

// EmailMessage represents a message from the queue
type EmailMessage struct {
	NotificationID   string                 `json:"notification_id"`
	NotificationType string                 `json:"notification_type"`
	UserID           string                 `json:"user_id"`
	Recipient        string                 `json:"recipient"`
	Subject          string                 `json:"subject"`
	Body             string                 `json:"body"`
	TemplateCode     string                 `json:"template_code"`
	Variables        map[string]interface{} `json:"variables"`
	Priority         int                    `json:"priority"`
	Metadata         struct {
		Timestamp  string `json:"timestamp"`
		RetryCount int    `json:"retry_count"`
	} `json:"metadata"`
}

// EmailTemplate represents template data from Template Service
type EmailTemplate struct {
	Subject   string   `json:"subject"`
	Body      string   `json:"body"`
	Variables []string `json:"variables"`
}

// StatusMessage represents a status update message
type StatusMessage struct {
	NotificationID string    `json:"notification_id"`
	UserID         string    `json:"user_id"`
	Status         string    `json:"status"` // "sent", "delivered", "failed"
	Timestamp      time.Time `json:"timestamp"`
	Error          string    `json:"error,omitempty"`
	Provider       string    `json:"provider"`
}

// TemplateResponse represents the response from Template Service
type TemplateResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Template struct {
			ID           string    `json:"id"`
			TemplateKey  string    `json:"template_key"`
			Name         string    `json:"name"`
			TemplateType string    `json:"template_type"`
			IsActive     bool      `json:"is_active"`
			CreatedAt    time.Time `json:"created_at"`
		} `json:"template"`
		Version struct {
			ID            string    `json:"id"`
			Subject       string    `json:"subject"`
			Body          string    `json:"body"`
			Variables     []string  `json:"variables"`
			Language      string    `json:"language"`
			VersionNumber int       `json:"version_number"`
			IsPublished   bool      `json:"is_published"`
			CreatedAt     time.Time `json:"created_at"`
		} `json:"version"`
	} `json:"data"`
	Message string `json:"message"`
}
