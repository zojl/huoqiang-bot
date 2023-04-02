package plain

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"huoqiang/bot/messageHandler/plain"

	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/SevereCloud/vksdk/v2/events"
)

func Handle(messageObject events.MessageNewObject) {
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
	messageDate := time.Unix(int64(message.Date), 0)

	if (isProfileMessage(message.Text)) {
		fmt.Println("That's a profile message from " + strconv.Itoa(senderId))
		profile.HandleProfile(message.Text, senderId, messageDate)
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

