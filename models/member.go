package models

import (
	"errors"
	"gorm.io/gorm"
	"liewell.fun/v2ex/core"
	"time"
)

var (
	memberTableName = "member"
	EmptyMember     = &Member{}
)

type Member struct {
	Id         int       `gorm:"primaryKey;autoIncrement:true;column:id" json:"id"`
	Number     int       `gorm:"column:number" json:"number"`
	Name       string    `gorm:"column:name" json:"name"`
	Website    string    `gorm:"column:website;type:varchar(2048)" json:"website"`
	Twitter    string    `gorm:"column:twitter;type:varchar(2048)" json:"twitter"`
	Github     string    `gorm:"column:github;type:varchar(2048)" json:"github"`
	Location   string    `gorm:"column:location;type:varchar(2048)" json:"location"`
	Tagline    string    `gorm:"column:tag_line;type:text" json:"tagline"`
	Avatar     string    `gorm:"column:avatar;type:varchar(2048)" json:"avatar"`
	Status     string    `gorm:"column:status" json:"status"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
}

func (m *Member) TableName() string {
	return memberTableName
}

func NewFakeMember(number int) *Member {
	return &Member{
		Number:     number,
		Name:       "@faker@",
		Status:     "not found",
		CreateTime: time.Unix(0, 0),
	}
}

func SaveMember(m *Member) (int, error) {
	err := core.MYSQL.Save(m).Error
	return m.Id, err
}

// FindLastMember 获取最新的会员数据
func FindLastMember() (*Member, error) {
	var m Member
	err := core.MYSQL.Order("number desc").First(&m).Error
	if err != nil {
		// 如果数据不存在,则返回一个假的数据,防止后续处理异常
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return EmptyMember, nil
		}
		// 否则抛出真正的错误
		return nil, err
	}
	return &m, err
}

func FindMembers(equalCondition *Member, offset, limit int, createTimeRange []time.Time) (int64, []*Member, error) {
	records := make([]*Member, 0)
	var count int64
	tx := core.MYSQL.Offset(offset).Limit(limit).Order("id desc")
	if len(createTimeRange) == 2 {
		tx.Where("create_time > ? and create_time < ?", createTimeRange[0], createTimeRange[1])
	}
	err := tx.Find(&records, equalCondition).Limit(-1).Offset(-1).Count(&count).Error
	return count, records, err
}
