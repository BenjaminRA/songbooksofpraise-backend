package main

import (
	routes "github.com/BenjaminRA/himnario-backend/routes/songs"
	"github.com/gin-gonic/gin"
)

func main() {
	// Migrate()

	router := gin.Default()
	router.GET("/songs", routes.GetSongs)

	router.Run("localhost:8080")
}
