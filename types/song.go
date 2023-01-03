package types

import (
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/graphql-go/graphql"
)

var Song = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Song",
		Fields: graphql.Fields{
			"_id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Source.(models.Song).ID.Hex()
					return id, nil
				},
			},
			"songbook_id": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Source.(models.Song).SongbookID.Hex()
					return id, nil
				},
			},
			"categories_id": &graphql.Field{
				Type: graphql.NewList(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					temp := []string{}
					ids := p.Source.(models.Song).CategoriesID

					for _, id := range ids {
						temp = append(temp, id.Hex())
					}

					return temp, nil
				},
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"chords": &graphql.Field{
				Type: graphql.Boolean,
			},
			"music_sheet": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Source.(models.Song).MusicSheet.Hex()
					return id, nil
				},
			},
			"music": &graphql.Field{
				Type: graphql.String,
			},
			"author": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"bible_verse": &graphql.Field{
				Type: graphql.String,
			},
			"number": &graphql.Field{
				Type: graphql.Int,
			},
			"text": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)

var NewSong = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "NewSong",
	Fields: graphql.InputObjectConfigFieldMap{
		"songbook_id": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"categories_id": &graphql.InputObjectFieldConfig{
			Type: graphql.NewList(graphql.String),
		},
		"title": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"chords": &graphql.InputObjectFieldConfig{
			Type: graphql.Boolean,
		},
		"music_sheet": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"music": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"author": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"description": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"bible_verse": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"number": &graphql.InputObjectFieldConfig{
			Type: graphql.Int,
		},
		"text": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"created_at": &graphql.InputObjectFieldConfig{
			Type: graphql.DateTime,
		},
		"updated_at": &graphql.InputObjectFieldConfig{
			Type: graphql.DateTime,
		},
	},
})

func init() {
	Song.AddFieldConfig("songbook", &graphql.Field{
		Type: Songbook,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			songbook_id := p.Source.(models.Song).SongbookID.Hex()
			lang := p.Context.Value("language").(string)

			if songbook_id == "000000000000000000000000" {
				return nil, nil
			}

			songbook := new(models.Songbook).GetSongbookByID(songbook_id, lang)

			if songbook.ID.Hex() == "000000000000000000000000" {
				return nil, nil
			}

			return songbook, nil
		},
	})

	Song.AddFieldConfig("categories", &graphql.Field{
		Type: graphql.NewList(Category),
	})
}
