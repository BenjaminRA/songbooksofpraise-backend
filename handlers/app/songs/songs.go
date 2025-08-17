package app_songs

import (
	"net/http"
	"strconv"

	models "github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

func GetSongByID(c *gin.Context) {
	songIDStr := c.Param("song_id")

	if songIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Song ID is required"})
		return
	}

	songID, err := strconv.Atoi(songIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID format"})
		return
	}

	song, err := (&models.Song{}).GetSongByID(songID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"song": song,
	})
}
