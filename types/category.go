package types

import (
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/graphql-go/graphql"
)

var Category = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Category",
		Fields: graphql.Fields{
			"_id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Source.(models.Category).ID.Hex()
					return id, nil
				},
			},
			"all": &graphql.Field{
				Type: graphql.Boolean,
			},
			"category": &graphql.Field{
				Type: graphql.String,
			},
			"songbook_id": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Source.(models.Category).SongbookID.Hex()
					return id, nil
				},
			},
			"parent_id": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Source.(models.Category).ParentID.Hex()
					return id, nil
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

func init() {
	Category.AddFieldConfig("children", &graphql.Field{
		Type: graphql.NewList(Category),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id := p.Source.(models.Category).ID
			category := new(models.Category).GetCategoryById(id.Hex())
			return category.Children, nil
		},
	})

	Category.AddFieldConfig("songs", &graphql.Field{
		Type: graphql.NewList(Song),
	})
}
