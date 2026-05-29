# IIS Deployment Guide: CareerHubV2 Frontend

This guide is designed for the frontend developer to build and deploy the Next.js frontend as a static website on **IIS (Internet Information Services)**.

---

## 🛠️ Step 1: Verify Configuration
Before building, configure your `next.config.js` (at the root of the project) for Static HTML Export. Open the file and ensure it matches the configuration below:

```javascript
/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'export',          // Tells Next.js to export static HTML/CSS/JS
  images: {
    unoptimized: true,       // Required: IIS cannot optimize images on-the-fly
  },
};

module.exports = nextConfig;
```

---

## 📦 Step 2: Build the Frontend Locally
Run the build script on your development machine to compile your code:

```bash
npm run build
```

*   This will run the Next.js compiler and generate a folder named **`out`** at the project root. 
*   The `out` folder contains 100% pure static HTML, CSS, images, and JavaScript assets.

---

## 🚚 Step 3: Copy Files to the Windows Server
1.  Compress (Zip) the contents of the generated **`out`** folder.
2.  Log into your Windows Server.
3.  Create or navigate to your hosting directory, for example:  
    `C:\inetpub\wwwroot\careerhub-web`
4.  Extract the zipped files directly into this directory so that `index.html` resides at the root level of `careerhub-web`.

---

## ⚙️ Step 4: Add the Routing configuration (`web.config`)
Because this is a Single Page Application (SPA), we need to tell IIS to route deep URLs (like `/jobs/123`) back to `index.html` on browser refresh, letting React handle the nested paths.

1.  In your server deployment folder (`C:\inetpub\wwwroot\careerhub-web`), create a new text file named **`web.config`** (ensure it does not end in `.txt`).
2.  Paste the following XML configuration:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<configuration>
    <system.webServer>
        <rewrite>
            <rules>
                <!-- Redirects all requests that are not physical files or folders back to index.html -->
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

## 🖥️ Step 5: Configure IIS Manager
Now, tell IIS to serve this directory:

1.  Open **IIS Manager** on the Windows Server.
2.  Under **Sites**, click **Add Website...**
3.  Configure:
    *   **Site name**: `CareerHub Frontend`
    *   **Physical path**: `C:\inetpub\wwwroot\careerhub-web`
    *   **Binding/Port**: Port `80` (or `443` for HTTPS) with your host domain (e.g., `careerhub.uow.edu.my`).
4.  Ensure the **URL Rewrite Module** is installed on the IIS Server (otherwise, the rewrite rule in Step 4 will throw an error).
5.  Click **OK** and start the website!

---

## 💡 Developer Guidelines & Rules
*   **No Node.js Runtime Code**: Since there is no running Node.js server, you **cannot** use server-side runtime functions such as `getServerSideProps` or React Server Components that fetch data during server request-time.
*   **Data Fetching**: All data must be fetched client-side (in the browser) using standard libraries like `axios` or fetching utilities (`useEffect`, SWR, or React Query) querying the Go REST API.
