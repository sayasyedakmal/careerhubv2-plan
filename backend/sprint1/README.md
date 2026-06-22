# Sprint 1: Quickstart & Execution Guide

Welcome to Sprint 1! This guide outlines the exact roadmap we will follow when you return to start the implementation.

---

## 📂 Sprint 1 Documents Reference

Before you start writing code, make sure you are familiar with the helper files in this directory:
1. **[plan.md](file:///C:/repo/personal/careerhubv2-plan/backend/sprint1/plan.md):** The core database tables, endpoint specs, and implementation steps.
2. **[auth-concepts.md](file:///C:/repo/personal/careerhubv2-plan/backend/sprint1/auth-concepts.md):** Clear explanation of the token design (Microsoft ID Token vs. CareerHub Access/Refresh Tokens).
3. **[entra-setup-guide.md](file:///C:/repo/personal/careerhubv2-plan/backend/sprint1/entra-setup-guide.md):** Step-by-step instructions for registering your temporary test app in Azure to get credentials.
4. **[pre-implementation-checklist.md](file:///C:/repo/personal/careerhubv2-plan/backend/sprint1/pre-implementation-checklist.md):** Verification list for your local `.env` values, SQL Server readiness, and target folder structure.

---

## 🚀 Step-by-Step Execution Plan (For Next Session)

When you start the next session, we will execute these steps one by one:

### Phase 1: Bootstrap & Connectivity (Day 1)
1. **[ ] Initialize Go Module:**  
   Initialize the project directory in `backend/`:
   ```bash
   go mod init careerhubv2-backend
   ```
2. **[ ] Install Core Dependencies:**  
   Download the required library packages:
   ```bash
   go get github.com/gin-gonic/gin
   go get github.com/microsoft/go-mssqldb
   go get github.com/golang-jwt/jwt/v5
   go get github.com/joho/godotenv
   ```
3. **[ ] Run DDL Migrations:**  
   Run the database schema script from [plan.md](file:///C:/repo/personal/careerhubv2-plan/backend/sprint1/plan.md#L16-L78) directly in your SQL Server instance to create the tables.
4. **[ ] Write Database Connection Helper (`config/database.go`):**  
   Implement connection pool validation utilizing `database/sql` and verify the backend can connect to SQL Server.

### Phase 2: Security & JWT Validation (Day 2)
5. **[ ] Build Microsoft JWKS Caching Resolver:**  
   Write a utility function to fetch and cache Microsoft's public keys (`AZURE_AD_JWKS_URL`) for verifying the `id_token` signature.
6. **[ ] Implement Local JWT Middleware (`middleware/auth.go`):**  
   Create middleware to intercept requests, verify the Custom Access Token, and bind the user context.

### Phase 3: Controller Endpoints & Testing (Day 3)
7. **[ ] Write Auth Handlers (`handlers/auth.go`):**  
   Implement `/auth/login/microsoft` (parsing, JWKS validation, auto-registration, custom JWT generation) and `/auth/refresh` (silent token renewal).
8. **[ ] Write Student Dashboard Handlers:**  
   Implement the endpoints:
   * `GET /api/v1/users/me` (Profile detail welcome banner)
   * `GET /api/v1/categories` (Dashboard categories)
   * `GET /api/v1/jobs` (Paginated, filtered jobs listing)
   * `GET /api/v1/jobs/{id}` (Segments tab overview for details page)
9. **[ ] Verify Endpoints:**  
   Use `oidcdebugger.com` to sign in, grab a real Microsoft ID token, and test the endpoints via Postman or cURL.
