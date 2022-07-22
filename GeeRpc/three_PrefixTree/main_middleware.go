package main

import (
	"log"
	"net/http"
	"personalCode/GeeRpc/three_PrefixTree/gee"
	"time"
)

func onlyForV2() gee.HandlerFunc {
	return func(context *gee.Context) {
		t := time.Now()
		context.Fail(500, "Internet Server Error")
		log.Printf("[%d] %s in %v for group v2", context.StatusCode, context.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := gee.New()
	r.Use(gee.Logger())
	r.GET("/", func(context *gee.Context) {
		context.HTML(http.StatusOK, "<h1>hello</h1>")
	})
	v2 := r.Group("/v2")
	v2.Use(onlyForV2())
	{
		v2.GET("/hello/:name", func(context *gee.Context) {

			context.String(http.StatusOK, "hello %s,you are at %s\n", context.Param("name"), context.Path)
		})
	}
	r.Run(":9999")

}
