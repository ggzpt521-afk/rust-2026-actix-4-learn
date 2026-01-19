// ============================================================================
// 10_errors.go - 错误处理
// ============================================================================
// 运行: go run 10_errors.go
//
// 【本文件学习目标】
// 1. 理解 Go 的错误处理哲学
// 2. 掌握 error 接口和自定义错误类型
// 3. 学会使用 errors.Is 和 errors.As 检查错误
// 4. 理解错误包装和错误链
// 5. 掌握 panic/recover 的正确使用场景
// 6. 了解 Go 1.13+ 和 Go 1.20+ 的错误处理新特性
//
// 【Go 错误处理哲学】
// - 错误是值（Errors are values）
// - 显式处理错误，不隐藏错误
// - 没有异常机制，使用返回值
// - 调用者决定如何处理错误
//
// 【error 接口定义】
// type error interface {
//     Error() string
// }
// 任何实现了 Error() string 方法的类型都是 error
//
// 【Go 错误处理 vs 其他语言】
// | 特性       | Go             | Java/Python      |
// |------------|----------------|------------------|
// | 机制       | 返回值         | 异常 (Exception)  |
// | 检查方式   | if err != nil  | try-catch        |
// | 强制处理   | 否（但推荐）   | 否（checked除外）|
// | 性能       | 无额外开销     | 栈展开有开销     |
// | 控制流     | 线性           | 跳转             |
// ============================================================================

package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

// ============================================================================
// 【自定义错误类型】
// ============================================================================
// 当需要携带额外信息时，定义自定义错误类型
//
// 【设计原则】
// - 包含足够的上下文信息
// - 实现 error 接口（Error() string）
// - 如果需要支持错误链，实现 Unwrap() error
// - 错误类型名通常以 Error 结尾
//
// 【何时使用自定义错误】
// - 需要携带额外字段（如错误码、字段名）
// - 需要程序化地检查错误类型
// - 需要实现特定的错误行为
// ============================================================================

// ValidationError: 简单的自定义错误
// 用于表示字段验证失败
// 【字段说明】
// - Field: 验证失败的字段名
// - Message: 错误描述信息
type ValidationError struct {
	Field   string // 验证失败的字段
	Message string // 错误消息
}

// Error: 实现 error 接口
// 【返回值格式】
// 应该是小写开头、无标点的简洁描述
// 方便与其他错误消息拼接
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// NetworkError: 带更多上下文的错误
// 用于表示网络操作失败
// 【字段说明】
// - Op: 操作类型（如 GET、POST）
// - URL: 请求的 URL
// - Err: 底层错误（支持错误链）
type NetworkError struct {
	Op  string // 操作名称
	URL string // 请求 URL
	Err error  // 底层错误（用于错误链）
}

// Error: 实现 error 接口
// 包含操作、URL 和底层错误的完整上下文
func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error during %s to %s: %v", e.Op, e.URL, e.Err)
}

// Unwrap: 实现 Unwrap 方法支持错误链
// 【作用】
// - 让 errors.Is 和 errors.As 可以检查底层错误
// - 支持错误链的遍历
//
// 【Go 1.13+ 错误链】
// 如果错误包含底层错误，实现 Unwrap() 方法
// 这样 errors.Is/As 会递归检查整个错误链
func (e *NetworkError) Unwrap() error {
	return e.Err
}

// ============================================================================
// 【哨兵错误（Sentinel Errors）】
// ============================================================================
// 预定义的错误值，用于表示已知的错误情况
//
// 【命名约定】
// - 以 Err 开头：ErrNotFound, ErrTimeout
// - 包级变量，可被外部包引用
//
// 【标准库中的哨兵错误】
// - io.EOF: 文件/流结束
// - io.ErrUnexpectedEOF: 意外的流结束
// - os.ErrNotExist: 文件不存在
// - os.ErrPermission: 权限不足
// - sql.ErrNoRows: 查询无结果
//
// 【使用场景】
// - 表示可预期的、调用者需要处理的错误
// - 使用 errors.Is 进行检查
// ============================================================================

// 业务错误常量（哨兵错误）
// 使用 errors.New 创建简单的错误值
var (
	ErrNotFound     = errors.New("resource not found")  // 资源未找到
	ErrUnauthorized = errors.New("unauthorized access") // 未授权访问
	ErrInvalidInput = errors.New("invalid input")       // 无效输入
)

func main() {
	fmt.Println("=== Go 错误处理 ===\n")

	// ========================================================================
	// 【基本错误处理】
	// ========================================================================
	// Go 的惯用模式：函数返回 (result, error)
	//
	// 【处理流程】
	// 1. 调用函数，接收返回值
	// 2. 检查 error 是否为 nil
	// 3. 如果不为 nil，处理错误
	// 4. 如果为 nil，使用结果
	//
	// 【为什么 error 是最后一个返回值？】
	// - Go 惯例：error 总是最后一个返回值
	// - 方便检查：if err != nil
	// - 链式调用时可以立即检查
	// ========================================================================
	fmt.Println("--- 基本错误处理 ---")

	// Go 函数通常返回 (result, error)
	// 【标准模式】
	// result, err := someFunc()
	// if err != nil {
	//     // 处理错误
	//     return err  // 或 return fmt.Errorf("context: %w", err)
	// }
	// // 使用 result
	result, err := divide(10, 2)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("10 / 2 = %.2f\n", result)
	}

	// 除零错误
	// 【注意】
	// - 不要忽略错误：_ , _ = divide(10, 0)  // 错误！
	// - 即使不需要结果，也要检查错误
	result, err = divide(10, 0)
	if err != nil {
		fmt.Printf("错误: %v\n", err) // division by zero
	}

	// ========================================================================
	// 【创建错误的方式】
	// ========================================================================
	// Go 提供多种创建错误的方式
	//
	// 【方式对比】
	// | 方式           | 用途                     | 示例                    |
	// |----------------|--------------------------|-------------------------|
	// | errors.New     | 简单静态错误             | errors.New("failed")    |
	// | fmt.Errorf     | 带格式化的错误           | fmt.Errorf("x=%d", x)   |
	// | fmt.Errorf %w  | 包装错误（保留原错误）   | fmt.Errorf("...: %w",e) |
	// | 自定义类型     | 需要额外信息             | &MyError{...}           |
	// ========================================================================
	fmt.Println("\n--- 创建错误的方式 ---")

	// 方式1: errors.New
	// 【特点】
	// - 创建简单的静态错误
	// - 每次调用创建新的错误值
	// - 适合包级哨兵错误
	err1 := errors.New("something went wrong")
	fmt.Printf("errors.New: %v\n", err1)

	// 方式2: fmt.Errorf（支持格式化）
	// 【特点】
	// - 可以包含动态信息
	// - 格式化语法与 fmt.Sprintf 相同
	// - 适合需要上下文信息的错误
	name := "config.json"
	err2 := fmt.Errorf("failed to load file: %s", name)
	fmt.Printf("fmt.Errorf: %v\n", err2)

	// 方式3: 错误包装（Go 1.13+）
	// 【%w 格式动词】
	// - 将原始错误包装到新错误中
	// - 保留原始错误，可通过 Unwrap 获取
	// - errors.Is/As 可以检查被包装的错误
	//
	// 【重要区别】
	// %v: 只包含错误文本，丢失原始错误
	// %w: 包装错误，保留原始错误（推荐）
	originalErr := errors.New("connection refused")
	wrappedErr := fmt.Errorf("database error: %w", originalErr)
	fmt.Printf("wrapped error: %v\n", wrappedErr)

	// ========================================================================
	// 【错误检查】
	// ========================================================================
	// Go 1.13+ 引入了 errors.Is 和 errors.As
	//
	// 【errors.Is vs ==】
	// - ==: 只比较当前错误
	// - errors.Is: 检查整个错误链
	//
	// 【errors.Is vs errors.As】
	// - errors.Is: 检查是否是特定错误值（用于哨兵错误）
	// - errors.As: 检查是否是特定错误类型（用于自定义错误）
	//
	// 【何时使用哪个】
	// | 场景                    | 使用           |
	// |-------------------------|----------------|
	// | 检查 io.EOF             | errors.Is      |
	// | 检查 os.ErrNotExist     | errors.Is      |
	// | 检查 *ValidationError   | errors.As      |
	// | 检查 *os.PathError      | errors.As      |
	// ========================================================================
	fmt.Println("\n--- 错误检查 ---")

	// errors.Is: 检查错误链中是否包含特定错误
	// 【工作原理】
	// 1. 比较 target == err
	// 2. 如果不等，尝试 err.Unwrap()
	// 3. 递归检查整个错误链
	if errors.Is(wrappedErr, originalErr) {
		fmt.Println("errors.Is: wrappedErr 包含 originalErr")
	}

	// 检查哨兵错误
	// 【最佳实践】
	// - 使用 errors.Is 而非 == 比较
	// - 因为错误可能被包装
	_, err = findUser("unknown")
	if errors.Is(err, ErrNotFound) {
		fmt.Println("errors.Is: 用户未找到")
	}

	// errors.As: 检查错误链中是否包含特定类型
	// 【语法】
	// var target *ErrorType
	// if errors.As(err, &target) {
	//     // 使用 target
	// }
	//
	// 【注意】
	// - 第二个参数必须是指向错误类型指针的指针（**ErrorType）
	// - 如果找到，会将值赋给 target
	validationErr := &ValidationError{Field: "email", Message: "invalid format"}
	wrappedValidation := fmt.Errorf("user creation failed: %w", validationErr)

	var ve *ValidationError
	if errors.As(wrappedValidation, &ve) {
		fmt.Printf("errors.As: 找到 ValidationError, field=%s\n", ve.Field)
	}

	// ========================================================================
	// 【自定义错误使用】
	// ========================================================================
	// 自定义错误可以携带额外的上下文信息
	// 调用者可以通过类型断言或 errors.As 获取详细信息
	// ========================================================================
	fmt.Println("\n--- 自定义错误 ---")

	err = validateUser("", 15)
	if err != nil {
		fmt.Printf("验证失败: %v\n", err)

		// 类型断言获取详细信息
		// 【使用 errors.As 而非直接类型断言】
		// - errors.As 可以检查错误链
		// - 直接类型断言只能检查当前错误
		var ve *ValidationError
		if errors.As(err, &ve) {
			fmt.Printf("  字段: %s\n", ve.Field)
			fmt.Printf("  消息: %s\n", ve.Message)
		}
	}

	// ========================================================================
	// 【错误链】
	// ========================================================================
	// 错误链是通过 Unwrap 方法连接的一系列错误
	//
	// 【构建错误链】
	// 1. 使用 fmt.Errorf("%w", err) 包装
	// 2. 自定义错误实现 Unwrap() error
	//
	// 【遍历错误链】
	// for e := err; e != nil; e = errors.Unwrap(e) {
	//     // 处理每个错误
	// }
	//
	// 【错误链的好处】
	// - 保留完整的错误上下文
	// - 可以检查链中的任何错误
	// - 调试时能看到完整的错误路径
	// ========================================================================
	fmt.Println("\n--- 错误链 ---")

	err = fetchData("https://api.example.com/data")
	if err != nil {
		fmt.Printf("获取数据失败: %v\n", err)

		// 解包错误链
		// 【errors.Unwrap】
		// - 返回被包装的错误
		// - 如果没有 Unwrap 方法，返回 nil
		fmt.Println("错误链:")
		for e := err; e != nil; e = errors.Unwrap(e) {
			fmt.Printf("  -> %v\n", e)
		}
	}

	// ========================================================================
	// 【多错误处理 (Go 1.20+)】
	// ========================================================================
	// errors.Join 可以将多个错误合并为一个
	//
	// 【语法】
	// err := errors.Join(err1, err2, err3)
	//
	// 【特点】
	// - 合并后的错误包含所有原始错误
	// - errors.Is 可以检查任何一个原始错误
	// - 错误消息用换行符分隔
	//
	// 【使用场景】
	// - 批量操作，需要收集所有错误
	// - 验证多个字段
	// - 并行操作的错误收集
	// ========================================================================
	fmt.Println("\n--- 多错误处理 (Go 1.20+) ---")

	err = processItems([]string{"valid", "invalid1", "invalid2"})
	if err != nil {
		fmt.Printf("处理失败: %v\n", err)
	}

	// ========================================================================
	// 【错误处理模式】
	// ========================================================================
	// Go 社区总结了几种常见的错误处理模式
	//
	// 【模式选择指南】
	// | 模式         | 适用场景                       |
	// |--------------|--------------------------------|
	// | 立即返回     | 大多数情况                     |
	// | 添加上下文   | 需要更多信息帮助调试           |
	// | 记录并继续   | 错误不影响主流程               |
	// | 重试         | 临时性错误（网络超时等）       |
	// | 降级         | 有备选方案                     |
	// ========================================================================
	fmt.Println("\n--- 错误处理模式 ---")

	// 模式1: 立即返回
	// 【最常见的模式】
	// if err := doSomething(); err != nil {
	//     return err
	// }
	fmt.Println("模式1: 立即返回 (最常见)")
	if err := doSomething(); err != nil {
		fmt.Printf("  错误: %v\n", err)
	}

	// 模式2: 添加上下文后返回
	// 【推荐用于底层错误】
	// if err := doSomething(); err != nil {
	//     return fmt.Errorf("操作X失败: %w", err)
	// }
	fmt.Println("模式2: 添加上下文")
	if err := processWithContext(); err != nil {
		fmt.Printf("  错误: %v\n", err)
	}

	// 模式3: 记录并继续
	// 【用于非关键错误】
	// if err := doSomething(); err != nil {
	//     log.Printf("警告: %v", err)
	//     // 继续执行
	// }
	fmt.Println("模式3: 记录并继续")
	processWithLogging()

	// ========================================================================
	// 【类型转换错误】
	// ========================================================================
	// strconv 包的函数返回特定的错误类型
	//
	// 【strconv.NumError 结构】
	// type NumError struct {
	//     Func string // 函数名（如 "Atoi"）
	//     Num  string // 输入字符串
	//     Err  error  // 底层错误（如 ErrSyntax, ErrRange）
	// }
	//
	// 【常见的 strconv 错误】
	// - strconv.ErrSyntax: 语法错误
	// - strconv.ErrRange: 超出范围
	// ========================================================================
	fmt.Println("\n--- 类型转换错误 ---")

	// strconv 错误
	_, err = strconv.Atoi("not a number")
	if err != nil {
		// 使用 errors.As 获取详细错误信息
		var numErr *strconv.NumError
		if errors.As(err, &numErr) {
			fmt.Printf("转换错误: Func=%s, Num=%s, Err=%v\n",
				numErr.Func, numErr.Num, numErr.Err)
		}
	}

	// ========================================================================
	// 【文件操作错误】
	// ========================================================================
	// os 包的函数返回 *os.PathError
	//
	// 【os.PathError 结构】
	// type PathError struct {
	//     Op   string // 操作（如 "open"）
	//     Path string // 文件路径
	//     Err  error  // 底层错误
	// }
	//
	// 【检查文件错误的两种方式】
	// 1. os.IsNotExist(err) - 旧方式，仍然可用
	// 2. errors.Is(err, os.ErrNotExist) - Go 1.13+ 推荐方式
	// ========================================================================
	fmt.Println("\n--- 文件操作错误 ---")

	_, err = os.Open("nonexistent.txt")
	if err != nil {
		// 旧方式：使用 os.IsNotExist
		// 【注意】这些函数仍然有效，但 errors.Is 更通用
		if os.IsNotExist(err) {
			fmt.Println("文件不存在")
		}

		// Go 1.13+ 方式：使用 errors.Is
		// 【推荐】更一致的 API，可检查错误链
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("errors.Is: 文件不存在")
		}
	}

	// ========================================================================
	// 【defer 与错误处理】
	// ========================================================================
	// defer 常用于确保资源清理
	//
	// 【常见模式】
	// f, err := os.Open(filename)
	// if err != nil {
	//     return err
	// }
	// defer f.Close()  // 确保文件关闭
	//
	// 【处理 Close 错误】
	// defer func() {
	//     if cerr := f.Close(); cerr != nil {
	//         // 处理关闭错误
	//     }
	// }()
	//
	// 【命名返回值模式】
	// func process() (err error) {
	//     defer func() {
	//         if cerr := cleanup(); cerr != nil && err == nil {
	//             err = cerr  // 只有主操作成功时才报告清理错误
	//         }
	//     }()
	//     // ...
	// }
	// ========================================================================
	fmt.Println("\n--- defer 与错误处理 ---")

	err = processFile("test.txt")
	if err != nil {
		fmt.Printf("处理文件错误: %v\n", err)
	}

	// ========================================================================
	// 【panic 与 recover】
	// ========================================================================
	// panic 和 recover 是 Go 的异常处理机制
	//
	// 【panic】
	// - 停止当前函数执行
	// - 执行所有 defer
	// - 向上传播到调用者
	// - 最终导致程序崩溃（如果不 recover）
	//
	// 【recover】
	// - 只能在 defer 函数中调用
	// - 捕获 panic 值
	// - 阻止 panic 继续传播
	//
	// 【何时使用 panic】
	// ✗ 不要用于普通错误（用 error 返回）
	// ✓ 程序初始化失败（无法继续）
	// ✓ 程序员错误（如数组越界）
	// ✓ 不可能发生的情况（断言失败）
	//
	// 【何时使用 recover】
	// ✓ 库代码中防止 panic 逃逸
	// ✓ HTTP 处理器中防止一个请求崩溃整个服务
	// ✓ goroutine 中防止 panic 影响其他 goroutine
	// ========================================================================
	fmt.Println("\n--- panic 与 recover ---")

	// 正常情况不应该使用 panic
	// panic 适用于不可恢复的错误

	// 使用 recover 捕获 panic
	// 【safeCall 模式】
	// 在 goroutine 或库函数中使用，防止 panic 传播
	result2 := safeCall(func() {
		panic("something terrible happened")
	})
	fmt.Printf("safeCall 结果: %v\n", result2)

	// ========================================================================
	// 【错误处理最佳实践】
	// ========================================================================
	// 1. 总是检查错误
	//    - 不要忽略：_, _ = doSomething()
	//    - 不知道怎么处理就返回给调用者
	//
	// 2. 使用 %w 包装错误
	//    - 保留原始错误信息
	//    - 添加有意义的上下文
	//    - 避免丢失错误链
	//
	// 3. 使用哨兵错误
	//    - 定义可预期的错误
	//    - 使用 errors.Is 检查
	//
	// 4. 使用自定义错误类型
	//    - 携带额外信息
	//    - 使用 errors.As 检查
	//
	// 5. 错误消息格式
	//    - 小写开头，无标点
	//    - 简洁但有信息量
	//    - 避免 "failed to" 等冗余前缀
	//
	// 6. 避免过度包装
	//    - 每层只添加新的上下文
	//    - 避免重复信息
	// ========================================================================
	fmt.Println("\n--- 错误处理最佳实践 ---")
	fmt.Println("1. 总是检查错误，不要忽略")
	fmt.Println("2. 使用 %w 包装错误以保留上下文")
	fmt.Println("3. 使用哨兵错误（如 io.EOF）进行已知错误检查")
	fmt.Println("4. 使用自定义错误类型携带额外信息")
	fmt.Println("5. 错误消息应该小写开头，不加标点")
	fmt.Println("6. 避免过度包装导致消息冗余")
}

// ============================================================================
// 【辅助函数】
// ============================================================================

// divide: 除法运算，演示基本错误返回
// 【返回值模式】
// - 成功：返回结果和 nil
// - 失败：返回零值和 error
//
// 【注意】
// - 即使出错，也要返回有意义的零值
// - 错误消息简洁明了
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero") // 小写开头，无标点
	}
	return a / b, nil
}

// findUser: 查找用户，演示哨兵错误的使用
// 【返回值】
// - 找到：返回用户结构体指针
// - 未找到：返回 nil 和 ErrNotFound
//
// 【调用者应该】
// _, err := findUser(name)
// if errors.Is(err, ErrNotFound) {
//     // 处理未找到的情况
// }
func findUser(name string) (*struct{ Name string }, error) {
	if name == "unknown" {
		return nil, ErrNotFound // 返回哨兵错误
	}
	return &struct{ Name string }{Name: name}, nil
}

// validateUser: 验证用户，演示自定义错误的使用
// 【返回自定义错误的好处】
// - 调用者可以知道具体哪个字段验证失败
// - 可以构建更友好的错误消息
func validateUser(name string, age int) error {
	if name == "" {
		return &ValidationError{Field: "name", Message: "cannot be empty"}
	}
	if age < 18 {
		return &ValidationError{Field: "age", Message: "must be at least 18"}
	}
	return nil
}

// fetchData: 模拟网络请求，演示错误链
// 【错误链结构】
// NetworkError
//   └── "retry failed: connection timeout"
//         └── "connection timeout"
//
// 【好处】
// - 完整的错误上下文
// - 可以检查链中的任何错误
func fetchData(url string) error {
	// 模拟网络错误
	innerErr := errors.New("connection timeout")
	return &NetworkError{
		Op:  "GET",
		URL: url,
		Err: fmt.Errorf("retry failed: %w", innerErr), // 包装底层错误
	}
}

// processItems: 处理多个项目，演示 errors.Join (Go 1.20+)
// 【使用场景】
// - 批量操作需要收集所有错误
// - 不想因为一个失败而停止处理
//
// 【errors.Join 特点】
// - 返回的错误包含所有原始错误
// - errors.Is 可以检查任何一个
// - nil 错误会被忽略
func processItems(items []string) error {
	var errs []error
	for _, item := range items {
		if item != "valid" {
			errs = append(errs, fmt.Errorf("invalid item: %s", item))
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...) // 合并所有错误
	}
	return nil
}

// doSomething: 模拟可能失败的操作
func doSomething() error {
	return errors.New("operation failed")
}

// processWithContext: 演示添加上下文的错误处理模式
// 【模式】
// 调用底层函数，失败时添加上下文后返回
// 这样错误消息会包含完整的调用链信息
func processWithContext() error {
	err := doSomething()
	if err != nil {
		// 使用 %w 包装，保留原始错误
		return fmt.Errorf("processing failed: %w", err)
	}
	return nil
}

// processWithLogging: 演示记录并继续的错误处理模式
// 【适用场景】
// - 错误不是致命的
// - 有降级或默认行为
// - 不想因为一个错误停止整个流程
func processWithLogging() {
	err := doSomething()
	if err != nil {
		// 记录错误但继续
		fmt.Printf("  [WARN] 操作失败但继续: %v\n", err)
	}
	fmt.Println("  继续执行其他操作...")
}

// processFile: 演示 defer 与错误处理
// 【命名返回值】
// 使用命名返回值 err，可以在 defer 中修改它
//
// 【资源清理模式】
// 1. 获取资源
// 2. 检查错误
// 3. defer 清理
// 4. 使用资源
func processFile(filename string) (err error) {
	// 使用 defer 确保资源清理
	// 实际应用中这里会打开文件
	fmt.Printf("  打开文件: %s\n", filename)

	defer func() {
		fmt.Println("  关闭文件（defer）")
		// 可以在这里处理关闭错误
		// 如果需要，可以修改命名返回值 err
	}()

	// 模拟处理
	// 返回 io.EOF 包装的错误
	return fmt.Errorf("read error: %w", io.EOF)
}

// safeCall: 安全调用函数，捕获 panic
// 【recover 使用要点】
// 1. 必须在 defer 函数中调用
// 2. 只能捕获当前 goroutine 的 panic
// 3. 返回 panic 的值（interface{} 类型）
//
// 【典型用法】
// - HTTP 中间件：防止一个请求崩溃整个服务
// - goroutine 包装：防止 panic 传播
// - 库函数：防止 panic 逃逸给调用者
func safeCall(f func()) (err error) {
	defer func() {
		// recover 捕获 panic
		// 如果没有 panic，返回 nil
		if r := recover(); r != nil {
			// 将 panic 值转换为 error
			err = fmt.Errorf("panic recovered: %v", r)
		}
	}()
	f()      // 调用可能 panic 的函数
	return nil // 正常返回
}
