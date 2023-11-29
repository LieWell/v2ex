package apis

import (
	"github.com/gin-gonic/gin"
	"liewell.fun/v2ex/core"
	"net/http"
)

func ListAPIs(ctx *gin.Context) {
	apiList := []string{
		"/api",
		"/api/members/:pageNo/:pageSize",
	}
	ctx.JSON(http.StatusOK, core.NewWithSuccess(apiList))
}
