# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Repository Is

This is a **planning and specification repository** for CareerHubV2 — a unified job board for UOW Malaysia students, alumni, SAC staff, and system admins. No application code exists here yet. The repo holds the implementation plans, database schema, API contracts, UI designs, and architecture docs that two developers (one backend, one frontend) will build from.

The actual implementation will live in separate repositories. When writing code for this project, follow the specs below exactly.

---

## Local Development Setup

### Start the SQL Server database (Docker required)
```sh
docker-compose up -d
```
The container is named `careerhub_sqlserver` and exposes SQL Server on `localhost:1433` with SA password `YourStrong!Passw0rd`.

### Backend (Go)
```sh
cd backend
cp .env.example .env        # Fill in real values
go mod tidy
go run ./cmd/server/main.go # Port from .env (default: 8080)
```

### Frontend (Next.js)
```sh
cd frontend
npm install
npm run dev                 # Dev server on localhost:3000
npm run build               # Produces static export in out/
```
The frontend is configured for `output: 'export'` — no Node.js runtime in production. All data fetching runs client-side via Axios.

### Database initialization
Apply the DDL in `docs/init.sql` against the SQL Server instance to create all tables and seed categories.

---

## Architecture

### Tech Stack
| Layer | Technology |
|---|---|
| Backend | Go 1.21+, Gin Gonic, `database/sql` + `go-mssqldb` |
| Frontend | Next.js (App Router, static export), TypeScript, Tailwind CSS, Heroicons |
| Database | Microsoft SQL Server 2022 |
| Auth | Hybrid: Microsoft Entra ID (MSAL) + local bcrypt/JWT for Alumni |
| Deployment | IIS with HttpPlatformHandler (backend) and static files from `out/` (frontend) |

### Authentication Flow
The backend is the sole JWT authority. All users — regardless of login method — receive a CareerHub local JWT (Access: 15 min, Refresh: 7 days):

1. **Entra ID path** (Students, Staff, Admins): Frontend sends Microsoft `id_token` to `POST /api/v1/auth/login/microsoft`. Backend validates JWKS signature, upserts user, issues local JWT.
2. **Local path** (Alumni): `POST /api/v1/auth/login` with email + password. Backend verifies bcrypt hash, issues local JWT.

Auto-provisioning rules on first Entra ID login:
- `@student.uow.edu.my` → `UserType=Student`, auto-assigned `Student` role
- `@uow.edu.my` → `UserType=Staff`, no default role
- other domains → `UserType=External`, `RegistrationStatus=Pending`
- email matches `BOOTSTRAP_ADMIN_EMAIL` env var → `UserType=SystemAdmin`, auto-assigned `System Admin` role

See `docs/auth-process.md` for the full sequence diagram.

### RBAC
Five roles: **System Admin**, **SAC Department**, **Student**, **Alumni**, **External**. Permission strings (e.g. `manage_jobs`, `approve_alumni`, `view_jobs`) are checked by a `RequirePermission("...")` middleware. Full matrix in `docs/roles-permissions-matrix.md`.

### Database Schema
Full table definitions in `docs/db-schema.md` and `docs/init.sql`. Key tables:
- `Users` — single table for all user types; `MicrosoftObjectID` null for Alumni, `PasswordHash` null for Entra ID users
- `Jobs` — content split into `JobDescription`, `Responsibilities`, `Requirements`, `HowToApply` to map to the four frontend tabs (Overview / Requirements / Details / Apply)
- `Roles`, `Permissions`, `RolePermissions`, `UserRoles` — RBAC junction tables
- `AuditLogs` — required for all mutating admin actions
- `PasswordResets` — short-lived hashed tokens for Alumni password setup and recovery

---

## Key Contracts

### API Specification
`shared/openapi.yaml` is the **source of truth** for all API shapes. Backend must implement exactly these endpoints; frontend must consume exactly these responses. Never diverge from it without updating it first.

### API Base Path
All endpoints are under `/api/v1/`.

### Error Response Shape
```json
{ "error": "human message", "code": 401, "error_code": "TOKEN_EXPIRED" }
```
Validation errors (422) include a `details` array with `field` + `message` per failing field.

### Job Filtering (`GET /api/v1/jobs`)
Accepts `page`, `limit`, `search`, `category_id`, `job_type_id`. Always filters out `IsActive=0` and expired jobs for student requests. Returns empty `data: []` with `200 OK` (never 404) so the frontend can render "No results" screens.

---

## Implementation Plans

| Area | Document |
|---|---|
| Backend full plan | `backend/PLAN.md` |
| Backend Sprint 1 detail | `backend/sprint1/plan.md` |
| Frontend full plan | `frontend/PLAN.md` |
| 5-week milestone schedule | `docs/weekly_milestone.md` |

---

## UI Design Reference

High-fidelity PNG designs for every screen live in `frontend/design/`. Sub-folders:
- `Careerhub-login-design/` — Sign-in page and 401 Unauthorized state
- `Careerhub-student-portal-design/` — All student views including filter drawers and tabbed job detail
- `Student Portal Design v2/` — Revised student designs (use these over v1 where they conflict)
- `Careerhub-SAC-portal/` — SAC Department dashboard
- `Careerhub-SystemAdmin-portal/` — System Admin dashboard

**Branding**: Primary Navy `#001C41`, Primary Red `#E21F26`. Font: Montserrat. All layouts are mobile-first.

---

## Backend Testing

There is no Node.js runtime in production, so backend endpoints are tested in isolation with Postman or cURL. To test Entra ID token exchange during development, generate a real `id_token` from [oidcdebugger.com](https://oidcdebugger.com) using the registered `MICROSOFT_CLIENT_ID` and `MICROSOFT_TENANT_ID`, then POST it to `http://localhost:8080/api/v1/auth/login/microsoft`.

Frontend uses **MSW (Mock Service Worker)** during parallel development before the backend is ready.

---

## Deployment

- **Backend**: Compile Go binary, run as Windows Service under IIS `HttpPlatformHandler`. All config from environment variables (see `backend/.env.example`).
- **Frontend**: `npm run build` → copy `out/` to IIS root. Apply `web.config` SPA rewrite rule (detailed in `docs/iis-deployment.md` and `frontend/iis-deployment.md`) to redirect all paths to `index.html`.
