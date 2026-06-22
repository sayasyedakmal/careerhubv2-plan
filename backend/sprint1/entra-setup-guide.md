# Microsoft Entra ID Application Registration Guide

This guide provides step-by-step instructions for registering CareerHubV2 in the Microsoft Entra admin center to obtain the **Application (Client) ID** and **Directory (Tenant) ID** required for MSAL authentication.

---

## Step 1: Sign In to the Admin Center
1. Open your web browser and go to the [Microsoft Entra admin center](https://entra.microsoft.com/) (or the [Azure Portal](https://portal.azure.com/)).
2. Log in using your developer credentials or university administrator account.

---

## Step 2: Navigate to App Registrations
1. In the left-hand sidebar navigation, click to expand **Identity**.
2. Select **Applications** > **App registrations**.
3. Click the **`+ New registration`** button at the top toolbar.

```
  Identity
  ├── Users
  ├── Groups
  └── Applications
      └── App registrations  <── Click here, then click "+ New registration"
```

---

## Step 3: Register the Application
On the **Register an application** form, configure the following settings:

1. **Name:** Enter a user-facing name for the app (e.g. `CareerHubV2-Portal`).
2. **Supported account types:**
   * Select **"Accounts in this organizational directory only (Single tenant)"** (Recommended for university portals to restrict access solely to students/staff with university emails like `@uow.edu.my`).
3. **Redirect URI:**
   * In the dropdown, select **Single-page application (SPA)** (since the Next.js/React frontend will run MSAL).
   * In the text box, enter your local frontend URL (e.g. `http://localhost:3000` or `http://localhost:3000/redirect`).
4. Click **Register** at the bottom of the page.

---

## Step 4: Retrieve Your Client & Tenant IDs
Once the registration is complete, Microsoft will redirect you to the application's **Overview** page. Look under the **Essentials** section at the top of the page.

Copy the values for these two fields:

```
┌─────────────────────────────────────────────────────────────┐
│ Essentials                                                  │
│                                                             │
│  Application (client) ID  :  [ 1a2b3c4d-5e6f-7a8b-9c0d... ] │ ◄── AZURE_AD_CLIENT_ID
│  Directory (tenant) ID    :  [ 9f8e7d6c-5b4a-3f2e-1d0c... ] │ ◄── AZURE_AD_TENANT_ID
│  Object ID                :  [ 8c7b6a5d-4e3f-2a1b-0c9d... ] │
└─────────────────────────────────────────────────────────────┘
```

---

## Step 5: Configure your Backend Environment File
Create or update your `backend/.env` file with these values:

```bash
# Microsoft Entra ID Authentication Settings
AZURE_AD_CLIENT_ID="your-copied-application-client-id"
AZURE_AD_TENANT_ID="your-copied-directory-tenant-id"
AZURE_AD_JWKS_URL="https://login.microsoftonline.com/your-copied-directory-tenant-id/discovery/v2.0/keys"
```

---

## Step 6: Configure Platform Settings (Authentication)
If your colleague needs to add additional redirect URIs later (e.g. staging or production URLs):
1. In the left menu of your registered application, select **Authentication** under the **Manage** section.
2. Under **Single-page application**, click **Add URI** and enter any additional frontend addresses.
3. *Note on Implicit Flow:* If the frontend uses modern **MSAL v2/v3** (Auth Code Flow with PKCE), leave the checkboxes under "Implicit grant and hybrid flows" unchecked.
4. Click **Save** at the top of the page.
