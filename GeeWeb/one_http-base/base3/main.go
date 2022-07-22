package main

import (
	"fmt"
	"net/http"
	"personalCode/GeeRPC/one_http-base/base3/gee"
)

func main() {

	r := gee.New()
	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL.Path=%q\n", r.URL.Path)
	})
	r.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%q]=%q", k, v)
		}
	})
	r.Run(":9999")

}

//1.首先对于HandlerFunc来说，是框架提供用户的，用来定义路由映射的处理方法。在Engine中，添加了一个路由映射表(map[string]HandlerFunc).针对相同的路由，请求方法
//不同，可以映射不同的处理方法(handler),value是用户映射的处理方法。
//2. 用户调用(*Engine)GET()方法时，会将路由和方法注册到映射表router中。
//3. Engine实现的HandlerFunc：用来解析请求路径，如果存在路由映射表中，就执行注册的处理方法，如果查不到，就返回404 NOT FOUND
//4. 此时实现的小demo，已经实现了路由映射表、用户静态注册路由、启动服务等功能，还可以加入动态路由、中间件等
