package web

import (
	"github.com/gin-gonic/gin"
	"liewell.fun/v2ex/core"
	"liewell.fun/v2ex/models"
	"net/http"
)

func RenderIndex(context *gin.Context) {

	// 查询最后更新事件
	lastDrawTime, err := models.FindSystemConfig(models.SystemConfigKeyLastDrawTime)
	if err != nil {
		core.Logger.Warnf("[RenderIndex] query last draw time error: %v", err)
	}

	// 统计当前有效会员数与无效会员数量
	var normalCount, notFoundCount int
	status, err := models.CountMemberStatus()
	if err != nil {
		core.Logger.Warnf("[RenderIndex] count member status error: %v", err)
	}
	for _, kv := range status {
		switch kv.KeyOne {
		case models.MemberStatusFound:
			normalCount = kv.Count
		case models.MemberStatusNotFound:
			notFoundCount = kv.Count
		}
	}

	// 设置到页面
	var pageMap = map[string]any{
		"normalCount":   normalCount,
		"notFoundCount": notFoundCount,
		"totalCount":    normalCount + notFoundCount,
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
