package model

import "time"

type Profile struct {
	Id uint `gorm:"primaryKey"`
	Username string
	TeamId uint
	Team Team `gorm:"foreignKey:TeamId"`
	FractionId uint
	Fraction Fraction `gorm:"foreignKey:FractionId"`
	Lead uint
	Level uint
	Experience uint
	Money uint
	Vkcoin float64 `gorm:"type:decimal(16,6);"`
	Points uint
	Bitcoins uint
	Disks uint
	Pages uint
	Chips uint
	Instructions uint
	Stocks uint
	Motivation uint
	MotivationLimit uint
	Practice uint
	Theory uint
	Cunning uint
	Wisdom uint
	Stamina uint
	UserId uint
	User User `gorm:"foreignKey:UserId"`
	MessageDate time.Time
	Created time.Time `gorm:"autoCreateTime"`
}