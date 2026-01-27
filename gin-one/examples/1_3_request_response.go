// ============================================================================
// 1.3 请求与响应
// ============================================================================
// 运行方式: go run examples/1_3_request_response.go
// ============================================================================

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ============================================================================
// 核心概念：Gin 的请求与响应处理
// ============================================================================
//
// 【gin.Context 是什么？】
//
// Context 是 Gin 最重要的结构体，每个请求都会创建一个 Context 实例
// 它封装了：
// - http.Request: 原始请求对象
// - http.ResponseWriter: 响应写入器
// - Params: 路径参数
// - Keys: 请求级别的 KV 存储（中间件传递数据）
// - Errors: 错误收集
// - handlers: 中间件/处理函数链
//
// 【为什么不直接用 http.Request？】
//
// Gin 的 Context 提供了更便捷的 API：
// - 统一的数据绑定 (JSON, XML, Form)
// - 多种响应格式 (JSON, XML, HTML, String, File)
// - 请求数据验证
// - 中间件支持
//
// ============================================================================

// User 用户结构体
type User struct {
	ID       int    `json:"id" xml:"id"`
	Name     string `json:"name" xml:"name"`
	Email    string `json:"email" xml:"email"`
	Age      int    `json:"age,omitempty" xml:"age,omitempty"` // omitempty: 零值不输出
	IsActive bool   `json:"is_active" xml:"is_active"`
}

func main() {
	r := gin.Default()

	// ========================================================================
	// 一、JSON 响应
	// ========================================================================

	// 1. 使用 gin.H (快速构造)
	r.GET("/json/simple", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello",
			"status":  "success",
		})
	})

	// 2. 使用结构体 (推荐: 类型安全、可复用)
	r.GET("/json/struct", func(c *gin.Context) {
		user := User{
			ID:       1,
			Name:     "张三",
			Email:    "zhangsan@example.com",
			Age:      25,
			IsActive: true,
		}
		c.JSON(http.StatusOK, user)
	})

	// 3. 使用切片返回列表
	r.GET("/json/list", func(c *gin.Context) {
		users := []User{
			{ID: 1, Name: "张三", Email: "zhangsan@example.com"},
			{ID: 2, Name: "李四", Email: "lisi@example.com"},
		}
		c.JSON(http.StatusOK, gin.H{
			"data":  users,
			"total": len(users),
		})
	})

	// 4. PureJSON - 不转义 HTML 字符
	r.GET("/json/pure", func(c *gin.Context) {
		// JSON() 会把 < > & 转义为 \u003c \u003e \u0026
		// PureJSON() 保持原样
		c.PureJSON(http.StatusOK, gin.H{
			"html": "<b>Hello</b>",
		})
	})

	// 5. IndentedJSON - 格式化输出 (调试用)
	r.GET("/json/indented", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{
			"nested": gin.H{
				"key": "value",
			},
		})
	})

	// 6. SecureJSON - 防止 JSON 劫持 (在数组前加 while(1);)
	r.GET("/json/secure", func(c *gin.Context) {
		names := []string{"alice", "bob", "charlie"}
		// 输出: while(1);["alice","bob","charlie"]
		c.SecureJSON(http.StatusOK, names)
	})

	// 7. JSONP - 跨域回调 (已被 CORS 替代，了解即可)
	r.GET("/json/jsonp", func(c *gin.Context) {
		// URL: /json/jsonp?callback=myFunc
		// 输出: myFunc({"message":"hello"});
		c.JSONP(http.StatusOK, gin.H{"message": "hello"})
	})

	// 8. AsciiJSON - 非 ASCII 字符转义
	r.GET("/json/ascii", func(c *gin.Context) {
		c.AsciiJSON(http.StatusOK, gin.H{
			"message": "你好世界", // 输出: \u4f60\u597d\u4e16\u754c
		})
	})

	// ========================================================================
	// 二、XML 响应
	// ========================================================================

	r.GET("/xml", func(c *gin.Context) {
		user := User{
			ID:       1,
			Name:     "张三",
			Email:    "zhangsan@example.com",
			IsActive: true,
		}
		c.XML(http.StatusOK, user)
	})

	// ========================================================================
	// 三、YAML 响应
	// ========================================================================

	r.GET("/yaml", func(c *gin.Context) {
		c.YAML(http.StatusOK, gin.H{
			"name":    "gin",
			"version": "1.9.0",
		})
	})

	// ========================================================================
	// 四、纯文本响应
	// ========================================================================

	r.GET("/string", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, %s!", "World")
	})

	// ========================================================================
	// 五、HTML 响应 (模板渲染)
	// ========================================================================

	// 加载模板 (需要先创建 templates 目录和模板文件)
	// r.LoadHTMLGlob("templates/*")
	// 或者加载多个目录
	// r.LoadHTMLGlob("templates/**/*")

	r.GET("/html", func(c *gin.Context) {
		// 由于没有模板文件，这里用 Data 返回 HTML
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
			<!DOCTYPE html>
			<html>
			<head><title>Gin Demo</title></head>
			<body><h1>Hello Gin!</h1></body>
			</html>
		`))
	})

	// ========================================================================
	// 六、文件响应
	// ========================================================================

	// 下载文件
	r.GET("/file/download", func(c *gin.Context) {
		// c.File("./files/report.pdf")
		// 或者指定下载文件名
		// c.FileAttachment("./files/report.pdf", "月度报告.pdf")

		// 演示：返回动态生成的内容作为文件
		c.Header("Content-Disposition", "attachment; filename=demo.txt")
		c.Data(http.StatusOK, "text/plain", []byte("This is file content"))
	})

	// ========================================================================
	// 七、重定向
	// ========================================================================

	// HTTP 重定向
	r.GET("/redirect/external", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://www.google.com/")
	})

	// 内部路由重定向
	r.GET("/redirect/internal", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/json/simple")
	})

	// ========================================================================
	// 八、获取请求数据
	// ========================================================================

	// ---------- 8.1 获取 Header ----------
	r.GET("/header", func(c *gin.Context) {
		// 获取单个 Header (不区分大小写)
		userAgent := c.GetHeader("User-Agent")
		contentType := c.GetHeader("Content-Type")

		// 获取自定义 Header
		customHeader := c.GetHeader("X-Custom-Header")

		c.JSON(http.StatusOK, gin.H{
			"user_agent":    userAgent,
			"content_type":  contentType,
			"custom_header": customHeader,
		})
	})

	// ---------- 8.2 获取 Cookie ----------
	r.GET("/cookie/get", func(c *gin.Context) {
		// 获取 Cookie
		value, err := c.Cookie("session_id")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": "cookie not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"session_id": value})
	})

	// ---------- 8.3 设置 Cookie ----------
	r.GET("/cookie/set", func(c *gin.Context) {
		// SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
		c.SetCookie(
			"session_id", // Cookie 名称
			"abc123",     // Cookie 值
			3600,         // 过期时间(秒)，0 表示会话结束时过期
			"/",          // 路径
			"",           // 域名，空表示当前域
			false,        // Secure (仅 HTTPS)
			true,         // HttpOnly (JS 不可访问)
		)
		c.JSON(http.StatusOK, gin.H{"message": "cookie set"})
	})

	// ---------- 8.4 获取 Form 表单数据 ----------
	r.POST("/form", func(c *gin.Context) {
		// Content-Type: application/x-www-form-urlencoded
		// 或 Content-Type: multipart/form-data

		// PostForm: 获取 POST 表单数据
		name := c.PostForm("name")
		email := c.DefaultPostForm("email", "default@example.com")

		// 同时获取 Query 和 PostForm
		// page := c.Query("page")       // 从 URL 查询参数
		// name := c.PostForm("name")    // 从表单

		c.JSON(http.StatusOK, gin.H{
			"name":  name,
			"email": email,
		})
	})

	// ---------- 8.5 获取原始 Body ----------
	r.POST("/raw", func(c *gin.Context) {
		// 读取原始 body (只能读一次)
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"raw_body":     string(body),
			"content_type": c.ContentType(),
		})
	})

	// ========================================================================
	// 九、设置响应头
	// ========================================================================

	r.GET("/custom-headers", func(c *gin.Context) {
		// 设置单个 Header
		c.Header("X-Custom-Header", "custom-value")

		// 设置多个 Header
		c.Header("X-Request-Id", "12345")
		c.Header("Cache-Control", "no-cache")

		c.JSON(http.StatusOK, gin.H{"message": "check response headers"})
	})

	// ========================================================================
	// 十、统一响应格式 (最佳实践)
	// ========================================================================

	// 成功响应
	r.GET("/api/success", func(c *gin.Context) {
		data := User{ID: 1, Name: "张三", Email: "zhangsan@example.com"}
		Success(c, data)
	})

	// 错误响应
	r.GET("/api/error", func(c *gin.Context) {
		Error(c, http.StatusBadRequest, "INVALID_PARAM", "参数错误")
	})

	// 分页响应
	r.GET("/api/list", func(c *gin.Context) {
		users := []User{
			{ID: 1, Name: "张三"},
			{ID: 2, Name: "李四"},
		}
		SuccessWithPagination(c, users, 100, 1, 10)
	})

	// ========================================================================
	// 十一、Content Negotiation (内容协商)
	// ========================================================================

	r.GET("/negotiate", func(c *gin.Context) {
		// 根据 Accept Header 返回不同格式
		data := User{ID: 1, Name: "张三", Email: "zhangsan@example.com"}

		// c.Negotiate 会根据 Accept Header 自动选择格式
		c.Negotiate(http.StatusOK, gin.Negotiate{
			Offered:  []string{gin.MIMEJSON, gin.MIMEXML, gin.MIMEYAML},
			Data:     data,
			HTMLName: "", // 如果是 HTML，使用的模板名
		})
	})

	r.Run(":8080")
}

// ============================================================================
// 统一响应格式封装 (生产级最佳实践)
// ============================================================================

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`              // 业务状态码
	Message string      `json:"message"`           // 提示信息
	Data    interface{} `json:"data,omitempty"`    // 数据
	Error   string      `json:"error,omitempty"`   // 错误详情
}

// PaginatedResponse 分页响应
type PaginatedResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta"`
}

// Meta 分页元数据
type Meta struct {
	Total   int `json:"total"`
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, httpCode int, errCode string, message string) {
	c.JSON(httpCode, Response{
		Code:    -1,
		Message: message,
		Error:   errCode,
	})
}

// SuccessWithPagination 分页成功响应
func SuccessWithPagination(c *gin.Context, data interface{}, total, page, perPage int) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Code:    0,
		Message: "success",
		Data:    data,
		Meta: Meta{
			Total:   total,
			Page:    page,
			PerPage: perPage,
		},
	})
}

// ============================================================================
// 测试命令
// ============================================================================
//
// # JSON 响应测试
// curl http://localhost:8080/json/simple
// curl http://localhost:8080/json/struct
// curl http://localhost:8080/json/list
// curl http://localhost:8080/json/pure
// curl http://localhost:8080/json/indented
//
// # XML 响应测试
// curl http://localhost:8080/xml
//
// # YAML 响应测试
// curl http://localhost:8080/yaml
//
// # Header 测试
// curl -H "X-Custom-Header: my-value" http://localhost:8080/header
//
// # Cookie 测试
// curl -c cookies.txt http://localhost:8080/cookie/set
// curl -b cookies.txt http://localhost:8080/cookie/get
//
// # Form 测试
// curl -X POST -d "name=张三&email=test@example.com" http://localhost:8080/form
//
// # Raw Body 测试
// curl -X POST -H "Content-Type: text/plain" -d "raw content" http://localhost:8080/raw
//
// # 统一响应格式测试
// curl http://localhost:8080/api/success
// curl http://localhost:8080/api/error
// curl http://localhost:8080/api/list
//
// # 内容协商测试
// curl -H "Accept: application/json" http://localhost:8080/negotiate
// curl -H "Accept: application/xml" http://localhost:8080/negotiate
//
// ============================================================================

// ============================================================================
// 易错点总结
// ============================================================================
//
// 1. 【Body 只能读一次】
//    c.GetRawData() 会消耗 Body，第二次读取返回空
//    解决: 使用 c.ShouldBindBodyWith() 或自行缓存
//
// 2. 【Header 名称大小写】
//    HTTP Header 不区分大小写，Gin 会自动处理
//    c.GetHeader("content-type") == c.GetHeader("Content-Type")
//
// 3. 【JSON 零值问题】
//    int 类型零值 0 会输出，如需省略用 omitempty
//    但 omitempty 对 bool 的 false 也会省略，需注意
//
// 4. 【PureJSON vs JSON】
//    JSON(): < > & 会转义为 \u003c \u003e \u0026 (安全但可能影响前端)
//    PureJSON(): 不转义，适合内部 API
//
// 5. 【Cookie 安全设置】
//    生产环境务必设置:
//    - Secure: true (仅 HTTPS 传输)
//    - HttpOnly: true (JS 不可访问，防 XSS)
//    - SameSite: Strict/Lax (防 CSRF)
//
// 6. 【响应已发送后不能再写】
//    调用 c.JSON() 后再调用 c.String() 会报错
//    中间件用 c.Abort() 后注意不要重复响应
//
// ============================================================================

// ============================================================================
// 练习题
// ============================================================================
//
// 1. 实现一个接口，根据 Query 参数 format=json/xml/yaml 返回不同格式
//
// 2. 封装一个通用的错误响应函数，支持:
//    - HTTP 状态码
//    - 业务错误码
//    - 错误信息
//    - 错误详情 (仅开发环境显示)
//
// 3. 实现一个日志追踪功能:
//    - 每个请求生成唯一 Request-Id
//    - 响应头返回 X-Request-Id
//    - 所有日志输出包含 Request-Id
//
// ============================================================================
