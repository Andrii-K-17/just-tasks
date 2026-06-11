package router

import (
	"net/http"
	"time"

	"github.com/Andrii-K-17/just-tasks/internal/config"
	"github.com/Andrii-K-17/just-tasks/internal/handlers"
	"github.com/Andrii-K-17/just-tasks/internal/middleware"
	"github.com/Andrii-K-17/just-tasks/internal/repository"
	"github.com/Andrii-K-17/just-tasks/internal/services"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
)

// New initializes and configures the main application router.
func New(db *sqlx.DB, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.AllowedOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)

	authSvc := services.NewAuthService(userRepo, refreshTokenRepo)
	taskSvc := services.NewTaskService(taskRepo)
	categorySvc := services.NewCategoryService(categoryRepo)
	aiSvc := services.NewAIService(cfg.GroqAPIKey)

	auth := handlers.NewAuthHandler(
		authSvc,
		cfg.JWTSecret,
		cfg.JWTExpiry,
		cfg.RefreshExpiry,
		cfg.IsProd(),
	)
	tasks := handlers.NewTaskHandler(taskSvc)
	categories := handlers.NewCategoryHandler(categorySvc)
	ai := handlers.NewAIHandler(aiSvc)

	r.Route("/api", func(r chi.Router) {
		// Public authentication routes.
		r.Post("/register", auth.Register)
		r.Post("/login", auth.Login)
		r.Post("/logout", auth.Logout)
		r.Post("/refresh", auth.Refresh)

		// Protected routes requiring JWT authentication.
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(cfg.JWTSecret))

			r.Get("/me", auth.Me)
			r.Delete("/account", auth.DeleteAccount)

			r.Get("/tasks", tasks.GetTasks)
			r.Post("/tasks", tasks.CreateTask)
			r.Put("/tasks/reorder", tasks.ReorderTasks)
			r.Patch("/tasks/{id}", tasks.UpdateTask)
			r.Delete("/tasks/{id}", tasks.DeleteTask)

			r.Get("/categories", categories.GetCategories)
			r.Post("/categories", categories.CreateCategory)
			r.Delete("/categories/{id}", categories.DeleteCategory)

			r.Post("/tasks/{id}/collaborators", tasks.AddCollaborator)
			r.Delete("/tasks/{id}/collaborators/{collabId}", tasks.RemoveCollaborator)

			r.Post("/ai/generate", ai.GenerateTasks)
		})
	})

	return r
}
