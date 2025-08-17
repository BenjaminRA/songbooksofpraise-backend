package songs

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

func UpdateSong(c *gin.Context) {
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

	var updatedSong models.Song
	if err := c.ShouldBindJSON(&updatedSong); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	// Ensure the ID from the URL is used
	updatedSong.ID = songID

	if err := updatedSong.UpdateSong(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Song updated successfully",
		"song":    updatedSong,
	})
}

func CreateSong(c *gin.Context) {
	var newSong models.Song
	if err := c.ShouldBindJSON(&newSong); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	if err := newSong.CreateSong(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create song: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Song created successfully",
		"song":    newSong,
	})
}

func DeleteSong(c *gin.Context) {
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

	if err := song.DeleteSong(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Song deleted successfully",
	})
}
