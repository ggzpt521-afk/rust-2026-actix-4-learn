// ============================================================================
// 11_generics.go - 泛型 (Go 1.18+)
// ============================================================================
// 运行: go run 11_generics.go
//
// 【本文件学习目标】
// 1. 理解泛型的基本概念和语法
// 2. 掌握类型参数和类型约束
// 3. 学会定义和使用泛型函数
// 4. 学会定义和使用泛型类型
// 5. 理解 ~ 操作符和类型近似
// 6. 掌握标准库提供的约束
//
// 【泛型的核心概念】
// - 类型参数（Type Parameter）：用方括号声明，如 [T any]
// - 类型约束（Type Constraint）：限制类型参数可以是哪些类型
// - 类型推断（Type Inference）：编译器自动推断类型参数
// - 实例化（Instantiation）：用具体类型替换类型参数
//
// 【为什么需要泛型？】
// - 代码复用：一份代码适用于多种类型
// - 类型安全：编译时检查类型错误
// - 性能：无运行时开销（不像 interface{} 需要类型断言）
//
// 【泛型 vs interface{}】
// | 特性       | 泛型              | interface{}        |
// |------------|-------------------|--------------------|
// | 类型检查   | 编译时            | 运行时             |
// | 性能       | 无额外开销        | 装箱/拆箱开销      |
// | 灵活性     | 约束内的类型      | 任意类型           |
// | 错误发现   | 编译错误          | 运行时 panic       |
//
// 【Go 1.18+ 引入的泛型特性】
// - 泛型函数
// - 泛型类型（结构体、接口）
// - 类型约束
// - any 和 comparable 内置约束
// - ~ 波浪线操作符（类型近似）
// ============================================================================

package main

import (
	"cmp" // Go 1.21+ 提供的比较包
	"fmt"
)

// ============================================================================
// 【类型参数基础】
// ============================================================================
// 泛型函数在函数名后用方括号声明类型参数
//
// 【语法】
// func 函数名[T 约束](参数 T) T { ... }
// func 函数名[T, U 约束](参数1 T, 参数2 U) { ... }
//
// 【类型参数命名约定】
// - T: Type（单个类型参数）
// - K, V: Key, Value（map 相关）
// - E: Element（切片元素）
// - S: Slice
// - 大写单字母是惯例，也可以用有意义的名字
// ============================================================================

// Min: 泛型函数，返回两个值中的最小值
// [T cmp.Ordered] 表示：
// - T 是类型参数
// - cmp.Ordered 是约束，T 必须是可排序的类型
//
// 【cmp.Ordered 包含的类型】
// - 整数类型：int, int8, int16, int32, int64
// - 无符号整数：uint, uint8, uint16, uint32, uint64
// - 浮点数：float32, float64
// - 字符串：string
func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Max: 泛型函数，返回两个值中的最大值
func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Swap: 泛型函数，交换两个值
// [T any] 表示 T 可以是任意类型
// any 是 interface{} 的别名（Go 1.18+）
//
// 【返回多值】
// 泛型函数可以返回多个同类型或不同类型的值
func Swap[T any](a, b T) (T, T) {
	return b, a
}

// ============================================================================
// 【类型约束】
// ============================================================================
// 类型约束定义了类型参数允许的类型集合
//
// 【约束的定义方式】
// 1. 使用内置约束：any, comparable
// 2. 使用标准库约束：cmp.Ordered
// 3. 自定义接口约束
//
// 【接口作为约束的新语法】
// type 约束名 interface {
//     Type1 | Type2 | Type3  // 类型联合
//     ~BaseType              // 类型近似（包含基础类型及其派生类型）
//     Method()               // 方法要求
// }
//
// 【| 操作符】
// 表示类型联合，T 可以是列出的任意一种类型
//
// 【~ 操作符（波浪线）】
// ~int 表示：底层类型是 int 的所有类型
// 包括 int 本身，以及 type MyInt int 这样的派生类型
// ============================================================================

// Number: 自定义类型约束
// 定义了可以进行数值运算的类型集合
//
// 【~ 的作用】
// ~int 不仅包括 int，还包括：
// - type MyInt int
// - type Age int
// 等所有以 int 为底层类型的自定义类型
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Sum: 使用自定义约束的泛型函数
// 计算切片中所有元素的和
// 只有满足 Number 约束的类型才能使用此函数
func Sum[T Number](nums []T) T {
	var total T // 零值初始化
	for _, n := range nums {
		total += n
	}
	return total
}

// Stringer: 带方法约束的接口
// 这是传统的接口用法，但现在也可以作为类型约束
// T 必须实现 String() string 方法
type Stringer interface {
	String() string
}

// PrintAll: 使用方法约束的泛型函数
// T 必须实现 Stringer 接口
func PrintAll[T Stringer](items []T) {
	for _, item := range items {
		fmt.Println(item.String())
	}
}

// ============================================================================
// 【泛型类型（结构体）】
// ============================================================================
// 结构体也可以是泛型的
//
// 【语法】
// type 类型名[T 约束] struct {
//     字段 T
// }
//
// 【泛型类型的方法】
// func (s *类型名[T]) 方法名() { }
// 注意：方法不能有自己的类型参数
// ============================================================================

// Stack: 泛型栈数据结构
// [T any] 表示栈可以存储任意类型的元素
//
// 【为什么用 any 而不是特定约束？】
// - 栈只需要存储和取出元素
// - 不需要比较或排序元素
// - any 提供最大的灵活性
type Stack[T any] struct {
	items []T // 使用切片存储元素
}

// NewStack: 泛型构造函数
// 【注意】
// - 构造函数是普通函数，不是方法
// - 需要声明自己的类型参数 [T any]
// - 返回类型必须指定类型参数 *Stack[T]
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{items: make([]T, 0)}
}

// Push: 入栈操作
// 【方法的类型参数】
// - 方法使用结构体定义的类型参数 T
// - 方法不能有额外的类型参数
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop: 出栈操作
// 【返回零值的处理】
// - 当栈为空时，需要返回类型 T 的零值
// - 使用 var zero T 获取任意类型的零值
// - 返回 bool 表示操作是否成功
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T // 任意类型的零值
		return zero, false
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, true
}

// Peek: 查看栈顶元素（不移除）
func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// Len: 返回栈中元素数量
func (s *Stack[T]) Len() int {
	return len(s.items)
}

// ============================================================================
// 【多类型参数】
// ============================================================================
// 泛型类型可以有多个类型参数
//
// 【语法】
// type 类型名[K 约束1, V 约束2] struct { ... }
//
// 【使用时必须指定所有类型参数】
// Pair[string, int]{...}
// ============================================================================

// Pair: 泛型键值对
// 两个类型参数 K 和 V，都使用 any 约束
type Pair[K, V any] struct {
	Key   K
	Value V
}

// NewPair: 创建键值对的泛型构造函数
// 【类型推断】
// 通常不需要显式指定类型参数
// NewPair("name", 42) 会自动推断为 Pair[string, int]
func NewPair[K, V any](key K, value V) Pair[K, V] {
	return Pair[K, V]{Key: key, Value: value}
}

// ============================================================================
// 【Result 类型】
// ============================================================================
// 类似 Rust 的 Result<T, E> 类型
// 用于表示可能成功或失败的操作结果
//
// 【设计模式】
// - Ok(value): 成功时返回
// - Err(err): 失败时返回
// - 提供 Unwrap, UnwrapOr 等方法处理结果
//
// 【与 Go 错误处理的区别】
// - 传统 Go：return value, err
// - Result：return Result[T]{...}
// 这只是一种封装，不是推荐的 Go 风格
// ============================================================================

// Result: 泛型结果类型
type Result[T any] struct {
	value T     // 成功时的值
	err   error // 错误（nil 表示成功）
}

// Ok: 创建成功的 Result
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value}
}

// Err: 创建失败的 Result
// 【注意类型参数】
// 即使失败了，也需要指定 T 的类型
// 因为 Result[T] 的类型是固定的
func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

// IsOk: 检查是否成功
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// Unwrap: 获取值，失败时 panic
// 【警告】只有在确定成功时才调用
func (r Result[T]) Unwrap() T {
	if r.err != nil {
		panic(r.err)
	}
	return r.value
}

// UnwrapOr: 获取值，失败时返回默认值
// 【推荐】比 Unwrap 更安全
func (r Result[T]) UnwrapOr(defaultVal T) T {
	if r.err != nil {
		return defaultVal
	}
	return r.value
}

// ============================================================================
// 【泛型切片操作】
// ============================================================================
// 函数式编程风格的切片操作
// 这些函数在 Go 1.21+ 的 slices 包中有官方实现
//
// 【函数式编程三大操作】
// - Map: 转换每个元素
// - Filter: 筛选元素
// - Reduce: 归约为单个值
// ============================================================================

// Map: 转换切片中的每个元素
// [T, U any] 表示输入类型 T 和输出类型 U 可以不同
//
// 【示例】
// Map([]int{1,2,3}, func(n int) string { return fmt.Sprint(n) })
// 结果：[]string{"1", "2", "3"}
func Map[T, U any](slice []T, f func(T) U) []U {
	result := make([]U, len(slice)) // 预分配容量
	for i, v := range slice {
		result[i] = f(v)
	}
	return result
}

// Filter: 过滤切片
// 返回满足条件的元素
//
// 【predicate 参数】
// 谓词函数，返回 true 表示保留该元素
func Filter[T any](slice []T, predicate func(T) bool) []T {
	var result []T // 不预分配，因为不知道最终大小
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Reduce: 归约切片
// 将切片元素依次累积到初始值上
//
// 【参数】
// - slice: 输入切片
// - initial: 初始值（累积器的初始状态）
// - f: 归约函数，接收累积器和当前元素，返回新的累积器
//
// 【示例】
// Reduce([]int{1,2,3}, 0, func(acc, n int) int { return acc + n })
// 计算过程：0 + 1 = 1, 1 + 2 = 3, 3 + 3 = 6
// 结果：6
func Reduce[T, U any](slice []T, initial U, f func(U, T) U) U {
	acc := initial
	for _, v := range slice {
		acc = f(acc, v)
	}
	return acc
}

// Find: 查找第一个满足条件的元素
// 返回 (元素, true) 或 (零值, false)
func Find[T any](slice []T, predicate func(T) bool) (T, bool) {
	for _, v := range slice {
		if predicate(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// Contains: 检查切片是否包含指定元素
// 【comparable 约束】
// T 必须是可比较的类型（支持 == 操作）
// 包括：整数、浮点数、字符串、指针、channel、接口、数组、结构体（如果所有字段可比较）
// 不包括：切片、map、函数
func Contains[T comparable](slice []T, target T) bool {
	for _, v := range slice {
		if v == target { // 需要 comparable 约束才能使用 ==
			return true
		}
	}
	return false
}

// ============================================================================
// 【泛型 Map 操作】
// ============================================================================
// 对 map 的泛型操作
// ============================================================================

// Keys: 获取 map 的所有键
// 【约束要求】
// - K: comparable（map 的键必须可比较）
// - V: any（值可以是任意类型）
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m)) // 预分配容量
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values: 获取 map 的所有值
func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// ============================================================================
// 【用于演示的类型】
// ============================================================================

// Person: 演示用的结构体
// 实现了 Stringer 接口，可以用于 PrintAll 函数
type Person struct {
	Name string
	Age  int
}

// String: 实现 Stringer 接口
func (p Person) String() string {
	return fmt.Sprintf("%s (%d)", p.Name, p.Age)
}

func main() {
	fmt.Println("=== Go 泛型 (Go 1.18+) ===\n")

	// ========================================================================
	// 【基本泛型函数】
	// ========================================================================
	// 泛型函数可以处理多种类型，编译器自动推断类型参数
	//
	// 【类型推断】
	// - 大多数情况下不需要显式指定类型
	// - 编译器从参数推断类型参数
	// - 当无法推断时，需要显式指定
	// ========================================================================
	fmt.Println("--- 基本泛型函数 ---")

	// 类型推断：编译器自动推断 T 的类型
	// Min(3, 5) -> Min[int](3, 5)
	fmt.Printf("Min(3, 5) = %d\n", Min(3, 5))

	// T 推断为 float64
	fmt.Printf("Min(3.14, 2.71) = %.2f\n", Min(3.14, 2.71))

	// T 推断为 string（字符串按字典序比较）
	fmt.Printf("Min(\"apple\", \"banana\") = %s\n", Min("apple", "banana"))

	fmt.Printf("Max(10, 20) = %d\n", Max(10, 20))

	// 显式指定类型参数
	// 【何时需要显式指定】
	// - 无法从参数推断类型时
	// - 想要明确类型时
	// - 需要类型转换时
	fmt.Printf("Min[int](7, 3) = %d\n", Min[int](7, 3))

	// Swap 函数：交换两个值
	a, b := Swap(1, 2)
	fmt.Printf("Swap(1, 2) = %d, %d\n", a, b)

	s1, s2 := Swap("hello", "world")
	fmt.Printf("Swap(\"hello\", \"world\") = %s, %s\n", s1, s2)

	// ========================================================================
	// 【自定义约束】
	// ========================================================================
	// Number 约束允许所有数值类型（包括自定义数值类型）
	// ========================================================================
	fmt.Println("\n--- 自定义类型约束 ---")

	ints := []int{1, 2, 3, 4, 5}
	floats := []float64{1.1, 2.2, 3.3}

	// Sum 可以处理任何满足 Number 约束的类型
	fmt.Printf("Sum(%v) = %d\n", ints, Sum(ints))
	fmt.Printf("Sum(%v) = %.1f\n", floats, Sum(floats))

	// ========================================================================
	// 【泛型栈】
	// ========================================================================
	// 泛型数据结构可以存储任意类型的元素
	// 同一份代码适用于所有类型
	// ========================================================================
	fmt.Println("\n--- 泛型栈 ---")

	// 整数栈
	// 【实例化】NewStack[int]() 创建存储 int 的栈
	intStack := NewStack[int]()
	intStack.Push(1)
	intStack.Push(2)
	intStack.Push(3)

	fmt.Printf("栈长度: %d\n", intStack.Len())
	if val, ok := intStack.Pop(); ok {
		fmt.Printf("Pop: %d\n", val)
	}
	if val, ok := intStack.Peek(); ok {
		fmt.Printf("Peek: %d\n", val)
	}

	// 字符串栈
	// 同样的 Stack 类型，不同的类型参数
	strStack := NewStack[string]()
	strStack.Push("Go")
	strStack.Push("is")
	strStack.Push("awesome")

	// 依次出栈
	for strStack.Len() > 0 {
		if val, ok := strStack.Pop(); ok {
			fmt.Printf("Pop string: %s\n", val)
		}
	}

	// ========================================================================
	// 【泛型键值对】
	// ========================================================================
	// Pair 是一个有两个类型参数的泛型类型
	// ========================================================================
	fmt.Println("\n--- 泛型键值对 ---")

	// 类型推断：Pair[string, string]
	pair1 := NewPair("name", "Gopher")
	// 类型推断：Pair[int, float64]
	pair2 := NewPair(1, 3.14)

	fmt.Printf("Pair1: Key=%s, Value=%s\n", pair1.Key, pair1.Value)
	fmt.Printf("Pair2: Key=%d, Value=%.2f\n", pair2.Key, pair2.Value)

	// ========================================================================
	// 【Result 类型】
	// ========================================================================
	// 演示类似 Rust 的 Result 类型
	// 这不是 Go 的标准模式，仅作为泛型示例
	// ========================================================================
	fmt.Println("\n--- Result 类型 ---")

	// 成功的 Result
	r1 := Ok(42)
	fmt.Printf("Ok(42).Unwrap() = %d\n", r1.Unwrap())

	// 失败的 Result
	// 【注意】必须指定类型参数 [int]
	r2 := Err[int](fmt.Errorf("something went wrong"))
	fmt.Printf("Err.IsOk() = %v\n", r2.IsOk())
	fmt.Printf("Err.UnwrapOr(0) = %d\n", r2.UnwrapOr(0))

	// ========================================================================
	// 【泛型切片操作】
	// ========================================================================
	// 函数式编程风格的切片操作
	// ========================================================================
	fmt.Println("\n--- 泛型切片操作 ---")

	numbers := []int{1, 2, 3, 4, 5}

	// Map: 转换元素
	// 每个元素乘以 2
	doubled := Map(numbers, func(n int) int { return n * 2 })
	fmt.Printf("Map (double): %v -> %v\n", numbers, doubled)

	// Map: 类型转换
	// int -> string
	strings := Map(numbers, func(n int) string { return fmt.Sprintf("#%d", n) })
	fmt.Printf("Map (to string): %v -> %v\n", numbers, strings)

	// Filter: 筛选元素
	// 只保留偶数
	evens := Filter(numbers, func(n int) bool { return n%2 == 0 })
	fmt.Printf("Filter (evens): %v -> %v\n", numbers, evens)

	// Reduce: 求和
	sum := Reduce(numbers, 0, func(acc, n int) int { return acc + n })
	fmt.Printf("Reduce (sum): %v -> %d\n", numbers, sum)

	// Reduce: 求积
	product := Reduce(numbers, 1, func(acc, n int) int { return acc * n })
	fmt.Printf("Reduce (product): %v -> %d\n", numbers, product)

	// Find: 查找第一个大于 3 的元素
	if found, ok := Find(numbers, func(n int) bool { return n > 3 }); ok {
		fmt.Printf("Find (>3): %d\n", found)
	}

	// Contains: 检查元素是否存在
	fmt.Printf("Contains(3): %v\n", Contains(numbers, 3))
	fmt.Printf("Contains(10): %v\n", Contains(numbers, 10))

	// ========================================================================
	// 【泛型 Map 操作】
	// ========================================================================
	fmt.Println("\n--- 泛型 Map 操作 ---")

	ages := map[string]int{
		"Alice":   25,
		"Bob":     30,
		"Charlie": 35,
	}

	// 【注意】map 的遍历顺序是随机的
	// 所以 Keys 和 Values 的顺序可能每次不同
	fmt.Printf("Keys: %v\n", Keys(ages))
	fmt.Printf("Values: %v\n", Values(ages))

	// ========================================================================
	// 【方法约束】
	// ========================================================================
	// Stringer 约束要求类型必须实现 String() 方法
	// ========================================================================
	fmt.Println("\n--- 方法约束 ---")

	// Person 实现了 String() string 方法
	// 所以可以用于 PrintAll[Stringer]
	people := []Person{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
	}
	fmt.Println("PrintAll:")
	PrintAll(people)

	// ========================================================================
	// 【内置约束】
	// ========================================================================
	// Go 提供了几个重要的内置约束
	//
	// 【any】
	// - 等价于 interface{}
	// - 允许任意类型
	// - 不能进行任何操作（除了赋值）
	//
	// 【comparable】
	// - 允许使用 == 和 != 比较的类型
	// - 包括：基本类型、指针、channel、接口、可比较的数组和结构体
	// - 不包括：切片、map、函数
	//
	// 【cmp.Ordered (Go 1.21+)】
	// - 允许使用 <, >, <=, >= 比较的类型
	// - 包括：整数、浮点数、字符串
	// ========================================================================
	fmt.Println("\n--- 内置约束 ---")
	fmt.Println("any       - 任意类型 (interface{})")
	fmt.Println("comparable - 可比较类型 (支持 == 和 !=)")
	fmt.Println("cmp.Ordered - 可排序类型 (支持 < > <= >=)")

	// ========================================================================
	// 【~ 波浪线操作符】
	// ========================================================================
	// ~ 表示类型近似（type approximation）
	// ~T 匹配所有底层类型是 T 的类型
	//
	// 【示例】
	// type MyInt int
	// ~int 匹配 int 和 MyInt
	// int 只匹配 int（不匹配 MyInt）
	//
	// 【Number 约束使用 ~】
	// type Number interface { ~int | ~float64 | ... }
	// 所以 type MyInt int 的 []MyInt 也可以用 Sum 函数
	// ========================================================================
	fmt.Println("\n--- ~ 波浪线操作符 ---")

	// ~ 表示基础类型，允许自定义类型
	// MyInt 的底层类型是 int
	// 因为 Number 约束使用了 ~int，所以 MyInt 满足 Number 约束
	type MyInt int
	var myInts []MyInt = []MyInt{1, 2, 3}
	fmt.Printf("Sum(MyInt): %d\n", Sum(myInts))

	// 【如果 Number 没有使用 ~】
	// type Number interface { int | float64 | ... }
	// 那么 Sum(myInts) 会编译错误
	// 因为 MyInt 不等于 int（虽然底层类型相同）
}
