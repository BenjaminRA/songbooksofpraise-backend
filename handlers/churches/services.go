package churches

import (
	"net/http"
	"strconv"

	models "github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

// GetServices returns services for a specific church
func GetServices(c *gin.Context) {
	churchID := c.Param("church_id")
	churchIDInt, err := strconv.Atoi(churchID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid church ID"})
		return
	}

	service := &models.Service{}
	services, err := service.GetServicesByChurchID(churchIDInt)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"services":  services,
		"count":     len(services),
		"church_id": churchIDInt,
	})
}

// GetServiceByID returns a specific service by ID
func GetServiceByID(c *gin.Context) {
	id := c.Param("id")
	serviceID, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	service := &models.Service{}
	result, err := service.GetServiceByID(serviceID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"service": result,
	})
}

// CreateService creates a new service
func CreateService(c *gin.Context) {
	var req models.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service := &models.Service{}
	result, err := service.CreateService(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"service": result,
	})
}

// UpdateService updates an existing service
func UpdateService(c *gin.Context) {
	id := c.Param("id")
	serviceID, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	var req models.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service := &models.Service{}
	result, err := service.UpdateService(serviceID, &req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"service": result,
	})
}

// DeleteService deletes a service
func DeleteService(c *gin.Context) {
	id := c.Param("id")
	serviceID, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	service := &models.Service{}
	err = service.DeleteService(serviceID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Service deleted successfully",
	})
}
