// ============================================================================
// 2.1 模型绑定 (Model Binding)
// ============================================================================
// 运行方式: go run examples/2_1_model_binding.go
// ============================================================================

package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// ============================================================================
// 核心概念：什么是模型绑定？
// ============================================================================
//
// 模型绑定是将 HTTP 请求数据自动映射到 Go 结构体的过程
//
// 【请求数据来源】
// 1. Query 参数: /user?name=tom&age=18
// 2. Path 参数: /user/:id
// 3. Form 表单: application/x-www-form-urlencoded 或 multipart/form-data
// 4. JSON Body: application/json
// 5. XML Body: application/xml
// 6. Header: 请求头
//
// 【绑定的本质】
// Gin 使用反射读取结构体的 tag，然后从对应来源获取数据并赋值
//
// ============================================================================
//
// 【Bind vs ShouldBind 系列对比】
//
// ┌─────────────────┬──────────────────┬───────────────────────────────────┐
// │ 方法             │ 绑定失败时        │ 适用场景                           │
// ├─────────────────┼──────────────────┼───────────────────────────────────┤
// │ Bind            │ 返回 400 并中止   │ 不推荐 (无法自定义错误响应)         │
// │ ShouldBind      │ 仅返回 error     │ 推荐 (可自定义错误处理)             │
// │ MustBind        │ panic            │ 几乎不用                           │
// └─────────────────┴──────────────────┴───────────────────────────────────┘
//
// 【ShouldBind 系列方法】
//
// | 方法                | 数据来源                    | Content-Type            |
// |---------------------|----------------------------|-------------------------|
// | ShouldBind          | 自动检测                    | 根据 Content-Type 自动选择 |
// | ShouldBindJSON      | Body                       | application/json         |
// | ShouldBindXML       | Body                       | application/xml          |
// | ShouldBindQuery     | URL Query                  | -                        |
// | ShouldBindUri       | Path 参数                   | -                        |
// | ShouldBindHeader    | Request Header             | -                        |
// | ShouldBindBodyWith  | Body (可重复读取)           | 指定类型                  |
//
// ============================================================================

// ============================================================================
// 结构体定义与 Tag 详解
// ============================================================================

// User 用户注册请求
// 【Tag 格式】`json:"字段名" form:"字段名" binding:"校验规则"`
type User struct {
	// json: JSON Body 中的字段名
	// form: Form 表单和 Query 参数中的字段名
	// binding: 校验规则 (后面章节详细讲)
	Name  string `json:"name" form:"name" binding:"required"`
	Email string `json:"email" form:"email" binding:"required,email"`
	Age   int    `json:"age" form:"age" binding:"gte=0,lte=150"`
}

// UserUri 路径参数绑定
type UserUri struct {
	ID   int    `uri:"id" binding:"required,gt=0"`
	Name string `uri:"name"` // 可选参数
}

// UserQuery 查询参数绑定
type UserQuery struct {
	Page     int    `form:"page,default=1"`       // 默认值
	PageSize int    `form:"page_size,default=10"` // 下划线命名
	Keyword  string `form:"keyword"`              // 可选
	Status   string `form:"status"`
}

// UserHeader 请求头绑定
type UserHeader struct {
	Authorization string `header:"Authorization" binding:"required"`
	UserAgent     string `header:"User-Agent"`
	RequestID     string `header:"X-Request-ID"`
}

// CreateUserRequest 复合绑定示例
type CreateUserRequest struct {
	// 可以组合多种来源
	Name     string `json:"name" form:"name" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,min=6"`
}

// UpdateUserRequest 更新请求 (部分字段可选)
type UpdateUserRequest struct {
	Name  *string `json:"name" form:"name"`   // 使用指针，nil 表示未提供
	Email *string `json:"email" form:"email"` // 使用指针区分 "未提供" 和 "空值"
	Age   *int    `json:"age" form:"age"`
}

func main() {
	r := gin.Default()

	// ========================================================================
	// 一、JSON Body 绑定
	// ========================================================================

	// POST /users/json
	// Content-Type: application/json
	// Body: {"name": "张三", "email": "test@example.com", "age": 25}
	r.POST("/users/json", func(c *gin.Context) {
		var user User

		// 【推荐】ShouldBindJSON - 明确指定绑定 JSON
		if err := c.ShouldBindJSON(&user); err != nil {
			// 自定义错误响应
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_request",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "user created",
			"user":    user,
		})
	})

	// ========================================================================
	// 二、Form 表单绑定
	// ========================================================================

	// POST /users/form
	// Content-Type: application/x-www-form-urlencoded
	// Body: name=张三&email=test@example.com&age=25
	r.POST("/users/form", func(c *gin.Context) {
		var user User

		// ShouldBind 会根据 Content-Type 自动选择绑定方式
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_request",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "user created (form)",
			"user":    user,
		})
	})

	// ========================================================================
	// 三、Query 参数绑定
	// ========================================================================

	// GET /users?page=2&page_size=20&keyword=test
	r.GET("/users", func(c *gin.Context) {
		var query UserQuery

		// ShouldBindQuery - 仅绑定 Query 参数
		if err := c.ShouldBindQuery(&query); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_query",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"query": query,
			// 演示默认值生效
			"note": fmt.Sprintf("Page: %d, PageSize: %d", query.Page, query.PageSize),
		})
	})

	// ========================================================================
	// 四、Path 参数绑定 (Uri)
	// ========================================================================

	// GET /users/123
	// GET /users/123/tom
	r.GET("/users/:id", func(c *gin.Context) {
		var uri UserUri

		// ShouldBindUri - 绑定路径参数
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_uri",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id": uri.ID,
		})
	})

	// 多个路径参数
	r.GET("/users/:id/:name", func(c *gin.Context) {
		var uri UserUri

		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id":   uri.ID,
			"user_name": uri.Name,
		})
	})

	// ========================================================================
	// 五、Header 绑定
	// ========================================================================

	r.GET("/profile", func(c *gin.Context) {
		var headers UserHeader

		// ShouldBindHeader - 绑定请求头
		if err := c.ShouldBindHeader(&headers); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "missing_header",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"authorization": headers.Authorization,
			"user_agent":    headers.UserAgent,
			"request_id":    headers.RequestID,
		})
	})

	// ========================================================================
	// 六、多来源复合绑定
	// ========================================================================

	// POST /users/:id?version=1
	// Body: {"name": "new name"}
	// Header: Authorization: Bearer xxx
	r.PUT("/users/:id", func(c *gin.Context) {
		// 分别绑定不同来源
		var uri UserUri
		var body UpdateUserRequest
		var header UserHeader

		// 绑定 URI
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uri", "detail": err.Error()})
			return
		}

		// 绑定 Body
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body", "detail": err.Error()})
			return
		}

		// 绑定 Header (跳过验证错误，Header 为可选)
		_ = c.ShouldBindHeader(&header)

		c.JSON(http.StatusOK, gin.H{
			"user_id":       uri.ID,
			"update_fields": body,
			"request_id":    header.RequestID,
		})
	})

	// ========================================================================
	// 七、ShouldBindBodyWith - Body 重复读取
	// ========================================================================
	//
	// 【问题】HTTP Body 是流式的，只能读取一次
	// 第一次 ShouldBindJSON 后，第二次读取会得到空数据
	//
	// 【解决方案】使用 ShouldBindBodyWith
	//

	r.POST("/users/double-bind", func(c *gin.Context) {
		var user1 User
		var user2 User

		// 使用 ShouldBindBodyWith，Gin 会缓存 Body
		if err := c.ShouldBindBodyWith(&user1, binding.JSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 第二次绑定仍然有效
		if err := c.ShouldBindBodyWith(&user2, binding.JSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user1": user1,
			"user2": user2,
			"equal": user1.Name == user2.Name,
		})
	})

	// ========================================================================
	// 八、指针字段处理 PATCH 更新
	// ========================================================================

	r.PATCH("/users/:id/partial", func(c *gin.Context) {
		var uri UserUri
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var req UpdateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 使用指针判断哪些字段需要更新
		updates := make(map[string]interface{})

		if req.Name != nil {
			updates["name"] = *req.Name
		}
		if req.Email != nil {
			updates["email"] = *req.Email
		}
		if req.Age != nil {
			updates["age"] = *req.Age
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id":       uri.ID,
			"fields_to_update": updates,
		})
	})

	// ========================================================================
	// 九、自定义绑定器 (实现 binding.Binding 接口)
	// ========================================================================
	// 这是高级用法，通常不需要，了解即可

	r.Run(":8080")
}

// ============================================================================
// 测试命令
// ============================================================================
//
// # JSON 绑定测试
// curl -X POST http://localhost:8080/users/json \
//   -H "Content-Type: application/json" \
//   -d '{"name":"张三","email":"test@example.com","age":25}'
//
// # 缺少必填字段
// curl -X POST http://localhost:8080/users/json \
//   -H "Content-Type: application/json" \
//   -d '{"name":"张三"}'
//
// # Form 表单绑定
// curl -X POST http://localhost:8080/users/form \
//   -d "name=张三&email=test@example.com&age=25"
//
// # Query 参数绑定
// curl "http://localhost:8080/users?page=2&page_size=20&keyword=test"
// curl "http://localhost:8080/users"  # 测试默认值
//
// # Path 参数绑定
// curl http://localhost:8080/users/123
// curl http://localhost:8080/users/abc  # 会报错，id 需要是数字
// curl http://localhost:8080/users/123/tom
//
// # Header 绑定
// curl -H "Authorization: Bearer token123" \
//   -H "X-Request-ID: req-001" \
//   http://localhost:8080/profile
//
// # 复合绑定
// curl -X PUT http://localhost:8080/users/123 \
//   -H "Content-Type: application/json" \
//   -H "X-Request-ID: req-002" \
//   -d '{"name":"新名字"}'
//
// # PATCH 部分更新
// curl -X PATCH http://localhost:8080/users/123/partial \
//   -H "Content-Type: application/json" \
//   -d '{"name":"只更新名字"}'
//
// curl -X PATCH http://localhost:8080/users/123/partial \
//   -H "Content-Type: application/json" \
//   -d '{"age":30}'
//
// ============================================================================

// ============================================================================
// 易错点总结
// ============================================================================
//
// 1. 【Bind vs ShouldBind】
//    Bind: 失败时自动返回 400，无法自定义错误格式
//    ShouldBind: 返回 error，可自定义处理 (推荐)
//
// 2. 【Body 只能读一次】
//    第一次 ShouldBindJSON 后，再次调用会失败
//    解决: 使用 ShouldBindBodyWith
//
// 3. 【Tag 名称必须与来源匹配】
//    JSON Body: 用 json tag
//    Form/Query: 用 form tag
//    Path 参数: 用 uri tag
//    Header: 用 header tag
//
// 4. 【默认值不支持 required】
//    `form:"page,default=1" binding:"required"` - 矛盾！
//    有默认值就不要 required
//
// 5. 【指针 vs 值类型】
//    PATCH 更新时用指针，可区分 "未提供" 和 "提供空值"
//    name: null  vs  不传 name  是不同的
//
// 6. 【ShouldBind 只能调用一次】
//    对于 Body 数据，ShouldBind/ShouldBindJSON 只能调用一次
//    第二次调用返回 EOF 错误
//
// 7. 【Content-Type 必须正确】
//    JSON Body 必须设置 Content-Type: application/json
//    否则 ShouldBindJSON 可能绑定失败
//
// ============================================================================

// ============================================================================
// 练习题
// ============================================================================
//
// 1. 实现一个搜索接口 GET /search，支持以下参数:
//    - keyword (必填)
//    - category (可选)
//    - min_price, max_price (可选，数字)
//    - page, page_size (可选，有默认值)
//
// 2. 实现一个用户信息更新接口 PATCH /users/:id
//    - 使用指针字段，只更新用户提供的字段
//    - 返回实际更新了哪些字段
//
// 3. 实现一个接口，同时从 Header、URI、Query、Body 获取数据
//    演示复合绑定的完整流程
//
// ============================================================================
