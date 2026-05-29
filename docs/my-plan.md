# CareerHubV2 Plan

The app name is **CareerHubV2**. This repository is a monorepo containing both the backend and frontend.

## MVP Scope (Aligned with Frontend Designs)
The MVP covers the end-to-end flow for job searching, management, and user authentication with custom mobile/web UI layouts:
*   **User Authentication**: Web and Mobile homepage layouts with dual-authentication (Microsoft Entra ID login + Local password credentials for Alumni). Includes explicit 401 Unauthorized home states.
*   **Job Discovery**:
    *   **All Categories Directory**: Browse jobs by organized classifications.
    *   **Search & Multi-Step Filters**: Advanced job search with specific criteria filters (multi-step drawer panels) and custom "No results" screens.
    *   **Interactive Job Details**: A tabbed detail interface split into **Overview**, **Requirements**, **Details**, and **Apply** tabs.
*   **User Identity**: Dynamic profiles containing customized student **Namecards** (Design B/B1 format).
*   **Core Error States**: Standardized custom HTTP error screens for **403 Forbidden**, **404 Not Found**, **429 Too Many Requests**, and **500 Internal Server Error**.

---

## Architecture & Technology Stack
*   **Backend**: REST API built with **Go** (Gin, SQL driver, JWT, bcrypt) handling business logic, authentication, and logging.
*   **Frontend**: Built with **Next.js**, **Tailwind CSS** (mobile-first approach matching the Figma layouts), and **Heroicons**.
*   **Database**: SQL Server. Migrations managed via `golang-migrate`.

---

## User Groups & Access Control

1.  **Students**
    *   **Current Active Students**: Authenticate via Microsoft Entra ID connected to the `uow.edu.my` domain (`uowstudent` AD tenant).
    *   **Alumni (Old Students)**: Authenticate using locally registered credentials, as their Active Directory accounts have been deactivated.
2.  **SAC Department (Student Affairs Centre)**
    *   Manages job listings and performs manual verification/approval of Alumni registrations.
3.  **System Admin**
    *   Manages global configurations, assigns user groups/roles, and configures role-based permissions.

---

## Detailed Flows

### 1. Alumni (Old Student) Registration & Onboarding
1.  **Submission**: Alumni fills out a registration form containing:
    *   Full Name
    *   Student ID
    *   Email & Contact Number
2.  **Pending State**: Upon submission, the account is set to `Pending`. The alumnus receives an email confirming registration receipt.
3.  **Verification**: The request appears in the SAC Department Staff Dashboard.
4.  **Decision**:
    *   **If Approved**: System sends an email containing a secure registration link to continue setting up their local password. Once set, they can log in.
    *   **If Denied**: System sends a rejection notification email.

### 2. User Profiles & Namecards
All users have a Profile section to update their contact information (Name, Email, Phone Number, Password).
*   **Student Profile Layout**: Features a customized student **Namecard** card design (Design B/B1 layout) indicating registration status and summary details.
*   **Secure Password Reset**: A standardized, secure password reset flow using short-lived tokens generated on the backend and verified via email links.

### 3. Student Portal Features
Students interact with the career board via a highly optimized mobile/web workspace:
*   **Job Category Directory**: Browse all categories page to filter down opportunities.
*   **Multi-Step Filters**: Dynamic modal overlay with tabs and toggles to filter jobs by type, salary, or location.
*   **Tabbed Job Detail View**:
    *   **Overview**: Summary of role, compensation, and timing.
    *   **Requirements**: Bulleted qualifications and education level.
    *   **Details**: Description of tasks and day-to-day work.
    *   **Apply**: Submission panel/external redirection.

### 4. SAC Department Features
*   **Job Posting Management**:
    *   Create job postings with rich field schemas (categorized tags, detailed description sections).
    *   Edit postings and toggle active/inactive states (inactive jobs are hidden from student listings but saved in the DB).
    *   Delete job postings permanently (physically removed from DB, audited).
*   **Registration Management**:
    *   Approve or deny pending Alumni registrations from a specialized dashboard list.

### 5. System Admin Features
Modular Role-Based Access Control (RBAC):
*   Create and manage user roles/groups.
*   Assign/unassign users to roles dynamically.
*   Manage role permissions.

---

## Non-Functional & Security Requirements
*   **Audit Logs**: All critical mutating actions (approvals, role changes, job creation/deletion) must create a record in the `AuditLogs` table.
*   **Standardized Error Handling**: Any failure state or route exception redirects to the corresponding custom design view (401, 403, 404, 429, 500).
*   **Database Migrations**: All tables and seed values must be versioned.
