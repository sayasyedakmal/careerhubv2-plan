# Sprint 1 Task List

Reference docs in this folder before starting each task:
- `plan.md` — DB schema, endpoint specs, implementation steps
- `api-specification.md` — exact request/response shapes and error codes
- `auth-concepts.md` — token flow diagrams
- `plan-review.md` — resolved ambiguities and contradictions
- `pre-implementation-checklist.md` — folder structure and env vars
- `entra-setup-guide.md` — Azure app registration steps

---

## Phase 1: Bootstrap & Connectivity

### Task 1 — Initialize Go module and project folder structure ✅
- [x] Run `go mod init careerhubv2-backend` inside `backend/`
- [x] Create folders: `config/`, `middleware/`, `models/`, `handlers/`, `internal/jwks/`
- [x] Create `main.go` with Gin engine, `godotenv` loader, CORS, rate limit, auth middleware, and all `/api/v1` routes registered

### Task 2 — Install core Go dependencies ✅
- [x] `go get github.com/gin-gonic/gin`
- [x] `go get github.com/microsoft/go-mssqldb`
- [x] `go get github.com/golang-jwt/jwt/v5`
- [x] `go get github.com/joho/godotenv`
- [x] `go get golang.org/x/time/rate` (rate limiter)
- [x] `go.mod` and `go.sum` updated, `go build ./...` passes with zero errors

### Task 3 — Run DDL migrations on SQL Server ✅
- [x] `docs/init.sql` applied to the `careerhubv2` named instance
- [x] All tables created: `Categories`, `JobTypes`, `Users`, `PasswordResets`, `Jobs`, `Roles`, `Permissions`, `RolePermissions`, `UserRoles`, `AuditLogs`
- [x] `PasswordHash` and `PasswordResets` present in schema but dormant — local auth out of scope for Sprint 1

### Task 4 — Seed Categories, JobTypes, and sample Jobs ✅
- [x] `backend/scripts/seed.sql` run successfully
- [x] 6 sample Jobs inserted across Computing & IT, Engineering, Business & Finance, and Hospitality & Tourism categories

### Task 5 — Implement database connection pool (`config/database.go`) ✅
- [x] Reads `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` from `.env`
- [x] Opens `*sql.DB` using `go-mssqldb` with `sqlserver://` DSN
- [x] `db.Ping()` on startup; fatal exit if unreachable
- [x] Exports `config.DB` for use by all handlers

### Task 6 — Build health check endpoint (`GET /api/v1/health`) ✅
- [x] Route registered in `main.go`, no auth required
- [x] Returns `200 OK` with `{"status": "ok"}`

---

## Phase 2: Security & Middleware

### Task 7 — Implement CORS middleware (`middleware/cors.go`) ✅
- [x] `Access-Control-Allow-Origin: http://localhost:3000`
- [x] `Access-Control-Allow-Methods: GET, POST, OPTIONS`
- [x] `Access-Control-Allow-Headers: Content-Type, Authorization`
- [x] OPTIONS preflight returns `200 OK`
- [x] Registered globally in `main.go`

### Task 8 — Implement rate limiting middleware ✅
- [x] Per-IP token bucket limiter in `middleware/ratelimit.go`
- [x] 30 req/min on `POST /auth/login/microsoft`
- [x] 60 req/min on `POST /auth/refresh`
- [x] Returns `429` with `{"error": "...", "code": "RATE_LIMIT_EXCEEDED"}` and `Retry-After: 60` header
- [x] Stale IP entries cleaned up every minute

### Task 9 — Build Microsoft JWKS caching resolver ✅
- [x] `internal/jwks/resolver.go` — fetches from `MICROSOFT_ENTRA_JWKS_URL`
- [x] In-memory cache; re-fetches only on cache miss
- [x] Uses tenant-specific URL, not `/common/` (per `plan-review.md` Section 4A)
- [x] Returns `*rsa.PublicKey` for a given `kid`

### Task 10 — Implement JWT auth middleware (`middleware/auth.go`) ✅
- [x] Extracts `Authorization: Bearer <token>` header
- [x] Verifies CareerHub access token (HS256) with `JWT_SECRET`
- [x] Injects `userID`, `email`, `display_name`, `role` into Gin context
- [x] `MISSING_TOKEN` (401), `TOKEN_EXPIRED` (401), `INVALID_TOKEN` (401) returned correctly
- [x] Applied to all protected routes; auth/refresh/health are unprotected

### Task 11 — Define Go struct models (`models/`) ✅
- [x] `models/user.go` — `User` struct with `json:` tags; sensitive fields (`PasswordHash`, `MicrosoftObjectID`) excluded from JSON output
- [x] `models/job.go` — `JobSummary` (list), `JobDetail` (full tabs), `Pagination`, `JobsResponse`
- [x] `models/category.go` — `Category` struct with `ActiveJobCount`

---

## Phase 3: Auth Handlers

### Task 12 — Implement `POST /api/v1/auth/login/microsoft` ✅
- [x] Decodes request body; `400 INVALID_PAYLOAD` if malformed
- [x] Extracts `kid` from JWT header without verification
- [x] Resolves RSA public key via JWKS resolver
- [x] Verifies token signature (RS256)
- [x] Validates `aud`, `iss`, `tid` (403 if tenant mismatch), `exp` (401 MICROSOFT_TOKEN_EXPIRED)
- [x] Extracts `oid`, `email` (falls back to `preferred_username`), `name`
- [x] Upserts user with all domain rules and bootstrap admin override
- [x] Syncs `DisplayName`/`Email` on repeat login
- [x] Returns `{access_token, refresh_token, expires_in: 3600}`

### Task 13 — Implement `POST /api/v1/auth/refresh` ✅
- [x] `400 INVALID_PAYLOAD` if body malformed or field missing
- [x] Validates refresh token JWT (HS256) signature and expiration
- [x] `401 REFRESH_TOKEN_EXPIRED` if expired; `401 INVALID_REFRESH_TOKEN` if invalid
- [x] Fetches current role from DB to embed in new access token
- [x] Returns fresh access token + still-valid refresh token (stateless, no rotation — deferred to future sprint)

---

## Phase 4: Dashboard Handlers

### Task 14 — Implement `GET /api/v1/users/me` ✅
- [x] Reads `userID` from Gin context (injected by JWT middleware)
- [x] Queries `Users` table by `UserID`
- [x] Returns `{id, email, display_name, phone, user_type, student_id}`
- [x] `404 RESOURCE_NOT_FOUND` if record deleted mid-session

### Task 15 — Implement `GET /api/v1/categories` ✅
- [x] Validates `limit` (positive int) and `sort` (only `job_count` accepted)
- [x] `400 INVALID_QUERY_PARAMETER` for invalid values
- [x] LEFT JOIN counts only active + unexpired jobs; ORDER BY count DESC when `sort=job_count`
- [x] Returns `[{id, name, icon_name, active_job_count}]`

### Task 16 — Implement `GET /api/v1/jobs` (paginated + filtered) ✅
- [x] Validates `page`, `limit` (default 10, max 100), `category_id`, `job_type_id`
- [x] `400 INVALID_QUERY_PARAMETER` on invalid params
- [x] Always filters `IsActive=1` and non-expired deadlines
- [x] `search` does case-insensitive LIKE on `JobTitle`, `CompanyName`, `CategoryName`
- [x] `OFFSET / FETCH NEXT` pagination; separate COUNT query for totals
- [x] Returns `{data: [...], pagination: {...}}`; empty results return `200` with `data: []`

### Task 17 — Implement `GET /api/v1/jobs/:id` ✅
- [x] Validates `:id` is a positive integer; `400 INVALID_PATH_PARAMETER` if not
- [x] Queries `Jobs` JOIN `Categories` JOIN `JobTypes`
- [x] Returns full detail shape: `job_description`, `responsibilities`, `requirements`, `additional_information`, `how_to_apply`
- [x] `404 RESOURCE_NOT_FOUND` if job does not exist

---

## Phase 5: Testing & Verification

### Task 18 — Test all endpoints with Postman (OIDC Debugger flow) ✅
- [x] Got real `id_token` via OIDC Debugger (had to enable implicit grant / ID tokens in Azure App Registration first)
- [x] `POST /auth/login/microsoft` — returned `access_token` + `refresh_token`; bootstrap admin logic confirmed (`role: System Admin`, `email: akmal.o@uow.edu.my`)
- [x] `GET /users/me` — profile returned correctly
- [x] `GET /categories?limit=3&sort=job_count` — top 3 categories by active job count returned
- [x] `GET /jobs?limit=5` — paginated envelope returned
- [x] `GET /jobs/1` — full job detail with all tab fields returned
- [x] `GET /jobs?search=nonexistent` — `200` with `data: []`
- [x] `GET /jobs/abc` → `400 INVALID_PATH_PARAMETER`
- [x] `GET /categories?sort=invalid` → `400 INVALID_QUERY_PARAMETER`
