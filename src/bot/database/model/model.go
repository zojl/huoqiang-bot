package model

import (
	"gorm.io/gorm"
)

func GetAllModels() []interface{} {
	return []interface{} {
		Fraction{},
		User{},
		Team{},
		Profile{},
	}
}

func CreateInitialValues(db *gorm.DB) {
	fillFractions(db)
}

func fillFractions(db *gorm.DB) {
	var existingFractionsCount int64
	db.Model(&Fraction{}).Count(&existingFractionsCount)
	if (existingFractionsCount >= 6) {
		return
	}

	fractions := []Fraction {
		{Name: "Aegis", Code: "ae", Icon: "💠"},
		{Name: "V-hack", Code: "vh", Icon: "🚧"},
		{Name: "Phantoms", Code: "ph", Icon: "🎭"},
		{Name: "Huǒqiáng", Code: "hu", Icon: "🈵"},
		{Name: "NetKings", Code: "nk", Icon: "🔱"},
		{Name: "NHS", Code: "nh", Icon: "🇺🇸"},
	}

	tableName := db.Model(&Fraction{}).Name()
	db.Exec("TRUNCATE TABLE " + tableName + " RESTART IDENTITY;")

	db.Create(&fractions)
}