package email

import (
	"fmt"
	"os"

	"github.com/BenjaminRA/himnario-backend/locale"
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

func SendPasswordResetEmail(c *gin.Context, user models.User, token string) error {
	lang := c.Request.Context().Value("language").(string)

	link := fmt.Sprintf("%s/reset-password?token=%s", os.Getenv("FRONTEND_URL"), token)

	// Build email content using template helpers
	content := Paragraph(locale.GetLocalizedMessage(lang, "email.password_reset.greeting")) +
		Paragraph(locale.GetLocalizedMessage(lang, "email.password_reset.body")) +
		InfoBox(locale.GetLocalizedMessage(lang, "email.password_reset.info"), "warning") +
		Paragraph(locale.GetLocalizedMessage(lang, "email.password_reset.ignore"))

	// Build footer with the actual link
	footerText := fmt.Sprintf("%s<br><br><a href=\"%s\" style=\"color: #5048E5; word-break: break-all;\">%s</a>",
		locale.GetLocalizedMessage(lang, "email.password_reset.footer"),
		link,
		link,
	)

	// Generate HTML email using template
	htmlContent := EmailTemplate(
		locale.GetLocalizedMessage(lang, "email.password_reset.title"),
		locale.GetLocalizedMessage(lang, "email.password_reset.preheader"),
		content,
		locale.GetLocalizedMessage(lang, "email.password_reset.button"),
		link,
		footerText,
	)

	err := SendEmail(
		user.Email,
		locale.GetLocalizedMessage(lang, "email.password_reset.subject"),
		htmlContent,
	)

	if err != nil {
		return err
	}

	return nil
}
