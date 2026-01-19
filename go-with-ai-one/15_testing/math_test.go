// ============================================================================
// 15_testing - 单元测试
// ============================================================================
// 运行: cd 15_testing && go test -v
// 或者: go test -v ./15_testing/...
//
// 【本文件学习目标】
// 1. 掌握 Go 测试的基本语法和命名规范
// 2. 学会编写表格驱动测试（Table-Driven Tests）
// 3. 理解子测试（Subtests）的使用
// 4. 掌握基准测试（Benchmarks）的编写
// 5. 了解模糊测试（Fuzz Testing）
// 6. 学会使用测试辅助函数
//
// 【Go 测试文件规范】
// - 文件名以 _test.go 结尾
// - 与被测试文件在同一目录
// - 包名相同（或加 _test 后缀）
//
// 【测试函数类型】
// | 类型       | 函数签名                      | 作用               |
// |------------|-------------------------------|--------------------|
// | 单元测试   | TestXxx(t *testing.T)         | 功能测试           |
// | 基准测试   | BenchmarkXxx(b *testing.B)    | 性能测试           |
// | 示例测试   | ExampleXxx()                  | 文档示例           |
// | 模糊测试   | FuzzXxx(f *testing.F)         | 随机输入测试       |
//
// 【testing.T 常用方法】
// | 方法           | 作用                           |
// |----------------|--------------------------------|
// | t.Error/Errorf | 报告失败但继续执行             |
// | t.Fatal/Fatalf | 报告失败并立即停止当前测试     |
// | t.Skip/Skipf   | 跳过测试                       |
// | t.Run          | 运行子测试                     |
// | t.Parallel     | 标记为可并行执行               |
// | t.Helper       | 标记为辅助函数（错误定位更准确）|
// ============================================================================
package testing_demo

import (
	"testing"
)

// ============================================================================
// 【基本测试】
// ============================================================================
// 最简单的测试形式
//
// 【命名规范】
// - 函数名以 Test 开头（首字母大写）
// - Test 后面的名字首字母大写
// - 例如：TestAdd, TestDivide, TestIsPrime
//
// 【测试逻辑】
// 1. 调用被测函数
// 2. 比较实际结果和期望结果
// 3. 不相等则报告错误
// ============================================================================

// TestAdd: 基本的单元测试
// 【参数】*testing.T 提供测试控制和报告功能
func TestAdd(t *testing.T) {
	result := Add(2, 3) // 调用被测函数
	expected := 5       // 期望结果

	// 【断言】比较实际和期望
	if result != expected {
		// t.Errorf 报告错误但继续执行其他测试
		// 格式：got X; want Y（Go 社区惯例）
		t.Errorf("Add(2, 3) = %d; want %d", result, expected)
	}
}

// TestSubtract: 另一个基本测试
func TestSubtract(t *testing.T) {
	result := Subtract(5, 3)
	expected := 2

	if result != expected {
		t.Errorf("Subtract(5, 3) = %d; want %d", result, expected)
	}
}

// ============================================================================
// 【表格驱动测试】(Table-Driven Tests)
// ============================================================================
// Go 最推荐的测试模式！
//
// 【优点】
// - 易于添加新的测试用例
// - 代码复用，减少重复
// - 测试用例一目了然
// - 便于覆盖边界条件
//
// 【结构】
// 1. 定义测试用例结构体切片
// 2. 遍历测试用例
// 3. 使用 t.Run 运行子测试
// ============================================================================

// TestAddTable: 表格驱动测试示例
func TestAddTable(t *testing.T) {
	// 【测试用例表格】
	// 每个结构体代表一个测试用例
	tests := []struct {
		name     string // 测试用例名称（便于识别失败的用例）
		a, b     int    // 输入
		expected int    // 期望输出
	}{
		// 【设计测试用例的思路】
		// 1. 正常情况
		// 2. 边界值
		// 3. 特殊值
		// 4. 错误情况
		{"positive numbers", 2, 3, 5},    // 正数
		{"negative numbers", -2, -3, -5}, // 负数
		{"mixed numbers", -2, 3, 1},      // 正负混合
		{"zero", 0, 0, 0},                // 零
		{"with zero", 5, 0, 5},           // 带零
	}

	// 【遍历测试用例】
	for _, tt := range tests {
		// 【t.Run】运行子测试
		// 第一个参数是子测试名称（会显示在输出中）
		// 第二个参数是测试函数
		t.Run(tt.name, func(t *testing.T) {
			result := Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// TestIsPrimeTable: 素数测试（表格驱动）
func TestIsPrimeTable(t *testing.T) {
	// 【简化的表格】当测试用例简单时可以省略 name
	tests := []struct {
		input    int
		expected bool
	}{
		{0, false},   // 0 不是素数
		{1, false},   // 1 不是素数
		{2, true},    // 2 是最小的素数
		{3, true},    // 3 是素数
		{4, false},   // 4 = 2 * 2
		{5, true},    // 5 是素数
		{10, false},  // 10 = 2 * 5
		{17, true},   // 17 是素数
		{100, false}, // 100 = 10 * 10
	}

	for _, tt := range tests {
		// 当没有 name 时可以用空字符串
		t.Run("", func(t *testing.T) {
			result := IsPrime(tt.input)
			if result != tt.expected {
				t.Errorf("IsPrime(%d) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// 【错误测试】
// ============================================================================
// 测试返回 error 的函数
// 需要同时测试成功和失败的情况
// ============================================================================

// TestDivide: 测试带错误返回的函数
func TestDivide(t *testing.T) {
	// 【正常情况】
	result, err := Divide(10, 2)
	if err != nil {
		// t.Fatalf 报告致命错误并停止当前测试
		// 用于：错误发生后继续测试没有意义的情况
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 5 {
		t.Errorf("Divide(10, 2) = %d; want 5", result)
	}

	// 【错误情况】
	_, err = Divide(10, 0)
	if err == nil {
		// 期望有错误但没有
		t.Error("Divide(10, 0) should return error")
	}
}

// ============================================================================
// 【子测试】(Subtests)
// ============================================================================
// 使用 t.Run 组织相关的测试
//
// 【优点】
// - 逻辑分组
// - 可以单独运行：go test -run TestMath/Addition
// - 可以共享 setup/teardown
// - 输出更清晰
// ============================================================================

// TestMath: 使用子测试组织相关测试
func TestMath(t *testing.T) {
	// 【子测试】每个 t.Run 是一个独立的子测试
	t.Run("Addition", func(t *testing.T) {
		if Add(1, 1) != 2 {
			t.Error("1 + 1 should equal 2")
		}
	})

	t.Run("Subtraction", func(t *testing.T) {
		if Subtract(2, 1) != 1 {
			t.Error("2 - 1 should equal 1")
		}
	})

	t.Run("Division", func(t *testing.T) {
		result, _ := Divide(10, 2)
		if result != 5 {
			t.Error("10 / 2 should equal 5")
		}
	})
}

// ============================================================================
// 【跳过测试】
// ============================================================================
// 有时需要跳过某些测试：
// - 依赖外部服务
// - 运行时间太长
// - 特定环境才能运行
// ============================================================================

// TestSkipped: 演示跳过测试
func TestSkipped(t *testing.T) {
	// 【t.Skip】立即跳过测试
	// 测试会被标记为 SKIP，不是失败
	t.Skip("跳过这个测试")

	// Skip 后的代码不会执行
	if Add(1, 1) != 2 {
		t.Error("failed")
	}
}

// TestConditionalSkip: 条件跳过
func TestConditionalSkip(t *testing.T) {
	// 【testing.Short()】检查是否运行在 short 模式
	// 运行：go test -short
	if testing.Short() {
		t.Skip("在 short 模式下跳过")
	}
	// 这里可以放长时间运行的测试
}

// ============================================================================
// 【并行测试】
// ============================================================================
// 标记测试为可并行执行，加快测试速度
//
// 【使用场景】
// - 测试之间没有依赖
// - 不共享可变状态
// - 不依赖执行顺序
//
// 【注意】
// - t.Parallel() 必须在测试开始时调用
// - 并行测试不会与顺序测试同时运行
// ============================================================================

// TestParallel: 标记为可并行执行
func TestParallel(t *testing.T) {
	t.Parallel() // 标记为可并行执行

	result := Add(1, 2)
	if result != 3 {
		t.Errorf("expected 3, got %d", result)
	}
}

// TestParallelGroup: 并行运行的子测试组
func TestParallelGroup(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"test1", 1, 2, 3},
		{"test2", 2, 3, 5},
		{"test3", 3, 4, 7},
	}

	for _, tt := range tests {
		// 【重要】Go 1.22 之前，必须在循环内捕获变量！
		// 否则所有 goroutine 会共享同一个 tt 变量
		tt := tt // 创建副本（Go 1.22+ 不需要）
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // 子测试并行执行
			result := Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// 【测试辅助函数】
// ============================================================================
// 抽取重复的测试逻辑
// 使用 t.Helper() 标记，让错误定位更准确
// ============================================================================

// assertEqual: 测试辅助函数
// 【t.Helper()】标记为辅助函数
// 作用：当断言失败时，错误位置指向调用者，而不是这个函数内部
func assertEqual(t *testing.T, got, want int) {
	t.Helper() // 关键！标记为辅助函数
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

// TestWithHelper: 使用辅助函数的测试
func TestWithHelper(t *testing.T) {
	// 使用辅助函数简化测试代码
	assertEqual(t, Add(2, 3), 5)
	assertEqual(t, Subtract(5, 3), 2)
}

// ============================================================================
// 【示例测试】(Example Tests)
// ============================================================================
// 示例测试有双重作用：
// 1. 作为文档出现在 godoc 中
// 2. 作为测试验证输出
//
// 【命名规范】
// ExampleFuncName      - 函数示例
// ExampleTypeName      - 类型示例
// ExampleTypeName_MethodName - 方法示例
//
// 【Output 注释】
// // Output: 期望的输出
// 如果实际输出与注释不符，测试失败
// ============================================================================

// ExampleAdd: Add 函数的示例
// 这会出现在 godoc 文档中
func ExampleAdd() {
	result := Add(2, 3)
	println(result)
	// Output: 5
}

// ExampleReverseString: ReverseString 函数的示例
func ExampleReverseString() {
	result := ReverseString("Hello")
	println(result)
	// Output: olleH
}

// ============================================================================
// 【基准测试】(Benchmarks)
// ============================================================================
// 测量函数的性能
// 运行: go test -bench=. -benchmem
//
// 【函数签名】
// BenchmarkXxx(b *testing.B)
//
// 【关键】
// - 循环 b.N 次
// - b.N 由测试框架自动调整
// - 直到测量结果稳定
//
// 【输出解读】
// BenchmarkAdd-8    1000000000    0.3 ns/op    0 B/op    0 allocs/op
// - -8: 使用 8 个 CPU
// - 1000000000: 运行次数
// - 0.3 ns/op: 每次操作耗时
// - 0 B/op: 每次操作分配的字节
// - 0 allocs/op: 每次操作的内存分配次数
// ============================================================================

// BenchmarkAdd: Add 函数的基准测试
func BenchmarkAdd(b *testing.B) {
	// 【关键】必须循环 b.N 次
	// b.N 由框架自动调整，确保测量足够准确
	for i := 0; i < b.N; i++ {
		Add(1, 2)
	}
}

// BenchmarkFibonacci: Fibonacci 函数的基准测试
func BenchmarkFibonacci(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fibonacci(10)
	}
}

// BenchmarkReverseString: ReverseString 函数的基准测试
func BenchmarkReverseString(b *testing.B) {
	s := "Hello, World!"
	for i := 0; i < b.N; i++ {
		ReverseString(s)
	}
}

// BenchmarkFibonacciSub: 带子基准测试
// 可以比较不同输入的性能
func BenchmarkFibonacciSub(b *testing.B) {
	// 【子基准测试表格】
	benchmarks := []struct {
		name string
		n    int
	}{
		{"Fib5", 5},   // 小输入
		{"Fib10", 10}, // 中等输入
		{"Fib15", 15}, // 较大输入
	}

	for _, bm := range benchmarks {
		// 【b.Run】运行子基准测试
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Fibonacci(bm.n)
			}
		})
	}
}

// ============================================================================
// 【模糊测试】(Fuzz Testing) - Go 1.18+
// ============================================================================
// 自动生成随机输入来发现边界情况
// 运行: go test -fuzz=FuzzReverse
//
// 【函数签名】
// FuzzXxx(f *testing.F)
//
// 【工作原理】
// 1. 提供种子语料（seed corpus）
// 2. 框架基于种子生成随机输入
// 3. 测试函数检验输入是否导致崩溃或错误
//
// 【用途】
// - 发现边界情况
// - 安全测试
// - 验证不变量（invariants）
// ============================================================================

// FuzzReverse: 反转字符串的模糊测试
func FuzzReverse(f *testing.F) {
	// 【种子语料】初始的测试输入
	// 框架会基于这些种子变异生成更多输入
	f.Add("hello")  // ASCII
	f.Add("world")  // ASCII
	f.Add("Go语言") // Unicode

	// 【模糊测试函数】
	// t: 用于报告错误
	// s: 框架生成的随机输入
	f.Fuzz(func(t *testing.T, s string) {
		// 【测试不变量】
		// 双重反转应该等于原字符串
		rev := ReverseString(s)
		doubleRev := ReverseString(rev)
		if s != doubleRev {
			t.Errorf("双重反转后不相等: %q -> %q -> %q", s, rev, doubleRev)
		}
	})
}

/*
=== 测试命令汇总 ===

【基本命令】
  go test              # 运行当前包的测试
  go test ./...        # 运行所有包的测试
  go test -v           # 详细输出（显示每个测试的结果）
  go test -run TestAdd # 运行匹配 "TestAdd" 的测试
  go test -run Test/Addition # 运行子测试

【基准测试】
  go test -bench=.           # 运行所有基准测试
  go test -bench=BenchmarkAdd
  go test -bench=. -benchmem # 显示内存分配信息
  go test -bench=. -count=5  # 运行 5 次取平均

【覆盖率】
  go test -cover                           # 显示覆盖率百分比
  go test -coverprofile=coverage.out       # 生成覆盖率数据文件
  go tool cover -html=coverage.out         # 生成 HTML 报告
  go tool cover -func=coverage.out         # 按函数显示覆盖率

【模糊测试】
  go test -fuzz=FuzzReverse         # 运行模糊测试
  go test -fuzz=FuzzReverse -fuzztime=30s  # 限制时间

【其他选项】
  go test -short         # 短模式（可以跳过长测试）
  go test -race          # 检测数据竞争
  go test -timeout 30s   # 设置超时时间
  go test -count 3       # 每个测试运行 3 次
  go test -failfast      # 第一个失败后立即停止
  go test -shuffle=on    # 随机顺序运行测试

【常用组合】
  go test -v -race -cover ./...  # 详细输出 + 竞争检测 + 覆盖率
  go test -bench=. -benchmem -count=3  # 基准测试完整输出
*/
