package messageHandler

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"jiqiren/bot/messageHandler/plain"
	"jiqiren/bot/messageHandler/reply"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/SevereCloud/vksdk/v2/events"
)

const battleReportStart = "Результаты битвы за "
const battleReportEnd = "Или встань на 🔐 защиту своей компании."
const teamProjectResultPrefix = "Ты вложился в запил командного проекта."
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
		replyParams := reply.MakeParams(parentMessage, vk)

		if (profileHandleResult.IsInserted) {
			message := "Профиль принят."
			for _, messageNote := range profileHandleResult.Messages {
				message += "\n" + messageNote
			}
			replyParams.Reply(message)
			return
		}

		replyParams.Reply("Произошла ошибка, профиль не принят. Повторите ещё раз позже.")
		return
	}

	if (isBattleReport(message.Text)) {
		if (os.Getenv("ENV") == "dev") {
			log.Println("Battle report from " + strconv.Itoa(senderId))
			log.Println("Contents: " + message.Text)
		}

		reportResult := plain.HandleReport(message.Text, senderId, messageDate)

		replyParams := reply.MakeParams(parentMessage, vk)
		if (!reportResult.IsUserExist) {
			replyParams.Reply("Отчёт не принят, сначала нужно отправить боту профиль.")
			return
		}

		if (!reportResult.IsFirst) {
			replyParams.Reply("Отчёт не принят, повторный отчёт за эту дату.")
			return
		}

		if (!reportResult.IsStored) {
			replyParams.Reply("Отчёт не принят. Неизвестная ошибка.")
			return
		}

		if (!reportResult.IsParticipated) {
			replyParams.Reply("Отчёт принят. Участие в битвах важно для фракции, не пропускайте битвы!")
			return
		}

		if (!reportResult.IsProfileFound) {
			replyParams.Reply("Отчёт принят, но не найден профиль за период между этой битвой и прошлой. Не забывайте сдавать профиль.")
			return
		}

		message := "Отчёт принят. Спасибо за участие в битве!"
		if (reportResult.ContestResult != nil) {
			message = message + "\n" + reportResult.ContestResult.Message
		}
		replyParams.Reply(message)

		return
	}

	if (isTeamProjectResult(message.Text)) {
		if (os.Getenv("ENV") == "dev") {
			log.Println("team project result from " + strconv.Itoa(senderId))
			log.Println("Contents: " + message.Text)
		}

		projectResult := plain.HandleTeamProject(message, senderId, messageDate)

		if (projectResult) {
			replyParams := reply.MakeParams(parentMessage, vk)
			replyParams.Reply("Вклад в проект принят")
		}
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

func isTeamProjectResult(messageText string) bool {
	return strings.HasPrefix(messageText, teamProjectResultPrefix)
}