package email

import (
	"fmt"

	"github.com/BenjaminRA/himnario-backend/locale"
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

func SendSongbookToVerifiedEmail(c *gin.Context, songbook_id string) error {
	lang := c.Request.Context().Value("language").(string)

	songbook, err := new(models.Songbook).GetSongbookByID(songbook_id, "EN")
	if err != nil {
		return err
	}

	users, err := new(models.User).GetAllModerators()
	if err != nil {
		return err
	}

	for _, user := range users {
		err = SendEmail(
			user.Email,
			locale.GetLocalizedMessage(lang, "email.songbook.sent_to_verify.subject"),
			fmt.Sprintf(locale.GetLocalizedMessage(lang, "email.songbook.sent_to_verify.content"), songbook.Title),
		)

		if err != nil {
			return err
		}
	}

	return nil
}
