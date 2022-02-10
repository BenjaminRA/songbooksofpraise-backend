package main

import (
	route_songbooks "github.com/BenjaminRA/himnario-backend/routes/songbooks"
	route_songs "github.com/BenjaminRA/himnario-backend/routes/songs"
	"github.com/gin-gonic/gin"
)

func main() {
	// Migrate()

	router := gin.Default()
	router.GET("/songs", route_songs.GetSongs)
	router.GET("/songs/:id", route_songs.GetSongsById)

	router.GET("/songbooks", route_songbooks.GetSongbooks)

	router.Run("localhost:8080")
}
