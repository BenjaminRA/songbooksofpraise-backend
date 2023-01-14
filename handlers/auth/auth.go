package auth

import (
	"net/http"

	auth "github.com/BenjaminRA/himnario-backend/auth"
	"github.com/BenjaminRA/himnario-backend/helpers"
	"github.com/BenjaminRA/himnario-backend/locale"
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	lang := c.Request.Context().Value("language").(string)
	var user models.User

	var body gin.H
	c.BindJSON(&body)

	if err := helpers.BindJSON(body, &user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	if err := user.Register(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": locale.GetLocalizedMessage(lang, err.Error()),
		})
		return
	}

	token, err := auth.CreateToken(user.ID.Hex())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	token.SendToken(c)

	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func Login(c *gin.Context) {
	lang := c.Request.Context().Value("language").(string)

	var request map[string]string
	c.BindJSON(&request)

	email := request["email"]
	password := request["password"]

	user, err := new(models.User).Login(email, password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": locale.GetLocalizedMessage(lang, "login.invalid.credentials"),
		})
		return
	}

	token, err := auth.CreateToken(user.ID.Hex())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	token.SendToken(c)
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func GetUser(c *gin.Context) {
	user, err := auth.RetrieveUser(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func Logout(c *gin.Context) {
	auth.UnsetToken(c)
	c.Status(http.StatusOK)
}
