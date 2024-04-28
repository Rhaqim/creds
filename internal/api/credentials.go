package api

import (
	"github.com/Rhaqim/creds/internal/models"
	"github.com/gin-gonic/gin"
)

func CreateCrendentials(c *gin.Context) {
	var req models.Credential

	user, err := GetUserFromToken(c)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	if err := c.Bind(&req); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error":   err.Error(),
			"message": "Invalid request body.",
		})

		return
	}

	if err := req.CreateCredential(user); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error":   err.Error(),
			"message": "Error creating credential.",
		})

		return
	}

	c.JSON(200, gin.H{
		"message":    "Credential created successfully",
		"credential": req,
	})
}
