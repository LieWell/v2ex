package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"liewell.fun/v2ex/core"
	"liewell.fun/v2ex/v2ex"
	"liewell.fun/v2ex/web"
)

func main() {

	// 全局 context 控制
	ctx, cancel := context.WithCancel(context.Background())
	go WaitTerm(cancel)

	// 读取配置文件并解析
	c := flag.String("c", "config.yaml", "config file(default: confiÂg.yaml)")
	flag.Parse()
	core.LoadYamlConfig(*c)

	// 初始化日志模块,必须是第一个被初始化的模块
	core.InitZap()

	// 初始化数据库
	core.InitMysql()

	// 启动爬虫任务
	if core.GlobalConfig.V2ex.StartSpider {
		go v2ex.StartMemberSpider()
		go v2ex.StartTopicSpider()
	} else {
		core.Logger.Infof("v2ex spider switch is off!")
	}

	// 启动分析任务
	go v2ex.StartDrawCharts()
	go v2ex.StartAvatarSpider()

	// 启动 web 服务
	web.StartAndWait(ctx)
}

func WaitTerm(cancel context.CancelFunc) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)
	<-quit
	cancel()
}
