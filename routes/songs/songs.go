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
	c.IndentedJSON(http.StatusNotFound, songs)

}
