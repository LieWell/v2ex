package web

import (
	"github.com/gin-gonic/gin"
	"liewell.fun/v2ex/core"
	"liewell.fun/v2ex/models"
	"net/http"
)

func RenderIndex(context *gin.Context) {

	// 统计当前有效会员数与无效会员数量
	status, err := models.CountMemberStatus()
	if err != nil {
		core.Logger.Warnf("[RenderIndex] count member status error: %v", err)
	}
	var pageMap = map[string]int{
		"normalCount":   0,
		"notFoundCount": 0,
		"totalCount":    0,
	}
	for _, kv := range status {
		switch kv.KeyOne {
		case models.MemberStatusFound:
			pageMap["normalCount"] = kv.Count
		case models.MemberStatusNotFound:
			pageMap["notFoundCount"] = kv.Count
		}
	}
	pageMap["totalCount"] = pageMap["normalCount"] + pageMap["notFoundCount"]
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
