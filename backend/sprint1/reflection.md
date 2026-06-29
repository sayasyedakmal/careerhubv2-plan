# Sprint 1 Reflection — Backend Development
**Date:** 29 June 2026  
**Scope:** Go backend bootstrapping, authentication, and student dashboard APIs  
**Outcome:** All 18 tasks completed. Server running, all endpoints tested and verified.

---

## 1. What We Set Out to Do

Sprint 1 targeted the minimum backend needed to support the student dashboard design (`Namecard_Design B.png`):

- A working Go REST API (Gin) connected to SQL Server
- Microsoft Entra ID token exchange and auto-provisioning
- Stateless JWT session management (access + refresh tokens)
- Four protected data endpoints: `/users/me`, `/categories`, `/jobs`, `/jobs/:id`
- Correct HTTP error codes and structured error response shapes throughout

Everything was pre-planned in `plan.md`, `api-specification.md`, and `plan-review.md` before a single line of code was written. The spec review (`plan-review.md`) had already resolved all ambiguities and contradictions, which made implementation straightforward with no backtracking on design decisions.

---

## 2. What We Built

### Project Structure
```
backend/
├── main.go                     # Gin engine, all routes, middleware chain
├── config/
│   └── database.go             # SQL Server connection pool (key-value DSN)
├── middleware/
│   ├── cors.go                 # CORS for localhost:3000
│   ├── ratelimit.go            # Per-IP token bucket (30/60 req/min)
│   └── auth.go                 # CareerHub JWT validation, context injection
├── internal/
│   └── jwks/
│       └── resolver.go         # Microsoft JWKS public key cache
├── models/
│   ├── user.go                 # User struct
│   ├── category.go             # Category struct with ActiveJobCount
│   └── job.go                  # JobSummary, JobDetail, Pagination, JobsResponse
├── handlers/
│   ├── health.go               # GET /api/v1/health
│   ├── auth.go                 # POST /auth/login/microsoft, POST /auth/refresh
│   ├── user.go                 # GET /users/me
│   ├── category.go             # GET /categories
│   └── job.go                  # GET /jobs, GET /jobs/:id
└── scripts/
    └── seed.sql                # 6 sample Jobs across 4 categories
```

### Dependencies
| Package | Purpose |
|---|---|
| `github.com/gin-gonic/gin` | HTTP router and middleware framework |
| `github.com/microsoft/go-mssqldb` | SQL Server driver |
| `github.com/golang-jwt/jwt/v5` | JWT signing and verification |
| `github.com/joho/godotenv` | `.env` file loader |
| `golang.org/x/time/rate` | Token bucket rate limiter |

### Endpoints Delivered
| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/api/v1/health` | No | Server liveness check |
| POST | `/api/v1/auth/login/microsoft` | No | Entra ID token exchange + auto-provisioning |
| POST | `/api/v1/auth/refresh` | No | Stateless JWT refresh |
| GET | `/api/v1/users/me` | Yes | Authenticated user profile |
| GET | `/api/v1/categories` | Yes | Category list with active job counts |
| GET | `/api/v1/jobs` | Yes | Paginated, filtered job listings |
| GET | `/api/v1/jobs/:id` | Yes | Full job detail (all 4 tabs) |

---

## 3. Issues Encountered & How We Resolved Them

### Issue 1 — `seed.sql` duplicate insert errors in SQL Server Management Studio
**What happened:** The first version of `seed.sql` tried to insert Categories, JobTypes, Roles, Permissions, and RolePermissions — but `docs/init.sql` had already seeded all of that. Running `seed.sql` produced unique constraint violations and primary key conflicts on every lookup table.

**Root cause:** `seed.sql` was written before realising that `docs/init.sql` was a combined DDL + seed script, not just DDL.

**Resolution:** Rewrote `seed.sql` to contain only the 6 sample `Jobs` INSERT statements. All lookup data (Categories, JobTypes, Roles, Permissions, RolePermissions) is the sole responsibility of `docs/init.sql`.

**Lesson:** Always check what an existing init script does before writing a separate seed script. `init.sql` was a combined DDL + seed — not a pure schema migration.

---

### Issue 2 — SQL Server connection refused on port 1433
**What happened:** Server started but immediately exited with:
```
Failed to reach SQL Server: dial tcp [::1]:1433: connectex: No connection could be made because the target machine actively refused it.
```

**Root cause:** The database is a **named SQL Server instance** (`localhost\careerhubv2`). Named instances do not use port 1433 by default — that port is reserved for the default SQL Server instance. Named instances are assigned a dynamic port at startup, discovered via the SQL Server Browser service.

**Resolution:** Opened SQL Server Configuration Manager, enabled TCP/IP on the named instance, confirmed the actual port it was listening on, and updated the environment accordingly. The Go driver then connected successfully.

**Lesson:** Never assume port 1433 for a named SQL Server instance. Always verify via SQL Server Configuration Manager → SQL Server Network Configuration → Protocols for [instance] → TCP/IP → IP Addresses → IPAll → TCP Port.

---

### Issue 3 — DSN URL format rejected for named instance and special-character password
**What happened:** The database connection string was initially built in URL format:
```
sqlserver://sa:Pa$$w0rd@localhost\careerhubv2:1433?database=careerhubv2
```
This produced: `unable to parse connection string: invalid URL format`

**Root cause:** Two compounding problems:
1. The backslash (`\`) in `localhost\careerhubv2` is not valid in a URL
2. The `$` characters in the password broke URL parsing

**Resolution:** Switched `config/database.go` to the **key-value connection string format**:
```
server=localhost\careerhubv2;user id=sa;password=Pa$$w0rd;database=careerhubv2
```
Added logic to omit the port entirely when the host contains a backslash (named instance), letting the driver discover the port via SQL Server Browser.

**Lesson:** Always use the key-value DSN format for SQL Server in Go when using named instances or passwords with special characters. URL format is too fragile for these cases.

---

### Issue 4 — OIDC Debugger failed with `AADSTS9002325`
**What happened:** When trying to get a test Microsoft ID token via [oidcdebugger.com](https://oidcdebugger.com), the request failed with:
```
AADSTS9002325: Proof Key for Code Exchange is required for cross-origin authorization code redemption.
```

**Root cause:** The Azure App Registration did not have the **implicit grant flow** enabled for ID tokens. OIDC Debugger uses the implicit flow (`response_type=id_token`) to get tokens directly in the browser, but the app was configured to only allow the authorization code + PKCE flow.

**Resolution:** Went to Azure Portal → App Registration → Authentication → enabled **ID tokens** under "Implicit grant and hybrid flows" → saved. OIDC Debugger then worked immediately.

**Important note for production:** This implicit grant setting is only needed for OIDC Debugger testing. The actual Next.js frontend will use MSAL (`@azure/msal-react`), which uses the auth code + PKCE flow internally and does **not** require implicit grant to be enabled. This setting can be reviewed and potentially disabled once frontend integration is live.

---

## 4. Key Design Decisions Confirmed in Practice

### Bootstrap Admin worked as designed
On first login with `akmal.o@uow.edu.my` (matching `BOOTSTRAP_ADMIN_EMAIL` in `.env`), the backend correctly:
- Set `UserType = 'SystemAdmin'`
- Auto-assigned the `System Admin` role in `UserRoles`
- Embedded `"role": "System Admin"` in the returned JWT

The access token decoded to:
```json
{
  "sub": 1,
  "email": "akmal.o@uow.edu.my",
  "display_name": "Syed Akmal bin Syed Othman",
  "role": "System Admin",
  "exp": 1782724871,
  "iat": 1782721271
}
```
This confirmed the bootstrap admin override takes priority over domain rules, as specified in `plan-review.md` Section 1B.

### Stateless refresh tokens
Refresh tokens are signed JWTs (HS256) with a 7-day expiry. No database table or session storage required. On refresh, the handler re-queries the DB to get the user's current role so any role changes made by an admin take effect on the next token refresh — a deliberate design choice documented in `plan-review.md` Section 3B.

### Named instance DSN pattern
The final connection string pattern for named instances:
```go
// Named instance — omit port, let Browser service discover it
dsn = fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", host, user, password, dbName)

// Default instance — include port
dsn = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s", host, user, password, port, dbName)
```
Detected automatically by checking if `DB_HOST` contains a backslash.

---

## 5. Test Results Summary

All endpoints tested via Postman after getting a real Entra ID `id_token` from OIDC Debugger.

| Test | Expected | Result |
|---|---|---|
| `POST /auth/login/microsoft` with valid token | `200` + token pair | ✅ Pass |
| `POST /auth/login/microsoft` with no body | `400 INVALID_PAYLOAD` | ✅ Pass |
| `GET /users/me` with valid token | `200` + profile | ✅ Pass |
| `GET /categories?limit=3&sort=job_count` | `200` + top 3 categories | ✅ Pass |
| `GET /jobs?limit=5` | `200` + paginated envelope | ✅ Pass |
| `GET /jobs/1` | `200` + full job detail | ✅ Pass |
| `GET /jobs?search=nonexistent` | `200` with `data: []` | ✅ Pass |
| `GET /jobs/abc` | `400 INVALID_PATH_PARAMETER` | ✅ Pass |
| `GET /categories?sort=invalid` | `400 INVALID_QUERY_PARAMETER` | ✅ Pass |
| `GET /users/me` with no token | `401 MISSING_TOKEN` | ✅ Pass |

---

## 6. What's Next — Sprint 2 (Week 2 per weekly_milestone.md)

### Backend
- `POST /api/v1/auth/login` — local credential login for Alumni (bcrypt password verification)
- `POST /api/v1/auth/register` — Alumni registration form submission (creates `Pending` account)
- RBAC authorization middleware (`RequirePermission("...")`) — checks permission strings from `docs/roles-permissions-matrix.md`
- Token bucket rate limiting already done — verify it covers the new auth endpoints

### Frontend (parallel)
- Dual-login homepage (`Homepage - sign in.png`)
- MSAL React provider setup
- `AuthContext` for JWT storage and claims
- Homepage 401 Unauthorized state

### Integration reminder
The `MICROSOFT_ENTRA_JWKS_URL` in `.env` uses the tenant-specific endpoint (not `/common/`), which is correct per the contradiction resolved in `plan-review.md` Section 4A. Do not change this.
