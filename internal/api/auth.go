package api

import (
	"github.com/Rhaqim/creds/internal/authentication"
	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	// Retrieve the dynamic value from the context
	provider := c.Param("provider")

	handler := authentication.NewOAuth2Adapter(provider)
	if handler == nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "Invalid provider, valid providers are google and github",
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
			"error": "Invalid provider, valid providers are google and github",
		})
		return
	}

	handler.OAuth2CallbackHandler(c)
}
