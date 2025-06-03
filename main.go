package main

import (
	"os"

	"go-product-service/config"
	"go-product-service/controllers"
	"go-product-service/middleware"
	"go-product-service/models"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func main() {
	// Cargar configuración
	config.Load()

	// Conexión a la base de datos
	db, err := gorm.Open(sqlserver.Open(config.Cfg.DSN), &gorm.Config{})
	if err != nil {
		panic("No se pudo conectar a la base de datos: " + err.Error())
	}

	// Migraciones
	db.AutoMigrate(&models.Product{})

	// Configurar Gin
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Origen del frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Static("/uploads", "./uploads")

	api := router.Group("/api/products")
	api.Use(middleware.JWTAuthMiddleware())

	// Rutas disponibles para todos los usuarios autenticados
	userRoutes := api.Group("/")
	userRoutes.Use(middleware.AuthorizeRole("USER", "ADMINISTRATOR", "SUPER ADMINISTRATOR"))
	userRoutes.GET("/", controllers.GetProducts(db))
	userRoutes.GET("/:id", controllers.GetProductById(db))

	// Rutas restringidas a administradores
	adminRoutes := api.Group("/")
	adminRoutes.Use(middleware.AuthorizeRole("ADMINISTRATOR", "SUPER ADMINISTRATOR"))
	adminRoutes.POST("/", controllers.CreateProduct(db))
	adminRoutes.PUT("/:id", controllers.UpdateProduct(db))
	adminRoutes.DELETE("/:id", controllers.DeleteProduct(db))

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "2222"
	}
	router.Run(":" + port)
}
