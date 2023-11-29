package echarts

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RenderLine(context *gin.Context) {
	context.HTML(http.StatusOK, "line.html", nil)
}
