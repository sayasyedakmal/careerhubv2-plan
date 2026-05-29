# Frontend Implementation Plan: CareerHubV2

This document provides the full technical specification for the React/Next.js User Interface based on the high-fidelity mobile and web design screens provided in `@frontend/design/`.

## 1. Tech Stack & Configuration
*   **Framework**: Next.js (using App Router, configured for **Static HTML Export** (`output: 'export'` in `next.config.js`))
*   **Language**: TypeScript
*   **Styling**: Tailwind CSS
*   **Icons**: Heroicons
*   **Authentication**: Hybrid (`@azure/msal-react` for Entra ID + Custom Context for local JWT storage)
*   **API Interceptor**: Axios with JWT automatic attachment
*   **Build Target**: Pure client-side compiled assets (`HTML`/`JS`/`CSS`) served directly by IIS. Developers **must not** use server-side runtime APIs (such as `getServerSideProps` or dynamic node-based routing) as there is no Node.js runtime in production. All state and data fetching must run inside the browser.

---

## 2. Global Styling & Layout Architecture
To align with the design guidelines and assets:
*   **Color Palette (UOW Malaysia Branding)**:
    *   Primary Navy: `#001C41`
    *   Primary Red: `#E21F26`
    *   Neutral Grays, White backgrounds, and subtle borders.
*   **Typography**: Montserrat (via Google Fonts).
*   **Responsiveness**: Mobile-first approach. All components must match the custom drawer headers, navigation elements, and layout flows seen in the `(Mobile-view)` designs.
*   **Shell Layout**:
    *   Mobile view contains a toggleable burger menu trigger opening a left/right slide-out **Drawer Navbar**.
    *   Navbar states: Default logged-out, Logged-in navigation, "Your Profile" navigation toggle, and a confirmation modal overlay for the **Logout** action.

---

## 3. High-Fidelity UI Screens & Component Mapping

Based on the design deliverables under `frontend/design/`, the following screens must be implemented exactly:

### A. Authentication & Home (Web & Mobile)
*   **Files**: `Homepage - sign in.png`, `Homepage - 401 Unauthorized.png`
*   **Implementation**:
    *   Dual login layout displaying a "Sign in with Microsoft" button (for Active Students/Staff) and a classic Username/Password form (for Alumni).
    *   When an unauthenticated user attempts to access protected routes, intercept and render the Homepage containing the **401 Unauthorized** error panel.

### B. Navigation & Interaction Overlays
*   **Files**: `Dashboard page Navbar.png`, `Dashboard page Navbar (Your Profile)).png`, `Dashboard page Navbar (Log out).png`, `Logout Confirmation popup.png`
*   **Implementation**:
    *   Mobile navigation uses a persistent top navbar with a hamburger menu and page title.
    *   Toggling menu slides open a side drawer. Under the logged-in state, the profile thumbnail opens a sub-level menu containing "Your Profile".
    *   Triggering log out presents a centered **Logout Confirmation Popup** (modal overlay) with customized gray and red confirmation buttons.

### C. Student Portal Views
*   **Files**: `All categories page.png`, `All Jobs page.png`, `Category-specific Jobs page.png`, `Filtered Jobs page.png`, `_Search_ Jobs page.png`
    *   **Category Browse**: An interactive page showing card blocks for all available job categories (e.g., Computing, Engineering, Business) with icon assets.
    *   **Empty States**: Custom "No results" screens for searches, filter states, and category selections (`All Jobs page (No results).png`, `Jobs Category page (No results).png`, `_Search_ Jobs page (No results).png`).
*   **Multi-Step Filter System**:
    *   **Files**: `Filter feature.png`, `Filter feature1.png`, `Filter feature2.png`, `Filter feature3.png`
    *   **Implementation**: Clicking "Filter" on the job list opens a multi-step slide-up modal (Filter drawer). The user can configure filters step-by-step (e.g., job type, location, experience levels) with clean toggle chips and validation.
*   **Tabbed Job Detail Page**:
    *   **Files**: `Job's Detail page_Overview.png`, `Job's Detail page_Requirements.png`, `Job's Detail page_Details.png`, `Job's Detail page_Apply.png`
    *   **Implementation**: Refactor the job detail page into a structured view with a sticky sub-navbar containing 4 distinct interactive tabs:
        1.  **Overview**: Title, company, salary band, location, and metadata.
        2.  **Requirements**: Grid or list of qualifications and skills.
        3.  **Details**: Full narrative of the job duties and company profile.
        4.  **Apply**: Actionable panel to submit application, upload CV, or redirect to external portals.

### D. Profiles & Identity
*   **Files**: `Your Profile page.png`, `Namecard_Design B.png`, `Namecard_Design B1.png`
*   **Implementation**:
    *   Implement user details editing forms.
    *   Render the student **Namecard** card widget. This card behaves as a digital badge, styled with custom colors, showcasing student major, ID, registration status, and profile photo.

### E. Portal Dashboards
*   **SAC Portal**:
    *   **Files**: `SAC Dashboard page Navbar.png`, `SAC Dashboard page2.png`, `SAC Dashboard page_Design C.png`, `SAC Your Profile page.png`
    *   **Implementation**: A professional tabular list layout tracking job listings, approval queues, active counters, and statuses. Fits within SAC Drawer Navbar.
*   **System Admin Portal**:
    *   **Files**: `Sys Admin Dashboard page.png`, `Sys Admin-Your Profile page.png`
    *   **Implementation**: Management panels for group assignments, roles, and system monitoring.

### F. Global Error Boundary Screens
*   **Files**: `Error page (403 Unauthorized).png`, `Error page (404 Not Found).png`, `Error page (429 Too Many Requests).png`, `Error page (500 Internal Server Error).png`
*   **Implementation**:
    *   Create dedicated routing templates/pages matching the designated mobile graphics for each of these HTTP response conditions. Each page must feature localized action links (e.g., "Return to Dashboard", "Try Again").

---

## 4. API Integration & State Management
*   **MSW (Mock Service Worker)**: Configure MSW handlers to return paginated search data, empty arrays (to test the "No results" screens), list of categories, and fake JWT structures.
*   **State Management**: Use React local state and context wrappers to manage multi-step filter configurations and currently selected job detail tab indices.

---

## 5. Build & IIS Static Deployment Specifications
To compile the site for direct, ultra-high-speed static delivery from an IIS directory:
1.  **Configuration**: Enable static output in `next.config.js`:
    ```javascript
    module.exports = {
      output: 'export',
      images: {
        unoptimized: true, // Required for static export
      },
    }
    ```
2.  **Compilation**: Run `npm run build`. This generates an `out/` directory containing purely static files.
3.  **Deployment**: Copy the contents of `out/` directly into the designated IIS physical path (e.g., `C:\inetpub\wwwroot\careerhub`).
4.  **Client-Side Routing**: Include an IIS `web.config` rewrite rule (detailed in `docs/iis-deployment.md`) to redirect dynamic route requests back to `index.html`, letting the client-side router handle nested paths seamlessly.
