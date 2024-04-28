package api

import (
	"github.com/Rhaqim/creds/internal/models"
	"github.com/gin-gonic/gin"
)

func CreateOrganization(c *gin.Context) {
	var req models.Organization

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

	if err := req.CreateOrganization(user); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error":   err.Error(),
			"message": "Error creating organization.",
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "Organization created successfully.",
		"org":     req,
	})
}
