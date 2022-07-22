package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

type HandlerFunc func(*Context)

type Engine struct {
	router *router
	*RouterGroup
	groups []*RouterGroup
	//html render
	//将所有的模版加载进内存
	htmlTemplate *template.Template
	//自定义模版渲染函数
	funcMap template.FuncMap
}

//html render

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

//加载模版

func (engine *Engine) LoadHTMLBlob(pattern string) {
	//Must 只取第一个返回值 不返回err
	//New()创建模版
	//Funcs() --> FuncMap key为模版使用的函数名，value是函数对象
	engine.htmlTemplate = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine, prefix: ""}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) CreateStaticHandle(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := group.prefix + relativePath
	//过滤指定处理程序对应请求的特定前缀以显示文件目录
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(context *Context) {
		file := context.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			context.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(context.Writer, context.Req)
	}
}

//暴露给用户 用户可以通过本地静态资源映射到路由
// /usr/xiaoli/blog/static/js/xiaoli.js--->http:localhost:9999/asserts/js/xiaoli.js

func (group *RouterGroup) Static(relativePath string, root string) {
	//http.Dir()返回http.Dir()类型，用于将字符串类型转换为文件系统
	handler := group.CreateStaticHandle(relativePath, http.Dir(root))
	//贪心匹配
	urlPattern := path.Join(relativePath, "/*filepath")
	group.GET(urlPattern, handler)
}

type RouterGroup struct {
	prefix      string        //依靠前缀进行分组的路由控制
	middlewares []HandlerFunc //根据中间件进行分组下的功能控制
	parent      *RouterGroup  //支持分组嵌套
	engine      *Engine       //指针指向路由，由路由统一协调，通过路由间接访问接口
}

//初始化RouterGroup

func (group *RouterGroup) Group(prefix string) *RouterGroup {

	engine := group.engine
	newGroup := &RouterGroup{
		engine: engine,
		prefix: group.prefix + prefix,
		parent: group,
	}
	log.Printf("routerGroup prefix :%s", newGroup.prefix)
	engine.groups = append(engine.groups, newGroup)

	//下次使用该分组控制会作为前缀放在路由前
	return newGroup
}

//func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
//engine.router.addRoute(method, pattern, handler)
//}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	//log.Printf("addRoute prefix :%s", group.prefix)
	pattern := group.prefix + comp
	log.Printf("Route %4s-%s", method, pattern)
	//实现路由映射
	group.engine.router.addRoute(method, pattern, handler)
}

//func (engine *Engine) GET(pattern string, handler HandlerFunc) {
//	engine.addRoute("GET", pattern, handler)
//}
//func (engine *Engine) POST(pattern string, handler HandlerFunc) {
//	engine.addRoute("POST", pattern, handler)
//}
//交给RouterGroup实现路由有关函数 既可以添加路由 也可以实现分组添加路由

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		//根据请求路由前缀判断所对应的中间件
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	//每次创建新的Context 不会对同一个进行写入
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}

//中间件

//增加中间件处理函数
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

//Default 使用日志和错误恢复中间件
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}
