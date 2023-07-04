package plain

import (
	"time"
	"github.com/SevereCloud/vksdk/v2/object"
	

	"jiqiren/bot/database"
	"jiqiren/bot/database/model"
	"jiqiren/bot/database/repository"
)

func HandleTeamProject(message object.MessagesMessage, senderId int, messageTime time.Time) bool {
	messageId := message.ConversationMessageID
	
	user := repository.FindOneUserByVkId(senderId)
	if (user == nil) {
		return false
	}
	
	lastProfile, profileErr := repository.FindLastProfileByVkId(uint(senderId))
	if (profileErr != nil) {
		return false
	}
	
	contest, contestErr := repository.FindOneContestByFractionIdTypeCodeAndDate(lastProfile.FractionId, "project", &messageTime)
	if (contestErr != nil) {
		return false
	}
	
	cpm := repository.FindOneContestProjectMessageByMessageId(messageId)
	
	if (cpm != nil) {
		return false
	}
	
	newCpm := model.ContestProjectMessages{
		UserId: user.Id,
		ContestId: contest.Id,
		MessageId: messageId,
		MessageDate: &messageTime,
	}
	
	db := database.GetDb()
	db.Create(&newCpm)
	
	return true
}