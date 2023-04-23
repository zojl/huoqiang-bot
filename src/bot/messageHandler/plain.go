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

		profileHandleResult := plain.HandleProfile(message.Text, senderId, messageDate)
		if (profileHandleResult.IsInserted) {
			message := "Профиль принят."
			for _, messageNote := range profileHandleResult.Messages {
				message += "\n" + messageNote
			}
			ReplyTo(parentMessage, message, vk)
			return
		}

		ReplyTo(parentMessage, "Произошла ошибка, профиль не принят. Повторите ещё раз позже.", vk)
		return
	}

	if (isBattleReport(message.Text)) {
		if (os.Getenv("ENV") == "dev") {
			log.Println("Battle report from " + strconv.Itoa(senderId))
			log.Println("Contents: " + message.Text)
		}

		reportResult := plain.HandleReport(message.Text, senderId, messageDate)

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

		if (!reportResult.IsProfileFound) {
			ReplyTo(parentMessage, "Отчёт принят, но не найден профиль за период между этой битвой и прошлой. Не забывайте сдавать профиль.", vk)
			return
		}

		message := "Отчёт принят. Спасибо за участие в битве!"
		if (reportResult.ContestResult != nil) {
			message = message + "\n" + reportResult.ContestResult.Message
		}
		ReplyTo(parentMessage, message, vk)

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
