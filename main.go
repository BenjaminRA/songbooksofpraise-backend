package main

import (
	"flag"
	"fmt"
	"time"

	route_categories "github.com/BenjaminRA/himnario-backend/routes/categories"
	route_countries "github.com/BenjaminRA/himnario-backend/routes/countries"
	route_languages "github.com/BenjaminRA/himnario-backend/routes/languages"
	route_songbooks "github.com/BenjaminRA/himnario-backend/routes/songbooks"
	route_songs "github.com/BenjaminRA/himnario-backend/routes/songs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	var migrate bool
	flag.BoolVar(&migrate, "m", false, "Migrate database")
	flag.Parse()

	if migrate {
		fmt.Println("Migrating database")
		Migrate()
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/songs", route_songs.GetSongs)
	router.GET("/songs/:id", route_songs.GetSongsById)
	router.GET("/songs/:id/music_sheet", route_songs.GetMusicSheet)
	router.GET("/songs/:id/voices/:voice", route_songs.GetVoicesByVoice)

	router.GET("/songbooks", route_songbooks.GetSongbooks)
	router.POST("/songbooks", route_songbooks.PostSongbook)
	router.GET("/songbooks/:id", route_songbooks.GetSongbooksById)
	router.PUT("/songbooks/:id", route_songbooks.UpdateSongbook)
	router.DELETE("/songbooks/:id", route_songbooks.DeleteSongbook)

	router.GET("/categories", route_categories.GetCategories)
	router.POST("/categories", route_categories.PostCategory)
	router.GET("/categories/:id", route_categories.GetCategoriesById)
	router.PUT("/categories/:id", route_categories.UpdateCategory)
	router.DELETE("/categories/:id", route_categories.DeleteCategory)

	router.GET("/languages", route_languages.GetLanguages)
	router.PUT("/languages/:code", route_languages.UpdateLanguage)

	router.GET("/countries", route_countries.GetCountries)

	router.Run("localhost:8080")
}
