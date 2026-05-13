# IIS Deployment Guide: CareerHubV2

This guide covers the manual deployment of the CareerHubV2 monorepo to a Windows Server running IIS.

## Prerequisites
Before starting, ensure the following are installed on the server:
1. **[HttpPlatformHandler v1.2+](https://www.microsoft.com/en-us/download/details.aspx?id=49493)** (For Go Backend)
2. **[URL Rewrite Module](https://www.iis.net/downloads/microsoft/url-rewrite)** (For React Frontend)

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
                <environmentVariable name="DB_CONNECTION_STRING" value="your_real_connection_string" />
            </environmentVariables>
        </httpPlatform>
    </system.webServer>
</configuration>
```

---

## 2. Frontend Deployment (React)

1. Build your React app: `npm run build`.
2. Create a folder: `C:\inetpub\wwwroot\careerhub-web`.
3. Copy the contents of the `dist` folder into that folder.
4. Add the following `web.config` to handle client-side routing.

### `web.config` for React
```xml
<?xml version="1.0" encoding="UTF-8"?>
<configuration>
    <system.webServer>
        <rewrite>
            <rules>
                <rule name="React Routes" stopProcessing="true">
                    <match url=".*" />
                    <conditions logicalGrouping="MatchAll">
                        <add input="{REQUEST_FILENAME}" matchType="IsFile" negate="true" />
                        <add input="{REQUEST_FILENAME}" matchType="IsDirectory" negate="true" />
                    </conditions>
                    <action type="Rewrite" url="./index.html" />
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
