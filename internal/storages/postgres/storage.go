package postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"test-work/internal/config"
)

type Storage struct {
	Db *gorm.DB
}

func New(storageCfg config.PostgresConfig) (*Storage, error) {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		storageCfg.Host, storageCfg.User, storageCfg.Password, storageCfg.DbName, storageCfg.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("connection error: %w", err)
	}

	return &Storage{Db: db}, nil
}
