package model

type Team struct {
	Id uint `gorm:"primaryKey"`
	Name string
	Code string `gorm:"size:2`
	FractionId uint
	Fraction Fraction `gorm:"foreignKey:FractionId"`
}