

这是一个“包装器”方法：
- 在 Go 里用结构体**组合（embedding）**来拦截 `WriteHeader`，记录状态码，再把调用转发给底层 `http.ResponseWriter`。

关键点（结合这段代码）：

1. `type responseWriter struct { http.ResponseWriter; statusCode int }`  
   这里把原始 `ResponseWriter` 嵌入进来（类似“组合/委托”），所以 `responseWriter` 也具备 `ResponseWriter` 的方法。

2. `func (rw *responseWriter) WriteHeader(code int)`  
   这是给 `*responseWriter` 定义的方法（方法接收者是指针），它**覆盖/拦截**了 `WriteHeader`：
    - 先记录 `rw.statusCode = code`
    - 再调用底层真正的 `rw.ResponseWriter.WriteHeader(code)` 把状态码写出去

3. `newResponseWriter` 默认 `statusCode` 为 `http.StatusOK`，因为很多 handler 只写 body 不显式调用 `WriteHeader`，Go 默认是 200。

一句话：这是 Go 里常见的“装饰器/代理”方式，用来在不改变原接口的前提下加上**观测/记录**能力。


总结如下（简洁版）：

- `func (rw *responseWriter) WriteHeader(...)` 是 **方法**，`(rw *responseWriter)` 表示方法接收者，说明这个方法属于 `*responseWriter` 类型。
- 只有 `*responseWriter` 实例才能调用它：`rw.WriteHeader(code)`。
- 这不是匿名函数，也不是继承，而是 Go 里用 **方法 + 接口** 实现多态的常见写法。