package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	//origin
	Writer http.ResponseWriter
	Req    *http.Request
	//request info
	Path   string
	Method string
	Params map[string]string //存储解析后的参数，提供对路由参数对访问
	//response
	StatusCode int
	//middleware
	handlers []HandlerFunc
	index    int
	//engine 可以访问html模版
	engine *Engine
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

//c.Param(part)获取参数
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

//init Context
func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
		index:  -1,
	}
}
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		//handler(c)
		c.handlers[c.index](c)
	}
}
func (c *Context) PostForm(key string) string {
	//获取href的查询参数
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	//返回路由前缀树中的节点名称到请求url的节点名称
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

////快速构造HTML响应
func (c *Context) HTML(code int, html string) {
	c.Status(code)
	c.Writer.Write([]byte(html))
}

func (context *Context) HTMLWithExecutor(code int, name string, data interface{}) {
	//自动调用htmlji解析器对文件进行处理
	context.SetHeader("Content-Type", "text/html")
	context.Status(code)
	if err := context.engine.htmlTemplate.ExecuteTemplate(context.Writer, name, data); err != nil {
		context.Fail(500, err.Error())
	}
}
