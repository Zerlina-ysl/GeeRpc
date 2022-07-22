package main

import (
	"net/http"
	"personalCode/GeeRpc/three_PrefixTree/gee"
)

//加入路由分组控制后进行测试

func main() {
	r := gee.New()
	r.GET("/GET", func(context *gee.Context) {
		context.HTML(http.StatusOK, "<h1>hello</h1>")
	})
	v1 := r.Group("/v1")

	v1.GET("/", func(context *gee.Context) {
		context.HTML(http.StatusOK, "<h1>heello</h1>")
	})
	v1.GET("/hellp", func(context *gee.Context) {
		context.String(http.StatusOK, "hello %s,you are in %s", context.Query("name"), context.Path)
	})

	v2 := r.Group("/v2")

	v2.GET("/hello/:name", func(context *gee.Context) {
		context.String(http.StatusOK, "hello %s,you are in %s", context.Param("name"), context.Path)
	})
	v2.POST("/hello/login", func(context *gee.Context) {
		context.JSON(http.StatusOK, gee.H{
			//获取href中参数
			"username": context.PostForm("username"),
			"password": context.PostForm("password"),
		})
	})

	r.Run(":9999")
}
