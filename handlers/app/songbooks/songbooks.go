package app_songbooks

import (
	"fmt"
	"net/http"
	"strconv"

	models "github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

func GetSongbooks(c *gin.Context) {
	songbooks, err := (&models.Songbook{}).GetAllSongbooksApp()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"songbooks": songbooks,
	})
}

func GetSongbookByID(c *gin.Context) {
	id := c.Param("id")

	songbookID, err := strconv.Atoi(id)
	if err != nil {
		panic("Invalid songbook ID")
	}

	songbook, err := (&models.Songbook{}).GetSongbookByID(songbookID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"songbook": songbook,
	})
}

func ExportSongbookByID(c *gin.Context) {
	id := c.Param("id")

	songbookID, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid songbook ID"})
		return
	}

	songbook, err := (&models.Songbook{}).GetSongbookByID(songbookID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	exportedData, err := songbook.ExportSongbookSQL()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// export data as sql script file
	c.FileAttachment(exportedData, fmt.Sprintf("songbook_%d_export.sql", songbookID))
}
