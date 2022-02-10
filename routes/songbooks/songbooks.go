package songbooks

import (
	"net/http"

	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

func GetSongbooks(c *gin.Context) {
	songbooks := new(models.Songbook).GetAllSongbooks()

	c.IndentedJSON(http.StatusOK, songbooks)
}

func GetSongbooksById(c *gin.Context) {
	id := c.Param("id")

	songs := new(models.Song).GetSongByID(id)

	if len(songs) == 0 {
		c.IndentedJSON(http.StatusNotFound, songs)
	}
	c.IndentedJSON(http.StatusNotFound, songs)

}
