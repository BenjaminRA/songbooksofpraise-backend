package email

import (
	"fmt"
	"os"

	"github.com/BenjaminRA/himnario-backend/locale"
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

func SendSongbookVerifiedEmail(c *gin.Context, songbook_id int) error {
	lang := c.Request.Context().Value("language").(string)

	songbook, err := new(models.Songbook).GetSongbookByID(songbook_id)
	if err != nil {
		return err
	}

	user, err := new(models.User).GetUserById(songbook.OwnerID)
	if err != nil {
		return err
	}

	songbookLink := fmt.Sprintf("%s/songbooks/%d", os.Getenv("FRONTEND_URL"), songbook.ID)

	// Send to owner
	content := Paragraph(locale.GetLocalizedMessage(lang, "email.songbook.verified.greeting")) +
		Paragraph(fmt.Sprintf(locale.GetLocalizedMessage(lang, "email.songbook.verified.body"), songbook.Title)) +
		InfoBox(locale.GetLocalizedMessage(lang, "email.songbook.verified.info"), "success")

	// Build footer with the actual link
	footerText := fmt.Sprintf("%s<br><br><a href=\"%s\" style=\"color: #5048E5; word-break: break-all;\">%s</a>", 
		locale.GetLocalizedMessage(lang, "email.songbook.verified.footer"),
		songbookLink,
		songbookLink,
	)

	htmlContent := EmailTemplate(
		locale.GetLocalizedMessage(lang, "email.songbook.verified.title"),
		locale.GetLocalizedMessage(lang, "email.songbook.verified.preheader"),
		content,
		locale.GetLocalizedMessage(lang, "email.songbook.verified.button"),
		songbookLink,
		footerText,
	)

	err = SendEmail(
		user.Email,
		locale.GetLocalizedMessage(lang, "email.songbook.verified.subject"),
		htmlContent,
	)
	if err != nil {
		return err
	}

	// Get songbook editors and send to them too
	editors, err := models.GetSongbookEditors(songbook.ID)
	if err != nil {
		return err
	}

	for _, editor := range editors {
		editorUser, err := new(models.User).GetUserById(editor.UserID)
		if err != nil {
			continue
		}

		err = SendEmail(
			editorUser.Email,
			locale.GetLocalizedMessage(lang, "email.songbook.verified.subject"),
			htmlContent,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
