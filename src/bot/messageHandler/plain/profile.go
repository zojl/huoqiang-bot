package profile

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"huoqiang/bot/database"
	"huoqiang/bot/database/model"
)

var regulars = []string {
	`(?P<Username>[a-zA-Z0-9а-яА-ЯёЁ ,_\-!?\$<>]{3,25})` +
		`\s\(` +
		`(?P<Lead>\S?)(?P<Squad>[A-Z]{2})\s\|\s\S+` +
		`(?P<Fraction>Huǒqiáng|Aegis|V-hack|Phantoms|NetKings|NHS)` +
		`\)`,
	`💻: (?P<Level>\d+)`,
	`💻Уровень: (?P<Level>\d+)`,
	`💡: (?P<Experience>\d+)`,
	`💡Опыт: (?P<Experience>\d+)`,
	`💵: (?P<Money>\d+)`,
	`💵Деньги: (?P<Money>\d+)`,
	`💸: (?P<Vkcoin>\d+\.?\d*)`,
	`💸VKCoin: (?P<Vkcoin>\d+\.?\d*)`,
	`🔘: (?P<Points>\d+)`,
	`🔘Поинты: (?P<Points>\d+)`,
	`💳: (?P<Bitcoins>\d+)`,
	`💳Биткоины: (?P<Bitcoins>\d+)`,
	`💿: (?P<Disks>\d+)`,
	`💿Радиодетали: (?P<Disks>\d+)`,
	`📄: (?P<Pages>\d+)`,
	`📄Страницы: (?P<Pages>\d+)`,
	`💽: (?P<Chips>\d+)`,
	`💽Микрочипы: (?P<Chips>\d+)`,
	`📑: (?P<Instructions>\d+)`,
	`📑Инструкции: (?P<Instructions>\d+)`,
	`📈: (?P<Stocks>\d+)`,
	`📈Акции: (?P<Stocks>\d+)`,
	`🔥: (?P<Motivation>\d+)/(?P<MotivationLimit>\d+)`,
	`🔥Мотивация: (?P<Motivation>\d+) из (?P<MotivationLimit>\d+)`,
	`📡: (?P<Practice>\d+)`,
	`📡Практика: (?P<Practice>\d+)`,
	`📡Практика: [\d]+\+[\d]+ \((?P<Practice>\d+)\)`,
	`💾: (?P<Theory>\d+)`,
	`💾Теория: (?P<Theory>\d+)`,
	`💾Теория: [\d]+\+[\d]+ \((?P<Theory>\d+)\)`,
	`📱: (?P<Cunning>\d+)`,
	`📱Хитрость: (?P<Cunning>\d+)`,
	`📱Хитрость: [\d]+\+[\d]+ \((?P<Cunning>\d+)\)`,
	`🔎: (?P<Wisdom>\d+)`,
	`🔎Мудрость: (?P<Wisdom>\d+)`,
	`🔎Мудрость: [\d]+\+[\d]+ \((?P<Wisdom>\d+)\)`,
	`🔋: (?P<Stamina>\d+)`,
	`🔋Выносливость: (?P<Stamina>\d+)`,
	`\n⚔: (?P<Target>[🔐💠🚧🎭🈵🔱🇺🇸]{1,2})`,
	`Занятие:\n[А-Яа-я ⚔]*(?P<Target>[🔐💠🚧🎭🈵🔱🇺🇸]{1,2})`,
}

var compiledRegulars = make([]*regexp.Regexp, len(regulars), len(regulars))

var areRegularsCompiled = false

func HandleProfile(messageText string, senderId int, messageTime time.Time) {
	parsedProfile := parseProfile(messageText)
	storeParsedProfile(parsedProfile, senderId, messageTime)
}

func storeParsedProfile(parsedProfile *ProfileParseResult, senderId int, messageTime time.Time) {
	user := getUser(parsedProfile, senderId)
	fraction := getFraction(parsedProfile)
	team := getTeam(parsedProfile, fraction)
	target := getTarget(parsedProfile, fraction)
	insertProfile(parsedProfile, user, fraction, team, target, messageTime)
}

func getUser(parsedProfile *ProfileParseResult, senderId int) *model.User {
	db := database.GetDb()
	user := model.User{}
	err := db.Unscoped().Where("vk_id = ?", senderId).First(&user).Error
	if (err != nil) {
		user.VkId = uint(senderId)
		db.Create(&user)
	}

	return &user
}

func getFraction(parsedProfile *ProfileParseResult) *model.Fraction {
	db := database.GetDb()
	fraction := model.Fraction{}
	err := db.Unscoped().Where("Name = ?", parsedProfile.Fraction).First(&fraction).Error
	if (err != nil) {
		log.Fatal("Erroneous fraction name: " + parsedProfile.Fraction)
	}

	return &fraction
}

func getTeam(parsedProfile *ProfileParseResult, fraction *model.Fraction) *model.Team {
	db := database.GetDb()
	team := model.Team{}
	if (len(parsedProfile.Squad) == 2) {
		err := db.Unscoped().Where("Code = ?", parsedProfile.Squad).First(&team).Error

		if (err != nil) {
			team.Code = parsedProfile.Squad
			team.Fraction = *fraction
			db.Create(&team)
		}
	}

	return &team
}

func getTarget(parsedProfile *ProfileParseResult, userFraction *model.Fraction) *model.Fraction {
	if (len(parsedProfile.Target) == 0) {
		var fraction model.Fraction
		return &fraction
	}

	if (parsedProfile.Target == "🔐") {
		return userFraction
	}

	db := database.GetDb()
	fraction := model.Fraction{}
	err := db.Unscoped().Where("Icon = ?", parsedProfile.Target).First(&fraction).Error
	if (err != nil) {
		log.Fatal("Erroneous fraction icon: " + parsedProfile.Target)
	}

	return &fraction
}

func insertProfile(
	parsedProfile *ProfileParseResult,
	user *model.User,
	fraction *model.Fraction,
	team *model.Team,
	target *model.Fraction,
	messageTime time.Time,
) {
	var (
		profile model.Profile
		val uint64
	)

	val, _ = strconv.ParseUint(parsedProfile.Level, 10, 64)
	profile.Level = uint(val)
	val, _ = strconv.ParseUint(parsedProfile.Experience, 10, 64)
	profile.Experience = uint(val)
	val, _ = strconv.ParseUint(parsedProfile.Money, 10, 64)
	profile.Money = uint(val)
	val, _ = strconv.ParseUint(parsedProfile.Points, 10, 64)
	profile.Points = uint(val)
	val, _ = strconv.ParseUint(parsedProfile.Bitcoins, 10, 64)
	profile.Bitcoins = uint(val)
	val, _ = strconv.ParseUint(parsedProfile.Disks, 10, 64)
	profile.Disks = uint(val)
	val, _ = strconv.ParseUint(parsedProfile.Pages, 10, 64)
	profile.Pages = uint(val)
	val, _ = strconv.ParseUint(parsedProfile.Chips, 10, 64)
	profile.Chips = uint(val)
	val, _ = strconv.ParseUint(parsedProfile.Instructions, 10, 64)
	profile.Instructions = uint(val)
	val, _ = strconv.ParseUint(parsedProfile.Stocks, 10, 64)
	profile.Stocks = uint(val)
	val, _ = strconv.ParseUint(parsedProfile.Motivation, 10, 64)
	profile.Motivation = uint(val)
	val, _ = strconv.ParseUint(parsedProfile.MotivationLimit, 10, 64)
	profile.MotivationLimit = uint(val)
	val, _ = strconv.ParseUint(parsedProfile.Practice, 10, 64)
	profile.Practice = uint(val)
	val, _ = strconv.ParseUint(parsedProfile.Theory, 10, 64)
	profile.Theory = uint(val)
	val, _ = strconv.ParseUint(parsedProfile.Cunning, 10, 64)
	profile.Cunning = uint(val)
	val, _ = strconv.ParseUint(parsedProfile.Wisdom, 10, 64)
	profile.Wisdom = uint(val)
	val, _ = strconv.ParseUint(parsedProfile.Stamina, 10, 64)
	profile.Stamina = uint(val)

	profile.Team = *team
	profile.Fraction = *fraction
	profile.User = *user
	profile.MessageDate = messageTime

	profile.Username = parsedProfile.Username
	profile.Vkcoin, _ = strconv.ParseFloat(parsedProfile.Vkcoin, 10)
	profile.Lead = getLeadType(parsedProfile.Lead)

	if (target.Id > 0) {
		profile.Target = *target
	} else {
		profile.TargetId = nil
	}

	db := database.GetDb()
	db.Create(&profile)
}

func getLeadType(icon string) uint {
	if (icon == "🔸") {
		return 2
	}

	if (icon == "▫") {
		return 1
	}

	return 0
}

func parseProfile(messageText string) *ProfileParseResult {
	var parsedProfile ProfileParseResult
	reflection := reflect.ValueOf(&parsedProfile)

	for _, regular := range getCompiledRegulars() {
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
			}

			value := match[key]
			if (len(value) == 0) {
				continue
			}

			field.SetString(value)
		}
	}

	return &parsedProfile
}

func getCompiledRegulars () []*regexp.Regexp {

	if (areRegularsCompiled) {
		return compiledRegulars
	}

	for i, regular := range regulars {
		compiledRegular, err := regexp.Compile(regular)
		if (err != nil) {
			fmt.Println(err)
		} else {
			compiledRegulars[i] = compiledRegular
		}
	}

	areRegularsCompiled = true
	return compiledRegulars
}

type ProfileParseResult struct {
	Username string
	Squad string
	Fraction string
	Lead string
	Level string
	Experience string
	Money string
	Vkcoin string
	Points string
	Bitcoins string
	Disks string
	Pages string
	Chips string
	Instructions string
	Stocks string
	Motivation string
	MotivationLimit string
	Practice string
	Theory string
	Cunning string
	Wisdom string
	Stamina string
	Target string
}
