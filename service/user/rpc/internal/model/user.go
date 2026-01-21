package model

import "time"

type User struct {
	Id         uint64            `gorm:"primaryKey"`
	Username   string            `gorm:"column:username"`
	Password   string            `gorm:"column:password"`
	Email      string            `gorm:"column:email"`
	Status     int64             `gorm:"column:status;default:0"`
	Score      int32             `gorm:"column:score"`
	ExtraInfo  map[string]string `gorm:"column:extra_info;serializer:json"`
	CreateTime time.Time         `gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time         `gorm:"column:update_time;autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}
