package plain


import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"huoqiang/bot/database"
	"huoqiang/bot/database/model"
	"huoqiang/bot/database/repository"
)

var reportRegulars = [...]string {
	`Результаты битвы за (?P<Date>\d{2}\.\d{2}\.\d{4})`,
	`(?:Атака на |Команда )(?P<Target>[💠🚧🎭🈵🔱🇺🇸]{1,2})\S+ (?P<BattleResult>прошла успешно|отбила вашу атаку)`, //this ones for attack
	`Твоя команда (?P<BattleResult>успешно отбила атаку.|не смогла отбить атаку.)`, //this ones for defence
	`Ты получил:(?:\nДеньги: ?(?P<RewardMoney>\d+)💵)?\nОпыт: (?P<RewardExperience>\d+)💡(?:\nVKCoin: (?P<RewardVkc>\d+\.?\d+)💸)?\nОсталось выносливости: 🔋(?P<Stamina>\d+)%`,
	`(?:Вы перехватили транзакцию \S+|Вы получили награду за защиту:)\nДеньги: (?P<TransactionMoney>\d+)💵\nОпыт: (?P<TransactionExperience>\d+)💡`,
}

var compiledReportRegulars = make([]*regexp.Regexp, len(reportRegulars), len(reportRegulars))
var areReportRegularsCompiled = false

func HandleReport(messageText string, senderId int, messageTime time.Time) *ReportResponse {
	parsedReport := parseReport(messageText)
	response := ReportResponse{}

	user := repository.FindOneUserByVkId(senderId)
	if (user == nil) {
		response.IsUserExist = false
		response.IsStored = false
		return &response
	}

	response.IsUserExist = true

	date := parseDate(parsedReport.Date)

	reportsCount := repository.CountReportsByUserAndDate(senderId, date)
	if (reportsCount > 0) {
		response.IsFirst = false
		response.IsStored = false
		return &response
	}

	response.IsFirst = true

	target := getReportTarget(parsedReport)
	if (len(parsedReport.BattleResult) > 0) {
		response.IsParticipated = true
	}

	response.IsStored = storeReport(parsedReport, user, date, target, messageTime)

	return &response
}

func parseReport(messageText string) *ReportParseResult {
	var parsedReport ReportParseResult
	reflection := reflect.ValueOf(&parsedReport)

	for _, regular := range getCompiledReportRegulars() {
		allMatches := regular.FindAllStringSubmatch(messageText, -1)
		if (len(allMatches) == 0) {
			continue
		}
		match := allMatches[len(allMatches) - 1][1:]
		groups := regular.SubexpNames()[1:]

		for key, _ := range groups {
			field := reflection.Elem().FieldByName(groups[key])
			if !field.IsValid() {
				fmt.Println("invalid field " + groups[key])
				continue
			}

			value := match[key]
			if (len(value) == 0) {
				continue
			}

			field.SetString(value)
		}
	}

	return &parsedReport
}

func storeReport(parsedReport *ReportParseResult, user *model.User, date *time.Time, target *model.Fraction, messageTime time.Time) bool {
	var val uint64
	report := model.Report{}
	report.Target = target
	report.User = user
	report.BattleDate = date
	report.MessageDate = messageTime

	if len(parsedReport.RewardExperience) == 0 {
		report.IsSkipped = true
	}

	report.IsAttack = (target != nil)
	if (report.IsAttack) {
		report.IsSuccess = (parsedReport.BattleResult == "прошла успешно")
	} else {
		report.IsSuccess = (parsedReport.BattleResult == "успешно отбила атаку.")
	}

	val, _ = strconv.ParseUint(parsedReport.RewardMoney, 10, 64)
	report.RewardMoney = uint(val)
	val, _ = strconv.ParseUint(parsedReport.RewardExperience, 10, 64)
	report.RewardExperience = uint(val)
	val, _ = strconv.ParseUint(parsedReport.RewardVkc, 10, 64)
	report.Stamina = uint(val)
	val, _ = strconv.ParseUint(parsedReport.TransactionMoney, 10, 64)
	report.TransactionMoney = uint(val)
	val, _ = strconv.ParseUint(parsedReport.TransactionExperience, 10, 64)
	report.TransactionExperience = uint(val)
	report.RewardVkc, _ = strconv.ParseFloat(parsedReport.Stamina, 10)


	db := database.GetDb()
	db.Create(&report)

	return report.Id != 0
}

func parseDate(date string) *time.Time {
	layout := "02.01.2006"
	result, err := time.Parse(layout, date)
	if err != nil {
		return nil
	}
	return &result
}

func getReportTarget(report *ReportParseResult) *model.Fraction {
	if (len(report.Target) == 0) {
		return nil
	}

	target, _ := repository.FindOneFractionByIcon(report.Target)
	return target
}

func getCompiledReportRegulars () []*regexp.Regexp {

	if (areReportRegularsCompiled) {
		return compiledReportRegulars
	}

	for i, regular := range reportRegulars {
		compiledRegular, err := regexp.Compile(regular)
		if (err != nil) {
			fmt.Println(err)
		} else {
			compiledReportRegulars[i] = compiledRegular
		}
	}

	areReportRegularsCompiled = true
	return compiledReportRegulars
}

type ReportParseResult struct {
	Date string
	BattleResult string
	Target string
	RewardMoney string
	RewardExperience string
	RewardVkc string
	Stamina string
	TransactionMoney string
	TransactionExperience string
}

type ReportResponse struct {
	//IsToday bool
	IsFirst bool
	IsParticipated bool
	IsStored bool
	IsUserExist bool
}