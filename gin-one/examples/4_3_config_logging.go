// ============================================================================
// 4.3 配置管理与日志系统
// ============================================================================
// 运行方式: go run examples/4_3_config_logging.go
// 需要先安装:
//   go get -u github.com/spf13/viper
//   go get -u go.uber.org/zap
//   go get -u gopkg.in/natefinch/lumberjack.v2
// ============================================================================

package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ============================================================================
// 一、Viper 配置管理
// ============================================================================
//
// 【Viper 特点】
//
// 1. 支持多种配置格式 (JSON, YAML, TOML, ENV)
// 2. 支持环境变量
// 3. 支持配置热加载
// 4. 支持远程配置 (etcd, consul)
// 5. 支持配置默认值
//
// 【配置优先级】(从高到低)
//
// 1. 显式调用 Set 设置
// 2. 命令行标志
// 3. 环境变量
// 4. 配置文件
// 5. 默认值
//
// ============================================================================

// Config 配置结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Log      LogConfig      `mapstructure:"log"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
	Mode         string        `mapstructure:"mode"`          // gin.ReleaseMode / gin.DebugMode
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Driver       string `mapstructure:"driver"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`       // debug, info, warn, error
	Format     string `mapstructure:"format"`      // json, console
	Filename   string `mapstructure:"filename"`    // 日志文件路径
	MaxSize    int    `mapstructure:"max_size"`    // 单个文件最大 MB
	MaxBackups int    `mapstructure:"max_backups"` // 保留旧文件数量
	MaxAge     int    `mapstructure:"max_age"`     // 保留天数
	Compress   bool   `mapstructure:"compress"`    // 是否压缩
}

type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	ExpireTime time.Duration `mapstructure:"expire_time"`
}

// 全局配置
var AppConfig Config

// 全局日志
var Logger *zap.Logger

// InitConfig 初始化配置
func InitConfig() error {
	// 设置配置文件名和路径
	viper.SetConfigName("config")        // 配置文件名（不带扩展名）
	viper.SetConfigType("yaml")          // 配置文件类型
	viper.AddConfigPath(".")             // 当前目录
	viper.AddConfigPath("./configs")     // configs 目录
	viper.AddConfigPath("$HOME/.myapp")  // home 目录

	// 设置默认值
	setDefaults()

	// 支持环境变量
	viper.AutomaticEnv()
	// 环境变量前缀，如 APP_SERVER_PORT
	viper.SetEnvPrefix("APP")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		// 配置文件不存在时使用默认值
		log.Printf("Config file not found, using defaults: %v", err)
	}

	// 解析到结构体
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return err
	}

	// 监听配置变化（热加载）
	// viper.WatchConfig()
	// viper.OnConfigChange(func(e fsnotify.Event) {
	// 	log.Println("Config file changed:", e.Name)
	// 	viper.Unmarshal(&AppConfig)
	// })

	return nil
}

// setDefaults 设置默认值
func setDefaults() {
	// Server
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "10s")
	viper.SetDefault("server.write_timeout", "10s")

	// Database
	viper.SetDefault("database.driver", "mysql")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)

	// Redis
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)

	// Log
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.filename", "logs/app.log")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_backups", 3)
	viper.SetDefault("log.max_age", 7)
	viper.SetDefault("log.compress", true)

	// JWT
	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("jwt.expire_time", "24h")
}

// ============================================================================
// 二、Zap 日志系统
// ============================================================================
//
// 【为什么用 Zap？】
//
// 1. 高性能：比标准库 log 快 10 倍以上
// 2. 结构化：原生支持 JSON 格式
// 3. 分级：支持 Debug/Info/Warn/Error 等级别
// 4. 可扩展：支持自定义输出、格式
//
// 【两种 Logger】
//
// 1. zap.Logger: 类型安全，性能最高
// 2. zap.SugaredLogger: API 友好，性能稍低
//
// ============================================================================

// InitLogger 初始化日志
func InitLogger() error {
	// 日志级别
	level := getLogLevel(AppConfig.Log.Level)

	// 日志编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,    // 小写级别
		EncodeTime:     zapcore.ISO8601TimeEncoder,       // ISO8601 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,       // 短调用者
	}

	// 选择编码器
	var encoder zapcore.Encoder
	if AppConfig.Log.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 日志轮转 (使用 lumberjack)
	lumberJackLogger := &lumberjack.Logger{
		Filename:   AppConfig.Log.Filename,
		MaxSize:    AppConfig.Log.MaxSize,    // MB
		MaxBackups: AppConfig.Log.MaxBackups, // 保留文件数
		MaxAge:     AppConfig.Log.MaxAge,     // 保留天数
		Compress:   AppConfig.Log.Compress,   // 压缩
	}

	// 创建 Core
	// 同时输出到文件和控制台
	core := zapcore.NewTee(
		// 文件输出
		zapcore.NewCore(encoder, zapcore.AddSync(lumberJackLogger), level),
		// 控制台输出
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
	)

	// 创建 Logger
	Logger = zap.New(core,
		zap.AddCaller(),                   // 添加调用者信息
		zap.AddCallerSkip(1),              // 跳过一层调用
		zap.AddStacktrace(zap.ErrorLevel), // Error 级别添加堆栈
	)

	return nil
}

// getLogLevel 获取日志级别
func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// ============================================================================
// 三、Gin 集成 Zap
// ============================================================================

// GinLogger 返回 Zap 日志中间件
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		cost := time.Since(start)

		Logger.Info("HTTP Request",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery 返回 Zap Recovery 中间件
func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				Logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.Stack("stack"),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "Internal Server Error",
				})
			}
		}()
		c.Next()
	}
}

// ============================================================================
// 主程序
// ============================================================================

func main() {
	// 1. 初始化配置
	if err := InitConfig(); err != nil {
		log.Fatalf("Failed to init config: %v", err)
	}

	// 2. 初始化日志
	if err := InitLogger(); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	defer Logger.Sync()

	// 3. 设置 Gin 模式
	gin.SetMode(AppConfig.Server.Mode)

	// 4. 创建 Gin Engine
	r := gin.New()

	// 5. 使用自定义中间件
	r.Use(GinLogger())
	r.Use(GinRecovery())

	// 6. 路由
	r.GET("/ping", func(c *gin.Context) {
		Logger.Info("Ping handler called",
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"config": gin.H{
				"server_mode": AppConfig.Server.Mode,
				"server_port": AppConfig.Server.Port,
				"log_level":   AppConfig.Log.Level,
			},
		})
	})

	r.GET("/config", func(c *gin.Context) {
		// 不要在生产环境暴露完整配置！
		c.JSON(http.StatusOK, gin.H{
			"server": gin.H{
				"mode": AppConfig.Server.Mode,
				"port": AppConfig.Server.Port,
			},
			"log": gin.H{
				"level":  AppConfig.Log.Level,
				"format": AppConfig.Log.Format,
			},
		})
	})

	r.GET("/panic", func(c *gin.Context) {
		panic("test panic!")
	})

	// 使用不同日志级别
	r.GET("/log-levels", func(c *gin.Context) {
		Logger.Debug("This is debug log")
		Logger.Info("This is info log",
			zap.String("key", "value"),
			zap.Int("count", 42),
		)
		Logger.Warn("This is warn log")
		Logger.Error("This is error log")

		c.JSON(http.StatusOK, gin.H{"message": "check your logs"})
	})

	// 7. 启动服务器
	Logger.Info("Server starting",
		zap.Int("port", AppConfig.Server.Port),
		zap.String("mode", AppConfig.Server.Mode),
	)

	addr := fmt.Sprintf(":%d", AppConfig.Server.Port)
	if err := r.Run(addr); err != nil {
		Logger.Fatal("Failed to start server", zap.Error(err))
	}
}

import "fmt"

// ============================================================================
// 配置文件示例: config.yaml
// ============================================================================
//
// server:
//   mode: debug
//   port: 8080
//   read_timeout: 10s
//   write_timeout: 10s
//
// database:
//   driver: mysql
//   host: localhost
//   port: 3306
//   username: root
//   password: password
//   database: myapp
//   max_idle_conns: 10
//   max_open_conns: 100
//
// redis:
//   host: localhost
//   port: 6379
//   password: ""
//   db: 0
//
// log:
//   level: debug
//   format: json
//   filename: logs/app.log
//   max_size: 100
//   max_backups: 3
//   max_age: 7
//   compress: true
//
// jwt:
//   secret: your-super-secret-key
//   expire_time: 24h
//
// ============================================================================

// ============================================================================
// 测试命令
// ============================================================================
//
// # 创建配置文件
// mkdir -p logs
// cat > config.yaml << 'EOF'
// server:
//   mode: debug
//   port: 8080
//
// log:
//   level: debug
//   format: console
//   filename: logs/app.log
// EOF
//
// # 运行
// go run examples/4_3_config_logging.go
//
// # 测试
// curl http://localhost:8080/ping
// curl http://localhost:8080/config
// curl http://localhost:8080/log-levels
// curl http://localhost:8080/panic
//
// # 使用环境变量覆盖配置
// APP_SERVER_PORT=9090 go run examples/4_3_config_logging.go
//
// ============================================================================

// ============================================================================
// 易错点总结
// ============================================================================
//
// 1. 【Viper 配置名大小写】
//    YAML 用小写下划线: max_idle_conns
//    环境变量用大写: APP_DATABASE_MAX_IDLE_CONNS
//    结构体 tag: mapstructure:"max_idle_conns"
//
// 2. 【时间类型配置】
//    YAML 中写 "10s", "1h" 等字符串
//    Viper 会自动解析为 time.Duration
//
// 3. 【日志文件目录】
//    lumberjack 不会自动创建目录
//    需要先 mkdir -p logs
//
// 4. 【Logger.Sync()】
//    程序退出前必须调用 Logger.Sync()
//    否则缓冲区日志可能丢失
//
// 5. 【生产环境配置】
//    不要暴露完整配置到 API
//    敏感信息用环境变量传递
//    配置文件不要提交到 git
//
// 6. 【配置热加载的坑】
//    WatchConfig 在某些平台上不稳定
//    修改配置后需要重新 Unmarshal
//
// ============================================================================

// ============================================================================
// 练习题
// ============================================================================
//
// 1. 实现多环境配置:
//    - config.dev.yaml
//    - config.prod.yaml
//    - 通过环境变量 APP_ENV 选择
//
// 2. 实现配置验证:
//    - 必填项检查
//    - 数值范围检查
//    - 格式验证
//
// 3. 实现日志采样:
//    - 只记录 10% 的 Debug 日志
//    - Error 日志全部记录
//
// ============================================================================
