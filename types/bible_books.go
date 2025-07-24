package types

import (
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/graphql-go/graphql"
)

var BibleBook = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "BibleBook",
		Fields: graphql.Fields{
			"_id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Source.(models.BibleBook).ID.Hex()
					return id, nil
				},
			},
			"code": &graphql.Field{
				Type: graphql.String,
			},
			"language_code": &graphql.Field{
				Type: graphql.String,
			},
			"book": &graphql.Field{
				Type: graphql.String,
			},
			"testament": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
