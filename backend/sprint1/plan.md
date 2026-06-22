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

We will provision the following tables in Microsoft SQL Server to support Sprint 1 requirements:

```sql
-- 1. Roles Table
CREATE TABLE Roles (
    id INT IDENTITY(1,1) PRIMARY KEY,
    role_name NVARCHAR(50) NOT NULL UNIQUE
);

-- 2. Users Table
CREATE TABLE Users (
    id UNIQUEIDENTIFIER DEFAULT NEWID() PRIMARY KEY,
    email NVARCHAR(256) NOT NULL UNIQUE,
    display_name NVARCHAR(256) NOT NULL,
    phone NVARCHAR(20) NULL,
    user_type NVARCHAR(50) NOT NULL DEFAULT 'ActiveStudent', -- ActiveStudent, Alumni, Staff
    registration_status NVARCHAR(50) NOT NULL DEFAULT 'N/A', -- Pending, Approved, Denied, N/A
    student_id NVARCHAR(50) NULL UNIQUE,
    major NVARCHAR(256) NULL,
    password_hash NVARCHAR(256) NULL, -- Deprecated/Unused (All authentication is handled via Microsoft Entra ID)
    created_at DATETIME2 DEFAULT GETDATE(),
    updated_at DATETIME2 DEFAULT GETDATE()
);

-- 3. UserRoles Mapping Table
CREATE TABLE UserRoles (
    user_id UNIQUEIDENTIFIER FOREIGN KEY REFERENCES Users(id) ON DELETE CASCADE,
    role_id INT FOREIGN KEY REFERENCES Roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

-- 4. Categories Table
CREATE TABLE Categories (
    id INT IDENTITY(1,1) PRIMARY KEY,
    name NVARCHAR(100) NOT NULL UNIQUE,
    icon_name NVARCHAR(100) NOT NULL -- Maps to frontend Heroicons (e.g. "CommandLineIcon")
);

-- 5. Jobs Table
CREATE TABLE Jobs (
    id UNIQUEIDENTIFIER DEFAULT NEWID() PRIMARY KEY,
    job_title NVARCHAR(256) NOT NULL,
    company_name NVARCHAR(256) NOT NULL,
    category_id INT FOREIGN KEY REFERENCES Categories(id),
    job_type NVARCHAR(50) NOT NULL, -- e.g., Full-Time, Part-Time, Internship, Contract
    location NVARCHAR(256) NOT NULL,
    location_type NVARCHAR(50) NOT NULL, -- e.g., Onsite, Remote, Hybrid
    experience_level NVARCHAR(50) NOT NULL, -- e.g., Junior, Mid, Senior
    salary_min DECIMAL(18,2) NULL,
    salary_max DECIMAL(18,2) NULL,
    deadline DATE NOT NULL,
    job_description NVARCHAR(MAX) NULL,
    responsibilities NVARCHAR(MAX) NULL,
    requirements NVARCHAR(MAX) NULL,
    additional_information NVARCHAR(MAX) NULL,
    how_to_apply NVARCHAR(MAX) NULL,
    is_active BIT NOT NULL DEFAULT 1,
    created_at DATETIME2 DEFAULT GETDATE(),
    updated_at DATETIME2 DEFAULT GETDATE()
);
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
  * `AZURE_AD_CLIENT_ID`: Frontend's App Application ID registered in Azure AD.
  * `AZURE_AD_TENANT_ID`: Directory/Tenant ID restricting login to the university's tenant.
  * `AZURE_AD_JWKS_URL`: Microsoft public keys URL (`https://login.microsoftonline.com/{tenant_id}/discovery/v2.0/keys`).
* **Backend Processing Pipeline:**
  1. **Token Parse & Header Check:** Parse the `id_token` (JWT format) without verification to read the header and extract the Key ID (`kid`).
  2. **JWKS Key Resolution:** Resolve the corresponding public key from the cached JWKS (JSON Web Key Set) endpoint. If missing, fetch from `AZURE_AD_JWKS_URL` and update the cache.
  3. **Signature Verification:** Verify the token's cryptographic signature against the public key.
  4. **Claims Validation:** Validate:
     * Audience (`aud`): Must match `AZURE_AD_CLIENT_ID`.
     * Issuer (`iss`): Must match `https://login.microsoftonline.com/{tenant_id}/v2.0` (or `https://sts.windows.net/{tenant_id}/`).
     * Expiration (`exp`): Current time must be before expiration time.
     * Not Before (`nbf`): Current time must be after active time.
  5. **User Registration / Mapping:** Extract claims: `email`, `name` (or `preferred_username`), and `oid` (unique MS object ID). Query `Users` database table:
     * **If User exists:** Fetch user record.
     * **If User does NOT exist:** Auto-register user. If the email contains the domain `@student.uow.edu.my`, set `user_type` = `ActiveStudent` and assign `Active Student` role.
     * **Bootstrap Admin Check:** If the email matches the environment variable `BOOTSTRAP_ADMIN_EMAIL`, automatically assign the `System Admin` role.
  6. **Token Generation:** Generate and return CareerHub custom JWT tokens (Access token and Refresh token).
* **Success Response (200 OK):**
  ```json
  {
    "access_token": "careerhub_jwt_access_token...",
    "refresh_token": "careerhub_jwt_refresh_token...",
    "expires_in": 3600
  }
  ```

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
* **Success Response (200 OK):** Returns new `access_token` and `refresh_token`.

---

### Group B: Student Dashboard APIs

#### 4. Get Current User Profile
* **Endpoint:** `GET /api/v1/users/me`
* **Authorization:** `Bearer <JWT_ACCESS_TOKEN>`
* **Usage:** Returns details to show student name in the welcome banner (*Welcome, [Student Name]*).
* **Success Response (200 OK):**
  ```json
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "student@uow.edu.my",
    "display_name": "Sarah Connor",
    "phone": "+60129998877",
    "user_type": "ActiveStudent",
    "student_id": "99887766",
    "major": "Bachelor of Software Engineering (Hons)"
  }
  ```

#### 5. Fetch Dashboard Categories
* **Endpoint:** `GET /api/v1/categories`
* **Authorization:** `Bearer <JWT_ACCESS_TOKEN>`
* **Usage:** Feeds the *"Top 3 Job Categories"* navigation block.
* **Success Response (200 OK):**
  ```json
  [
    { "id": "1", "name": "Computing & IT", "icon_name": "CommandLineIcon" },
    { "id": "2", "name": "Engineering", "icon_name": "WrenchIcon" },
    { "id": "3", "name": "Hospitality", "icon_name": "BriefcaseIcon" }
  ]
  ```

#### 6. Fetch Jobs (Filtered & Paginated)
* **Endpoint:** `GET /api/v1/jobs`
* **Authorization:** `Bearer <JWT_ACCESS_TOKEN>`
* **Usage:** 
  * Feeds the *"Latest Jobs"* carousel (e.g. via `GET /api/v1/jobs?limit=5`).
  * Powers the search bar input and advanced filters drawer.
* **Query Parameters:** `page`, `limit`, `search`, `category_id`, `job_type`, `location_type`, `experience_level`.
* **Success Response (200 OK):**
  ```json
  [
    {
      "id": "1a2b3c4d-...",
      "job_title": "Event Executive (Skincare Industry)",
      "company_name": "LAVIN PHARMA (M) SDN BHD",
      "category_id": "3",
      "category_name": "Hospitality",
      "job_type": "Internship",
      "location": "Selangor",
      "location_type": "Hybrid",
      "experience_level": "Junior",
      "salary_min": 700.00,
      "salary_max": 1200.00,
      "deadline": "2026-07-31",
      "is_active": true,
      "created_at": "2026-06-15T15:30:00Z"
    }
  ]
  ```

#### 7. Fetch Single Job Details
* **Endpoint:** `GET /api/v1/jobs/{id}`
* **Authorization:** `Bearer <JWT_ACCESS_TOKEN>`
* **Usage:** Triggers when clicking **"Find out more"** to fetch segmented tabs details.
* **Success Response (200 OK):** Returns detailed fields like `job_description`, `responsibilities`, `requirements`, and `how_to_apply`.

---

## 4. Implementation Steps (Timeline)

1. **Step 1: DB Migration Script**  
   Draft and execute the SQL Server DDL migrations to set up the 5 core tables. Set up initial mock data for Categories and Jobs to support dashboard previews.
2. **Step 2: Security & JWT Verification Middleware**  
   Implement JWT validation and parsing using `golang-jwt`. Write standard auth validation checks (HTTP `401` handling) and role check checks (HTTP `403` handling).
3. **Step 3: Microsoft Entra ID JWKS Validator**  
   Implement a helper function to fetch and cache Microsoft's public keys from `https://login.microsoftonline.com/common/discovery/v2.0/keys` to securely verify Microsoft Entra ID tokens.
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
     * **Authorize URI:** `https://login.microsoftonline.com/{AZURE_AD_TENANT_ID}/oauth2/v2.0/authorize` *(Replace `{AZURE_AD_TENANT_ID}` with your active tenant ID)*
     * **Redirect URI:** `https://oidcdebugger.com/redirect`
     * **Client ID:** Your `AZURE_AD_CLIENT_ID`
     * **Scope:** `openid profile email`
     * **Response type:** `id_token`
     * **Response mode:** `form_post`
   * Click **Send Request** and authenticate using your Microsoft/University account.
   * Copy the generated token from the results screen.

2. **Send Exchange Request:**
   * Using an API client (like Postman or cURL), send the token to the local server:
     ```bash
     POST http://localhost:8080/api/v1/auth/login/microsoft
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
   * Query protected endpoints (e.g. `GET http://localhost:8080/api/v1/jobs?limit=5`) to verify authorization checks, data payload formats, and correct HTTP response codes.
