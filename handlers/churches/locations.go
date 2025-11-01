package churches

import (
	"net/http"
	"strconv"

	models "github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

// GetCountries returns all countries
func GetCountries(c *gin.Context) {
	country := &models.Country{}
	countries, err := country.GetAllCountries()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"countries": countries,
		"count":     len(countries),
	})
}

// GetCountriesWithStates returns all countries with their states
func GetCountriesWithStates(c *gin.Context) {
	country := &models.Country{}
	countries, err := country.GetCountriesWithStates()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"countries": countries,
		"count":     len(countries),
	})
}

// GetCountryByID returns a specific country by ID
func GetCountryByID(c *gin.Context) {
	id := c.Param("id")
	countryID, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid country ID"})
		return
	}

	country := &models.Country{}
	result, err := country.GetCountryByID(countryID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Country not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"country": result,
	})
}

// GetStates returns all states
func GetStates(c *gin.Context) {
	state := &models.State{}
	states, err := state.GetAllStates()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"states": states,
		"count":  len(states),
	})
}

// GetStateByID returns a specific state by ID
func GetStateByID(c *gin.Context) {
	id := c.Param("id")
	stateID, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid state ID"})
		return
	}

	state := &models.State{}
	result, err := state.GetStateByID(stateID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "State not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"state": result,
	})
}

// GetStatesByCountry returns states for a specific country
func GetStatesByCountry(c *gin.Context) {
	countryID := c.Param("country_id")
	countryIDInt, err := strconv.Atoi(countryID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid country ID"})
		return
	}

	state := &models.State{}
	states, err := state.GetStatesByCountryID(countryIDInt)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"states":     states,
		"count":      len(states),
		"country_id": countryIDInt,
	})
}
