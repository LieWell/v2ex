package web

import (
	"errors"
	"github.com/gin-gonic/gin"
	"liewell.fun/v2ex/core"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

func middleWareRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne.Err, &se) {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				if brokenPipe {
					core.Logger.Infof("recovery: %v", err)
					_ = c.Error(err.(error))
					c.Abort()
					return
				}

				core.Logger.Infof("recovery: %s \n %s", err, debug.Stack())
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func middleWareCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, Attributes, Access-Token")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func middleLogger(skipPathArr []string) gin.HandlerFunc {

	skipPaths := make(map[string]bool, len(skipPathArr))
	for _, path := range skipPathArr {
		skipPaths[path] = true
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		c.Next()

		if _, ok := skipPaths[path]; !ok {
			end := time.Now()
			latency := end.Sub(start)

			if len(c.Errors) > 0 {
				// Append error field if this is an erroneous request.
				for _, e := range c.Errors.Errors() {
					core.Logger.Error(e)
				}
			} else {
				if c.Request.Method != http.MethodOptions {
					core.Logger.Debugf("%s %s status: %d cost: %v", strings.ToUpper(c.Request.Method), c.Request.RequestURI, c.Writer.Status(), latency)
				}
			}
		}
	}
}
