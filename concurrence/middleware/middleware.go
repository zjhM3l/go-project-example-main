package middleware

import (
	"context"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// 打印每个请求的request和response
func main() {
	h := server.New()

	h.POST("/login", func (c context.Context, ctx *app.RequestContext) {
		// print request
		log.Printf("Received RawRequest: %s", ctx.Request.RawRequest())

		// some biz logic
		ctx.JSON(200, "OK")

		//print response
		log.Printf("Send RawResponse: %s", ctx.Response.RawResponse())
	})

	h.POST("/logout", func (c context.Context, ctx *app.RequestContext) {
		// print request
		log.Printf("Received RawRequest: %s", ctx.Request.RawRequest())

		// some biz logic
		ctx.JSON(200, "OK")

		//print response
		log.Printf("Send RawResponse: %s", ctx.Response.RawResponse())
	})

	h.Spin()
}

// 使用中间件执行打印操作
func middleware() {
	h := server.New()

	h.Use(func(c context.Context, ctx *app.RequestContext) {
		// print request
		log.Printf("Received RawRequest: %s", ctx.Request.RawRequest())

		// next middleware or handler
		ctx.Next()

		//print response
		log.Printf("Send RawResponse: %s", ctx.Response.RawResponse())
	})

	h.POST("/login", func (c context.Context, ctx *app.RequestContext) {
		// some biz logic
		ctx.JSON(200, "OK")
	})

	h.POST("/logout", func (c context.Context, ctx *app.RequestContext) {
		// some biz logic
		ctx.JSON(200, "OK")
	})
}


