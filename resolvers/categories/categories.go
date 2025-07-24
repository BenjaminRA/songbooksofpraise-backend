package categories

import (
	"fmt"

	"github.com/BenjaminRA/himnario-backend/helpers"
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/graphql-go/graphql"
)

func GetCategories(p graphql.ResolveParams) (interface{}, error) {
	categories, err := new(models.Category).GetAllCategories()
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func GetCategory(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["_id"].(string)
	if !ok {
		return nil, nil
	}

	category, err := new(models.Category).GetCategoryById(id)
	if err != nil {
		return nil, err
	}

	if category.ID.Hex() == "000000000000000000000000" {
		return nil, nil
	}

	return category, nil
}

func CreateCategory(p graphql.ResolveParams) (interface{}, error) {
	var category models.Category

	if err := helpers.BindJSON(p.Args["category"], &category); err != nil {
		return nil, err
	}

	err := category.CreateCategory()

	if err != nil {
		return nil, err
	}

	return category, nil
}

func UpdateCategory(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["_id"].(string)
	if !ok {
		return nil, nil
	}

	category, err := new(models.Category).GetCategoryById(id)
	if err != nil {
		return nil, err
	}

	if category.ID.Hex() == "000000000000000000000000" {
		return nil, fmt.Errorf("category not found")
	}

	if err := helpers.BindJSON(p.Args["category"], &category); err != nil {
		return nil, err
	}

	if err := category.UpdateCategory(); err != nil {
		return nil, err
	}

	return category, nil
}

func DeleteCategory(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["_id"].(string)
	if !ok {
		return nil, nil
	}

	category, err := new(models.Category).GetCategoryById(id)
	if err != nil {
		return nil, err
	}

	if category.ID.Hex() == "000000000000000000000000" {
		return nil, fmt.Errorf("category not found")
	}

	if err := category.DeleteCategory(); err != nil {
		return nil, err
	}

	return category, nil

}
