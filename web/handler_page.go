package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"liewell.fun/v2ex/core"
	"liewell.fun/v2ex/models"
)

func RenderIndex(context *gin.Context) {

	// 查询最后制表时间
	lastDrawTime, err := models.FindSystemConfig(models.SystemConfigKeyLastDrawTime)
	if err != nil {
		core.Logger.Warnf("[RenderIndex] query last draw time error: %v", err)
	}

	// 使用最后会员数据作为数量以及统计截止时间
	member, err := models.FindLastMember()
	if err != nil {
		core.Logger.Warnf("[RenderIndex] query last member error: %v", err)
	}

	// 设置到页面
	var pageMap = map[string]any{
		"totalCount":    member.Number,
		"lastCrawlTime": member.CreateTime.Format(core.DefaultDayFormat),
		"lastDrawTime":  lastDrawTime,
	}
	context.HTML(http.StatusOK, "index.html", pageMap)
}

func RenderMembersCount(context *gin.Context) {
	context.HTML(http.StatusOK, "members_count.html", nil)
}

func RenderMembersTrend(context *gin.Context) {
	context.HTML(http.StatusOK, "members_trend.html", nil)
}

func RenderMembersLocation(context *gin.Context) {
	context.HTML(http.StatusOK, "members_location.html", nil)
}

func RenderMembersMosaic(context *gin.Context) {
	context.HTML(http.StatusOK, "mosaic.html", nil)
}
