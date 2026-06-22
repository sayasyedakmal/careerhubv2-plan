# CareerHubV2 Weekly Milestones

This document breaks down the development of the CareerHubV2 MVP into a 5-week schedule, updated to reflect the high-fidelity mobile-first design system layout from `@frontend/design/`.

---

## Week 1: Foundation, Styling & Shell Navigation
**Goal:** Establish the development environment, theme setup, and core layout components.

### Backend Developer
- [ ] Initialize Go module and install dependencies (Gin, SQL driver, JWT, bcrypt).
- [ ] Set up `docker-compose.yml` for SQL Server database.
- [ ] Configure `golang-migrate` and create the initial migration files (`db-schema.md`).
- [ ] Implement DB seeder script with pre-defined categories (Matching Category pages).
- [ ] Build `/api/v1/health` endpoint.

### Frontend Developer
- [ ] Initialize Next.js project with Tailwind CSS and TypeScript.
- [ ] Set up global color variables reflecting official branding (Navy `#001C41` and Red `#E21F26`).
- [ ] Configure MSW (Mock Service Worker) for API intercepting.
- [ ] Build the Core Layout containing the responsive **Mobile-first Drawer Navigation Header** (hamburger menu trigger, slide-out panel, and profile dropdown states).

---

## Week 2: Authentication & Homepage Views
**Goal:** Dual login integration, authentication middlewares, and homepage unauthorized state handling.

### Backend Developer
- [ ] Build local user credential registration and password validation (`POST /api/v1/auth/login`).
- [ ] Integrate JWT validation for Microsoft Entra ID tokens (`POST /api/v1/auth/login/microsoft`).
- [ ] Implement authentication and modular RBAC authorization middleware (Refer to `docs/roles-permissions-matrix.md`).
- [ ] Set up Token Bucket **Rate Limiting** middleware returning custom `429 Too Many Requests` states.

### Frontend Developer
- [ ] Build responsive dual-login Homepage (`Homepage - sign in.png`) supporting both Entra ID buttons and local forms.
- [ ] Setup Microsoft MSAL React provider.
- [ ] Create `AuthContext` to store JWT claims and permissions.
- [ ] Design the Homepage unauthorized landing state (`Homepage - 401 Unauthorized.png`) shown when unauthenticated users access protected paths.

---

## Week 3: Job Board, Advanced Filtering & Tabbed Details
**Goal:** Dynamic job discovery, multi-step filters, and category-focused listings.

### Backend Developer
- [ ] Implement `/api/v1/categories` endpoint returning classification options.
- [ ] Code `/api/v1/jobs` supporting dynamic WHERE parameters matching the multi-step filters (category, type, salary range).
- [ ] Ensure empty results query returns status 200 with an empty collection.
- [ ] Implement `/api/v1/jobs/:id` returning segmented payloads (separating Overview, Requirements, and Detail sections).

### Frontend Developer
- [ ] Build **All Categories Directory** page showcasing available job classifications.
- [ ] Build the primary Job Board page with specific "No results" screens for filtered/empty directories.
- [ ] Implement the **Multi-Step Filter modal** (Filter feature drawer transitions 1, 2, and 3).
- [ ] Design the dynamic **Tabbed Job Detail Page** containing independent, click-navigable tabs:
    *   **Overview**
    *   **Requirements**
    *   **Details**
    *   **Apply**

---

## Week 4: Profiles, Modals & Error Boundary Views
**Goal:** User Profiles, administrative forms, logout confirmations, and custom HTTP error status pages.

### Backend Developer
- [ ] Implement SMTP notification mailers for Alumni approval/rejection.
- [ ] Build user profile edit actions and the secure, tokenized password-reset flows.
- [ ] Build administrative endpoints for system management (`/api/v1/admin/registrations/:id/approve` & `deny`).
- [ ] Configure global recovery middleware to catch panic states and return standardized JSON 500 formats.

### Frontend Developer
- [ ] Build the User Profile page (including visual presentation card widget layouts like Design B/B1 format, handled independently on the frontend).
- [ ] Create the **Logout Confirmation Modal Popup** overlay.
- [ ] Implement global error pages mapping strictly to designed assets:
    *   `403 Unauthorized / Forbidden`
    *   `404 Page Not Found`
    *   `429 Too Many Requests`
    *   `500 Internal Server Error`
- [ ] Build SAC and Sys Admin Portal Dashboards (aligned with Design C / Admin views).

---

## Week 5: Integration, End-to-End Testing & Deployment
**Goal:** Connect backend/frontend, run full QA, and deploy to IIS servers.

### Both Developers
- [ ] Disable MSW mocks in frontend and route Axios client requests directly to Go REST API.
- [ ] Perform cross-browser and responsive testing across simulated mobile, tablet, and web aspect ratios.
- [ ] Compile Go backend as Windows Service, configure environment values, and deploy under IIS **HttpPlatformHandler**.
- [ ] Build Next.js in Static Export mode (`npm run build`), extract output files to the IIS root web directory, and apply custom SPA dynamic URL Rewrite rules via `web.config`.
- [ ] Final sign-off.
