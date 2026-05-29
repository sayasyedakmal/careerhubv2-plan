# Developer Onboarding & Guidance Notes

This document acts as a companion guide to `docs/weekly_milestone.md` for project managers and tech leads. It evaluates the learning curve for junior/beginner developers, identifies potential technical hurdles, and provides direct mitigation strategies to help them stay on track.

---

## 1. Why the Current Flow is Beginner-Friendly

The 5-week milestone path is structured to minimize cognitive overload and prevent developers from getting blocked:

* **Decoupled Frontend via MSW (Mock Service Worker):** 
  By mock-routing API endpoints in Weeks 1 and 2, frontend developers can build out the user interface, design systems, and responsive layouts without waiting for the backend APIs to be fully written.
* **Immediate Visual Feedback (Week 1):** 
  Starting with basic Tailwind shells, Docker initialization, and a simple backend `/api/v1/health` check gives beginners immediate "quick wins" to build momentum.
* **Logical Complexity Scaling:** 
  The project scales naturally from static layouts (Week 1), to user authentication (Week 2), to relational data fetching (Week 3), to complex business workflow state machines (Week 4), and finally to deployment (Week 5).

---

## 2. Anticipated Developer Hurdles & Mitigations

While the milestone schedule is sequential, beginners may face challenges in certain areas. Use the table below to anticipate these blocks and guide them:

| Week | Potential Hurdles for Beginners | Why It Is Difficult | How to Mitigate (Where to Direct Them) |
| :--- | :--- | :--- | :--- |
| **Week 1** | **Docker & SQL Server Setup** | Configuring local databases, managing volumes, and running migration tools can be confusing if they haven't worked with SQL Server. | Direct them to run `docker-compose up -d` using the root `docker-compose.yml` and execute the updated schema script in `docs/init.sql` directly. |
| **Week 2** | **Dual Authentication & Custom RBAC Middleware** | Writing an active token-validation middleware in Go (Gin) that checks claims, queries RBAC databases, and returns correct HTTP error states. | Refer them to `docs/auth-process.md` and `docs/rbac-implementation-guide.md` which contain step-by-step logic and string mappings. |
| **Week 2** | **MSAL React Context Integration** | Managing Azure AD/Entra ID SSO configurations and binding claims to a global React `AuthContext`. | Direct them to use official `@azure/msal-react` packages and review the frontend setup steps in `frontend/PLAN.md`. |
| **Week 3** | **Tabbed Payload Structuring** | Mapping single database records to split payload responses for the dynamic **Job Detail Tabs** (Overview, Requirements, Details, Apply). | Show them how the `Jobs` table in `docs/init.sql` maps directly to these tabs (`JobDescription`, `Requirements`, `Responsibilities`, `HowToApply`). |
| **Week 3** | **Dynamic Search Queries & Pagination** | Writing safe dynamic SQL queries in Go with `OFFSET` and `FETCH NEXT` based on optional URL parameters (category, type, salary range). | Emphasize using safe parameterized queries or query builders to prevent SQL Injection, and returning status `200` with an empty collection if no results match. |
| **Week 4** | **SMTP Mailers & State Handlers** | Configuring background email workers and managing state switches (e.g., transition from local registration requests $\rightarrow$ approved/denied). | Have them implement standard SMTP packages in Go utilizing environment values from `backend/.env.example`. |
| **Week 5** | **Deploying to IIS (Windows Server)** | Setting up Go as a Windows Service under `HttpPlatformHandler` and configuring Next.js rewrites via `web.config`. | Have them follow the exact deployment recipes in `docs/iis-deployment.md` and `frontend/iis-deployment.md`. They can copy-paste the `web.config` rewrite blocks directly. |

---

## 3. Recommended Mentoring Milestones

To keep junior developers aligned, set up 15-minute syncs at these key transition points:

1. **End of Week 1 Sync:** Verify that the SQL Server Docker container is running successfully and the Go health endpoint works.
2. **Mid-Week 2 Sync:** Review the custom Authorization middleware structure before they start mapping individual route permissions.
3. **End of Week 3 Sync:** Test the frontend multi-step filter panel using MSW mocked values to ensure smooth UX transitions.
4. **Mid-Week 5 Sync:** Run a dry-run production build (`npm run build` static export) locally to check for any broken dynamic links before transferring files to IIS.
