package routes

import (
	"hospital-system/handlers"
	"hospital-system/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	api := app.Group("/api/v1", middleware.GeneralAPIRateLimit())

	auth := api.Group("/auth")
	auth.Post("/register", middleware.RegisterRateLimit(), handlers.Register)
	auth.Post("/login", middleware.AuthRateLimit(), handlers.Login)
	auth.Post("/refresh", middleware.AuthRateLimit(), handlers.RefreshToken)

	protected := api.Group("/", middleware.JWTMiddleware(), middleware.MedicalRateLimit())



	mfa := protected.Group("/mfa")
	mfa.Get("/status", handlers.GetMFAStatus)
	mfa.Post("/enable", handlers.EnableMFA)
	mfa.Post("/verify", handlers.VerifyMFA)
	mfa.Post("/disable", handlers.DisableMFA)

	usuarios := protected.Group("/usuarios")
	usuarios.Get("/", handlers.GetUsuarios)
	usuarios.Get("/:id", handlers.GetUsuario)
	usuarios.Post("/", middleware.RequireRole("admin"), middleware.AdminRateLimit(), handlers.CreateUsuario)
	usuarios.Put("/:id", handlers.UpdateUsuario)
	usuarios.Delete("/:id", middleware.RequireRole("admin"), middleware.AdminRateLimit(), handlers.DeleteUsuario)

	consultorios := protected.Group("/consultorios")
	consultorios.Get("/", handlers.GetConsultorios)
	consultorios.Get("/:id", handlers.GetConsultorio)
	consultorios.Post("/", middleware.AdminRateLimit(), handlers.CreateConsultorio)
	consultorios.Put("/:id", middleware.AdminRateLimit(), handlers.UpdateConsultorio)
	consultorios.Delete("/:id", middleware.AdminRateLimit(), handlers.DeleteConsultorio)

	consultas := protected.Group("/consultas")
	consultas.Get("/", handlers.GetConsultas)
	consultas.Post("/", handlers.CreateConsulta)

	expedientes := protected.Group("/expedientes")
	expedientes.Get("/", handlers.GetExpedientes)
	expedientes.Get("/:id", handlers.GetExpediente)
	expedientes.Post("/", handlers.CreateExpediente)
	expedientes.Put("/:id", handlers.UpdateExpediente)
	expedientes.Delete("/:id", middleware.AdminRateLimit(), handlers.DeleteExpediente)

	horarios := protected.Group("/horarios")
	horarios.Get("/", handlers.GetHorarios)
	horarios.Get("/:id", handlers.GetHorario)
	horarios.Post("/", handlers.CreateHorario)
	horarios.Put("/:id", handlers.UpdateHorario)
	horarios.Delete("/:id", middleware.AdminRateLimit(), handlers.DeleteHorario)

	recetas := protected.Group("/recetas")
	recetas.Get("/", handlers.GetRecetas)
	recetas.Get("/:id", handlers.GetReceta)
	recetas.Post("/", handlers.CreateReceta)
	recetas.Put("/:id", handlers.UpdateReceta)
	recetas.Delete("/:id", middleware.AdminRateLimit(), handlers.DeleteReceta)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Hospital System API is running",
		})
	})

	api.Get("/rate-limit/status", middleware.JWTMiddleware(), middleware.RequireRole("admin"), func(c *fiber.Ctx) error {
		clientIP := c.Query("ip")
		if clientIP == "" {
			clientIP = c.IP()
		}

		status := middleware.GetRateLimitStatus(clientIP)
		return c.JSON(fiber.Map{
			"ip":     clientIP,
			"status": status,
		})
	})
}
