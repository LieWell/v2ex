package main

import (
	"context"
	"flag"
	"liewell.fun/v2ex/core"
	"liewell.fun/v2ex/v2ex"
	"liewell.fun/v2ex/web"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// 全局 context 控制
	ctx, cancel := context.WithCancel(context.Background())
	go WaitTerm(cancel)

	// 读取配置文件并解析
	c := flag.String("c", "config.yaml", "config file(default: config.yaml)")
	flag.Parse()
	core.LoadYamlConfig(*c)

	// 初始化日志模块,必须是第一个被初始化的模块
	core.InitZap()

	// 初始化数据库
	core.InitMysql()

	// 启动爬虫任务
	//go v2ex.StartSpider()

	// 启动分析任务
	go v2ex.StartDrawCharts()

	// 启动 web 服务
	web.StartAndWait(ctx)
}

func WaitTerm(cancel context.CancelFunc) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)
	<-quit
	cancel()
}
