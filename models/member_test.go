package models

import (
	"flag"
	"fmt"
	"liewell.fun/v2ex/core"
	"testing"
)

func prepareMysql() {
	c := flag.String("c", "../config.yaml", "")
	flag.Parse()
	core.LoadYamlConfig(*c)
	core.InitZap()
	core.InitMysql()
}

func TestFindLastMember(t *testing.T) {
	prepareMysql()
	member, err := FindLastMember()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v", member)
}
