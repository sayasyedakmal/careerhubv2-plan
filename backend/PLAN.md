# Backend Implementation Plan: CareerHubV2

This document outlines the Go-based REST API implementation, updated to fully support the high-fidelity designs received from the frontend team.

## 1. Tech Stack
*   **Language**: Go 1.21+
*   **Framework**: [Gin Gonic](https://gin-gonic.com/) (Web)
*   **Database**: Microsoft SQL Server
*   **Authentication**: Hybrid (Microsoft Entra ID JWT Validation + Local bcrypt/JWT)
*   **Deployment**: IIS with `HttpPlatformHandler`

---

## 2. Database Design & API Entities (SQL Server)
The database schema must support the rich, responsive attributes defined in the design documents:
*   **Jobs (Tab-structured Content)**:
    *   Instead of a single text field, store data in fields aligning with the designed Job Details tabs:
        *   `overview_text` (markdown or text for key characteristics, salary, location, timeline)
        *   `requirements_text` (markdown list of qualifications, experience level)
        *   `details_text` (full narrative, company profile, core duties)
    *   `category_id` linked to a `Categories` table to populate the new **All Categories** navigation structure.
*   **Categories**:
    *   `id` (PK)
    *   `name` (Computing, Engineering, Business, etc.)
    *   `icon_name` (string mapping to Heroicons)
*   **Users & Profiles**:
    *   Support dynamic details (Full Name, Student ID, Email, Phone Number, Password) linked to a digital badge or Student Namecard profile component.
*   **AuditLogs**: Audit trail tracking administrative/approval operations.

---

## 3. Advanced Filtering & Search API
To power the multi-step filter pages (`Filter feature` steps 1-3) on the frontend:
*   `GET /api/v1/jobs` must accept query parameters reflecting the advanced filter drawer:
    *   `category_id`: Filter by job categories.
    *   `job_type`: e.g., Full-time, Part-time, Internship.
    *   `location`: e.g., Onsite, Remote, Hybrid.
    *   `salary_min` / `salary_max`.
*   The handler must dynamically build the SQL query using `OFFSET / FETCH NEXT` for pagination, returning metadata like total record count, and supporting empty arrays (status `200` with empty body) so the frontend can display its design-calibrated "No results" pages.

---

## 4. API Endpoints (v1)

### Auth & Onboarding
*   `POST /api/v1/auth/login/microsoft` (Entra ID flow)
    *   **Auto-Registration & Role Mapping**: When validating an Entra ID token for a new user, upsert their record into the `Users` table. If the email contains the student domain `@uow.edu.my`, automatically assign the `Active Student` role in `UserRoles`.
    *   **Bootstrap Admin Override**: Read `BOOTSTRAP_ADMIN_EMAIL` from OS/environment variables. If the validated Entra ID email matches this variable, automatically assign them the `System Admin` role to avoid cold-start lockouts.
*   `POST /api/v1/auth/login` (Local login)
*   `POST /api/v1/auth/register` (Alumni onboarding request)
*   `POST /api/v1/auth/password-reset` & `POST /api/v1/auth/password-change`

### Categories & Jobs
*   `GET /api/v1/categories` (Returns list of categories to feed the "All Categories" browse page)
*   `GET /api/v1/jobs` (Paginated, filtered list)
*   `GET /api/v1/jobs/:id` (Returns segmented JSON structure matching Job Details tabs: Overview, Requirements, Details)
*   `POST /api/v1/jobs`, `PUT /api/v1/jobs/:id`, `PATCH /api/v1/jobs/:id/status`, `DELETE /api/v1/jobs/:id`

---

## 5. Security, Rate Limiting & Error Handling Middleware
To trigger the high-fidelity frontend custom error states, the backend must return clean, structured error formats with standard HTTP headers:
*   **Authentication (401 Unauthorized)**: Middleware returns standard JSON error bodies when token is missing/expired.
*   **Authorization (403 Forbidden)**: Middleware checks RBAC permissions (refer to `docs/roles-permissions-matrix.md`) and issues `403` standard status.
*   **Rate Limiting (429 Too Many Requests)**: Implement a token bucket or window-based rate limiter middleware using Go libraries (e.g. `golang.org/x/time/rate`). Return standard `429` status code with `Retry-After` headers.
*   **Error Recovery (500 Internal Error)**: Ensure Gin's custom recovery middleware catches panics and formats standard `500` HTTP payloads.
