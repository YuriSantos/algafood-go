package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(cfg *DatabaseConfig) (*gorm.DB, error) {
	// Create database if not exists
	if cfg.CreateDatabaseIfNotExist {
		if err := createDatabaseIfNotExists(cfg); err != nil {
			log.Printf("Warning: could not create database: %v", err)
		}
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func createDatabaseIfNotExists(cfg *DatabaseConfig) error {
	// Connect without specifying database name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=%s&parseTime=%t&loc=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Charset, cfg.ParseTime, cfg.Loc)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("falha ao conectar ao servidor MySQL: %w", err)
	}
	defer db.Close()

	// Create database if not exists
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.Name))
	if err != nil {
		return fmt.Errorf("falha ao criar banco de dados: %w", err)
	}

	log.Printf("Banco de dados '%s' está pronto", cfg.Name)
	return nil
}
