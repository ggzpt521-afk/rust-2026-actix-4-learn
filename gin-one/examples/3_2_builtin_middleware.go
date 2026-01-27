// ============================================================================
// 3.2 内置中间件详解
// ============================================================================
// 运行方式: go run examples/3_2_builtin_middleware.go
// ============================================================================

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// ============================================================================
// Gin 内置中间件
// ============================================================================
//
// Gin 提供了两个核心内置中间件：
//
// 1. gin.Logger() - 请求日志中间件
// 2. gin.Recovery() - panic 恢复中间件
//
// gin.Default() = gin.New() + Logger() + Recovery()
//
// ============================================================================

// ============================================================================
// 一、Logger 中间件详解
// ============================================================================
//
// 【默认日志格式】
//
// [GIN] 2024/01/15 - 10:30:45 | 200 |     1.234ms |       127.0.0.1 | GET      "/ping"
//
// 包含信息：
// - 时间戳
// - HTTP 状态码
// - 响应耗时
// - 客户端 IP
// - HTTP 方法
// - 请求路径
//
// 【源码简化版】
//
// func Logger() HandlerFunc {
//     return LoggerWithConfig(LoggerConfig{})
// }
//
// type LoggerConfig struct {
//     Formatter LogFormatter        // 日志格式化函数
//     Output    io.Writer           // 输出目标
//     SkipPaths []string            // 跳过的路径
// }
//
// ============================================================================

// ============================================================================
// 二、Recovery 中间件详解
// ============================================================================
//
// 【作用】
//
// 捕获 Handler 中的 panic，防止整个服务崩溃
// 返回 500 Internal Server Error
//
// 【源码简化版】
//
// func Recovery() HandlerFunc {
//     return RecoveryWithWriter(DefaultErrorWriter)
// }
//
// func CustomRecovery(handle RecoveryFunc) HandlerFunc {
//     return RecoveryWithWriter(DefaultErrorWriter, handle)
// }
//
// func RecoveryWithWriter(out io.Writer, recovery ...RecoveryFunc) HandlerFunc {
//     return func(c *Context) {
//         defer func() {
//             if err := recover(); err != nil {
//                 // 打印堆栈
//                 // 返回 500
//                 // 调用自定义 recovery 函数（如果有）
//             }
//         }()
//         c.Next()
//     }
// }
//
// ============================================================================

func main() {
	// ========================================================================
	// 一、使用默认 Logger
	// ========================================================================

	// 方式一：gin.Default() 自动包含 Logger 和 Recovery
	// r := gin.Default()

	// 方式二：手动添加（推荐，可控制）
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// ========================================================================
	// 二、自定义 Logger 配置
	// ========================================================================

	customLogger := gin.LoggerWithConfig(gin.LoggerConfig{
		// 自定义日志格式
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("[GIN] %s | %3d | %13v | %15s | %-7s %s\n%s",
				param.TimeStamp.Format("2006-01-02 15:04:05"),
				param.StatusCode,
				param.Latency,
				param.ClientIP,
				param.Method,
				param.Path,
				param.ErrorMessage,
			)
		},
		// 输出到标准输出
		Output: os.Stdout,
		// 跳过某些路径（如健康检查）
		SkipPaths: []string{"/health", "/metrics"},
	})

	// 使用自定义 Logger 的路由组
	customGroup := r.Group("/custom")
	customGroup.Use(customLogger)
	{
		customGroup.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "custom logger"})
		})
	}

	// ========================================================================
	// 三、Logger 输出到文件
	// ========================================================================

	// 创建日志文件
	logFile, err := os.OpenFile("gin.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// 同时输出到文件和控制台
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)

	// ========================================================================
	// 四、自定义 Recovery
	// ========================================================================

	customRecovery := gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// 记录错误日志
		log.Printf("[PANIC] %v", recovered)

		// 可以在这里：
		// 1. 发送告警通知（钉钉、企业微信等）
		// 2. 上报错误到监控系统（Sentry、ELK 等）
		// 3. 记录详细的错误信息

		// 返回友好的错误响应
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务器内部错误，请稍后重试",
			"error":   fmt.Sprintf("%v", recovered),
		})
	})

	// 使用自定义 Recovery 的路由组
	panicGroup := r.Group("/panic")
	panicGroup.Use(customRecovery)
	{
		// 这个接口会 panic，但不会导致服务崩溃
		panicGroup.GET("/test", func(c *gin.Context) {
			panic("something went wrong!")
		})

		// 数组越界 panic
		panicGroup.GET("/array", func(c *gin.Context) {
			arr := []int{1, 2, 3}
			_ = arr[10] // 越界
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		// 空指针 panic
		panicGroup.GET("/nil", func(c *gin.Context) {
			var m map[string]string
			m["key"] = "value" // nil map 写入
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})
	}

	// ========================================================================
	// 五、关闭日志颜色（生产环境）
	// ========================================================================

	// 生产环境通常将日志写入文件或 ELK，不需要颜色
	// gin.DisableConsoleColor()

	// 开发环境可以强制开启颜色
	// gin.ForceConsoleColor()

	// ========================================================================
	// 六、访问日志中间件（生产级）
	// ========================================================================

	productionLogger := ProductionLogger()
	prodGroup := r.Group("/prod")
	prodGroup.Use(productionLogger)
	{
		prodGroup.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "production log test"})
		})
	}

	// ========================================================================
	// 测试路由
	// ========================================================================

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// 健康检查（会被 skipPaths 跳过日志）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// 不同状态码测试
	r.GET("/status/:code", func(c *gin.Context) {
		code := 200
		fmt.Sscanf(c.Param("code"), "%d", &code)
		c.JSON(code, gin.H{"code": code})
	})

	log.Println("Server starting on :8080")
	log.Println("Test URLs:")
	log.Println("  curl http://localhost:8080/ping")
	log.Println("  curl http://localhost:8080/health")
	log.Println("  curl http://localhost:8080/panic/test")
	log.Println("  curl http://localhost:8080/status/404")
	log.Println("  curl http://localhost:8080/prod/test")

	r.Run(":8080")
}

// ============================================================================
// 生产级日志中间件
// ============================================================================

// ProductionLogger 生产环境日志中间件
func ProductionLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()

		// 请求路径
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()
		latency := end.Sub(start)

		// 获取状态码
		statusCode := c.Writer.Status()

		// 获取错误信息
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// JSON 格式日志（方便 ELK 解析）
		logEntry := map[string]interface{}{
			"timestamp":   end.Format(time.RFC3339),
			"status":      statusCode,
			"latency":     latency.String(),
			"latency_ms":  latency.Milliseconds(),
			"client_ip":   c.ClientIP(),
			"method":      c.Request.Method,
			"path":        path,
			"query":       query,
			"user_agent":  c.Request.UserAgent(),
			"error":       errorMessage,
			"body_size":   c.Writer.Size(),
			"request_id":  c.GetHeader("X-Request-ID"),
		}

		// 根据状态码选择日志级别
		if statusCode >= 500 {
			log.Printf("[ERROR] %v", logEntry)
		} else if statusCode >= 400 {
			log.Printf("[WARN] %v", logEntry)
		} else {
			log.Printf("[INFO] %v", logEntry)
		}
	}
}

// ============================================================================
// Recovery 工作原理演示
// ============================================================================
//
// func RecoveryMiddleware() gin.HandlerFunc {
//     return func(c *gin.Context) {
//         defer func() {
//             if r := recover(); r != nil {
//                 // 1. 获取堆栈信息
//                 stack := debug.Stack()
//
//                 // 2. 记录日志
//                 log.Printf("Panic recovered: %v\n%s", r, stack)
//
//                 // 3. 返回错误响应
//                 c.AbortWithStatusJSON(500, gin.H{
//                     "error": "Internal Server Error",
//                 })
//             }
//         }()
//
//         c.Next()
//     }
// }
//
// ============================================================================

// ============================================================================
// 测试命令
// ============================================================================
//
// # 正常请求
// curl http://localhost:8080/ping
//
// # 触发 panic
// curl http://localhost:8080/panic/test
// curl http://localhost:8080/panic/array
// curl http://localhost:8080/panic/nil
//
// # 不同状态码
// curl http://localhost:8080/status/200
// curl http://localhost:8080/status/404
// curl http://localhost:8080/status/500
//
// # 自定义日志
// curl http://localhost:8080/custom/test
//
// # 生产级日志
// curl http://localhost:8080/prod/test
//
// ============================================================================

// ============================================================================
// 易错点总结
// ============================================================================
//
// 1. 【Logger 不记录 Body】
//    默认 Logger 不记录请求/响应 Body
//    如需记录，需要自己实现（注意大 Body 问题）
//
// 2. 【Recovery 只捕获当前 goroutine】
//    如果 Handler 中启动了新的 goroutine，其中的 panic 不会被捕获
//    解决：在新 goroutine 中也加 defer recover
//
// 3. 【日志文件需要轮转】
//    直接写文件会导致文件无限增大
//    使用 lumberjack 等库实现日志轮转
//
// 4. 【生产环境关闭颜色】
//    颜色代码写入文件会产生乱码
//    gin.DisableConsoleColor()
//
// 5. 【SkipPaths 精确匹配】
//    SkipPaths 是精确匹配，不支持通配符
//    /health 和 /health/ 是不同的路径
//
// ============================================================================

// ============================================================================
// 练习题
// ============================================================================
//
// 1. 实现一个记录请求 Body 的日志中间件
//    （提示：需要读取 Body 后再写回）
//
// 2. 实现一个 Recovery 中间件，在 panic 时：
//    - 记录详细堆栈信息到文件
//    - 发送告警通知（模拟）
//    - 返回 request_id 供排查
//
// 3. 实现日志采样：只记录 10% 的正常请求日志
//
// ============================================================================
