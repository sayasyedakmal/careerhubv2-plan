# Frontend Implementation Plan: CareerHubV2

This document provides the full technical plan for the React-based User Interface.

## 1. Tech Stack
- **Framework**: React 18+ (Vite)
- **Language**: TypeScript
- **State Management**: [TanStack Query](https://tanstack.com/query/latest) (React Query)
- **Auth**: `@azure/msal-react` (Microsoft Entra ID)
- **Styling**: Vanilla CSS (Modern CSS variables, Flexbox/Grid)

## 2. Branding & Design (UOW Malaysia)
Follow the official branding guidelines:
- **Colors**: Navy (`#001C41`), Red (`#E21F26`).
- **Typography**: Montserrat (via Google Fonts).
- **Aesthetic**: Modern, clean, professional. Use glassmorphism for cards and subtle animations.

## 3. Core Development Tasks

### Authentication
- Setup `MsalProvider` in `App.tsx`.
- Implement a "Login" button that triggers the redirect flow.
- Ensure the app only shows content to authenticated users.

### API Layer & Mocking (MSW)
- **Mocking**: Before the backend is ready, use **MSW (Mock Service Worker)** to simulate the API. 
- Follow the instructions in `docs/auth-process.md` for MSW setup.
- Refer to `shared/openapi.yaml` for data structures.

### Key Components
- `Navbar`: Login/Logout status and Branding.
- `JobSearch`: Search bar and category filters.
- `JobCard`: Summary view of a job listing.
- `JobDetail`: Full-page view for the job description.

## 4. Deployment (IIS)
- Build the app: `npm run build`.
- Use the `web.config` from `docs/iis-deployment.md`.
- Install the **URL Rewrite Module** on the IIS server.
- Copy the `dist` folder to the server root.

---

## 5. How to Integrate with Backend
To ensure you and the backend developer are in sync:

1.  **Shared Contract**: Always refer to `shared/openapi.yaml` as your source of truth for API responses.
2.  **Environment**: Use `shared/.env.template` to share the `AZURE_CLIENT_ID` and `AZURE_TENANT_ID`.
3.  **API URL**: Set your API base URL to the Go server (usually `http://localhost:8080/api/v1`).
4.  **Authorization**: You MUST attach the Access Token to every API request using the `Authorization: Bearer <token>` header.
