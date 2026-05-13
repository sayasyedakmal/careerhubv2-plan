# Backend Implementation Plan: CareerHubV2

This document provides the full technical plan for the Go-based REST API.

## 1. Tech Stack
- **Language**: Go 1.21+
- **Framework**: [Gin Gonic](https://gin-gonic.com/) (Web)
- **Database**: Microsoft SQL Server
- **Authentication**: Microsoft Entra ID (JWT Validation)
- **Deployment**: IIS with `HttpPlatformHandler`

## 2. Database Design (SQL Server)
You are building from scratch. Run `docs/init.sql` to initialize:
- **Jobs**: Stores listings with fields from `list.csv`.
- **Users**: Tracks Microsoft `oid`, Email, and `LastLogin`.
- **Categories/JobTypes**: Normalized lookup tables for filters.

## 3. Core Development Tasks

### Authentication Middleware
- Implement a middleware to intercept `Authorization: Bearer <token>` headers.
- Validate tokens against Microsoft's JWKS endpoint.
- Verify `aud` (Audience) matches your `AZURE_CLIENT_ID`.
- Inject user claims into the Go request context.

### API Endpoints (v1)
Refer to `shared/openapi.yaml` for exact schemas:
- `GET /health`: Public health check.
- `GET /jobs`: Returns filtered/searched list of jobs. (Auth Required)
- `GET /jobs/:id`: Returns full details for a single job. (Auth Required)

### SQL Integration
- Use `database/sql` with the `denisenkom/go-mssqldb` driver.
- Implement the Repository pattern for clean data access.

## 4. Deployment (IIS)
- Build the binary: `go build -o careerhub-api.exe`.
- Use the `web.config` from `docs/iis-deployment.md`.
- Install **HttpPlatformHandler** on the server.
- Set App Pool to **"No Managed Code"**.

---

## 5. How to Integrate with Frontend
To ensure you and the frontend developer are in sync:

1.  **Shared Contract**: Always refer to `shared/openapi.yaml`. If you change an API response, update this file first.
2.  **Environment**: Use `shared/.env.template` to share the `AZURE_CLIENT_ID` and `AZURE_TENANT_ID`.
3.  **CORS**: Ensure your Gin setup allows requests from the frontend origin (e.g., `http://localhost:5173`).
4.  **Token Header**: Expect the token in the `Authorization: Bearer <token>` header.
