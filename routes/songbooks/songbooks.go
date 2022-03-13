package songbooks

import (
	"net/http"

	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

func GetSongbooks(c *gin.Context) {
	lang := c.GetHeader("Language")
	if lang == "" {
		lang = "EN"
	}

	songbooks := new(models.Songbook).GetAllSongbooks(lang)

	c.IndentedJSON(http.StatusOK, songbooks)
}

func GetSongbooksById(c *gin.Context) {
	id := c.Param("id")
	lang := c.GetHeader("Language")
	if lang == "" {
		lang = "EN"
	}

	song := new(models.Songbook).GetSongbookByID(id, lang)

	if song.ID.Hex() == "000000000000000000000000" {
		c.IndentedJSON(http.StatusNotFound, song)
	} else {
		c.IndentedJSON(http.StatusOK, song)
	}
}

func PostSongbook(c *gin.Context) {
	var songbook models.Songbook
	lang := c.GetHeader("Language")
	if lang == "" {
		lang = "EN"
	}

	if err := c.BindJSON(&songbook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}

	songbook, err := songbook.CreateSongbook(songbook, lang)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.IndentedJSON(http.StatusCreated, songbook)
}

func DeleteSongbook(c *gin.Context) {
	id := c.Param("id")
	lang := c.GetHeader("Language")
	if lang == "" {
		lang = "EN"
	}

	songbook := new(models.Songbook).GetSongbookByID(id, lang)

	if songbook.ID.Hex() == "000000000000000000000000" {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "Songbook not found",
		})
	} else {
		if err := songbook.DeleteSongbook(); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}

		c.IndentedJSON(http.StatusOK, songbook)
	}
}

func UpdateSongbook(c *gin.Context) {
	id := c.Param("id")
	lang := c.GetHeader("Language")
	if lang == "" {
		lang = "EN"
	}

	songbook := new(models.Songbook).GetSongbookByID(id, lang)

	if songbook.ID.Hex() == "000000000000000000000000" {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "Songbook not found",
		})
		return
	}

	if err := c.BindJSON(&songbook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := songbook.UpdateSongbook(); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, songbook)
}
