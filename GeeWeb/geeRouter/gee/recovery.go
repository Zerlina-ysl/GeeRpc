package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func Recovery() HandlerFunc {
	return func(context *Context) {
		defer func() {

			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				context.Fail(http.StatusInternalServerError, "Internal Server Error")

			}
		}()
		//调用下一个中间件
		context.Next()
	}

}

//触发panic的堆栈信息
func trace(message string) string {
	var pcs [32]uintptr
	//调用goroutine堆栈上函数调用的返回程序计数器填充切片pc
	//skip上要跳过的堆栈帧数
	//返回堆栈中的程序计数器 第0个是Callers本身 第一个是上一层trace 第二个是defer func
	n := runtime.Callers(3, pcs[:])

	var str strings.Builder
	str.WriteString(message + "\nTraceBack:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		//获取到调用函数的函数名和行号
		file, line := fn.FileLine(pc)
		//日志带你
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))

	}
	return str.String()
}
