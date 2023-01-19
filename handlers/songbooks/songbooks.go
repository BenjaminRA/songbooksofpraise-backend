package songbooks

import (
	"net/http"

	models "github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

func VerifySongbook(c *gin.Context) {
	id := c.Param("id")

	if err := models.SetSongbookVerified(id, true); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"songbook": id,
	})
}
