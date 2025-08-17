package songbooks

import (
	"net/http"
	"strconv"

	"github.com/BenjaminRA/himnario-backend/auth"
	"github.com/BenjaminRA/himnario-backend/email"
	models "github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

func VerifySongbook(c *gin.Context) {
	id := c.Param("id")

	songbookID, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid songbook ID"})
		return
	}

	if err := models.SetSongbookVerificationStatus(songbookID, true, false, false, true); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := email.SendSongbookVerifiedEmail(c, songbookID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"songbook": id,
	})
}

func SendToVerifySongbook(c *gin.Context) {
	id := c.Param("id")

	songbookID, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid songbook ID"})
		return
	}

	if err := models.SetSongbookVerificationStatus(songbookID, false, true, false, true); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := email.SendSongbookToVerifiedEmail(c, songbookID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"songbook": id,
	})
}

func RejectSongbook(c *gin.Context) {
	id := c.Param("id")

	songbookID, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid songbook ID"})
		return
	}

	if err := models.SetSongbookVerificationStatus(songbookID, false, false, true, true); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := email.SendSongbookRejectedEmail(c, songbookID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"songbook": id,
	})
}

func GetSongbooks(c *gin.Context) {
	songbooks, err := (&models.Songbook{}).GetAllSongbooks()
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

func CreateSongbook(c *gin.Context) {
	// Extract the raw JSON data to get editors
	var requestData map[string]interface{}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		panic("Invalid request body")
	}

	// Create songbook struct from request data
	var songbook models.Songbook
	if title, ok := requestData["title"].(string); ok {
		songbook.Title = title
	}

	// get user from context
	user, err := auth.RetrieveUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	songbook.OwnerID = user.ID

	// Create the songbook first to get the ID
	if err := songbook.CreateSongbook(); err != nil {
		panic(err)
	}

	// Extract and add editors if present
	if editorsData, ok := requestData["editors"].([]interface{}); ok {
		for _, editorInterface := range editorsData {
			if editor, ok := editorInterface.(string); ok {
				songbook.AddEditor(editor)
			}
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"songbook": songbook,
	})
}

func DeleteSongbook(c *gin.Context) {
	id := c.Param("id")

	songbookID, err := strconv.Atoi(id)
	if err != nil {
		panic("Invalid songbook ID")
	}

	songbook, err := (&models.Songbook{}).GetSongbookByID(songbookID)
	if err != nil {
		panic(err)
	}

	if err := songbook.DeleteSongbook(); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Songbook deleted successfully",
	})
}

func UpdateSongbook(c *gin.Context) {
	id := c.Param("id")

	print(id)

	songbookID, err := strconv.Atoi(id)
	if err != nil {
		panic("Invalid songbook ID")
	}

	// Extract the raw JSON data for later use
	var requestData map[string]interface{}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		panic("Invalid request body")
	}

	// Convert the map data to songbook struct
	songbook, err := (&models.Songbook{}).GetSongbookByID(songbookID)
	if err != nil {
		panic(err)
	}

	// Update fields if present in request data
	if title, ok := requestData["title"].(string); ok {
		songbook.Title = title
	}
	if verified, ok := requestData["verified"].(bool); ok {
		songbook.Verified = verified
	}
	if inVerification, ok := requestData["in_verification"].(bool); ok {
		songbook.InVerification = inVerification
	}
	if rejected, ok := requestData["rejected"].(bool); ok {
		songbook.Rejected = rejected
	}
	if ownerID, ok := requestData["owner_id"].(float64); ok {
		songbook.OwnerID = int(ownerID)
	}

	songbook.ID = songbookID
	if err := songbook.UpdateSongbook(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	songbook.RemoveAllEditors()

	// Update editors if provided
	if editorsData, ok := requestData["editors"].([]interface{}); ok {
		for _, editorInterface := range editorsData {
			if editor, ok := editorInterface.(string); ok {
				songbook.AddEditor(editor)
			}
		}
	}

	// You can now use requestData for any additional processing
	// For example, logging the original request:
	// log.Printf("Original request data: %+v", requestData)

	c.JSON(http.StatusOK, gin.H{
		"songbook": songbook,
	})
}
