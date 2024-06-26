package api

import (
	"github.com/Rhaqim/creds/internal/models"
	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	var newFile models.CredentialFile

	user, err := GetUserFromToken(c)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error":   err.Error(),
			"message": "File not found in request.",
		})
		return
	}

	credID := c.Param("cred_id")
	format := c.Query("format")

	err = newFile.Process(user, file, credID, format)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error":   err.Error(),
			"message": "Error processing file.",
		})
		return
	}

	c.JSON(200, gin.H{"message": "File uploaded successfully"})
}
