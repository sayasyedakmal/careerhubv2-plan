# CareerHubV2 API Contracts & Payload Mapping

This document provides a clean, standard API Endpoint Table directly mapped to the full specifications in `shared/openapi.yaml`. 

It serves as the development contract between the Go backend developer and the Next.js frontend developer.

---

| API endpoint | Request payload | Response payload | Error Handling / Status Codes |
| :--- | :--- | :--- | :--- |
| **`POST /api/v1/auth/login/microsoft`**<br><br>*Login using Microsoft Entra ID token.* | <pre>{<br>  "id_token": "eyJhbGciOi..."<br>}</pre> | <pre>{<br>  "access_token": "access_token_jwt...",<br>  "refresh_token": "refresh_token_jwt...",<br>  "expires_in": 3600<br>}</pre> | <ul><li>**`400`**: Malformed payload.</li><li>**`401`**: Invalid or expired Microsoft token, or domain unauthorized.</li><li>**`500`**: Identity Provider connection failed.</li></ul> |
| **`POST /api/v1/auth/login`**<br><br>*Login using Local Credentials (for Alumni).* | <pre>{<br>  "email": "alumni@example.com",<br>  "password": "SecurePassword123"<br>}</pre> | <pre>{<br>  "access_token": "access_token_jwt...",<br>  "refresh_token": "refresh_token_jwt...",<br>  "expires_in": 3600<br>}</pre> | <ul><li>**`400`**: Malformed payload.</li><li>**`401`**: Invalid email/password, or account is pending/denied (returns `ALUMNI_PENDING` or `ALUMNI_DENIED` `error_code`).</li></ul> |
| **`POST /api/v1/auth/refresh`**<br><br>*Refresh user Access Token.* | <pre>{<br>  "refresh_token": "refresh_token_jwt..."<br>}</pre> | <pre>{<br>  "access_token": "new_access_token_jwt...",<br>  "refresh_token": "new_refresh_token_jwt...",<br>  "expires_in": 3600<br>}</pre> | <ul><li>**`400`**: Missing token.</li><li>**`401`**: Invalid or expired refresh token (returns `TOKEN_INVALID`).</li></ul> |
| **`POST /api/v1/auth/register`**<br><br>*Register Alumni account (Triggers pending verification status).* | <pre>{<br>  "email": "alumni@example.com",<br>  "display_name": "John Doe",<br>  "phone": "+60123456789",<br>  "student_id": "0123456"<br>}</pre> | <pre>{<br>  "message": "Registration submitted and pending approval."<br>}</pre> | <ul><li>**`400`**: Invalid formatting.</li><li>**`409`**: Email or StudentID already registered in the system.</li><li>**`422`**: Input validation failed (e.g. empty name, invalid phone format).</li></ul> |
| **`POST /api/v1/auth/password-reset`**<br><br>*Request tokenized password reset link via SMTP email.* | <pre>{<br>  "email": "alumni@example.com"<br>}</pre> | <pre>{<br>  "message": "If the email is registered, a reset link has been sent."<br>}</pre> | <ul><li>**`400`**: Missing or invalid email address.</li><li>**`500`**: Mail server (SMTP) unreachable.</li></ul> |
| **`POST /api/v1/auth/password-change`**<br><br>*Complete password setup/reset using secure token.* | <pre>{<br>  "token": "reset_token_hash_here",<br>  "password": "NewStrongSecurePassword123!"<br>}</pre> | <pre>{<br>  "message": "Password updated successfully."<br>}</pre> | <ul><li>**`400`**: Invalid or expired setup token.</li><li>**`422`**: Password strength requirements validation failed.</li></ul> |
| **`GET /api/v1/users/me`**<br><br>*Get current user profile and enrollment details (Requires JWT).* | *None (Header: Bearer JWT)* | <pre>{<br>  "id": "123",<br>  "email": "student@uow.edu.my",<br>  "display_name": "Sarah Connor",<br>  "phone": "+60129998877",<br>  "user_type": "Student",<br>  "registration_status": "N/A",<br>  "student_id": "99887766",<br>  "major": "Bachelor of Software Engineering (Hons)"<br>}</pre> | <ul><li>**`401`**: Token missing, invalid, or expired (`TOKEN_EXPIRED` / `TOKEN_INVALID`).</li></ul> |
| **`PUT /api/v1/users/me`**<br><br>*Update personal profile attributes (Requires JWT).* | <pre>{<br>  "display_name": "Sarah Connor",<br>  "phone": "+60121112233",<br>  "major": "Bachelor of Software Engineering (Hons)"<br>}</pre> | <pre>{<br>  "id": "123",<br>  "email": "student@uow.edu.my",<br>  "display_name": "Sarah Connor",<br>  "phone": "+60121112233",<br>  "user_type": "Student",<br>  "registration_status": "N/A",<br>  "student_id": "99887766",<br>  "major": "Bachelor of Software Engineering (Hons)"<br>}</pre> | <ul><li>**`401`**: Token missing or expired.</li><li>**`422`**: Validation failed.</li></ul> |
| **`GET /api/v1/categories`**<br><br>*Fetch All Categories for Dashboard Directories.* | *None (Header: Bearer JWT)* | <pre>[<br>  {<br>    "id": "1",<br>    "name": "Computing & IT",<br>    "icon_name": "CommandLineIcon",<br>    "active_job_count": 15<br>  },<br>  {<br>    "id": "2",<br>    "name": "Engineering",<br>    "icon_name": "WrenchIcon",<br>    "active_job_count": 9<br>  }<br>]</pre> | <ul><li>**`401`**: Token missing or expired.</li></ul> |
| **`POST /api/v1/admin/registrations/{id}/approve`**<br><br>*Approve a pending Alumni onboarding request (SAC Only).* | *None (Header: Bearer JWT)* | <pre>{<br>  "message": "Alumni registration approved. Welcome email dispatched."<br>}</pre> | <ul><li>**`401`**: Token missing/invalid.</li><li>**`403`**: User lacks administrative permissions.</li><li>**`404`**: Target Registration ID not found.</li><li>**`500`**: SMTP delivery failed.</li></ul> |
| **`POST /api/v1/admin/registrations/{id}/deny`**<br><br>*Reject a pending Alumni onboarding request (SAC Only).* | *None (Header: Bearer JWT)* | <pre>{<br>  "message": "Alumni registration denied."<br>}</pre> | <ul><li>**`401`**: Token missing/invalid.</li><li>**`403`**: User lacks administrative permissions.</li><li>**`404`**: Target Registration ID not found.</li></ul> |
| **`GET /api/v1/jobs`**<br><br>*Fetch jobs with filters (Query params: `page`, `limit`, `search`, `category_id`, `job_type`, `location_type`, `experience_level`, `salary_min`, `salary_max`).* | *None (Header: Bearer JWT)* | <pre>{<br>  "data": [<br>    {<br>      "id": "45",<br>      "job_title": "Junior Full Stack Engineer",<br>      "company_name": "Tech Corp",<br>      "category_id": "1",<br>      "category_name": "Computing & IT",<br>      "job_type": "Full-Time",<br>      "location": "Kuala Lumpur",<br>      "location_type": "Hybrid",<br>      "experience_level": "Junior",<br>      "salary_min": 3500.00,<br>      "salary_max": 4500.00,<br>      "deadline": "2026-10-31",<br>      "is_active": true,<br>      "created_at": "2026-05-29T08:00:00Z"<br>    }<br>  ],<br>  "pagination": {<br>    "page": 1,<br>    "limit": 10,<br>    "total_records": 45,<br>    "total_pages": 5<br>  }<br>}</pre> | <ul><li>**`401`**: Token missing/invalid.</li><li>**`429`**: Rate limit exceeded (triggers 429 UI screen, returns `Retry-After` header).</li><li>**`500`**: Database query execution failure.</li></ul> |
| **`POST /api/v1/jobs`**<br><br>*Create a new Job Posting (SAC Only).* | <pre>{<br>  "job_title": "DevOps Engineer",<br>  "company_name": "Cloud Systems",<br>  "category_id": "1",<br>  "job_type": "Full-Time",<br>  "location": "Penang",<br>  "location_type": "Onsite",<br>  "experience_level": "Mid",<br>  "has_salary": true,<br>  "salary_min": 5000.00,<br>  "salary_max": 7000.00,<br>  "deadline": "2026-09-30",<br>  "position_count": 2,<br>  "job_description": "# Role Description...",<br>  "responsibilities": "- Maintain pipelines...",<br>  "requirements": "- 2 years experience...",<br>  "additional_information": "Dress code: Smart Casual",<br>  "how_to_apply": "Apply via Portal...",<br>  "is_active": true<br>}</pre> | <pre>{<br>  "id": "46",<br>  "job_title": "DevOps Engineer",<br>  "company_name": "Cloud Systems",<br>  "category_id": "1",<br>  "job_type": "Full-Time",<br>  "location": "Penang",<br>  "location_type": "Onsite",<br>  "experience_level": "Mid",<br>  "has_salary": true,<br>  "salary_min": 5000.00,<br>  "salary_max": 7000.00,<br>  "deadline": "2026-09-30",<br>  "position_count": 2,<br>  "job_description": "# Role Description...",<br>  "responsibilities": "- Maintain pipelines...",<br>  "requirements": "- 2 years experience...",<br>  "additional_information": "Dress code: Smart Casual",<br>  "how_to_apply": "Apply via Portal...",<br>  "is_active": true,<br>  "created_at": "2026-05-29T10:45:00Z"<br>}</pre> | <ul><li>**`401`**: Token missing/invalid.</li><li>**`403`**: User lacks SAC permissions.</li><li>**`422`**: Validation failed (e.g., negative salary values, missing required fields).</li></ul> |
| **`GET /api/v1/jobs/{id}`**<br><br>*Fetch full detail specs for single job tabs (Overview, Requirements, Details, Apply).* | *None (Header: Bearer JWT)* | *(Returns full Job JSON structure as shown in POST /jobs response payload)* | <ul><li>**`401`**: Token missing/invalid.</li><li>**`404`**: Job ID not found (triggers 404 UI error screen).</li></ul> |
| **`PUT /api/v1/jobs/{id}`**<br><br>*Replace a Job details structure (SAC Only).* | *(Same payload structure as POST /jobs)* | *(Returns updated full Job JSON structure)* | <ul><li>**`401`**: Token missing/invalid.</li><li>**`403`**: User lacks SAC permissions.</li><li>**`404`**: Job ID not found.</li><li>**`422`**: Validation failed on job properties.</li></ul> |
| **`DELETE /api/v1/jobs/{id}`**<br><br>*Permanently delete a Job Posting (SAC Only).* | *None (Header: Bearer JWT)* | <pre>{<br>  "message": "Job posting permanently deleted."<br>}</pre> | <ul><li>**`401`**: Token missing/invalid.</li><li>**`403`**: User lacks SAC permissions.</li><li>**`404`**: Job ID not found.</li></ul> |
| **`PATCH /api/v1/jobs/{id}/status`**<br><br>*Toggle active/inactive job state (SAC Only).* | *None (Header: Bearer JWT)* | <pre>{<br>  "is_active": false,<br>  "message": "Job status updated successfully."<br>}</pre> | <ul><li>**`401`**: Token missing/invalid.</li><li>**`403`**: User lacks SAC permissions.</li><li>**`404`**: Job ID not found.</li></ul> |
| **`GET /api/v1/admin/roles`**<br><br>*List all user roles (System Admin Only).* | *None (Header: Bearer JWT)* | <pre>[<br>  { "id": "1", "role_name": "System Admin" },<br>  { "id": "2", "role_name": "SAC Department" }<br>]</pre> | <ul><li>**`401`**: Token missing/invalid.</li><li>**`403`**: User lacks System Admin permissions.</li></ul> |
| **`POST /api/v1/admin/roles`**<br><br>*Create a new user role (System Admin Only).* | <pre>{<br>  "role_name": "Company Partner"<br>}</pre> | <pre>{<br>  "id": "5",<br>  "role_name": "Company Partner"<br>}</pre> | <ul><li>**`401`**: Token missing/invalid.</li><li>**`403`**: User lacks System Admin permissions.</li><li>**`422`**: Validation failed (e.g. empty or duplicate role name).</li></ul> |
| **`PUT /api/v1/admin/users/{id}/roles`**<br><br>*Assign target user to roles list (System Admin Only).* | <pre>{<br>  "role_ids": ["1", "3"]<br>}</pre> | <pre>{<br>  "message": "User roles updated successfully."<br>}</pre> | <ul><li>**`401`**: Token missing/invalid.</li><li>**`403`**: User lacks System Admin permissions.</li><li>**`404`**: Target User ID not found.</li></ul> |
| **`GET /api/v1/health`**<br><br>*Check API system status.* | *None* | <pre>{<br>  "status": "up"<br>}</pre> | <ul><li>**`500`**: Database connection offline or system experiencing heavy disruption.</li></ul> |

---

## Standard Error Response Structures

### 1. Standard Error (400, 401, 403, 404, 409, 500)
```json
{
  "error": "Short-lived token expired.",
  "code": 401,
  "error_code": "TOKEN_EXPIRED"
}
```

### 2. Validation Failures (422 Unprocessable Entity)
```json
{
  "error": "Validation failed",
  "code": 422,
  "details": [
    {
      "field": "email",
      "message": "Invalid email format. Must contain a valid domain."
    },
    {
      "field": "password",
      "message": "Password is too weak. Must contain at least one uppercase letter and number."
    }
  ]
}
```
