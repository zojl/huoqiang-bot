package repository

import (
	"fmt"
	"time"

	"huoqiang/bot/database"
	"huoqiang/bot/database/model"
)


func FindOneContestByFractionIdTypeCodeAndDate(fractionId uint, contestTypeCode string, date *time.Time) (*model.Contest, error) {
	db := database.GetDb()
	contest := model.Contest{}

	err := db.Model(&contest).
		Joins("JOIN contest_types ON contest_types.id = contests.type_id AND contest_types.code = ?", contestTypeCode).
		Where(
			"contests.fraction_id = ? AND contests.start_at < ? AND contests.end_at > ?",
			fractionId,
			date,
			date,
		).
		First(&contest).Error

	if (err != nil) {
		return nil, fmt.Errorf("No contest found by criteria: %d, %s, %+v", fractionId, contestTypeCode, date)
	}

	return &contest, nil
}