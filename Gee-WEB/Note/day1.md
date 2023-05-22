### Gee
#### Day1 HTTP 基础
第一天实现内容包括:
+ 设计并编写Gee结构体组成
+ 实现基本的路由GET，POST

基本设计思路：
+ 主结构体 **Engine**：
```go
// HandlerFunc defines the request handler used by gee
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Engine implement the interface of ServeHTTP
type Engine struct {
	router map[string]HandlerFunc
}
```
在这个版本的设计中，
**HandlerFunc** 是 gee 中 request 的执行函数。
**Engine** 结构体具有属性router，
+ router 是 string $\rightarrow$ HandlerFunc 的映射，在后面的实现中，
  **string** 的类型为 "method-pattern" 类型，其中
  **method** 是请求方式如 GET, POST $\dots$，
  **pattern** 为URL的路由如 "/", "/xxx" $\dots$。

进一步，完成 Engine 接受 http 请求路由之后的执行函数（实际上实在后面实现的 Run 函数中的监听中调用的。
```go
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
``` 
**ServeHTTP** 函数是实现了具体路由执行的方法。

接下来编写的内容是：
+ Gee 实例的构造函数
+ GET、POST 方法的实现
+ Gee 实例的运行函数 Run

在这里需要注意的是，GET、POST的作用应该是向 **Engine** 结构体中的 **router** 添加路由记录，而非执行，**真正的路由执行在 Run 后接受到 http 请求才进行。**
代码如下：
```go
// New is the constructor of gee.Engine
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

// GET defines the method to add GET request
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
```