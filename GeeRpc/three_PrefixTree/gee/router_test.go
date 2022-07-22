package gee

import (
	"fmt"
	"reflect"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	//method pattern handler
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/asserts/*filepath", nil)
	return r
}

//1. func(c *T)Fatal() == Log()+FailNow()
//FailNow 将当前测试函数标识为失败 并停止执行 之后 测试过程会在下一个测试或下一个基准测试继续
//必须在运行测试函数或基准测试函数的一个goroutine中调用，而不能在测试时创建的goroutine中调用 调用goroutine不会导致其他goroutinue停止
//Log 使用与Printf相同的格式化语法对参数进行格式化 记录到错误日志
//2 reflect.DeepEqual() 对array struct map slice等结构体判断值是否相同

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}

}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.getRoute("GET", "/hello/xiaoli")
	if n == nil {
		t.Fatal("nil should't be returned")
	}
	if n.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}
	if ps["name"] != "xiaoli" {
		t.Fatal("should be equal to 'xiaoli' ")
	}
	fmt.Printf("matched path:%s,params['name']:%s\n", n.pattern, ps["name"])
}
