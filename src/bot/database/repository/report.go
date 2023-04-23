package repository

import (
	"time"

	"jiqiren/bot/database"
	"jiqiren/bot/database/model"
)

func CountReportsByUserAndDate(vkId int, date *time.Time) int64 {
	db := database.GetDb()
	var count int64
	err := db.Model(&model.Report{}).Joins("JOIN users ON users.id = reports.user_id").
		Where("users.vk_id = ? AND reports.battle_date = ?", vkId, date).Count(&count).Error

	if (err != nil) {
		return 0
	}

	return count
}