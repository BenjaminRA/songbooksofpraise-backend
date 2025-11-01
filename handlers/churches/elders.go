package churches

import (
	"net/http"
	"strconv"

	models "github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

// GetElders returns elders for a specific church
func GetElders(c *gin.Context) {
	churchID := c.Param("church_id")
	churchIDInt, err := strconv.Atoi(churchID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid church ID"})
		return
	}

	elder := &models.Elder{}
	elders, err := elder.GetEldersByChurchID(churchIDInt)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"elders":    elders,
		"count":     len(elders),
		"church_id": churchIDInt,
	})
}

// GetElderByID returns a specific elder by ID
func GetElderByID(c *gin.Context) {
	id := c.Param("id")
	elderID, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid elder ID"})
		return
	}

	elder := &models.Elder{}
	result, err := elder.GetElderByID(elderID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Elder not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"elder": result,
	})
}

// CreateElder creates a new elder
func CreateElder(c *gin.Context) {
	var req models.CreateElderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	elder := &models.Elder{}
	result, err := elder.CreateElder(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"elder": result,
	})
}

// UpdateElder updates an existing elder
func UpdateElder(c *gin.Context) {
	id := c.Param("id")
	elderID, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid elder ID"})
		return
	}

	var req models.UpdateElderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	elder := &models.Elder{}
	result, err := elder.UpdateElder(elderID, &req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Elder not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"elder": result,
	})
}

// DeleteElder deletes an elder
func DeleteElder(c *gin.Context) {
	id := c.Param("id")
	elderID, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid elder ID"})
		return
	}

	elder := &models.Elder{}
	err = elder.DeleteElder(elderID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Elder deleted successfully",
	})
}
