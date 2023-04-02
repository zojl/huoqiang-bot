package database

import (
	"fmt"
	"log"
	"os"

	"huoqiang/bot/database/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() {
	dsn := fmt.Sprintf(
		"host=db user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Europe/Moscow",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if (err != nil) {
		log.Fatal(err)
	}

	for _, entityName := range model.GetAllModels() {
		db.AutoMigrate(entityName)
	}

	model.CreateInitialValues(db)
}

func GetDb() *gorm.DB {
	fmt.Println(db)
	return db
}