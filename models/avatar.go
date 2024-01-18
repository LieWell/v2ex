package models

import (
	"liewell.fun/v2ex/core"
)

var (
	avatarTableName = "avatar"
	EmptyAvatar     = &Avatar{}
)

type Avatar struct {
	Id     int    `gorm:"primaryKey;autoIncrement:true;column:id" json:"id"`
	Name   string `gorm:"column:name" json:"name"`
	Avatar string `gorm:"column:avatar;type:varchar(2048)" json:"avatar"`
}

func (a *Avatar) TableName() string {
	return avatarTableName
}

func SaveAvatar(a *Avatar) (int, error) {
	err := core.MYSQL.Save(a).Error
	return a.Id, err
}

func FindAllAvatar() ([]Avatar, error) {
	var avatars []Avatar
	err := core.MYSQL.Find(&avatars).Error
	return avatars, err
}

func DeleteAvatar(id int) error {
	return core.MYSQL.Delete(EmptyAvatar, id).Error
}
