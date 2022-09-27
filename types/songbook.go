package types

import (
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/graphql-go/graphql"
)

var Songbook = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Songbook",
		Fields: graphql.Fields{
			"_id": &graphql.Field{
				Type: graphql.ID,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Source.(models.Songbook).ID.Hex()
					return id, nil
				},
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"language_code": &graphql.Field{
				Type: graphql.String,
			},
			"country_code": &graphql.Field{
				Type: graphql.String,
			},
			"numeration": &graphql.Field{
				Type: graphql.Boolean,
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

var NewSongbook = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "NewSongbook",
	Fields: graphql.InputObjectConfigFieldMap{
		"title": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"description": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"language_code": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"country_code": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"numeration": &graphql.InputObjectFieldConfig{
			Type: graphql.Boolean,
		},
	},
})

func init() {
	Songbook.AddFieldConfig("language", &graphql.Field{
		Type: Language,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			language_code := p.Source.(models.Songbook).LanguageCode
			lang := p.Context.Value("language").(string)
			if language_code == "" {
				return nil, nil
			}

			language := new(models.Language).GetLanguageByCode(language_code, lang)
			if language.ID.Hex() == "000000000000000000000000" {
				return nil, nil
			}

			return language, nil
		},
	})

	Songbook.AddFieldConfig("country", &graphql.Field{
		Type: Country,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			country_code := p.Source.(models.Songbook).CountryCode
			lang := p.Context.Value("language").(string)
			if country_code == "" {
				return nil, nil
			}

			country := new(models.Country).GetCountryByCode(country_code, lang)
			if country.ID.Hex() == "000000000000000000000000" {
				return nil, nil
			}

			return country, nil
		},
	})

	Songbook.AddFieldConfig("songs", &graphql.Field{
		Type: graphql.NewList(Song),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id := p.Source.(models.Songbook).ID.Hex()
			songs := new(models.Songbook).GetSongs(id)
			return songs, nil
		},
	})

	Songbook.AddFieldConfig("categories", &graphql.Field{
		Type: graphql.NewList(Category),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id := p.Source.(models.Songbook).ID.Hex()
			lang := p.Context.Value("language").(string)
			songbook := new(models.Songbook).GetSongbookByID(id, lang)

			return songbook.Categories, nil
		},
	})
}
