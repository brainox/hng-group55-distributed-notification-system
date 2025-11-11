package handler

import (
	"net/http"

	"github.com/brainox/hng-group55-distributed-notification-system/services/template_service/internal/models"
	"github.com/brainox/hng-group55-distributed-notification-system/services/template_service/internal/service"
	"github.com/brainox/hng-group55-distributed-notification-system/services/template_service/pkg/logger"
	"github.com/brainox/hng-group55-distributed-notification-system/services/template_service/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TemplateHandler struct {
	service service.TemplateService
}

func NewTemplateHandler(service service.TemplateService) *TemplateHandler {
	return &TemplateHandler{service: service}
}

// CreateTemplate handles POST /api/v1/templates
func (h *TemplateHandler) CreateTemplate(c *gin.Context) {
	var req models.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("failed to bind request", zap.Error(err))
		response.ErrorMessage(c, http.StatusBadRequest, err.Error(), "Invalid request body")
		return
	}

	result, err := h.service.CreateTemplate(c.Request.Context(), req)
	if err != nil {
		logger.Log.Error("failed to create template", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, err, "Failed to create template")
		return
	}

	response.Success(c, http.StatusCreated, result, "Template created successfully")
}

// GetTemplateByID handles GET /api/v1/templates/:id
func (h *TemplateHandler) GetTemplateByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "Invalid template ID", "Invalid template ID format")
		return
	}

	language := c.DefaultQuery("language", "en")
	version := c.DefaultQuery("version", "latest")

	result, err := h.service.GetTemplateByID(c.Request.Context(), id, language, version)
	if err != nil {
		logger.Log.Error("failed to get template", zap.Error(err), zap.String("id", idParam))
		response.Error(c, http.StatusNotFound, err, "Template not found")
		return
	}

	response.Success(c, http.StatusOK, result, "Template retrieved successfully")
}

// GetTemplateByKey handles GET /api/v1/templates/key/:key
func (h *TemplateHandler) GetTemplateByKey(c *gin.Context) {
	key := c.Param("key")
	language := c.DefaultQuery("language", "en")
	version := c.DefaultQuery("version", "latest")

	result, err := h.service.GetTemplateByKey(c.Request.Context(), key, language, version)
	if err != nil {
		logger.Log.Error("failed to get template by key", zap.Error(err), zap.String("key", key))
		response.Error(c, http.StatusNotFound, err, "Template not found")
		return
	}

	response.Success(c, http.StatusOK, result, "Template retrieved successfully")
}

// ListTemplates handles GET /api/v1/templates
func (h *TemplateHandler) ListTemplates(c *gin.Context) {
	var query models.ListTemplatesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, err.Error(), "Invalid query parameters")
		return
	}

	templates, total, err := h.service.ListTemplates(c.Request.Context(), query)
	if err != nil {
		logger.Log.Error("failed to list templates", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, err, "Failed to list templates")
		return
	}

	// Set defaults for pagination
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 20
	}

	meta := response.CalculateMeta(total, query.Limit, query.Page)
	response.SuccessWithMeta(c, http.StatusOK, templates, "Templates retrieved successfully", meta)
}

// UpdateTemplate handles PUT /api/v1/templates/:id
func (h *TemplateHandler) UpdateTemplate(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "Invalid template ID", "Invalid template ID format")
		return
	}

	var req models.UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, err.Error(), "Invalid request body")
		return
	}

	if err := h.service.UpdateTemplate(c.Request.Context(), id, req); err != nil {
		logger.Log.Error("failed to update template", zap.Error(err), zap.String("id", idParam))
		response.Error(c, http.StatusInternalServerError, err, "Failed to update template")
		return
	}

	response.Success(c, http.StatusOK, nil, "Template updated successfully")
}

// DeleteTemplate handles DELETE /api/v1/templates/:id
func (h *TemplateHandler) DeleteTemplate(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "Invalid template ID", "Invalid template ID format")
		return
	}

	if err := h.service.DeleteTemplate(c.Request.Context(), id); err != nil {
		logger.Log.Error("failed to delete template", zap.Error(err), zap.String("id", idParam))
		response.Error(c, http.StatusInternalServerError, err, "Failed to delete template")
		return
	}

	response.Success(c, http.StatusOK, nil, "Template deleted successfully")
}

// CreateVersion handles POST /api/v1/templates/:id/versions
func (h *TemplateHandler) CreateVersion(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "Invalid template ID", "Invalid template ID format")
		return
	}

	var req models.CreateVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, err.Error(), "Invalid request body")
		return
	}

	version, err := h.service.CreateVersion(c.Request.Context(), id, req)
	if err != nil {
		logger.Log.Error("failed to create version", zap.Error(err), zap.String("template_id", idParam))
		response.Error(c, http.StatusInternalServerError, err, "Failed to create version")
		return
	}

	response.Success(c, http.StatusCreated, version, "Version created successfully")
}

// PublishVersion handles POST /api/v1/templates/:id/versions/:version_id/publish
func (h *TemplateHandler) PublishVersion(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "Invalid template ID", "Invalid template ID format")
		return
	}

	versionIDParam := c.Param("version_id")
	versionID, err := uuid.Parse(versionIDParam)
	if err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "Invalid version ID", "Invalid version ID format")
		return
	}

	if err := h.service.PublishVersion(c.Request.Context(), id, versionID); err != nil {
		logger.Log.Error("failed to publish version", zap.Error(err), zap.String("version_id", versionIDParam))
		response.Error(c, http.StatusInternalServerError, err, "Failed to publish version")
		return
	}

	response.Success(c, http.StatusOK, nil, "Version published successfully")
}

// GetVersionHistory handles GET /api/v1/templates/:id/versions
func (h *TemplateHandler) GetVersionHistory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "Invalid template ID", "Invalid template ID format")
		return
	}

	versions, err := h.service.GetVersionHistory(c.Request.Context(), id)
	if err != nil {
		logger.Log.Error("failed to get version history", zap.Error(err), zap.String("template_id", idParam))
		response.Error(c, http.StatusInternalServerError, err, "Failed to get version history")
		return
	}

	response.Success(c, http.StatusOK, versions, "Version history retrieved successfully")
}

// ValidateTemplate handles POST /api/v1/templates/validate
func (h *TemplateHandler) ValidateTemplate(c *gin.Context) {
	var req models.ValidateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, err.Error(), "Invalid request body")
		return
	}

	valid, missing, err := h.service.ValidateTemplate(c.Request.Context(), req)
	if err != nil {
		logger.Log.Error("failed to validate template", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, err, "Failed to validate template")
		return
	}

	result := gin.H{
		"valid":             valid,
		"missing_variables": missing,
	}

	response.Success(c, http.StatusOK, result, "Template validation completed")
}

// PreviewTemplate handles POST /api/v1/templates/:id/preview
func (h *TemplateHandler) PreviewTemplate(c *gin.Context) {
	idParam := c.Param("id")
	_, err := uuid.Parse(idParam)
	if err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "Invalid template ID", "Invalid template ID format")
		return
	}

	var req models.PreviewTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, err.Error(), "Invalid request body")
		return
	}

	subject, body, err := h.service.PreviewTemplate(c.Request.Context(), req)
	if err != nil {
		logger.Log.Error("failed to preview template", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, err, "Failed to preview template")
		return
	}

	result := gin.H{
		"subject": subject,
		"body":    body,
	}

	response.Success(c, http.StatusOK, result, "Template preview generated successfully")
}
