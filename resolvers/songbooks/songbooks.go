package songbooks

import (
	"fmt"

	"github.com/BenjaminRA/himnario-backend/helpers"
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/graphql-go/graphql"
)

func GetSongbooks(p graphql.ResolveParams) (interface{}, error) {
	lang := p.Context.Value("language").(string)
	songbooks, err := new(models.Songbook).GetAllSongbooks(lang)
	if err != nil {
		return nil, err
	}

	return songbooks, nil
}

func GetSongbook(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["_id"].(string)
	if !ok {
		return nil, nil
	}

	lang := p.Context.Value("language").(string)
	songbook, err := new(models.Songbook).GetSongbookByID(id, lang)
	if err != nil {
		return nil, err
	}

	if songbook.ID.Hex() == "000000000000000000000000" {
		return nil, nil
	}

	return songbook, nil
}

func UpdateSongbook(p graphql.ResolveParams) (interface{}, error) {
	id := p.Args["_id"].(string)
	lang := p.Context.Value("language").(string)

	songbook, err := new(models.Songbook).GetSongbookByID(id, lang)
	if err != nil {
		return nil, err
	}

	if songbook.ID.Hex() == "000000000000000000000000" {
		return nil, fmt.Errorf("songbook not found")
	}

	if err := helpers.BindJSON(p.Args["songbook"], &songbook); err != nil {
		return nil, err
	}

	if err := songbook.UpdateSongbook(); err != nil {
		return nil, err
	}

	return songbook, nil
}

func CreateSongbook(p graphql.ResolveParams) (interface{}, error) {
	var songbook models.Songbook
	lang := p.Context.Value("language").(string)

	if err := helpers.BindJSON(p.Args["songbook"], &songbook); err != nil {
		return nil, err
	}

	err := songbook.CreateSongbook(lang)
	if err != nil {
		return nil, err
	}

	return songbook, nil
}

func DeleteSongbook(p graphql.ResolveParams) (interface{}, error) {
	id := p.Args["_id"].(string)
	lang := p.Context.Value("language").(string)

	songbook, err := new(models.Songbook).GetSongbookByID(id, lang)
	if err != nil {
		return nil, err
	}

	if songbook.ID.Hex() == "000000000000000000000000" {
		return nil, fmt.Errorf("songbook not found")
	}
	if err := songbook.DeleteSongbook(); err != nil {
		return nil, err
	}

	return songbook, nil
}
