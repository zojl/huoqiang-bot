package profile

import (
	"fmt"
	"reflect"
	"regexp"
)

var regulars = []string {
	`(?P<Username>[a-zA-Z0-9а-яА-ЯёЁ\s,_\-!?\$<>]{3,25})\s\(\S*(?P<Squad>[A-Z]{2})\s\|\s\S+(?P<Fraction>Huǒqiáng|Aegis|V-hack|Phantoms|NetKings|NHS)\)`,
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
}

var compiledRegulars = make([]*regexp.Regexp, len(regulars), len(regulars))

var areRegularsCompiled = false

func HandleProfile(messageText string, senderId int) {
	parsedProfile := parseProfile(messageText)
	fmt.Printf("%+v\n", parsedProfile)
}

func parseProfile(messageText string) ProfileParseResult {
	var result ProfileParseResult
	reflection := reflect.ValueOf(&result)

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
			field.SetString(value)
		}
	}

	return result
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
}