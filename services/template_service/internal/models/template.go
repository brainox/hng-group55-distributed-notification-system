package models

import (
	"time"

	"github.com/google/uuid"
)

type Template struct {
	ID           uuid.UUID `json:"id" db:"id"`
	TemplateKey  string    `json:"template_key" db:"template_key"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	TemplateType string    `json:"template_type" db:"template_type"` // email, push, sms
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type TemplateVersion struct {
	ID            uuid.UUID `json:"id" db:"id"`
	TemplateID    uuid.UUID `json:"template_id" db:"template_id"`
	VersionNumber int       `json:"version_number" db:"version_number"`
	Language      string    `json:"language" db:"language"`
	Subject       string    `json:"subject" db:"subject"`
	Body          string    `json:"body" db:"body"`
	Variables     []string  `json:"variables" db:"variables"`
	IsPublished   bool      `json:"is_published" db:"is_published"`
	CreatedBy     string    `json:"created_by" db:"created_by"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// Request/Response DTOs
type CreateTemplateRequest struct {
	TemplateKey  string   `json:"template_key" binding:"required"`
	Name         string   `json:"name" binding:"required"`
	Description  string   `json:"description"`
	TemplateType string   `json:"template_type" binding:"required,oneof=email push sms"`
	Subject      string   `json:"subject"`
	Body         string   `json:"body" binding:"required"`
	Language     string   `json:"language"`
	Variables    []string `json:"variables"`
}

type UpdateTemplateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

type CreateVersionRequest struct {
	Language  string   `json:"language" binding:"required"`
	Subject   string   `json:"subject"`
	Body      string   `json:"body" binding:"required"`
	Variables []string `json:"variables"`
}

type ValidateTemplateRequest struct {
	TemplateKey string                 `json:"template_key" binding:"required"`
	Variables   map[string]interface{} `json:"variables" binding:"required"`
}

type PreviewTemplateRequest struct {
	VersionID string                 `json:"version_id" binding:"required"`
	Variables map[string]interface{} `json:"variables" binding:"required"`
}

type TemplateResponse struct {
	Template *Template        `json:"template"`
	Version  *TemplateVersion `json:"version,omitempty"`
}

type ListTemplatesQuery struct {
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
	Type     string `form:"type"`
	Language string `form:"language"`
	Search   string `form:"search"`
}
