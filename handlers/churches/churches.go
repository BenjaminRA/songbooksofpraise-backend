package churches

import (
	"net/http"
	"strconv"

	models "github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

// GetChurches returns all churches
func GetChurches(c *gin.Context) {
	church := &models.Church{}
	churches, err := church.GetAllChurches()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"churches": churches,
		"count":    len(churches),
	})
}

// GetChurchByID returns a specific church by ID
func GetChurchByID(c *gin.Context) {
	id := c.Param("id")
	churchID, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid church ID"})
		return
	}

	church := &models.Church{}
	result, err := church.GetChurchByID(churchID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Church not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"church": result,
	})
}

// CreateChurch creates a new church
func CreateChurch(c *gin.Context) {
	var newChurch models.Church
	if err := c.ShouldBindJSON(&newChurch); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := newChurch.CreateChurch(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create church: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Church created successfully",
		"church":  newChurch,
	})
}

// UpdateChurch updates an existing church
func UpdateChurch(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Church ID is required"})
		return
	}

	churchID, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid church ID"})
		return
	}

	var updatedChurch models.Church
	if err := c.ShouldBindJSON(&updatedChurch); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedChurch.ID = churchID

	if err := updatedChurch.UpdateChurch(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update church"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Church updated successfully",
		"church":  updatedChurch,
	})
}

// DeleteChurch deletes a church
func DeleteChurch(c *gin.Context) {
	id := c.Param("id")
	churchID, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid church ID"})
		return
	}

	church := &models.Church{}
	err = church.DeleteChurch(churchID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Church deleted successfully",
	})
}

// GetChurchesByState returns churches in a specific state
func GetChurchesByState(c *gin.Context) {
	stateID := c.Param("state_id")
	stateIDInt, err := strconv.Atoi(stateID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid state ID"})
		return
	}

	church := &models.Church{}
	churches, err := church.GetChurchesByStateID(stateIDInt)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"churches": churches,
		"count":    len(churches),
		"state_id": stateIDInt,
	})
}
