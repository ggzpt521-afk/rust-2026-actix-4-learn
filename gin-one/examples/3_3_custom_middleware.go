// ============================================================================
// 3.3 自定义中间件实战
// ============================================================================
// 运行方式: go run examples/3_3_custom_middleware.go
// ============================================================================

package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// ============================================================================
// 本节实现以下生产级中间件：
//
// 1. CORS 跨域中间件
// 2. 请求 ID 中间件
// 3. 耗时统计中间件
// 4. 限流中间件（令牌桶）
// 5. JWT 认证中间件（简化版）
// 6. 请求日志中间件
//
// ============================================================================

func main() {
	r := gin.New()

	// ========================================================================
	// 全局中间件
	// ========================================================================

	// 1. Recovery (必须)
	r.Use(gin.Recovery())

	// 2. 请求 ID (方便追踪)
	r.Use(RequestIDMiddleware())

	// 3. CORS 跨域
	r.Use(CORSMiddleware())

	// 4. 请求日志
	r.Use(AccessLogMiddleware())

	// 5. 耗时统计
	r.Use(TimingMiddleware())

	// ========================================================================
	// 公开路由
	// ========================================================================

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":    "pong",
			"request_id": c.GetString("request_id"),
		})
	})

	// 登录接口（返回 token）
	r.POST("/login", func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 简化：假设任何用户名/密码都有效
		token := generateToken()

		c.JSON(http.StatusOK, gin.H{
			"token":      token,
			"token_type": "Bearer",
			"expires_in": 3600,
		})
	})

	// ========================================================================
	// 需要认证的路由
	// ========================================================================

	authorized := r.Group("/api")
	authorized.Use(JWTAuthMiddleware())
	{
		authorized.GET("/profile", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"user_id":    c.GetInt("user_id"),
				"username":   c.GetString("username"),
				"request_id": c.GetString("request_id"),
			})
		})

		authorized.GET("/settings", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"settings": "user settings"})
		})
	}

	// ========================================================================
	// 限流测试路由
	// ========================================================================

	// 创建限流器：每秒 5 个请求，最大突发 10 个
	rateLimiter := NewRateLimiter(5, 10)

	limited := r.Group("/limited")
	limited.Use(RateLimitMiddleware(rateLimiter))
	{
		limited.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "rate limit test"})
		})
	}

	// ========================================================================
	// 特定路由中间件
	// ========================================================================

	// 慢接口 - 演示耗时统计
	r.GET("/slow", func(c *gin.Context) {
		time.Sleep(200 * time.Millisecond)
		c.JSON(http.StatusOK, gin.H{"message": "slow response"})
	})

	log.Println("Server starting on :8080")
	log.Println("Test URLs:")
	log.Println("  curl http://localhost:8080/ping")
	log.Println("  curl -X POST http://localhost:8080/login -H 'Content-Type: application/json' -d '{\"username\":\"test\",\"password\":\"123\"}'")
	log.Println("  curl -H 'Authorization: Bearer <token>' http://localhost:8080/api/profile")
	log.Println("  for i in {1..15}; do curl http://localhost:8080/limited/test; done")

	r.Run(":8080")
}

// ============================================================================
// 1. CORS 跨域中间件
// ============================================================================
//
// 【什么是 CORS？】
//
// Cross-Origin Resource Sharing (跨域资源共享)
// 浏览器出于安全考虑，限制从脚本发起的跨域 HTTP 请求
//
// 【CORS 工作流程】
//
// 1. 简单请求：直接发送，服务器响应中带 CORS 头
// 2. 预检请求：先发 OPTIONS 请求，通过后再发实际请求
//
// 【预检请求条件】（满足任一）
// - 使用 PUT/DELETE/PATCH 等方法
// - Content-Type 不是 form-urlencoded/multipart/text-plain
// - 有自定义 Header
//

// CORSMiddleware CORS 中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// 允许的域名列表（生产环境应该配置具体域名）
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"https://example.com",
		}

		// 检查是否允许该 Origin
		allowed := false
		for _, o := range allowedOrigins {
			if o == origin || o == "*" {
				allowed = true
				break
			}
		}

		if allowed {
			// 允许的源
			c.Header("Access-Control-Allow-Origin", origin)
			// 允许的方法
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			// 允许的请求头
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")
			// 允许携带凭证（Cookie）
			c.Header("Access-Control-Allow-Credentials", "true")
			// 预检请求缓存时间
			c.Header("Access-Control-Max-Age", "86400")
			// 暴露给前端的响应头
			c.Header("Access-Control-Expose-Headers", "X-Request-ID, X-Response-Time")
		}

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// ============================================================================
// 2. 请求 ID 中间件
// ============================================================================
//
// 【作用】
//
// 为每个请求生成唯一 ID，用于：
// 1. 日志追踪
// 2. 错误排查
// 3. 分布式链路追踪
//

// RequestIDMiddleware 请求 ID 中间件
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 优先从请求头获取（支持分布式追踪）
		requestID := c.GetHeader("X-Request-ID")

		// 如果没有则生成新的
		if requestID == "" {
			requestID = generateRequestID()
		}

		// 设置到 Context 供后续使用
		c.Set("request_id", requestID)

		// 设置响应头
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// generateRequestID 生成请求 ID
func generateRequestID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// ============================================================================
// 3. 耗时统计中间件
// ============================================================================

// TimingMiddleware 耗时统计中间件
func TimingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		// 设置响应头
		c.Header("X-Response-Time", duration.String())

		// 慢请求告警（超过 1 秒）
		if duration > time.Second {
			log.Printf("[SLOW] %s %s took %v", c.Request.Method, c.Request.URL.Path, duration)
		}
	}
}

// ============================================================================
// 4. 限流中间件（令牌桶算法）
// ============================================================================
//
// 【令牌桶算法】
//
// 1. 桶中有固定数量的令牌
// 2. 每个请求需要获取一个令牌
// 3. 令牌以固定速率补充
// 4. 桶满时，新令牌被丢弃
//
// 优点：允许一定程度的突发流量
//

// RateLimiter 限流器
type RateLimiter struct {
	rate       float64    // 令牌生成速率（每秒）
	maxTokens  float64    // 桶容量
	tokens     float64    // 当前令牌数
	lastUpdate time.Time  // 上次更新时间
	mu         sync.Mutex // 互斥锁
}

// NewRateLimiter 创建限流器
func NewRateLimiter(rate float64, maxTokens float64) *RateLimiter {
	return &RateLimiter{
		rate:       rate,
		maxTokens:  maxTokens,
		tokens:     maxTokens, // 初始满桶
		lastUpdate: time.Now(),
	}
}

// Allow 尝试获取令牌
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastUpdate).Seconds()

	// 补充令牌
	rl.tokens += elapsed * rl.rate
	if rl.tokens > rl.maxTokens {
		rl.tokens = rl.maxTokens
	}
	rl.lastUpdate = now

	// 尝试获取令牌
	if rl.tokens >= 1 {
		rl.tokens--
		return true
	}

	return false
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":      "rate_limit_exceeded",
				"message":    "请求过于频繁，请稍后重试",
				"retry_after": 1,
			})
			return
		}

		c.Next()
	}
}

// ============================================================================
// 5. JWT 认证中间件（简化版）
// ============================================================================
//
// 【实际生产环境应该】
//
// 1. 使用 github.com/golang-jwt/jwt 库
// 2. 实现 token 刷新机制
// 3. 实现 token 黑名单
// 4. 存储用户信息到 Redis
//

// JWTAuthMiddleware JWT 认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization Header
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "missing_token",
				"message": "Authorization header is required",
			})
			return
		}

		// 解析 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_token_format",
				"message": "Authorization header must be: Bearer <token>",
			})
			return
		}

		token := parts[1]

		// 验证 Token（简化版，实际应该用 JWT 库）
		if len(token) < 20 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_token",
				"message": "Token is invalid or expired",
			})
			return
		}

		// 模拟解析用户信息
		// 实际应该从 JWT claims 中获取
		c.Set("user_id", 1001)
		c.Set("username", "test_user")

		c.Next()
	}
}

// generateToken 生成 Token（简化版）
func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// ============================================================================
// 6. 请求日志中间件
// ============================================================================

// AccessLogMiddleware 访问日志中间件
func AccessLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 处理请求
		c.Next()

		// 记录日志
		log.Printf("[ACCESS] %s | %3d | %13v | %15s | %-7s %s | %s",
			start.Format("2006-01-02 15:04:05"),
			c.Writer.Status(),
			time.Since(start),
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.Path,
			c.GetString("request_id"),
		)
	}
}

// ============================================================================
// 测试命令
// ============================================================================
//
// # 基本测试
// curl http://localhost:8080/ping
//
// # 登录获取 token
// curl -X POST http://localhost:8080/login \
//   -H "Content-Type: application/json" \
//   -d '{"username":"test","password":"123"}'
//
// # 使用 token 访问受保护接口
// curl -H "Authorization: Bearer <token>" http://localhost:8080/api/profile
//
// # CORS 预检请求测试
// curl -X OPTIONS http://localhost:8080/ping \
//   -H "Origin: http://localhost:3000" \
//   -H "Access-Control-Request-Method: POST" \
//   -v
//
// # 限流测试（快速发送 15 个请求）
// for i in {1..15}; do curl http://localhost:8080/limited/test; echo; done
//
// # 慢接口测试
// curl http://localhost:8080/slow -v
//
// ============================================================================

// ============================================================================
// 易错点总结
// ============================================================================
//
// 1. 【CORS Allow-Origin 不能同时用 * 和 Credentials】
//    如果 Allow-Credentials=true，Origin 必须是具体值
//
// 2. 【限流器并发安全】
//    多个请求同时访问限流器，需要加锁
//    或使用 golang.org/x/time/rate 库
//
// 3. 【JWT Token 不要放 URL】
//    不要 /api/users?token=xxx
//    应该放 Authorization Header
//
// 4. 【请求 ID 要透传】
//    调用下游服务时，要把 Request ID 传递下去
//    实现分布式链路追踪
//
// 5. 【中间件顺序很重要】
//    Recovery 应该最先（捕获所有 panic）
//    RequestID 要在日志之前（日志需要用）
//    认证中间件要在业务路由之前
//
// ============================================================================

// ============================================================================
// 练习题
// ============================================================================
//
// 1. 实现 IP 黑白名单中间件:
//    - 白名单模式：只允许指定 IP
//    - 黑名单模式：禁止指定 IP
//
// 2. 实现滑动窗口限流:
//    - 每分钟最多 100 次请求
//    - 统计最近 60 秒内的请求次数
//
// 3. 实现请求签名验证中间件:
//    - 请求带 timestamp 和 signature
//    - signature = MD5(timestamp + secret_key)
//    - timestamp 超过 5 分钟则拒绝
//
// ============================================================================
