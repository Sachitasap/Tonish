package routes

import (
	"tonish/backend/handlers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// API prefix
	api := app.Group("/api")
	
	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", handlers.Register)
	auth.Post("/login", handlers.Login)
	
	// Skip authentication for local use
	// api.Use(middleware.AuthRequired)
	
	// User routes
	api.Get("/user/me", handlers.GetCurrentUser)
	
	// Task routes
	tasks := api.Group("/tasks")
	tasks.Get("/", handlers.GetAllTasks)
	tasks.Get("/archived", handlers.GetArchivedTasks)
	tasks.Get("/status", handlers.GetTasksByStatus)
	tasks.Get("/quadrant/:quadrant", handlers.GetTasksByQuadrant)
	tasks.Post("/", handlers.CreateTask)
	tasks.Post("/:id/archive", handlers.ArchiveTask)
	tasks.Post("/:id/restore", handlers.RestoreTask)
	tasks.Delete("/:id/permanent", handlers.PermanentDeleteTask)
	tasks.Get("/:id", handlers.GetTask)
	tasks.Put("/:id", handlers.UpdateTask)
	tasks.Delete("/:id", handlers.DeleteTask)
	
	// Notebook routes
	notebooks := api.Group("/notebooks")
	notebooks.Get("/", handlers.GetAllNotebooks)
	notebooks.Get("/:id", handlers.GetNotebook)
	notebooks.Post("/", handlers.CreateNotebook)
	notebooks.Put("/:id", handlers.UpdateNotebook)
	notebooks.Delete("/:id", handlers.DeleteNotebook)
	
	// Page routes
	pages := api.Group("/pages")
	pages.Get("/search", handlers.SearchPages)
	pages.Get("/:id", handlers.GetPage)
	pages.Post("/", handlers.CreatePage)
	pages.Put("/:id", handlers.UpdatePage)
	pages.Delete("/:id", handlers.DeletePage)
}
