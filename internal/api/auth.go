package api

import (
	"github.com/Rhaqim/creds/internal/authentication"
	err "github.com/Rhaqim/creds/internal/errors"
	"github.com/Rhaqim/creds/internal/models"
	"github.com/gin-gonic/gin"
)

func GetUserFromToken(c *gin.Context) (models.User, error) {
	check, ok := c.Get("user") // Check if user is logged in
	if !ok {
		return models.User{}, err.ErrUnauthorized
	}

	user := check.(models.User)

	return user, nil
}
func LoginHandler(c *gin.Context) {
	// Retrieve the dynamic value from the context
	provider := c.Param("provider")

	handler := authentication.NewOAuth2Adapter(provider)
	if handler == nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.ErrInvalidProvider,
		})
		return
	}

	handler.OAuth2LoginHandler(c)
}

func CallbackHandler(c *gin.Context) {
	provider := c.Param("provider")

	handler := authentication.NewOAuth2Adapter(provider)
	if handler == nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.ErrInvalidProvider,
		})
		return
	}

	handler.OAuth2CallbackHandler(c)
}

func LogoutHandler(c *gin.Context) {
	authentication.Logout(c)
	c.JSON(200, gin.H{"message": "Logged out successfully"})
}
