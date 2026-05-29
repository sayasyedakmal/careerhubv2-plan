# IIS Deployment Guide: CareerHubV2

This guide covers the manual deployment of the CareerHubV2 monorepo to a Windows Server running IIS.

## Prerequisites
Before starting, ensure the following are installed on the server:
1. **[HttpPlatformHandler v1.2+](https://www.microsoft.com/en-us/download/details.aspx?id=49493)** (For Go Backend)
2. **[URL Rewrite Module](https://www.iis.net/downloads/microsoft/url-rewrite)** (For Frontend client-side SPA routing)

---

## 1. Backend Deployment (Go)

1. Build your Go app: `go build -o careerhub-api.exe`.
2. Create a folder: `C:\inetpub\wwwroot\careerhub-api`.
3. Place the `.exe` and the following `web.config` in that folder.

### `web.config` for Go
```xml
<?xml version="1.0" encoding="UTF-8"?>
<configuration>
    <system.webServer>
        <handlers>
            <add name="httpplatformhandler" path="*" verb="*" modules="httpPlatformHandler" resourceType="Unspecified" />
        </handlers>
        <httpPlatform processPath=".\careerhub-api.exe" 
                      arguments="" 
                      stdoutLogEnabled="true" 
                      stdoutLogFile=".\logs\stdout.log" 
                      startupTimeLimit="60">
            <environmentVariables>
                <environmentVariable name="PORT" value="%HTTP_PLATFORM_PORT%" />
                <!-- 
                SECURITY WARNING: Do not store plain text passwords here.
                Instead, read from a `.env` file secured by NTFS permissions, 
                or set these as System Environment Variables at the OS level.
                -->
            </environmentVariables>
        </httpPlatform>
    </system.webServer>
</configuration>
```

---

## 2. Frontend Deployment (Static HTML Export)

To maximize performance, stability, and eliminate runtime overhead, the frontend compiles into purely static assets served natively by IIS.

1. Configure `next.config.js` to include `output: 'export'` and `images: { unoptimized: true }`.
2. Build the Next.js app: `npm run build`. This generates static HTML, CSS, and JS output in the `out/` folder.
3. Create a folder: `C:\inetpub\wwwroot\careerhub-web`.
4. Copy the entire contents of the `out/` directory into `C:\inetpub\wwwroot\careerhub-web`.
5. Add the following `web.config` to the folder. This ensures IIS rewrites dynamic routes (like `/jobs/12` or `/profile`) back to `index.html` on refresh, allowing client-side routing to function correctly.

### `web.config` for Frontend (SPA Routing)
```xml
<?xml version="1.0" encoding="UTF-8"?>
<configuration>
    <system.webServer>
        <rewrite>
            <rules>
                <rule name="SPA Client-Side Routing" stopProcessing="true">
                    <match url=".*" />
                    <conditions logicalGrouping="MatchAll">
                        <add input="{REQUEST_FILENAME}" matchType="IsFile" negate="true" />
                        <add input="{REQUEST_DIR}" matchType="IsDirectory" negate="true" />
                    </conditions>
                    <action type="Rewrite" url="index.html" />
                </rule>
            </rules>
        </rewrite>
    </system.webServer>
</configuration>
```

---

## 3. IIS Manager Configuration

1. **Application Pools**:
   - Create a new App Pool for the Backend.
   - Set **.NET CLR Version** to **"No Managed Code"**.
2. **Websites**:
   - Create a site for the Frontend pointing to `C:\inetpub\wwwroot\careerhub-web`.
   - Create a site (or sub-application) for the Backend pointing to `C:\inetpub\wwwroot\careerhub-api`.
3. **Permissions**:
   - Ensure the `IIS_IUSRS` group has **Read & Execute** permissions on both folders.
   - For the Backend, `IIS_IUSRS` also needs **Write** permission for the `\logs` folder if you enabled logging.
