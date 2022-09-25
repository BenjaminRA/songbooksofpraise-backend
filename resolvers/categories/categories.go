package categories

import (
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/graphql-go/graphql"
)

func GetCategories(p graphql.ResolveParams) (interface{}, error) {
	categories := new(models.Category).GetAllCategories()

	return categories, nil
}

func GetCategory(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["_id"].(string)
	if !ok {
		return nil, nil
	}

	category := new(models.Category).GetCategoryById(id)
	if category.ID.Hex() == "000000000000000000000000" {
		return nil, nil
	}

	return category, nil
}

// func GetCategoriesById(c *gin.Context) {
// 	id := c.Param("id")

// 	category := new(models.Category).GetCategoryById(id)

// 	if category.ID.Hex() == "000000000000000000000000" {
// 		c.IndentedJSON(http.StatusNotFound, category)
// 	} else {
// 		c.IndentedJSON(http.StatusOK, category)
// 	}

// }

// func DeleteCategory(c *gin.Context) {
// 	id := c.Param("id")

// 	category := new(models.Category).GetCategoryById(id)

// 	if category.ID.Hex() == "000000000000000000000000" {
// 		c.IndentedJSON(http.StatusNotFound, category)
// 	} else {
// 		if err := category.DeleteCategory(); err != nil {
// 			c.IndentedJSON(http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		c.IndentedJSON(http.StatusOK, category)
// 	}
// }

// func UpdateCategory(c *gin.Context) {
// 	id := c.Param("id")

// 	category := new(models.Category).GetCategoryById(id)

// 	if category.ID.Hex() == "000000000000000000000000" {
// 		c.IndentedJSON(http.StatusNotFound, gin.H{
// 			"message": "Songbook not found",
// 		})
// 		return
// 	}

// 	if err := c.BindJSON(&category); err != nil {
// 		c.IndentedJSON(http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	if err := category.UpdateCategory(); err != nil {
// 		c.IndentedJSON(http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	c.IndentedJSON(http.StatusOK, category)
// }

// func PostCategory(c *gin.Context) {
// 	var category models.Category

// 	if err := c.BindJSON(&category); err != nil {
// 		c.IndentedJSON(http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	category, err := category.CreateCategory()
// 	if err != nil {
// 		c.IndentedJSON(http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	c.IndentedJSON(http.StatusCreated, category)
// }
