package messageHandler

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"huoqiang/bot/messageHandler/plain"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/SevereCloud/vksdk/v2/events"
)

func HandlePlain(messageObject events.MessageNewObject, vk *api.VK) {
	if (len(messageObject.Message.FwdMessages) > 0) {
		for _, message := range messageObject.Message.FwdMessages {
			if (strconv.Itoa(message.FromID) == os.Getenv("HW_ID")) {
				handleHwForward(message, &messageObject.Message, vk)
			}
		}
	}
}

func handleHwForward(message object.MessagesMessage, parentMessage *object.MessagesMessage, vk *api.VK) {
	messageDate := time.Unix(int64(message.Date), 0)

	if (isProfileMessage(message.Text)) {
		senderId := parentMessage.FromID
		if (os.Getenv("ENV") == "dev") {
			fmt.Println("That's a profile message from " + strconv.Itoa(senderId))
			fmt.Println("Contents: " + message.Text)
		}
		if (plain.HandleProfile(message.Text, senderId, messageDate)) {
			ReplyTo(parentMessage, "Профиль принят", vk)
		}
		
		return
	}

	fmt.Println("bad message " + message.Text)
}

func isProfileMessage(messageText string) bool {
	icons := []string{"💻", "💡", "💵", "📈", "💿", "📄", "💽", "📑", "🔘", "💸", "🔥", "🔋", "📡", "💾", "📱", "🔎"}

	for _, icon := range icons {
		if (!strings.Contains(messageText, icon)) {
			return false;
		}
	}

	return true
}

