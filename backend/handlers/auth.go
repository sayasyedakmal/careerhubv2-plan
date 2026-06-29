package handlers

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"careerhubv2-backend/config"
	"careerhubv2-backend/internal/jwks"
)

// jwksResolver is a package-level singleton initialised on first use.
var jwksResolver *jwks.Resolver

func getJWKSResolver() *jwks.Resolver {
	if jwksResolver == nil {
		jwksResolver = jwks.NewResolver(os.Getenv("MICROSOFT_ENTRA_JWKS_URL"))
	}
	return jwksResolver
}

// ─── Request / Response types ────────────────────────────────────────────────

type microsoftLoginRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ─── Microsoft ID token claims ────────────────────────────────────────────────

type microsoftClaims struct {
	OID               string `json:"oid"`
	Email             string `json:"email"`
	PreferredUsername string `json:"preferred_username"`
	Name              string `json:"name"`
	TenantID          string `json:"tid"`
	jwt.RegisteredClaims
}

// ─── CareerHub JWT claims ─────────────────────────────────────────────────────

type careerhubClaims struct {
	UserID      int    `json:"sub"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`
	jwt.RegisteredClaims
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func extractKID(tokenStr string) (string, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("malformed JWT")
	}
	headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", err
	}
	var header struct {
		Kid string `json:"kid"`
	}
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		return "", err
	}
	if header.Kid == "" {
		return "", fmt.Errorf("kid not found in JWT header")
	}
	return header.Kid, nil
}

func generateTokenPair(userID int, email, displayName, role string) (tokenResponse, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	accessExp := time.Now().Add(time.Hour)
	accessClaims := careerhubClaims{
		UserID:      userID,
		Email:       email,
		DisplayName: displayName,
		Role:        role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(secret)
	if err != nil {
		return tokenResponse{}, err
	}

	refreshExp := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := careerhubClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(secret)
	if err != nil {
		return tokenResponse{}, err
	}

	return tokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600,
	}, nil
}

// ─── POST /api/v1/auth/login/microsoft ───────────────────────────────────────

func LoginMicrosoft(c *gin.Context) {
	var req microsoftLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to decode request JSON payload",
			"code":  "INVALID_PAYLOAD",
		})
		return
	}

	// 1. Extract kid from JWT header
	kid, err := extractKID(req.IDToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Malformed Microsoft ID token",
			"code":  "INVALID_PAYLOAD",
		})
		return
	}

	// 2. Resolve RSA public key from JWKS
	resolver := getJWKSResolver()
	pubKey, err := resolver.GetKey(kid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Microsoft token signature verification failed: could not resolve public key",
			"code":  "INVALID_MICROSOFT_TOKEN",
		})
		return
	}

	// 3. Parse and verify token signature + claims
	msClaims := &microsoftClaims{}
	token, err := jwt.ParseWithClaims(req.IDToken, msClaims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return pubKey, nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "The Microsoft Entra ID token has expired",
				"code":  "MICROSOFT_TOKEN_EXPIRED",
			})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Microsoft token signature verification failed: invalid issuer or audience claim",
			"code":  "INVALID_MICROSOFT_TOKEN",
		})
		return
	}
	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Microsoft token is not valid",
			"code":  "INVALID_MICROSOFT_TOKEN",
		})
		return
	}

	// 4. Validate aud, iss, tid
	clientID := os.Getenv("MICROSOFT_ENTRA_CLIENT_ID")
	tenantID := os.Getenv("MICROSOFT_ENTRA_TENANT_ID")
	expectedIssuer := fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", tenantID)

	aud, _ := msClaims.GetAudience()
	iss, _ := msClaims.GetIssuer()

	audienceOK := false
	for _, a := range aud {
		if a == clientID {
			audienceOK = true
			break
		}
	}
	if !audienceOK || iss != expectedIssuer {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Microsoft token signature verification failed: invalid issuer or audience claim",
			"code":  "INVALID_MICROSOFT_TOKEN",
		})
		return
	}

	if msClaims.TenantID != tenantID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Authenticated Microsoft tenant is not authorized for this application",
			"code":  "INSUFFICIENT_PERMISSIONS",
		})
		return
	}

	// 5. Resolve email and display name from claims
	email := msClaims.Email
	if email == "" {
		email = msClaims.PreferredUsername
	}
	displayName := msClaims.Name
	oid := msClaims.OID

	// 6. Upsert user in database
	userID, role, err := upsertUser(oid, email, displayName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save auto-provisioned user details to database",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	// 7. Generate CareerHub token pair
	resp, err := generateTokenPair(userID, email, displayName, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate session tokens",
			"code":  "INTERNAL_SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// upsertUser creates or updates the user record and returns (userID, roleName, error).
func upsertUser(oid, email, displayName string) (int, string, error) {
	db := config.DB
	bootstrapAdmin := os.Getenv("BOOTSTRAP_ADMIN_EMAIL")

	// Determine user type and role based on domain rules
	userType := "External"
	registrationStatus := "Pending"
	roleName := ""
	var studentID *string

	emailLower := strings.ToLower(email)

	if strings.EqualFold(email, bootstrapAdmin) {
		userType = "SystemAdmin"
		registrationStatus = "N/A"
		roleName = "System Admin"
	} else if strings.HasSuffix(emailLower, "@student.uow.edu.my") {
		userType = "Student"
		registrationStatus = "N/A"
		roleName = "Student"
		parts := strings.Split(emailLower, "@")
		if len(parts) == 2 {
			sid := parts[0]
			studentID = &sid
		}
	} else if strings.HasSuffix(emailLower, "@uow.edu.my") {
		userType = "Staff"
		registrationStatus = "N/A"
	}

	// Check if user already exists
	var existingID int
	var existingDisplayName, existingEmail string
	err := db.QueryRow(
		`SELECT UserID, DisplayName, Email FROM Users WHERE MicrosoftObjectID = @p1`,
		oid,
	).Scan(&existingID, &existingDisplayName, &existingEmail)

	if err == nil {
		// User exists — sync display name and email if changed
		if existingDisplayName != displayName || existingEmail != email {
			_, err = db.Exec(
				`UPDATE Users SET DisplayName = @p1, Email = @p2 WHERE UserID = @p3`,
				displayName, email, existingID,
			)
			if err != nil {
				return 0, "", fmt.Errorf("failed to sync user profile: %w", err)
			}
		}
		// Fetch their current role
		var currentRole string
		_ = db.QueryRow(
			`SELECT r.RoleName FROM Roles r
			 INNER JOIN UserRoles ur ON r.RoleID = ur.RoleID
			 WHERE ur.UserID = @p1`,
			existingID,
		).Scan(&currentRole)
		return existingID, currentRole, nil
	}

	if err != sql.ErrNoRows {
		return 0, "", fmt.Errorf("failed to query user: %w", err)
	}

	// New user — insert
	var newUserID int
	err = db.QueryRow(
		`INSERT INTO Users (MicrosoftObjectID, Email, DisplayName, UserType, StudentID, RegistrationStatus)
		 OUTPUT INSERTED.UserID
		 VALUES (@p1, @p2, @p3, @p4, @p5, @p6)`,
		oid, email, displayName, userType, studentID, registrationStatus,
	).Scan(&newUserID)
	if err != nil {
		return 0, "", fmt.Errorf("failed to insert new user: %w", err)
	}

	// Assign role if applicable
	if roleName != "" {
		var roleID int
		err = db.QueryRow(`SELECT RoleID FROM Roles WHERE RoleName = @p1`, roleName).Scan(&roleID)
		if err == nil {
			_, _ = db.Exec(
				`INSERT INTO UserRoles (UserID, RoleID) VALUES (@p1, @p2)`,
				newUserID, roleID,
			)
		}
	}

	return newUserID, roleName, nil
}

// ─── POST /api/v1/auth/refresh ────────────────────────────────────────────────

func RefreshToken(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to decode refresh token JSON payload",
			"code":  "INVALID_PAYLOAD",
		})
		return
	}

	secret := []byte(os.Getenv("JWT_SECRET"))
	claims := &careerhubClaims{}

	token, err := jwt.ParseWithClaims(req.RefreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secret, nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Session has expired. Please sign in again.",
				"code":  "REFRESH_TOKEN_EXPIRED",
			})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "The provided refresh token is invalid or has been revoked",
			"code":  "INVALID_REFRESH_TOKEN",
		})
		return
	}
	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "The provided refresh token is invalid",
			"code":  "INVALID_REFRESH_TOKEN",
		})
		return
	}

	// Fetch user's current role from DB to embed in new access token
	var role string
	_ = config.DB.QueryRow(
		`SELECT r.RoleName FROM Roles r
		 INNER JOIN UserRoles ur ON r.RoleID = ur.RoleID
		 WHERE ur.UserID = @p1`,
		claims.UserID,
	).Scan(&role)

	resp, err := generateTokenPair(claims.UserID, claims.Email, claims.DisplayName, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate new session tokens",
			"code":  "INTERNAL_SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
