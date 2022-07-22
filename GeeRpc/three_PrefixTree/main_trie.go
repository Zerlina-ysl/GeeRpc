package main

import (
	"net/http"
	"personalCode/GeeRpc/three_PrefixTree/gee"
)

func main() {
	r := gee.New()
	r.GET("/", func(context *gee.Context) {
		context.HTML(http.StatusOK, "<h1>hellowolrd</h1>")
	})
	r.GET("/hello", func(context *gee.Context) {
		context.String(http.StatusOK, "hello %s,you're at %s", context.Query("name"), context.Path)
	})
	r.GET("/hello/:name", func(context *gee.Context) {
		context.String(http.StatusOK, "hello %s,you're at %s", context.Query("name"), context.Path)
	})
	r.GET("/assets/*filepath", func(context *gee.Context) {
		context.JSON(http.StatusOK, gee.H{"filepath": context.Param("filepath")})
	})
	r.Run(":9999")

}
