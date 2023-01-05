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
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
	filenameArray := strings.Split(file.Filename, ".")
	secret := helpers.GetSecretString()

	path := fmt.Sprintf("tmp/%s.%s", secret, filenameArray[len(filenameArray)-1])

	c.SaveUploadedFile(file, path)
	c.IndentedJSON(http.StatusOK, map[string]string{
		"path": path,
	})
}
