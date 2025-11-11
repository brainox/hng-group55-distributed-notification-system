-- Drop triggers
DROP TRIGGER IF EXISTS update_templates_updated_at ON templates;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_template_versions_language;
DROP INDEX IF EXISTS idx_template_versions_published;
DROP INDEX IF EXISTS idx_template_versions_template_id;
DROP INDEX IF EXISTS idx_templates_active;
DROP INDEX IF EXISTS idx_templates_type;
DROP INDEX IF EXISTS idx_templates_template_key;

-- Drop tables
DROP TABLE IF EXISTS template_versions;
DROP TABLE IF EXISTS templates;
