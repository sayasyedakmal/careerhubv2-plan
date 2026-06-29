# Sprint 1 Spec Review: Auth & Student Dashboard APIs

This document outlines the design gaps, ambiguities, missing edge cases, contradictions, and critical questions identified during the review of the [Sprint 1 Implementation Plan](file:///c:/repo/personal/careerhubv2-plan/backend/sprint1/plan.md).

---

## 1. Ambiguities

### A. User Auto-Registration Rules & Roles
* **The Issue:** The plan states: *“If the email contains the domain `@student.uow.edu.my`, set `UserType = 'ActiveStudent'`... and assign the `Active Student` role”* (Line 163).
* **Ambiguity:** What happens if a user signs in with a Microsoft account that does **not** have the `@student.uow.edu.my` domain? (e.g., a staff member with `@uow.edu.my`, or a personal email/external tenant email). 
  * Are non-student emails blocked from logging in?
  * Or are they registered with a different default `UserType` (e.g., `'Alumni'`, `'Staff'`, `'SystemAdmin'`)? If so, what role are they assigned?
* **Resolution:** 
  * Active Student is renamed to **Student** for clarity (avoiding active/inactive confusion).
  * Emails ending in `@student.uow.edu.my` are registered with `UserType = 'Student'` and auto-assigned the `Student` role.
  * Emails ending in `@uow.edu.my` (staff) are registered with `UserType = 'Staff'` and have **no default roles** (they must be manually assigned roles).
  * Any other emails are registered with `UserType = 'External'` and have **no default roles** (pending registration/approval workflows).

### B. Bootstrap Admin User Type
* **The Issue:** If the email matches the environment variable `BOOTSTRAP_ADMIN_EMAIL`, the backend automatically assigns the `System Admin` role.
* **Ambiguity:** What is their primary `UserType` stored in the `Users` table (`'Staff'`, `'SystemAdmin'`, or `'ActiveStudent'`)?
* **Resolution:** The bootstrap admin is registered with `UserType = 'SystemAdmin'` and auto-assigned the `System Admin` role.

### C. Refresh Token Architecture & Storage
* **The Issue:** The plan specifies a `POST /api/v1/auth/refresh` endpoint but does not define where refresh tokens are stored. The SQL schema in Section 2 does not contain a `RefreshTokens` or `UserSessions` table.
* **Ambiguity:**
  * Are refresh tokens intended to be stateless JWTs? If so, is there a mechanism to revoke or blacklist them (e.g., on logout)?
  * Or should they be stateful and stored in the database? If stateful, we are missing a schema definition for token tracking.
* **Resolution:** Refresh tokens are stateless JWTs validated "on the fly" by the backend. No database tables or tracking schemas are needed. On logout, the frontend simply discards the token.

### D. Student ID Resolution
* **The Issue:** The `/api/v1/users/me` response returns a `"student_id": "99887766"`.
* **Ambiguity:** How is the `StudentID` resolved during Microsoft Entra ID registration? 
  * Is it extracted from a custom claim inside the Microsoft ID token?
  * Or is it parsed from the email prefix (e.g., extracting `99887766` from `99887766@student.uow.edu.my`)?
* **Resolution:** The `StudentID` is extracted from the username prefix of the `@student.uow.edu.my` email address (e.g. `99887766` is extracted from `99887766@student.uow.edu.my`).

---

## 2. Missing Edge Cases

### A. User Profile Synchronization
* **The Issue:** If a user updates their Display Name or Email in Azure AD / Entra ID after their initial registration, does the backend update the matching database record on subsequent logins, or is the record only written once at initial registration?
* **Resolution:** On every successful Entra ID token exchange (`POST /api/v1/auth/login/microsoft`), the backend will check if the claims (Display Name, Email) in the incoming token match the existing record in the database. If any changes are detected, the backend will sync and update the database record automatically before returning the session JWT.

### B. Expired and Inactive Jobs
* **For `/api/v1/categories?sort=job_count`:** Does `job_count` count all jobs in the category, or only **active and unexpired** jobs (`IsActive = 1` and `DeadlineAt > GETDATE()`)?
* **For `/api/v1/jobs`:** Should this endpoint filter out inactive (`IsActive = 0`) and expired jobs by default for students? 
* **Resolution:**
  * For `/api/v1/categories`, only active and unexpired jobs (`IsActive = 1` AND (`DeadlineAt IS NULL` OR `DeadlineAt > GETDATE()`)) are counted. The field name in the JSON response is renamed from `job_count` to `active_job_count` (or `job_count_active`) for explicit clarity.
  * For `/api/v1/jobs`, the endpoint will filter out inactive (`IsActive = 0`) and expired (`DeadlineAt <= GETDATE()`) jobs by default for students.

### C. Pagination Metadata Format
* **The Issue:** The success response for `/api/v1/jobs` shows a plain list: `[ { "id": 101, ... } ]`.
* **Missing Feature:** When fetching paginated jobs (e.g., via `page` and `limit`), does the response format include metadata like `total_records`, `total_pages`, or `current_page`? Returning just a flat array makes it difficult for the frontend to render proper pagination controls.
* **Resolution:** All paginated list endpoints (including `/api/v1/jobs`) will return a standard envelope containing `data` and `pagination` metadata (`page`, `limit`, `total_records`, `total_pages`).

---

## 3. Unstated Assumptions

### A. Tenant Restriction Enforcement
* **The Assumption:** The configuration includes `MICROSOFT_ENTRA_TENANT_ID` to restrict logins.
* **Unstated Rule:** If a user logs in with a valid Microsoft ID token but from an external/personal tenant (different Tenant ID), does the backend reject the login? We need to verify if the backend explicitly validates the `tid` claim (tenant ID) or issuer against our environment variable.
* **Resolution:** The backend will strictly validate the `tid` claim (tenant ID) and the `iss` claim (issuer) in the Microsoft ID token against `MICROSOFT_ENTRA_TENANT_ID`. If they do not match, the request will be rejected with an `HTTP 403 Forbidden` status code.

### B. Token Rotation & Reuse
* **The Assumption:** `/auth/refresh` returns a new access and refresh token.
* **Unstated Rule:** Is refresh token rotation active (i.e., the old refresh token is invalidated immediately upon use), or can the same refresh token be reused until expiration?
* **Resolution:** 
  * **Sprint 1 Scope (Option 1):** Refresh tokens are stateless JWTs that can be reused until their expiration time. Calling `/api/v1/auth/refresh` returns a new access token and returns the valid refresh token.
  * **Future Backlog (Option 2 - Industry Standard Refresh Token Rotation):** For future sprints requiring heightened security, we will transition to full Refresh Token Rotation with Token Family tracking and Reuse Detection (requiring a `UserSessions` DB table / Redis cache and a 10–30 second concurrency grace period).

---

## 4. Contradictions

### A. JWKS Endpoint Mismatch
* **Section 3 (Line 151):** Lists `MICROSOFT_ENTRA_JWKS_URL` as `https://login.microsoftonline.com/{tenant_id}/discovery/v2.0/keys`.
* **Section 4 (Line 267):** Says the validator fetches keys from `https://login.microsoftonline.com/common/discovery/v2.0/keys` (using `common` instead of `{tenant_id}`).
* **Contradiction:** Using `/common/` allows keys from any Microsoft tenant to verify successfully, whereas using the tenant-specific endpoint restricts keys to just those valid for the university.
* **Resolution:** Standardize on `https://login.microsoftonline.com/{tenant_id}/discovery/v2.0/keys` across all backend code and documentation so public keys are scoped strictly to our university tenant.

### B. Out-of-Scope Local Credentials vs. Password Schema
* **Section 1 (Line 10):** *"Local credentials login and registration are out of scope (all authentication is routed through Microsoft)"*.
* **Section 2:** Includes the `PasswordResets` table (Line 56) and a `PasswordHash` field in `Users` (Line 46).
* **Contradiction:** If local auth is entirely out of scope for Sprint 1, these tables and fields are redundant for the initial database setup.
* **Resolution:** Keep `PasswordHash` (nullable) and the `PasswordResets` table in the initial database migration script to avoid schema migration overhead in Sprint 2, but mark them explicitly as dormant/unused in Sprint 1 backend code.

---

## 5. Missing Error States

### A. Token Exchange Failures (`POST /api/v1/auth/login/microsoft`)
* **What happens if the Microsoft token signature verification fails or the token is expired?**
  * **Resolution:** Return `401 Unauthorized` with `{"error_code": "INVALID_MICROSOFT_TOKEN", "message": "Microsoft ID token verification failed or token expired"}`.
* **What happens if the Entra ID token is valid, but the database connection is down?**
  * **Resolution:** Return `500 Internal Server Error` with `{"error_code": "DATABASE_ERROR", "message": "Internal service error while processing user login"}`.
* **What happens if a user logs in with an unauthorized tenant?**
  * **Resolution:** Return `403 Forbidden` with `{"error_code": "INVALID_TENANT", "message": "Authenticated Microsoft tenant is not authorized for this application"}`.

### B. Refresh Token Failures (`POST /api/v1/auth/refresh`)
* **What error code/payload is returned if the refresh token is expired or tampered with?**
  * **Resolution:** Return `401 Unauthorized` with `{"error_code": "INVALID_REFRESH_TOKEN", "message": "Refresh token is invalid or expired. Please sign in again."}`.

---

## 6. Spec Review Status: COMPLETED

All identified ambiguities, missing edge cases, unstated assumptions, contradictions, and missing error states in [plan.md](file:///c:/repo/personal/careerhubv2-plan/backend/sprint1/plan.md) have been thoroughly analyzed and resolved above. 

The main plan document `backend/sprint1/plan.md` can now be updated to incorporate all approved resolutions.
