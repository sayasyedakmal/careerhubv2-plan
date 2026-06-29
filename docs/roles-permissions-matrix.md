# Roles and Permissions Matrix

This document defines the base roles and their exact permission strings for the CareerHubV2 RBAC system. This should be used to seed the database upon initial deployment.

## Defined Roles and Assignments

1. **System Admin**: Has full control over the system, roles, and settings.
   * *Assignment*: Bootstrapped via the environment variable `BOOTSTRAP_ADMIN_EMAIL` on first login, after which they can assign other Admins.
2. **SAC Department**: Manages job postings and approves alumni/external registrations.
   * *Assignment*: Assigned manually by a System Admin to users with a `Staff` UserType.
3. **Student**: Standard student user. Can view jobs and apply.
   * *Assignment*: **Auto-assigned** on first login if the Entra ID authenticated user's email ends with `@student.uow.edu.my`.
4. **Alumni**: Registered alumni user. Can view jobs and apply.
   * *Assignment*: Promoted from `External` user status once their registration details are Approved by the SAC Department or System Admin.
5. **External**: Any user signing in with a non-university email address (e.g. personal email).
   * *Assignment*: Auto-assigned `UserType = 'External'` on first login. They start with **no permissions/roles** and must complete a registration form to apply for approval (e.g., as Alumni).

## Permissions Matrix

| Permission String | Description | System Admin | SAC Dept | Student | Alumni | External |
| :--- | :--- | :---: | :---: | :---: | :---: | :---: |
| `manage_user_groups` | Create group, assign users to group | ✅ | ❌ | ❌ | ❌ | ❌ |
| `manage_users` | View all users, edit user profiles | ✅ | ❌ | ❌ | ❌ | ❌ |
| `approve_alumni` | Approve or deny pending alumni registrations | ✅ | ✅ | ❌ | ❌ | ❌ |
| `create_job` | Create new job postings | ✅ | ✅ | ❌ | ❌ | ❌ |
| `edit_job` | Edit existing job postings | ✅ | ✅ | ❌ | ❌ | ❌ |
| `delete_job` | Delete job postings permanently | ✅ | ✅ | ❌ | ❌ | ❌ |
| `toggle_job` | Change job status (active/inactive) | ✅ | ✅ | ❌ | ❌ | ❌ |
| `view_jobs` | View active job listings and details | ✅ | ✅ | ✅ | ✅ | ❌ |
| `apply_job` | (Future) Apply to a job posting | ✅ | ✅ | ✅ | ✅ | ❌ |

## Implementation Note
The backend `RequirePermission("...")` middleware should strictly check these string values against the `Permissions` assigned to the `Group/Role` of the currently authenticated `UserID`.
