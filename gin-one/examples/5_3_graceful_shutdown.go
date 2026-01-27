// ============================================================================
// 5.3 优雅关闭与 Docker 部署
// ============================================================================
// 运行方式: go run examples/5_3_graceful_shutdown.go
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
// 优雅关闭核心概念
// ============================================================================
//
// 【为什么需要优雅关闭？】
//
// 直接关闭服务器 (kill -9) 会导致：
// 1. 正在处理的请求被强制中断
// 2. 数据库连接可能未正确关闭
// 3. 文件句柄泄露
// 4. 事务未提交
// 5. 消息队列消息丢失
//
// 【优雅关闭的流程】
//
// 1. 收到关闭信号 (SIGTERM/SIGINT)
// 2. 停止接受新的请求
// 3. 等待正在处理的请求完成
// 4. 关闭数据库连接、消息队列等资源
// 5. 超时后强制退出
//
// ┌─────────────────────────────────────────────────────────┐
// │                     关闭信号                             │
// │                        │                                 │
// │                        ▼                                 │
// │              ┌─────────────────┐                        │
// │              │ 停止接受新请求  │                        │
// │              └────────┬────────┘                        │
// │                       │                                 │
// │                       ▼                                 │
// │              ┌─────────────────┐                        │
// │              │ 等待请求完成    │◄──── 超时退出          │
// │              └────────┬────────┘                        │
// │                       │                                 │
// │                       ▼                                 │
// │              ┌─────────────────┐                        │
// │              │ 关闭资源连接    │                        │
// │              └────────┬────────┘                        │
// │                       │                                 │
// │                       ▼                                 │
// │                    退出                                  │
// └─────────────────────────────────────────────────────────┘
//
// ============================================================================

// ============================================================================
// Kubernetes 中的优雅关闭
// ============================================================================
//
// 【K8s Pod 关闭流程】
//
// 1. K8s 发送 SIGTERM 信号
// 2. 等待 terminationGracePeriodSeconds（默认 30 秒）
// 3. 如果还没退出，发送 SIGKILL 强制杀死
//
// 【最佳实践】
//
// 1. 应用启动时注册信号处理
// 2. 收到 SIGTERM 后开始优雅关闭
// 3. 关闭时间 < terminationGracePeriodSeconds
// 4. 实现健康检查接口
//
// ============================================================================

func main() {
	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// ========================================================================
	// 健康检查接口
	// ========================================================================

	// Liveness Probe - 存活检查
	// K8s 用来判断是否需要重启容器
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "alive",
		})
	})

	// Readiness Probe - 就绪检查
	// K8s 用来判断是否可以接收流量
	var isReady = true // 实际应该检查数据库、依赖服务等
	r.GET("/ready", func(c *gin.Context) {
		if isReady {
			c.JSON(http.StatusOK, gin.H{
				"status": "ready",
			})
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
			})
		}
	})

	// ========================================================================
	// 业务接口
	// ========================================================================

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// 模拟慢接口（用于测试优雅关闭）
	r.GET("/slow", func(c *gin.Context) {
		log.Println("Slow request started...")
		time.Sleep(10 * time.Second) // 模拟耗时操作
		log.Println("Slow request completed.")
		c.JSON(http.StatusOK, gin.H{
			"message": "slow response completed",
		})
	})

	// ========================================================================
	// 优雅关闭实现
	// ========================================================================

	// 配置服务器
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
		// 超时配置
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// 在 goroutine 中启动服务器
	go func() {
		log.Printf("Server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	// SIGINT: Ctrl+C
	// SIGTERM: kill 命令（K8s 默认发送）
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	log.Printf("Received signal: %v", sig)
	log.Println("Shutting down server...")

	// 标记为未就绪（K8s 会停止发送新流量）
	isReady = false

	// 创建超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// 关闭其他资源（数据库、Redis、消息队列等）
	log.Println("Closing database connections...")
	// db.Close()

	log.Println("Closing Redis connections...")
	// redis.Close()

	log.Println("Server exited gracefully")
}

// ============================================================================
// Dockerfile
// ============================================================================
//
// # 多阶段构建
//
// # 阶段 1: 构建
// FROM golang:1.21-alpine AS builder
//
// # 设置工作目录
// WORKDIR /app
//
// # 安装依赖（利用缓存）
// COPY go.mod go.sum ./
// RUN go mod download
//
// # 复制源码
// COPY . .
//
// # 构建
// # CGO_ENABLED=0: 禁用 CGO，生成静态链接的二进制文件
// # -ldflags="-s -w": 去除调试信息，减小体积
// RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server
//
// # 阶段 2: 运行
// FROM alpine:3.18
//
// # 安装时区数据和证书
// RUN apk --no-cache add ca-certificates tzdata
//
// # 设置时区
// ENV TZ=Asia/Shanghai
//
// # 创建非 root 用户
// RUN adduser -D -g '' appuser
//
// WORKDIR /app
//
// # 从构建阶段复制二进制文件
// COPY --from=builder /app/server .
//
// # 复制配置文件
// COPY configs/ ./configs/
//
// # 切换到非 root 用户
// USER appuser
//
// # 暴露端口
// EXPOSE 8080
//
// # 健康检查
// HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
//     CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1
//
// # 启动命令
// CMD ["./server"]
//
// ============================================================================

// ============================================================================
// docker-compose.yaml
// ============================================================================
//
// version: '3.8'
//
// services:
//   app:
//     build: .
//     ports:
//       - "8080:8080"
//     environment:
//       - GIN_MODE=release
//       - DB_HOST=mysql
//       - REDIS_HOST=redis
//     depends_on:
//       mysql:
//         condition: service_healthy
//       redis:
//         condition: service_started
//     restart: unless-stopped
//     healthcheck:
//       test: ["CMD", "wget", "--spider", "http://localhost:8080/healthz"]
//       interval: 30s
//       timeout: 3s
//       retries: 3
//
//   mysql:
//     image: mysql:8.0
//     environment:
//       MYSQL_ROOT_PASSWORD: rootpassword
//       MYSQL_DATABASE: myapp
//     volumes:
//       - mysql_data:/var/lib/mysql
//     healthcheck:
//       test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
//       interval: 10s
//       timeout: 5s
//       retries: 5
//
//   redis:
//     image: redis:7-alpine
//     volumes:
//       - redis_data:/data
//
// volumes:
//   mysql_data:
//   redis_data:
//
// ============================================================================

// ============================================================================
// Kubernetes Deployment (deployment.yaml)
// ============================================================================
//
// apiVersion: apps/v1
// kind: Deployment
// metadata:
//   name: gin-app
//   labels:
//     app: gin-app
// spec:
//   replicas: 3
//   selector:
//     matchLabels:
//       app: gin-app
//   template:
//     metadata:
//       labels:
//         app: gin-app
//     spec:
//       # 优雅关闭等待时间
//       terminationGracePeriodSeconds: 30
//       containers:
//       - name: gin-app
//         image: your-registry/gin-app:latest
//         ports:
//         - containerPort: 8080
//         # 资源限制
//         resources:
//           requests:
//             memory: "64Mi"
//             cpu: "100m"
//           limits:
//             memory: "256Mi"
//             cpu: "500m"
//         # 存活检查
//         livenessProbe:
//           httpGet:
//             path: /healthz
//             port: 8080
//           initialDelaySeconds: 5
//           periodSeconds: 10
//           timeoutSeconds: 3
//           failureThreshold: 3
//         # 就绪检查
//         readinessProbe:
//           httpGet:
//             path: /ready
//             port: 8080
//           initialDelaySeconds: 5
//           periodSeconds: 5
//           timeoutSeconds: 3
//           failureThreshold: 3
//         # 环境变量
//         env:
//         - name: GIN_MODE
//           value: "release"
//         - name: DB_HOST
//           valueFrom:
//             configMapKeyRef:
//               name: gin-app-config
//               key: db_host
//
// ---
// apiVersion: v1
// kind: Service
// metadata:
//   name: gin-app
// spec:
//   selector:
//     app: gin-app
//   ports:
//   - port: 80
//     targetPort: 8080
//   type: ClusterIP
//
// ============================================================================

// ============================================================================
// Makefile
// ============================================================================
//
// .PHONY: build run test clean docker-build docker-run
//
// APP_NAME=gin-app
// VERSION=1.0.0
//
// # 构建
// build:
// 	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/$(APP_NAME) ./cmd/server
//
// # 运行
// run:
// 	go run ./cmd/server
//
// # 测试
// test:
// 	go test -v ./...
//
// # 清理
// clean:
// 	rm -rf bin/
//
// # 构建 Docker 镜像
// docker-build:
// 	docker build -t $(APP_NAME):$(VERSION) .
//
// # 运行 Docker 容器
// docker-run:
// 	docker run -p 8080:8080 $(APP_NAME):$(VERSION)
//
// # 使用 docker-compose 启动
// docker-up:
// 	docker-compose up -d
//
// # 停止 docker-compose
// docker-down:
// 	docker-compose down
//
// ============================================================================

// ============================================================================
// 测试优雅关闭
// ============================================================================
//
// 终端 1: 启动服务器
// $ go run examples/5_3_graceful_shutdown.go
//
// 终端 2: 发送慢请求
// $ curl http://localhost:8080/slow
//
// 终端 1: 在慢请求处理期间按 Ctrl+C
// 观察输出：
// - 收到信号
// - 等待慢请求完成
// - 然后退出
//
// ============================================================================

// ============================================================================
// 易错点总结
// ============================================================================
//
// 1. 【信号处理必须注册】
//    不注册信号处理，Ctrl+C 会直接杀死进程
//    没有优雅关闭的机会
//
// 2. 【Shutdown 超时设置】
//    超时太短：请求被强制中断
//    超时太长：部署更新慢
//    建议：比最长请求稍长即可
//
// 3. 【K8s terminationGracePeriodSeconds】
//    必须大于应用的 Shutdown 超时
//    否则 K8s 会发送 SIGKILL
//
// 4. 【健康检查区分】
//    liveness: 判断是否存活（失败则重启）
//    readiness: 判断是否就绪（失败则不接收流量）
//
// 5. 【资源关闭顺序】
//    先关闭 HTTP 服务器
//    再关闭数据库等依赖
//    顺序错误可能导致请求失败
//
// 6. 【Docker 信号传递】
//    使用 exec 形式的 CMD，而不是 shell 形式
//    CMD ["./server"] ✓
//    CMD ./server     ✗ (信号发给 shell，不是应用)
//
// ============================================================================

// ============================================================================
// 练习题
// ============================================================================
//
// 1. 实现一个生产级的优雅关闭:
//    - 先标记为未就绪
//    - 等待一段时间（让 LB 摘除流量）
//    - 再开始关闭
//
// 2. 实现资源清理:
//    - 关闭数据库连接池
//    - 关闭 Redis 连接
//    - 发送完消息队列中的消息
//
// 3. 实现优雅关闭的单元测试
//
// ============================================================================
