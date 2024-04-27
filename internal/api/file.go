package api

import (
	"strconv"

	"github.com/Rhaqim/creds/internal/models"
	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	var newFile models.CredentialFile

	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	credID := c.Param("id")

	filedata := make([]byte, file.Size)

	// Open the file
	src, err := file.Open()
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	// Read the file
	_, err = src.Read(filedata)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	newFile.FileName = file.Filename
	newFile.FileSize = file.Size
	newFile.FileFormat = models.JSON
	newFile.FileData = filedata

	id, err := strconv.ParseUint(credID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid credential id"})
		return
	}

	err = newFile.Save(uint(id))
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "File uploaded successfully"})
}
