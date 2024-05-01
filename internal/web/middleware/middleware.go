package middleware

import (
	"net/http"

	"github.com/Rhaqim/creds/internal/authentication"
	"github.com/Rhaqim/creds/internal/models"
	"github.com/gin-gonic/gin"
)

func AuthGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var user models.User
		var token string

		var jwt *authentication.JWT = authentication.NewJWT()

		// get token from cookie
		token, err = c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

			return
		}

		// decode token
		claims, err := jwt.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

			return
		}

		// get user by email
		err = user.GetByEmail(claims.Email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// set user in context
		c.Set("user", user)

		c.Next()
	}
}
