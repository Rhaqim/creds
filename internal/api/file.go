package api

import (
	"fmt"

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

	fmt.Println(file.Filename)
	fmt.Println(file.Size)
	fmt.Println(file.Header)
	fmt.Println(file.Header.Get("Content-Type"))

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

	err = newFile.Insert()
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "File uploaded successfully"})
}
