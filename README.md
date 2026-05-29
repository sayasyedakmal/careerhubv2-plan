# CareerHubV2

CareerHubV2 is a unified job board and registration workspace for current active students, alumni, SAC (Student Affairs Centre) department staff, and system administrators at UOW Malaysia.

This repository is organized as a monorepo containing:
*   `backend/` - High-performance REST API built with Go (Gin, SQL Server, JWT).
*   `frontend/` - Responsive, interactive Next.js + Tailwind CSS application.
*   `shared/` - Shared contracts, including the OpenAPI specification.
*   `docs/` - System architecture, authentication flows, and implementation guides.

---

## 🎨 Design System & High-Fidelity Deliverables
The frontend user experience maps directly to high-fidelity mobile and desktop designs located in:
📂 **`frontend/design/`**

### UI Layouts Included
1.  **Careerhub Login & Landing Homepage**: Dual-authentication page views (Entra ID + local sign-in) with native `401 Unauthorized` states.
2.  **Student Portal**: 
    *   **All Categories Directory** page.
    *   Advanced **multi-step filter** workflows.
    *   Interactive, tabbed **Job Detail view** split into `Overview`, `Requirements`, `Details`, and `Apply`.
    *   Student Profile page utilizing the custom digital **Namecard Widget** (Design B/B1 format).
    *   Mobile drawer navigators and logout dialog modals.
3.  **Department Portals**:
    *   **SAC Dashboard** & Your Profile page.
    *   **System Admin Dashboard** & Your Profile page.
4.  **Error States**: Calibrated templates for `403 (Forbidden)`, `404 (Not Found)`, `429 (Too Many Requests)`, and `500 (Internal Server Error)`.

---

## 🛠️ Project Structure
```text
careerhubv2-plan/
├── backend/                  # Go REST API source code
│   └── PLAN.md               # Backend implementation tasks & endpoints
├── frontend/                 # Next.js web application
│   ├── design/               # High-fidelity developer screen designs (.png)
│   └── PLAN.md               # Frontend page routes & component definitions
├── shared/                   # Monorepo OpenAPI specifications & configs
│   └── openapi.yaml          # OpenAPI Specification (Source of Truth)
└── docs/                     # System architecture & guidelines
    ├── my-plan.md            # Overall MVP specifications
    ├── weekly_milestone.md   # Concurrent 5-week milestone targets
    ├── auth-process.md       # Hybrid Authentication strategy (Entra ID + Local JWT)
    ├── db-schema.md          # SQL Server relational database schema
    ├── roles-permissions-matrix.md   # RBAC permissions matrix
    └── iis-deployment.md     # Production deployment instructions on IIS
```

---

## 🚀 How to Navigate & Run
1.  **Overall Plan**: Check out `docs/my-plan.md` for a summary of features and user roles.
2.  **Implementation Goals**: Review `docs/weekly_milestone.md` to track current development progress for both backend and frontend.
3.  **Contracts**: Follow `shared/openapi.yaml` when writing API endpoints and mock handlers.
