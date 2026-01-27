// ============================================================================
// 5.2 Swagger 接口文档自动生成
// ============================================================================
// 运行方式:
//   1. 安装 swag: go install github.com/swaggo/swag/cmd/swag@latest
//   2. 生成文档: swag init -g examples/5_2_swagger.go
//   3. 运行服务: go run examples/5_2_swagger.go
//   4. 访问文档: http://localhost:8080/swagger/index.html
//
// 需要先安装:
//   go get -u github.com/swaggo/gin-swagger
//   go get -u github.com/swaggo/files
// ============================================================================

package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	// 导入 swagger 相关包
	// swaggerFiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
	// _ "your-project/docs" // 导入生成的 docs 包
)

// ============================================================================
// Swagger 核心概念
// ============================================================================
//
// 【什么是 Swagger/OpenAPI？】
//
// OpenAPI (原 Swagger) 是一种 API 描述规范
// 可以自动生成:
// 1. API 文档 (在线查看)
// 2. 客户端 SDK
// 3. 服务端桩代码
// 4. API 测试工具
//
// 【swag 工具工作原理】
//
// 1. 在代码中写注释
// 2. swag init 扫描注释
// 3. 生成 docs/swagger.json
// 4. gin-swagger 读取 JSON 展示 UI
//
// ============================================================================
//
// 【常用注释格式】
//
// @Summary      简短描述
// @Description  详细描述
// @Tags         标签（用于分组）
// @Accept       接受的 MIME 类型
// @Produce      返回的 MIME 类型
// @Param        参数名 位置 类型 必填 "描述"
// @Success      状态码 {数据类型} 模型 "描述"
// @Failure      状态码 {数据类型} 模型 "描述"
// @Router       路径 [方法]
// @Security     安全方案名
//
// ============================================================================

// ============================================================================
// API 通用信息（必须在 main 函数前定义）
// ============================================================================

// @title           Gin Learning API
// @version         1.0
// @description     Gin 框架学习项目 API 文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Bearer Token 认证，格式: Bearer {token}

// ============================================================================
// 数据模型定义
// ============================================================================

// User 用户模型
// @Description 用户信息
type User struct {
	ID        uint   `json:"id" example:"1"`                            // 用户ID
	Username  string `json:"username" example:"zhangsan"`               // 用户名
	Email     string `json:"email" example:"zhangsan@example.com"`      // 邮箱
	CreatedAt string `json:"created_at" example:"2024-01-15T10:30:00Z"` // 创建时间
}

// CreateUserRequest 创建用户请求
// @Description 创建用户的请求参数
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50" example:"zhangsan"` // 用户名
	Email    string `json:"email" binding:"required,email" example:"zhangsan@example.com"` // 邮箱
	Password string `json:"password" binding:"required,min=6" example:"password123"`       // 密码
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Username string `json:"username,omitempty" example:"lisi"` // 用户名
	Email    string `json:"email,omitempty" example:"lisi@example.com"` // 邮箱
}

// Response 通用响应
type Response struct {
	Code    int         `json:"code" example:"0"`       // 状态码
	Message string      `json:"message" example:"成功"` // 消息
	Data    interface{} `json:"data,omitempty"`         // 数据
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`    // 错误码
	Message string `json:"message" example:"参数错误"` // 错误信息
	Error   string `json:"error,omitempty" example:"username is required"` // 错误详情
}

// PaginatedResponse 分页响应
type PaginatedResponse struct {
	Code    int         `json:"code" example:"0"`
	Message string      `json:"message" example:"成功"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total" example:"100"`
	Page    int         `json:"page" example:"1"`
	Size    int         `json:"size" example:"10"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`    // 用户名
	Password string `json:"password" binding:"required" example:"admin123"` // 密码
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIs..."`  // Access Token
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIs..."` // Refresh Token
	ExpiresIn    int    `json:"expires_in" example:"7200"`                        // 过期时间(秒)
}

// ============================================================================
// 模拟数据
// ============================================================================

var userList = []User{
	{ID: 1, Username: "admin", Email: "admin@example.com", CreatedAt: "2024-01-01T00:00:00Z"},
	{ID: 2, Username: "zhangsan", Email: "zhangsan@example.com", CreatedAt: "2024-01-15T10:30:00Z"},
}

// ============================================================================
// Handler 函数（带 Swagger 注释）
// ============================================================================

// Login godoc
// @Summary      用户登录
// @Description  使用用户名密码登录，获取 JWT Token
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest  true  "登录参数"
// @Success      200      {object}  Response{data=LoginResponse}  "登录成功"
// @Failure      400      {object}  ErrorResponse  "参数错误"
// @Failure      401      {object}  ErrorResponse  "用户名或密码错误"
// @Router       /auth/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "参数错误",
			Error:   err.Error(),
		})
		return
	}

	// 模拟验证
	if req.Username == "admin" && req.Password == "admin123" {
		c.JSON(http.StatusOK, Response{
			Code:    0,
			Message: "登录成功",
			Data: LoginResponse{
				AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				ExpiresIn:    7200,
			},
		})
		return
	}

	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Code:    401,
		Message: "用户名或密码错误",
	})
}

// GetUsers godoc
// @Summary      获取用户列表
// @Description  获取所有用户的列表，支持分页
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        page   query     int     false  "页码"       default(1)
// @Param        size   query     int     false  "每页数量"   default(10)
// @Param        keyword query   string  false  "搜索关键字"
// @Success      200    {object}  PaginatedResponse{data=[]User}  "成功"
// @Failure      500    {object}  ErrorResponse  "服务器错误"
// @Security     BearerAuth
// @Router       /users [get]
func GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	c.JSON(http.StatusOK, PaginatedResponse{
		Code:    0,
		Message: "成功",
		Data:    userList,
		Total:   int64(len(userList)),
		Page:    page,
		Size:    size,
	})
}

// GetUser godoc
// @Summary      获取用户详情
// @Description  根据用户 ID 获取用户详细信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "用户ID"
// @Success      200  {object}  Response{data=User}  "成功"
// @Failure      400  {object}  ErrorResponse  "参数错误"
// @Failure      404  {object}  ErrorResponse  "用户不存在"
// @Security     BearerAuth
// @Router       /users/{id} [get]
func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "参数错误",
			Error:   "id must be a number",
		})
		return
	}

	for _, user := range userList {
		if user.ID == uint(id) {
			c.JSON(http.StatusOK, Response{
				Code:    0,
				Message: "成功",
				Data:    user,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, ErrorResponse{
		Code:    404,
		Message: "用户不存在",
	})
}

// CreateUser godoc
// @Summary      创建用户
// @Description  创建新用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        request  body      CreateUserRequest  true  "用户信息"
// @Success      201      {object}  Response{data=User}  "创建成功"
// @Failure      400      {object}  ErrorResponse  "参数错误"
// @Security     BearerAuth
// @Router       /users [post]
func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "参数错误",
			Error:   err.Error(),
		})
		return
	}

	newUser := User{
		ID:        uint(len(userList) + 1),
		Username:  req.Username,
		Email:     req.Email,
		CreatedAt: "2024-01-20T12:00:00Z",
	}
	userList = append(userList, newUser)

	c.JSON(http.StatusCreated, Response{
		Code:    0,
		Message: "创建成功",
		Data:    newUser,
	})
}

// UpdateUser godoc
// @Summary      更新用户
// @Description  根据 ID 更新用户信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id       path      int                true  "用户ID"
// @Param        request  body      UpdateUserRequest  true  "更新内容"
// @Success      200      {object}  Response{data=User}  "更新成功"
// @Failure      400      {object}  ErrorResponse  "参数错误"
// @Failure      404      {object}  ErrorResponse  "用户不存在"
// @Security     BearerAuth
// @Router       /users/{id} [put]
func UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "参数错误",
		})
		return
	}

	for i, user := range userList {
		if user.ID == uint(id) {
			if req.Username != "" {
				userList[i].Username = req.Username
			}
			if req.Email != "" {
				userList[i].Email = req.Email
			}
			c.JSON(http.StatusOK, Response{
				Code:    0,
				Message: "更新成功",
				Data:    userList[i],
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, ErrorResponse{
		Code:    404,
		Message: "用户不存在",
	})
}

// DeleteUser godoc
// @Summary      删除用户
// @Description  根据 ID 删除用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "用户ID"
// @Success      200  {object}  Response  "删除成功"
// @Failure      404  {object}  ErrorResponse  "用户不存在"
// @Security     BearerAuth
// @Router       /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	for i, user := range userList {
		if user.ID == uint(id) {
			userList = append(userList[:i], userList[i+1:]...)
			c.JSON(http.StatusOK, Response{
				Code:    0,
				Message: "删除成功",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, ErrorResponse{
		Code:    404,
		Message: "用户不存在",
	})
}

// ============================================================================
// 主程序
// ============================================================================

func main() {
	r := gin.Default()

	// ========================================================================
	// Swagger 路由
	// ========================================================================

	// 取消下面的注释以启用 Swagger UI
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ========================================================================
	// API 路由
	// ========================================================================

	// 认证接口
	r.POST("/api/v1/auth/login", Login)

	// 用户接口
	v1 := r.Group("/api/v1")
	{
		v1.GET("/users", GetUsers)
		v1.GET("/users/:id", GetUser)
		v1.POST("/users", CreateUser)
		v1.PUT("/users/:id", UpdateUser)
		v1.DELETE("/users/:id", DeleteUser)
	}

	// 打印说明
	println("Server starting on :8080")
	println("")
	println("Swagger UI: http://localhost:8080/swagger/index.html")
	println("(需要先运行 swag init 生成文档)")
	println("")
	println("生成文档命令:")
	println("  swag init -g examples/5_2_swagger.go")
	println("")
	println("测试命令:")
	println(`  curl http://localhost:8080/api/v1/users`)
	println(`  curl http://localhost:8080/api/v1/users/1`)

	r.Run(":8080")
}

// ============================================================================
// Swag 命令
// ============================================================================
//
// # 安装 swag
// go install github.com/swaggo/swag/cmd/swag@latest
//
// # 生成文档
// swag init -g examples/5_2_swagger.go -o docs
//
// # 格式化注释
// swag fmt
//
// ============================================================================

// ============================================================================
// 常用注释格式示例
// ============================================================================
//
// 【路径参数】
// @Param   id   path   int   true   "用户ID"
//
// 【查询参数】
// @Param   page   query   int   false   "页码"   default(1)
// @Param   keyword   query   string   false   "搜索关键字"
//
// 【请求体】
// @Param   request   body   CreateUserRequest   true   "请求体"
//
// 【表单参数】
// @Param   file   formData   file   true   "上传文件"
//
// 【请求头】
// @Param   Authorization   header   string   true   "Bearer Token"
//
// 【枚举类型】
// @Param   status   query   string   false   "状态"   Enums(active, inactive)
//
// 【嵌套对象响应】
// @Success   200   {object}   Response{data=User}   "成功"
//
// 【数组响应】
// @Success   200   {object}   Response{data=[]User}   "成功"
//
// ============================================================================

// ============================================================================
// 易错点总结
// ============================================================================
//
// 1. 【注释必须紧贴函数】
//    注释和函数之间不能有空行
//
// 2. 【类型名大小写】
//    Swagger 通过反射读取类型
//    类型必须是导出的（首字母大写）
//
// 3. 【example 标签】
//    结构体字段的 example 标签用于文档展示
//    example:"张三" 不是 example:"张三"
//
// 4. 【嵌套类型】
//    使用 Response{data=User} 表示泛型响应
//
// 5. 【数组类型】
//    使用 []User 或 Response{data=[]User}
//
// 6. 【安全定义】
//    securityDefinitions 在文件顶部定义
//    @Security 在每个需要认证的接口添加
//
// ============================================================================

// ============================================================================
// 练习题
// ============================================================================
//
// 1. 为项目中的所有接口添加 Swagger 注释
//
// 2. 添加文件上传接口的文档
//
// 3. 实现分组标签（Tags）来组织 API
//
// ============================================================================
