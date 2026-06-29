# Sprint 1 Implementation Plan: Auth & Student Dashboard APIs

This document outlines the step-by-step development plan for **Sprint 1** of the CareerHubV2 Go backend. 

The primary objective of Sprint 1 is to implement the core authentication services and backend APIs necessary to fully support the student dashboard page depicted in [(Mobile-view) Namecard_Design B.png](file:///C:/repo/personal/careerhubv2-plan/frontend/design/Careerhub-student-portal-design/(Mobile-view)%20Namecard_Design%20B.png).

---

## 1. Objectives & Deliverables
* **Authentication Flow:** Support MSAL (Microsoft Authentication Library) frontend integration using Entra ID token validation and silent token refresh. Local credentials login and registration are out of scope (all authentication is routed through Microsoft).
* **Student Dashboard APIs:** Deliver user profile info, available categories, and latest job list feeds to populate the dashboard.
* **Database Setup:** Define SQL Server schema tables needed for users, categories, and jobs.

### 1.1 Tech Stack & Dependencies (Sprint 1)
* **Go Module Name:** `careerhubv2-backend`
* **Web Framework:** Gin Gonic (`github.com/gin-gonic/gin`)
* **Database Library:** Go Standard Library `database/sql` using Microsoft's official driver (`github.com/microsoft/go-mssqldb`).
* **JWT Management:** `github.com/golang-jwt/jwt/v5`
* **Environment Configuration:** `github.com/joho/godotenv`

---

## 2. Database Schema (Sprint 1 Scope)

We will provision the following tables in Microsoft SQL Server to support Sprint 1 requirements, aligned with [init.sql](file:///c:/repo/personal/careerhubv2-plan/docs/init.sql):

```sql
-- 1. Categories Table
CREATE TABLE Categories (
    CategoryID INT IDENTITY(1,1) PRIMARY KEY,
    CategoryName NVARCHAR(100) NOT NULL,
    IconName NVARCHAR(100) NULL
);

-- 2. JobTypes Table
CREATE TABLE JobTypes (
    JobTypeID INT IDENTITY(1,1) PRIMARY KEY,
    TypeName NVARCHAR(50) NOT NULL
);

-- 3. Users Table (Supports both local & Entra ID auth)
CREATE TABLE Users (
    UserID INT IDENTITY(1,1) PRIMARY KEY,
    MicrosoftObjectID NVARCHAR(100) UNIQUE NULL, -- Null for Alumni
    Email NVARCHAR(255) UNIQUE NOT NULL,
    PasswordHash NVARCHAR(255) NULL,            -- Null for Entra ID SSO users
    DisplayName NVARCHAR(255) NULL,
    Phone NVARCHAR(50) NULL,
    UserType NVARCHAR(50) NOT NULL,              -- 'Student', 'Alumni', 'Staff', 'SystemAdmin', 'External'
    StudentID NVARCHAR(50) UNIQUE NULL,          -- Null for Staff/Admin/External
    RegistrationStatus NVARCHAR(50) DEFAULT 'N/A', -- 'Pending', 'Approved', 'Denied', 'N/A'
    CreatedAt DATETIME DEFAULT GETDATE(),
    LastLoginAt DATETIME NULL
);

-- 4. PasswordResets Table
CREATE TABLE PasswordResets (
    ResetID INT IDENTITY(1,1) PRIMARY KEY,
    UserID INT NOT NULL,
    TokenHash NVARCHAR(255) NOT NULL,
    ExpiresAt DATETIME NOT NULL,
    CONSTRAINT FK_PasswordResets_Users FOREIGN KEY (UserID) REFERENCES Users(UserID) ON DELETE CASCADE
);

-- 5. Jobs Table
CREATE TABLE Jobs (
    JobID INT IDENTITY(1,1) PRIMARY KEY,
    JobTitle NVARCHAR(255) NOT NULL,
    CompanyName NVARCHAR(255) NOT NULL,
    CategoryID INT NOT NULL,
    JobTypeID INT NOT NULL,
    IsActive BIT DEFAULT 1,
    HasSalary BIT DEFAULT 0,
    SalaryMin DECIMAL(18, 2) NULL,
    SalaryMax DECIMAL(18, 2) NULL,
    DeadlineAt DATETIME NULL,
    State NVARCHAR(255) NULL,
    City NVARCHAR(255) NULL,
    PositionCount INT DEFAULT 1,
    JobDescription NVARCHAR(MAX) NULL,     -- Matches Overview Tab detail
    Responsibilities NVARCHAR(MAX) NULL,   -- Matches Details Tab narrative
    Requirements NVARCHAR(MAX) NULL,       -- Matches Requirements Tab checklist
    AdditionalInformation NVARCHAR(MAX) NULL,
    HowToApply NVARCHAR(MAX) NULL,          -- Matches Apply Tab instructions
    CreatedAt DATETIME DEFAULT GETDATE(),
    CreatedBy INT NULL,
    UpdatedAt DATETIME DEFAULT GETDATE(),
    UpdatedBy INT NULL,
    
    CONSTRAINT FK_Jobs_Categories FOREIGN KEY (CategoryID) REFERENCES Categories(CategoryID),
    CONSTRAINT FK_Jobs_JobTypes FOREIGN KEY (JobTypeID) REFERENCES JobTypes(JobTypeID),
    CONSTRAINT FK_Jobs_CreatedBy FOREIGN KEY (CreatedBy) REFERENCES Users(UserID),
    CONSTRAINT FK_Jobs_UpdatedBy FOREIGN KEY (UpdatedBy) REFERENCES Users(UserID)
);

-- 6. Roles Table
CREATE TABLE Roles (
    RoleID INT IDENTITY(1,1) PRIMARY KEY,
    RoleName NVARCHAR(100) UNIQUE NOT NULL
);

-- 7. Permissions Table
CREATE TABLE Permissions (
    PermissionID INT IDENTITY(1,1) PRIMARY KEY,
    PermissionName NVARCHAR(100) UNIQUE NOT NULL
);

-- 8. RolePermissions Junction Table
CREATE TABLE RolePermissions (
    RoleID INT NOT NULL,
    PermissionID INT NOT NULL,
    CONSTRAINT PK_RolePermissions PRIMARY KEY (RoleID, PermissionID),
    CONSTRAINT FK_RolePermissions_Roles FOREIGN KEY (RoleID) REFERENCES Roles(RoleID) ON DELETE CASCADE,
    CONSTRAINT FK_RolePermissions_Permissions FOREIGN KEY (PermissionID) REFERENCES Permissions(PermissionID) ON DELETE CASCADE
);

-- 9. UserRoles Junction Table
CREATE TABLE UserRoles (
    UserID INT NOT NULL,
    RoleID INT NOT NULL,
    CONSTRAINT PK_UserRoles PRIMARY KEY (UserID, RoleID),
    CONSTRAINT FK_UserRoles_Users FOREIGN KEY (UserID) REFERENCES Users(UserID) ON DELETE CASCADE,
    CONSTRAINT FK_UserRoles_Roles FOREIGN KEY (RoleID) REFERENCES Roles(RoleID) ON DELETE CASCADE
);

-- 10. Performance Indexes
-- For counting active jobs per category
CREATE INDEX IX_Jobs_CategoryID_IsActive ON Jobs(CategoryID, IsActive);

-- For fetching latest active jobs efficiently (homepage feeds)
CREATE INDEX IX_Jobs_CreatedAt_IsActive ON Jobs(CreatedAt DESC, IsActive);
```

---

## 3. Sprint 1 Endpoints Details

### Group A: Authentication APIs

#### 1. Exchange Microsoft Token
* **Endpoint:** `POST /api/v1/auth/login/microsoft`
* **Input Payload:**
  ```json
  {
    "id_token": "eyJhbGciOi..." // JWT from Microsoft Entra ID via frontend MSAL
  }
  ```
* **Required Configuration (Environment Variables):**
  * `MICROSOFT_ENTRA_CLIENT_ID`: Frontend's App Application ID registered in Azure AD.
  * `MICROSOFT_ENTRA_TENANT_ID`: Directory/Tenant ID restricting login to the university's tenant.
  * `MICROSOFT_ENTRA_JWKS_URL`: Microsoft public keys URL (`https://login.microsoftonline.com/{tenant_id}/discovery/v2.0/keys`).
* **Backend Processing Pipeline:**
  1. **Token Parse & Header Check:** Parse the `id_token` (JWT format) without verification to read the header and extract the Key ID (`kid`).
  2. **JWKS Key Resolution:** Resolve the corresponding public key from the cached JWKS (JSON Web Key Set) endpoint. If missing, fetch from `MICROSOFT_ENTRA_JWKS_URL` and update the cache.
  3. **Signature Verification:** Verify the token's cryptographic signature against the public key.
  4. **Claims Validation:** Validate:
     * Audience (`aud`): Must match `MICROSOFT_ENTRA_CLIENT_ID`.
     * Issuer (`iss`): Must match `https://login.microsoftonline.com/{tenant_id}/v2.0` (or `https://sts.windows.net/{tenant_id}/`).
     * Tenant ID (`tid`): Must match `MICROSOFT_ENTRA_TENANT_ID`. (If mismatch, return `403 Forbidden` with error code `INVALID_TENANT`).
     * Expiration (`exp`): Current time must be before expiration time.
     * Not Before (`nbf`): Current time must be after active time.
  5. **User Registration / Mapping & Profile Sync:** Extract claims: `email`, `name` (or `preferred_username`), and `oid` (unique MS object ID). Query `Users` database table:
     * **If User exists:** Fetch user record. Check if incoming claims (`DisplayName`, `Email`) differ from the DB record. If changes are detected, sync and update the database record automatically.
     * **If User does NOT exist:** Auto-register user based on domain rules:
       * **Student Domain (`@student.uow.edu.my`):** Set `UserType = 'Student'`, `RegistrationStatus = 'N/A'`, extract `StudentID` from email username prefix (e.g., `99887766`), and assign the `Student` role in `UserRoles`.
       * **Staff Domain (`@uow.edu.my`):** Set `UserType = 'Staff'`, `RegistrationStatus = 'N/A'`, and assign **no default roles** (roles assigned manually by Admin).
       * **Other Domains:** Set `UserType = 'External'`, `RegistrationStatus = 'Pending'`, and assign **no default roles**.
     * **Bootstrap Admin Check:** If the email matches the environment variable `BOOTSTRAP_ADMIN_EMAIL`, set `UserType = 'SystemAdmin'`, `RegistrationStatus = 'N/A'`, and automatically assign the `System Admin` role.
  6. **Token Generation:** Generate and return CareerHub custom JWT tokens (Access token and Refresh token). Refresh tokens are stateless JWTs.
* **Success Response (200 OK):**
  ```json
  {
    "access_token": "careerhub_jwt_access_token...",
    "refresh_token": "careerhub_jwt_refresh_token...",
    "expires_in": 3600
  }
  ```
* **Error Responses:**
  * `401 Unauthorized`: `{"error_code": "INVALID_MICROSOFT_TOKEN", "message": "Microsoft ID token verification failed or token expired"}`
  * `403 Forbidden`: `{"error_code": "INVALID_TENANT", "message": "Authenticated Microsoft tenant is not authorized for this application"}`
  * `500 Internal Server Error`: `{"error_code": "DATABASE_ERROR", "message": "Internal service error while processing user login"}`

#### 2. Local Credentials Login (Out of Scope)
* **Endpoint:** `POST /api/v1/auth/login` (Not implemented in Sprint 1. All users authenticate via Microsoft MSAL.)

#### 3. Refresh Access Token
* **Endpoint:** `POST /api/v1/auth/refresh`
* **Input Payload:**
  ```json
  {
    "refresh_token": "careerhub_jwt_refresh_token..."
  }
  ```
* **Success Response (200 OK):** Returns new `access_token` and the current/valid `refresh_token`.
* **Error Response:**
  * `401 Unauthorized`: `{"error_code": "INVALID_REFRESH_TOKEN", "message": "Refresh token is invalid or expired. Please sign in again."}`

---

### Group B: Student Dashboard APIs

#### 4. Get Current User Profile
* **Endpoint:** `GET /api/v1/users/me`
* **Authorization:** `Bearer <JWT_ACCESS_TOKEN>`
* **Usage:** Returns details to show student name in the welcome banner.
* **Success Response (200 OK):**
  ```json
  {
    "id": 1,
    "email": "student@uow.edu.my",
    "display_name": "Sarah Connor",
    "phone": "+60129998877",
    "user_type": "Student",
    "student_id": "99887766"
  }
  ```

#### 5. Fetch Dashboard Categories
* **Endpoint:** `GET /api/v1/categories`
* **Authorization:** `Bearer <JWT_ACCESS_TOKEN>`
* **Usage:** Feeds the *"Top 3 Job Categories"* navigation block. Supports sorting by job counts to fetch top categories.
* **Query Parameters:**
  * `limit` (int, optional) - Limit the number of categories returned (e.g. `limit=3`).
  * `sort` (string, optional) - Set to `job_count` to sort by categories with the most active and unexpired jobs first (`IsActive = 1` AND (`DeadlineAt IS NULL` OR `DeadlineAt > GETDATE()`)).
* **Success Response (200 OK):**
  ```json
  [
    { "id": 1, "name": "Computing & IT", "icon_name": "CommandLineIcon", "active_job_count": 15 },
    { "id": 2, "name": "Engineering", "icon_name": "WrenchIcon", "active_job_count": 9 },
    { "id": 3, "name": "Hospitality", "icon_name": "BriefcaseIcon", "active_job_count": 4 }
  ]
  ```

#### 6. Fetch Jobs (Filtered & Paginated)
* **Endpoint:** `GET /api/v1/jobs`
* **Authorization:** `Bearer <JWT_ACCESS_TOKEN>`
* **Usage:** 
  * Feeds the *"Latest Jobs"* carousel (e.g. via `GET /api/v1/jobs?limit=5`).
  * Powers the search bar input and advanced filters drawer.
* **Filtering Logic:** Automatically filters out inactive jobs (`IsActive = 0`) and expired jobs (`DeadlineAt <= GETDATE()`) by default for student requests.
* **Query Parameters:** `page`, `limit`, `search`, `category_id`, `job_type_id`.
* **Success Response (200 OK):** Returns a standardized envelope containing `data` list and `pagination` metadata:
  ```json
  {
    "data": [
      {
        "id": 101,
        "job_title": "Event Executive (Skincare Industry)",
        "company_name": "LAVIN PHARMA (M) SDN BHD",
        "category_id": 3,
        "category_name": "Hospitality & Tourism",
        "job_type": "Internship",
        "state": "Selangor",
        "city": "Shah Alam",
        "salary_min": 700.00,
        "salary_max": 1200.00,
        "deadline": "2026-07-31T00:00:00Z",
        "is_active": true,
        "created_at": "2026-06-15T15:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total_records": 45,
      "total_pages": 5
    }
  }
  ```

#### 7. Fetch Single Job Details
* **Endpoint:** `GET /api/v1/jobs/{id}`
* **Authorization:** `Bearer <JWT_ACCESS_TOKEN>`
* **Usage:** Triggers when clicking **"Find out more"** to fetch segmented tabs details.
* **Success Response (200 OK):** Returns detailed fields like `job_description`, `responsibilities`, `requirements`, `additional_information`, and `how_to_apply`.

---

## 4. Implementation Steps (Timeline)

1. **Step 1: DB Migration Script**  
   Draft and execute the SQL Server DDL migrations to set up the 5 core tables. Set up initial mock data for Categories and Jobs to support dashboard previews.
2. **Step 2: Security & JWT Verification Middleware**  
   Implement JWT validation and parsing using `golang-jwt`. Write standard auth validation checks (HTTP `401` handling) and role check checks (HTTP `403` handling).
3. **Step 3: Microsoft Entra ID JWKS Validator**  
   Implement a helper function to fetch and cache Microsoft's public keys from `https://login.microsoftonline.com/{tenant_id}/discovery/v2.0/keys` to securely verify Microsoft Entra ID tokens.
4. **Step 4: Auth Controllers Implementation**  
   Write Gin handlers for `/auth/login/microsoft` and `/auth/refresh`. (Local login/registration endpoints `/auth/login` and `/auth/register` are out of scope).
5. **Step 5: Dashboard Content Controllers**  
   Write Gin handlers for `/users/me`, `/categories`, and `/jobs` (with pagination, keywords, and filters).
6. **Step 6: Endpoint Testing & Validation**  
   Test endpoints locally with Postman or Go test cases, checking success payloads and HTTP custom error codes (401, 422, 429, 500).

---

## 5. Testing Strategy (Isolated Backend Testing)

Since the frontend is being developed in parallel, the backend endpoints must be tested in isolation. We will use **Strategy 2 (OIDC Debugger)** to test with real Microsoft tokens.

### A. Testing the MSAL Exchange (`POST /api/v1/auth/login/microsoft`)

To test the Microsoft ID Token validation logic with a real Entra ID token:

1. **Generate ID Token via OIDC Debugger:**
   * Go to [oidcdebugger.com](https://oidcdebugger.com/).
   * Enter the following details:
     * **Authorize URI:** `https://login.microsoftonline.com/{MICROSOFT_ENTRA_TENANT_ID}/oauth2/v2.0/authorize` *(Replace `{MICROSOFT_ENTRA_TENANT_ID}` with your active tenant ID)*
     * **Redirect URI:** `https://oidcdebugger.com/redirect`
     * **Client ID:** Your `MICROSOFT_ENTRA_CLIENT_ID`
     * **Scope:** `openid profile email`
     * **Response type:** `id_token`
     * **Response mode:** `form_post`
   * Click **Send Request** and authenticate using your Microsoft/University account.
   * Copy the generated token from the results screen.

2. **Send Exchange Request:**
   * Using an API client (like Postman or cURL), send the token to the local server:
     ```bash
     POST http://localhost:8085/api/v1/auth/login/microsoft
     Content-Type: application/json

     {
       "id_token": "<PASTE_COPIED_TOKEN>"
     }
     ```
   * **Expected Result:** The backend successfully verifies the token signature/claims, registers the user if they are new, and returns a JSON payload containing the `access_token` and `refresh_token`.

### B. Testing Protected Dashboard Endpoints

To test endpoints requiring authentication (such as `GET /api/v1/users/me` or `GET /api/v1/jobs`):

1. **Extract Access Token:**
   * Copy the `access_token` returned by the successful MSAL exchange endpoint.
2. **Authorize Requests:**
   * In your API client, add an authorization header to your request:
     ```http
     Authorization: Bearer <PASTE_ACCESS_TOKEN>
     ```
3. **Execute Requests:**
   * Query protected endpoints (e.g. `GET http://localhost:8085/api/v1/jobs?limit=5`) to verify authorization checks, data payload formats, and correct HTTP response codes.
