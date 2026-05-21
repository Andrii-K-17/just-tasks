package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"

	"github.com/andriik17/just-tasks/internal/handlers"
	"github.com/andriik17/just-tasks/internal/middleware"
)

// New initializes and configures the main application router.
func New(
	db *sqlx.DB,
	jwtSecret string,
	jwtExpiry time.Duration,
	allowedOrigin string,
	groqAPIKey string,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{allowedOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	auth := handlers.NewAuthHandler(db, jwtSecret, jwtExpiry)
	tasks := handlers.NewTaskHandler(db)
	categories := handlers.NewCategoryHandler(db)
	ai := handlers.NewAIHandler(groqAPIKey)

	r.Route("/api", func(r chi.Router) {
		// Public authentication routes.
		r.Post("/register", auth.Register)
		r.Post("/login", auth.Login)
		r.Post("/logout", auth.Logout)

		// Protected routes requiring JWT authentication.
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(jwtSecret))

			r.Get("/me", auth.Me)
			r.Delete("/account", auth.DeleteAccount)

			r.Get("/tasks", tasks.GetTasks)
			r.Post("/tasks", tasks.CreateTask)
			r.Put("/tasks/reorder", tasks.ReorderTasks)
			r.Put("/tasks/{id}", tasks.UpdateTask)
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
