package model

type User struct {
	Id uint `gorm:"primaryKey"`
	VkId uint `gorm:"unique"`
}