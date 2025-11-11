package repository

import (
	"context"
	"fmt"

	"github.com/brainox/hng-group55-distributed-notification-system/services/template_service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TemplateRepository interface {
	Create(ctx context.Context, template *models.Template) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Template, error)
	GetByKey(ctx context.Context, key string) (*models.Template, error)
	List(ctx context.Context, query models.ListTemplatesQuery) ([]*models.Template, int, error)
	Update(ctx context.Context, id uuid.UUID, req models.UpdateTemplateRequest) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

type templateRepository struct {
	db *pgxpool.Pool
}

func NewTemplateRepository(db *pgxpool.Pool) TemplateRepository {
	return &templateRepository{db: db}
}

func (r *templateRepository) Create(ctx context.Context, template *models.Template) error {
	query := `
		INSERT INTO templates (template_key, name, description, template_type, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(ctx, query,
		template.TemplateKey,
		template.Name,
		template.Description,
		template.TemplateType,
		template.IsActive,
	).Scan(&template.ID, &template.CreatedAt, &template.UpdatedAt)
}

func (r *templateRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Template, error) {
	query := `
		SELECT id, template_key, name, description, template_type, is_active, created_at, updated_at
		FROM templates
		WHERE id = $1 AND is_active = true
	`

	var template models.Template
	err := r.db.QueryRow(ctx, query, id).Scan(
		&template.ID,
		&template.TemplateKey,
		&template.Name,
		&template.Description,
		&template.TemplateType,
		&template.IsActive,
		&template.CreatedAt,
		&template.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("template not found")
	}
	if err != nil {
		return nil, err
	}

	return &template, nil
}

func (r *templateRepository) GetByKey(ctx context.Context, key string) (*models.Template, error) {
	query := `
		SELECT id, template_key, name, description, template_type, is_active, created_at, updated_at
		FROM templates
		WHERE template_key = $1 AND is_active = true
	`

	var template models.Template
	err := r.db.QueryRow(ctx, query, key).Scan(
		&template.ID,
		&template.TemplateKey,
		&template.Name,
		&template.Description,
		&template.TemplateType,
		&template.IsActive,
		&template.CreatedAt,
		&template.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("template not found")
	}
	if err != nil {
		return nil, err
	}

	return &template, nil
}

func (r *templateRepository) List(ctx context.Context, query models.ListTemplatesQuery) ([]*models.Template, int, error) {
	// Set defaults
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 20
	}

	offset := (query.Page - 1) * query.Limit

	// Build WHERE clause
	whereClause := "WHERE is_active = true"
	args := []interface{}{}
	argCount := 1

	if query.Type != "" {
		whereClause += fmt.Sprintf(" AND template_type = $%d", argCount)
		args = append(args, query.Type)
		argCount++
	}

	if query.Search != "" {
		whereClause += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+query.Search+"%")
		argCount++
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM templates %s", whereClause)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Fetch templates
	listQuery := fmt.Sprintf(`
		SELECT id, template_key, name, description, template_type, is_active, created_at, updated_at
		FROM templates
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argCount, argCount+1)

	args = append(args, query.Limit, offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	templates := []*models.Template{}
	for rows.Next() {
		var template models.Template
		err := rows.Scan(
			&template.ID,
			&template.TemplateKey,
			&template.Name,
			&template.Description,
			&template.TemplateType,
			&template.IsActive,
			&template.CreatedAt,
			&template.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		templates = append(templates, &template)
	}

	return templates, total, nil
}

func (r *templateRepository) Update(ctx context.Context, id uuid.UUID, req models.UpdateTemplateRequest) error {
	updates := []string{}
	args := []interface{}{}
	argCount := 1

	if req.Name != "" {
		updates = append(updates, fmt.Sprintf("name = $%d", argCount))
		args = append(args, req.Name)
		argCount++
	}

	if req.Description != "" {
		updates = append(updates, fmt.Sprintf("description = $%d", argCount))
		args = append(args, req.Description)
		argCount++
	}

	if req.IsActive != nil {
		updates = append(updates, fmt.Sprintf("is_active = $%d", argCount))
		args = append(args, *req.IsActive)
		argCount++
	}

	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE templates SET %s WHERE id = $%d", joinStrings(updates, ", "), argCount)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("template not found")
	}

	return nil
}

func (r *templateRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	query := "UPDATE templates SET is_active = false WHERE id = $1"
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("template not found")
	}

	return nil
}

func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
