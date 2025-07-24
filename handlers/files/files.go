package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/BenjaminRA/himnario-backend/helpers"
	"github.com/gin-gonic/gin"
)

func PostFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	filenameArray := strings.Split(file.Filename, ".")
	secret := helpers.GetSecretString()

	path := fmt.Sprintf("tmp/%s.%s", secret, filenameArray[len(filenameArray)-1])

	c.SaveUploadedFile(file, path)
	c.JSON(http.StatusOK, gin.H{
		"path": path,
	})
}
