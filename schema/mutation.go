package schema

import (
	resolver_songbooks "github.com/BenjaminRA/himnario-backend/resolvers/songbooks"
	"github.com/BenjaminRA/himnario-backend/types"
	"github.com/graphql-go/graphql"
)

var Mutation = graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{

		// Songbooks
		"updateSongbook": &graphql.Field{
			Type: types.Songbook,
			Args: graphql.FieldConfigArgument{
				"_id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"songbook": &graphql.ArgumentConfig{
					Type: types.NewSongbook,
				},
			},
			Description: "Update an existing songbook",
			Resolve:     resolver_songbooks.UpdateSongbook,
		},
		"createSongbook": &graphql.Field{
			Type: types.Songbook,
			Args: graphql.FieldConfigArgument{
				"songbook": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(types.NewSongbook),
				},
			},
			Description: "Create a new songbook",
			Resolve:     resolver_songbooks.CreateSongbook,
		},
		"deleteSongbook": &graphql.Field{
			Type: types.Songbook,
			Args: graphql.FieldConfigArgument{
				"_id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Description: "Delete a songbook",
			Resolve:     resolver_songbooks.DeleteSongbook,
		},

		// Categories
	},
}
