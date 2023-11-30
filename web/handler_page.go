package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RenderMembersCount(context *gin.Context) {
	context.HTML(http.StatusOK, "members_count.html", nil)
}

func RenderMembersTrend(context *gin.Context) {
	context.HTML(http.StatusOK, "members_trend.html", nil)
}

func RenderMembersLocation(context *gin.Context) {
	context.HTML(http.StatusOK, "members_location.html", nil)
}
