package main

import (
	"database/sql"
	"github.com/charlesozo/blog-aggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("database environment variable is not set")
	}

	conn, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("Connect to database", err)
	}

	queries := database.New(conn)
	apicfg := apiConfig{
		DB: queries,
	}
	go startScraping(queries, 10, time.Minute)

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	v1Router := chi.NewRouter()
	v1Router.Post("/user", apicfg.handlerCreateUser)
	v1Router.Get("/user", apicfg.MiddleWareAuth(apicfg.handlerGetUser))
	v1Router.Post("/feeds", apicfg.MiddleWareAuth(apicfg.handleCreateFeeds))
	v1Router.Get("/feeds", apicfg.handleGetAllFeeds)
	v1Router.Delete("/feed_follows/{feedFollowID}", apicfg.MiddleWareAuth(apicfg.handleDeleteFeedFollow))
	v1Router.Post("/feed_follows", apicfg.MiddleWareAuth(apicfg.handleCreateFeedFollow))
	v1Router.Get("/feed_follows", apicfg.MiddleWareAuth(apicfg.handleGetUserFeedFollows))
	v1Router.Get("/posts", apicfg.MiddleWareAuth(apicfg.handlerGetPostsForUser))
	router.Mount("/v1", v1Router)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	log.Printf("Serving on port: %s", port)
	log.Fatal(srv.ListenAndServe())
}
