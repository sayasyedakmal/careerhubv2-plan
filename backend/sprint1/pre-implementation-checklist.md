# Sprint 1: Pre-Implementation Checklist

Please review and resolve the following items in your local setup before we begin writing the Go backend code. You can use this file as a checklist in your editor.

---

## 1. Environment Configuration (`backend/.env`)

Ensure your local `.env` file contains all the following environment variables.

### A. Microsoft Entra ID (MSAL validation)
* [ ] **`AZURE_AD_CLIENT_ID`**: Your temporary testing Client ID from Azure App Registration.
* [ ] **`AZURE_AD_TENANT_ID`**: Your Directory Tenant ID.
* [ ] **`AZURE_AD_JWKS_URL`**: Set to `https://login.microsoftonline.com/{YOUR_TENANT_ID}/discovery/v2.0/keys`.

### B. Database Connection (SQL Server)
* [ ] **`DB_HOST`**: Host address (e.g. `localhost` or remote server IP).
* [ ] **`DB_PORT`**: SQL Server port (default is `1433`).
* [ ] **`DB_USER`**: Database username.
* [ ] **`DB_PASSWORD`**: Database password.
* [ ] **`DB_NAME`**: Target database name (e.g. `CareerHubV2`).
* [ ] *Alternative Option:* **`DB_CONNECTION_STRING`**: Fully formatted connection string (e.g. `sqlserver://username:password@host:port?database=dbname`).

### C. JWT Configuration (Session Management)
* [ ] **`JWT_SECRET`**: A long, random string used by the Go backend to sign session tokens.
* [ ] **`JWT_ACCESS_EXPIRATION_SEC`**: Access token lifetime (default `3600` for 1 hour).
* [ ] **`JWT_REFRESH_EXPIRATION_SEC`**: Refresh token lifetime (e.g. `604800` for 7 days).

### D. Application Settings
* [ ] **`PORT`**: Port the backend runs on (e.g. `8080`).
* [ ] **`ENV`**: Set to `development` for local testing.
* [ ] **`BOOTSTRAP_ADMIN_EMAIL`**: The email of the developer or admin account to automatically receive admin permissions upon first Microsoft login.

---

## 2. Go Setup Decisions

* [x] **Go Module Name:** Decided. Use:
  ```bash
  go mod init careerhubv2-backend
  ```
* [x] **Database Library Preference:** Decided. Use:
  * **Standard library:** `database/sql` using the Microsoft SQL Server driver `github.com/microsoft/go-mssqldb`.

---

## 3. Database Server Readiness

* [ ] **Instance Running:** Confirm Microsoft SQL Server is running locally or remotely.
* [ ] **Network Access:** Ensure SQL Server is configured to allow TCP/IP connections on port `1433`.
* [ ] **Database Created:** Ensure the target database name (e.g., `CareerHubV2`) is created. The backend migration scripts will create tables inside this database.

---

## 4. Recommended Project Directory Structure

Once the checklist is complete, we will organize the codebase as follows:

```
backend/
├── .env                  # Environment configurations
├── main.go               # Application entry point
├── go.mod                # Go module specifications
├── config/
│   └── database.go       # SQL Server connection pool initialization
├── middleware/
│   ├── auth.go           # JWT authorization verification middleware
│   └── cors.go           # CORS configuration for frontend hosting
├── models/
│   ├── user.go           # User database structures & schemas
│   ├── job.go            # Job database structures & schemas
│   └── category.go       # Category database structures & schemas
└── handlers/
    ├── auth.go           # MSAL exchange & token refresh handlers
    ├── user.go           # Profile fetch/update handlers
    ├── job.go            # Paginated job listing & detail handlers
    └── category.go       # Job categories list handlers
```
