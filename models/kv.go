package models

// KV 用来接收特定返回值
// 无需真的定义与真实字段匹配的字段,只要使用别名可以接收就完了
type KV struct {
	IntOne    int    `gorm:"column:int_one"`    // 使用 IntOne 接收第一个整数字段
	IntTwo    int    `gorm:"column:int_two"`    // 使用 IntTwo 接收第二个整数字段
	StringOne string `gorm:"column:string_one"` // 使用 StringOne 接收第一个字符串字段
	StringTwo string `gorm:"column:string_two"` // 使用 StringTwo 接收第二个字符串字段
	Date      string `gorm:"column:date"`       // 使用 date 字段接收时间
	Count     int    `gorm:"column:count"`      // 使用 count 字段接收统计的数据
}
