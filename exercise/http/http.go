package http

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	h := server.New()

	// 对应的路由注册到server上，/sis是url，下面是业务逻辑
	h.POST("/sis", func (c context.Context, ctx *app.RequestContext) {
		ctx.Data(200, "text/plain; charset=utf-8", []byte("OK"))
	})

	h.Spin()
}