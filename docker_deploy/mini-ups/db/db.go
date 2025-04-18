package db

import (
	"log"
	"mini-ups/config"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dbconfig := config.AppConfig.Postgres
	dsn := "host=" + dbconfig.Host + " user=" + dbconfig.User + " password=" + dbconfig.Password + " dbname=" + dbconfig.DBName + " port=" + strconv.Itoa(dbconfig.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	DB = db
}
