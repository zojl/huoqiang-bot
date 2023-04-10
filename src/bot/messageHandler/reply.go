package messageHandler

import (
	"encoding/json"
	"log"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/SevereCloud/vksdk/v2/api/params"
)

func ReplyTo(sourceMessage *object.MessagesMessage, text string, vk *api.VK) {
	messageBuilder := params.NewMessagesSendBuilder()
	messageBuilder.Message(text)
	messageBuilder.RandomID(0)
	messageBuilder.PeerID(sourceMessage.PeerID)

	forward := Forward{}
	forward.IsReply = 1
	forward.PeerID = sourceMessage.PeerID
	messageIds := []int{sourceMessage.ConversationMessageID}
	forward.ConversationMessageIDs = messageIds
	jsonMessageRaw, _ := json.Marshal(forward)
	jsonMessage := string(jsonMessageRaw)
	messageBuilder.Forward(jsonMessage)

	_, err := vk.MessagesSend(messageBuilder.Params)
	if err != nil {
		log.Fatal(err)
	}
}

type Forward struct {
	IsReply int `json:"is_reply"`
	PeerID int `json:"peer_id"`
	ConversationMessageIDs []int `json:"conversation_message_ids"`
}
