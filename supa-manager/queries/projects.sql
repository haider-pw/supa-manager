-- name: GetProjectsForAccountId :many
SELECT p.*
FROM organization_membership om
         JOIN project p on om.organization_id = p.organization_id
WHERE account_id = $1;

-- name: CreateProject :one
INSERT INTO project (project_ref, project_name, organization_id, status, jwt_secret, cloud_provider, region)
VALUES ($1, $2, $3, 'PROVISIONING', $4, $5, $6)
RETURNING *;

-- name: GetProjectByRef :one
SELECT *
FROM project
WHERE project_ref = $1;

-- name: UpdateProjectStatus :one
UPDATE project
SET status = $2, updated_at = now()
WHERE project_ref = $1
RETURNING *;

-- name: UpdateProjectInfrastructure :one
UPDATE project
SET docker_compose_path = $2,
    docker_network_name = $3,
    postgres_port = $4,
    kong_http_port = $5,
    kong_https_port = $6,
    anon_key = $7,
    service_role_key = $8,
    provisioned_at = now(),
    updated_at = now()
WHERE project_ref = $1
RETURNING *;

-- name: GetProjectsByStatus :many
SELECT *
FROM project
WHERE status = $1
ORDER BY created_at DESC;

-- name: DeleteProject :exec
DELETE FROM project
WHERE project_ref = $1;