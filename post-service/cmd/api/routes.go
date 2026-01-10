package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-TOKEN"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	// Post routes
	mux.Get("/posts", app.GetAllPosts)
	mux.Get("/posts/{id}", app.GetPostByID)
	mux.Get("/posts/author/{authorId}", app.GetPostsByAuthor)
	mux.Post("/posts", app.CreatePost)
	mux.Put("/posts/{id}", app.UpdatePost)
	mux.Delete("/posts/{id}", app.DeletePost)

	return mux
}
