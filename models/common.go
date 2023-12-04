package models

// KV 用来接收特定返回值
// 无需真的定义与真实字段匹配的字段,只要使用别名可以接收就完了
type KV struct {
	KeyOne string `gorm:"column:key_one"` // 使用 keyOne 字段第一个字段
	Date   string `gorm:"column:date"`    // 使用 date 字段接收时间
	Count  int    `gorm:"column:count"`   // 使用 count 字段接收统计的数据
}
