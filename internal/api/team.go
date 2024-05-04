package api

import (
	"github.com/Rhaqim/creds/internal/models"
	"github.com/gin-gonic/gin"
)

func CreateTeam(c *gin.Context) {
	var req models.OrganizationTeam

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

	if err := req.CreateTeam(user); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error":   err.Error(),
			"message": "Error creating team.",
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "Team created successfully.",
		"team":    req,
	})

}

func AddMember(c *gin.Context) {
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

func GetMembers(c *gin.Context) {
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

func GetMember(c *gin.Context) {
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

func InviteMember(c *gin.Context) {
	var mem models.OrganizationMember

	user, err := GetUserFromToken(c)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	email := c.Query("email")

	err = c.Bind(&mem)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error":   err.Error(),
			"message": "Invalid request body.",
		})

		return
	}

	if err := mem.InviteMember(user, email); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error":   err.Error(),
			"message": "Error inviting member.",
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "Invitation sent successfully.",
	})
}
