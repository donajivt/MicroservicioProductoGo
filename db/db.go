package db

import (
	"go-product-service/config"
	"go-product-service/models"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	var err error
	DB, err = gorm.Open(sqlserver.Open(config.Cfg.DSN), &gorm.Config{})
	if err != nil {
		return err
	}

	return DB.AutoMigrate(&models.Product{})
}
