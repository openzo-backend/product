package main

import (
	"fmt"

	"github.com/tanush-128/openzo_backend/product/config"
	"github.com/tanush-128/openzo_backend/product/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func connectToDB(cfg *config.Config) (*gorm.DB, error) {

	if cfg.MODE == "prdoduction" {
		dsn := cfg.DB_URL

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to open database connection: %w", err)
		}
		return db, nil

	}
	db, err := gorm.Open(
		sqlite.Open("test.db"),

		&gorm.Config{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	db.Migrator().AutoMigrate(&models.Product{})
	db.Migrator().AutoMigrate(&models.ProductImage{})
	db.Migrator().AutoMigrate(&models.SizeVariant{})
	db.Migrator().AutoMigrate(&models.ColorVariant{})

	return db, nil
}
