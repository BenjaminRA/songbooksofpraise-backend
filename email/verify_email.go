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

	// Build email content using template helpers
	content := Paragraph(locale.GetLocalizedMessage(lang, "email.verify.greeting")) +
		Paragraph(locale.GetLocalizedMessage(lang, "email.verify.body")) +
		InfoBox(locale.GetLocalizedMessage(lang, "email.verify.info"), "info")

	// Build footer with the actual link
	footerText := fmt.Sprintf("%s<br><br><a href=\"%s\" style=\"color: #5048E5; word-break: break-all;\">%s</a>", 
		locale.GetLocalizedMessage(lang, "email.verify.footer"),
		link,
		link,
	)

	// Generate HTML email using template
	htmlContent := EmailTemplate(
		locale.GetLocalizedMessage(lang, "email.verify.title"),
		locale.GetLocalizedMessage(lang, "email.verify.preheader"),
		content,
		locale.GetLocalizedMessage(lang, "email.verify.button"),
		link,
		footerText,
	)

	err = SendEmail(
		user.Email,
		locale.GetLocalizedMessage(lang, "email.verify.subject"),
		htmlContent,
	)

	if err != nil {
		return err
	}

	return nil
}
