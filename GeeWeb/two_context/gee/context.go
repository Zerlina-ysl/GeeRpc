package gee

//封装Context结构体，构造常用基本方法
import (
	"encoding/json"
	"fmt"
	"net/http"
)

//简洁数据
type H map[string]interface{}

//此次请求中的信息
type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	StatusCode int
}

//init Context
func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
	}
}

func (c *Context) PostForm(key string) string {
	//获取href的查询参数
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	//Query()返回的时参数名到参数值数组的map关系
	return c.Req.URL.Query().Get(key)
}
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}
func (c *Context) SetHeader(key string, value string) {
	c.Req.Header.Set(key, value)
}

//构造字符串输出
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

//快速构造JSON响应
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

//快速构造Data响应
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

//快速构造HTML响应
func (c *Context) HTML(code int, html string) {
	c.Status(code)
	c.Writer.Write([]byte(html))
}
