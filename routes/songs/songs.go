package songs

import (
	"net/http"

	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

func GetSongs(c *gin.Context) {
	songs := new(models.Song).GetAllSongs()

	c.IndentedJSON(http.StatusOK, songs)
}

func GetSongsById(c *gin.Context) {
	id := c.Param("id")

	songs := new(models.Song).GetSongByID(id)

	if len(songs) == 0 {
		c.IndentedJSON(http.StatusNotFound, songs)
	}
	c.IndentedJSON(http.StatusOK, songs)

}

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
