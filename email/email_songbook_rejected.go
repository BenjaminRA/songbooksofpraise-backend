package email

import (
	"fmt"

	"github.com/BenjaminRA/himnario-backend/locale"
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/gin-gonic/gin"
)

func SendSongbookRejectedEmail(c *gin.Context, songbook_id int) error {
	lang := c.Request.Context().Value("language").(string)

	songbook, err := new(models.Songbook).GetSongbookByID(songbook_id)
	if err != nil {
		return err
	}

	user, err := new(models.User).GetUserById(songbook.OwnerID)
	if err != nil {
		return err
	}

	err = SendEmail(
		user.Email,
		locale.GetLocalizedMessage(lang, "email.songbook.rejected.subject"),
		fmt.Sprintf(locale.GetLocalizedMessage(lang, "email.songbook.rejected.content"), songbook.Title),
	)
	if err != nil {
		return err
	}

	// Get songbook editors
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
			locale.GetLocalizedMessage(lang, "email.songbook.rejected.subject"),
			fmt.Sprintf(locale.GetLocalizedMessage(lang, "email.songbook.rejected.content"), songbook.Title),
		)

		if err != nil {
			return err
		}
	}

	return nil
}
