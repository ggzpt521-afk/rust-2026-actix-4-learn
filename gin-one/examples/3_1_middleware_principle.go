// ============================================================================
// 3.1 中间件原理
// ============================================================================
// 运行方式: go run examples/3_1_middleware_principle.go
// ============================================================================

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ============================================================================
// 核心概念：什么是中间件？
// ============================================================================
//
// 中间件 (Middleware) 是在请求到达最终 Handler 之前/之后执行的代码
// 常见用途：日志、认证、限流、跨域、压缩、错误处理等
//
// 【洋葱模型 (Onion Model)】
//
//                     请求进入
//                        │
//                        ▼
//    ┌─────────────────────────────────────────┐
//    │           Middleware 1 (进入)            │
//    │    ┌─────────────────────────────────┐  │
//    │    │       Middleware 2 (进入)        │  │
//    │    │    ┌─────────────────────────┐  │  │
//    │    │    │    Middleware 3 (进入)   │  │  │
//    │    │    │    ┌─────────────────┐  │  │  │
//    │    │    │    │                 │  │  │  │
//    │    │    │    │    Handler      │  │  │  │
//    │    │    │    │                 │  │  │  │
//    │    │    │    └─────────────────┘  │  │  │
//    │    │    │    Middleware 3 (返回)   │  │  │
//    │    │    └─────────────────────────┘  │  │
//    │    │       Middleware 2 (返回)        │  │
//    │    └─────────────────────────────────┘  │
//    │           Middleware 1 (返回)            │
//    └─────────────────────────────────────────┘
//                        │
//                        ▼
//                     响应返回
//
// ============================================================================
//
// 【执行流程】
//
// 1. 请求进入 → Middleware1 前半部分
// 2. 调用 c.Next() → 进入 Middleware2 前半部分
// 3. 调用 c.Next() → 进入 Middleware3 前半部分
// 4. 调用 c.Next() → 执行 Handler
// 5. Handler 返回 → Middleware3 后半部分
// 6. 返回 → Middleware2 后半部分
// 7. 返回 → Middleware1 后半部分
// 8. 响应发送给客户端
//
// ============================================================================

// ============================================================================
// Gin 中间件的本质
// ============================================================================
//
// type HandlerFunc func(*Context)
//
// 中间件就是一个 HandlerFunc，和普通 Handler 没有本质区别
// 区别在于：
// 1. 中间件通常会调用 c.Next() 继续后续处理
// 2. 中间件通过 c.Use() 注册，会对多个路由生效
//
// 【c.Next() 的作用】
//
// c.Next() 会暂停当前中间件，执行后续的中间件和 Handler
// 当后续全部执行完毕后，继续执行当前中间件 c.Next() 之后的代码
//
// 【c.Abort() 的作用】
//
// c.Abort() 会阻止后续中间件和 Handler 的执行
// 但当前中间件的剩余代码仍会执行（除非 return）
//
// ============================================================================

func main() {
	// 使用 gin.New() 创建空 Engine，不带任何中间件
	r := gin.New()

	// ========================================================================
	// 一、演示中间件执行顺序
	// ========================================================================

	// 注册全局中间件（按顺序执行）
	r.Use(middleware1())
	r.Use(middleware2())
	r.Use(middleware3())

	// 这个路由会经过所有全局中间件
	r.GET("/order", func(c *gin.Context) {
		log.Println(">>> Handler 执行")
		c.JSON(http.StatusOK, gin.H{
			"message": "handler executed",
		})
	})

	// ========================================================================
	// 二、演示 c.Abort() 的效果
	// ========================================================================

	r.GET("/abort", abortMiddleware(), func(c *gin.Context) {
		// 这个 Handler 不会执行
		log.Println(">>> Handler 执行 (这条不会打印)")
		c.JSON(http.StatusOK, gin.H{"message": "handler"})
	})

	// ========================================================================
	// 三、演示 c.Set() 和 c.Get() 传递数据
	// ========================================================================

	r.GET("/data", setDataMiddleware(), func(c *gin.Context) {
		// 获取中间件设置的数据
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found"})
			return
		}

		username, _ := c.Get("username")
		requestTime, _ := c.Get("request_time")

		c.JSON(http.StatusOK, gin.H{
			"user_id":      userID,
			"username":     username,
			"request_time": requestTime,
		})
	})

	// ========================================================================
	// 四、演示耗时统计中间件
	// ========================================================================

	r.GET("/slow", timingMiddleware(), func(c *gin.Context) {
		// 模拟慢接口
		time.Sleep(100 * time.Millisecond)
		c.JSON(http.StatusOK, gin.H{"message": "slow response"})
	})

	// ========================================================================
	// 五、路由组中间件
	// ========================================================================

	// 公开路由组（无需认证）
	public := r.Group("/public")
	{
		public.GET("/info", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "public info"})
		})
	}

	// 需要认证的路由组
	authorized := r.Group("/api")
	authorized.Use(authMiddleware())
	{
		authorized.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			c.JSON(http.StatusOK, gin.H{
				"message": "protected resource",
				"user_id": userID,
			})
		})

		authorized.GET("/settings", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "user settings"})
		})
	}

	// ========================================================================
	// 六、单个路由的中间件
	// ========================================================================

	// 只有这个路由使用 adminMiddleware
	r.GET("/admin/dashboard",
		authMiddleware(),    // 先认证
		adminMiddleware(),   // 再检查管理员权限
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin dashboard"})
		},
	)

	// ========================================================================
	// 七、中间件执行顺序可视化
	// ========================================================================

	r.GET("/visual", visualMiddleware("A"), visualMiddleware("B"), visualMiddleware("C"),
		func(c *gin.Context) {
			log.Println(">>> Handler")
			c.String(http.StatusOK, "OK")
		},
	)

	log.Println("Server starting on :8080")
	log.Println("Test URLs:")
	log.Println("  curl http://localhost:8080/order")
	log.Println("  curl http://localhost:8080/abort")
	log.Println("  curl http://localhost:8080/data")
	log.Println("  curl http://localhost:8080/slow")
	log.Println("  curl http://localhost:8080/visual")
	log.Println("  curl http://localhost:8080/public/info")
	log.Println("  curl -H 'Authorization: Bearer token' http://localhost:8080/api/profile")

	r.Run(":8080")
}

// ============================================================================
// 中间件实现
// ============================================================================

// middleware1 第一个中间件
func middleware1() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("==> Middleware 1: 进入")

		c.Next() // 执行后续中间件和 Handler

		log.Println("<== Middleware 1: 返回")
	}
}

// middleware2 第二个中间件
func middleware2() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("==> Middleware 2: 进入")

		c.Next()

		log.Println("<== Middleware 2: 返回")
	}
}

// middleware3 第三个中间件
func middleware3() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("==> Middleware 3: 进入")

		c.Next()

		log.Println("<== Middleware 3: 返回")
	}
}

// abortMiddleware 演示 Abort
func abortMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("==> abortMiddleware: 准备 Abort")

		// Abort 阻止后续执行，但当前函数的代码会继续
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "access denied",
		})

		// 这行仍然会执行！
		log.Println("==> abortMiddleware: Abort 之后的代码仍执行")

		// 如果不想执行后续代码，需要 return
		// return

		// 这行也会执行
		log.Println("==> abortMiddleware: 函数结束")
	}
}

// setDataMiddleware 演示数据传递
func setDataMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置数据，Handler 和后续中间件可以读取
		c.Set("user_id", 12345)
		c.Set("username", "张三")
		c.Set("request_time", time.Now())

		c.Next()
	}
}

// timingMiddleware 耗时统计中间件
func timingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		// 计算耗时
		duration := time.Since(start)

		// 设置响应头
		c.Header("X-Response-Time", duration.String())

		log.Printf("[TIMING] %s %s - %v", c.Request.Method, c.Request.URL.Path, duration)
	}
}

// authMiddleware 认证中间件
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Authorization header is required",
			})
			return
		}

		// 简化的 token 验证（实际应该验证 JWT）
		if len(token) < 10 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_token",
				"message": "Invalid token format",
			})
			return
		}

		// 设置用户信息供后续使用
		c.Set("user_id", 1001)
		c.Set("username", "authenticated_user")

		c.Next()
	}
}

// adminMiddleware 管理员权限中间件
func adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "user not authenticated",
			})
			return
		}

		// 简化的权限检查（实际应该查数据库）
		if userID != 1001 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "admin access required",
			})
			return
		}

		c.Next()
	}
}

// visualMiddleware 可视化执行顺序
func visualMiddleware(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("==> [%s] 进入", name)
		c.Next()
		log.Printf("<== [%s] 返回", name)
	}
}

// ============================================================================
// c.Next() 源码解析
// ============================================================================
//
// func (c *Context) Next() {
//     c.index++
//     for c.index < int8(len(c.handlers)) {
//         c.handlers[c.index](c)
//         c.index++
//     }
// }
//
// 关键点：
// 1. c.index 记录当前执行到第几个 handler
// 2. c.handlers 是所有中间件和 Handler 的数组
// 3. Next() 会递增 index 并执行下一个 handler
//
// ============================================================================
//
// c.Abort() 源码解析
//
// const abortIndex int8 = math.MaxInt8 >> 1  // 63
//
// func (c *Context) Abort() {
//     c.index = abortIndex
// }
//
// 关键点：
// 1. Abort 将 index 设为一个很大的值
// 2. 这导致 Next() 的循环条件不满足，停止执行后续 handler
//
// ============================================================================

// ============================================================================
// 测试输出示例
// ============================================================================
//
// 访问 /order 会输出：
//
// ==> Middleware 1: 进入
// ==> Middleware 2: 进入
// ==> Middleware 3: 进入
// >>> Handler 执行
// <== Middleware 3: 返回
// <== Middleware 2: 返回
// <== Middleware 1: 返回
//
// 访问 /visual 会输出：
//
// ==> [A] 进入
// ==> [B] 进入
// ==> [C] 进入
// >>> Handler
// <== [C] 返回
// <== [B] 返回
// <== [A] 返回
//
// ============================================================================

// ============================================================================
// 易错点总结
// ============================================================================
//
// 1. 【Abort 后代码仍执行】
//    c.Abort() 只阻止后续 handler，当前函数代码继续执行
//    需要配合 return 使用：
//    c.AbortWithStatusJSON(...)
//    return
//
// 2. 【中间件注册顺序很重要】
//    先注册的中间件先执行（进入时）
//    后注册的中间件先返回（返回时）
//
// 3. 【c.Set/c.Get 的类型】
//    c.Get 返回 interface{}，需要类型断言
//    推荐用 c.GetString, c.GetInt 等便捷方法
//
// 4. 【不要在中间件外使用 context】
//    中间件执行完后，context 可能被回收
//    如需异步使用，应该 copy := c.Copy()
//
// 5. 【响应只能写一次】
//    如果中间件已经写了响应，Handler 再写会报错
//    使用 c.Writer.Written() 检查是否已写入
//
// ============================================================================

// ============================================================================
// 练习题
// ============================================================================
//
// 1. 实现一个请求 ID 中间件:
//    - 从 Header 获取 X-Request-ID，如果没有则生成
//    - 将 Request ID 设置到 context 和响应 Header
//
// 2. 实现一个限流中间件:
//    - 每秒最多处理 10 个请求
//    - 超出限制返回 429 Too Many Requests
//
// 3. 实现一个黑名单中间件:
//    - 维护一个 IP 黑名单
//    - 黑名单中的 IP 返回 403 Forbidden
//
// ============================================================================
