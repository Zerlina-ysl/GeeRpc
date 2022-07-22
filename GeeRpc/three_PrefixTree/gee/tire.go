package gee

import (
	"fmt"
	"strings"
)

//树节点
type node struct {
	pattern  string  //待匹配路由 如/p/:lang 完整url
	part     string  //该节点在该层的路由名称	如 /:lang
	children []*node //该节点的子节点
	isWild   bool    //精准匹配-->含有：或*时 为true

}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s,part=%s,isWild=%t}", n.pattern, n.part, n.isWild)
}

//返回当前节点匹配成功的子节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		//动态匹配做强校验 防止路由注册覆盖
		if child.part == part || child.isWild || part[0] == ':' || part[0] == '*' {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	wildNodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part {
			//静态路由节点优先
			nodes = append(nodes, child)
		}
		if child.isWild {
			//动态路由由节点延后
			wildNodes = append(wildNodes, child)
		}
	}
	nodes = append(nodes, wildNodes...)
	return nodes
}

//边匹配边插入
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		//最后一层路由 直接插入
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	//节点未匹配，则新建子节点
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		if child.isWild && len(n.children) > 0 {
			panic(part + "already has ambigious router")
		}
		if part[0] == '*' && len(parts) > height+1 {
			panic(part + "fuzzy matching *  can't exist between the path")
		}
		if (part[0] == '*' || part[0] == ':') && len(part) == 1 {
			panic(part + "fuzzy symbols cant exist alone")
		}
		n.children = append(n.children, child)
	}
	//递归查找每一层节点
	child.insert(pattern, parts, height+1)
}

//根据方法在roots中找到的对应的路由前缀树 返回对应请求路由的节点
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		//指针影响原list
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}
