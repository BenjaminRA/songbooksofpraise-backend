package schema

import (
	resolver_categories "github.com/BenjaminRA/himnario-backend/resolvers/categories"
	resolver_songbooks "github.com/BenjaminRA/himnario-backend/resolvers/songbooks"
	resolver_songs "github.com/BenjaminRA/himnario-backend/resolvers/songs"
	"github.com/BenjaminRA/himnario-backend/types"
	"github.com/graphql-go/graphql"
)

var fields = graphql.Fields{
	"songbooks": &graphql.Field{
		Type:        graphql.NewList(types.Songbook),
		Description: "Get a list of all songbooks",
		Resolve:     resolver_songbooks.GetSongbooks,
	},
	"songbook": &graphql.Field{
		Type:        types.Songbook,
		Description: "Get a specific songbook",
		Args: graphql.FieldConfigArgument{
			"_id": &graphql.ArgumentConfig{
				Type: graphql.ID,
			},
		},
		Resolve: resolver_songbooks.GetSongbook,
	},
	"categories": &graphql.Field{
		Type:        graphql.NewList(types.Category),
		Description: "Get a list of all the categories",
		Resolve:     resolver_categories.GetCategories,
	},
	"category": &graphql.Field{
		Type:        types.Category,
		Description: "Get a specific category",
		Args: graphql.FieldConfigArgument{
			"_id": &graphql.ArgumentConfig{
				Type: graphql.ID,
			},
		},
		Resolve: resolver_categories.GetCategory,
	},
	"songs": &graphql.Field{
		Type:        graphql.NewList(types.Song),
		Description: "Get a list of all the songs",
		Args: graphql.FieldConfigArgument{
			"songbook_id": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"category_id": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: resolver_songs.GetSongs,
	},
	"song": &graphql.Field{
		Type:        types.Song,
		Description: "Get a specific song",
		Args: graphql.FieldConfigArgument{
			"_id": &graphql.ArgumentConfig{
				Type: graphql.ID,
			},
		},
		Resolve: resolver_songs.GetSongById,
	},
}
var Query = graphql.ObjectConfig{Name: "Query", Fields: fields}
