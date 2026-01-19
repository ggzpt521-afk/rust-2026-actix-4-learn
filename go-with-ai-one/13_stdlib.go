// ============================================================================
// 13_stdlib.go - 常用标准库
// ============================================================================
// 运行: go run 13_stdlib.go
//
// 【本文件学习目标】
// 1. 掌握 fmt 包的格式化输入输出
// 2. 学会使用 strings 和 strconv 包处理字符串
// 3. 理解 time 包的时间操作
// 4. 掌握 os、io、filepath 包的文件操作
// 5. 学会使用 encoding/json 进行 JSON 处理
// 6. 了解 regexp、sort、log、flag、net/http、math/rand 等常用包
//
// 【Go 标准库的特点】
// - 丰富而实用：覆盖大多数常见需求
// - 设计一致：API 风格统一
// - 文档完善：godoc 格式的文档
// - 高质量：经过充分测试
//
// 【常用标准库分类】
// | 类别       | 包名                           |
// |------------|--------------------------------|
// | 格式化     | fmt                            |
// | 字符串     | strings, strconv, unicode      |
// | 时间       | time                           |
// | 文件/IO    | os, io, bufio, filepath        |
// | 编码       | encoding/json, encoding/xml    |
// | 网络       | net, net/http, net/url         |
// | 并发       | sync, sync/atomic, context     |
// | 数学       | math, math/rand                |
// | 排序       | sort, slices (Go 1.21+)        |
// | 正则       | regexp                         |
// | 日志       | log, log/slog (Go 1.21+)       |
// | 命令行     | flag, os                       |
// | 测试       | testing                        |
// ============================================================================

package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Println("=== Go 常用标准库 ===\n")

	// 依次演示各个标准库
	fmtDemo()
	stringsDemo()
	strconvDemo()
	timeDemo()
	osDemo()
	ioDemo()
	filepathDemo()
	jsonDemo()
	regexpDemo()
	sortDemo()
	contextDemo()
	logDemo()
	flagDemo()
	httpDemo()
	randDemo()
}

// ============================================================================
// 【fmt 包】
// ============================================================================
// fmt 包实现了格式化 I/O
//
// 【Print 系列】
// Print, Println, Printf - 输出到标准输出
// Sprint, Sprintln, Sprintf - 返回字符串
// Fprint, Fprintln, Fprintf - 输出到 io.Writer
//
// 【Scan 系列】
// Scan, Scanln, Scanf - 从标准输入读取
// Sscan, Sscanln, Sscanf - 从字符串读取
// Fscan, Fscanln, Fscanf - 从 io.Reader 读取
// ============================================================================
func fmtDemo() {
	fmt.Println("--- fmt 包 ---")

	// 格式化输出
	name := "Gopher"
	age := 10
	score := 95.5

	// Print 系列
	// 【区别】
	// Print: 不换行，参数间无空格
	// Println: 换行，参数间有空格
	// Printf: 格式化输出
	fmt.Print("Print: 不换行")
	fmt.Println(" Println: 换行")
	fmt.Printf("Printf: name=%s, age=%d, score=%.1f\n", name, age, score)

	// Sprint 系列（返回字符串）
	// 【用途】
	// - 构建字符串
	// - 不直接输出，而是保存结果
	str := fmt.Sprintf("name=%s, age=%d", name, age)
	fmt.Printf("Sprintf: %s\n", str)

	// 常用格式化动词
	// 【格式化动词表】
	// | 动词  | 说明                           | 示例                    |
	// |-------|--------------------------------|-------------------------|
	// | %v    | 默认格式                       | {name: "Go"}            |
	// | %+v   | 结构体带字段名                 | {Name:Go}               |
	// | %#v   | Go 语法表示                    | main.Person{Name:"Go"}  |
	// | %T    | 类型                           | main.Person             |
	// | %d    | 十进制整数                     | 42                      |
	// | %b    | 二进制                         | 101010                  |
	// | %o    | 八进制                         | 52                      |
	// | %x    | 十六进制（小写）               | 2a                      |
	// | %X    | 十六进制（大写）               | 2A                      |
	// | %f    | 浮点数                         | 3.141593                |
	// | %.2f  | 浮点数（2 位小数）             | 3.14                    |
	// | %e    | 科学计数法                     | 3.141593e+00            |
	// | %s    | 字符串                         | hello                   |
	// | %q    | 带引号字符串                   | "hello"                 |
	// | %p    | 指针                           | 0xc0000...              |
	// | %%    | 百分号                         | %                       |
	fmt.Println("\n常用格式化动词:")
	fmt.Printf("  %%v  通用: %v\n", map[string]int{"a": 1})
	fmt.Printf("  %%+v 带字段名: %+v\n", struct{ Name string }{"Go"})
	fmt.Printf("  %%#v Go语法: %#v\n", []int{1, 2})
	fmt.Printf("  %%T  类型: %T\n", 3.14)
	fmt.Printf("  %%d  十进制: %d\n", 42)
	fmt.Printf("  %%b  二进制: %b\n", 42)
	fmt.Printf("  %%o  八进制: %o\n", 42)
	fmt.Printf("  %%x  十六进制: %x\n", 42)
	fmt.Printf("  %%f  浮点数: %f\n", 3.14159)
	fmt.Printf("  %%.2f 精度: %.2f\n", 3.14159)
	fmt.Printf("  %%e  科学计数: %e\n", 123456.789)
	fmt.Printf("  %%s  字符串: %s\n", "hello")
	fmt.Printf("  %%q  带引号: %q\n", "hello")
	fmt.Printf("  %%p  指针: %p\n", &name)
}

// ============================================================================
// 【strings 包】
// ============================================================================
// strings 包提供字符串操作函数
//
// 【常用函数分类】
// 查找：Contains, ContainsAny, HasPrefix, HasSuffix, Index, LastIndex
// 转换：ToUpper, ToLower, Title
// 修改：Replace, ReplaceAll, Trim, TrimSpace, TrimPrefix, TrimSuffix
// 分割：Split, SplitN, Fields
// 连接：Join, Repeat
// 比较：Compare, EqualFold（忽略大小写）
//
// 【strings.Builder】
// 高效的字符串构建器，避免多次字符串拼接的性能问题
// ============================================================================
func stringsDemo() {
	fmt.Println("\n--- strings 包 ---")

	s := "Hello, World!"

	fmt.Printf("原字符串: %s\n", s)

	// 查找
	fmt.Printf("Contains: %v\n", strings.Contains(s, "World"))   // 包含
	fmt.Printf("HasPrefix: %v\n", strings.HasPrefix(s, "Hello")) // 前缀
	fmt.Printf("HasSuffix: %v\n", strings.HasSuffix(s, "!"))     // 后缀
	fmt.Printf("Index: %d\n", strings.Index(s, "World"))         // 位置

	// 转换
	fmt.Printf("ToUpper: %s\n", strings.ToUpper(s)) // 大写
	fmt.Printf("ToLower: %s\n", strings.ToLower(s)) // 小写

	// 修改
	fmt.Printf("Replace: %s\n", strings.Replace(s, "World", "Go", 1)) // 替换

	// 分割和连接
	fmt.Printf("Split: %v\n", strings.Split(s, ", "))                  // 分割
	fmt.Printf("Join: %s\n", strings.Join([]string{"a", "b", "c"}, "-")) // 连接

	// 清理
	fmt.Printf("TrimSpace: [%s]\n", strings.TrimSpace("  hello  ")) // 去空白

	// 其他
	fmt.Printf("Repeat: %s\n", strings.Repeat("Go", 3)) // 重复
	fmt.Printf("Count: %d\n", strings.Count(s, "l"))    // 计数

	// Builder（高效字符串构建）
	// 【为什么使用 Builder】
	// - 字符串是不可变的，+ 拼接会产生新字符串
	// - Builder 内部使用 []byte，减少内存分配
	// - 适合大量字符串拼接场景
	var builder strings.Builder
	builder.WriteString("Hello")
	builder.WriteString(", ")
	builder.WriteString("Builder!")
	fmt.Printf("Builder: %s\n", builder.String())
}

// ============================================================================
// 【strconv 包】
// ============================================================================
// strconv 包提供字符串与基本类型之间的转换
//
// 【字符串 -> 数字】
// Atoi(s) - 字符串转 int
// ParseInt(s, base, bitSize) - 字符串转整数（指定进制和位数）
// ParseFloat(s, bitSize) - 字符串转浮点数
// ParseBool(s) - 字符串转布尔
//
// 【数字 -> 字符串】
// Itoa(i) - int 转字符串
// FormatInt(i, base) - 整数转字符串（指定进制）
// FormatFloat(f, fmt, prec, bitSize) - 浮点数转字符串
// FormatBool(b) - 布尔转字符串
//
// 【Quote 系列】
// Quote(s) - 添加 Go 字符串引号
// QuoteRune(r) - 添加 Go 字符引号
// Unquote(s) - 移除引号
// ============================================================================
func strconvDemo() {
	fmt.Println("\n--- strconv 包 ---")

	// 字符串转数字
	// 【Atoi】= ParseInt(s, 10, 0) 的简写
	// 返回 (int, error)
	i, _ := strconv.Atoi("42")
	fmt.Printf("Atoi: \"42\" -> %d\n", i)

	// ParseFloat: 第二个参数是位数（32 或 64）
	f, _ := strconv.ParseFloat("3.14", 64)
	fmt.Printf("ParseFloat: \"3.14\" -> %f\n", f)

	// ParseBool: 接受 "1", "t", "T", "TRUE", "true", "True" 为 true
	b, _ := strconv.ParseBool("true")
	fmt.Printf("ParseBool: \"true\" -> %v\n", b)

	// 数字转字符串
	// 【Itoa】= FormatInt(int64(i), 10) 的简写
	str := strconv.Itoa(42)
	fmt.Printf("Itoa: 42 -> %q\n", str)

	// FormatFloat: 参数为 (value, format, precision, bitSize)
	// format: 'f'=小数, 'e'=科学计数法, 'g'=自动选择
	str = strconv.FormatFloat(3.14159, 'f', 2, 64)
	fmt.Printf("FormatFloat: 3.14159 -> %q\n", str)

	str = strconv.FormatBool(true)
	fmt.Printf("FormatBool: true -> %q\n", str)

	// Quote: 为字符串添加 Go 语法的引号，转义特殊字符
	fmt.Printf("Quote: %s\n", strconv.Quote("Hello\tWorld"))
}

// ============================================================================
// 【time 包】
// ============================================================================
// time 包提供时间的测量和显示功能
//
// 【核心类型】
// Time: 时间点
// Duration: 时间段
// Location: 时区
//
// 【Go 的时间格式化】
// Go 使用特殊的参考时间: 2006-01-02 15:04:05 MST
// 这是 Go 诞生的时间（2006年1月2日15:04:05）
// 记忆：1月2日下午3点4分5秒 -> 01/02 03:04:05 -> 1 2 3 4 5
//
// 【常用格式】
// time.RFC3339: "2006-01-02T15:04:05Z07:00"
// time.RFC822:  "02 Jan 06 15:04 MST"
// time.Kitchen: "3:04PM"
// ============================================================================
func timeDemo() {
	fmt.Println("\n--- time 包 ---")

	// 当前时间
	now := time.Now()
	fmt.Printf("当前时间: %v\n", now)

	// 格式化（Go 使用特殊的参考时间: 2006-01-02 15:04:05）
	// 【重要】不是随意的日期，而是 Go 的诞生时间
	// 记忆：1月2日下午3点4分5秒2006年
	fmt.Printf("格式化: %s\n", now.Format("2006-01-02 15:04:05"))
	fmt.Printf("RFC3339: %s\n", now.Format(time.RFC3339))

	// 时间组件
	fmt.Printf("年月日: %d-%d-%d\n", now.Year(), now.Month(), now.Day())
	fmt.Printf("时分秒: %d:%d:%d\n", now.Hour(), now.Minute(), now.Second())
	fmt.Printf("星期: %s\n", now.Weekday())

	// 解析时间
	// 【Parse】第一个参数是格式，第二个参数是要解析的字符串
	t, _ := time.Parse("2006-01-02", "2024-03-15")
	fmt.Printf("解析: %v\n", t)

	// 时间计算
	// 【Add】添加 Duration
	tomorrow := now.Add(24 * time.Hour)
	fmt.Printf("明天: %s\n", tomorrow.Format("2006-01-02"))

	// 【AddDate】添加年/月/日
	yesterday := now.AddDate(0, 0, -1)
	fmt.Printf("昨天: %s\n", yesterday.Format("2006-01-02"))

	// 时间比较
	fmt.Printf("now.Before(tomorrow): %v\n", now.Before(tomorrow))
	fmt.Printf("now.After(yesterday): %v\n", now.After(yesterday))

	// Duration
	// 【Duration 常量】
	// time.Nanosecond, time.Microsecond, time.Millisecond
	// time.Second, time.Minute, time.Hour
	duration := 2*time.Hour + 30*time.Minute
	fmt.Printf("Duration: %v\n", duration)

	// 计时
	// 【time.Since】= time.Now().Sub(start)
	start := time.Now()
	time.Sleep(10 * time.Millisecond)
	elapsed := time.Since(start)
	fmt.Printf("耗时: %v\n", elapsed)

	// 定时器
	// 【time.NewTimer】创建单次定时器
	// timer.C 是一个 channel，到时后发送当前时间
	timer := time.NewTimer(50 * time.Millisecond)
	<-timer.C
	fmt.Println("Timer 触发")
}

// ============================================================================
// 【os 包】
// ============================================================================
// os 包提供操作系统功能的平台无关接口
//
// 【环境变量】
// Getenv(key) - 获取环境变量
// Setenv(key, value) - 设置环境变量
// Environ() - 所有环境变量
//
// 【文件操作】
// Create(name) - 创建文件
// Open(name) - 打开文件（只读）
// OpenFile(name, flag, perm) - 打开文件（指定模式）
// ReadFile(name) - 读取整个文件（Go 1.16+）
// WriteFile(name, data, perm) - 写入整个文件（Go 1.16+）
// Remove(name) - 删除文件
// Rename(old, new) - 重命名
//
// 【目录操作】
// Mkdir(name, perm) - 创建目录
// MkdirAll(path, perm) - 递归创建目录
// ReadDir(name) - 读取目录（Go 1.16+）
// ============================================================================
func osDemo() {
	fmt.Println("\n--- os 包 ---")

	// 环境变量
	fmt.Printf("HOME: %s\n", os.Getenv("HOME"))
	fmt.Printf("PATH 存在: %v\n", os.Getenv("PATH") != "")

	// 工作目录
	wd, _ := os.Getwd()
	fmt.Printf("工作目录: %s\n", wd)

	// 临时目录
	fmt.Printf("临时目录: %s\n", os.TempDir())

	// 命令行参数
	// os.Args[0] 是程序名
	// os.Args[1:] 是参数
	fmt.Printf("程序名: %s\n", os.Args[0])
	fmt.Printf("参数数量: %d\n", len(os.Args))

	// 文件操作示例（仅展示 API）
	fmt.Println("\n文件操作 API:")
	fmt.Println("  os.Create(name)      - 创建文件")
	fmt.Println("  os.Open(name)        - 打开文件（只读）")
	fmt.Println("  os.OpenFile(...)     - 打开文件（指定模式）")
	fmt.Println("  os.ReadFile(name)    - 读取整个文件")
	fmt.Println("  os.WriteFile(...)    - 写入整个文件")
	fmt.Println("  os.Remove(name)      - 删除文件")
	fmt.Println("  os.Rename(old, new)  - 重命名")
	fmt.Println("  os.Mkdir(name, perm) - 创建目录")
	fmt.Println("  os.MkdirAll(...)     - 递归创建目录")
	fmt.Println("  os.Stat(name)        - 获取文件信息")
}

// ============================================================================
// 【io 包】
// ============================================================================
// io 包提供 I/O 原语的基本接口
//
// 【核心接口】
// io.Reader - Read(p []byte) (n int, err error)
// io.Writer - Write(p []byte) (n int, err error)
// io.Closer - Close() error
// io.Seeker - Seek(offset int64, whence int) (int64, error)
//
// 【常用函数】
// io.Copy(dst, src) - 复制数据
// io.ReadAll(r) - 读取所有数据（Go 1.16+）
// io.WriteString(w, s) - 写入字符串
//
// 【bufio 包】
// 带缓冲的 I/O，提高性能
// bufio.NewReader(r) - 带缓冲的 Reader
// bufio.NewWriter(w) - 带缓冲的 Writer
// bufio.NewScanner(r) - 按行读取
// ============================================================================
func ioDemo() {
	fmt.Println("\n--- io 包 ---")

	// bytes.Buffer 实现了 io.Reader 和 io.Writer
	// 【bytes.Buffer】
	// - 可读可写的字节缓冲区
	// - 常用于测试和构建数据
	var buf bytes.Buffer

	// 写入
	buf.WriteString("Hello, ")
	buf.Write([]byte("World!"))
	fmt.Printf("Buffer 内容: %s\n", buf.String())

	// 读取
	data := make([]byte, 5)
	buf.Read(data)
	fmt.Printf("读取 5 字节: %s\n", data)

	// io.Copy
	// 【用途】从 Reader 复制到 Writer
	src := strings.NewReader("Copy this text")
	var dst bytes.Buffer
	io.Copy(&dst, src)
	fmt.Printf("io.Copy 结果: %s\n", dst.String())

	// bufio
	// 【bufio.Reader】
	// - 带缓冲，减少系统调用
	// - 提供便捷的读取方法
	reader := bufio.NewReader(strings.NewReader("line1\nline2\nline3"))
	line, _ := reader.ReadString('\n') // 读取到换行符（包含换行符）
	fmt.Printf("bufio.ReadString: %q\n", line)
}

// ============================================================================
// 【filepath 包】
// ============================================================================
// filepath 包提供文件路径操作，兼容不同操作系统
//
// 【常用函数】
// Dir(path) - 目录部分
// Base(path) - 文件名部分
// Ext(path) - 扩展名
// Join(elem...) - 连接路径
// Split(path) - 分割为目录和文件名
// Clean(path) - 清理路径
// Abs(path) - 绝对路径
// Rel(basepath, targpath) - 相对路径
// Match(pattern, name) - 模式匹配
// Walk(root, fn) - 遍历目录树
// ============================================================================
func filepathDemo() {
	fmt.Println("\n--- filepath 包 ---")

	path := "/home/user/documents/file.txt"

	fmt.Printf("原路径: %s\n", path)
	fmt.Printf("Dir: %s\n", filepath.Dir(path))   // /home/user/documents
	fmt.Printf("Base: %s\n", filepath.Base(path)) // file.txt
	fmt.Printf("Ext: %s\n", filepath.Ext(path))   // .txt
	fmt.Printf("Join: %s\n", filepath.Join("home", "user", "file.txt"))
	fmt.Printf("Split: %v\n", func() (string, string) { return filepath.Split(path) }())

	// 路径清理
	// 【Clean】规范化路径，处理 . 和 ..
	fmt.Printf("Clean: %s\n", filepath.Clean("/a/b/../c/./d")) // /a/c/d
}

// ============================================================================
// 【encoding/json 包】
// ============================================================================
// json 包提供 JSON 编码和解码
//
// 【序列化（Go -> JSON）】
// Marshal(v) - 序列化为 []byte
// MarshalIndent(v, prefix, indent) - 带缩进的序列化
// Encoder.Encode(v) - 流式编码
//
// 【反序列化（JSON -> Go）】
// Unmarshal(data, v) - 反序列化到变量
// Decoder.Decode(v) - 流式解码
//
// 【结构体标签】
// `json:"name"` - 指定 JSON 字段名
// `json:"name,omitempty"` - 空值时省略
// `json:"-"` - 忽略此字段
// `json:"name,string"` - 数字以字符串形式编码
// ============================================================================
func jsonDemo() {
	fmt.Println("\n--- encoding/json 包 ---")

	// 定义结构体，使用 JSON 标签
	type Person struct {
		Name    string   `json:"name"`            // 指定 JSON 字段名
		Age     int      `json:"age"`             //
		Email   string   `json:"email,omitempty"` // 空值时省略
		Hobbies []string `json:"hobbies"`         //
	}

	// 序列化
	person := Person{
		Name:    "Alice",
		Age:     25,
		Hobbies: []string{"reading", "coding"},
		// Email 为空，会被省略
	}

	// 【Marshal】返回紧凑的 JSON
	jsonData, _ := json.Marshal(person)
	fmt.Printf("Marshal: %s\n", jsonData)

	// 【MarshalIndent】返回格式化的 JSON
	prettyJSON, _ := json.MarshalIndent(person, "", "  ")
	fmt.Printf("MarshalIndent:\n%s\n", prettyJSON)

	// 反序列化
	jsonStr := `{"name":"Bob","age":30,"hobbies":["music"]}`
	var p2 Person
	// 【Unmarshal】第二个参数必须是指针
	json.Unmarshal([]byte(jsonStr), &p2)
	fmt.Printf("Unmarshal: %+v\n", p2)

	// 解析到 map
	// 【动态 JSON】不知道结构时，可以解析到 map[string]interface{}
	var m map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &m)
	fmt.Printf("Unmarshal to map: %v\n", m)
}

// ============================================================================
// 【regexp 包】
// ============================================================================
// regexp 包提供正则表达式功能
//
// 【编译正则】
// Compile(expr) - 编译，返回错误
// MustCompile(expr) - 编译，失败则 panic
//
// 【匹配方法】
// MatchString(s) - 是否匹配
// FindString(s) - 第一个匹配
// FindAllString(s, n) - 所有匹配（n=-1 表示全部）
// FindStringSubmatch(s) - 带捕获组的匹配
//
// 【替换方法】
// ReplaceAllString(s, repl) - 替换所有匹配
// ReplaceAllStringFunc(s, f) - 用函数替换
//
// 【常用正则语法】
// .  - 任意字符
// *  - 0 或多个
// +  - 1 或多个
// ?  - 0 或 1 个
// \d - 数字
// \w - 字母数字下划线
// \s - 空白字符
// [] - 字符类
// () - 捕获组
// ============================================================================
func regexpDemo() {
	fmt.Println("\n--- regexp 包 ---")

	// 编译正则表达式
	// 【MustCompile】失败会 panic，适合静态正则
	re := regexp.MustCompile(`\d+`)

	text := "abc123def456"
	fmt.Printf("原文本: %s\n", text)
	fmt.Printf("FindString: %s\n", re.FindString(text))         // 第一个匹配
	fmt.Printf("FindAllString: %v\n", re.FindAllString(text, -1)) // 所有匹配
	fmt.Printf("MatchString: %v\n", re.MatchString(text))       // 是否匹配
	fmt.Printf("ReplaceAllString: %s\n", re.ReplaceAllString(text, "#")) // 替换

	// 邮箱验证
	// 【实用正则示例】
	emailRe := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	fmt.Printf("邮箱验证 'test@example.com': %v\n", emailRe.MatchString("test@example.com"))
}

// ============================================================================
// 【sort 包】
// ============================================================================
// sort 包提供排序功能
//
// 【基本排序】
// Ints(x) - 整数排序
// Float64s(x) - 浮点数排序
// Strings(x) - 字符串排序
//
// 【自定义排序】
// Slice(x, less) - 使用 less 函数排序
// SliceStable(x, less) - 稳定排序
//
// 【搜索】
// SearchInts(a, x) - 二分查找（需要已排序）
// Search(n, f) - 通用二分查找
//
// 【检查】
// IntsAreSorted(x) - 是否已排序
// ============================================================================
func sortDemo() {
	fmt.Println("\n--- sort 包 ---")

	// 整数排序
	ints := []int{5, 2, 8, 1, 9}
	sort.Ints(ints)
	fmt.Printf("Ints: %v\n", ints)

	// 字符串排序
	strs := []string{"banana", "apple", "cherry"}
	sort.Strings(strs)
	fmt.Printf("Strings: %v\n", strs)

	// 浮点数排序
	floats := []float64{3.14, 1.41, 2.72}
	sort.Float64s(floats)
	fmt.Printf("Float64s: %v\n", floats)

	// 自定义排序
	// 【sort.Slice】第二个参数是 less 函数
	// less(i, j) 返回 true 表示 i 应该在 j 前面
	type Person struct {
		Name string
		Age  int
	}
	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
	}
	sort.Slice(people, func(i, j int) bool {
		return people[i].Age < people[j].Age // 按年龄升序
	})
	fmt.Printf("自定义排序: %v\n", people)

	// 检查是否已排序
	fmt.Printf("IsSorted: %v\n", sort.IntsAreSorted(ints))

	// 二分查找
	// 【注意】切片必须已排序
	idx := sort.SearchInts(ints, 5)
	fmt.Printf("SearchInts(5): 索引 %d\n", idx)
}

// ============================================================================
// 【context 包】
// ============================================================================
// context 包定义了 Context 类型，用于在 goroutine 之间传递截止日期、
// 取消信号和请求范围的值
//
// 【创建 Context】
// Background() - 根 context
// TODO() - 占位 context
// WithCancel(parent) - 可取消的 context
// WithTimeout(parent, d) - 带超时的 context
// WithDeadline(parent, t) - 带截止时间的 context
// WithValue(parent, k, v) - 携带值的 context
//
// 【使用规则】
// 1. Context 应该作为函数的第一个参数
// 2. 不要存储 Context
// 3. 不要传递 nil Context
// ============================================================================
func contextDemo() {
	fmt.Println("\n--- context 包 ---")

	// Background 和 TODO
	fmt.Println("context.Background() - 根 context")
	fmt.Println("context.TODO() - 占位符 context")

	// WithCancel
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			fmt.Println("WithCancel: 收到取消信号")
			return
		}
	}(ctx)
	cancel() // 调用 cancel 取消 context
	time.Sleep(10 * time.Millisecond)

	// WithTimeout
	// 【特点】指定时间后自动取消
	ctx2, cancel2 := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel2() // 即使超时，也要调用 cancel 释放资源
	select {
	case <-time.After(100 * time.Millisecond):
		fmt.Println("WithTimeout: 操作完成")
	case <-ctx2.Done():
		fmt.Printf("WithTimeout: %v\n", ctx2.Err())
	}

	// WithValue
	// 【用途】传递请求范围的值（如用户 ID、请求 ID）
	// 【注意】不要用于传递可选参数，只用于跨 API 边界的请求数据
	ctx3 := context.WithValue(context.Background(), "userID", 123)
	fmt.Printf("WithValue: userID=%v\n", ctx3.Value("userID"))
}

// ============================================================================
// 【log 包】
// ============================================================================
// log 包提供简单的日志功能
//
// 【基本函数】
// Print, Println, Printf - 普通日志
// Fatal, Fatalln, Fatalf - 日志后 os.Exit(1)
// Panic, Panicln, Panicf - 日志后 panic
//
// 【自定义 Logger】
// New(out, prefix, flag) - 创建 Logger
//
// 【日志标志】
// Ldate - 日期
// Ltime - 时间
// Lmicroseconds - 微秒
// Llongfile - 完整文件路径
// Lshortfile - 文件名和行号
// LUTC - UTC 时间
// LstdFlags - Ldate | Ltime
// ============================================================================
func logDemo() {
	fmt.Println("\n--- log 包 ---")

	// 基本日志
	log.Println("这是一条日志")

	// 自定义 logger
	// 【参数】(输出目标, 前缀, 标志)
	var buf bytes.Buffer
	logger := log.New(&buf, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println("自定义日志")
	fmt.Printf("日志输出: %s", buf.String())

	// 日志标志
	fmt.Println("\n日志标志:")
	fmt.Println("  log.Ldate      - 日期")
	fmt.Println("  log.Ltime      - 时间")
	fmt.Println("  log.Lmicroseconds - 微秒")
	fmt.Println("  log.Llongfile  - 完整文件路径")
	fmt.Println("  log.Lshortfile - 文件名和行号")
	fmt.Println("  log.LUTC       - UTC 时间")
}

// ============================================================================
// 【flag 包】
// ============================================================================
// flag 包提供命令行参数解析
//
// 【定义标志】
// String(name, default, usage) - 字符串标志
// Int(name, default, usage) - 整数标志
// Bool(name, default, usage) - 布尔标志
// Duration(name, default, usage) - 时间段标志
//
// 【解析】
// Parse() - 解析命令行参数
// Args() - 非标志参数
//
// 【使用模式】
// 1. 定义标志（返回指针）
// 2. 调用 flag.Parse()
// 3. 使用 *flag 获取值
// ============================================================================
func flagDemo() {
	fmt.Println("\n--- flag 包 ---")

	// 注意：这里只是展示 API，实际解析需要命令行参数
	fmt.Println("flag 包用于解析命令行参数:")
	fmt.Println(`
  // 定义标志
  name := flag.String("name", "default", "用户名")
  age := flag.Int("age", 0, "年龄")
  verbose := flag.Bool("v", false, "详细模式")

  // 解析
  flag.Parse()

  // 使用
  fmt.Println(*name, *age, *verbose)

  // 非标志参数
  args := flag.Args()
`)
}

// ============================================================================
// 【net/http 包】
// ============================================================================
// http 包提供 HTTP 客户端和服务器实现
//
// 【服务器】
// HandleFunc(pattern, handler) - 注册处理函数
// ListenAndServe(addr, handler) - 启动服务器
// ListenAndServeTLS(addr, certFile, keyFile, handler) - HTTPS 服务器
//
// 【客户端】
// Get(url) - GET 请求
// Post(url, contentType, body) - POST 请求
// Client.Do(req) - 自定义请求
//
// 【常用类型】
// Request - 请求
// Response - 响应
// ResponseWriter - 响应写入器
// Handler - 处理器接口
// ServeMux - 路由器
// ============================================================================
func httpDemo() {
	fmt.Println("\n--- net/http 包 ---")

	fmt.Println("HTTP 服务器示例:")
	fmt.Println(`
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
      fmt.Fprintf(w, "Hello, World!")
  })
  http.ListenAndServe(":8080", nil)
`)

	fmt.Println("\nHTTP 客户端示例:")
	fmt.Println(`
  resp, err := http.Get("https://api.example.com/data")
  if err != nil {
      log.Fatal(err)
  }
  defer resp.Body.Close()
  body, _ := io.ReadAll(resp.Body)
`)
}

// ============================================================================
// 【math/rand 包】
// ============================================================================
// rand 包提供伪随机数生成
//
// 【Go 1.20+ 变化】
// 不再需要手动设置种子，默认使用随机种子
//
// 【常用函数】
// Int() - 随机 int
// Intn(n) - [0, n) 的随机 int
// Int63() - 随机 int64
// Float64() - [0.0, 1.0) 的随机 float64
// Shuffle(n, swap) - 随机打乱
//
// 【注意】
// - math/rand 不是加密安全的
// - 加密用途请使用 crypto/rand
// ============================================================================
func randDemo() {
	fmt.Println("\n--- math/rand 包 ---")

	// 注意：Go 1.20+ 默认自动设置随机种子
	// 【Intn】返回 [0, n) 的随机整数
	fmt.Printf("随机整数: %d\n", rand.Intn(100))

	// 【Float64】返回 [0.0, 1.0) 的随机浮点数
	fmt.Printf("随机浮点数: %.4f\n", rand.Float64())

	// 随机打乱切片
	// 【Shuffle】第一个参数是长度，第二个是交换函数
	nums := []int{1, 2, 3, 4, 5}
	rand.Shuffle(len(nums), func(i, j int) {
		nums[i], nums[j] = nums[j], nums[i]
	})
	fmt.Printf("打乱后: %v\n", nums)

	// 生成随机字符串
	const charset = "abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, 8)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	fmt.Printf("随机字符串: %s\n", string(result))
}
