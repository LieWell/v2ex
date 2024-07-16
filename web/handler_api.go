package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"liewell.fun/v2ex/core"
	"liewell.fun/v2ex/models"
	"liewell.fun/v2ex/v2ex"
)

func ListAPIs(ctx *gin.Context) {
	apiList := []string{
		"/api",
		// "/api/members/:pageNo/:pageSize",
	}
	ctx.JSON(http.StatusOK, core.NewWithSuccess(apiList))
}

func ListMembers(context *gin.Context) {

	pageCondition, err := core.BuildOffsetLimitWithGinContext(context)
	if err != nil {
		core.Logger.Errorf("[ListMembers] BuildOffsetLimitWithGinContext error: %v", err)
		context.JSON(http.StatusBadRequest, core.SimpleBadRequestError())
		return
	}
	offset, limit := pageCondition.Offset, pageCondition.Limit

	total, data, err := models.FindMembers(models.EmptyMember, offset, limit, nil)
	if err != nil {
		core.Logger.Errorf("[ListMembers] models.FindMembers error: %v", err)
		context.JSON(http.StatusInternalServerError, core.SimpleInternalServerError())
		return
	}

	response := core.HttpPageResponse{
		Data:       data,
		PageNo:     pageCondition.PageNo,
		PageSize:   pageCondition.PageSize,
		TotalCount: int(total),
	}
	context.JSON(http.StatusOK, core.NewWithSuccess(response))
}

func DrawPic(context *gin.Context) {
	v2ex.DrawMemberCountBar("members_count.html")
	v2ex.DrawMemberTrendLine("members_trend.html")
	context.JSON(http.StatusOK, core.SuccessResponse)
}
