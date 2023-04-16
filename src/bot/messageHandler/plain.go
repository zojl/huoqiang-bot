package messageHandler

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"huoqiang/bot/messageHandler/plain"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/SevereCloud/vksdk/v2/events"
)

const battleReportStart = "Результаты битвы за "
const battleReportEnd = "Или встань на 🔐 защиту своей компании."
var profileIcons = [...]string {"💻", "💡", "💵", "📈", "💿", "📄", "💽", "📑", "🔘", "💸", "🔥", "🔋", "📡", "💾", "📱", "🔎"}

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
	senderId := parentMessage.FromID

	if (isProfileMessage(message.Text)) {
		if (os.Getenv("ENV") == "dev") {
			log.Println("That's a profile message from " + strconv.Itoa(senderId))
			log.Println("Contents: " + message.Text)
		}
		if (plain.HandleProfile(message.Text, senderId, messageDate)) {
			ReplyTo(parentMessage, "Профиль принят.", vk)
		}

		return
	}

	if (isBattleReport(message.Text)) {
		if (os.Getenv("ENV") == "dev") {
			log.Println("Battle report from " + strconv.Itoa(senderId))
			log.Println("Contents: " + message.Text)
		}

		reportResult := plain.HandleReport(message.Text, senderId, messageDate)
		log.Printf("%+v\n", reportResult)

		if (!reportResult.IsUserExist) {
			ReplyTo(parentMessage, "Отчёт не принят, сначала нужно отправить боту профиль.", vk)
			return
		}

		if (!reportResult.IsFirst) {
			ReplyTo(parentMessage, "Отчёт не принят, повторный отчёт за эту дату.", vk)
			return
		}

		if (!reportResult.IsStored) {
			ReplyTo(parentMessage, "Отчёт не принят. Неизвестная ошибка.", vk)
			return
		}

		if (!reportResult.IsParticipated) {
			ReplyTo(parentMessage, "Отчёт принят. Участие в битвах важно для фракции, не пропускайте битвы!", vk)
			return
		}

		ReplyTo(parentMessage, "Отчёт принят. Спасибо за участие в битве!", vk)

		return
	}

	log.Println("bad message " + message.Text)
}

func isProfileMessage(messageText string) bool {
	for _, icon := range profileIcons {
		if (!strings.Contains(messageText, icon)) {
			return false;
		}
	}

	return true
}

func isBattleReport(messageText string) bool {
	if !strings.HasPrefix(messageText, battleReportStart) {
		return false
	}

	return strings.HasSuffix(messageText, battleReportEnd)
}
