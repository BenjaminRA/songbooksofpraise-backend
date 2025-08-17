package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/BenjaminRA/himnario-backend/helpers"
	"github.com/gin-gonic/gin"
)

func PostFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		panic(err)
	}

	song_id, _ := c.GetQuery("song_id")
	file_type, _ := c.GetQuery("file_type")

	filenameArray := strings.Split(file.Filename, ".")
	filename := filenameArray[len(filenameArray)-2]
	fileExtension := filenameArray[len(filenameArray)-1]

	path := "tmp"

	if song_id != "" {
		path = path + "/songs/" + song_id
		filename = song_id
	}

	if file_type != "" {
		// If the file type is voices, we save the files in the same folder and work with the filename
		// to differentiate them
		if strings.Contains(file_type, "_voice_audio_file") {
			path = path + "/voices"
			filename = song_id + "_" + strings.ReplaceAll(file_type, "_voice_audio_file", "")
		} else {
			path = path + "/" + file_type
		}

	}

	path = fmt.Sprintf("%s/%s-%s.%s", path, filename, helpers.GetUniqueMD5ID(), fileExtension)

	// Create directories if they don't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory: " + err.Error()})
		panic(err)
	}

	// Now save the file
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file: " + err.Error()})
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"path": path,
	})
}
