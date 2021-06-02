package group_lit

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
)

var ginApp *gin.Engine

func init() {
	initGinServer()
}

func initGinServer() {
	gin.SetMode(gin.ReleaseMode)

	ginApp = gin.New()
	ginAddMiddlewareHandlers(ginApp, middlewareNum)
	addRoutes(func(url string) {
		ginApp.Handle(http.MethodGet, url, ginHandle2)
	})
}

func ginHandle2(c *gin.Context) {
	// println(unsafe.Sizeof(*c))
}

// add gin middlewares
func ginAddMiddlewareHandlers(ginApp *gin.Engine, ct int) {
	for i := 0; i < ct; i++ {
		ginApp.Use(func(context *gin.Context) {
			ginMiddlewareHandle(context)
		})
	}
}

func ginMiddlewareHandle(ctx *gin.Context) {
	ctx.Next()
}

// start benchmark
func BenchmarkGinWebRouter(b *testing.B) {
	benchRequest(b, ginApp)
}