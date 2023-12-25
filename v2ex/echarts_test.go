package v2ex

import (
	"flag"
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

func TestDrawMemberCountBar(t *testing.T) {
	prepareMysql()
	rootRenderDir = "../static/template/"
	DrawMemberCountBar("members_count.html")
}

func TestDrawMemberTrendLine(t *testing.T) {
	prepareMysql()
	rootRenderDir = "../static/template/"
	DrawMemberTrendLine("members_trend.html")
}
