// ============================================================================
// 1.2 路由详解
// ============================================================================
// 运行方式: go run examples/1_2_routing.go
// ============================================================================

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ============================================================================
// 核心概念：Gin 的路由系统
// ============================================================================
//
// 【Radix Tree (基数树/压缩前缀树)】
//
// Gin 使用 Radix Tree 存储路由，相比普通的 Map 存储：
// - 内存占用更小（共享公共前缀）
// - 查找效率 O(k)，k 是路径长度
// - 支持路径参数和通配符
//
// 示例路由树结构：
//
//                    [/]
//                     |
//            [api/v1/]─────────[health]
//                     |
//     [users]────────[orders]
//        |              |
//   [:id]─[/]      [:id]─[/items]
//
// ============================================================================
//
// 【RESTful API 设计原则】
//
// | HTTP 方法 | 语义       | 幂等性 | 安全性 | 示例              |
// |-----------|------------|--------|--------|-------------------|
// | GET       | 查询资源   | 是     | 是     | GET /users/123    |
// | POST      | 创建资源   | 否     | 否     | POST /users       |
// | PUT       | 全量更新   | 是     | 否     | PUT /users/123    |
// | PATCH     | 部分更新   | 否     | 否     | PATCH /users/123  |
// | DELETE    | 删除资源   | 是     | 否     | DELETE /users/123 |
//
// 幂等性：多次执行结果相同
// 安全性：不会修改资源
//
// ============================================================================

func main() {
	r := gin.Default()

	// ========================================================================
	// 一、基础路由注册
	// ========================================================================

	// 1. 基本的 HTTP 方法
	r.GET("/get", handleGet)
	r.POST("/post", handlePost)
	r.PUT("/put", handlePut)
	r.PATCH("/patch", handlePatch)
	r.DELETE("/delete", handleDelete)
	r.OPTIONS("/options", handleOptions)
	r.HEAD("/head", handleHead)

	// 2. Any - 匹配所有 HTTP 方法
	r.Any("/any", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"method": c.Request.Method,
		})
	})

	// ========================================================================
	// 二、路径参数 (Path Parameters)
	// ========================================================================
	//
	// 【语法】
	// :param  - 必选参数，匹配一个路径段
	// *param  - 通配符参数，匹配剩余所有路径
	//
	// 【底层原理】
	// Gin 在路由匹配时，会把 :param 部分提取出来，存入 c.Params
	// c.Params 是 []Param 类型，Param 结构体包含 Key 和 Value
	//

	// 示例 1：单个路径参数
	// GET /users/123 -> id = "123"
	r.GET("/users/:id", func(c *gin.Context) {
		// 方式一：c.Param() - 获取路径参数
		id := c.Param("id")

		c.JSON(http.StatusOK, gin.H{
			"user_id": id,
		})
	})

	// 示例 2：多个路径参数
	// GET /users/123/posts/456 -> userId = "123", postId = "456"
	r.GET("/users/:userId/posts/:postId", func(c *gin.Context) {
		userId := c.Param("userId")
		postId := c.Param("postId")

		c.JSON(http.StatusOK, gin.H{
			"user_id": userId,
			"post_id": postId,
		})
	})

	// 示例 3：通配符参数 (*)
	// GET /files/path/to/file.txt -> filepath = "/path/to/file.txt"
	//
	// 【注意】通配符必须在路径末尾
	r.GET("/files/*filepath", func(c *gin.Context) {
		filepath := c.Param("filepath")
		c.JSON(http.StatusOK, gin.H{
			"filepath": filepath, // 包含前导斜杠
		})
	})

	// ========================================================================
	// 三、查询参数 (Query Parameters)
	// ========================================================================
	//
	// URL 示例: /search?keyword=gin&page=1&limit=10
	//
	// 【常用方法对比】
	//
	// | 方法                    | 不存在时返回 | 用途                |
	// |-------------------------|--------------|---------------------|
	// | c.Query("key")          | ""           | 简单获取            |
	// | c.DefaultQuery("k","d") | 默认值 d     | 需要默认值          |
	// | c.GetQuery("key")       | "", false    | 需要判断是否存在    |
	// | c.QueryArray("key")     | []string{}   | 获取数组参数        |
	// | c.QueryMap("key")       | map[string]string{} | 获取 map     |
	//

	// GET /search?keyword=gin&page=1&limit=10&tags=go&tags=web
	r.GET("/search", func(c *gin.Context) {
		// 基本获取
		keyword := c.Query("keyword") // 不存在返回 ""

		// 带默认值
		page := c.DefaultQuery("page", "1")
		limit := c.DefaultQuery("limit", "10")

		// 判断参数是否存在
		sort, exists := c.GetQuery("sort")
		if !exists {
			sort = "created_at" // 手动设置默认值
		}

		// 获取数组参数 (同名多个值)
		// URL: ?tags=go&tags=web&tags=api
		tags := c.QueryArray("tags")

		// 获取 Map 参数
		// URL: ?filter[status]=active&filter[type]=premium
		filters := c.QueryMap("filter")

		c.JSON(http.StatusOK, gin.H{
			"keyword": keyword,
			"page":    page,
			"limit":   limit,
			"sort":    sort,
			"tags":    tags,
			"filters": filters,
		})
	})

	// ========================================================================
	// 四、路由组 (Router Group)
	// ========================================================================
	//
	// 【为什么需要路由组？】
	//
	// 1. 统一路径前缀，减少重复
	// 2. 组级别中间件，如：/api/v1 需要鉴权，/public 不需要
	// 3. 版本管理，如：/api/v1, /api/v2
	//
	// 【路由组的本质】
	//
	// RouterGroup 结构体包含：
	// - basePath: 基础路径
	// - Handlers: 组级别中间件
	// - engine:   指向 Engine
	//

	// API v1 版本
	v1 := r.Group("/api/v1")
	{
		// GET /api/v1/users
		v1.GET("/users", listUsersV1)

		// 嵌套路由组
		userGroup := v1.Group("/users")
		{
			// GET /api/v1/users/:id
			userGroup.GET("/:id", getUserV1)
			// POST /api/v1/users
			userGroup.POST("", createUserV1)
			// PUT /api/v1/users/:id
			userGroup.PUT("/:id", updateUserV1)
			// DELETE /api/v1/users/:id
			userGroup.DELETE("/:id", deleteUserV1)
		}

		// 订单相关路由
		orderGroup := v1.Group("/orders")
		{
			orderGroup.GET("", listOrders)
			orderGroup.GET("/:id", getOrder)
			orderGroup.POST("", createOrder)
		}
	}

	// API v2 版本 (演示版本演进)
	v2 := r.Group("/api/v2")
	{
		v2.GET("/users", listUsersV2) // 新版本可能返回不同的数据结构
	}

	// ========================================================================
	// 五、路由组 + 中间件
	// ========================================================================

	// 需要认证的路由组
	authorized := r.Group("/admin")
	authorized.Use(authMiddleware()) // 组级别中间件
	{
		authorized.GET("/dashboard", adminDashboard)
		authorized.GET("/settings", adminSettings)
	}

	// 公开路由（无需认证）
	public := r.Group("/public")
	{
		public.GET("/info", publicInfo)
	}

	// ========================================================================
	// 六、静态文件服务
	// ========================================================================

	// 单个静态文件
	// r.StaticFile("/favicon.ico", "./resources/favicon.ico")

	// 静态目录
	// r.Static("/static", "./static")

	// 静态文件系统 (嵌入式文件等)
	// r.StaticFS("/assets", http.Dir("assets"))

	// ========================================================================
	// 七、路由冲突与优先级
	// ========================================================================
	//
	// 【重要】Gin 路由匹配优先级：
	//
	// 1. 精确匹配 > 参数匹配 > 通配符匹配
	//
	// 示例：
	// /users/new     <- 精确匹配
	// /users/:id     <- 参数匹配
	// /users/*path   <- 通配符匹配
	//
	// 请求 /users/new 会匹配到第一个
	// 请求 /users/123 会匹配到第二个
	//
	// 【易错点】路由注册顺序
	//

	// 正确顺序：精确路由在前
	r.GET("/items/new", createItemForm)  // 精确匹配
	r.GET("/items/:id", getItem)         // 参数匹配

	// ========================================================================
	// 八、获取完整 URL 信息
	// ========================================================================

	r.GET("/debug/request", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			// 请求路径
			"path": c.Request.URL.Path, // /debug/request
			// 完整 URL
			"full_url": c.Request.URL.String(), // /debug/request?foo=bar
			// 查询字符串
			"query_string": c.Request.URL.RawQuery, // foo=bar
			// 请求方法
			"method": c.Request.Method, // GET
			// 客户端 IP
			"client_ip": c.ClientIP(),
			// 匹配的路由模式
			"full_path": c.FullPath(), // /debug/request
		})
	})

	// ========================================================================
	// 九、NoRoute 和 NoMethod 处理
	// ========================================================================

	// 404 - 路由不存在
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "not_found",
			"message": "The requested resource was not found",
			"path":    c.Request.URL.Path,
		})
	})

	// 405 - 方法不允许
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error":   "method_not_allowed",
			"message": "The requested method is not allowed",
			"method":  c.Request.Method,
		})
	})

	// 启动服务器
	r.Run(":8080")
}

// ============================================================================
// Handler 函数定义
// ============================================================================

func handleGet(c *gin.Context)     { c.String(http.StatusOK, "GET") }
func handlePost(c *gin.Context)    { c.String(http.StatusOK, "POST") }
func handlePut(c *gin.Context)     { c.String(http.StatusOK, "PUT") }
func handlePatch(c *gin.Context)   { c.String(http.StatusOK, "PATCH") }
func handleDelete(c *gin.Context)  { c.String(http.StatusOK, "DELETE") }
func handleOptions(c *gin.Context) { c.String(http.StatusOK, "OPTIONS") }
func handleHead(c *gin.Context)    { c.Status(http.StatusOK) }

// V1 版本的用户接口
func listUsersV1(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": "v1",
		"users":   []string{"user1", "user2"},
	})
}

func getUserV1(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"version": "v1", "id": id})
}

func createUserV1(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"version": "v1", "action": "created"})
}

func updateUserV1(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"version": "v1", "id": id, "action": "updated"})
}

func deleteUserV1(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"version": "v1", "id": id, "action": "deleted"})
}

// V2 版本（演示版本演进）
func listUsersV2(c *gin.Context) {
	// V2 返回更详细的数据结构
	c.JSON(http.StatusOK, gin.H{
		"version": "v2",
		"data": []gin.H{
			{"id": 1, "name": "user1", "email": "user1@example.com"},
			{"id": 2, "name": "user2", "email": "user2@example.com"},
		},
		"meta": gin.H{
			"total": 2,
			"page":  1,
		},
	})
}

// 订单相关
func listOrders(c *gin.Context)  { c.JSON(http.StatusOK, gin.H{"orders": []string{}}) }
func getOrder(c *gin.Context)    { c.JSON(http.StatusOK, gin.H{"order": c.Param("id")}) }
func createOrder(c *gin.Context) { c.JSON(http.StatusCreated, gin.H{"action": "created"}) }

// 管理后台
func adminDashboard(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"page": "dashboard"}) }
func adminSettings(c *gin.Context)  { c.JSON(http.StatusOK, gin.H{"page": "settings"}) }
func publicInfo(c *gin.Context)     { c.JSON(http.StatusOK, gin.H{"info": "public"}) }

// 商品
func createItemForm(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"form": "new item"}) }
func getItem(c *gin.Context)        { c.JSON(http.StatusOK, gin.H{"item": c.Param("id")}) }

// 简单的认证中间件（示例）
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}
		c.Next()
	}
}

// ============================================================================
// 测试命令
// ============================================================================
//
// # 基础路由测试
// curl http://localhost:8080/get
// curl -X POST http://localhost:8080/post
//
// # 路径参数测试
// curl http://localhost:8080/users/123
// curl http://localhost:8080/users/123/posts/456
// curl http://localhost:8080/files/path/to/file.txt
//
// # 查询参数测试
// curl "http://localhost:8080/search?keyword=gin&page=2&limit=20&tags=go&tags=web"
// curl "http://localhost:8080/search?keyword=test&filter[status]=active&filter[type]=premium"
//
// # 路由组测试
// curl http://localhost:8080/api/v1/users
// curl http://localhost:8080/api/v1/users/123
// curl http://localhost:8080/api/v2/users
//
// # 认证测试
// curl http://localhost:8080/admin/dashboard  # 401
// curl -H "Authorization: Bearer token" http://localhost:8080/admin/dashboard  # 200
//
// # 404/405 测试
// curl http://localhost:8080/not-exist
// curl -X DELETE http://localhost:8080/get
//
// ============================================================================

// ============================================================================
// 易错点总结
// ============================================================================
//
// 1. 【路径参数 vs 查询参数】
//    - /users/:id -> c.Param("id")      路径参数
//    - /users?id=123 -> c.Query("id")   查询参数
//
// 2. 【路由顺序】
//    精确路由必须在参数路由之前注册
//    r.GET("/users/new", ...)  // 先
//    r.GET("/users/:id", ...)  // 后
//
// 3. 【通配符位置】
//    *param 必须在路径末尾
//    正确: /files/*filepath
//    错误: /files/*filepath/info
//
// 4. 【路由组不要遗漏斜杠】
//    v1 := r.Group("/api/v1")  // 正确
//    v1 := r.Group("api/v1")   // 可能导致路径错误
//
// 5. 【Param 返回值包含前导斜杠】
//    /files/*filepath 匹配 /files/a/b/c
//    c.Param("filepath") 返回 "/a/b/c" (注意有斜杠)
//
// ============================================================================

// ============================================================================
// 练习题
// ============================================================================
//
// 1. 设计一个博客系统的 RESTful API 路由：
//    - 文章 CRUD: /api/v1/posts
//    - 文章评论: /api/v1/posts/:id/comments
//    - 标签管理: /api/v1/tags
//
// 2. 实现一个路由组，/admin 下的所有路由需要认证，/public 下的不需要
//
// 3. 使用 QueryMap 实现一个高级搜索接口：
//    /search?filter[category]=tech&filter[author]=john&sort=-created_at
//
// ============================================================================
