package api

import (
	"github.com/Rhaqim/creds/internal/models"
	"github.com/gin-gonic/gin"
)

func CreateOrganization(c *gin.Context) {
	var req models.Organization

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	user, err := GetUserFromToken(c)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	if err := req.CreateOrganization(user); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(200, req)
}
