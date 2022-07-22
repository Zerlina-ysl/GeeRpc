package main

import (
	"net/http"
	"personalCode/GeeRpc/two_context/gee"
)

func main() {

	r := gee.New()
	r.GET("/", func(context *gee.Context) {
		context.HTML(http.StatusOK, "<p>hello gee</p>")
	})
	r.GET("/hello", func(context *gee.Context) {
		context.String(http.StatusOK, "hello %s,you're at %s\n", context.Query("name"), context.Path)
	})
	r.POST("/login", func(context *gee.Context) {
		context.JSON(http.StatusOK, gee.H{
			"username": context.PostForm("username"),
			"password": context.PostForm("password"),
		})
	})
	r.Run(":9999")
}
