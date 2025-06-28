package main

import (
    "menchaca-api/config"
    "menchaca-api/routes"
    "log"
    "os"

    "github.com/gofiber/fiber/v2"
    "github.com/joho/godotenv"
)

func main() {
    // Cargar variables de entorno
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found")
    }

    // Conectar a la base de datos
    config.ConnectDB()

    // Crear aplicación Fiber
    app := fiber.New(fiber.Config{
        AppName: "Proyecto Menchaca",
    })

    // Configurar rutas
    routes.SetupRoutes(app)

    // Obtener puerto del entorno o usar 3000 por defecto
    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }

    log.Printf("🚀 Server starting on port %s", port)
    log.Fatal(app.Listen(":" + port))
}