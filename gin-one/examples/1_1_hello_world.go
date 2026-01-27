// ============================================================================
// 1.1 Gin 安装与 Hello World
// ============================================================================
// 运行方式: go run examples/1_1_hello_world.go
// 测试命令: curl http://localhost:8080/ping
//          curl http://localhost:8080/health
// ============================================================================

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// ============================================================================
// 核心概念解析
// ============================================================================
//
// 【什么是 Gin Engine？】
//
// Engine 是 Gin 框架的核心结构体，它：
// 1. 实现了 http.Handler 接口 (ServeHTTP 方法)
// 2. 包含路由树 (基于 Radix Tree，高性能路由匹配)
// 3. 管理中间件链
// 4. 存储配置信息
//
// 【gin.New() vs gin.Default() 源码对比】
//
// gin.New() 源码:
//   func New() *Engine {
//       engine := &Engine{...}
//       engine.RouterGroup.engine = engine
//       engine.pool.New = func() any { return engine.allocateContext() }
//       return engine
//   }
//
// gin.Default() 源码:
//   func Default() *Engine {
//       engine := New()
//       engine.Use(Logger(), Recovery())  // 关键区别：自动添加两个中间件
//       return engine
//   }
//
// 【为什么生产环境推荐 gin.New()？】
//
// 1. Logger 中间件的日志格式是固定的，生产环境通常需要 JSON 格式日志
// 2. 生产环境需要更细粒度的中间件控制
// 3. 可能需要自定义 Recovery 行为（如：上报错误到监控系统）
//
// ============================================================================

func main() {
	// ========================================================================
	// 第一步：设置运行模式
	// ========================================================================
	//
	// Gin 有三种模式:
	// - gin.DebugMode   ("debug")   : 开发模式，输出详细日志
	// - gin.ReleaseMode ("release") : 生产模式，关闭调试信息
	// - gin.TestMode    ("test")    : 测试模式，用于单元测试
	//
	// 【最佳实践】从环境变量读取，而不是硬编码
	//
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = gin.DebugMode // 默认开发模式
	}
	gin.SetMode(mode)

	// ========================================================================
	// 第二步：创建 Engine 实例
	// ========================================================================
	//
	// 【对比演示】
	//
	// 方式一：gin.Default() - 适合快速开发
	// r := gin.Default()
	//
	// 方式二：gin.New() - 推荐生产环境使用
	r := gin.New()

	// ========================================================================
	// 第三步：手动添加中间件
	// ========================================================================
	//
	// 使用 gin.New() 后需要手动添加中间件
	// 这样可以：
	// 1. 控制中间件顺序
	// 2. 使用自定义的 Logger/Recovery
	// 3. 按需添加其他中间件
	//
	r.Use(gin.Logger())   // 请求日志中间件
	r.Use(gin.Recovery()) // panic 恢复中间件

	// ========================================================================
	// 第四步：注册路由
	// ========================================================================
	//
	// 【gin.H 是什么？】
	//
	// type H map[string]any
	//
	// 它只是 map[string]interface{} 的类型别名，用于方便构造 JSON 响应
	// any 是 Go 1.18 引入的，等价于 interface{}
	//

	// 健康检查接口 - K8s/Docker 部署必备
	// 用于：
	// 1. 负载均衡器健康检查
	// 2. K8s 的 livenessProbe / readinessProbe
	// 3. 监控系统存活检测
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// 基本的 ping-pong 接口
	r.GET("/ping", func(c *gin.Context) {
		// c.JSON() 会自动设置 Content-Type: application/json
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// ========================================================================
	// 第五步：配置并启动服务器
	// ========================================================================
	//
	// 【r.Run() vs http.Server 对比】
	//
	// r.Run(":8080") 内部实现:
	//   func (engine *Engine) Run(addr ...string) error {
	//       address := resolveAddress(addr)
	//       return http.ListenAndServe(address, engine)
	//   }
	//
	// 问题：
	// 1. 无法设置读写超时（可能被慢客户端攻击）
	// 2. 无法优雅关闭（强制断开连接）
	//
	// 【推荐方式】使用 http.Server
	//

	// 从环境变量读取端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 创建自定义 Server，配置超时参数
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r, // Engine 实现了 http.Handler 接口

		// 【重要】超时配置，防止慢客户端攻击
		ReadTimeout:  10 * time.Second, // 读取请求的超时时间
		WriteTimeout: 10 * time.Second, // 写入响应的超时时间
		IdleTimeout:  60 * time.Second, // Keep-Alive 连接的空闲超时
	}

	// ========================================================================
	// 第六步：优雅启动与关闭
	// ========================================================================
	//
	// 【为什么需要优雅关闭？】
	//
	// 直接关闭服务器会导致：
	// 1. 正在处理的请求被强制中断
	// 2. 数据库事务可能未提交
	// 3. 用户体验差
	//
	// 优雅关闭会：
	// 1. 停止接受新请求
	// 2. 等待现有请求处理完成
	// 3. 设置超时时间，超时后强制关闭
	//

	// 异步启动服务器
	go func() {
		log.Printf("[INFO] Server starting on http://localhost:%s", port)
		log.Printf("[INFO] Health check: http://localhost:%s/health", port)
		log.Printf("[INFO] Ping test: http://localhost:%s/ping", port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[FATAL] Server failed to start: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	// SIGINT: Ctrl+C
	// SIGTERM: kill 命令 (K8s 默认发送)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[INFO] Shutting down server...")

	// 创建超时上下文，最多等待 5 秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[FATAL] Server forced to shutdown: %v", err)
	}

	log.Println("[INFO] Server exited gracefully")
}

// ============================================================================
// 常见错误示例 (反面教材)
// ============================================================================
//
// 【错误 1】忽略 Run() 的返回值
//
//   r.Run(":8080")  // 错误！端口被占用时程序会静默失败
//
// 【正确做法】
//
//   if err := r.Run(":8080"); err != nil {
//       log.Fatal(err)
//   }
//
// ============================================================================
//
// 【错误 2】硬编码配置
//
//   r.Run(":8080")  // 错误！部署时需要改代码
//
// 【正确做法】
//
//   port := os.Getenv("PORT")
//   if port == "" {
//       port = "8080"
//   }
//
// ============================================================================
//
// 【错误 3】生产环境使用 Debug 模式
//
//   // 忘记设置模式，默认是 debug
//
// 【正确做法】
//
//   gin.SetMode(gin.ReleaseMode)  // 或从环境变量读取
//
// ============================================================================

// ============================================================================
// 练习题
// ============================================================================
//
// 1. 添加一个 /version 接口，返回应用版本号（从环境变量 APP_VERSION 读取）
// 2. 添加一个 /ready 接口，用于 K8s readinessProbe
// 3. 尝试把 ReadTimeout 设置为 1 秒，然后用 curl 发送一个大请求，观察行为
//
// ============================================================================
