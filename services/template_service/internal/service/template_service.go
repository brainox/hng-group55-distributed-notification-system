package service

import (
	"context"
	"fmt"
	"time"

	"github.com/brainox/hng-group55-distributed-notification-system/services/template_service/internal/models"
	"github.com/brainox/hng-group55-distributed-notification-system/services/template_service/internal/renderer"
	"github.com/brainox/hng-group55-distributed-notification-system/services/template_service/internal/repository"
	"github.com/google/uuid"
)

type TemplateService interface {
	CreateTemplate(ctx context.Context, req models.CreateTemplateRequest) (*models.TemplateResponse, error)
	GetTemplateByID(ctx context.Context, id uuid.UUID, language string, version string) (*models.TemplateResponse, error)
	GetTemplateByKey(ctx context.Context, key string, language string, version string) (*models.TemplateResponse, error)
	ListTemplates(ctx context.Context, query models.ListTemplatesQuery) ([]*models.Template, int, error)
	UpdateTemplate(ctx context.Context, id uuid.UUID, req models.UpdateTemplateRequest) error
	DeleteTemplate(ctx context.Context, id uuid.UUID) error
	CreateVersion(ctx context.Context, templateID uuid.UUID, req models.CreateVersionRequest) (*models.TemplateVersion, error)
	PublishVersion(ctx context.Context, templateID uuid.UUID, versionID uuid.UUID) error
	GetVersionHistory(ctx context.Context, templateID uuid.UUID) ([]*models.TemplateVersion, error)
	ValidateTemplate(ctx context.Context, req models.ValidateTemplateRequest) (bool, []string, error)
	PreviewTemplate(ctx context.Context, req models.PreviewTemplateRequest) (string, string, error)
}

type templateService struct {
	templateRepo repository.TemplateRepository
	versionRepo  repository.VersionRepository
	cache        CacheService
	cacheTTL     time.Duration
}

func NewTemplateService(
	templateRepo repository.TemplateRepository,
	versionRepo repository.VersionRepository,
	cache CacheService,
	cacheTTL time.Duration,
) TemplateService {
	return &templateService{
		templateRepo: templateRepo,
		versionRepo:  versionRepo,
		cache:        cache,
		cacheTTL:     cacheTTL,
	}
}

func (templateService *templateService) CreateTemplate(ctx context.Context, req models.CreateTemplateRequest) (*models.TemplateResponse, error) {
	// Set defaults
	if req.Language == "" {
		req.Language = "en"
	}

	// Extract variables from body
	if req.Variables == nil || len(req.Variables) == 0 {
		req.Variables = renderer.ExtractVariables(req.Body)
		if req.Subject != "" {
			subjectVars := renderer.ExtractVariables(req.Subject)
			// Merge unique variables
			varMap := make(map[string]bool)
			for _, v := range req.Variables {
				varMap[v] = true
			}
			for _, v := range subjectVars {
				if !varMap[v] {
					req.Variables = append(req.Variables, v)
				}
			}
		}
	}

	// Create template
	template := &models.Template{
		TemplateKey:  req.TemplateKey,
		Name:         req.Name,
		Description:  req.Description,
		TemplateType: req.TemplateType,
		IsActive:     true,
	}

	if err := templateService.templateRepo.Create(ctx, template); err != nil {
		return nil, fmt.Errorf("failed to create template: %w", err)
	}

	// Create first version
	version := &models.TemplateVersion{
		TemplateID:    template.ID,
		VersionNumber: 1,
		Language:      req.Language,
		Subject:       req.Subject,
		Body:          req.Body,
		Variables:     req.Variables,
		IsPublished:   true,
		CreatedBy:     "system",
	}

	if err := templateService.versionRepo.Create(ctx, version); err != nil {
		return nil, fmt.Errorf("failed to create version: %w", err)
	}

	return &models.TemplateResponse{
		Template: template,
		Version:  version,
	}, nil
}

func (templateService *templateService) GetTemplateByID(ctx context.Context, id uuid.UUID, language string, version string) (*models.TemplateResponse, error) {
	if language == "" {
		language = "en"
	}

	// Get template
	template, err := templateService.templateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get version
	var templateVersion *models.TemplateVersion
	if version == "latest" || version == "" {
		templateVersion, err = templateService.versionRepo.GetPublished(ctx, id, language)
	} else {
		// For specific version, would need version number parsing
		templateVersion, err = templateService.versionRepo.GetPublished(ctx, id, language)
	}

	if err != nil {
		return nil, err
	}

	return &models.TemplateResponse{
		Template: template,
		Version:  templateVersion,
	}, nil
}

func (templateService *templateService) GetTemplateByKey(ctx context.Context, key string, language string, version string) (*models.TemplateResponse, error) {
	if language == "" {
		language = "en"
	}
	if version == "" {
		version = "latest"
	}

	// Check cache first
	cached, err := templateService.cache.GetTemplate(ctx, key, language, version)
	if err == nil && cached != nil {
		return cached, nil
	}

	// Get template
	template, err := templateService.templateRepo.GetByKey(ctx, key)
	if err != nil {
		return nil, err
	}

	// Get published version
	templateVersion, err := templateService.versionRepo.GetPublished(ctx, template.ID, language)
	if err != nil {
		return nil, err
	}

	response := &models.TemplateResponse{
		Template: template,
		Version:  templateVersion,
	}

	// Cache the result
	_ = templateService.cache.SetTemplate(ctx, key, language, version, response, templateService.cacheTTL)

	return response, nil
}

func (s *templateService) ListTemplates(ctx context.Context, query models.ListTemplatesQuery) ([]*models.Template, int, error) {
	return s.templateRepo.List(ctx, query)
}

func (s *templateService) UpdateTemplate(ctx context.Context, id uuid.UUID, req models.UpdateTemplateRequest) error {
	if err := s.templateRepo.Update(ctx, id, req); err != nil {
		return err
	}

	// Invalidate cache
	template, _ := s.templateRepo.GetByID(ctx, id)
	if template != nil {
		_ = s.cache.InvalidateTemplate(ctx, template.TemplateKey)
	}

	return nil
}

func (s *templateService) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	// Get template key for cache invalidation
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.templateRepo.SoftDelete(ctx, id); err != nil {
		return err
	}

	// Invalidate cache
	_ = s.cache.InvalidateTemplate(ctx, template.TemplateKey)

	return nil
}

func (s *templateService) CreateVersion(ctx context.Context, templateID uuid.UUID, req models.CreateVersionRequest) (*models.TemplateVersion, error) {
	if req.Language == "" {
		req.Language = "en"
	}

	// Extract variables if not provided
	if req.Variables == nil || len(req.Variables) == 0 {
		req.Variables = renderer.ExtractVariables(req.Body)
		if req.Subject != "" {
			subjectVars := renderer.ExtractVariables(req.Subject)
			varMap := make(map[string]bool)
			for _, v := range req.Variables {
				varMap[v] = true
			}
			for _, v := range subjectVars {
				if !varMap[v] {
					req.Variables = append(req.Variables, v)
				}
			}
		}
	}

	// Get next version number
	nextVersion, err := s.versionRepo.GetNextVersionNumber(ctx, templateID, req.Language)
	if err != nil {
		return nil, err
	}

	version := &models.TemplateVersion{
		TemplateID:    templateID,
		VersionNumber: nextVersion,
		Language:      req.Language,
		Subject:       req.Subject,
		Body:          req.Body,
		Variables:     req.Variables,
		IsPublished:   false,
		CreatedBy:     "system",
	}

	if err := s.versionRepo.Create(ctx, version); err != nil {
		return nil, err
	}

	return version, nil
}

func (s *templateService) PublishVersion(ctx context.Context, templateID uuid.UUID, versionID uuid.UUID) error {
	// Get version to determine language
	version, err := s.versionRepo.GetByID(ctx, versionID)
	if err != nil {
		return err
	}

	// Unpublish other versions for this language
	if err := s.versionRepo.UnpublishOthers(ctx, templateID, version.Language, versionID); err != nil {
		return err
	}

	// Publish this version
	if err := s.versionRepo.Publish(ctx, versionID); err != nil {
		return err
	}

	// Invalidate cache
	template, _ := s.templateRepo.GetByID(ctx, templateID)
	if template != nil {
		_ = s.cache.InvalidateTemplate(ctx, template.TemplateKey)
	}

	return nil
}

func (s *templateService) GetVersionHistory(ctx context.Context, templateID uuid.UUID) ([]*models.TemplateVersion, error) {
	return s.versionRepo.ListByTemplateID(ctx, templateID)
}

func (s *templateService) ValidateTemplate(ctx context.Context, req models.ValidateTemplateRequest) (bool, []string, error) {
	// Get template
	template, err := s.templateRepo.GetByKey(ctx, req.TemplateKey)
	if err != nil {
		return false, nil, err
	}

	// Get published version
	version, err := s.versionRepo.GetPublished(ctx, template.ID, "en")
	if err != nil {
		return false, nil, err
	}

	// Check missing variables
	missing := renderer.ValidateVariables(version.Variables, req.Variables)

	return len(missing) == 0, missing, nil
}

func (s *templateService) PreviewTemplate(ctx context.Context, req models.PreviewTemplateRequest) (string, string, error) {
	versionID, err := uuid.Parse(req.VersionID)
	if err != nil {
		return "", "", fmt.Errorf("invalid version ID")
	}

	// Get version
	version, err := s.versionRepo.GetByID(ctx, versionID)
	if err != nil {
		return "", "", err
	}

	// Render subject
	renderedSubject := version.Subject
	if version.Subject != "" {
		renderedSubject, err = renderer.RenderTemplate(version.Subject, req.Variables)
		if err != nil {
			return "", "", fmt.Errorf("failed to render subject: %w", err)
		}
	}

	// Render body
	renderedBody, err := renderer.RenderTemplate(version.Body, req.Variables)
	if err != nil {
		return "", "", fmt.Errorf("failed to render body: %w", err)
	}

	return renderedSubject, renderedBody, nil
}
