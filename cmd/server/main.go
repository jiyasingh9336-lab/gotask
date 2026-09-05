package main

import (
	"log"
	"net/http"

	"github.com/sudarshanmg/gotask/internal/auth"
	"github.com/sudarshanmg/gotask/internal/task"
	"github.com/sudarshanmg/gotask/pkg/config"
	"github.com/sudarshanmg/gotask/pkg/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	// middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	cfg := config.Load()
	db, err := db.Connect(cfg.URL)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully connected to the database!")
	}

	repo := task.NewRepository(db)
	service := task.NewService(repo)
	taskHandler := task.NewHandler(service)

	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, cfg.JWTSecret)
	authHandler := auth.NewHandler(authService)
	auth.RegisterRoutes(r, authHandler)

	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware(cfg.JWTSecret))
		task.RegisterRoutes(r, taskHandler)
		r.Post("/auth/logout-all", authHandler.LogoutAll)
	})

	log.Printf("Server is listening on port %s...\n", cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Fatal(err)
	}

}
