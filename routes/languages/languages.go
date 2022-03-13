package languages

import (
	"net/http"

	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

func GetLanguages(c *gin.Context) {
	reader_code := c.GetHeader("Language")
	languages := new(models.Language).GetAllLanguages(reader_code)

	c.IndentedJSON(http.StatusOK, languages)
}

func UpdateLanguage(c *gin.Context) {
	code := c.Param("code")

	var language models.Language

	if err := c.BindJSON(&language); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := language.UpdateLanguage(code); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.IndentedJSON(http.StatusCreated, language)
}
