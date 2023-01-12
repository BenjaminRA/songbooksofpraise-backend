package schema

import (
	resolver_categories "github.com/BenjaminRA/himnario-backend/resolvers/categories"
	resolver_songbooks "github.com/BenjaminRA/himnario-backend/resolvers/songbooks"
	resolver_songs "github.com/BenjaminRA/himnario-backend/resolvers/songs"
	"github.com/BenjaminRA/himnario-backend/types"
	"github.com/graphql-go/graphql"
)

var Mutation = graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{

		// Songbooks
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
		"createCategory": &graphql.Field{
			Type: types.Category,
			Args: graphql.FieldConfigArgument{
				"category": &graphql.ArgumentConfig{
					Type: types.NewCategory,
				},
			},
			Description: "Create a new category",
			Resolve:     resolver_categories.CreateCategory,
		},
		"updateCategory": &graphql.Field{
			Type: types.Category,
			Args: graphql.FieldConfigArgument{
				"_id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"category": &graphql.ArgumentConfig{
					Type: types.NewCategory,
				},
			},
			Description: "Updates an existing category",
			Resolve:     resolver_categories.UpdateCategory,
		},
		"deleteCategory": &graphql.Field{
			Type: types.Category,
			Args: graphql.FieldConfigArgument{
				"_id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Description: "Delete a category",
			Resolve:     resolver_categories.DeleteCategory,
		},

		// Languages
		"createLanguage": &graphql.Field{
			Type: types.Language,
			Args: graphql.FieldConfigArgument{
				"language": &graphql.ArgumentConfig{
					Type: graphql.NewList(types.NewLanguage),
				},
			},
			Description: "Create a new language",
		},
		"updateLanguage": &graphql.Field{
			Type: types.Language,
			Args: graphql.FieldConfigArgument{
				"_id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"language": &graphql.ArgumentConfig{
					Type: types.NewLanguage,
				},
			},
			Description: "Update an existing language",
		},
		"deleteLanguage": &graphql.Field{
			Type: types.Language,
			Args: graphql.FieldConfigArgument{
				"_id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Description: "Delete a language",
		},

		// Countries
		"createCountry": &graphql.Field{
			Type: types.Country,
			Args: graphql.FieldConfigArgument{
				"country": &graphql.ArgumentConfig{
					Type: types.NewCountry,
				},
			},
			Description: "Create a new country",
		},
		"updateCountry": &graphql.Field{
			Type: types.Country,
			Args: graphql.FieldConfigArgument{
				"_id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"country": &graphql.ArgumentConfig{
					Type: types.NewCountry,
				},
			},
			Description: "Update an existing country",
		},
		"deleteCountry": &graphql.Field{
			Type: types.Country,
			Args: graphql.FieldConfigArgument{
				"_id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Description: "Delete a country",
		},

		// Songs
		"createSong": &graphql.Field{
			Type: types.Song,
			Args: graphql.FieldConfigArgument{
				"song": &graphql.ArgumentConfig{
					Type: types.NewSong,
				},
			},
			Description: "Create a new song",
			Resolve:     resolver_songs.CreateSong,
		},

		"updateSong": &graphql.Field{
			Type: types.Song,
			Args: graphql.FieldConfigArgument{
				"_id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"song": &graphql.ArgumentConfig{
					Type: types.NewSong,
				},
			},
			Description: "Update an existing song",
			Resolve:     resolver_songs.UpdateSong,
		},

		"deleteSong": &graphql.Field{
			Type: types.Song,
			Args: graphql.FieldConfigArgument{
				"_id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Description: "Delete a song",
		},
	},
}
