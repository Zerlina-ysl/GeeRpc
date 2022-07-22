package main

import (
	"net/http"
	"personalCode/GeeRpc/three_PrefixTree/gee"
)

func main() {
	r := gee.Default()
	r.GET("/", func(context *gee.Context) {
		context.String(http.StatusOK, "hello world\n")
	})
	r.GET("/panic", func(context *gee.Context) {
		name := []string{"xoaoli"}
		context.String(http.StatusOK, name[100])
	})
	r.Run(":9999")
}
