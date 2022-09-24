package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/BenjaminRA/himnario-backend/middlewares"
	resolver_songbooks "github.com/BenjaminRA/himnario-backend/resolvers/songbooks"
	"github.com/BenjaminRA/himnario-backend/types"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func main() {
	var migrate bool
	flag.BoolVar(&migrate, "m", false, "Migrate database")
	flag.Parse()

	if migrate {
		fmt.Println("Migrating database")
		Migrate()
	}

	fields := graphql.Fields{
		"songbooks": &graphql.Field{
			Type:        graphql.NewList(types.Songbook),
			Description: "Get the list of all songbooks",
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
	}
	rootQuery := graphql.ObjectConfig{Name: "Query", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, _ := graphql.NewSchema(schemaConfig)

	h := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		Playground: true,
	})

	http.Handle("/", middlewares.FinalMiddleware(h))
	http.ListenAndServe(":8080", nil)

	// router := gin.Default()
	// router.Use(cors.New(cors.Config{
	// 	AllowAllOrigins:  true,
	// 	AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
	// 	AllowHeaders:     []string{"*"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))

	// router.GET("/songs", route_songs.GetSongs)
	// router.GET("/songs/:id", route_songs.GetSongsById)
	// router.GET("/songs/:id/music_sheet", route_songs.GetMusicSheet)
	// router.GET("/songs/:id/voices/:voice", route_songs.GetVoicesByVoice)

	// router.GET("/songbooks", route_songbooks.GetSongbooks)
	// router.POST("/songbooks", route_songbooks.PostSongbook)
	// router.GET("/songbooks/:id", route_songbooks.GetSongbooksById)
	// router.PUT("/songbooks/:id", route_songbooks.UpdateSongbook)
	// router.DELETE("/songbooks/:id", route_songbooks.DeleteSongbook)

	// router.GET("/categories", route_categories.GetCategories)
	// router.POST("/categories", route_categories.PostCategory)
	// router.GET("/categories/:id", route_categories.GetCategoriesById)
	// router.PUT("/categories/:id", route_categories.UpdateCategory)
	// router.DELETE("/categories/:id", route_categories.DeleteCategory)

	// router.GET("/languages", route_languages.GetLanguages)
	// router.PUT("/languages/:code", route_languages.UpdateLanguage)

	// router.GET("/countries", route_countries.GetCountries)

	// router.Run("localhost:8080")
}
