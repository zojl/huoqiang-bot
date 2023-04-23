package repository

import (
	"fmt"
	"time"

	"jiqiren/bot/database"
	"jiqiren/bot/database/model"
)

func FindLastProfileByVkIdForBattleDate(vkId uint, battleDate *time.Time) (*model.Profile, error) {
	db := database.GetDb()
	profile := model.Profile{}

	battleHourDiff := 15 //battle is at 18:00 msk, means 15:00 utc
	battleDuration := time.Duration(battleHourDiff) * time.Hour
	battleMoment := battleDate.Add(+battleDuration)

	previousBattleHourDiff := 9 //difference between previous battle and battle day midnight: 15+9 = 24
	previousBattleDuration := time.Duration(previousBattleHourDiff) * time.Hour
	previousBattleMoment := battleDate.Add(-previousBattleDuration)

	err := db.Model(&profile).Joins("JOIN users ON users.id = profiles.user_id").
		Where("users.vk_id = ? AND profiles.message_date < ? AND profiles.message_date > ?", vkId, battleMoment, previousBattleMoment).Last(&profile).Error

	if (err != nil) {
		return nil, fmt.Errorf("No profiles associated with vk id: %d for battle at %+v", vkId, battleDate)
	}

	return &profile, nil
}