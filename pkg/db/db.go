package db

import (
	"fmt"
	"github.com/hovhannesyan/RiskIndex-AuthSVC/pkg/models"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func Init(cfg Config) Handler {
	url := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.Username, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		logrus.Fatalln(err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		logrus.Fatalln(err)
	}

	return Handler{db}
}
