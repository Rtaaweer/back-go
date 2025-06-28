package routes

import (
	"menchaca-api/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// API routes
	api := app.Group("/api")

	// Usuarios routes
	usuarios := api.Group("/usuarios")
	usuarios.Get("/", handlers.GetUsuarios)
	usuarios.Get("/:id", handlers.GetUsuario)
	usuarios.Post("/", handlers.CreateUsuario)
	usuarios.Put("/:id", handlers.UpdateUsuario)
	usuarios.Delete("/:id", handlers.DeleteUsuario)

	// Consultorios routes
	consultorios := api.Group("/consultorios")
	consultorios.Get("/", handlers.GetConsultorios)
	consultorios.Get("/:id", handlers.GetConsultorio)
	consultorios.Post("/", handlers.CreateConsultorio)
	consultorios.Put("/:id", handlers.UpdateConsultorio)
	consultorios.Delete("/:id", handlers.DeleteConsultorio)

	// Consultas routes
	consultas := api.Group("/consultas")
	consultas.Get("/", handlers.GetConsultas)
	consultas.Post("/", handlers.CreateConsulta)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "Hospital System API is running",
		})
	})
}