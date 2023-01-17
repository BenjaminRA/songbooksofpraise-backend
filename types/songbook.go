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
			"owners_id": &graphql.Field{
				Type: graphql.NewList(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					temp := []string{}
					ids := p.Source.(models.Songbook).OwnersID

					for _, id := range ids {
						temp = append(temp, id.Hex())
					}

					return temp, nil
				},
			},
			"editors_id": &graphql.Field{
				Type: graphql.NewList(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					temp := []string{}
					ids := p.Source.(models.Songbook).EditorsID

					for _, id := range ids {
						temp = append(temp, id.Hex())
					}

					return temp, nil
				},
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
		"editors_id": &graphql.InputObjectFieldConfig{
			Type: graphql.NewList(graphql.String),
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

			language, err := new(models.Language).GetLanguageByCode(language_code, lang)
			if err != nil {
				return nil, err
			}

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

			country, err := new(models.Country).GetCountryByCode(country_code, lang)
			if err != nil {
				return nil, err
			}

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
			songs, err := new(models.Songbook).GetSongs(id)
			if err != nil {
				return nil, err
			}

			return songs, nil
		},
	})

	Songbook.AddFieldConfig("editors", &graphql.Field{
		Type: graphql.NewList(User),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			ids := p.Source.(models.Songbook).EditorsID
			users := []models.User{}

			for _, id := range ids {
				user, err := new(models.User).GetUserById(id.Hex())
				if err != nil {
					return nil, err
				}

				users = append(users, user)
			}

			return users, nil
		},
	})

	Songbook.AddFieldConfig("owners", &graphql.Field{
		Type: graphql.NewList(User),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			ids := p.Source.(models.Songbook).OwnersID
			users := []models.User{}

			for _, id := range ids {
				user, err := new(models.User).GetUserById(id.Hex())
				if err != nil {
					return nil, err
				}

				users = append(users, user)
			}

			return users, nil
		},
	})

	Songbook.AddFieldConfig("categories", &graphql.Field{
		Type: graphql.NewList(Category),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id := p.Source.(models.Songbook).ID.Hex()
			lang := p.Context.Value("language").(string)
			songbook, err := new(models.Songbook).GetSongbookByID(id, lang)
			if err != nil {
				return nil, err
			}

			return songbook.Categories, nil
		},
	})
}
