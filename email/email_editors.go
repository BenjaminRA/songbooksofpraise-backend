package email

import (
	"fmt"

	"github.com/BenjaminRA/himnario-backend/locale"
	"github.com/BenjaminRA/himnario-backend/models"
)

func EmailEditors(current []string, updated models.Songbook, lang string) error {
	fmt.Printf("current: %v\n", current)
	fmt.Printf("new_editors: %v\n", updated.Editors)
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
	for _, editor := range updated.Editors {
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

	fmt.Printf("deleted: %v\n", deleted)
	for _, email := range deleted {
		err := SendEmail(
			email,
			locale.GetLocalizedMessage(lang, "email.editor.removed.subject"),
			fmt.Sprintf(locale.GetLocalizedMessage(lang, "email.editor.removed.content"), updated.Title),
		)

		if err != nil {
			return err
		}
	}

	fmt.Printf("added: %v\n", added)
	for _, email := range added {
		err := SendEmail(
			email,
			locale.GetLocalizedMessage(lang, "email.editor.added.subject"),
			fmt.Sprintf(locale.GetLocalizedMessage(lang, "email.editor.added.content"), updated.Title),
		)

		if err != nil {
			return err
		}
	}

	return nil
}
