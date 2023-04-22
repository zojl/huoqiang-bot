package model

import "time"

type Contest struct {
	Id uint `gorm:"primaryKey"`
	Name string
	TypeId uint
	Type ContestType `gorm:"foreignKey:TypeId"`
	FractionId *uint `gorm:"null"`
	Fraction Fraction `gorm:"foreignKey:FractionId;null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	StartAt time.Time
	EndAt time.Time
}