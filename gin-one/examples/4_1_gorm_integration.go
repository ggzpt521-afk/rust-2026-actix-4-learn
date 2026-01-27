// ============================================================================
// 4.1 GORM 集成与 CRUD
// ============================================================================
// 运行方式: go run examples/4_1_gorm_integration.go
// 需要先安装: go get -u gorm.io/gorm gorm.io/driver/sqlite
// ============================================================================

package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ============================================================================
// GORM 核心概念
// ============================================================================
//
// 【什么是 ORM？】
//
// Object-Relational Mapping (对象关系映射)
// 将数据库表映射为 Go 结构体，用面向对象的方式操作数据库
//
// 【GORM 特点】
//
// 1. 全功能 ORM（关联、事务、迁移、钩子等）
// 2. 支持多种数据库（MySQL、PostgreSQL、SQLite、SQL Server）
// 3. 自动迁移（Auto Migration）
// 4. 链式调用
//
// ============================================================================

// ============================================================================
// 模型定义
// ============================================================================

// User 用户模型
type User struct {
	// gorm.Model 包含 ID, CreatedAt, UpdatedAt, DeletedAt
	// type Model struct {
	//     ID        uint           `gorm:"primaryKey"`
	//     CreatedAt time.Time
	//     UpdatedAt time.Time
	//     DeletedAt gorm.DeletedAt `gorm:"index"`
	// }
	gorm.Model

	// 用户名：唯一索引，非空
	Username string `gorm:"uniqueIndex;not null;size:50" json:"username"`

	// 邮箱：唯一索引
	Email string `gorm:"uniqueIndex;size:100" json:"email"`

	// 密码：不返回给前端
	Password string `gorm:"not null" json:"-"`

	// 年龄：默认值
	Age int `gorm:"default:0" json:"age"`

	// 状态：枚举
	Status string `gorm:"type:varchar(20);default:'active'" json:"status"`

	// 关联：一个用户有多篇文章
	Posts []Post `gorm:"foreignKey:UserID" json:"posts,omitempty"`
}

// Post 文章模型
type Post struct {
	gorm.Model
	Title   string `gorm:"not null;size:200" json:"title"`
	Content string `gorm:"type:text" json:"content"`
	UserID  uint   `gorm:"index" json:"user_id"`

	// 属于某个用户
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 自定义表名
func (User) TableName() string {
	return "users"
}

func (Post) TableName() string {
	return "posts"
}

// ============================================================================
// 全局数据库连接
// ============================================================================

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB() error {
	var err error

	// SQLite 连接（开发环境）
	// 生产环境换成 MySQL/PostgreSQL
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		// 日志配置
		Logger: logger.Default.LogMode(logger.Info),
		// 禁用默认事务（提升性能）
		// SkipDefaultTransaction: true,
		// 预编译语句缓存
		PrepareStmt: true,
	})
	if err != nil {
		return err
	}

	// 获取底层 *sql.DB 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// 连接池配置
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大存活时间

	// 自动迁移（开发环境使用，生产环境用 migrate 工具）
	err = DB.AutoMigrate(&User{}, &Post{})
	if err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

func main() {
	// 初始化数据库
	if err := InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	r := gin.Default()

	// ========================================================================
	// 用户 CRUD 接口
	// ========================================================================

	users := r.Group("/users")
	{
		users.POST("", CreateUser)      // 创建用户
		users.GET("", ListUsers)        // 用户列表
		users.GET("/:id", GetUser)      // 获取用户
		users.PUT("/:id", UpdateUser)   // 更新用户
		users.DELETE("/:id", DeleteUser) // 删除用户
	}

	// ========================================================================
	// 文章接口（演示关联）
	// ========================================================================

	posts := r.Group("/posts")
	{
		posts.POST("", CreatePost)
		posts.GET("", ListPosts)
		posts.GET("/:id", GetPost)
	}

	// ========================================================================
	// 高级查询演示
	// ========================================================================

	r.GET("/advanced/query", AdvancedQuery)

	// ========================================================================
	// 事务演示
	// ========================================================================

	r.POST("/transaction", TransactionDemo)

	r.Run(":8080")
}

// ============================================================================
// 请求/响应结构体
// ============================================================================

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Age      int    `json:"age" binding:"gte=0,lte=150"`
}

type UpdateUserRequest struct {
	Username *string `json:"username" binding:"omitempty,min=3,max=50"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Age      *int    `json:"age" binding:"omitempty,gte=0,lte=150"`
	Status   *string `json:"status" binding:"omitempty,oneof=active inactive banned"`
}

type ListUsersQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Status   string `form:"status"`
	Keyword  string `form:"keyword"`
}

// ============================================================================
// 用户 CRUD Handler
// ============================================================================

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // 实际应该加密
		Age:      req.Age,
	}

	// Create 创建记录
	result := DB.Create(&user)
	if result.Error != nil {
		// 处理唯一键冲突
		c.JSON(http.StatusConflict, gin.H{"error": "username or email already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created",
		"user":    user,
	})
}

// ListUsers 用户列表
func ListUsers(c *gin.Context) {
	var query ListUsersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var users []User
	var total int64

	// 构建查询
	db := DB.Model(&User{})

	// 条件过滤
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.Keyword != "" {
		db = db.Where("username LIKE ? OR email LIKE ?",
			"%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	// 统计总数
	db.Count(&total)

	// 分页查询
	offset := (query.Page - 1) * query.PageSize
	db.Offset(offset).Limit(query.PageSize).Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"data":  users,
		"total": total,
		"page":  query.Page,
		"size":  query.PageSize,
	})
}

// GetUser 获取单个用户
func GetUser(c *gin.Context) {
	id := c.Param("id")

	var user User
	// First 查询第一条，未找到返回 ErrRecordNotFound
	result := DB.First(&user, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	// 先检查用户是否存在
	var user User
	if err := DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 构建更新 map（只更新非空字段）
	updates := make(map[string]interface{})
	if req.Username != nil {
		updates["username"] = *req.Username
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Age != nil {
		updates["age"] = *req.Age
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	// Updates 更新多个字段
	if err := DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 重新查询返回最新数据
	DB.First(&user, id)

	c.JSON(http.StatusOK, gin.H{
		"message": "user updated",
		"user":    user,
	})
}

// DeleteUser 删除用户（软删除）
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	// Delete 软删除（设置 deleted_at）
	result := DB.Delete(&User{}, id)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

// ============================================================================
// 文章 Handler（关联查询）
// ============================================================================

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
	UserID  uint   `json:"user_id" binding:"required"`
}

// CreatePost 创建文章
func CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户是否存在
	var user User
	if err := DB.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	post := Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  req.UserID,
	}

	DB.Create(&post)

	c.JSON(http.StatusCreated, post)
}

// ListPosts 文章列表（带关联用户）
func ListPosts(c *gin.Context) {
	var posts []Post

	// Preload 预加载关联数据
	DB.Preload("User").Find(&posts)

	c.JSON(http.StatusOK, posts)
}

// GetPost 获取文章详情
func GetPost(c *gin.Context) {
	id := c.Param("id")

	var post Post
	// Preload 预加载用户信息
	result := DB.Preload("User").First(&post, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// ============================================================================
// 高级查询演示
// ============================================================================

// AdvancedQuery 高级查询示例
func AdvancedQuery(c *gin.Context) {
	// 1. Select 指定字段
	var usernames []string
	DB.Model(&User{}).Pluck("username", &usernames)

	// 2. 原生 SQL
	var count int64
	DB.Raw("SELECT COUNT(*) FROM users WHERE status = ?", "active").Scan(&count)

	// 3. 聚合查询
	type Result struct {
		Status string
		Count  int64
	}
	var results []Result
	DB.Model(&User{}).Select("status, count(*) as count").Group("status").Scan(&results)

	// 4. 子查询
	subQuery := DB.Model(&Post{}).Select("user_id")
	var usersWithPosts []User
	DB.Where("id IN (?)", subQuery).Find(&usersWithPosts)

	// 5. Join 查询
	type UserPostCount struct {
		Username  string
		PostCount int64
	}
	var userPostCounts []UserPostCount
	DB.Model(&User{}).
		Select("users.username, count(posts.id) as post_count").
		Joins("LEFT JOIN posts ON posts.user_id = users.id").
		Group("users.id").
		Scan(&userPostCounts)

	c.JSON(http.StatusOK, gin.H{
		"usernames":        usernames,
		"active_count":     count,
		"status_stats":     results,
		"users_with_posts": usersWithPosts,
		"user_post_counts": userPostCounts,
	})
}

// ============================================================================
// 事务演示
// ============================================================================

// TransactionDemo 事务示例
func TransactionDemo(c *gin.Context) {
	// 方式一：自动事务
	err := DB.Transaction(func(tx *gorm.DB) error {
		// 创建用户
		user := User{Username: "tx_user", Email: "tx@example.com", Password: "123456"}
		if err := tx.Create(&user).Error; err != nil {
			return err // 返回错误会自动回滚
		}

		// 创建文章
		post := Post{Title: "Transaction Post", UserID: user.ID}
		if err := tx.Create(&post).Error; err != nil {
			return err
		}

		// 返回 nil 自动提交
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 方式二：手动事务
	// tx := DB.Begin()
	// defer func() {
	//     if r := recover(); r != nil {
	//         tx.Rollback()
	//     }
	// }()
	// if err := tx.Create(&user).Error; err != nil {
	//     tx.Rollback()
	//     return
	// }
	// tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "transaction success"})
}

// ============================================================================
// 测试命令
// ============================================================================
//
// # 创建用户
// curl -X POST http://localhost:8080/users \
//   -H "Content-Type: application/json" \
//   -d '{"username":"zhangsan","email":"zhangsan@example.com","password":"123456","age":25}'
//
// # 用户列表
// curl "http://localhost:8080/users?page=1&page_size=10&keyword=zhang"
//
// # 获取用户
// curl http://localhost:8080/users/1
//
// # 更新用户
// curl -X PUT http://localhost:8080/users/1 \
//   -H "Content-Type: application/json" \
//   -d '{"age":26,"status":"active"}'
//
// # 删除用户
// curl -X DELETE http://localhost:8080/users/1
//
// # 创建文章
// curl -X POST http://localhost:8080/posts \
//   -H "Content-Type: application/json" \
//   -d '{"title":"Hello GORM","content":"GORM is great!","user_id":1}'
//
// # 文章列表（带用户信息）
// curl http://localhost:8080/posts
//
// # 高级查询
// curl http://localhost:8080/advanced/query
//
// # 事务
// curl -X POST http://localhost:8080/transaction
//
// ============================================================================

// ============================================================================
// 易错点总结
// ============================================================================
//
// 1. 【Find vs First】
//    Find: 未找到返回空切片，不报错
//    First: 未找到返回 ErrRecordNotFound
//
// 2. 【软删除的坑】
//    默认查询会自动加 WHERE deleted_at IS NULL
//    查询已删除记录: DB.Unscoped().Find(&users)
//    永久删除: DB.Unscoped().Delete(&user)
//
// 3. 【Update 零值问题】
//    Save: 更新所有字段（包括零值）
//    Updates(struct): 跳过零值字段
//    Updates(map): 包含零值
//    推荐用 map 更新部分字段
//
// 4. 【Preload N+1 问题】
//    不用 Preload: 每个用户单独查文章 = N+1 次查询
//    用 Preload: 只需 2 次查询
//
// 5. 【连接池配置】
//    生产环境必须配置连接池
//    MaxIdleConns < MaxOpenConns
//    ConnMaxLifetime 根据数据库超时设置
//
// 6. 【事务并发】
//    同一事务中的操作是串行的
//    不要在事务中做耗时操作（如调用外部 API）
//
// ============================================================================

// ============================================================================
// 练习题
// ============================================================================
//
// 1. 实现用户密码加密存储:
//    - 创建时用 bcrypt 加密
//    - 使用 GORM Hook (BeforeCreate)
//
// 2. 实现文章标签多对多关联:
//    - Post has many Tags
//    - Tag has many Posts
//    - 中间表 post_tags
//
// 3. 实现用户行为日志:
//    - 使用 GORM Hook (AfterCreate, AfterUpdate, AfterDelete)
//    - 记录操作类型、操作时间、操作人
//
// ============================================================================
