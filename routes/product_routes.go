package routes

import (
	"go-product-service/controllers"
	"go-product-service/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/api/products", controllers.GetProducts(db))
	r.GET("/api/products/:id", controllers.GetProductById(db))

	admin := r.Group("/api/products")
	admin.Use(middleware.AuthorizeRole("ADMINISTRATOR", "SUPER ADMINISTRATOR"))
	admin.POST("", controllers.CreateProduct(db))
	admin.DELETE("/:id", controllers.DeleteProduct(db))
}
