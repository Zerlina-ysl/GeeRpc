## Context
1. 对于web服务来说，通常都是根据请求*http.Request构造响应http.ResponseWriter.但这两个接口提供的接口粒度太细。例如，如果我们要构造一个完整的响应，需要考虑：
   消息头和消息体，而Header里通常包括了状态码、消息类型等几乎每次请求都要设置的属性。如果不进行有效的封装，那么用户会大量的重复造轮子。  
拿JSON工具包举例子，以展示框架封装所带来的便利 
```go
//封装前
obj = map[string]interface{}{
"name": "geektutu",
"password": "1234",
}
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
encoder := json.NewEncoder(w)
if err := encoder.Encode(obj); err != nil {
http.Error(w, err.Error(), 500)
}
//封装后的json数据返回
c.JSON(http.StatusOK, gee.H{
"username": c.PostForm("username"),
"password": c.PostForm("password"),
})
```
2. 针对具体应用场景，封装请求和响应，简化设计，这是Context的出发点之一。但是对于框架还需要额外的功能：如动态解析路由时的
参数缓存、中间件所产生的信息缓存。**Context随着每一次请求的出现而产生，请求的结束而销毁**，因此当前请求有关的信息都可以使用Context来承载
