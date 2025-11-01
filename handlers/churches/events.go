package churches

import (
	"net/http"
	"strconv"

	models "github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

// GetEvents returns events for a specific church
func GetEvents(c *gin.Context) {
	churchID := c.Param("id")
	churchIDInt, err := strconv.Atoi(churchID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid church ID"})
		return
	}

	event := &models.Event{}
	events, err := event.GetEventsByChurchID(churchIDInt)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"events":    events,
		"count":     len(events),
		"church_id": churchIDInt,
	})
}

// GetEventByID returns a specific event by ID
func GetEventByID(c *gin.Context) {
	id := c.Param("id")
	eventID, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event := &models.Event{}
	result, err := event.GetEventByID(eventID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"event": result,
	})
}

// CreateEvent creates a new event
func CreateEvent(c *gin.Context) {
	var req models.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := &models.Event{}
	result, err := event.CreateEvent(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"event": result,
	})
}

// UpdateEvent updates an existing event
func UpdateEvent(c *gin.Context) {
	id := c.Param("id")
	eventID, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var req models.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := &models.Event{}
	result, err := event.UpdateEvent(eventID, &req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"event": result,
	})
}

// DeleteEvent deletes an event
func DeleteEvent(c *gin.Context) {
	id := c.Param("id")
	eventID, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event := &models.Event{}
	err = event.DeleteEvent(eventID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event deleted successfully",
	})
}
