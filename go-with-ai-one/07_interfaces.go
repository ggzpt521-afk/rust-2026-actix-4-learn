// ============================================================================
// 07_interfaces.go - 接口与类型断言
// ============================================================================
// 运行: go run 07_interfaces.go
//
// 【本文件学习目标】
// 1. 理解 Go 接口的隐式实现机制
// 2. 掌握接口的多态特性
// 3. 学会使用类型断言和类型 switch
// 4. 理解空接口 interface{}/any 的用法
// 5. 掌握接口组合的设计模式
// 6. 了解接口值的内部结构
//
// 【Go 接口的核心特点】
// - 隐式实现：不需要 implements 关键字
// - 结构化类型（鸭子类型）：只要方法签名匹配就实现了接口
// - 接口是一种类型：可以作为参数、返回值、变量类型
// - 小接口设计哲学：Go 推崇单方法接口（如 io.Reader, io.Writer）
//
// 【接口 vs 其他语言】
// | 特性       | Go           | Java/C#        |
// |------------|--------------|----------------|
// | 实现方式   | 隐式         | 显式 implements |
// | 接口大小   | 推崇小接口   | 通常较大        |
// | 多重继承   | 组合代替继承 | 支持多接口      |
// | 默认方法   | 不支持       | Java 8+ 支持    |
// ============================================================================

package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// ============================================================================
// 【接口定义】
// ============================================================================
// 接口定义了一组方法签名，任何实现了这些方法的类型都实现了该接口
//
// 【语法】
// type 接口名 interface {
//     方法名(参数列表) 返回值列表
//     ...
// }
//
// 【命名约定】
// - 单方法接口通常以方法名+er 结尾：Reader, Writer, Stringer
// - 接口名首字母大写表示可导出
// ============================================================================

// Shape: 基本接口示例
// 任何实现了 Area() 和 Perimeter() 方法的类型都是 Shape
type Shape interface {
	Area() float64      // 计算面积
	Perimeter() float64 // 计算周长
}

// ============================================================================
// 【单方法接口】
// ============================================================================
// Go 标准库大量使用单方法接口，这是 Go 的设计哲学
//
// 【优点】
// - 易于实现：只需实现一个方法
// - 高度可组合：多个小接口可以组合成大接口
// - 解耦性强：最小化依赖
//
// 【标准库中的单方法接口】
// - io.Reader: Read(p []byte) (n int, err error)
// - io.Writer: Write(p []byte) (n int, err error)
// - io.Closer: Close() error
// - fmt.Stringer: String() string
// - error: Error() string
// ============================================================================

// Stringer: 字符串表示接口
// fmt.Println 等会自动调用实现了此接口的类型的 String() 方法
type Stringer interface {
	String() string
}

// Writer: 写入接口
// 与 io.Writer 相同的签名
type Writer interface {
	Write(data []byte) (int, error) // 返回写入字节数和可能的错误
}

// Reader: 读取接口
// 与 io.Reader 相同的签名
type Reader interface {
	Read(buf []byte) (int, error) // 返回读取字节数和可能的错误
}

// ============================================================================
// 【接口组合（嵌入）】
// ============================================================================
// 接口可以嵌入其他接口，形成更大的接口
//
// 【语法】
// type 大接口 interface {
//     小接口1
//     小接口2
//     其他方法()
// }
//
// 【标准库示例】
// io.ReadWriter = io.Reader + io.Writer
// io.ReadCloser = io.Reader + io.Closer
// io.ReadWriteCloser = io.Reader + io.Writer + io.Closer
//
// 【设计原则】
// - 组合优于继承
// - 从小接口组合出大接口
// - 需要什么就组合什么
// ============================================================================

// ReadWriter: 组合接口示例
// 同时具有读和写的能力
type ReadWriter interface {
	Reader // 嵌入 Reader 接口
	Writer // 嵌入 Writer 接口
}

// ============================================================================
// 【空接口】
// ============================================================================
// interface{} 是空接口，没有任何方法要求
// 任何类型都实现了空接口，所以它可以存储任何值
//
// 【Go 1.18+ any】
// any 是 interface{} 的类型别名：type any = interface{}
//
// 【用途】
// - 存储任意类型的值（如 []any）
// - 泛型出现前的"万能类型"
// - JSON 解析的中间结果
//
// 【注意】
// - 使用空接口会丢失类型信息
// - 需要类型断言才能取出具体值
// - 尽量避免滥用，优先使用具体类型或泛型
// ============================================================================

// Any: 空接口的类型别名（演示用）
// 实际代码中直接使用 any 或 interface{}
type Any interface{}

// ============================================================================
// 【实现接口的类型】
// ============================================================================
// Go 中实现接口是隐式的：
// - 不需要 implements 关键字
// - 只需要实现接口定义的所有方法
// - 编译器自动检查类型是否满足接口
//
// 【值接收者 vs 指针接收者】
// | 实现方式    | T 实现接口 | *T 实现接口 |
// |-------------|------------|-------------|
// | 值接收者    | ✓          | ✓           |
// | 指针接收者  | ✗          | ✓           |
//
// 如果方法使用指针接收者，只有指针类型才能赋值给接口
// ============================================================================

// Rectangle: 矩形，实现 Shape 和 Stringer 接口
type Rectangle struct {
	Width, Height float64 // 宽度和高度
}

// Area: 计算矩形面积（实现 Shape.Area）
// 使用值接收者：不修改结构体，Rectangle 和 *Rectangle 都能赋值给 Shape
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter: 计算矩形周长（实现 Shape.Perimeter）
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// String: 字符串表示（实现 fmt.Stringer 接口）
// 当使用 fmt.Println 打印时会自动调用此方法
func (r Rectangle) String() string {
	return fmt.Sprintf("Rectangle(%.2f x %.2f)", r.Width, r.Height)
}

// Circle: 圆形，实现 Shape 接口
type Circle struct {
	Radius float64 // 半径
}

// Area: 计算圆的面积
// 公式: π × r²
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Perimeter: 计算圆的周长（圆周）
// 公式: 2 × π × r
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// String: 圆的字符串表示
func (c Circle) String() string {
	return fmt.Sprintf("Circle(radius=%.2f)", c.Radius)
}

// Triangle: 三角形，实现 Shape 接口
type Triangle struct {
	A, B, C float64 // 三边长度
}

// Area: 计算三角形面积
// 使用海伦公式（Heron's formula）
// s = (a + b + c) / 2
// Area = √(s(s-a)(s-b)(s-c))
func (t Triangle) Area() float64 {
	s := (t.A + t.B + t.C) / 2 // 半周长
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}

// Perimeter: 计算三角形周长
func (t Triangle) Perimeter() float64 {
	return t.A + t.B + t.C
}

// ConsoleWriter: 简单的 Writer 实现
// 将数据输出到控制台
type ConsoleWriter struct{}

// Write: 实现 Writer 接口
// 将字节数据打印到标准输出
func (w ConsoleWriter) Write(data []byte) (int, error) {
	fmt.Print(string(data))   // 转换为字符串并打印
	return len(data), nil     // 返回写入的字节数，无错误
}

// ============================================================================
// 【sort.Interface 接口示例】
// ============================================================================
// sort.Interface 定义了三个方法：
//   Len() int           - 返回元素数量
//   Swap(i, j int)      - 交换两个元素
//   Less(i, j int) bool - 比较两个元素
//
// 实现这三个方法后，可以使用 sort.Sort() 进行排序
//
// 【Go 1.8+ 更简单的方式】
// sort.Slice(slice, func(i, j int) bool { return ... })
// 不需要实现完整接口，只需提供比较函数
// ============================================================================

// Person: 人员信息
type Person struct {
	Name string // 姓名
	Age  int    // 年龄
}

// ByAge: 按年龄排序的 Person 切片类型
// 定义新类型是为了给它添加方法
type ByAge []Person

// Len: 返回切片长度（实现 sort.Interface）
func (a ByAge) Len() int { return len(a) }

// Swap: 交换两个元素（实现 sort.Interface）
func (a ByAge) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less: 比较两个元素（实现 sort.Interface）
// 返回 true 表示 i 应该排在 j 前面
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }

func main() {
	fmt.Println("=== Go 接口与类型断言 ===\n")

	// ========================================================================
	// 【接口基础】
	// ========================================================================
	// 接口变量可以存储任何实现了该接口的值
	//
	// 【原理】
	// 接口变量在底层由两部分组成：
	// 1. 类型信息（type）：存储的具体类型
	// 2. 值信息（value）：存储的具体值
	//
	// 【赋值规则】
	// - 具体类型 -> 接口：总是可以（如果实现了接口）
	// - 接口 -> 具体类型：需要类型断言
	// ========================================================================
	fmt.Println("--- 接口基础 ---")

	// 声明接口变量
	var shape Shape // 零值为 nil

	// 将 Rectangle 赋值给 Shape 接口
	shape = Rectangle{Width: 10, Height: 5}
	// 通过接口调用方法（多态）
	fmt.Printf("Rectangle: Area=%.2f, Perimeter=%.2f\n", shape.Area(), shape.Perimeter())

	// 同一个接口变量可以存储不同的实现类型
	shape = Circle{Radius: 3}
	fmt.Printf("Circle: Area=%.2f, Perimeter=%.2f\n", shape.Area(), shape.Perimeter())

	// ========================================================================
	// 【接口值的内部结构】
	// ========================================================================
	// 接口值 = (type, value)
	//
	// 【内存布局】
	// ┌────────────────┐
	// │   Interface    │
	// ├────────────────┤
	// │ type: *_type   │ -> 指向类型信息
	// │ data: unsafe.Pointer │ -> 指向实际数据
	// └────────────────┘
	//
	// 【nil 接口 vs 包含 nil 的接口】
	// - nil 接口：type 和 value 都是 nil
	// - 包含 nil 的接口：type 不为 nil，value 为 nil
	// - 这两者不相等！
	// ========================================================================
	fmt.Println("\n--- 接口值的内部结构 ---")

	var s Shape // nil 接口
	fmt.Printf("nil 接口: value=%v, type=%T\n", s, s)

	s = Rectangle{Width: 5, Height: 3}
	fmt.Printf("赋值后: value=%v, type=%T\n", s, s)
	// 输出显示接口存储的具体类型是 main.Rectangle

	// 接口值由 (type, value) 组成
	// 即使 value 是 nil，只要 type 不是 nil，接口就不是 nil

	// ========================================================================
	// 【多态（Polymorphism）】
	// ========================================================================
	// 多态是面向对象编程的核心特性之一
	// 同一个接口，不同的实现，产生不同的行为
	//
	// 【优点】
	// - 代码复用：对接口编程，不依赖具体实现
	// - 扩展性：新增实现无需修改现有代码
	// - 可测试性：可以用 mock 实现替换真实实现
	// ========================================================================
	fmt.Println("\n--- 多态 ---")

	// 创建不同形状的切片
	shapes := []Shape{
		Rectangle{Width: 10, Height: 5}, // 矩形
		Circle{Radius: 3},               // 圆形
		Triangle{A: 3, B: 4, C: 5},      // 三角形（3-4-5 直角三角形）
	}

	// 统一处理所有形状
	totalArea := 0.0
	for i, shape := range shapes {
		// 多态：同一个 Area() 调用，不同类型有不同的实现
		fmt.Printf("Shape %d: Area=%.2f\n", i+1, shape.Area())
		totalArea += shape.Area()
	}
	fmt.Printf("Total Area: %.2f\n", totalArea)

	// ========================================================================
	// 【类型断言（Type Assertion）】
	// ========================================================================
	// 从接口值中提取具体类型的值
	//
	// 【语法】
	// value := interfaceValue.(ConcreteType)      // 不安全，失败会 panic
	// value, ok := interfaceValue.(ConcreteType)  // 安全，失败返回 false
	//
	// 【原理】
	// 类型断言检查接口的动态类型是否匹配
	// 如果匹配，返回底层值；否则 panic 或返回 false
	//
	// 【最佳实践】
	// 总是使用带 ok 的形式，避免 panic
	// ========================================================================
	fmt.Println("\n--- 类型断言 ---")

	var i interface{} = "hello" // 空接口可以存储任何值

	// 不安全的类型断言（失败会 panic）
	str := i.(string) // 我们确定它是 string
	fmt.Printf("断言成功: %s\n", str)

	// 安全的类型断言（带 ok 检查）
	if str, ok := i.(string); ok {
		fmt.Printf("是 string: %s\n", str)
	}

	// 断言为错误的类型
	if num, ok := i.(int); ok {
		fmt.Printf("是 int: %d\n", num)
	} else {
		fmt.Println("不是 int 类型") // 执行这里
	}

	// 危险！不带 ok 的断言失败会 panic
	// num := i.(int) // panic: interface conversion: interface {} is string, not int

	// ========================================================================
	// 【类型 switch】
	// ========================================================================
	// 用于判断接口的动态类型并执行相应代码
	//
	// 【语法】
	// switch v := x.(type) {
	// case Type1:
	//     // v 是 Type1 类型
	// case Type2:
	//     // v 是 Type2 类型
	// default:
	//     // v 是 interface{} 类型
	// }
	//
	// 【注意】
	// - x.(type) 只能在 switch 语句中使用
	// - 每个 case 中 v 的类型不同
	// - 可以合并多个类型：case int, int64:
	// ========================================================================
	fmt.Println("\n--- 类型 switch ---")

	values := []interface{}{42, "hello", 3.14, true, []int{1, 2, 3}, nil}

	for _, v := range values {
		describeType(v) // 使用类型 switch 判断类型
	}

	// ========================================================================
	// 【接口断言】
	// ========================================================================
	// 不仅可以断言为具体类型，还可以断言为其他接口
	//
	// 【场景】
	// - 检查类型是否实现了某个接口
	// - 获取类型的其他能力（如 Stringer）
	// ========================================================================
	fmt.Println("\n--- 接口断言 ---")

	var sh Shape = Rectangle{Width: 5, Height: 3}

	// 断言为具体类型
	if rect, ok := sh.(Rectangle); ok {
		// 可以访问 Rectangle 特有的字段
		fmt.Printf("是 Rectangle: Width=%.2f, Height=%.2f\n", rect.Width, rect.Height)
	}

	// 断言为另一个接口（fmt.Stringer）
	if stringer, ok := sh.(fmt.Stringer); ok {
		// Rectangle 也实现了 Stringer 接口
		fmt.Printf("实现了 Stringer: %s\n", stringer.String())
	}

	// ========================================================================
	// 【空接口 interface{} / any】
	// ========================================================================
	// 空接口没有任何方法要求，因此任何类型都实现了空接口
	//
	// 【Go 1.18+ any】
	// any 是 interface{} 的类型别名，更简洁
	// type any = interface{} // 在 builtin 包中定义
	//
	// 【典型用途】
	// - fmt.Println(a ...interface{}) - 接受任意参数
	// - json.Unmarshal 解析未知结构的 JSON
	// - 泛型容器（Go 1.18 前）
	//
	// 【注意】
	// - 空接口会丢失类型信息
	// - 需要类型断言才能使用具体类型的方法
	// - Go 1.18+ 优先使用泛型而非空接口
	// ========================================================================
	fmt.Println("\n--- 空接口 (interface{} / any) ---")

	// any 是 interface{} 的别名（Go 1.18+）
	var anything any

	// 存储各种类型
	anything = 42
	fmt.Printf("存储 int: %v (%T)\n", anything, anything)

	anything = "Go语言"
	fmt.Printf("存储 string: %v (%T)\n", anything, anything)

	anything = []float64{1.1, 2.2, 3.3}
	fmt.Printf("存储 slice: %v (%T)\n", anything, anything)

	// 空接口切片可以存储混合类型
	mixed := []any{1, "two", 3.0, true}
	fmt.Printf("混合切片: %v\n", mixed)

	// ========================================================================
	// 【常用接口示例】
	// ========================================================================
	// Go 标准库定义了许多重要接口
	//
	// 【常见接口】
	// | 接口           | 包     | 方法                           |
	// |----------------|--------|--------------------------------|
	// | Stringer       | fmt    | String() string                |
	// | error          | builtin| Error() string                 |
	// | Reader         | io     | Read([]byte) (int, error)      |
	// | Writer         | io     | Write([]byte) (int, error)     |
	// | Closer         | io     | Close() error                  |
	// | sort.Interface | sort   | Len, Swap, Less                |
	// ========================================================================
	fmt.Println("\n--- 常用接口示例 ---")

	// fmt.Stringer 接口
	rect := Rectangle{Width: 10, Height: 5}
	fmt.Printf("Stringer: %s\n", rect) // 自动调用 String() 方法

	// sort.Interface 接口
	people := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
	}

	fmt.Println("排序前:", people)
	sort.Sort(ByAge(people)) // 使用实现了 sort.Interface 的 ByAge 类型
	fmt.Println("按年龄排序后:", people)

	// io.Writer 接口
	var writer Writer = ConsoleWriter{}
	writer.Write([]byte("Hello from ConsoleWriter\n"))

	// ========================================================================
	// 【接口组合使用】
	// ========================================================================
	// 标准库中有很多组合接口的例子
	// 如 io.ReadWriter, io.ReadCloser 等
	// ========================================================================
	fmt.Println("\n--- 接口组合 ---")

	// 使用 strings.Reader（实现了 io.Reader）
	reader := strings.NewReader("Hello, World!")
	buf := make([]byte, 5)
	n, _ := reader.Read(buf) // 读取最多 5 字节
	fmt.Printf("读取 %d 字节: %s\n", n, buf)

	// ========================================================================
	// 【接口最佳实践】
	// ========================================================================
	// 1. 接受接口，返回具体类型
	//    - 函数参数用接口：灵活，可测试
	//    - 返回值用具体类型：清晰，避免不必要的抽象
	//
	// 2. 小接口优于大接口
	//    - 单方法接口最灵活
	//    - 需要时再组合
	//
	// 3. 不要过早定义接口
	//    - 先写具体实现
	//    - 当需要解耦或多态时再抽象出接口
	//
	// 4. 接口由使用方定义
	//    - Go 的接口是使用方驱动的
	//    - 不同于 Java 的提供方驱动
	// ========================================================================
	fmt.Println("\n--- 接口最佳实践 ---")

	fmt.Println("1. 函数应该接受接口类型，返回具体类型")
	fmt.Println("2. 定义小而专注的接口（通常1-3个方法）")
	fmt.Println("3. 先写具体实现，需要时再抽象出接口")

	// ========================================================================
	// 【接口的零值与 nil 陷阱】
	// ========================================================================
	// 这是 Go 接口最容易出错的地方！
	//
	// 【两种情况】
	// 1. nil 接口：type = nil, value = nil
	// 2. 包含 nil 值的接口：type ≠ nil, value = nil
	//
	// 【陷阱】
	// var p *MyType = nil
	// var i MyInterface = p  // i 不是 nil！因为 type 信息存在
	//
	// 【检查方法】
	// - 接口 == nil：只检查整个接口是否为 nil
	// - 反射：reflect.ValueOf(i).IsNil() 检查值是否为 nil
	// ========================================================================
	fmt.Println("\n--- 接口的零值 ---")

	// 真正的 nil 接口
	var nilInterface interface{}
	fmt.Printf("nil 接口 == nil: %v\n", nilInterface == nil) // true

	// 包含 nil 指针的接口
	var nilPointer *Rectangle = nil
	var shapeInterface Shape = nilPointer
	fmt.Printf("包含 nil 指针的接口 == nil: %v\n", shapeInterface == nil) // false!
	// 注意：虽然底层值是 nil，但接口本身不是 nil
	// 因为接口的 type 部分存储了 *Rectangle 类型信息

	// ========================================================================
	// 【编译时检查接口实现】
	// ========================================================================
	// Go 是隐式实现接口，编译器不会强制检查
	// 可以使用以下技巧在编译时验证类型实现了接口
	//
	// 【语法】
	// var _ InterfaceName = TypeName{}      // 值类型
	// var _ InterfaceName = (*TypeName)(nil) // 指针类型
	//
	// 【原理】
	// - _ 是空白标识符，忽略值
	// - 如果类型没有实现接口，编译会报错
	// - 不产生运行时开销
	// ========================================================================
	fmt.Println("\n--- 编译时检查接口实现 ---")

	// 编译时检查 Rectangle 是否实现了 Shape
	var _ Shape = Rectangle{}           // 值类型实现检查
	var _ Shape = (*Rectangle)(nil)     // 指针类型实现检查

	fmt.Println("Rectangle 实现了 Shape 接口 ✓")
}

// ============================================================================
// 【类型 switch 示例函数】
// ============================================================================
// 使用 switch v.(type) 判断接口的动态类型
//
// 【特点】
// - 每个 case 中 v 的类型是对应的具体类型
// - default 中 v 保持 interface{} 类型
// - 可以用 case nil: 处理 nil 值
// ============================================================================
func describeType(v interface{}) {
	switch t := v.(type) { // t 在每个 case 中有不同的具体类型
	case nil:
		fmt.Println("  nil 值")
	case int:
		fmt.Printf("  int: %d\n", t) // t 是 int 类型
	case string:
		fmt.Printf("  string: %q\n", t) // t 是 string 类型
	case float64:
		fmt.Printf("  float64: %.2f\n", t) // t 是 float64 类型
	case bool:
		fmt.Printf("  bool: %v\n", t) // t 是 bool 类型
	case []int:
		fmt.Printf("  []int: %v\n", t) // t 是 []int 类型
	default:
		fmt.Printf("  未知类型: %T\n", t) // %T 打印类型名
	}
}
