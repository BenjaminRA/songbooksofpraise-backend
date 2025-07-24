package schema

import (
	resolver_authors "github.com/BenjaminRA/himnario-backend/resolvers/authors"
	resolver_bible_books "github.com/BenjaminRA/himnario-backend/resolvers/bible_books"
	resolver_categories "github.com/BenjaminRA/himnario-backend/resolvers/categories"
	resolver_countries "github.com/BenjaminRA/himnario-backend/resolvers/countries"
	resolver_languages "github.com/BenjaminRA/himnario-backend/resolvers/languages"
	resolver_songbooks "github.com/BenjaminRA/himnario-backend/resolvers/songbooks"
	resolver_songs "github.com/BenjaminRA/himnario-backend/resolvers/songs"
	"github.com/BenjaminRA/himnario-backend/types"
	"github.com/graphql-go/graphql"
)

var Query = graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{

		// Songbooks
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

		// Categories
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

		// Songs
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

		// Languages
		"languages": &graphql.Field{
			Type:        graphql.NewList(types.Language),
			Description: "Get a list of all the languages",
			Resolve:     resolver_languages.GetLanguages,
		},

		// Countries
		"countries": &graphql.Field{
			Type:        graphql.NewList(types.Country),
			Description: "Get a list of all the countries",
			Resolve:     resolver_countries.GetCountries,
		},

		// Bible Books
		"bible_books": &graphql.Field{
			Type:        graphql.NewList(types.BibleBook),
			Description: "Get a list of all bible books",
			Resolve:     resolver_bible_books.GetAllBibleBooks,
		},
		"bible_book": &graphql.Field{
			Type:        graphql.NewList(types.BibleBook),
			Description: "Get a specific bible books",
			Args: graphql.FieldConfigArgument{
				"code": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: resolver_bible_books.GetAllBibleBooks,
		},

		// Authors
		"authors": &graphql.Field{
			Type:        graphql.NewList(types.Author),
			Description: "Get a list of all authors",
			Resolve:     resolver_authors.GetAuthors,
		},
	},
}
