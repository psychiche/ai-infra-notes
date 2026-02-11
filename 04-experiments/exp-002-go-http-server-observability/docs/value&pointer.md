

# 一、函数返回 `*Config` 还是 `Config`

## 1️⃣ 本质区别

| 返回值       | 语义       |
| --------- | -------- |
| `Config`  | 值语义（拷贝）  |
| `*Config` | 指针语义（共享） |

---

## 2️⃣ 选择原则（直接记这个）

### ✅ 用值（Config）

* struct 很小（几个字段）
* 只读数据
* 不需要共享状态
* 更偏“数据结构”

例：

```go
func LoadConfig() Config
```

---

### ✅ 用指针（*Config）

* struct 较大
* 需要共享
* 需要修改
* 更偏“实体”

例：

```go
func NewServer(cfg *Config) *Server
```

---

## 3️⃣ 工程建议

* Config / DTO → 通常返回值
* Server / DB / Client → 通常返回指针

---


# 二、工程级推荐结构

```go
type App struct {
    Config *Config
}

func main() {
    conf, err := LoadConfig(...)
    app := NewApp(conf)
    app.Run()
}
```


# 三、核心总结（最重要）

### 1️⃣ error 返回规则

> 成功一定返回 `nil`

---

### 2️⃣ 值 vs 指针判断

> 数据 → 值
> 实体 → 指针

---

### 3️⃣ 配置加载推荐

> 不要用全局变量
> 返回 (*Config, error)

---

