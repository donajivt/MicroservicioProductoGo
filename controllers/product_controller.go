package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"go-product-service/dtos"
	"go-product-service/models"
	"go-product-service/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetProducts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var products []models.Product
		db.Find(&products)
		c.JSON(http.StatusOK, dtos.ResponseDto{IsSuccess: true, Result: products})
	}
}

func GetProductById(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var product models.Product
		if err := db.First(&product, id).Error; err != nil {
			c.JSON(http.StatusNotFound, dtos.ResponseDto{IsSuccess: false, Message: "Producto no encontrado"})
			return
		}
		c.JSON(http.StatusOK, dtos.ResponseDto{IsSuccess: true, Result: product})
	}
}

func CreateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto dtos.ProductDto
		contentType := c.ContentType()

		if strings.HasPrefix(contentType, "application/json") {
			if err := c.ShouldBindJSON(&dto); err != nil {
				c.JSON(http.StatusBadRequest, dtos.ResponseDto{IsSuccess: false, Message: err.Error()})
				return
			}
		} else if strings.HasPrefix(contentType, "multipart/form-data") {
			dto.Name = c.PostForm("name")
			dto.Price, _ = strconv.ParseFloat(c.PostForm("price"), 64)
			dto.Description = c.PostForm("description")
			dto.CategoryName = c.PostForm("category_name")
			dto.Stock, _ = strconv.Atoi(c.PostForm("stock"))
			file, _ := c.FormFile("image")
			if file != nil {
				filename := dto.Name + "_" + file.Filename
				path, err := utils.SaveImage(file, filename)
				if err == nil {
					dto.ImageLocalPath = path
					dto.ImageUrl = "http://localhost:2222/uploads/images/" + filename
				}
			}
		} else {
			c.JSON(http.StatusUnsupportedMediaType, dtos.ResponseDto{IsSuccess: false, Message: "Formato de contenido no soportado"})
			return
		}

		product := models.Product{
			Name:           dto.Name,
			Price:          dto.Price,
			Description:    dto.Description,
			CategoryName:   dto.CategoryName,
			ImageUrl:       dto.ImageUrl,
			ImageLocalPath: dto.ImageLocalPath,
			Stock:          dto.Stock,
		}

		db.Create(&product)
		c.JSON(http.StatusCreated, dtos.ResponseDto{IsSuccess: true, Result: product})
	}
}

func DeleteProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var product models.Product
		if err := db.First(&product, id).Error; err != nil {
			c.JSON(http.StatusNotFound, dtos.ResponseDto{IsSuccess: false, Message: "No encontrado"})
			return
		}
		db.Delete(&product)
		c.JSON(http.StatusOK, dtos.ResponseDto{IsSuccess: true, Message: "Eliminado"})
	}
}

func UpdateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var product models.Product
		if err := db.First(&product, id).Error; err != nil {
			c.JSON(http.StatusNotFound, dtos.ResponseDto{IsSuccess: false, Message: "Producto no encontrado"})
			return
		}

		contentType := c.ContentType()

		if strings.HasPrefix(contentType, "application/json") {
			var dto dtos.ProductDto
			if err := c.ShouldBindJSON(&dto); err != nil {
				c.JSON(http.StatusBadRequest, dtos.ResponseDto{IsSuccess: false, Message: err.Error()})
				return
			}
			product.Name = dto.Name
			product.Price = dto.Price
			product.Description = dto.Description
			product.CategoryName = dto.CategoryName
			product.Stock = dto.Stock
		} else if strings.HasPrefix(contentType, "multipart/form-data") {
			// Parsear el formulario primero
			if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32 MB
				c.JSON(http.StatusBadRequest, dtos.ResponseDto{IsSuccess: false, Message: "Error al procesar el formulario"})
				return
			}

			// Actualizar campos básicos
			product.Name = c.PostForm("name")
			product.Price, _ = strconv.ParseFloat(c.PostForm("price"), 64)
			product.Description = c.PostForm("description")
			product.CategoryName = c.PostForm("category_name")
			product.Stock, _ = strconv.Atoi(c.PostForm("stock"))

			// Manejo de imágenes - CORRECCIÓN AQUÍ
			file, err := c.FormFile("image")
			if err == nil {
				// Eliminar imagen antigua solo si existe una local
				if product.ImageLocalPath != "" {
					if err := utils.DeleteImage(product.ImageLocalPath); err != nil {
						c.JSON(http.StatusInternalServerError, dtos.ResponseDto{
							IsSuccess: false,
							Message:   "Error al eliminar la imagen anterior",
						})
						return
					}
				}

				// Guardar nueva imagen
				filename := product.Name + "_" + file.Filename
				path, err := utils.SaveImage(file, filename)
				if err != nil {
					c.JSON(http.StatusInternalServerError, dtos.ResponseDto{
						IsSuccess: false,
						Message:   "Error al guardar la imagen",
					})
					return
				}

				// Actualizar solo los campos de imagen local
				product.ImageLocalPath = path
				// Mantener el image_url original si existe
				if product.ImageUrl == "" {
					product.ImageUrl = "http://localhost:2222/uploads/images/" + filename
				}
			} else {
				// Si no se subió imagen, mantener las existentes
				imageUrl := c.PostForm("image_url")
				if imageUrl != "" {
					product.ImageUrl = imageUrl
				}
			}
		} else {
			c.JSON(http.StatusUnsupportedMediaType, dtos.ResponseDto{IsSuccess: false, Message: "Formato de contenido no soportado"})
			return
		}

		if err := db.Save(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, dtos.ResponseDto{
				IsSuccess: false,
				Message:   "Error al guardar el producto",
			})
			return
		}

		c.JSON(http.StatusOK, dtos.ResponseDto{IsSuccess: true, Result: product})
	}
}
