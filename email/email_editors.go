package email

import (
	"github.com/BenjaminRA/himnario-backend/models"
)

func EmailEditors(current models.Songbook, updated_raw map[string]interface{}) error {
	return nil
	// var updated models.Songbook

	// if err := helpers.BindJSON(updated_raw, &updated); err != nil {
	// 	return err
	// }

	// added := []string{}
	// deleted := []string{}

	// old_editors := map[string]bool{}
	// for _, editor := range current.Editors {
	// 	old_editors[editor] = true
	// }

	// l9
	// tmp := updated.Editors
	// for _, editor := range updated.Editors {
	// 	_, ok := old_editors[editor]
	// 	if !ok {
	// 		deleted[]
	// 	}
	// }
}
