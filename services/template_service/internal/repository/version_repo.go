package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/brainox/hng-group55-distributed-notification-system/services/template_service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VersionRepository interface {
	Create(ctx context.Context, version *models.TemplateVersion) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.TemplateVersion, error)
	GetPublished(ctx context.Context, templateID uuid.UUID, language string) (*models.TemplateVersion, error)
	GetLatestVersion(ctx context.Context, templateID uuid.UUID, language string) (*models.TemplateVersion, error)
	ListByTemplateID(ctx context.Context, templateID uuid.UUID) ([]*models.TemplateVersion, error)
	Publish(ctx context.Context, versionID uuid.UUID) error
	UnpublishOthers(ctx context.Context, templateID uuid.UUID, language string, exceptVersionID uuid.UUID) error
	GetNextVersionNumber(ctx context.Context, templateID uuid.UUID, language string) (int, error)
}

type versionRepository struct {
	db *pgxpool.Pool
}

func NewVersionRepository(db *pgxpool.Pool) VersionRepository {
	return &versionRepository{db: db}
}

func (r *versionRepository) Create(ctx context.Context, version *models.TemplateVersion) error {
	variablesJSON, err := json.Marshal(version.Variables)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO template_versions (template_id, version_number, language, subject, body, variables, is_published, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`

	return r.db.QueryRow(ctx, query,
		version.TemplateID,
		version.VersionNumber,
		version.Language,
		version.Subject,
		version.Body,
		variablesJSON,
		version.IsPublished,
		version.CreatedBy,
	).Scan(&version.ID, &version.CreatedAt)
}

func (r *versionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.TemplateVersion, error) {
	query := `
		SELECT id, template_id, version_number, language, subject, body, variables, is_published, created_by, created_at
		FROM template_versions
		WHERE id = $1
	`

	var version models.TemplateVersion
	var variablesJSON []byte

	err := r.db.QueryRow(ctx, query, id).Scan(
		&version.ID,
		&version.TemplateID,
		&version.VersionNumber,
		&version.Language,
		&version.Subject,
		&version.Body,
		&variablesJSON,
		&version.IsPublished,
		&version.CreatedBy,
		&version.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("version not found")
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(variablesJSON, &version.Variables); err != nil {
		return nil, err
	}

	return &version, nil
}

func (r *versionRepository) GetPublished(ctx context.Context, templateID uuid.UUID, language string) (*models.TemplateVersion, error) {
	query := `
		SELECT id, template_id, version_number, language, subject, body, variables, is_published, created_by, created_at
		FROM template_versions
		WHERE template_id = $1 AND language = $2 AND is_published = true
		ORDER BY version_number DESC
		LIMIT 1
	`

	var version models.TemplateVersion
	var variablesJSON []byte

	err := r.db.QueryRow(ctx, query, templateID, language).Scan(
		&version.ID,
		&version.TemplateID,
		&version.VersionNumber,
		&version.Language,
		&version.Subject,
		&version.Body,
		&variablesJSON,
		&version.IsPublished,
		&version.CreatedBy,
		&version.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		// Try default language
		if language != "en" {
			return r.GetPublished(ctx, templateID, "en")
		}
		return nil, fmt.Errorf("published version not found")
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(variablesJSON, &version.Variables); err != nil {
		return nil, err
	}

	return &version, nil
}

func (r *versionRepository) GetLatestVersion(ctx context.Context, templateID uuid.UUID, language string) (*models.TemplateVersion, error) {
	query := `
		SELECT id, template_id, version_number, language, subject, body, variables, is_published, created_by, created_at
		FROM template_versions
		WHERE template_id = $1 AND language = $2
		ORDER BY version_number DESC
		LIMIT 1
	`

	var version models.TemplateVersion
	var variablesJSON []byte

	err := r.db.QueryRow(ctx, query, templateID, language).Scan(
		&version.ID,
		&version.TemplateID,
		&version.VersionNumber,
		&version.Language,
		&version.Subject,
		&version.Body,
		&variablesJSON,
		&version.IsPublished,
		&version.CreatedBy,
		&version.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("version not found")
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(variablesJSON, &version.Variables); err != nil {
		return nil, err
	}

	return &version, nil
}

func (r *versionRepository) ListByTemplateID(ctx context.Context, templateID uuid.UUID) ([]*models.TemplateVersion, error) {
	query := `
		SELECT id, template_id, version_number, language, subject, body, variables, is_published, created_by, created_at
		FROM template_versions
		WHERE template_id = $1
		ORDER BY version_number DESC, language ASC
	`

	rows, err := r.db.Query(ctx, query, templateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	versions := []*models.TemplateVersion{}
	for rows.Next() {
		var version models.TemplateVersion
		var variablesJSON []byte

		err := rows.Scan(
			&version.ID,
			&version.TemplateID,
			&version.VersionNumber,
			&version.Language,
			&version.Subject,
			&version.Body,
			&variablesJSON,
			&version.IsPublished,
			&version.CreatedBy,
			&version.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(variablesJSON, &version.Variables); err != nil {
			return nil, err
		}

		versions = append(versions, &version)
	}

	return versions, nil
}

func (r *versionRepository) Publish(ctx context.Context, versionID uuid.UUID) error {
	query := "UPDATE template_versions SET is_published = true WHERE id = $1"
	result, err := r.db.Exec(ctx, query, versionID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("version not found")
	}

	return nil
}

func (r *versionRepository) UnpublishOthers(ctx context.Context, templateID uuid.UUID, language string, exceptVersionID uuid.UUID) error {
	query := `
		UPDATE template_versions
		SET is_published = false
		WHERE template_id = $1 AND language = $2 AND id != $3
	`

	_, err := r.db.Exec(ctx, query, templateID, language, exceptVersionID)
	return err
}

func (r *versionRepository) GetNextVersionNumber(ctx context.Context, templateID uuid.UUID, language string) (int, error) {
	query := `
		SELECT COALESCE(MAX(version_number), 0) + 1
		FROM template_versions
		WHERE template_id = $1 AND language = $2
	`

	var nextVersion int
	err := r.db.QueryRow(ctx, query, templateID, language).Scan(&nextVersion)
	return nextVersion, err
}
