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

	req.EncryptionKey = nil

	c.JSON(200, gin.H{
		"message":    "Credential created successfully",
		"credential": req,
	})
}

func GetCredential(c *gin.Context) {
	var cred models.Credential

	user, err := GetUserFromToken(c)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	id := c.Param("credId")

	resp, err := cred.FetchCredential(user, id)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error":   err.Error(),
			"message": "Error fetching credential.",
		})

		return
	}

	cred.EncryptionKey = nil

	c.JSON(200, gin.H{
		"credential": resp,
	})
}

func AddFields(c *gin.Context) {
	var req []models.CredentialField
	var cred models.Credential

	user, err := GetUserFromToken(c)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	id := c.Param("credId")

	if err := c.Bind(&req); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error":   err.Error(),
			"message": "Invalid request body.",
		})

		return
	}

	if err := cred.AddFields(user, req, id); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error":   err.Error(),
			"message": "Error adding fields.",
		})

		return
	}

	c.JSON(200, gin.H{
		"message":    "Fields added successfully",
		"credential": req,
	})
}
