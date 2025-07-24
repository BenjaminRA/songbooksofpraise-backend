package email

import (
	"fmt"
	"os"

	"github.com/BenjaminRA/himnario-backend/auth"
	"github.com/BenjaminRA/himnario-backend/locale"
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

func SendVerificationEmail(c *gin.Context, user models.User) error {
	lang := c.Request.Context().Value("language").(string)

	token, err := auth.VerificationToken(user)
	if err != nil {
		return err
	}

	link := fmt.Sprintf("%s/login?token=%s", os.Getenv("FRONTEND_URL"), token)

	err = SendEmail(
		user.Email,
		locale.GetLocalizedMessage(lang, "email.verify.subject"),
		fmt.Sprintf(locale.GetLocalizedMessage(lang, "email.verify.message"), link),
	)

	if err != nil {
		return err
	}

	return nil
}
