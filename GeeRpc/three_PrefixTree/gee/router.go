package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

//roots['GET'] roots['POST"]
//handlers['GET-/p/:lang/doc'] handlers['POST-/p/book']
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, iterm := range vs {
		if iterm != "" {
			parts = append(parts, iterm)
			if iterm[0] == '*' {
				break
			}
		}
	}
	return parts
}
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	//传来的请求路径
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0)
	if n != nil {
		//处理的是方法对应的路由前缀树
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				//将前缀树的对应节点和传来的请求根据深度相映射
				log.Printf("tire node name is %s,request node name is %s", part[1:], searchParts[index])
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		//暂存解析出的路由参数

		c.Params = params
		key := c.Method + "-" + n.pattern
		//将通过路由前缀匹配到的中间件handler列表加入到handlers的列表中 等待执行
		c.handlers = append(c.handlers, r.handlers[key])
		//r.handlers[key](c)
	} else {
		//c.String(http.StatusNotFound, "404 NOT FOUND:%s\n", c.Path)
		c.handlers = append(c.handlers, func(context *Context) {
			context.String(http.StatusNotFound, "404 not found:%s\n", context.Path)
		})
	}
	c.Next()
}
