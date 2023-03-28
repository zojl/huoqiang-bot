package messageHandler

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"huoqiang/bot/messageHandler/plain"

	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/SevereCloud/vksdk/v2/events"
)

func HandlePlain(messageObject events.MessageNewObject) {
	if (len(messageObject.Message.FwdMessages) > 0) {
		for _, message := range messageObject.Message.FwdMessages {
			if (strconv.Itoa(message.FromID) == os.Getenv("HW_ID")) {
				//fmt.Printf("%+v\n", messageObject)
				handleHwForward(message, messageObject.Message.FromID)
			}
		}
	}
}

func handleHwForward(message object.MessagesMessage, senderId int) {
	if (isProfileMessage(message.Text)) {
		fmt.Println("That's a profile message from " + strconv.Itoa(senderId))
		plain.HandleProfile(message.Text, senderId)
		return
	}

	fmt.Println("fuck!")
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

