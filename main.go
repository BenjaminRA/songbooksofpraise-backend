package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	auth_handler "github.com/BenjaminRA/himnario-backend/handlers/auth"
	files_handler "github.com/BenjaminRA/himnario-backend/handlers/files"
	songbooks_handler "github.com/BenjaminRA/himnario-backend/handlers/songbooks"
	song_handler "github.com/BenjaminRA/himnario-backend/handlers/songs"
	"github.com/BenjaminRA/himnario-backend/middlewares"
	migration "github.com/BenjaminRA/himnario-backend/migration"
	"github.com/BenjaminRA/himnario-backend/schema"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/joho/godotenv"
)

func main() {
	var migrate bool
	flag.BoolVar(&migrate, "m", false, "Migrate database")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	if migrate {
		fmt.Println("Migrating database")
		migration.Migrate()
	}

	schemaConfig := graphql.SchemaConfig{
		Query:    graphql.NewObject(schema.Query),
		Mutation: graphql.NewObject(schema.Mutation),
	}
	schema, _ := graphql.NewSchema(schemaConfig)

	// http.Handle("/graphql", middlewares.FinalMiddleware(h))

	// http.ListenAndServe(":8080", nil)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://admin:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Accept-Language"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	h := handler.New(&handler.Config{
		Schema: &schema,
	})

	// Takes the http handler for the graphl schema and serves it in a gin handler
	graphqlHandler := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			h.ServeHTTP(c.Writer, c.Request)
		}
	}

	// Sets the language variable to retrieve the information accordingly
	router.Use(middlewares.LanguageParser())
	router.Use(middlewares.CheckAuthentication())

	router.POST("/graphql", graphqlHandler())

	// playgroundH := handler.New(&handler.Config{
	// 	Schema:     &schema,
	// 	Pretty:     true,
	// 	Playground: true,
	// })

	// playgroundHandler := func() gin.HandlerFunc {
	// 	return func(c *gin.Context) {
	// 		playgroundH.ServeHTTP(c.Writer, c.Request)
	// 	}
	// }
	// router.GET("/graphql", playgroundHandler())

	router.GET("/songs/:id/music", song_handler.GetMusic)
	router.GET("/songs/:id/music_only", song_handler.GetMusicOnly)
	router.GET("/songs/:id/music_sheet", song_handler.GetMusicSheet)
	router.GET("/songs/:id/voices/:voice", song_handler.GetVoicesByVoice)
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

	router.Run(fmt.Sprintf("0.0.0.0:%s", os.Getenv("BACKEND_PORT")))
}
