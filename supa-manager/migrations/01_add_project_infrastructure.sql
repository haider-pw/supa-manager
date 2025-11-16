-- Migration to add infrastructure-related columns to project table
ALTER TABLE project ADD COLUMN IF NOT EXISTS docker_compose_path TEXT;
ALTER TABLE project ADD COLUMN IF NOT EXISTS docker_network_name TEXT;
ALTER TABLE project ADD COLUMN IF NOT EXISTS postgres_port INTEGER;
ALTER TABLE project ADD COLUMN IF NOT EXISTS kong_http_port INTEGER;
ALTER TABLE project ADD COLUMN IF NOT EXISTS kong_https_port INTEGER;
ALTER TABLE project ADD COLUMN IF NOT EXISTS anon_key TEXT;
ALTER TABLE project ADD COLUMN IF NOT EXISTS service_role_key TEXT;
ALTER TABLE project ADD COLUMN IF NOT EXISTS provisioned_at TIMESTAMPTZ;
