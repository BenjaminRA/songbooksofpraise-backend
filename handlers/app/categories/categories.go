package app_categories

import (
	"net/http"
	"strconv"

	"github.com/BenjaminRA/himnario-backend/locale"
	models "github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	songbookIDStr := c.Param("id")
	if songbookIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Songbook ID is required"})
		return
	}

	songbookID, err := strconv.Atoi(songbookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid songbook ID format"})
		return
	}

	categories, err := (&models.Category{}).GetCategoriesBySongbookID(songbookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}

func GetCategoryByID(c *gin.Context) {
	lang := c.Request.Context().Value("language").(string)

	songbookIDStr := c.Param("id")
	categoryIDStr := c.Param("category_id")
	if categoryIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category ID is required"})
		return
	}

	if categoryIDStr == "all" {
		if songbookIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Songbook ID is required"})
			return
		}

		songbookID, err := strconv.Atoi(songbookIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid songbook ID format"})
			return
		}

		category := &models.Category{
			ID:         -1, // Use -1 to indicate all categories
			Name:       locale.GetLocalizedMessage(lang, "all"),
			SongbookID: &songbookID,
		}

		category.Songs, err = models.GetSongsBySongbookID(songbookID) // Fetch all songs for all categories
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve songs for all categories"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"category": category,
		})
		return
	}

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID format"})
		return
	}

	category, err := (&models.Category{}).GetCategoryById(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve category"})
		return
	}

	category.Songs, err = models.GetSongsByCategoryID(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve songs for category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"category": category,
	})
}
