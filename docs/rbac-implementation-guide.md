# RBAC Implementation Guide for Developers

Welcome! If you are new to Role-Based Access Control (RBAC), this guide will walk you step-by-step through exactly how we secure the CareerHubV2 application. 

At its core, RBAC is simply asking: **"Does this user have the correct permission string to do this action?"** If they do, allow it. If they don't, block it.

---

## Step 1: The Database Layer

Before anyone logs in, we need a list of rules. We have mapped out exactly what each role can do in the `docs/roles-permissions-matrix.md` file. 

During the backend database setup, you will create a seeding script that inserts these rules into the SQL Server database. For example, it will insert a rule saying that the "System Admin" role contains the `"create_job"` permission string, but the "Student" role does not.

---

## Step 2: The Login Process (Backend)

The most important part of RBAC happens the exact moment a user logs in. We use a **JSON Web Token (JWT)** to store the user's permissions.

When a user successfully logs in (either via Microsoft or Email/Password), the Go backend does the following:

1. Finds the `UserID` in the database.
2. Looks up the database to find what permissions belong to that user.
3. Embeds those permissions inside the JWT payload (claims).
4. Signs the JWT cryptographically using a secret key so no one can fake it.

**Example Backend Code (Go):**
```go
// 1. Fetch user's role and permissions from SQL Server
roleName := db.GetUserRole(userID) // e.g., "SAC Department"
permissionsArray := db.GetUserPermissions(userID) // e.g., ["view_jobs", "create_job", "edit_job"]

// 2. Define the JWT "Claims" (The data embedded inside the token)
claims := jwt.MapClaims{
    "sub":         userID,
    "role":        roleName,
    "permissions": permissionsArray, // We put the array directly inside the token!
    "exp":         time.Now().Add(time.Minute * 15).Unix(), // Expires in 15 mins
}

// 3. Create and sign the token
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
signedTokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

// 4. Return this string to the React frontend
return signedTokenString
```

---

## Step 3: Securing the API (Backend)

Now that the user has a token, they will send it with every API request (like trying to delete a job). The backend needs to read the token and block unauthorized requests.

You will write a **Middleware Guard** in Go that runs before the actual API endpoint:

```go
// The RBAC Middleware Guard
func RequirePermission(requiredPermission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Read the user's permissions that were extracted from their JWT token
        userPermissions := c.GetStringSlice("user_permissions")

        // Loop through their permissions to see if they have the required one
        hasPermission := false
        for _, p := range userPermissions {
            if p == requiredPermission {
                hasPermission = true
                break
            }
        }

        if !hasPermission {
            // Block the request! Send a 403 Forbidden error.
            c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You lack the required permission"})
            c.Abort()
            return
        }

        // If they have it, let the request proceed successfully
        c.Next()
    }
}
```

**How to use it on your routes:**
```go
// Only users with "create_job" permission can access this route
router.POST("/jobs", RequireAuth(), RequirePermission("create_job"), CreateJobHandler)
```

---

## Step 4: Securing the UI (Frontend)

The frontend's job is to provide a good User Experience. If a Student doesn't have permission to delete a job, the frontend should **hide the delete button** so they don't click it and get an error.

When the React frontend receives the JWT from the login API, it can decode the token to read the JSON payload. It saves the `permissions` array into React Context.

You will create a helper component to hide/show UI elements:

```tsx
// A Helper Component to conditionally render UI
const HasPermission = ({ required, children }) => {
    const { permissions } = useAuth(); // Gets the current user's permissions array
    
    // If they have the required string, show the button
    if (permissions.includes(required)) {
        return children;
    }
    
    // Otherwise, render absolutely nothing (hides the button)
    return null;
}
```

**How to use it in your React pages:**
```tsx
function JobDetail({ job }) {
    return (
        <div>
            <h1>{job.title}</h1>
            
            {/* ONLY SAC and System Admins will see this Delete button */}
            <HasPermission required="delete_job">
                <button className="bg-red-500">Delete Posting</button>
            </HasPermission>
        </div>
    )
}
```

---

## Step 5: How Frontend Developers Test This (Mocking)

Since the frontend and backend developers are building at the same time, the frontend developer can't wait for the backend login API to be finished. They will use **MSW (Mock Service Worker)** to fake the login process.

The frontend developer can use a site like [jwt.io](https://jwt.io) to generate a few fake, hardcoded tokens. One token will contain Student permissions, and another will contain Admin permissions.

**Example MSW Mock Handler (`src/mocks/handlers.js`):**
```javascript
import { http, HttpResponse } from 'msw'

// Hardcoded fake tokens
const MOCK_STUDENT_TOKEN = "eyJhbG... (student token)";
const MOCK_SAC_TOKEN = "eyJhbG... (sac token)";

export const handlers = [
  // Intercept the Login request
  http.post('http://localhost:8080/api/v1/auth/login', async ({ request }) => {
    const body = await request.json();
    
    // Test Scenario: If the developer types 'sac@uow.edu.my' into the login screen, return the SAC token
    if (body.email === 'sac@uow.edu.my') {
      return HttpResponse.json({ access_token: MOCK_SAC_TOKEN });
    }

    // Default: Log in as a standard Student
    return HttpResponse.json({ access_token: MOCK_STUDENT_TOKEN });
  }),
]
```

By typing different emails into the login screen, the frontend developer receives different fake tokens, allowing them to test if buttons correctly hide and show based on the user's role!
