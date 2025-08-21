package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	app_categories_handler "github.com/BenjaminRA/himnario-backend/handlers/app/categories"
	app_songbooks_handler "github.com/BenjaminRA/himnario-backend/handlers/app/songbooks"
	app_songs_handler "github.com/BenjaminRA/himnario-backend/handlers/app/songs"
	auth_handler "github.com/BenjaminRA/himnario-backend/handlers/auth"
	categories_handler "github.com/BenjaminRA/himnario-backend/handlers/categories"
	files_handler "github.com/BenjaminRA/himnario-backend/handlers/files"
	songbooks_handler "github.com/BenjaminRA/himnario-backend/handlers/songbooks"
	songs_handler "github.com/BenjaminRA/himnario-backend/handlers/songs"
	"github.com/BenjaminRA/himnario-backend/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// var migrate bool
	// flag.BoolVar(&migrate, "m", false, "Migrate database")
	// flag.Parse()

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// if migrate {
	// 	fmt.Println("Migrating database")
	// 	migration.Migrate()
	// }

	// schemaConfig := graphql.SchemaConfig{
	// 	Query:    graphql.NewObject(schema.Query),
	// 	Mutation: graphql.NewObject(schema.Mutation),
	// }
	// schema, _ := graphql.NewSchema(schemaConfig)

	// http.Handle("/graphql", middlewares.FinalMiddleware(h))

	// http.ListenAndServe(":8080", nil)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://admin:3000",
			"https://admin.songbooksofpraise.com",
			"https://backend.songbooksofpraise.com",
			"https://songbooksofpraise.com",
			"https://www.songbooksofpraise.com",
		},
		AllowMethods: []string{"PUT", "PATCH", "POST", "GET", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Content-Type",
			"Accept-Language",
			"Authorization",
			"X-API-Token",
			"Accept",
			"Origin",
			"X-Requested-With",
			"Cache-Control",
			"Cookie",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// h := handler.New(&handler.Config{
	// 	Schema: &schema,
	// })

	// Takes the http handler for the graphl schema and serves it in a gin handler
	// graphqlHandler := func() gin.HandlerFunc {
	// 	return func(c *gin.Context) {
	// 		h.ServeHTTP(c.Writer, c.Request)
	// 	}
	// }

	// Sets the language variable to retrieve the information accordingly
	router.Use(middlewares.LanguageParser())
	router.Use(middlewares.CheckAuthentication())
	router.Use(gin.RecoveryWithWriter(gin.DefaultErrorWriter, func(c *gin.Context, err interface{}) {
		// Log the panic and stack trace
		fmt.Printf("Panic recovered: %v\nStack trace:\n%s\n", err, debug.Stack())

		// Return a custom error response with 500 status
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "An unexpected error occurred. Please try again later.",
		})
		c.Abort() // Stop further processing of the request
	}))

	// Health check endpoint (ADD THIS)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
		})
	})

	router.POST("/files", files_handler.PostFile)

	router.POST("/login", auth_handler.Login)
	router.POST("/register", auth_handler.Register)
	router.POST("/logout", auth_handler.Logout)
	router.POST("/auth/user", auth_handler.GetUser)
	router.PUT("/auth/user", auth_handler.UpdateUser)
	router.DELETE("/auth/user", auth_handler.DeleteUser)
	router.GET("/auth/users", auth_handler.GetUsers)
	router.POST("/auth/verification", auth_handler.VerifyUserEmail)
	router.POST("/auth/verification/resend", auth_handler.EmailVerificationResend)

	router.POST("/songbooks/:id/verify", songbooks_handler.VerifySongbook)
	router.POST("/songbooks/:id/send-to-verify", songbooks_handler.SendToVerifySongbook)
	router.POST("/songbooks/:id/reject", songbooks_handler.RejectSongbook)

	// Songbooks endpoints
	router.GET("/songbooks", songbooks_handler.GetSongbooks)
	router.GET("/songbooks/:id", songbooks_handler.GetSongbookByID)
	router.POST("/songbooks", songbooks_handler.CreateSongbook)
	router.PUT("/songbooks/:id", songbooks_handler.UpdateSongbook)
	router.DELETE("/songbooks/:id", songbooks_handler.DeleteSongbook)

	// Categories endpoints
	router.GET("/songbooks/:id/categories", categories_handler.GetCategories)
	router.POST("/songbooks/:id/categories", categories_handler.CreateCategory)
	router.GET("/songbooks/:id/categories/:category_id", categories_handler.GetCategoryByID)
	router.PUT("/songbooks/:id/categories/:category_id", categories_handler.UpdateCategory)
	router.DELETE("/songbooks/:id/categories/:category_id", categories_handler.DeleteCategory)

	// Songs endpoints
	router.GET("/songbooks/:id/categories/:category_id/songs/:song_id", songs_handler.GetSongByID)
	router.PUT("/songs/:song_id", songs_handler.UpdateSong)
	router.DELETE("/songs/:song_id", songs_handler.DeleteSong)
	router.POST("/songs", songs_handler.CreateSong)

	// App Endpoints
	router.GET("/app/songbooks", app_songbooks_handler.GetSongbooks)
	router.GET("/app/songbooks/:id/export", app_songbooks_handler.ExportSongbookByID)

	router.GET("/app/songbooks/:id/categories", app_categories_handler.GetCategories)

	router.GET("/app/songbooks/:id/categories/:category_id", app_categories_handler.GetCategoryByID)

	router.GET("/app/songs/:song_id", app_songs_handler.GetSongByID)

	router.Run(fmt.Sprintf("0.0.0.0:%s", os.Getenv("BACKEND_PORT")))
}
