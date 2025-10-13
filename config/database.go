package config

import (
	"context"
	"fmt"
	"log"
	"practice/env"
	"practice/models"
	"practice/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	ctx context.Context
	db  *gorm.DB
}

func NewDB(ctx context.Context, logger *logger.Logger) *DB {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
		env.DBHost, env.DBPort, env.DBUser, env.DBPassword, env.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Todo{},
	)

	if err != nil {
		logger.Fatal("failed to migrate database: %v", err)
	}

	logger.Info("✅ Database connected! Host: %s Port: %d DB: %s", env.DBHost, env.DBPort, env.DBName)

	return &DB{ctx: ctx, db: db}
}

func (d *DB) Instance() *gorm.DB {
	return d.db
}

func (d *DB) Close() {
	sqlDB, err := d.db.DB()
	if err != nil {
		log.Fatalf("failed to get db instance: %v", err)
	}

	if err := sqlDB.Close(); err != nil {
		log.Printf("⚠️ error closing DB: %v", err)
	} else {
		log.Println("🛑 Database connection closed.")
	}
}
