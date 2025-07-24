package categories

import (
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/graphql-go/graphql"
)

func GetAuthors(p graphql.ResolveParams) (interface{}, error) {
	authors, err := new(models.Author).GetAllAuthors()
	if err != nil {
		return nil, err
	}

	return authors, nil
}
