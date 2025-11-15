# API Reference

Complete reference for the SupaManager REST API.

---

## Base URL

```
http://localhost:8080
```

All API endpoints are prefixed with `/platform`.

---

## Authentication

### Login

**Endpoint:** `POST /platform/auth/login`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "account": {
    "id": 1,
    "email": "user@example.com",
    "created_at": "2024-11-15T10:00:00Z"
  }
}
```

**Status Codes:**
- `200` - Success
- `401` - Invalid credentials
- `400` - Bad request

### Register (if enabled)

**Endpoint:** `POST /platform/auth/register`

**Request Body:**
```json
{
  "email": "newuser@example.com",
  "password": "securePassword123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "account": {
    "id": 2,
    "email": "newuser@example.com",
    "created_at": "2024-11-15T10:30:00Z"
  }
}
```

**Note:** Registration is only available if `ALLOW_SIGNUP=true` in configuration.

---

## Profile

### Get Current User Profile

**Endpoint:** `GET /platform/profile`

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "id": 1,
  "email": "user@example.com",
  "created_at": "2024-11-15T10:00:00Z"
}
```

---

## Organizations

### List Organizations

**Endpoint:** `GET /platform/organizations`

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "My Organization",
    "slug": "my-organization",
    "created_at": "2024-11-15T10:00:00Z",
    "role": "owner"
  }
]
```

### Create Organization

**Endpoint:** `POST /platform/organizations`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "New Organization",
  "slug": "new-organization"
}
```

**Response:**
```json
{
  "id": 2,
  "name": "New Organization",
  "slug": "new-organization",
  "created_at": "2024-11-15T11:00:00Z"
}
```

### Get Organization

**Endpoint:** `GET /platform/organizations/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "id": 1,
  "name": "My Organization",
  "slug": "my-organization",
  "created_at": "2024-11-15T10:00:00Z",
  "members": [
    {
      "account_id": 1,
      "email": "user@example.com",
      "role": "owner"
    }
  ]
}
```

### Update Organization

**Endpoint:** `PATCH /platform/organizations/:id`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Updated Organization Name"
}
```

**Response:**
```json
{
  "id": 1,
  "name": "Updated Organization Name",
  "slug": "my-organization",
  "updated_at": "2024-11-15T12:00:00Z"
}
```

### Delete Organization

**Endpoint:** `DELETE /platform/organizations/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "message": "Organization deleted successfully"
}
```

---

## Projects

### List Projects

**Endpoint:** `GET /platform/projects`

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `org_id` (optional) - Filter by organization

**Response:**
```json
[
  {
    "id": 1,
    "project_ref": "upending-expectoration",
    "project_name": "Production API",
    "organization_id": 1,
    "status": "UNKNOWN",
    "cloud_provider": "K8S",
    "region": "US-EAST-1",
    "created_at": "2024-11-15T10:00:00Z"
  }
]
```

### Create Project

**Endpoint:** `POST /platform/projects`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "New Project",
  "organization_id": 1,
  "cloud_provider": "K8S",
  "db_region": "US-EAST-1"
}
```

**Response:**
```json
{
  "id": 2,
  "project_ref": "flying-rocket",
  "project_name": "New Project",
  "organization_id": 1,
  "status": "UNKNOWN",
  "cloud_provider": "K8S",
  "region": "US-EAST-1",
  "created_at": "2024-11-15T11:00:00Z"
}
```

**Note:** Currently returns status "UNKNOWN" as provisioning is not implemented.

### Get Project

**Endpoint:** `GET /platform/projects/:ref`

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "id": 1,
  "project_ref": "upending-expectoration",
  "project_name": "Production API",
  "organization_id": 1,
  "organization_name": "My Organization",
  "status": "UNKNOWN",
  "cloud_provider": "K8S",
  "region": "US-EAST-1",
  "created_at": "2024-11-15T10:00:00Z",
  "updated_at": "2024-11-15T10:00:00Z"
}
```

### Get Project API Credentials

**Endpoint:** `GET /platform/projects/:ref/api`

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "anon_key": "a.b.c",
  "service_key": "a.b.c",
  "endpoint": "https://upending-expectoration.supamanager.io"
}
```

**Note:** Currently returns placeholder keys. Real JWT keys will be generated when provisioning is implemented.

### Get Project Status

**Endpoint:** `GET /platform/projects/:ref/status`

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "UNKNOWN",
  "services": []
}
```

**Future Response:**
```json
{
  "status": "ACTIVE",
  "services": [
    {
      "name": "postgres",
      "status": "healthy",
      "uptime": "2h 30m"
    },
    {
      "name": "kong",
      "status": "healthy",
      "uptime": "2h 29m"
    }
  ]
}
```

### Delete Project

**Endpoint:** `DELETE /platform/projects/:ref`

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "message": "Project deleted successfully"
}
```

---

## Project Settings

### Get Project Settings

**Endpoint:** `GET /platform/projects/:ref/settings`

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "project_ref": "upending-expectoration",
  "project_name": "Production API",
  "region": "US-EAST-1",
  "postgres_version": "14.2",
  "auto_backup": true,
  "backup_retention_days": 7
}
```

### Update Project Settings

**Endpoint:** `PATCH /platform/projects/:ref/settings`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "project_name": "Updated Project Name",
  "auto_backup": false
}
```

**Response:**
```json
{
  "message": "Settings updated successfully"
}
```

---

## Usage & Billing

### Get Organization Usage

**Endpoint:** `GET /platform/organizations/:id/usage`

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "organization_id": 1,
  "period": "2024-11",
  "projects_count": 3,
  "total_database_size": "2.5 GB",
  "total_storage_size": "1.2 GB",
  "total_bandwidth": "45 GB"
}
```

---

## Health & Status

### API Health Check

**Endpoint:** `GET /health`

**No authentication required**

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2024-11-15T12:00:00Z"
}
```

### Get Platform Status

**Endpoint:** `GET /platform/status`

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "api_version": "1.0.0",
  "database_status": "healthy",
  "total_projects": 15,
  "active_projects": 12
}
```

---

## Error Responses

All error responses follow this format:

```json
{
  "error": "Error message description",
  "code": "ERROR_CODE",
  "status": 400
}
```

### Common Error Codes

| Status Code | Meaning |
|-------------|---------|
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Invalid or missing token |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource doesn't exist |
| 409 | Conflict - Resource already exists |
| 422 | Unprocessable Entity - Validation error |
| 500 | Internal Server Error |

---

## Rate Limiting

**Current:** No rate limiting implemented

**Future:**
- 100 requests per minute per IP
- 1000 requests per hour per user
- Headers: `X-RateLimit-Limit`, `X-RateLimit-Remaining`

---

## Pagination

**Current:** No pagination (returns all results)

**Future:**
```
GET /platform/projects?page=1&per_page=20
```

**Response Headers:**
```
X-Total-Count: 150
X-Page: 1
X-Per-Page: 20
X-Total-Pages: 8
```

---

## Webhooks (Planned)

Future support for webhooks on events:
- `project.created`
- `project.deleted`
- `project.status_changed`
- `organization.member_added`
- `organization.member_removed`

---

## API Client Examples

### cURL

```bash
# Login
curl -X POST http://localhost:8080/platform/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'

# Create Project
curl -X POST http://localhost:8080/platform/projects \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"My Project","organization_id":1,"cloud_provider":"K8S","db_region":"US-EAST-1"}'
```

### JavaScript/Fetch

```javascript
// Login
const response = await fetch('http://localhost:8080/platform/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'user@example.com',
    password: 'password'
  })
});
const { token } = await response.json();

// Get Projects
const projects = await fetch('http://localhost:8080/platform/projects', {
  headers: { 'Authorization': `Bearer ${token}` }
}).then(r => r.json());
```

### Python

```python
import requests

# Login
response = requests.post('http://localhost:8080/platform/auth/login', json={
    'email': 'user@example.com',
    'password': 'password'
})
token = response.json()['token']

# Get Projects
projects = requests.get(
    'http://localhost:8080/platform/projects',
    headers={'Authorization': f'Bearer {token}'}
).json()
```

---

## SDK Support

**Current:** No official SDK

**Planned:**
- JavaScript/TypeScript SDK
- Python SDK
- Go SDK

---

## Related Documentation

- [Architecture Overview](Architecture-Overview) - System architecture
- [Configuration Reference](Configuration-Reference) - API configuration
- [Development Guide](Development-Guide) - Local development

---

**Next:** Learn about [Configuration Options](Configuration-Reference)
