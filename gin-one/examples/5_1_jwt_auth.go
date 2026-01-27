// ============================================================================
// 5.1 JWT 身份认证实战
// ============================================================================
// 运行方式: go run examples/5_1_jwt_auth.go
// 需要先安装: go get -u github.com/golang-jwt/jwt/v5
// ============================================================================

package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ============================================================================
// JWT 核心概念
// ============================================================================
//
// 【什么是 JWT？】
//
// JSON Web Token - 一种开放标准 (RFC 7519)
// 用于在各方之间安全地传输信息
//
// 【JWT 结构】
//
// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
//
// 由三部分组成，用 . 分隔：
//
// 1. Header (头部) - 算法和类型
//    {"alg": "HS256", "typ": "JWT"}
//
// 2. Payload (载荷) - 数据
//    {"sub": "1234567890", "name": "John Doe", "iat": 1516239022}
//
// 3. Signature (签名) - 验证完整性
//    HMACSHA256(base64UrlEncode(header) + "." + base64UrlEncode(payload), secret)
//
// ============================================================================
//
// 【JWT vs Session 对比】
//
// ┌─────────────────┬───────────────────────────────────────────────────────┐
// │ 方式             │ JWT                    │ Session                      │
// ├─────────────────┼────────────────────────┼──────────────────────────────┤
// │ 存储位置         │ 客户端 (localStorage)   │ 服务端 (Redis/内存)          │
// │ 扩展性           │ 天然支持分布式          │ 需要共享 Session 存储        │
// │ 带宽             │ 每次请求都带 Token      │ 只传 Session ID              │
// │ 安全性           │ 需防 XSS               │ 需防 CSRF                    │
// │ 时效控制         │ 无法主动失效            │ 可随时删除                   │
// └─────────────────┴────────────────────────┴──────────────────────────────┘
//
// ============================================================================

// ============================================================================
// 配置
// ============================================================================

var (
	// 【重要】生产环境必须从配置/环境变量读取
	JWTSecret          = []byte("your-super-secret-key-must-be-at-least-32-bytes")
	AccessTokenExpire  = 2 * time.Hour  // Access Token 有效期
	RefreshTokenExpire = 7 * 24 * time.Hour // Refresh Token 有效期
)

// ============================================================================
// Claims 定义
// ============================================================================

// CustomClaims 自定义 Claims
type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// ============================================================================
// Token 生成
// ============================================================================

// GenerateToken 生成 JWT Token
func GenerateToken(userID uint, username, role string) (accessToken, refreshToken string, err error) {
	// 创建 Access Token
	accessClaims := CustomClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "gin-app",
			Subject:   "access_token",
		},
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString(JWTSecret)
	if err != nil {
		return "", "", err
	}

	// 创建 Refresh Token（只包含 UserID）
	refreshClaims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "refresh_token",
		},
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString(JWTSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ParseToken 解析 JWT Token
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ============================================================================
// JWT 中间件
// ============================================================================

// JWTAuthMiddleware JWT 认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Authorization header is required",
			})
			return
		}

		// 解析 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Authorization header format must be: Bearer {token}",
			})
			return
		}

		tokenString := parts[1]

		// 解析 Token
		claims, err := ParseToken(tokenString)
		if err != nil {
			// 区分错误类型
			message := "Invalid token"
			if errors.Is(err, jwt.ErrTokenExpired) {
				message = "Token has expired"
			} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
				message = "Token not valid yet"
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": message,
				"error":   err.Error(),
			})
			return
		}

		// 检查是否是 Access Token
		if claims.Subject != "access_token" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid token type",
			})
			return
		}

		// 将用户信息存入 Context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("claims", claims)

		c.Next()
	}
}

// RoleMiddleware 角色权限中间件
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")

		allowed := false
		for _, r := range allowedRoles {
			if role == r {
				allowed = true
				break
			}
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "Permission denied",
			})
			return
		}

		c.Next()
	}
}

// ============================================================================
// 模拟用户数据
// ============================================================================

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // 不返回密码
	Role     string `json:"role"`
}

// 模拟数据库
var users = map[string]*User{
	"admin": {ID: 1, Username: "admin", Password: "admin123", Role: "admin"},
	"user":  {ID: 2, Username: "user", Password: "user123", Role: "user"},
}

// Token 黑名单（生产环境应该用 Redis）
var tokenBlacklist = make(map[string]bool)

// ============================================================================
// 主程序
// ============================================================================

func main() {
	r := gin.Default()

	// ========================================================================
	// 公开接口
	// ========================================================================

	// 登录
	r.POST("/login", func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 验证用户
		user, exists := users[req.Username]
		if !exists || user.Password != req.Password {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid username or password",
			})
			return
		}

		// 生成 Token
		accessToken, refreshToken, err := GenerateToken(user.ID, user.Username, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "Login successful",
			"data": gin.H{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
				"token_type":    "Bearer",
				"expires_in":    AccessTokenExpire.Seconds(),
			},
		})
	})

	// 刷新 Token
	r.POST("/refresh", func(c *gin.Context) {
		var req struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 检查是否在黑名单
		if tokenBlacklist[req.RefreshToken] {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Token has been revoked",
			})
			return
		}

		// 解析 Refresh Token
		claims, err := ParseToken(req.RefreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid refresh token",
			})
			return
		}

		// 检查 Token 类型
		if claims.Subject != "refresh_token" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid token type",
			})
			return
		}

		// 获取用户信息（实际应该从数据库查询）
		var user *User
		for _, u := range users {
			if u.ID == claims.UserID {
				user = u
				break
			}
		}

		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "User not found",
			})
			return
		}

		// 生成新的 Token
		accessToken, refreshToken, err := GenerateToken(user.ID, user.Username, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		// 将旧的 Refresh Token 加入黑名单
		tokenBlacklist[req.RefreshToken] = true

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "Token refreshed",
			"data": gin.H{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
				"token_type":    "Bearer",
				"expires_in":    AccessTokenExpire.Seconds(),
			},
		})
	})

	// ========================================================================
	// 需要认证的接口
	// ========================================================================

	authorized := r.Group("/api")
	authorized.Use(JWTAuthMiddleware())
	{
		// 获取当前用户信息
		authorized.GET("/me", func(c *gin.Context) {
			claims, _ := c.Get("claims")
			customClaims := claims.(*CustomClaims)

			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"data": gin.H{
					"user_id":  customClaims.UserID,
					"username": customClaims.Username,
					"role":     customClaims.Role,
				},
			})
		})

		// 登出
		authorized.POST("/logout", func(c *gin.Context) {
			// 获取当前 Token 并加入黑名单
			authHeader := c.GetHeader("Authorization")
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 {
				tokenBlacklist[parts[1]] = true
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "Logout successful",
			})
		})

		// 普通用户和管理员都可以访问
		authorized.GET("/profile", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "Profile data",
				"user_id": c.GetUint("user_id"),
			})
		})
	}

	// ========================================================================
	// 需要管理员权限的接口
	// ========================================================================

	admin := r.Group("/admin")
	admin.Use(JWTAuthMiddleware())
	admin.Use(RoleMiddleware("admin"))
	{
		admin.GET("/users", func(c *gin.Context) {
			var userList []User
			for _, u := range users {
				userList = append(userList, *u)
			}
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"data": userList,
			})
		})

		admin.DELETE("/users/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "User deleted (simulated)",
			})
		})
	}

	// 打印测试说明
	println("Server starting on :8080")
	println("")
	println("Test accounts:")
	println("  admin / admin123 (role: admin)")
	println("  user / user123 (role: user)")
	println("")
	println("Test commands:")
	println("")
	println("# Login")
	println(`curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123"}'`)
	println("")
	println("# Access protected resource")
	println(`curl http://localhost:8080/api/me -H "Authorization: Bearer <access_token>"`)
	println("")
	println("# Admin only")
	println(`curl http://localhost:8080/admin/users -H "Authorization: Bearer <access_token>"`)

	r.Run(":8080")
}

// ============================================================================
// 测试命令
// ============================================================================
//
// # 登录获取 Token
// curl -X POST http://localhost:8080/login \
//   -H "Content-Type: application/json" \
//   -d '{"username":"admin","password":"admin123"}'
//
// # 访问受保护接口
// curl http://localhost:8080/api/me \
//   -H "Authorization: Bearer <access_token>"
//
// # 刷新 Token
// curl -X POST http://localhost:8080/refresh \
//   -H "Content-Type: application/json" \
//   -d '{"refresh_token":"<refresh_token>"}'
//
// # 登出
// curl -X POST http://localhost:8080/api/logout \
//   -H "Authorization: Bearer <access_token>"
//
// # 管理员接口（需要 admin 角色）
// curl http://localhost:8080/admin/users \
//   -H "Authorization: Bearer <admin_access_token>"
//
// # 用普通用户访问管理员接口（会返回 403）
// curl http://localhost:8080/admin/users \
//   -H "Authorization: Bearer <user_access_token>"
//
// ============================================================================

// ============================================================================
// 易错点总结
// ============================================================================
//
// 1. 【Secret 太短】
//    HS256 的 Secret 至少 32 字节
//    否则容易被暴力破解
//
// 2. 【Token 存储位置】
//    浏览器: 推荐 httpOnly Cookie，防 XSS
//    移动端: 安全存储 (Keychain/Keystore)
//    不要存 localStorage (容易被 XSS 攻击)
//
// 3. 【无法主动失效】
//    JWT 一旦签发就无法撤销
//    解决方案:
//    - Token 黑名单 (需要存储)
//    - 短有效期 + Refresh Token
//
// 4. 【Token 泄露】
//    Access Token 泄露只影响短期
//    Refresh Token 泄露影响大，需要安全存储
//
// 5. 【时钟偏差】
//    分布式系统中服务器时钟可能不同步
//    可以在验证时加入一定的容错时间
//
// 6. 【Payload 不加密】
//    JWT 的 Payload 只是 Base64 编码，不是加密
//    不要存敏感信息（如密码）
//
// ============================================================================

// ============================================================================
// 练习题
// ============================================================================
//
// 1. 实现 Token 黑名单存储到 Redis
//
// 2. 实现多设备登录管理:
//    - 记录每个 Token 的设备信息
//    - 支持踢出指定设备
//
// 3. 实现 Token 自动续期:
//    - Access Token 快过期时自动刷新
//    - 返回新 Token 在响应头中
//
// ============================================================================
