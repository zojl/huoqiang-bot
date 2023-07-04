package model

import "time"

type ContestProjectMessages struct {
	Id uint `gorm:"primaryKey"`
	UserId uint
	User *User `gorm:"foreignKey:UserId"`
	ContestId uint
	Contest *Contest `gorm:"foreignKey:ContestId"`
	MessageId int
	MessageDate *time.Time
}