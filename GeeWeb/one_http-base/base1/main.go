package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	//设置了两个路由，分别访问/和/hello,绑定了indexHandler和helloHandler
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	//print+os.Exit()
	log.Fatal(http.ListenAndServe(":9999", nil))

}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	//该响应返回路径
	fmt.Fprintf(w, "URL.Path=%q\n", req.URL.Path)
}
func helloHandler(w http.ResponseWriter, req *http.Request) {
	//获取请求头
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q]=%q\n", k, v)
	}
}

// 1. fmt.Sprint(w,"") ---> 格式化字符串并写入w，返回字符串的字节数
// 2. http.HandlerFunc(ResponseWriter,*Request)------>HTTP适配器。如果f是合适的函数，HandlerFunc(f)就是调用f的适配器
// 3. %q 字符串或整数占位符。双引号围绕的字符串，由Go语法安全转义。
