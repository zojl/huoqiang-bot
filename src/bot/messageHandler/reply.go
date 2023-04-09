package messageHandler

import (
	"log"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
)

func ReplyTo(peerId int, messageId int, text string, vk *api.VK) {
	messageBuilder := params.NewMessagesSendBuilder()
	messageBuilder.Message(text)
	messageBuilder.RandomID(0)
	messageBuilder.PeerID(peerId)
	messageBuilder.ReplyTo(messageId)

	_, err := vk.MessagesSend(messageBuilder.Params)
	if err != nil {
		log.Fatal(err)
	}
}