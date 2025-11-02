package auth

import (
	"net/http"

	auth "github.com/BenjaminRA/himnario-backend/auth"
	"github.com/BenjaminRA/himnario-backend/email"
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": locale.GetLocalizedMessage(lang, err.Error())})
		return
	}

	if err := user.Register(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": locale.GetLocalizedMessage(lang, err.Error()),
		})
		return
	}

	token, err := auth.CreateToken(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	email.SendVerificationEmail(c, user)

	token.SendToken(c)

	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func VerifyUserEmail(c *gin.Context) {
	lang := c.Request.Context().Value("language").(string)
	var body gin.H
	c.BindJSON(&body)

	verificationToken := body["token"].(string)

	_, err := auth.VerifyVerificationToken(verificationToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": locale.GetLocalizedMessage(lang, "email.verify.invalid"),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": locale.GetLocalizedMessage(lang, "email.verify.success"),
	})
}

func EmailVerificationResend(c *gin.Context) {
	lang := c.Request.Context().Value("language").(string)
	user, err := auth.RetrieveUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": locale.GetLocalizedMessage(lang, "email.verify.invalid"),
		})
		return
	}

	email.SendVerificationEmail(c, user)

	c.JSON(http.StatusOK, gin.H{})
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

	token, err := auth.CreateToken(user)
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

	user.Password = ""

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
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

func GetUsers(c *gin.Context) {
	users, err := new(models.User).GetAllUsers()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func UpdateUser(c *gin.Context) {
	lang := c.Request.Context().Value("language").(string)
	var user models.User
	var body gin.H
	c.BindJSON(&body)

	if err := helpers.BindJSON(body, &user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": locale.GetLocalizedMessage(lang, err.Error())})
		return
	}

	if err := user.UpdateUser(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": locale.GetLocalizedMessage(lang, err.Error())})
		return
	}

	user.Password = "" // Don't send password back
	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"message": locale.GetLocalizedMessage(lang, "user.update.success"),
	})
}

func DeleteUser(c *gin.Context) {
	lang := c.Request.Context().Value("language").(string)
	var user models.User
	var body gin.H
	c.BindJSON(&body)

	if err := helpers.BindJSON(body, &user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": locale.GetLocalizedMessage(lang, err.Error())})
		return
	}

	if err := user.DeleteUser(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": locale.GetLocalizedMessage(lang, err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": locale.GetLocalizedMessage(lang, "user.delete.success"),
	})
}

// ForgotPassword handles forgot password requests
func ForgotPassword(c *gin.Context) {
	lang := c.Request.Context().Value("language").(string)

	var request map[string]string
	c.BindJSON(&request)

	emailAddress := request["email"]

	// Check if user exists
	user, err := new(models.User).GetUserByEmail(emailAddress)
	if err != nil {
		// Don't reveal if email exists or not for security
		c.JSON(http.StatusOK, gin.H{
			"message": locale.GetLocalizedMessage(lang, "password_reset.email_sent"),
		})
		return
	}

	// Generate password reset token (valid for 24 hours)
	token, err := auth.ResetToken(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Send password reset email
	err = email.SendPasswordResetEmail(c, user, token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": locale.GetLocalizedMessage(lang, "password_reset.email_sent"),
	})
}

// ResetPassword handles password reset with token
func ResetPassword(c *gin.Context) {
	lang := c.Request.Context().Value("language").(string)

	var request map[string]string
	c.BindJSON(&request)

	token := request["token"]
	newPassword := request["password"]

	// Validate the token
	user, err := auth.VerifyResetToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": locale.GetLocalizedMessage(lang, "password_reset.token_invalid"),
		})
		return
	}

	// Update the password
	user.Password = newPassword
	err = user.UpdatePassword()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": locale.GetLocalizedMessage(lang, "password_reset.success"),
	})
}

// ChangePassword handles password change for authenticated users
func ChangePassword(c *gin.Context) {
	lang := c.Request.Context().Value("language").(string)

	user, err := auth.RetrieveUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": locale.GetLocalizedMessage(lang, "session_token.invalid"),
		})
		return
	}

	var request map[string]string
	c.BindJSON(&request)

	currentPassword := request["current_password"]
	newPassword := request["new_password"]

	// Verify current password
	_, err = user.Login(user.Email, currentPassword)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": locale.GetLocalizedMessage(lang, "change_password.current_password_incorrect"),
		})
		return
	}

	// Update to new password
	user.Password = newPassword
	err = user.UpdatePassword()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": locale.GetLocalizedMessage(lang, "change_password.success"),
	})
}
