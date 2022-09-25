package types

import (
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/graphql-go/graphql"
)

var Language = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Language",
		Fields: graphql.Fields{
			"_id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Source.(models.Language).ID.Hex()
					return id, nil
				},
			},
			"language": &graphql.Field{
				Type: graphql.String,
			},
			"reader_code": &graphql.Field{
				Type: graphql.String,
			},
			"code": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
