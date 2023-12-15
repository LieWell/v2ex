package models

import "liewell.fun/v2ex/core"

var (
	systemConfigTableName = "system_config"
	EmptySystemConfig     = &SystemConfig{}
)

var (
	SystemConfigKeyLastDrawTime = "last_draw_time"
)

type SystemConfig struct {
	Id    int    `gorm:"primaryKey;autoIncrement:true;column:id"`
	Key   string `gorm:"column:key"`
	Value string `gorm:"column:value"`
}

func (s *SystemConfig) TableName() string {
	return systemConfigTableName
}

func UpdateSystemConfig(key, value string) error {
	return core.MYSQL.Update("`value`=?", value).Where("`key`=?", key).Error
}

func FindSystemConfig(key string) (string, error) {
	var systemConfig SystemConfig
	err := core.MYSQL.Model(EmptySystemConfig).Where("`key`=?", key).Find(&systemConfig).Error
	return systemConfig.Value, err
}
