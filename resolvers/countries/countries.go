package countries

import (
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/graphql-go/graphql"
)

func GetCountries(p graphql.ResolveParams) (interface{}, error) {
	reader_code := p.Context.Value("language").(string)
	countries := new(models.Country).GetAllCountries(reader_code)

	return countries, nil
}

// func UpdateLanguage(c *gin.Context) {
// 	code := c.Param("code")

// 	var language models.Language

// 	if err := c.BindJSON(&language); err != nil {
// 		c.IndentedJSON(http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	if err := language.UpdateLanguage(code); err != nil {
// 		c.IndentedJSON(http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	c.IndentedJSON(http.StatusCreated, language)
// }
