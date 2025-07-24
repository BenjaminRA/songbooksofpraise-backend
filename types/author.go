package types

import (
	"github.com/graphql-go/graphql"
)

var Author = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Author",
		Fields: graphql.Fields{
			"author": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
