# Roles and Permissions Matrix

This document defines the base roles and their exact permission strings for the CareerHubV2 RBAC system. This should be used to seed the database upon initial deployment.

## Defined Roles

1. **System Admin**: Has full control over the system, roles, and settings. 
   * *Assignment*: Bootstrapped via the environment variable `BOOTSTRAP_ADMIN_EMAIL` on first login, after which they can assign other Admins.
2. **SAC Department**: Manages job postings and approves alumni registrations.
   * *Assignment*: Assigned manually by a System Admin.
3. **Active Student**: Standard user with a Microsoft account. Can view jobs.
   * *Assignment*: **Auto-assigned** on first login if the Entra ID authenticated user's email ends with `@uow.edu.my`.
4. **Alumni**: Standard user who registered locally and was approved. Can view jobs.
   * *Assignment*: Assigned automatically once their local registration is Approved by SAC or Admin.

## Permissions Matrix

| Permission String | Description | System Admin | SAC Dept | Active Student | Alumni |
| :--- | :--- | :---: | :---: | :---: | :---: |
| `manage_user_group: default=student` | Create group, assign users to group | ✅ | ❌ | ❌ | ❌ |
| `manage_users ` | View all users, edit user profiles | ✅ | ❌ | ❌ | ❌ |
| `approve_alumni` | Approve or deny pending alumni | ✅ | ✅ | ❌ | ❌ |
| `create_job` | Create new job postings | ✅ | ✅ | ❌ | ❌ |
| `edit_job` | Edit existing job postings | ✅ | ✅ | ❌ | ❌ |
| `delete_job` | Delete job postings permanently | ✅ | ✅ | ❌ | ❌ |
| `toggle_job` | Change job status (active/inactive) | ✅ | ✅ | ❌ | ❌ |
| `view_jobs` | View active job listings and details | ✅ | ✅ | ✅ | ✅ |
| `apply_job` | (Future) Apply to a job posting | ✅ | ✅ | ✅ | ✅ |

## Implementation Note
The backend `RequirePermission("...")` middleware should strictly check these string values against the `Permissions` assigned to the `Group/Role` of the currently authenticated `UserID`.
