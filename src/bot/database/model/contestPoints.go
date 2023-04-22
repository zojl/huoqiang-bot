package model

import "time"

type ContestPoints struct {
	Id uint `gorm:"primaryKey"`
	UserId uint
	User *User `gorm:"foreignKey:UserId"`
	ContestId uint
	Contest *Contest `gorm:"foreignKey:ContestId"`
	Points int
	BattleDate *time.Time
}