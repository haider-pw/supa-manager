-- Add encrypted database password storage for connection strings
ALTER TABLE project ADD COLUMN IF NOT EXISTS db_password_encrypted TEXT;

-- Index for faster lookups
CREATE INDEX IF NOT EXISTS idx_project_ref ON project(project_ref);
