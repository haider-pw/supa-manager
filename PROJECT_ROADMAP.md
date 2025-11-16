# SupaManager Project Roadmap

**Last Updated:** November 16, 2025
**Project Status:** Phase 3 Complete, Production-Ready for Single-Server Deployments
**Version:** 1.0.0-beta

---

## Executive Summary

### Current State

SupaManager is a functional system for managing self-hosted Supabase instances through the Supabase Studio UI. As of November 16, 2025:

**Completed:**
- Core API infrastructure (Go/Gin, PostgreSQL, sqlc)
- User authentication and organization management
- Supabase Studio UI integration (patched v1.24.04)
- Docker Compose-based provisioning (Phase 3)
- Project lifecycle management (create, pause, resume, delete)
- Real JWT key generation and management
- Database migrations and schema management
- Basic health monitoring
- 38 API endpoints (90% Studio-compatible)

**Infrastructure:**
- 55 Go source files
- 38 API handlers
- 2 database migrations
- Zero test files (critical gap)
- No CI/CD pipeline
- No monitoring/alerting

### Vision

Transform SupaManager into a production-grade, enterprise-ready platform for managing multiple Supabase instances with:
- Full lifecycle management
- Comprehensive security
- High availability
- Multi-region support
- Advanced monitoring
- Complete test coverage
- Professional documentation

### Timeline Estimate

- **Phase 4 (Testing & Quality):** 3-4 weeks
- **Phase 5 (Security Hardening):** 2-3 weeks
- **Phase 6 (Production Readiness):** 3-4 weeks
- **Phase 7 (Feature Completeness):** 4-6 weeks
- **Phase 8 (Enterprise Features):** 6-8 weeks

**Total estimated time to full production readiness:** 18-25 weeks (4-6 months)

---

## Phase 4: Testing, Quality Assurance & Stability

**Priority:** P0 (Critical)
**Estimated Effort:** 3-4 weeks
**Dependencies:** Phase 3 completion
**Status:** Not Started

### Overview

Establish comprehensive testing infrastructure and improve code quality. This is critical before any production deployment.

### 4.1 Testing Infrastructure Setup (P0)

**Estimated Effort:** 1 week

- [ ] Set up Go testing framework structure
- [ ] Configure testcontainers for PostgreSQL integration tests
- [ ] Set up Docker test environment
- [ ] Create test database seed data
- [ ] Establish test data factories/fixtures
- [ ] Configure code coverage reporting (target: 80%+)
- [ ] Set up test logging and debugging

**Files to Create:**
```
supa-manager/
├── testing/
│   ├── fixtures/
│   │   ├── accounts.go
│   │   ├── organizations.go
│   │   └── projects.go
│   ├── helpers.go
│   └── testcontainers.go
└── Makefile (for test commands)
```

### 4.2 Unit Tests (P0)

**Estimated Effort:** 1.5 weeks
**Target Coverage:** 80%+

**Priority Areas:**

1. **Database Layer (database/)**
   - [ ] Test all sqlc-generated queries
   - [ ] Test transaction handling
   - [ ] Test error cases (constraints, foreign keys)
   - [ ] Test NULL handling for nullable fields

2. **Utilities (utils/)**
   - [ ] Test slug generation (uniqueness, format)
   - [ ] Test SQL type conversions
   - [ ] Test edge cases (empty strings, special characters)

3. **Permissions (permisions/)**
   - [ ] Test organization membership checks
   - [ ] Test role-based access control
   - [ ] Test permission denial cases

4. **Provisioner (provisioner/)**
   - [ ] Test port allocation logic
   - [ ] Test secret generation (randomness, length)
   - [ ] Test template rendering
   - [ ] Test configuration validation
   - [ ] Mock Docker client for unit tests

**Example Test Structure:**
```go
// supa-manager/provisioner/ports_test.go
package provisioner

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestAllocatePort(t *testing.T) {
    tests := []struct {
        name          string
        projectID     int32
        basePort      int
        expectedPort  int
    }{
        {"First project", 1, 5433, 5433},
        {"Second project", 2, 5433, 5434},
        {"Port overflow", 1000, 5433, 6433},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            port := AllocatePostgresPort(tt.projectID, tt.basePort)
            assert.Equal(t, tt.expectedPort, port)
        })
    }
}
```

### 4.3 Integration Tests (P0)

**Estimated Effort:** 1 week

- [ ] **API Endpoint Tests**
  - [ ] Test all 38 API handlers
  - [ ] Test authentication/authorization
  - [ ] Test request validation
  - [ ] Test error responses
  - [ ] Test CORS configuration

- [ ] **Database Integration Tests**
  - [ ] Test migrations (up/down)
  - [ ] Test constraint enforcement
  - [ ] Test transaction rollback
  - [ ] Test connection pooling

- [ ] **Provisioner Integration Tests**
  - [ ] Test full project provisioning flow
  - [ ] Test pause/resume functionality
  - [ ] Test project deletion and cleanup
  - [ ] Test concurrent provisioning
  - [ ] Test error recovery

**Example Integration Test:**
```go
// supa-manager/api/postPlatformProjects_test.go
func TestCreateProject_Integration(t *testing.T) {
    // Setup test database
    db := testing.SetupTestDB(t)
    defer db.Cleanup()

    // Create test user and org
    user := testing.CreateTestUser(db)
    org := testing.CreateTestOrg(db, user.ID)
    token := testing.GenerateJWT(user.GotrueID)

    // Make request
    resp := testing.PostJSON("/platform/projects", map[string]interface{}{
        "name": "test-project",
        "organization_id": org.ID,
    }, token)

    assert.Equal(t, 201, resp.StatusCode)

    // Verify project created
    project := db.GetProjectByRef(resp.Body["ref"].(string))
    assert.NotNil(t, project)
    assert.Equal(t, "PROVISIONING", project.Status)
}
```

### 4.4 End-to-End Tests (P1)

**Estimated Effort:** 0.5 weeks

- [ ] Test complete user workflow (signup → create org → create project)
- [ ] Test project lifecycle (create → pause → resume → delete)
- [ ] Test multi-user scenarios
- [ ] Test Studio UI integration

### 4.5 Error Handling Improvements (P0)

**Estimated Effort:** 0.5 weeks

**Current Issues:**
- Inconsistent error responses
- Some errors leak internal details
- Missing validation in many endpoints

**Tasks:**
- [ ] Create standardized error response format
- [ ] Implement proper error logging (without exposing to clients)
- [ ] Add input validation to all endpoints
- [ ] Sanitize error messages (no SQL, no stack traces)
- [ ] Add error codes for client handling

**Standard Error Format:**
```go
type APIError struct {
    Code    string `json:"code"`    // e.g., "PROJECT_NOT_FOUND"
    Message string `json:"message"` // User-friendly message
    Details any    `json:"details,omitempty"` // Optional context
}
```

### 4.6 Code Quality Improvements (P1)

**Estimated Effort:** 0.5 weeks

- [ ] Run `go vet` and fix all issues
- [ ] Run `golangci-lint` and address findings
- [ ] Add godoc comments to all exported functions
- [ ] Remove all commented-out code
- [ ] Address all TODO/FIXME comments (20+ found)
- [ ] Implement consistent naming conventions
- [ ] Extract magic numbers to constants

**Priority TODOs to Address:**
```go
// api/postGotrueToken.go
// TODO: use a real secret (currently using hardcoded test secret)
// TODO: support refresh tokens
// TODO: look at GoTrue code for more fields

// api/postPlatformOrganizations.go
// TODO: use tx (currently no transaction)

// api/postPlatformProjects.go
// TODO: Uncomment when ProvisionProject is implemented (DONE in Phase 3)
```

---

## Phase 5: Security Hardening

**Priority:** P0 (Critical)
**Estimated Effort:** 2-3 weeks
**Dependencies:** Phase 4 testing complete
**Status:** Not Started

### Overview

Address critical security vulnerabilities and implement production-grade security measures.

### 5.1 Critical Security Fixes (P0)

**Estimated Effort:** 1 week

#### 5.1.1 JWT Secret Management

**Current Issue:** JWT_SECRET from environment variable only

**Fix:**
- [ ] Implement secret rotation mechanism
- [ ] Store secret history for token validation during rotation
- [ ] Add JWT token expiration (currently no expiry!)
- [ ] Implement refresh token mechanism
- [ ] Add JWT token revocation (blacklist)

```go
// Example implementation
type JWTManager struct {
    currentSecret  string
    previousSecret string // For rotation period
    secretRotatedAt time.Time
}

func (j *JWTManager) ValidateToken(token string) (*Claims, error) {
    // Try current secret
    claims, err := validateWithSecret(token, j.currentSecret)
    if err == nil {
        return claims, nil
    }

    // Try previous secret (grace period)
    if time.Since(j.secretRotatedAt) < 24*time.Hour {
        return validateWithSecret(token, j.previousSecret)
    }

    return nil, ErrInvalidToken
}
```

#### 5.1.2 Encryption Secret Management

**Current Issue:** ENCRYPTION_SECRET stored in plain text .env

**Fix:**
- [ ] Integrate with HashiCorp Vault or AWS Secrets Manager
- [ ] Implement secret encryption at rest
- [ ] Add secret audit logging
- [ ] Implement least-privilege access
- [ ] Document secret rotation procedure

#### 5.1.3 Input Validation

**Current Issue:** Minimal input validation in API handlers

**Fix:**
- [ ] Add validation middleware for all endpoints
- [ ] Validate all string inputs (length, format, allowed characters)
- [ ] Validate numeric ranges
- [ ] Sanitize all user inputs
- [ ] Add SQL injection protection validation
- [ ] Add XSS protection for any HTML rendering

```go
// Example validation
type CreateProjectRequest struct {
    Name           string `json:"name" binding:"required,min=3,max=64,alphanum"`
    OrganizationID int32  `json:"organization_id" binding:"required,gt=0"`
    Region         string `json:"region" binding:"required,oneof=us-east us-west eu-west"`
}
```

#### 5.1.4 SQL Injection Protection

**Status:** Mostly protected (using sqlc)

**Verify:**
- [ ] Audit all sqlc queries for dynamic SQL
- [ ] Ensure no string concatenation in SQL
- [ ] Verify prepared statement usage
- [ ] Add integration tests for SQL injection attempts

#### 5.1.5 Authentication Vulnerabilities

**Issues Found:**
```go
// api/postGotrueToken.go:
// TODO: use a real secret (CRITICAL!)
```

**Fix:**
- [ ] Replace hardcoded test secret with real secret
- [ ] Add rate limiting to auth endpoints (prevent brute force)
- [ ] Implement account lockout after failed attempts
- [ ] Add MFA support (optional but recommended)
- [ ] Log all authentication attempts

### 5.2 Authorization & RBAC (P0)

**Estimated Effort:** 1 week

**Current State:** Basic org membership checks, incomplete RBAC

**Improvements:**
- [ ] Define complete permission model
- [ ] Implement role hierarchy (admin > owner > member > viewer)
- [ ] Add per-project permissions
- [ ] Implement API key scoping
- [ ] Add audit trail for all permission checks

**Permission Model:**
```go
type Role string

const (
    RoleAdmin   Role = "admin"   // Full system access
    RoleOwner   Role = "owner"   // Organization owner
    RoleMember  Role = "member"  // Project member
    RoleViewer  Role = "viewer"  // Read-only access
)

type Permission string

const (
    PermProjectCreate  Permission = "project:create"
    PermProjectRead    Permission = "project:read"
    PermProjectUpdate  Permission = "project:update"
    PermProjectDelete  Permission = "project:delete"
    PermProjectPause   Permission = "project:pause"
    PermBackupCreate   Permission = "backup:create"
    PermBackupRestore  Permission = "backup:restore"
)
```

### 5.3 Rate Limiting & DDoS Protection (P1)

**Estimated Effort:** 0.5 weeks

- [ ] Implement rate limiting middleware
- [ ] Add per-IP rate limits
- [ ] Add per-user rate limits
- [ ] Add per-endpoint rate limits
- [ ] Implement exponential backoff for repeated violations
- [ ] Add rate limit response headers (X-RateLimit-*)

```go
// Example rate limits
var RateLimits = map[string]RateLimit{
    "/auth/token":           {Requests: 5, Window: time.Minute},
    "/platform/projects":    {Requests: 10, Window: time.Minute},
    "/platform/projects/*":  {Requests: 100, Window: time.Minute},
}
```

### 5.4 Secrets Management in Containers (P0)

**Estimated Effort:** 0.5 weeks

**Current Issue:** Project secrets stored in plain text .env files

**Fix:**
- [ ] Implement secrets encryption at rest
- [ ] Use Docker secrets for sensitive values
- [ ] Encrypt backup files containing secrets
- [ ] Add secret rotation for provisioned projects
- [ ] Implement secure secret distribution to containers

### 5.5 Security Audit & Penetration Testing (P1)

**Estimated Effort:** 0.5 weeks

- [ ] Run automated security scanner (gosec, snyk)
- [ ] Perform manual code review for security issues
- [ ] Test common OWASP Top 10 vulnerabilities
- [ ] Document security findings and remediation
- [ ] Create security incident response plan

### 5.6 Audit Logging (P1)

**Estimated Effort:** 0.5 weeks

- [ ] Implement audit log storage (database table)
- [ ] Log all authentication events
- [ ] Log all authorization failures
- [ ] Log all project lifecycle events
- [ ] Log all configuration changes
- [ ] Add audit log retention policy
- [ ] Create audit log query API

```sql
CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id TEXT NOT NULL,
    action TEXT NOT NULL,
    resource_type TEXT NOT NULL,
    resource_id TEXT,
    ip_address INET,
    user_agent TEXT,
    success BOOLEAN NOT NULL,
    error_message TEXT,
    metadata JSONB
);
```

---

## Phase 6: Production Readiness

**Priority:** P0 (Critical)
**Estimated Effort:** 3-4 weeks
**Dependencies:** Phase 5 security complete
**Status:** Not Started

### Overview

Prepare system for production deployment with monitoring, documentation, and operational procedures.

### 6.1 Monitoring & Observability (P0)

**Estimated Effort:** 1.5 weeks

#### 6.1.1 Prometheus Metrics

- [ ] Install Prometheus Go client
- [ ] Expose /metrics endpoint
- [ ] Add application metrics:
  - [ ] HTTP request duration histogram
  - [ ] HTTP request count by status code
  - [ ] Database connection pool stats
  - [ ] Active projects gauge
  - [ ] Provisioning duration histogram
  - [ ] Provisioning success/failure counter
  - [ ] Docker container health status
  - [ ] API endpoint latency by endpoint

```go
// Example metrics
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )

    provisioningDuration = prometheus.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "provisioning_duration_seconds",
            Help:    "Time taken to provision a project",
            Buckets: []float64{30, 60, 120, 300, 600},
        },
    )
)
```

#### 6.1.2 Structured Logging

**Current:** Basic slog usage
**Improve:**
- [ ] Add correlation IDs to all requests
- [ ] Log levels per environment (debug/info/warn/error)
- [ ] Add contextual logging (user ID, project ID, org ID)
- [ ] Implement log aggregation (e.g., Loki)
- [ ] Add log sampling for high-volume endpoints
- [ ] Create log parsing rules for alerts

```go
// Example structured logging
logger.Info("Project provisioned successfully",
    slog.String("project_id", projectID),
    slog.String("user_id", userID),
    slog.String("correlation_id", correlationID),
    slog.Duration("duration", elapsed),
    slog.Int("containers_started", 12),
)
```

#### 6.1.3 Grafana Dashboards

- [ ] Create dashboard: System Overview
  - Request rate, error rate, latency
  - Database connections, query duration
  - Resource usage (CPU, memory, disk)

- [ ] Create dashboard: Projects
  - Total projects by status
  - Provisioning success rate
  - Project resource usage
  - Failed provisioning count

- [ ] Create dashboard: Security
  - Authentication failures
  - Rate limit violations
  - Permission denials
  - Suspicious activity

### 6.2 Alerting (P0)

**Estimated Effort:** 0.5 weeks

- [ ] Configure Alertmanager
- [ ] Define alert rules:
  - [ ] High error rate (>1% of requests)
  - [ ] Slow API responses (p95 > 1s)
  - [ ] Provisioning failures
  - [ ] Database connection issues
  - [ ] Disk space low (<10%)
  - [ ] Memory usage high (>85%)
  - [ ] Authentication failures spike
  - [ ] Service down

- [ ] Set up notification channels:
  - [ ] Email alerts
  - [ ] Slack integration
  - [ ] PagerDuty for critical alerts

### 6.3 High Availability (P1)

**Estimated Effort:** 1 week

**Current:** Single instance, single database

**Improvements:**
- [ ] Database replication (PostgreSQL streaming replication)
- [ ] API horizontal scaling (multiple instances behind load balancer)
- [ ] Shared session storage (Redis)
- [ ] Health check endpoints
- [ ] Graceful shutdown handling
- [ ] Connection draining
- [ ] Circuit breaker pattern for external services

### 6.4 Disaster Recovery (P0)

**Estimated Effort:** 0.5 weeks

- [ ] Document backup procedures
- [ ] Automate database backups (daily)
- [ ] Test restore procedures
- [ ] Document recovery time objectives (RTO)
- [ ] Document recovery point objectives (RPO)
- [ ] Create disaster recovery runbook
- [ ] Implement point-in-time recovery

**Recovery Procedures:**
```bash
# Database Recovery
1. Stop API service
2. Restore database from backup
3. Verify data integrity
4. Start API service
5. Verify system health

# Full System Recovery
1. Provision new infrastructure
2. Restore database
3. Restore project volumes
4. Update DNS
5. Verify all projects accessible
```

### 6.5 Documentation (P0)

**Estimated Effort:** 1 week

#### 6.5.1 API Documentation

- [ ] Generate OpenAPI/Swagger spec
- [ ] Document all 38+ endpoints
- [ ] Add request/response examples
- [ ] Document error codes
- [ ] Add authentication guide
- [ ] Publish API docs (Swagger UI)

#### 6.5.2 Deployment Documentation

- [ ] Docker deployment guide
- [ ] Kubernetes deployment guide
- [ ] Configuration reference (all env vars)
- [ ] Scaling guide
- [ ] Upgrade procedures
- [ ] Rollback procedures

#### 6.5.3 Operational Runbooks

- [ ] Incident response procedures
- [ ] Troubleshooting guide (common issues)
- [ ] Performance tuning guide
- [ ] Security incident response
- [ ] Backup and restore procedures
- [ ] Database maintenance tasks

#### 6.5.4 Developer Documentation

- [ ] Architecture overview
- [ ] Code structure guide
- [ ] Contributing guidelines
- [ ] Local development setup
- [ ] Testing guide
- [ ] Release process

### 6.6 CI/CD Pipeline (P1)

**Estimated Effort:** 0.5 weeks

**Current:** Manual builds and deploys

**Implement:**
- [ ] GitHub Actions workflows
  - [ ] Build and test on PR
  - [ ] Run linters and security scans
  - [ ] Run all tests (unit, integration, e2e)
  - [ ] Build Docker images
  - [ ] Deploy to staging on merge to main
  - [ ] Deploy to production on tag

```yaml
# .github/workflows/test.yml
name: Test
on: [pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: go test -v -race -coverprofile=coverage.txt ./...
      - run: go vet ./...
      - run: golangci-lint run
      - uses: codecov/codecov-action@v3
```

---

## Phase 7: Feature Completeness

**Priority:** P1 (High)
**Estimated Effort:** 4-6 weeks
**Dependencies:** Phase 6 production readiness
**Status:** Not Started

### Overview

Implement missing features to achieve parity with Supabase Cloud platform.

### 7.1 Backup & Restore System (P0)

**Estimated Effort:** 2 weeks

**Reference:** BACKUP_DESIGN.md

**Implementation:**
- [ ] Implement backup creation
  - [ ] Database backup (pg_dump)
  - [ ] Storage files backup
  - [ ] Configuration backup
  - [ ] Full project export

- [ ] Implement restore functionality
  - [ ] Database restore
  - [ ] Storage restore
  - [ ] Configuration restore
  - [ ] Full project import

- [ ] Automated backup scheduling
  - [ ] Daily backups
  - [ ] Configurable retention (default: 7 days)
  - [ ] Backup verification

- [ ] Backup storage options
  - [ ] Local storage (default)
  - [ ] S3-compatible storage
  - [ ] Backup encryption

**API Endpoints:**
```
POST   /platform/projects/:ref/backups
GET    /platform/projects/:ref/backups
GET    /platform/projects/:ref/backups/:id
DELETE /platform/projects/:ref/backups/:id
POST   /platform/projects/:ref/backups/:id/restore
GET    /platform/projects/:ref/backups/:id/download
```

### 7.2 Resource Quotas & Limits (P1)

**Estimated Effort:** 2 weeks

**Reference:** QUOTAS_DESIGN.md

**Implementation:**
- [ ] Database schema for quotas
- [ ] Quota enforcement middleware
- [ ] Usage tracking system
- [ ] Quota violation handling
- [ ] User notifications
- [ ] Admin quota management UI

**Quota Types:**
- [ ] Storage quotas (database, files, backups)
- [ ] Compute quotas (CPU, memory)
- [ ] Network quotas (bandwidth, requests/hour)
- [ ] Feature limits (max backups, max users)

### 7.3 Analytics & Metrics (P1)

**Estimated Effort:** 1.5 weeks

**Current:** Stub implementation returning empty data

**Implement:**
- [ ] Database for analytics storage
- [ ] API usage tracking
- [ ] Request count aggregation
- [ ] Error rate tracking
- [ ] Response time tracking
- [ ] Database query statistics
- [ ] Storage usage over time

**Endpoints to Implement:**
```
GET /projects/:ref/analytics/endpoints/usage.api-counts
GET /projects/:ref/analytics/endpoints/usage.api-requests-count
GET /projects/:ref/analytics/endpoints/usage.api-response-times
```

### 7.4 pg-meta Integration (P1)

**Estimated Effort:** 1 week

**Current:** Endpoints exist but return empty data

**Implement:**
- [ ] Connect to project databases
- [ ] Query PostgreSQL metadata
- [ ] Return table schemas
- [ ] Return custom types
- [ ] Return publications
- [ ] Return functions/procedures
- [ ] Return extensions

**Endpoints:**
```
POST /platform/pg-meta/:ref/query
GET  /platform/pg-meta/:ref/tables
GET  /platform/pg-meta/:ref/columns
GET  /platform/pg-meta/:ref/types
GET  /platform/pg-meta/:ref/publications
GET  /platform/pg-meta/:ref/functions
```

### 7.5 Project Settings Management (P2)

**Estimated Effort:** 1 week

- [ ] Database configuration updates
- [ ] JWT secret rotation
- [ ] Custom domains
- [ ] SSL certificate management
- [ ] Project transfer between organizations
- [ ] Project rename
- [ ] Project region migration

### 7.6 Integrations (P2)

**Estimated Effort:** 1.5 weeks

**Current:** Stub implementations

**Implement:**
- [ ] GitHub OAuth integration
- [ ] GitHub repository listing
- [ ] Vercel deployment integration
- [ ] Netlify integration
- [ ] Custom webhooks
- [ ] Slack notifications
- [ ] Email notifications (SMTP)

---

## Phase 8: Enterprise Features

**Priority:** P2 (Medium)
**Estimated Effort:** 6-8 weeks
**Dependencies:** Phase 7 feature completeness
**Status:** Future

### 8.1 Kubernetes Support (P2)

**Estimated Effort:** 3 weeks

**Current:** Docker Compose only (single server)

**Implement:**
- [ ] Kubernetes provisioner implementation
- [ ] Helm charts for projects
- [ ] Namespace isolation
- [ ] Service mesh integration (Istio/Linkerd)
- [ ] Horizontal pod autoscaling
- [ ] Persistent volume claims
- [ ] StatefulSets for databases

### 8.2 Multi-Region Support (P2)

**Estimated Effort:** 2 weeks

- [ ] Region configuration
- [ ] Region-specific provisioning
- [ ] Cross-region project migration
- [ ] Region health monitoring
- [ ] Geo-routing (DNS-based)
- [ ] Data residency compliance

### 8.3 Advanced Monitoring (P2)

**Estimated Effort:** 1.5 weeks

- [ ] Distributed tracing (Jaeger/Zipkin)
- [ ] Application Performance Monitoring (APM)
- [ ] Real-time dashboards
- [ ] Custom metrics API
- [ ] Anomaly detection
- [ ] Capacity planning tools

### 8.4 Team Management (P2)

**Estimated Effort:** 1 week

- [ ] Team creation
- [ ] Team permissions
- [ ] Project sharing with teams
- [ ] Team activity logs
- [ ] Team billing

### 8.5 Billing & Usage (P3)

**Estimated Effort:** 2 weeks

- [ ] Usage metering
- [ ] Billing plans
- [ ] Stripe integration
- [ ] Invoice generation
- [ ] Usage-based pricing
- [ ] Payment methods

---

## Quick Wins (Immediate Actions)

**These can be done in 1-2 days and provide immediate value:**

### 1. Fix Critical TODOs (4 hours)

```go
// api/postGotrueToken.go
- Replace hardcoded test secret with real JWT generation
- Add token expiration
- Add refresh token support
```

### 2. Add Basic Rate Limiting (4 hours)

```go
import "github.com/didip/tollbooth"

// Add to api.go Router()
limiter := tollbooth.NewLimiter(10, nil) // 10 req/sec
r.Use(tollbooth.LimitHandler(limiter))
```

### 3. Implement Health Check Endpoint (2 hours)

```go
func (a *Api) healthz(c *gin.Context) {
    // Check database
    if err := a.pgPool.Ping(c); err != nil {
        c.JSON(503, gin.H{"status": "unhealthy", "database": "down"})
        return
    }

    // Check Docker (if provisioner enabled)
    if a.provisioner != nil {
        // Check Docker daemon
    }

    c.JSON(200, gin.H{"status": "healthy"})
}
```

### 4. Add Correlation IDs (3 hours)

```go
func CorrelationIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        correlationID := c.GetHeader("X-Correlation-ID")
        if correlationID == "" {
            correlationID = uuid.New().String()
        }
        c.Set("correlation_id", correlationID)
        c.Header("X-Correlation-ID", correlationID)
        c.Next()
    }
}
```

### 5. Add Makefile for Common Tasks (2 hours)

```makefile
.PHONY: test build run lint

test:
	go test -v -race -coverprofile=coverage.txt ./...

build:
	go build -o bin/supa-manager ./supa-manager

run:
	go run ./supa-manager/main.go

lint:
	golangci-lint run

docker-build:
	docker build -t supa-manager:latest .
```

### 6. Add .gitignore Improvements (1 hour)

```gitignore
# Binaries
bin/
*.exe

# Test coverage
coverage.txt
*.out

# Project data
projects/
backups/

# IDE
.vscode/
.idea/
*.swp
```

---

## Technical Debt Backlog

**Address these in future iterations:**

### High Priority Debt

1. **Zero Test Coverage** (P0)
   - No tests for 55 Go files
   - High risk for regressions
   - Estimated: 3 weeks to achieve 80% coverage

2. **Inconsistent Error Handling** (P0)
   - Some endpoints leak internal errors
   - No standard error format
   - Estimated: 1 week

3. **Missing Input Validation** (P0)
   - Many endpoints don't validate inputs
   - SQL injection risk mitigation needed
   - Estimated: 1 week

4. **Hardcoded Configuration** (P1)
   - Listen address hardcoded (":8080")
   - Some secrets in code
   - Estimated: 2 days

5. **No Transaction Management** (P1)
   - Some operations should be atomic
   - TODO in postPlatformOrganizations.go
   - Estimated: 3 days

### Medium Priority Debt

6. **CORS Too Permissive** (P1)
   - Currently allows all origins ("*")
   - Should be configurable
   - Estimated: 2 hours

7. **No Request Timeout** (P1)
   - Long-running requests can hang
   - Need timeout middleware
   - Estimated: 4 hours

8. **Poor Logging Context** (P2)
   - Logs lack correlation IDs
   - Hard to trace request flows
   - Estimated: 1 day

9. **Magic Numbers** (P2)
   - Port calculations have magic numbers
   - Extract to constants
   - Estimated: 4 hours

10. **Code Duplication** (P2)
    - Auth checking duplicated across handlers
    - Need middleware
    - Estimated: 1 day

---

## Security Issues to Fix ASAP

**Critical (Fix Before Production):**

1. **JWT Token Has No Expiration** (P0)
   - Current: Tokens never expire
   - Risk: Stolen tokens valid forever
   - Fix: Add exp claim, default 24 hours
   - Estimated: 4 hours

2. **Hardcoded Test Secret** (P0)
   - Location: api/postGotrueToken.go
   - Risk: Predictable tokens
   - Fix: Use proper JWT generation
   - Estimated: 2 hours

3. **Weak CORS Configuration** (P0)
   - Current: AllowOrigins: ["*"]
   - Risk: CSRF attacks
   - Fix: Configure specific origins
   - Estimated: 1 hour

4. **No Rate Limiting** (P0)
   - Risk: Brute force, DoS attacks
   - Fix: Implement per-IP rate limiting
   - Estimated: 4 hours

5. **Secrets in Plain Text Files** (P0)
   - Current: Project .env files unencrypted
   - Risk: Credential exposure
   - Fix: Encrypt at rest
   - Estimated: 1 day

**High (Fix Soon):**

6. **No Input Length Limits** (P1)
   - Risk: Buffer overflow, DoS
   - Fix: Add max length validation
   - Estimated: 1 day

7. **No Account Lockout** (P1)
   - Risk: Brute force attacks
   - Fix: Lock after 5 failed attempts
   - Estimated: 4 hours

8. **Missing HTTPS Enforcement** (P1)
   - Current: HTTP allowed
   - Risk: Man-in-the-middle
   - Fix: Redirect HTTP to HTTPS
   - Estimated: 2 hours

9. **No Audit Logging** (P1)
   - Risk: Can't track security events
   - Fix: Implement audit log
   - Estimated: 1 week

10. **Docker Socket Exposure** (P1)
    - Current: Mounted in container
    - Risk: Container escape
    - Fix: Use Docker API proxy
    - Estimated: 1 day

---

## Installation & Usage Guide

**Target Audience:** DevOps engineers, system administrators

### Prerequisites

- Ubuntu 20.04+ or similar Linux distribution
- Docker 20.10+ with Compose V2
- 8GB RAM minimum (16GB recommended)
- 100GB free disk space
- PostgreSQL 14+ (for management database)

### Quick Start (Development)

```bash
# 1. Clone repository
git clone https://github.com/haider-pw/supa-manager.git
cd supa-manager

# 2. Build Studio image
cd studio
./build.sh v1.24.04 supa-manager/studio:v1.24.04 .env
cd ..

# 3. Configure environment
cp .env.example .env
# Edit .env - MUST change JWT_SECRET and ENCRYPTION_SECRET for production!

# 4. Start services
docker compose up -d

# 5. Access Studio
open http://localhost:3000
```

### Production Deployment

**See detailed guide:** [wiki/Deployment.md]

**Key Steps:**
1. Generate strong secrets (openssl rand -base64 48)
2. Configure external PostgreSQL database
3. Set up HTTPS/TLS certificates
4. Configure firewall rules
5. Set up monitoring
6. Configure backups
7. Test disaster recovery

### Configuration Reference

**Essential Environment Variables:**

```bash
# Database
DATABASE_URL=postgres://user:pass@host:5432/supabase

# Security (REQUIRED - change from defaults!)
JWT_SECRET=your-strong-secret-min-64-chars
ENCRYPTION_SECRET=your-different-strong-secret-min-64-chars

# Provisioning
PROVISIONING_ENABLED=true
PROVISIONING_DOCKER_HOST=unix:///var/run/docker.sock
PROVISIONING_PROJECTS_DIR=/var/lib/supamanager/projects
PROVISIONING_BASE_POSTGRES_PORT=5433
PROVISIONING_BASE_KONG_HTTP_PORT=54321

# Application
ALLOW_SIGNUP=true
DOMAIN_STUDIO_URL=https://studio.yourdomain.com
DOMAIN_BASE=yourdomain.com

# DNS Webhook (optional)
DOMAIN_DNS_HOOK_URL=http://dns-service:8081
DOMAIN_DNS_HOOK_KEY=your-webhook-secret
```

**See full reference:** [wiki/Configuration-Reference.md]

### Troubleshooting

**Common Issues:**

1. **Provisioning stuck in PROVISIONING status**
   - Check Docker daemon: `docker info`
   - Check logs: `docker compose logs supa-manager`
   - Check port availability: `lsof -i :5433`

2. **Studio can't connect to API**
   - Verify NEXT_PUBLIC_API_URL in studio/.env
   - Check CORS configuration in api/api.go
   - Check network connectivity

3. **Database migration errors**
   - Check migrations table: `SELECT * FROM migrations;`
   - Manually run migration if needed
   - Restart supa-manager service

**See full guide:** [wiki/Troubleshooting.md]

---

## Best Practices & Standards

### Code Organization

**Follow these patterns:**

1. **API Handlers** - One file per endpoint
   - Location: `api/`
   - Naming: `{method}{Path}.go` (e.g., `getProjectStatus.go`)
   - Pattern: Extract logic to separate functions

2. **Database Queries** - sqlc only
   - Location: `queries/*.sql`
   - Run `sqlc generate` after changes
   - Never write raw SQL in Go code

3. **Configuration** - Environment variables
   - Define in `conf/config.go`
   - Use `envconfig` tags
   - Provide sensible defaults

4. **Errors** - Structured error handling
   - Log internally with context
   - Return user-friendly messages
   - Use standard error codes

### Error Handling Pattern

```go
func (a *Api) someHandler(c *gin.Context) {
    // Parse request
    var req SomeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        a.logger.Error("Failed to parse request",
            slog.String("error", err.Error()),
            slog.String("correlation_id", c.GetString("correlation_id")),
        )
        c.JSON(400, gin.H{
            "code": "INVALID_REQUEST",
            "message": "Invalid request format",
        })
        return
    }

    // Validate request
    if err := req.Validate(); err != nil {
        c.JSON(400, gin.H{
            "code": "VALIDATION_ERROR",
            "message": err.Error(),
        })
        return
    }

    // Get authenticated user
    account, err := a.GetAccountFromRequest(c)
    if err != nil {
        c.JSON(401, gin.H{
            "code": "UNAUTHORIZED",
            "message": "Invalid or missing authentication",
        })
        return
    }

    // Check permissions
    if !hasPermission(account, req.ProjectID) {
        c.JSON(403, gin.H{
            "code": "FORBIDDEN",
            "message": "You don't have access to this resource",
        })
        return
    }

    // Perform operation
    result, err := a.performOperation(c.Request.Context(), req)
    if err != nil {
        a.logger.Error("Operation failed",
            slog.String("error", err.Error()),
            slog.String("user_id", account.GotrueID),
        )
        c.JSON(500, gin.H{
            "code": "INTERNAL_ERROR",
            "message": "An error occurred processing your request",
        })
        return
    }

    c.JSON(200, result)
}
```

### Testing Requirements

**All new code must include:**

1. **Unit Tests** - Test logic in isolation
   - Minimum 80% code coverage
   - Mock external dependencies
   - Test error cases

2. **Integration Tests** - Test API endpoints
   - Test happy path
   - Test error responses
   - Test authentication/authorization

3. **Test Naming** - Clear and descriptive
   ```go
   func TestCreateProject_WithValidData_ReturnsCreated(t *testing.T)
   func TestCreateProject_WithoutAuth_ReturnsUnauthorized(t *testing.T)
   func TestCreateProject_WithInvalidName_ReturnsBadRequest(t *testing.T)
   ```

### Documentation Standards

**Every exported function/type must have godoc:**

```go
// CreateProject provisions a new Supabase instance for the given organization.
// It generates all required secrets, allocates ports, and starts Docker containers.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - config: Project configuration including name, org, and resources
//
// Returns:
//   - ProjectInfo: Information about the provisioned project
//   - error: Error if provisioning fails
//
// Example:
//   info, err := provisioner.CreateProject(ctx, &ProjectConfig{
//       ProjectName: "my-app",
//       OrganizationID: "org_123",
//   })
func (p *DockerProvisioner) CreateProject(ctx context.Context, config *ProjectConfig) (*ProjectInfo, error)
```

### Git Workflow

**Branch Strategy:**

- `main` - Production-ready code
- `develop` - Integration branch
- `feature/*` - Feature branches
- `fix/*` - Bug fix branches
- `release/*` - Release preparation

**Commit Messages:**

```
type(scope): subject

body (optional)

footer (optional)

Types: feat, fix, docs, style, refactor, test, chore
Examples:
  feat(api): add backup creation endpoint
  fix(auth): prevent JWT token without expiration
  docs(readme): update installation instructions
```

**Pull Request Requirements:**

- [ ] All tests pass
- [ ] Code coverage maintained or improved
- [ ] Documentation updated
- [ ] No linter warnings
- [ ] Security scan passes
- [ ] Reviewed by at least one maintainer

---

## Dependencies & Prerequisites

### Technology Stack

**Backend:**
- Go 1.21+ (language)
- Gin 1.9+ (HTTP framework)
- pgx/v5 (PostgreSQL driver)
- sqlc (SQL code generation)
- jwt-go (JWT authentication)
- argon2 (password hashing)
- Docker SDK (container management)

**Database:**
- PostgreSQL 14+ (management database)

**Frontend:**
- Supabase Studio v1.24.04 (patched)
- Next.js 13+
- React 18+

**Infrastructure:**
- Docker 20.10+ (container runtime)
- Docker Compose V2 (orchestration)

**Optional:**
- Kubernetes 1.24+ (for K8s provisioner)
- Prometheus (metrics)
- Grafana (dashboards)
- Alertmanager (alerting)

### External Services

**Required:**
- None (fully self-hosted)

**Optional:**
- DNS provider (for automatic DNS updates)
- SMTP server (for email notifications)
- S3-compatible storage (for backup storage)
- HashiCorp Vault (for secrets management)
- Slack (for notifications)

---

## Metrics for Success

### Phase 4 Success Criteria

- [ ] 80%+ code coverage
- [ ] All API endpoints have tests
- [ ] Zero critical security vulnerabilities
- [ ] All linters pass
- [ ] Documentation complete

### Phase 5 Success Criteria

- [ ] JWT tokens have expiration
- [ ] Rate limiting implemented
- [ ] Audit logging functional
- [ ] Security scan passes
- [ ] Secrets encrypted at rest

### Phase 6 Success Criteria

- [ ] Monitoring dashboards live
- [ ] Alerts configured
- [ ] Documentation published
- [ ] CI/CD pipeline working
- [ ] Backup/restore tested

### Phase 7 Success Criteria

- [ ] Backup system functional
- [ ] Quota system implemented
- [ ] Analytics collecting data
- [ ] pg-meta returning real data
- [ ] All stub endpoints implemented

### Production Readiness Checklist

**Before going live:**

- [ ] All Phase 4-6 tasks complete
- [ ] Security audit passed
- [ ] Load testing completed
- [ ] Disaster recovery tested
- [ ] Monitoring and alerting active
- [ ] Documentation reviewed
- [ ] Secrets rotated from defaults
- [ ] HTTPS configured
- [ ] Backups automated
- [ ] Team trained on operations

---

## Risks & Mitigation

### Technical Risks

**Risk 1: Docker Socket Security**
- **Impact:** High - Container escape possible
- **Probability:** Medium
- **Mitigation:** Use Docker API proxy, implement least privilege
- **Owner:** DevOps team

**Risk 2: Database Corruption**
- **Impact:** Critical - Data loss
- **Probability:** Low
- **Mitigation:** Implement backups, test restores, use transactions
- **Owner:** Database team

**Risk 3: Port Exhaustion**
- **Impact:** Medium - Can't create new projects
- **Probability:** Medium
- **Mitigation:** Implement port pooling, add monitoring
- **Owner:** Platform team

**Risk 4: Zero Test Coverage**
- **Impact:** High - Regressions, bugs in production
- **Probability:** High (already occurring)
- **Mitigation:** Phase 4 addresses this - top priority
- **Owner:** Development team

### Operational Risks

**Risk 5: No Monitoring**
- **Impact:** High - Can't detect issues
- **Probability:** High (already occurring)
- **Mitigation:** Phase 6 implements full monitoring
- **Owner:** Operations team

**Risk 6: Insufficient Documentation**
- **Impact:** Medium - Team can't operate system
- **Probability:** Medium
- **Mitigation:** Phase 6 comprehensive documentation
- **Owner:** Documentation team

**Risk 7: Single Point of Failure**
- **Impact:** Critical - System downtime
- **Probability:** Medium
- **Mitigation:** Phase 6 high availability
- **Owner:** Architecture team

---

## Estimated Resource Requirements

### Phase 4 (Testing & Quality)
- **Engineers:** 2 backend developers
- **Duration:** 3-4 weeks
- **Cost Estimate:** $30-40K (contractor rates)

### Phase 5 (Security)
- **Engineers:** 1 security engineer, 1 backend developer
- **Duration:** 2-3 weeks
- **Cost Estimate:** $25-35K

### Phase 6 (Production Readiness)
- **Engineers:** 1 DevOps, 1 backend developer, 1 technical writer
- **Duration:** 3-4 weeks
- **Cost Estimate:** $35-45K

### Phase 7 (Features)
- **Engineers:** 2 backend developers
- **Duration:** 4-6 weeks
- **Cost Estimate:** $50-70K

### Phase 8 (Enterprise)
- **Engineers:** 2 backend developers, 1 DevOps
- **Duration:** 6-8 weeks
- **Cost Estimate:** $75-100K

**Total Estimated Cost:** $215-290K for full implementation

---

## Success Stories & Use Cases

### Target Use Cases

1. **Startup Development Environment**
   - Multiple projects per developer
   - Quick spin-up and tear-down
   - Cost-effective self-hosted alternative

2. **Agency Client Projects**
   - Isolated projects per client
   - White-label branding
   - Centralized management

3. **Enterprise On-Premises**
   - Data sovereignty requirements
   - Compliance (HIPAA, GDPR)
   - Air-gapped deployments

4. **Education & Training**
   - Student project environments
   - Workshop provisioning
   - Temporary project access

### Expected Outcomes

**For Developers:**
- 10x faster project provisioning (vs manual setup)
- Zero DevOps knowledge required
- Familiar Supabase Studio interface

**For Operations:**
- Centralized project management
- Automated backups and monitoring
- Consistent project configuration

**For Business:**
- Reduced infrastructure costs
- Improved developer productivity
- Better resource utilization

---

## Conclusion

SupaManager has successfully completed Phase 3 and is functional for single-server deployments. The Docker provisioning system works reliably and can provision complete Supabase instances.

**Current Strengths:**
- Solid architecture and code structure
- Working provisioning system
- Studio UI integration
- Comprehensive documentation

**Critical Gaps:**
- Zero test coverage (Phase 4)
- Security vulnerabilities (Phase 5)
- No monitoring/alerting (Phase 6)
- Missing features (Phase 7)

**Recommended Priority:**
1. **Phase 4** - Testing must come first before any production use
2. **Phase 5** - Security fixes are critical for any public deployment
3. **Phase 6** - Production readiness enables real-world usage
4. **Phase 7** - Feature completion achieves parity with Supabase Cloud
5. **Phase 8** - Enterprise features expand market opportunity

**Timeline to Production:**
- **Minimal Viable Production (MVP):** 8-10 weeks (Phases 4-6)
- **Feature Complete:** 12-16 weeks (through Phase 7)
- **Enterprise Ready:** 18-25 weeks (through Phase 8)

The project is well-architected and positioned for success. With focused execution on testing, security, and production readiness, SupaManager can become a premier self-hosted Supabase management platform.

---

**Document Version:** 1.0
**Author:** Claude Code (AI Assistant)
**Date:** November 16, 2025
**Next Review:** Upon Phase 4 completion
