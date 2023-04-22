package model

type ContestType struct {
	Id uint `gorm:"primaryKey"`
	Name string
	Code string `gorm:"size:16`
}