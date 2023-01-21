package songbooks

import (
	"net/http"

	"github.com/BenjaminRA/himnario-backend/email"
	models "github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

func VerifySongbook(c *gin.Context) {
	id := c.Param("id")

	if err := models.SetSongbookVerificationStatus(id, true, false, true); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := email.SendSongbookVerifiedEmail(c, id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"songbook": id,
	})
}

func SendToVerifySongbook(c *gin.Context) {
	id := c.Param("id")

	if err := models.SetSongbookVerificationStatus(id, false, true, true); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := email.SendSongbookToVerifiedEmail(c, id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"songbook": id,
	})
}

func RejectSongbook(c *gin.Context) {
	id := c.Param("id")

	if err := models.SetSongbookVerificationStatus(id, false, false, true); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := email.SendSongbookRejectedEmail(c, id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"songbook": id,
	})
}
