package main

import (
	"fmt"
	"html/template"
	"net/http"
	"personalCode/GeeRpc/three_PrefixTree/gee"
	"time"
)

type student struct {
	name string
	age  int
}

//funcMap中自定义的函数

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := gee.New()

	r.Use(gee.Logger())
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLBlob("template/*")
	//将相对路径./static中的文件映射为 :9999/assets
	r.Static("/assets", "./static")

	stu1 := &student{age: 11, name: "xiaoli"}
	stu2 := &student{age: 18, name: "xiaoyu"}
	r.GET("/", func(context *gee.Context) {
		context.HTMLWithExecutor(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/date", func(context *gee.Context) {
		context.HTMLWithExecutor(http.StatusOK, "arr.tmpl", gee.H{
			"title": "test",
			"now":   time.Date(2022, 7, 23, 0, 0, 0, 0, time.UTC),
		})
	})
	r.GET("/student", func(context *gee.Context) {
		context.HTMLWithExecutor(http.StatusOK, "student.tmpl", gee.H{
			"title":   "student",
			"student": [2]*student{stu1, stu2},
		})
	})
	r.Run(":9999")
}
