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
			file, _ := c.FormFile("image")
			if file != nil {
				filename := dto.Name + "_" + file.Filename
				path, err := utils.SaveImage(file, filename)
				if err == nil {
					dto.ImageLocalPath = path
					dto.ImageUrl = "/uploads/product_images/" + filename
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
		} else if strings.HasPrefix(contentType, "multipart/form-data") {
			product.Name = c.PostForm("name")
			product.Price, _ = strconv.ParseFloat(c.PostForm("price"), 64)
			product.Description = c.PostForm("description")
			product.CategoryName = c.PostForm("category_name")

			file, _ := c.FormFile("image")
			if file != nil {
				// Eliminar imagen antigua
				if product.ImageLocalPath != "" {
					_ = utils.DeleteImage(product.ImageLocalPath)
				}

				filename := product.Name + "_" + file.Filename
				path, err := utils.SaveImage(file, filename)
				if err == nil {
					product.ImageLocalPath = path
					product.ImageUrl = "/uploads/images/" + filename
				}
			}
		} else {
			c.JSON(http.StatusUnsupportedMediaType, dtos.ResponseDto{IsSuccess: false, Message: "Formato de contenido no soportado"})
			return
		}

		db.Save(&product)
		c.JSON(http.StatusOK, dtos.ResponseDto{IsSuccess: true, Result: product})
	}
}
