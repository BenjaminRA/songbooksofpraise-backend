package main

import (
	"flag"
	"fmt"
	"time"

	auth_handler "github.com/BenjaminRA/himnario-backend/handlers/auth"
	files_handler "github.com/BenjaminRA/himnario-backend/handlers/files"
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

	if migrate {
		fmt.Println("Migrating database")
		migration.Migrate()
	}

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// err = smtp.SendMail(
	// 	fmt.Sprintf("%s:smtp", os.Getenv("MAIL_HOST")),
	// 	smtp.PlainAuth(
	// 		os.Getenv("MAIL_IDENTITY"),
	// 		os.Getenv("MAIL_USERNAME"),
	// 		os.Getenv("MAIL_PASSWORD"),
	// 		os.Getenv("MAIL_HOST"),
	// 	),
	// 	"Songbooks Of Praise",
	// 	[]string{"success@simulator.amazonses.com"},
	// 	[]byte("Mensaje"),
	// )

	// if err != nil {
	// 	panic(err)
	// }

	// os.Exit(0)

	schemaConfig := graphql.SchemaConfig{
		Query:    graphql.NewObject(schema.Query),
		Mutation: graphql.NewObject(schema.Mutation),
	}
	schema, _ := graphql.NewSchema(schemaConfig)

	// http.Handle("/graphql", middlewares.FinalMiddleware(h))

	// http.ListenAndServe(":8080", nil)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
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

	router.Run("localhost:8080")
}
