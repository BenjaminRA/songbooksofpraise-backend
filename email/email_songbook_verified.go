package email

import (
	"fmt"

	"github.com/BenjaminRA/himnario-backend/locale"
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

func SendSongbookVerifiedEmail(c *gin.Context, songbook_id string) error {
	lang := c.Request.Context().Value("language").(string)

	songbook, err := new(models.Songbook).GetSongbookByID(songbook_id, "EN")
	if err != nil {
		return err
	}

	user, err := new(models.User).GetUserById(songbook.OwnerID.Hex())
	if err != nil {
		return err
	}

	err = SendEmail(
		user.Email,
		locale.GetLocalizedMessage(lang, "email.songbook.verified.subject"),
		fmt.Sprintf(locale.GetLocalizedMessage(lang, "email.songbook.verified.content"), songbook.Title),
	)
	if err != nil {
		return err
	}

	for _, email := range songbook.Editors {
		err = SendEmail(
			email,
			locale.GetLocalizedMessage(lang, "email.songbook.verified.subject"),
			fmt.Sprintf(locale.GetLocalizedMessage(lang, "email.songbook.verified.content"), songbook.Title),
		)

		if err != nil {
			return err
		}
	}

	return nil
}
