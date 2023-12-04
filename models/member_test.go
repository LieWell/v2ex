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

func TestCountMemberByYear(t *testing.T) {
	prepareMysql()
	kvList, err := CountMember()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v", kvList)
}

func TestFind404Members(t *testing.T) {
	prepareMysql()
	records, err := Find404Members(1, 99999)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v", len(records))
}

func TestCountStatus(t *testing.T) {
	prepareMysql()
	records, err := CountMemberStatus()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v", records)
}
