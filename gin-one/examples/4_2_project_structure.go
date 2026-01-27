// ============================================================================
// 4.2 项目目录结构设计
// ============================================================================
// 本文件是教学文档，不需要运行
// ============================================================================

package main

// ============================================================================
// 项目结构概述
// ============================================================================
//
// 【为什么需要分层？】
//
// 1. 职责分离：每层只做一件事
// 2. 可测试性：每层可以独立测试
// 3. 可维护性：修改一层不影响其他层
// 4. 可扩展性：新增功能只需添加文件
//
// 【常见分层模式】
//
// 1. MVC (Model-View-Controller)
// 2. 三层架构 (Controller-Service-Repository)
// 3. Clean Architecture (洋葱架构)
// 4. DDD (Domain-Driven Design)
//
// 对于中小型项目，推荐使用三层架构
//
// ============================================================================

// ============================================================================
// 推荐项目结构 (中小型项目)
// ============================================================================
//
// gin-project/
// ├── cmd/                      # 应用入口
// │   └── server/
// │       └── main.go           # 主程序入口
// │
// ├── internal/                 # 私有代码（不对外暴露）
// │   ├── config/               # 配置管理
// │   │   └── config.go
// │   │
// │   ├── handler/              # HTTP Handler (Controller 层)
// │   │   ├── user.go
// │   │   ├── post.go
// │   │   └── response.go       # 统一响应封装
// │   │
// │   ├── service/              # 业务逻辑层
// │   │   ├── user.go
// │   │   └── post.go
// │   │
// │   ├── repository/           # 数据访问层 (DAO)
// │   │   ├── user.go
// │   │   └── post.go
// │   │
// │   ├── model/                # 数据模型
// │   │   ├── user.go
// │   │   └── post.go
// │   │
// │   ├── middleware/           # 中间件
// │   │   ├── auth.go
// │   │   ├── cors.go
// │   │   └── logger.go
// │   │
// │   └── router/               # 路由注册
// │       └── router.go
// │
// ├── pkg/                      # 公共代码（可被外部引用）
// │   ├── database/             # 数据库连接
// │   │   └── mysql.go
// │   ├── logger/               # 日志
// │   │   └── zap.go
// │   └── utils/                # 工具函数
// │       └── hash.go
// │
// ├── api/                      # API 定义
// │   └── v1/
// │       └── user.go           # 请求/响应结构体
// │
// ├── configs/                  # 配置文件
// │   ├── config.yaml
// │   ├── config.dev.yaml
// │   └── config.prod.yaml
// │
// ├── scripts/                  # 脚本
// │   └── migrate.sh
// │
// ├── docs/                     # 文档
// │   └── swagger/
// │
// ├── test/                     # 测试
// │   └── integration/
// │
// ├── .gitignore
// ├── Dockerfile
// ├── docker-compose.yaml
// ├── Makefile
// ├── go.mod
// └── go.sum
//
// ============================================================================

// ============================================================================
// 各层职责详解
// ============================================================================
//
// 【Handler 层 (Controller)】
//
// 职责：
// 1. 接收 HTTP 请求
// 2. 参数绑定和校验
// 3. 调用 Service 层
// 4. 返回 HTTP 响应
//
// 不应该做：
// - 业务逻辑
// - 直接操作数据库
//
// ============================================================================
//
// 【Service 层】
//
// 职责：
// 1. 业务逻辑处理
// 2. 调用 Repository 层
// 3. 事务管理
// 4. 数据组装
//
// 不应该做：
// - HTTP 相关操作
// - SQL 语句
//
// ============================================================================
//
// 【Repository 层 (DAO)】
//
// 职责：
// 1. 数据库 CRUD 操作
// 2. SQL 语句封装
// 3. 数据映射
//
// 不应该做：
// - 业务逻辑
// - HTTP 相关操作
//
// ============================================================================

// ============================================================================
// 代码示例：分层架构
// ============================================================================

// --------------- api/v1/user.go ---------------

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// --------------- internal/model/user.go ---------------

import "gorm.io/gorm"

// User 用户模型
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null"`
	Email    string `gorm:"uniqueIndex"`
	Password string `gorm:"not null"`
}

// --------------- internal/repository/user.go ---------------

// UserRepository 用户数据访问层接口
type UserRepository interface {
	Create(user *User) error
	FindByID(id uint) (*User, error)
	FindByUsername(username string) (*User, error)
	Update(user *User) error
	Delete(id uint) error
}

// userRepository 接口实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建 UserRepository 实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*User, error) {
	var user User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*User, error) {
	var user User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&User{}, id).Error
}

// --------------- internal/service/user.go ---------------

import "errors"

// UserService 用户服务层接口
type UserService interface {
	CreateUser(req *CreateUserRequest) (*UserResponse, error)
	GetUser(id uint) (*UserResponse, error)
}

// userService 服务实现
type userService struct {
	userRepo UserRepository
}

// NewUserService 创建 UserService 实例
func NewUserService(repo UserRepository) UserService {
	return &userService{userRepo: repo}
}

func (s *userService) CreateUser(req *CreateUserRequest) (*UserResponse, error) {
	// 检查用户名是否已存在
	existing, _ := s.userRepo.FindByUsername(req.Username)
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	// 创建用户
	user := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashPassword(req.Password), // 密码加密
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 返回响应
	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *userService) GetUser(id uint) (*UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func hashPassword(password string) string {
	// 实际应该用 bcrypt
	return password
}

// --------------- internal/handler/user.go ---------------

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService UserService
}

// NewUserHandler 创建 UserHandler 实例
func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{userService: service}
}

// Create 创建用户
func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Get 获取用户
func (h *UserHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user, err := h.userService.GetUser(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// --------------- internal/router/router.go ---------------

// SetupRouter 设置路由
func SetupRouter(
	userHandler *UserHandler,
	// 其他 handler...
) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("", userHandler.Create)
			users.GET("/:id", userHandler.Get)
		}
	}

	return r
}

// --------------- cmd/server/main.go ---------------

// func main() {
// 	// 1. 加载配置
// 	cfg := config.Load()
//
// 	// 2. 初始化数据库
// 	db := database.NewMySQL(cfg.Database)
//
// 	// 3. 依赖注入
// 	userRepo := repository.NewUserRepository(db)
// 	userService := service.NewUserService(userRepo)
// 	userHandler := handler.NewUserHandler(userService)
//
// 	// 4. 设置路由
// 	r := router.SetupRouter(userHandler)
//
// 	// 5. 启动服务器
// 	r.Run(cfg.Server.Addr)
// }

// ============================================================================
// 依赖注入
// ============================================================================
//
// 【什么是依赖注入？】
//
// 将依赖通过参数传入，而不是在内部创建
//
// 【优点】
// 1. 解耦：组件之间通过接口交互
// 2. 可测试：可以传入 mock 对象
// 3. 可配置：可以替换不同实现
//
// 【实现方式】
// 1. 构造函数注入（推荐）
// 2. 使用依赖注入框架（wire、dig）
//
// ============================================================================

// ============================================================================
// 使用 wire 依赖注入（可选）
// ============================================================================
//
// 安装: go install github.com/google/wire/cmd/wire@latest
//
// --------------- wire.go ---------------
//
// //go:build wireinject
//
// package main
//
// import (
// 	"github.com/google/wire"
// )
//
// func InitializeApp() (*gin.Engine, error) {
// 	wire.Build(
// 		database.NewMySQL,
// 		repository.NewUserRepository,
// 		service.NewUserService,
// 		handler.NewUserHandler,
// 		router.SetupRouter,
// 	)
// 	return nil, nil
// }
//
// ============================================================================

// ============================================================================
// 易错点总结
// ============================================================================
//
// 1. 【循环依赖】
//    Service A 依赖 Service B，B 又依赖 A
//    解决：提取公共接口或合并服务
//
// 2. 【层级混乱】
//    Handler 直接调用 Repository
//    Handler 中写业务逻辑
//    务必保持: Handler → Service → Repository
//
// 3. 【全局变量滥用】
//    用全局 DB 变量而不是依赖注入
//    导致难以测试和替换
//
// 4. 【接口过度设计】
//    小项目不需要每个组件都定义接口
//    只在需要 mock 测试或多实现时才定义
//
// 5. 【目录结构过深】
//    不要为了分层而分层
//    小项目可以简化结构
//
// ============================================================================

// ============================================================================
// 简化版项目结构（小型项目）
// ============================================================================
//
// gin-simple/
// ├── main.go              # 入口 + 路由
// ├── config.go            # 配置
// ├── model.go             # 模型
// ├── handler.go           # 处理器
// ├── service.go           # 服务（可选）
// ├── middleware.go        # 中间件
// ├── config.yaml          # 配置文件
// └── go.mod
//
// ============================================================================

// ============================================================================
// 练习题
// ============================================================================
//
// 1. 创建一个新项目，按照推荐结构组织代码
//    - 实现用户 CRUD
//    - 使用依赖注入
//
// 2. 为项目添加单元测试
//    - Mock Repository 测试 Service
//    - Mock Service 测试 Handler
//
// 3. 使用 wire 实现自动依赖注入
//
// ============================================================================

func main() {
	// 这是文档文件，不需要运行
	println("This is a documentation file. See the code comments for project structure guide.")
}
