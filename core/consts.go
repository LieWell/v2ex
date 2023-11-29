package core

import "time"

const (
	ServerName           = "v2ex spider"
	ServerVersion        = "0.0.10"
	DefaultTimeFormat    = "2006-01-02 15:04:05"
	DefaultLogTimeFormat = "2006-01-02 15:04:05.000"
)

var (
	DefaultLocation, _ = time.LoadLocation("Asia/Shanghai")
)
