package main

import (
	"fmt"
	"log"
	"net/http"
)

type Engine struct{}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path=%q\n", req.URL)

	case "/hello":
		for k, v := range req.Header {
			fmt.Fprintf(w, "req[%q]=%q", k, v)
		}
	default:
		fmt.Fprintf(w, "404 Not Found:%s\n", req.URL)
	}

}

func main() {
	engine := &Engine{}
	log.Fatal(http.ListenAndServe(":9999", engine))
}

//1. go语言的接口实现是隐式的，非侵入式设计，go编译器自动在需要的时候检查。当一个接口的多个方法都被实现，接口才能被正常编译并使用。实现接口的类型方法要与接口的方法完全一致
//2. 定义了一个空的结构体Engine，Engine实现了Handler接口，重写了ServeHttp方法。
//param1: http.ResponseWriter,可以构造针对请求的响应
//param2: *http.Request,包含了HTTP请求的所有信息，包括请求地址、Header和Body等
//3. 向ListenAndServe传递engine实例，将所有的http请求转向了自己的处理逻辑。即 拦截了所有的http请求，拥有了统一的控制入口。在这里，我们可以实现一些扩展的功能：
//如：添加一些处理逻辑、自定义路由映射规则等
