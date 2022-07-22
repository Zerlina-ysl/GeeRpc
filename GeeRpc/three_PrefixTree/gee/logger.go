package gee

import (
	"log"
	"time"
)

//通用的Logger中间件，记录请求到响应所花费的时间

func Logger() HandlerFunc {
	return func(context *Context) {
		//start
		t := time.Now()
		//调用handler的下一个中间件
		context.Next()
		log.Printf("Logger:[%d] %s in %v", context.StatusCode, context.Req.RequestURI, time.Since(t))
	}
}
