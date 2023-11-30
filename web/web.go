package web

import (
	"context"
	"github.com/gin-gonic/gin"
	"liewell.fun/v2ex/core"
	"net/http"
)

func StartAndWait(ctx context.Context) {

	cfg := core.GlobalConfig.Http

	// 服务运行模式
	if core.GlobalConfig.Zap.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 新建实例并配置全局中间件
	r := gin.New()
	r.Use(middleLogger(nil), middleWareCORS(), middleWareRecovery())

	// 注册路由
	registryHandler(r)

	// 启动服务
	if len(cfg.ListenAddrTLS) > 0 {
		go func() {
			if err := r.RunTLS(cfg.ListenAddrTLS, cfg.CertFile, cfg.KeyFile); err != nil {
				core.Logger.Fatalf("[web] https server error: %v", err)
			}
		}()
	}
	if len(cfg.ListenAddr) > 0 {
		go func() {
			if err := r.Run(cfg.ListenAddr); err != nil {
				core.Logger.Fatalf("[web] http server error: %v", err)
			}
		}()
	}
	core.Logger.Infof("[web] start success with http[%v], https[%v]", cfg.ListenAddr, cfg.ListenAddrTLS)

	select {
	case <-ctx.Done():
		core.Logger.Fatalf("[web] server shutdown: %v", ctx.Err())
	}
}

func registryHandler(engine *gin.Engine) {

	// 静态目录
	engine.Static("/static", "static")

	// html 模板
	engine.LoadHTMLGlob("static/template/*")

	// 根目录定位到首页
	engine.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// 会员数量
	engine.GET("/members/count", RenderMembersCount)

	// 会员趋势
	engine.GET("/members/trend", RenderMembersTrend)

	// 会员地域分布
	engine.GET("/members/location", RenderMembersLocation)

	// API 分组
	api := engine.Group("/api")
	api.GET("/", ListAPIs)
	api.GET("/members/:pageNo/:pageSize", ListMembers)
}
