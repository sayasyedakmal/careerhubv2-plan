# CareerHubV2 Sprint 1: API Endpoint Specification

This document details the backend REST API specifications for Sprint 1. It is designed to help frontend and backend developers align on endpoint behavior, MSAL authentication validation, payload structures, query parameters, headers, and error handling.

---

## 🌐 Global Design Decisions & Standards

* **Base URL (Local Development):** `http://localhost:8085/api/v1`
* **Content-Type:** `application/json` (for all request and response bodies)
* **Date Format:** ISO 8601 extended format (`YYYY-MM-DDTHH:mm:ssZ` or `YYYY-MM-DD`).
* **Numeric Formats:** Financial metrics (e.g. salary range limits) are represented as decimals/floats in JSON.
* **CORS Settings:** 
  * Allowed Origins: `http://localhost:3000` (Local Next.js frontend)
  * Allowed Methods: `GET, POST, OPTIONS`
  * Allowed Headers: `Content-Type, Authorization`

---

## 🔑 Authentication & Tokens Lifecycle

All endpoints except **Microsoft Token Exchange** and **Token Refresh** require a custom JSON Web Token (JWT) issued by the CareerHub backend.

### 1. Token Types and Claims

* **Access Token (`access_token`):**
  * Type: JWT (Signed using HS256 with `JWT_SECRET`)
  * Validity: 1 hour (`JWT_ACCESS_EXPIRATION_SEC` = 3600)
  * Header: `Authorization: Bearer <access_token>`
  * Custom Claims:
    ```json
    {
      "sub": "1",                                    // User ID in CareerHub DB (stringified integer)
      "email": "student@uow.edu.my",                 // User email
      "display_name": "Sarah Connor",                // Display Name
      "role": "ActiveStudent",                       // User role (System Admin, Active Student, Alumni)
      "exp": 1782048000,                             // Expiration timestamp
      "iat": 1782044400                              // Issued-at timestamp
    }
    ```

* **Refresh Token (`refresh_token`):**
  * Type: Secure cryptographically signed token string (or high-entropy random token)
  * Validity: 7 days (`JWT_REFRESH_EXPIRATION_SEC` = 604800)
  * Scope: Restricted purely to requesting a new `access_token` via the `/auth/refresh` endpoint.

---

## 🚫 Standardized Error Response Format

All error responses from the backend use the same JSON schema to make error interception and user-friendly messaging uniform:

```json
{
  "error": "A clear, developer-friendly explanation of why the request failed.",
  "code": "UPPERCASE_ERROR_CODE",
  "details": {}
}
```

* `error` (string): Text detail for debugging.
* `code` (string): Unique system-wide error code used by the frontend to trigger specific UI logic.
* `details` (object, optional): Field-level errors (e.g., validation constraint failures).

### Global Error Codes Catalog

| HTTP Status | Error Code (`code`) | Description |
| :--- | :--- | :--- |
| `400` | `INVALID_PAYLOAD` | Request body is empty, malformed, or has invalid JSON syntax. |
| `400` | `INVALID_QUERY_PARAMETER` | One or more URL query string arguments failed format or boundary checks. |
| `400` | `INVALID_PATH_PARAMETER` | URL path variables (e.g. `:id`) fail pattern matching (e.g., not a valid integer). |
| `401` | `MISSING_TOKEN` | No `Authorization` header provided, or header is malformed. |
| `401` | `TOKEN_EXPIRED` | The custom CareerHub access token has expired. Triggers silent refresh flow. |
| `401` | `INVALID_TOKEN` | Custom access token signature is invalid, tampered with, or claims are malformed. |
| `401` | `MICROSOFT_TOKEN_EXPIRED` | The Entra ID `id_token` sent during exchange has expired. |
| `401` | `INVALID_MICROSOFT_TOKEN` | Microsoft signature, audience, or issuer validation checks failed. |
| `401` | `REFRESH_TOKEN_EXPIRED` | The refresh token has expired. Triggers redirect to sign-in page. |
| `401` | `INVALID_REFRESH_TOKEN` | The refresh token does not match records or signature is invalid. |
| `403` | `INSUFFICIENT_PERMISSIONS` | Authenticated user lacks the necessary roles (e.g., Student trying to access admin functions). |
| `404` | `RESOURCE_NOT_FOUND` | The requested entity (user profile, job details, category) does not exist in the system. |
| `422` | `VALIDATION_FAILED` | Payload is syntax-valid JSON but violates business constraints (e.g. negative values, missing required fields). |
| `429` | `RATE_LIMIT_EXCEEDED` | Request rate limit exceeded. Client must throttle requests. |
| `500` | `DATABASE_ERROR` | Internal SQL Server query failure or connection pool timeout. |
| `500` | `INTERNAL_SERVER_ERROR` | Unhandled backend exceptions, key resolution issues, or helper failure. |

---

## 🛠️ Endpoints Specification

### 1. Microsoft Token Exchange

Exchange the Microsoft Entra ID `id_token` returned by MSAL on the frontend for CareerHub access and refresh tokens.

* **HTTP Method:** `POST`
* **Path:** `/auth/login/microsoft`
* **Auth Required:** No
* **Rate Limit:** 30 requests per minute per IP

#### Request Headers
| Header Name | Value | Required | Description |
| :--- | :--- | :--- | :--- |
| `Content-Type` | `application/json` | Yes | Specifies request payload type. |

#### Request Body Schema
| Field Name | Type | Required | Description / Constraints |
| :--- | :--- | :--- | :--- |
| `id_token` | String | Yes | Base64 Encoded Microsoft Entra ID JSON Web Token. Must not be empty. |

##### Request Body Example
```json
{
  "id_token": "<MICROSOFT_ENTRA_ID_TOKEN>"
}
```

#### Success Response
* **HTTP Status:** `200 OK`
* **Headers:**
  | Header Name | Type | Description |
  | :--- | :--- | :--- |
  | `Content-Type` | String | `application/json` |
* **Body Schema:**
  | Field Name | Type | Description |
  | :--- | :--- | :--- |
  | `access_token` | String | JWT access token signed by CareerHub backend (valid 1 hour). |
  | `refresh_token` | String | Custom refresh token for silent renewal (valid 7 days). |
  | `expires_in` | Integer | Access token expiration lifetime in seconds (always `3600`). |

##### Success Response Example
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjNlNDU2Ny1lODliLTEyZDMtYTQ1Ni00MjY2MTQxNzQwMDAiLCJlbWFpbCI6ImFrbWFsLm9AdW93LmVkdS5teSIsImRpc3BsYXlfbmFtZSI6IlNhcmFoIENvbm5vciIsInJvbGUiOiJTeXN0ZW1BZG1pbiIsImV4cCI6MTc4MjA0ODAwMCwiaWF0IjoxNzgyMDQ0NDAwfQ.another-signature-bytes-here",
  "refresh_token": "ch_ref_7f16a5d4e3f2a1b0c9d8e7f6a5d4e3f2",
  "expires_in": 3600
}
```

#### Error Responses

##### Malformed Body (HTTP 400)
```json
{
  "error": "Failed to decode request JSON payload",
  "code": "INVALID_PAYLOAD"
}
```

##### Expired Entra ID Token (HTTP 401)
```json
{
  "error": "The Microsoft Entra ID token has expired",
  "code": "MICROSOFT_TOKEN_EXPIRED"
}
```

##### Invalid Signature or Issuer (HTTP 401)
```json
{
  "error": "Microsoft token signature verification failed: invalid issuer or audience claim",
  "code": "INVALID_MICROSOFT_TOKEN"
}
```

##### SQL Query Failure (HTTP 500)
```json
{
  "error": "Failed to save auto-provisioned student details to database",
  "code": "DATABASE_ERROR"
}
```

---

### 2. Token Refresh

Exchange the long-lived custom refresh token silently for a brand new access token.

* **HTTP Method:** `POST`
* **Path:** `/auth/refresh`
* **Auth Required:** No
* **Rate Limit:** 60 requests per minute per IP

#### Request Headers
| Header Name | Value | Required | Description |
| :--- | :--- | :--- | :--- |
| `Content-Type` | `application/json` | Yes | Specifies request payload type. |

#### Request Body Schema
| Field Name | Type | Required | Description / Constraints |
| :--- | :--- | :--- | :--- |
| `refresh_token` | String | Yes | Secure refresh token value previously issued. Must not be empty. |

##### Request Body Example
```json
{
  "refresh_token": "ch_ref_7f16a5d4e3f2a1b0c9d8e7f6a5d4e3f2"
}
```

#### Success Response
* **HTTP Status:** `200 OK`
* **Headers:**
  | Header Name | Type | Description |
  | :--- | :--- | :--- |
  | `Content-Type` | String | `application/json` |
* **Body Schema:** Same as Microsoft Token Exchange success response body.

##### Success Response Example
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjNlNDU2Ny1lODliLTEyZDMtYTQ1Ni00MjY2MTQxNzQwMDAiLCJlbWFpbCI6ImFrbWFsLm9AdW93LmVkdS5teSIsImRpc3BsYXlfbmFtZSI6IlNhcmFoIENvbm5vciIsInJvbGUiOiJTeXN0ZW1BZG1pbiIsImV4cCI6MTc4MjA1MTYwMCwiaWF0IjoxNzgyMDQ4MDAwfQ.refreshed-signature-bytes-here",
  "refresh_token": "ch_ref_new_8a9b0c1d2e3f4a5b6c7d8e9f0a1b2c3d",
  "expires_in": 3600
}
```

#### Error Responses

##### Malformed Body (HTTP 400)
```json
{
  "error": "Failed to decode refresh token JSON payload",
  "code": "INVALID_PAYLOAD"
}
```

##### Expired Refresh Token (HTTP 401)
```json
{
  "error": "Session has expired. Please sign in again.",
  "code": "REFRESH_TOKEN_EXPIRED"
}
```

##### Revoked or Altered Refresh Token (HTTP 401)
```json
{
  "error": "The provided refresh token is invalid or has been revoked",
  "code": "INVALID_REFRESH_TOKEN"
}
```

---

### 3. Get Current User Profile

Retrieves the profile metadata of the authenticated student to populate the welcome banner and verify session status.

* **HTTP Method:** `GET`
* **Path:** `/users/me`
* **Auth Required:** Yes (Bearer Access Token)

#### Request Headers
| Header Name | Value | Required | Description |
| :--- | :--- | :--- | :--- |
| `Authorization` | `Bearer <access_token>` | Yes | Custom CareerHub Access Token. |

#### Request Query Parameters
None.

#### Success Response
* **HTTP Status:** `200 OK`
* **Body Schema:**
  | Field Name | Type | Description |
  | :--- | :--- | :--- |
  | `id` | Integer | Unique user identifier in the database. |
  | `email` | String | University email address. |
  | `display_name` | String | Full name retrieved from Entra ID / student profile. |
  | `phone` | String / Null | Contact number. May be null. |
  | `user_type` | String | User type classification (`ActiveStudent`, `Alumni`, `Staff`). |
  | `student_id` | String / Null | University Student Card number. Null for non-student profiles. |

##### Success Response Example
```json
{
  "id": 1,
  "email": "akmal.o@uow.edu.my",
  "display_name": "Sarah Connor",
  "phone": "+60129998877",
  "user_type": "ActiveStudent",
  "student_id": "99887766"
}
```

#### Error Responses

##### Missing Authorization Header (HTTP 401)
```json
{
  "error": "Required Authorization Header is missing",
  "code": "MISSING_TOKEN"
}
```

##### Expired Access Token (HTTP 401)
```json
{
  "error": "The access token signature expired at 2026-06-22T17:30:00Z",
  "code": "TOKEN_EXPIRED"
}
```

##### Profile Deleted mid-session (HTTP 404)
```json
{
  "error": "Active user record not found in the database",
  "code": "RESOURCE_NOT_FOUND"
}
```

---

### 4. Fetch Job Categories

Fetches the job categories list. Supports limiting results and sorting by the count of active jobs in each category to populate the Top Categories section of the dashboard.

* **HTTP Method:** `GET`
* **Path:** `/categories`
* **Auth Required:** Yes (Bearer Access Token)

#### Request Headers
| Header Name | Value | Required | Description |
| :--- | :--- | :--- | :--- |
| `Authorization` | `Bearer <access_token>` | Yes | Custom CareerHub Access Token. |

#### Request Query Parameters
| Parameter Name | Type | Required | Default | Constraints / Description |
| :--- | :--- | :--- | :--- | :--- |
| `limit` | Integer | No | `None` | Max categories to return (e.g., `3`). Must be positive. |
| `sort` | String | No | `None` | Sort key. Supports only `job_count` to return categories with the most active jobs first. |

##### Request URL Example
`GET http://localhost:8085/api/v1/categories?limit=3&sort=job_count`

#### Success Response
* **HTTP Status:** `200 OK`
* **Body Schema:** Array of category objects.
  | Field Name | Type | Description |
  | :--- | :--- | :--- |
  | `id` | Integer | Unique category identifier in DB. |
  | `name` | String | User-facing Category Name. |
  | `icon_name` | String | Heroicon name to map on the frontend (e.g. `CommandLineIcon`). |
  | `job_count` | Integer | Number of active jobs currently associated with this category. |

##### Success Response Example
```json
[
  {
    "id": 1,
    "name": "Computing & IT",
    "icon_name": "CommandLineIcon",
    "job_count": 15
  },
  {
    "id": 2,
    "name": "Engineering",
    "icon_name": "WrenchIcon",
    "job_count": 9
  },
  {
    "id": 3,
    "name": "Hospitality",
    "icon_name": "BriefcaseIcon",
    "job_count": 4
  }
]
```

#### Error Responses

##### Invalid limit parameter value (HTTP 400)
```json
{
  "error": "Query parameter 'limit' must be a positive integer",
  "code": "INVALID_QUERY_PARAMETER"
}
```

##### Invalid sort criteria (HTTP 400)
```json
{
  "error": "Query parameter 'sort' supports only 'job_count'",
  "code": "INVALID_QUERY_PARAMETER"
}
```

---

### 5. Fetch Jobs (Filtered & Paginated)

Query, filter, and page through the active job catalog. Used both for the homepage "Latest Jobs" carousel and the full jobs listing/search page.

* **HTTP Method:** `GET`
* **Path:** `/jobs`
* **Auth Required:** Yes (Bearer Access Token)

#### Request Headers
| Header Name | Value | Required | Description |
| :--- | :--- | :--- | :--- |
| `Authorization` | `Bearer <access_token>` | Yes | Custom CareerHub Access Token. |

#### Request Query Parameters
| Parameter Name | Type | Required | Default | Constraints / Description |
| :--- | :--- | :--- | :--- | :--- |
| `page` | Integer | No | `1` | Pagination offset index. Must be positive. |
| `limit` | Integer | No | `10` | Size of pages returned. Must be between `1` and `100`. |
| `search` | String | No | `None` | Case-insensitive keyword matching `job_title`, `company_name`, or category name. |
| `category_id` | Integer | No | `None` | Filter by specific category ID. |
| `job_type_id` | Integer | No | `None` | Filter by job type ID (e.g. `1` for Full-Time, `3` for Internship). |

##### Request URL Example (Latest 5 Jobs Carousel)
`GET http://localhost:8085/api/v1/jobs?limit=5`

##### Request URL Example (Advanced Filter Search)
`GET http://localhost:8085/api/v1/jobs?search=Developer&job_type_id=3&page=2&limit=10`

#### Success Response
* **HTTP Status:** `200 OK`
* **Body Schema:** Array of job summary objects.
  | Field Name | Type | Description |
  | :--- | :--- | :--- |
  | `id` | Integer | Unique Job ID. |
  | `job_title` | String | Job Designation Title. |
  | `company_name` | String | Name of hiring company. |
  | `category_id` | Integer | Category identifier. |
  | `category_name` | String | Category classification name (joined). |
  | `job_type` | String | Employment type name (joined, e.g. "Internship"). |
  | `state` | String | State (e.g. "Selangor"). |
  | `city` | String | City (e.g. "Shah Alam"). |
  | `salary_min` | Float / Null | Minimum salary range. May be null if confidential. |
  | `salary_max` | Float / Null | Maximum salary range. May be null if confidential. |
  | `deadline` | String (DateTime) | Date format `YYYY-MM-DDTHH:mm:ssZ` until which applications are open. |
  | `is_active` | Boolean | True if the job listing is open. Defaults to true. |
  | `created_at` | String (DateTime) | Listing creation timestamp. |

##### Success Response Example
```json
[
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
    "created_at": "2026-06-22T14:38:00Z"
  }
]
```

#### Error Responses

##### Negative page number or out-of-bounds limit (HTTP 400)
```json
{
  "error": "Query parameter 'page' must be greater than zero",
  "code": "INVALID_QUERY_PARAMETER"
}
```

##### Invalid Category Filter (HTTP 400)
```json
{
  "error": "Query parameter 'category_id' must be a valid integer",
  "code": "INVALID_QUERY_PARAMETER"
}
```

---

### 6. Fetch Single Job Details

Fetches full details for a single job to populate the segments tabs (Job Details, Responsibilities, Requirements, Info, and How to Apply) on the job details view page.

* **HTTP Method:** `GET`
* **Path:** `/jobs/:id`
* **Auth Required:** Yes (Bearer Access Token)

#### Request Headers
| Header Name | Value | Required | Description |
| :--- | :--- | :--- | :--- |
| `Authorization` | `Bearer <access_token>` | Yes | Custom CareerHub Access Token. |

#### Request Path Parameters
| Parameter Name | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `id` | Integer | Yes | Target job listing unique identifier. |

##### Request URL Example
`GET http://localhost:8085/api/v1/jobs/101`

#### Success Response
* **HTTP Status:** `200 OK`
* **Body Schema:** Complete job entity details.
  | Field Name | Type | Description |
  | :--- | :--- | :--- |
  | `id` | Integer | Unique identifier. |
  | `job_title` | String | Job Designation Title. |
  | `company_name` | String | Hiring company. |
  | `category_id` | Integer | Category identifier. |
  | `category_name` | String | Category classification name (joined). |
  | `job_type` | String | Employment type (joined). |
  | `state` | String | State (e.g. "Selangor"). |
  | `city` | String | City (e.g. "Shah Alam"). |
  | `salary_min` | Float / Null | Minimum salary range. |
  | `salary_max` | Float / Null | Maximum salary range. |
  | `deadline` | String (DateTime) | Date format `YYYY-MM-DDTHH:mm:ssZ`. |
  | `job_description` | String / Null | Detailed summary block of job description. |
  | `responsibilities`| String / Null | Newline-separated lists of tasks and responsibilities. |
  | `requirements` | String / Null | Newline-separated lists of qualifications/skills required. |
  | `additional_information` | String / Null | Working hours, benefits, etc. |
  | `how_to_apply` | String / Null | Instructions on how to send resumes. |
  | `is_active` | Boolean | True if the job listing is open. |
  | `created_at` | String (DateTime) | Creation timestamp. |
  | `updated_at` | String (DateTime) | Last update timestamp. |

##### Success Response Example
```json
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
  "job_description": "We are seeking a energetic Event Executive to support beauty brand launch activations.",
  "responsibilities": "Plan event booth layouts;\nManage social media posting campaigns;\nCoordinate with caterers and vendors.",
  "requirements": "Pursuing or completed a diploma/degree in PR or Marketing;\nExcellent communication and organizing skills.",
  "additional_information": "Flexible working hours;\nProduct allowance benefits.",
  "how_to_apply": "Send your resume and portfolio to careers@lavinpharma.com.my",
  "is_active": true,
  "created_at": "2026-06-22T14:38:00Z",
  "updated_at": "2026-06-22T14:38:00Z"
}
```

#### Error Responses

##### Invalid integer path parameter format (HTTP 400)
```json
{
  "error": "Path parameter 'id' must be a valid integer",
  "code": "INVALID_PATH_PARAMETER"
}
```

##### Job listing not found (HTTP 404)
```json
{
  "error": "Job listing with ID '999' does not exist",
  "code": "RESOURCE_NOT_FOUND"
}
```
