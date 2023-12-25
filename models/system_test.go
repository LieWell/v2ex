package models

import (
	"liewell.fun/v2ex/core"
	"testing"
	"time"
)

func TestUpdateSystemConfig(t *testing.T) {
	prepareMysql()
	k := SystemConfigKeyLastDrawTime
	v := time.Now().Format(core.DefaultTimeFormat)
	err := UpdateSystemConfig(k, v)
	if err != nil {
		t.Error(err)
		return
	}
}
