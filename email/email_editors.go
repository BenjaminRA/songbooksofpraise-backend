package email

import (
	"fmt"
	"os"

	"github.com/BenjaminRA/himnario-backend/locale"
	"github.com/BenjaminRA/himnario-backend/models"
)

func EmailEditors(current []string, updated models.Songbook, lang string) error {
	// Get current editors from database
	editors, err := models.GetSongbookEditors(updated.ID)
	if err != nil {
		return err
	}

	// Convert to email addresses
	newEditors := make([]string, len(editors))
	for i, editor := range editors {
		user, err := new(models.User).GetUserById(editor.UserID)
		if err != nil {
			continue
		}
		newEditors[i] = user.Email
	}

	fmt.Printf("current: %v\n", current)
	fmt.Printf("new_editors: %v\n", newEditors)
	added := []string{}
	deleted := []string{}
	all_editors := []string{}

	old_editors := map[string]bool{}
	for _, editor := range current {
		if _, ok := old_editors[editor]; !ok {
			all_editors = append(all_editors, editor)
		}
		old_editors[editor] = true
	}

	new_editors := map[string]bool{}
	for _, editor := range newEditors {
		if _, ok := new_editors[editor]; !ok {
			if _, ok := old_editors[editor]; !ok {
				all_editors = append(all_editors, editor)
			}
		}
		new_editors[editor] = true
	}

	fmt.Printf("all_editors: %v\n", all_editors)

	for _, editor := range all_editors {
		_, old_ok := old_editors[editor]
		_, new_ok := new_editors[editor]
		fmt.Printf("editor: %s - old_ok: %v - new_ok: %v\n", editor, old_ok, new_ok)
		// if it exists in the new list and doesn't exists in the old list, means it has been added
		if new_ok && !old_ok {
			added = append(added, editor)
		}

		// if it doesn't exists in the new list and exists in the old list, means it has been deleted
		if !new_ok && old_ok {
			deleted = append(deleted, editor)
		}
	}

	// Send emails to removed editors
	fmt.Printf("deleted: %v\n", deleted)
	for _, email := range deleted {
		content := Paragraph(locale.GetLocalizedMessage(lang, "email.editor.removed.greeting")) +
			Paragraph(fmt.Sprintf(locale.GetLocalizedMessage(lang, "email.editor.removed.body"), updated.Title)) +
			InfoBox(locale.GetLocalizedMessage(lang, "email.editor.removed.info"), "info") +
			Paragraph(locale.GetLocalizedMessage(lang, "email.editor.removed.footer"))

		htmlContent := EmailTemplate(
			locale.GetLocalizedMessage(lang, "email.editor.removed.title"),
			locale.GetLocalizedMessage(lang, "email.editor.removed.preheader"),
			content,
			"",
			"",
			"",
		)

		err := SendEmail(
			email,
			locale.GetLocalizedMessage(lang, "email.editor.removed.subject"),
			htmlContent,
		)

		if err != nil {
			return err
		}
	}

	// Send emails to added editors
	fmt.Printf("added: %v\n", added)
	for _, email := range added {
		content := Paragraph(locale.GetLocalizedMessage(lang, "email.editor.added.greeting")) +
			Paragraph(fmt.Sprintf(locale.GetLocalizedMessage(lang, "email.editor.added.body"), updated.Title)) +
			InfoBox(locale.GetLocalizedMessage(lang, "email.editor.added.info"), "success")

		songbookLink := fmt.Sprintf("%s/songbooks/%d", os.Getenv("FRONTEND_URL"), updated.ID)

		// Build footer with the actual link
		footerText := fmt.Sprintf("%s<br><br><a href=\"%s\" style=\"color: #5048E5; word-break: break-all;\">%s</a>", 
			locale.GetLocalizedMessage(lang, "email.editor.added.footer"),
			songbookLink,
			songbookLink,
		)

		htmlContent := EmailTemplate(
			locale.GetLocalizedMessage(lang, "email.editor.added.title"),
			locale.GetLocalizedMessage(lang, "email.editor.added.preheader"),
			content,
			locale.GetLocalizedMessage(lang, "email.editor.added.button"),
			songbookLink,
			footerText,
		)

		err := SendEmail(
			email,
			locale.GetLocalizedMessage(lang, "email.editor.added.subject"),
			htmlContent,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
