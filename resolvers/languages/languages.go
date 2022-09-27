package languages

import (
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/graphql-go/graphql"
)

func GetLanguages(p graphql.ResolveParams) (interface{}, error) {
	reader_code := p.Context.Value("language").(string)
	languages := new(models.Language).GetAllLanguages(reader_code)

	return languages, nil
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
