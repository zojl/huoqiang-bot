package repository

import (
	"jiqiren/bot/database"
	"jiqiren/bot/database/model"
)


func FindOneContestProjectMessageByMessageId(messageId int) (*model.ContestProjectMessages) {
	db := database.GetDb()
	cpm := model.ContestProjectMessages{}

	err := db.Model(&cpm).
		Where(
			"message_id = ?",
			messageId,
		).
		First(&cpm).Error

	if (err != nil) {
		return nil
	}

	return &cpm
}

func CountContestProjectMessagesByVkIdAndContestId(contest *model.Contest, user *model.User) int64 {
	db := database.GetDb()
	var count int64
	err := db.Model(&model.ContestProjectMessages{}).
		Where("user_id = ? AND contest_id = ?", user.Id, contest.Id).Count(&count).Error

	if (err != nil) {
		return 0
	}

	return count
}