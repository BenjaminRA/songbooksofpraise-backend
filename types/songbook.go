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
				Type: graphql.NewNonNull(graphql.ID),
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
			"songs": &graphql.Field{
				Type: graphql.NewList(Song),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Source.(models.Songbook).ID.Hex()
					songs := new(models.Songbook).GetSongs(id)
					return songs, nil
				},
			},
			"categories": &graphql.Field{
				Type: graphql.NewList(Category),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Source.(models.Songbook).ID.Hex()
					lang := p.Context.Value("language").(string)
					songbook := new(models.Songbook).GetSongbookByID(id, lang)

					return songbook.Categories, nil
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
