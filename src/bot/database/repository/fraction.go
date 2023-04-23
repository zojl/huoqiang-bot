package repository

import (
	"fmt"

	"jiqiren/bot/database"
	"jiqiren/bot/database/model"
)

func FindOneFractionByIcon(icon string) (*model.Fraction, error) {
	if (icon == "🔐") {
		return nil, nil
	}

	db := database.GetDb()
	fraction := model.Fraction{}
	err := db.Unscoped().Where("Icon = ?", icon).First(&fraction).Error
	if (err != nil) {
		return nil, fmt.Errorf("Erroneous fraction icon: %s", icon)
	}

	return &fraction, nil
}