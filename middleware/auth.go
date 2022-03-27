package middleware

import (
	"fmt"
	"net/http"

	"github.com/abiiranathan/acada/auth"
	"github.com/abiiranathan/acada/models"
	"github.com/abiiranathan/acada/services"
	"github.com/gin-gonic/gin"
)

// AuthUser middleware
// Extracts the jwt token from the request and
// verifies it, sets the user in the context
func LoginRequired(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized. No auth token provided"})
			c.Abort()
			return
		}

		// Validate jwt token
		// strip Bearer from token
		token = token[7:]

		userId, err := auth.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Unauthorized: %s", err)})
			c.Abort()
			return
		}

		// Get user from repository
		user, err := userService.GetUser(userId)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden. User not found!"})
			c.Abort()
			return
		}

		if !user.IsActive {
			c.JSON(http.StatusForbidden, gin.H{"error": "Your account has been deactivated"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

// Enforces that a superuser is logged in.
// Dependes on LoginRequired middleware being applied first.
//
// Apply login required middleware first
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := c.Get("user")

		if !ok {
			c.JSON(http.StatusUnauthorized, "Must be logged in first!")
			c.Abort()
			return
		}

		// Make sure that user is a user struct
		// type assertion
		u, ok := user.(models.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, "Must be logged in first!")
			c.Abort()
			return
		}

		if !u.IsAdmin {
			c.JSON(http.StatusForbidden, "Must be an admin!")
			c.Abort()
			return
		}

		c.Next()
	}
}
