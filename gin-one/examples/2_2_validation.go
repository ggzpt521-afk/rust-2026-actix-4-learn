// ============================================================================
// 2.2 参数校验 (Validation)
// ============================================================================
// 运行方式: go run examples/2_2_validation.go
// ============================================================================

package main

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// ============================================================================
// 核心概念：Gin 的参数校验系统
// ============================================================================
//
// Gin 使用 go-playground/validator 库进行参数校验
//
// 【校验流程】
//
// 1. 请求到达 → 2. 模型绑定 → 3. 执行校验 → 4. 返回结果
//
// 校验规则写在 struct tag 的 binding 字段中
// 多个规则用逗号分隔: binding:"required,min=1,max=100"
//
// ============================================================================
//
// 【常用校验规则速查表】
//
// ┌─────────────────┬────────────────────────────────────────────────────────┐
// │ 规则             │ 说明                                                   │
// ├─────────────────┼────────────────────────────────────────────────────────┤
// │ required        │ 必填                                                    │
// │ omitempty       │ 空值时跳过后续校验                                       │
// ├─────────────────┼────────────────────────────────────────────────────────┤
// │ len=10          │ 长度等于 10 (string/slice/map)                          │
// │ min=5           │ 最小长度/最小值                                          │
// │ max=100         │ 最大长度/最大值                                          │
// │ eq=10           │ 等于                                                    │
// │ ne=0            │ 不等于                                                  │
// │ gt=0            │ 大于 (数字) / 晚于 (时间)                                │
// │ gte=0           │ 大于等于                                                │
// │ lt=100          │ 小于                                                    │
// │ lte=100         │ 小于等于                                                │
// ├─────────────────┼────────────────────────────────────────────────────────┤
// │ oneof=a b c     │ 枚举值 (空格分隔)                                        │
// │ contains=abc    │ 包含子串                                                │
// │ startswith=pre  │ 前缀                                                    │
// │ endswith=suf    │ 后缀                                                    │
// ├─────────────────┼────────────────────────────────────────────────────────┤
// │ email           │ 邮箱格式                                                │
// │ url             │ URL 格式                                                │
// │ uri             │ URI 格式                                                │
// │ uuid            │ UUID 格式                                               │
// │ ip              │ IP 地址 (v4 或 v6)                                      │
// │ ipv4            │ IPv4 地址                                               │
// │ ipv6            │ IPv6 地址                                               │
// │ datetime=layout │ 时间格式 (Go 时间格式)                                   │
// ├─────────────────┼────────────────────────────────────────────────────────┤
// │ alpha           │ 纯字母                                                  │
// │ alphanum        │ 字母+数字                                               │
// │ numeric         │ 数字字符串                                              │
// │ number          │ 数字 (包括小数)                                          │
// │ ascii           │ ASCII 字符                                              │
// │ json            │ JSON 格式字符串                                          │
// ├─────────────────┼────────────────────────────────────────────────────────┤
// │ eqfield=Field   │ 等于其他字段                                            │
// │ nefield=Field   │ 不等于其他字段                                          │
// │ gtfield=Field   │ 大于其他字段                                            │
// │ ltfield=Field   │ 小于其他字段                                            │
// └─────────────────┴────────────────────────────────────────────────────────┘
//
// ============================================================================

// ============================================================================
// 结构体定义：展示各种校验规则
// ============================================================================

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	// required: 必填
	// min/max: 长度限制
	Username string `json:"username" binding:"required,min=3,max=20,alphanum"`

	// email: 邮箱格式校验
	Email string `json:"email" binding:"required,email"`

	// min: 密码最小长度
	Password string `json:"password" binding:"required,min=8"`

	// eqfield: 必须等于 Password 字段
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`

	// gte/lte: 数值范围
	Age int `json:"age" binding:"required,gte=18,lte=120"`

	// oneof: 枚举值
	Gender string `json:"gender" binding:"required,oneof=male female other"`

	// 手机号: 自定义校验器 (后面实现)
	Phone string `json:"phone" binding:"required,phone"`

	// omitempty: 空值时跳过校验
	// url: URL 格式
	Website string `json:"website" binding:"omitempty,url"`
}

// ProductRequest 商品请求 (演示数值校验)
type ProductRequest struct {
	Name  string  `json:"name" binding:"required,min=1,max=200"`
	Price float64 `json:"price" binding:"required,gt=0,lte=999999.99"`
	Stock int     `json:"stock" binding:"gte=0"`

	// oneof 枚举
	Status string `json:"status" binding:"required,oneof=draft published archived"`

	// 数组校验: dive 深入校验每个元素
	Tags []string `json:"tags" binding:"max=5,dive,min=1,max=20"`

	// 嵌套结构体校验
	Details ProductDetails `json:"details" binding:"required"`
}

// ProductDetails 商品详情
type ProductDetails struct {
	Description string `json:"description" binding:"required,min=10,max=5000"`
	SKU         string `json:"sku" binding:"required,len=10,alphanum"`
}

// DateRangeRequest 日期范围请求
type DateRangeRequest struct {
	// datetime: 时间格式校验
	// ltfield: 必须小于 EndDate
	StartDate string `json:"start_date" binding:"required,datetime=2006-01-02"`
	EndDate   string `json:"end_date" binding:"required,datetime=2006-01-02,gtfield=StartDate"`
}

// ============================================================================
// 自定义校验器
// ============================================================================

// validatePhone 自定义手机号校验器
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	// 中国大陆手机号正则
	reg := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return reg.MatchString(phone)
}

// validateIDCard 自定义身份证校验器
func validateIDCard(fl validator.FieldLevel) bool {
	idCard := fl.Field().String()

	// 简单校验: 18位，最后一位可以是 X
	reg := regexp.MustCompile(`^\d{17}[\dXx]$`)
	return reg.MatchString(idCard)
}

// validatePassword 密码复杂度校验
// 至少包含: 大写字母、小写字母、数字、特殊字符中的三种
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var (
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber  = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial = regexp.MustCompile(`[!@#$%^&*]`).MatchString(password)
	)

	count := 0
	if hasUpper {
		count++
	}
	if hasLower {
		count++
	}
	if hasNumber {
		count++
	}
	if hasSpecial {
		count++
	}

	return count >= 3
}

// ============================================================================
// 错误信息处理
// ============================================================================

// ValidationError 校验错误
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// translateError 翻译校验错误为友好信息
func translateError(err error) []ValidationError {
	var errors []ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			var message string

			// 根据校验规则返回不同的错误信息
			switch e.Tag() {
			case "required":
				message = "该字段为必填项"
			case "email":
				message = "请输入有效的邮箱地址"
			case "min":
				if e.Type().Kind().String() == "string" {
					message = "长度至少为 " + e.Param() + " 个字符"
				} else {
					message = "值必须大于等于 " + e.Param()
				}
			case "max":
				if e.Type().Kind().String() == "string" {
					message = "长度最多为 " + e.Param() + " 个字符"
				} else {
					message = "值必须小于等于 " + e.Param()
				}
			case "gte":
				message = "值必须大于等于 " + e.Param()
			case "lte":
				message = "值必须小于等于 " + e.Param()
			case "gt":
				message = "值必须大于 " + e.Param()
			case "lt":
				message = "值必须小于 " + e.Param()
			case "oneof":
				message = "值必须是以下之一: " + e.Param()
			case "eqfield":
				message = "必须与 " + e.Param() + " 字段相同"
			case "phone":
				message = "请输入有效的手机号码"
			case "idcard":
				message = "请输入有效的身份证号码"
			case "password":
				message = "密码必须包含大写字母、小写字母、数字、特殊字符中的至少三种"
			case "alphanum":
				message = "只能包含字母和数字"
			case "url":
				message = "请输入有效的 URL"
			case "datetime":
				message = "日期格式不正确，应为 " + e.Param()
			case "len":
				message = "长度必须为 " + e.Param()
			case "gtfield":
				message = "必须大于 " + e.Param() + " 字段"
			case "ltfield":
				message = "必须小于 " + e.Param() + " 字段"
			default:
				message = "校验失败: " + e.Tag()
			}

			// 将字段名转为小写下划线格式
			field := toSnakeCase(e.Field())

			errors = append(errors, ValidationError{
				Field:   field,
				Message: message,
			})
		}
	}

	return errors
}

// toSnakeCase 驼峰转下划线
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

func main() {
	r := gin.Default()

	// ========================================================================
	// 注册自定义校验器
	// ========================================================================

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册自定义校验器
		v.RegisterValidation("phone", validatePhone)
		v.RegisterValidation("idcard", validateIDCard)
		v.RegisterValidation("password", validatePassword)
	}

	// ========================================================================
	// 一、基础校验示例
	// ========================================================================

	r.POST("/register", func(c *gin.Context) {
		var req RegisterRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			errors := translateError(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数校验失败",
				"errors":  errors,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "注册成功",
			"data":    req,
		})
	})

	// ========================================================================
	// 二、嵌套结构体 + 数组校验
	// ========================================================================

	r.POST("/products", func(c *gin.Context) {
		var req ProductRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			errors := translateError(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数校验失败",
				"errors":  errors,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "商品创建成功",
			"data":    req,
		})
	})

	// ========================================================================
	// 三、日期范围校验
	// ========================================================================

	r.POST("/reports/date-range", func(c *gin.Context) {
		var req DateRangeRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			errors := translateError(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数校验失败",
				"errors":  errors,
			})
			return
		}

		// 解析日期计算天数差
		startDate, _ := time.Parse("2006-01-02", req.StartDate)
		endDate, _ := time.Parse("2006-01-02", req.EndDate)
		days := int(endDate.Sub(startDate).Hours() / 24)

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "日期范围有效",
			"data": gin.H{
				"start_date": req.StartDate,
				"end_date":   req.EndDate,
				"days":       days,
			},
		})
	})

	// ========================================================================
	// 四、条件校验示例
	// ========================================================================

	type PaymentRequest struct {
		// 支付方式
		PaymentMethod string `json:"payment_method" binding:"required,oneof=credit_card bank_transfer alipay wechat"`

		// 信用卡信息 - 仅当 payment_method=credit_card 时必填
		// 注意: Gin 原生不支持条件校验，这里需要手动处理
		CardNumber string `json:"card_number"`
		ExpiryDate string `json:"expiry_date"`
		CVV        string `json:"cvv"`

		// 银行信息 - 仅当 payment_method=bank_transfer 时必填
		BankAccount string `json:"bank_account"`
		BankName    string `json:"bank_name"`
	}

	r.POST("/payments", func(c *gin.Context) {
		var req PaymentRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			errors := translateError(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数校验失败",
				"errors":  errors,
			})
			return
		}

		// 手动条件校验
		var fieldErrors []ValidationError

		switch req.PaymentMethod {
		case "credit_card":
			if req.CardNumber == "" {
				fieldErrors = append(fieldErrors, ValidationError{
					Field:   "card_number",
					Message: "信用卡号为必填项",
				})
			}
			if req.ExpiryDate == "" {
				fieldErrors = append(fieldErrors, ValidationError{
					Field:   "expiry_date",
					Message: "有效期为必填项",
				})
			}
			if req.CVV == "" {
				fieldErrors = append(fieldErrors, ValidationError{
					Field:   "cvv",
					Message: "CVV 为必填项",
				})
			}
		case "bank_transfer":
			if req.BankAccount == "" {
				fieldErrors = append(fieldErrors, ValidationError{
					Field:   "bank_account",
					Message: "银行账号为必填项",
				})
			}
			if req.BankName == "" {
				fieldErrors = append(fieldErrors, ValidationError{
					Field:   "bank_name",
					Message: "银行名称为必填项",
				})
			}
		}

		if len(fieldErrors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数校验失败",
				"errors":  fieldErrors,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "支付请求已提交",
		})
	})

	// ========================================================================
	// 五、展示所有校验规则
	// ========================================================================

	r.GET("/validation-demo", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "查看代码中的校验规则示例",
			"examples": gin.H{
				"required":    "binding:\"required\"",
				"email":       "binding:\"email\"",
				"min_max":     "binding:\"min=3,max=100\"",
				"range":       "binding:\"gte=1,lte=100\"",
				"enum":        "binding:\"oneof=a b c\"",
				"equal_field": "binding:\"eqfield=Password\"",
				"nested":      "binding:\"dive,required\"",
				"custom":      "binding:\"phone\"",
			},
		})
	})

	r.Run(":8080")
}

// ============================================================================
// 测试命令
// ============================================================================
//
// # 注册接口 - 成功
// curl -X POST http://localhost:8080/register \
//   -H "Content-Type: application/json" \
//   -d '{
//     "username": "testuser",
//     "email": "test@example.com",
//     "password": "password123",
//     "confirm_password": "password123",
//     "age": 25,
//     "gender": "male",
//     "phone": "13800138000"
//   }'
//
// # 注册接口 - 失败 (多个错误)
// curl -X POST http://localhost:8080/register \
//   -H "Content-Type: application/json" \
//   -d '{
//     "username": "ab",
//     "email": "invalid-email",
//     "password": "123",
//     "confirm_password": "456",
//     "age": 10,
//     "gender": "unknown",
//     "phone": "123"
//   }'
//
// # 商品接口 - 成功
// curl -X POST http://localhost:8080/products \
//   -H "Content-Type: application/json" \
//   -d '{
//     "name": "iPhone 15",
//     "price": 7999.00,
//     "stock": 100,
//     "status": "published",
//     "tags": ["电子", "手机"],
//     "details": {
//       "description": "这是一款非常棒的智能手机，功能强大。",
//       "sku": "IPHONE15PR"
//     }
//   }'
//
// # 日期范围 - 成功
// curl -X POST http://localhost:8080/reports/date-range \
//   -H "Content-Type: application/json" \
//   -d '{
//     "start_date": "2024-01-01",
//     "end_date": "2024-01-31"
//   }'
//
// # 日期范围 - 失败 (结束日期早于开始日期)
// curl -X POST http://localhost:8080/reports/date-range \
//   -H "Content-Type: application/json" \
//   -d '{
//     "start_date": "2024-01-31",
//     "end_date": "2024-01-01"
//   }'
//
// # 支付接口 - 信用卡
// curl -X POST http://localhost:8080/payments \
//   -H "Content-Type: application/json" \
//   -d '{
//     "payment_method": "credit_card",
//     "card_number": "4111111111111111",
//     "expiry_date": "12/25",
//     "cvv": "123"
//   }'
//
// # 支付接口 - 缺少信用卡信息
// curl -X POST http://localhost:8080/payments \
//   -H "Content-Type: application/json" \
//   -d '{
//     "payment_method": "credit_card"
//   }'
//
// ============================================================================

// ============================================================================
// 易错点总结
// ============================================================================
//
// 1. 【required 与 omitempty】
//    required: 必须有值
//    omitempty: 空值时跳过后续校验
//    两者不要同时使用！
//
// 2. 【数值类型的 min/max vs 字符串的 min/max】
//    string: min=3 表示最少 3 个字符
//    int: min=3 表示值 >= 3
//    自动根据类型判断
//
// 3. 【oneof 的值用空格分隔】
//    正确: oneof=a b c
//    错误: oneof=a,b,c
//
// 4. 【嵌套结构体默认不校验】
//    需要在嵌套字段上加 binding 标签
//    数组元素校验需要 dive 关键字
//
// 5. 【自定义校验器必须注册】
//    在 main() 开始时注册，否则会 panic
//    binding.Validator.Engine().(*validator.Validate).RegisterValidation(...)
//
// 6. 【字段名是 Go 结构体字段名】
//    错误信息中的 Field 是 Go 字段名 (如 ConfirmPassword)
//    而不是 JSON tag 名 (如 confirm_password)
//    需要手动转换
//
// 7. 【eqfield 使用 Go 字段名】
//    eqfield=Password (不是 eqfield=password)
//
// 8. 【datetime 使用 Go 时间格式】
//    datetime=2006-01-02 (Go 的参考时间格式)
//    不是 datetime=YYYY-MM-DD
//
// ============================================================================

// ============================================================================
// 练习题
// ============================================================================
//
// 1. 实现一个地址校验:
//    - province (必填，限定省份列表)
//    - city (必填)
//    - district (可选)
//    - detail (必填，10-200字符)
//    - postal_code (可选，6位数字)
//
// 2. 实现一个密码修改接口:
//    - old_password (必填)
//    - new_password (必填，8-20字符，复杂度校验)
//    - confirm_password (必填，等于 new_password)
//    - 新密码不能与旧密码相同
//
// 3. 实现一个时间段校验:
//    - start_time (必填，HH:mm 格式)
//    - end_time (必填，HH:mm 格式，必须大于 start_time)
//    - 时间跨度不能超过 8 小时
//
// ============================================================================
