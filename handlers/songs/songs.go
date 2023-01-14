package handlers

import (
	"net/http"

	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

func GetMusicSheet(c *gin.Context) {
	id := c.Param("id")
	data, filename, err := new(models.Song).GetMusicSheet(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, err)
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/octet-stream", data)
}

func GetMusic(c *gin.Context) {
	id := c.Param("id")
	data, filename, err := new(models.Song).GetMusic(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, err)
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/octet-stream", data)
}

func GetMusicOnly(c *gin.Context) {
	id := c.Param("id")
	data, filename, err := new(models.Song).GetMusicOnly(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, err)
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/octet-stream", data)
}

func GetVoicesByVoice(c *gin.Context) {
	id := c.Param("id")
	voice := c.Param("voice")
	data, filename, err := new(models.Song).GetVoice(id, voice)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, err)
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/octet-stream", data)
}
