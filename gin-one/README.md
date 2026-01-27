# Gin 框架从入门到生产级实战

> 学习笔记 & 核心知识点 & 易错点总结

---

## 快速开始

```bash
# 进入项目目录
cd gin-one

# 安装依赖
go mod tidy

# 运行任意示例文件
go run examples/1_1_hello_world.go

# 测试接口
curl http://localhost:8080/ping
```

---

## 文件索引

所有学习文件位于 `examples/` 目录下，每个文件可独立运行。

### 阶段一：基础与核心

| 文件 | 知识点 | 运行命令 |
|------|--------|---------|
| `1_1_hello_world.go` | Gin 安装、gin.New() vs Default()、优雅关闭 | `go run examples/1_1_hello_world.go` |
| `1_2_routing.go` | RESTful 路由、路由组、路径/查询参数 | `go run examples/1_2_routing.go` |
| `1_3_request_response.go` | JSON/XML 响应、Header/Cookie 处理 | `go run examples/1_3_request_response.go` |

### 阶段二：数据处理与验证

| 文件 | 知识点 | 运行命令 |
|------|--------|---------|
| `2_1_model_binding.go` | ShouldBind 系列、多来源绑定 | `go run examples/2_1_model_binding.go` |
| `2_2_validation.go` | validator 标签、自定义校验器 | `go run examples/2_2_validation.go` |
| `2_3_file_upload.go` | 单/多文件上传、流式处理 | `go run examples/2_3_file_upload.go` |

### 阶段三：中间件机制

| 文件 | 知识点 | 运行命令 |
|------|--------|---------|
| `3_1_middleware_principle.go` | 洋葱模型、Next()/Abort() 原理 | `go run examples/3_1_middleware_principle.go` |
| `3_2_builtin_middleware.go` | Logger、Recovery 源码解析 | `go run examples/3_2_builtin_middleware.go` |
| `3_3_custom_middleware.go` | CORS、限流、JWT 认证中间件 | `go run examples/3_3_custom_middleware.go` |

### 阶段四：工程化与数据库集成

| 文件 | 知识点 | 运行命令 |
|------|--------|---------|
| `4_1_gorm_integration.go` | GORM CRUD、关联查询、事务 | `go run examples/4_1_gorm_integration.go` |
| `4_2_project_structure.go` | 项目分层、依赖注入（文档） | 查看代码注释 |
| `4_3_config_logging.go` | Viper 配置、Zap 日志 | `go run examples/4_3_config_logging.go` |

### 阶段五：进阶与部署

| 文件 | 知识点 | 运行命令 |
|------|--------|---------|
| `5_1_jwt_auth.go` | JWT 生成/验证、Token 刷新 | `go run examples/5_1_jwt_auth.go` |
| `5_2_swagger.go` | Swagger 注解、自动文档生成 | `go run examples/5_2_swagger.go` |
| `5_3_graceful_shutdown.go` | 优雅关闭、Docker/K8s 部署 | `go run examples/5_3_graceful_shutdown.go` |

---

## 核心知识点速查

### 1. gin.New() vs gin.Default()

```
┌─────────────────────────────────────────┐
│           gin.Default()                  │
│  ┌───────────────────────────────────┐  │
│  │         gin.New()                  │  │
│  │    + Logger 中间件                  │  │
│  │    + Recovery 中间件               │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
```

**结论**：生产环境用 `gin.New()` + 自定义中间件

### 2. 中间件执行顺序（洋葱模型）

```
请求 → Middleware1 进入 → Middleware2 进入 → Handler
                                              ↓
响应 ← Middleware1 返回 ← Middleware2 返回 ←──┘
```

### 3. 数据绑定方法对比

| 方法 | 数据来源 | 失败时行为 |
|------|---------|-----------|
| `ShouldBindJSON` | Body (JSON) | 返回 error |
| `ShouldBindQuery` | URL Query | 返回 error |
| `ShouldBindUri` | Path 参数 | 返回 error |
| `ShouldBind` | 自动检测 | 返回 error |
| `Bind` | 自动检测 | 返回 400 |

**结论**：始终用 `ShouldBind` 系列

### 4. HTTP 方法与 RESTful 设计

| 方法 | 用途 | 幂等 | 示例 |
|------|-----|------|-----|
| GET | 查询 | 是 | `GET /users/123` |
| POST | 创建 | 否 | `POST /users` |
| PUT | 全量更新 | 是 | `PUT /users/123` |
| PATCH | 部分更新 | 否 | `PATCH /users/123` |
| DELETE | 删除 | 是 | `DELETE /users/123` |

---

## 易错点汇总

### 阶段一：基础

| 易错点 | 错误做法 | 正确做法 |
|--------|---------|---------|
| 忽略 Run() 返回值 | `r.Run(":8080")` | `if err := r.Run(); err != nil {...}` |
| 生产环境用 Default | `gin.Default()` | `gin.New()` + 自定义中间件 |
| 硬编码端口 | `:8080` | 从环境变量读取 |

### 阶段二：数据处理

| 易错点 | 错误做法 | 正确做法 |
|--------|---------|---------|
| Body 多次读取 | 连续两次 `ShouldBindJSON` | 用 `ShouldBindBodyWith` |
| PATCH 零值问题 | `struct { Age int }` | `struct { Age *int }` 用指针 |
| 校验 tag 错误 | `binding:"oneof=a,b,c"` | `binding:"oneof=a b c"` 用空格 |

### 阶段三：中间件

| 易错点 | 错误做法 | 正确做法 |
|--------|---------|---------|
| Abort 后不 return | `c.Abort()` | `c.Abort(); return` |
| 异步用 context | `go func() { c.JSON(...) }` | `copy := c.Copy()` |

### 阶段四：数据库

| 易错点 | 错误做法 | 正确做法 |
|--------|---------|---------|
| Find 判断空 | `if result.Error != nil` | Find 不报 ErrNotFound |
| Updates 零值 | `Updates(struct{Age: 0})` | `Updates(map[string]any{"age": 0})` |
| 未配置连接池 | 默认配置 | 设置 MaxIdleConns, MaxOpenConns |

### 阶段五：部署

| 易错点 | 错误做法 | 正确做法 |
|--------|---------|---------|
| JWT Secret 太短 | `"secret"` | 至少 32 字节 |
| 不监听信号 | 直接 `r.Run()` | 注册 SIGTERM/SIGINT |
| Docker CMD 格式 | `CMD ./server` | `CMD ["./server"]` |

---

## 常用测试命令

```bash
# 基本请求
curl http://localhost:8080/ping

# JSON 请求
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"张三","email":"test@example.com"}'

# 带认证
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/profile

# 文件上传
curl -X POST http://localhost:8080/upload -F "file=@test.txt"

# 查询参数
curl "http://localhost:8080/search?keyword=gin&page=1&limit=10"
```

---

## 推荐项目结构

```
gin-project/
├── cmd/server/main.go      # 入口
├── internal/
│   ├── handler/            # HTTP Handler
│   ├── service/            # 业务逻辑
│   ├── repository/         # 数据访问
│   ├── model/              # 数据模型
│   ├── middleware/         # 中间件
│   └── router/             # 路由注册
├── pkg/                    # 公共代码
├── configs/                # 配置文件
├── Dockerfile
├── docker-compose.yaml
└── Makefile
```

---

## 依赖安装

```bash
# Gin 框架
go get -u github.com/gin-gonic/gin

# GORM (数据库)
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
go get -u gorm.io/driver/sqlite

# 配置管理
go get -u github.com/spf13/viper

# 日志
go get -u go.uber.org/zap
go get -u gopkg.in/natefinch/lumberjack.v2

# JWT
go get -u github.com/golang-jwt/jwt/v5

# Swagger
go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files

# 参数校验
go get -u github.com/go-playground/validator/v10
```

---

## 学习进度追踪

- [x] 1.1 Gin 安装与 Hello World
- [x] 1.2 路由详解
- [x] 1.3 请求与响应
- [x] 2.1 模型绑定
- [x] 2.2 参数校验
- [x] 2.3 文件处理
- [x] 3.1 中间件原理
- [x] 3.2 内置中间件
- [x] 3.3 自定义中间件
- [x] 4.1 GORM 集成
- [x] 4.2 项目分层
- [x] 4.3 配置与日志
- [x] 5.1 JWT 认证
- [x] 5.2 Swagger 文档
- [x] 5.3 部署上线

---

## 参考资源

- [Gin 官方文档](https://gin-gonic.com/docs/)
- [GORM 官方文档](https://gorm.io/docs/)
- [Validator 文档](https://pkg.go.dev/github.com/go-playground/validator/v10)
- [Zap 日志库](https://pkg.go.dev/go.uber.org/zap)
- [Viper 配置库](https://github.com/spf13/viper)
