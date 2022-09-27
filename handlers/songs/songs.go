package handlers

import (
	"net/http"

	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

// func GetSongs(c *gin.Context) {
// 	songs := new(models.Song).GetAllSongs()

// 	c.IndentedJSON(http.StatusOK, songs)
// }

// func GetSongsById(c *gin.Context) {
// 	id := c.Param("id")

// 	song := new(models.Song).GetSongByID(id)

// 	if song.ID.Hex() == "000000000000000000000000" {
// 		c.IndentedJSON(http.StatusNotFound, song)
// 	} else {
// 		c.IndentedJSON(http.StatusOK, song)
// 	}

// }

func GetMusicSheet(c *gin.Context) {
	id := c.Param("id")
	data, filename := new(models.Song).GetMusicSheet(id)

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/octet-stream", data)
}

func GetVoicesByVoice(c *gin.Context) {
	id := c.Param("id")
	voice := c.Param("voice")
	data, filename, err := new(models.Song).GetVoice(id, voice)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, err)
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/octet-stream", data)
}
