package email

import (
	"fmt"
	"os"

	"github.com/BenjaminRA/himnario-backend/locale"
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

func SendSongbookToVerifiedEmail(c *gin.Context, songbook_id int) error {
	lang := c.Request.Context().Value("language").(string)

	songbook, err := new(models.Songbook).GetSongbookByID(songbook_id)
	if err != nil {
		return err
	}

	users, err := new(models.User).GetAllModerators()
	if err != nil {
		return err
	}

	songbookLink := fmt.Sprintf("%s/songbooks/%d", os.Getenv("FRONTEND_URL"), songbook.ID)

	for _, user := range users {
		content := Paragraph(locale.GetLocalizedMessage(lang, "email.songbook.sent_to_verify.greeting")) +
			Paragraph(fmt.Sprintf(locale.GetLocalizedMessage(lang, "email.songbook.sent_to_verify.body"), songbook.Title)) +
			InfoBox(locale.GetLocalizedMessage(lang, "email.songbook.sent_to_verify.info"), "info")

		// Build footer with the actual link
		footerText := fmt.Sprintf("%s<br><br><a href=\"%s\" style=\"color: #5048E5; word-break: break-all;\">%s</a>", 
			locale.GetLocalizedMessage(lang, "email.songbook.sent_to_verify.footer"),
			songbookLink,
			songbookLink,
		)

		htmlContent := EmailTemplate(
			locale.GetLocalizedMessage(lang, "email.songbook.sent_to_verify.title"),
			locale.GetLocalizedMessage(lang, "email.songbook.sent_to_verify.preheader"),
			content,
			locale.GetLocalizedMessage(lang, "email.songbook.sent_to_verify.button"),
			songbookLink,
			footerText,
		)

		err = SendEmail(
			user.Email,
			locale.GetLocalizedMessage(lang, "email.songbook.sent_to_verify.subject"),
			htmlContent,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
