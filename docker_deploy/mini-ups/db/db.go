package db

import (
	"log"
	"mini-ups/config"
	"mini-ups/model"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dbconfig := config.AppConfig.Postgres
	dsn := "host=" + dbconfig.Host + " user=" + dbconfig.User + " password=" + dbconfig.Password + " dbname=" + dbconfig.DBName + " port=" + strconv.Itoa(dbconfig.Port)
	log.Printf("host:%s", dbconfig.Host)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	DB = db

	initTables()
}

func initTables() {
	err := DB.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal("failed to auto migrate user:", err)
	}

	err = DB.AutoMigrate(&model.Truck{})
	if err != nil {
		log.Fatal("failed to auto migrate truck:", err)
	}

	err = DB.AutoMigrate(&model.Package{})
	if err != nil {
		log.Fatal("failed to auto migrate truck:", err)
	}
}
