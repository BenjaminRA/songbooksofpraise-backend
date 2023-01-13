package handlers

import (
	"fmt"
	"net/http"

	"github.com/BenjaminRA/himnario-backend/helpers"
	"github.com/BenjaminRA/himnario-backend/locale"
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	lang := c.Request.Context().Value("language").(string)
	var user models.User

	var body map[string]interface{}
	c.BindJSON(&body)

	if err := helpers.BindJSON(body, &user); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if err := user.Register(); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": locale.GetLocalizedMessage(lang, err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func Login(c *gin.Context) {
	lang := c.Request.Context().Value("language").(string)

	fmt.Println(c.Request.Cookie("SessionToken"))

	var request map[string]string
	c.BindJSON(&request)

	email := request["email"]
	password := request["password"]

	user, err := new(models.User).Login(email, password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": locale.GetLocalizedMessage(lang, "login.invalid.credentials"),
		})
		return
	}

	c.Header("Set-Cookie", "SessionToken=asdasasda; SameSite=Strict; Secure; HttpOnly")
	c.JSON(http.StatusOK, user)
}
