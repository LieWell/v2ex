package models

import (
	"fmt"
	"testing"
	"time"

	"liewell.fun/v2ex/core"
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

func TestGetSystemConfig(t *testing.T) {
	prepareMysql()
	k := SystemConfigKeyLastAvatarId
	v, err := FindSystemConfig(k)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(v)
}
