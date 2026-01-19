// ============================================================================
// 12_concurrency.go - 并发编程
// ============================================================================
// 运行: go run 12_concurrency.go
//
// 【本文件学习目标】
// 1. 理解 goroutine 的概念和使用
// 2. 掌握 channel 的创建和通信
// 3. 学会使用 select 进行多路复用
// 4. 理解 sync 包的同步原语
// 5. 掌握 context 的使用场景
// 6. 了解常见的并发模式
//
// 【Go 并发的核心理念】
// "不要通过共享内存来通信，而要通过通信来共享内存"
// - Don't communicate by sharing memory; share memory by communicating
// - 优先使用 channel 进行 goroutine 间通信
// - 必要时使用 sync 包的同步原语
//
// 【Go 并发模型：CSP】
// CSP = Communicating Sequential Processes（通信顺序进程）
// - goroutine: 轻量级线程（用户态）
// - channel: goroutine 间通信的管道
//
// 【goroutine vs 线程】
// | 特性       | goroutine        | OS 线程          |
// |------------|------------------|------------------|
// | 栈大小     | 2KB（可增长）    | 1-8 MB（固定）   |
// | 切换成本   | 几十纳秒         | 几微秒           |
// | 调度       | Go 运行时调度    | OS 调度          |
// | 数量       | 可以创建百万级   | 几千个就吃力     |
// | 创建成本   | 极低             | 较高             |
// ============================================================================

package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	fmt.Println("=== Go 并发编程 ===\n")

	// ========================================================================
	// 【Goroutine 基础】
	// ========================================================================
	// goroutine 是 Go 的轻量级线程
	//
	// 【启动语法】
	// go 函数调用  // 启动一个新的 goroutine
	//
	// 【特点】
	// - 由 Go 运行时调度，不是 OS 线程
	// - 初始栈只有 2KB，可以动态增长
	// - 创建成本极低，可以轻松创建成千上万个
	//
	// 【注意】
	// - goroutine 是非阻塞的，go 语句立即返回
	// - 主 goroutine 结束，程序就结束，不等待其他 goroutine
	// - 需要使用同步机制等待 goroutine 完成
	// ========================================================================
	fmt.Println("--- Goroutine 基础 ---")

	// 启动一个 goroutine
	// 【执行顺序】
	// - go 语句立即返回
	// - sayHello 在新 goroutine 中异步执行
	// - 主 goroutine 继续向下执行
	go sayHello("Goroutine")

	// 主 goroutine 继续执行
	fmt.Println("主 goroutine")

	// 等待一下让 goroutine 执行
	// 【警告】生产代码不要用 Sleep 做同步！
	// 应该使用 WaitGroup、channel 或 Context
	time.Sleep(100 * time.Millisecond)

	// 启动多个 goroutine
	// 【闭包陷阱】
	// 错误写法：
	// for i := 1; i <= 3; i++ {
	//     go func() { fmt.Println(i) }()  // 所有 goroutine 可能都打印 4！
	// }
	// 原因：闭包捕获的是变量 i 的引用，循环结束时 i=4
	//
	// 正确写法：
	// 1. 将 i 作为参数传递（推荐）
	// 2. 在循环内创建局部变量
	fmt.Println("\n启动多个 goroutine:")
	for i := 1; i <= 3; i++ {
		go func(n int) {
			fmt.Printf("  Goroutine %d\n", n)
		}(i) // 注意：传递 i 的副本避免闭包陷阱
	}
	time.Sleep(100 * time.Millisecond)

	// ========================================================================
	// 【Channel 基础】
	// ========================================================================
	// channel 是 goroutine 间通信的管道
	//
	// 【创建语法】
	// ch := make(chan Type)        // 无缓冲 channel
	// ch := make(chan Type, size)  // 带缓冲 channel
	//
	// 【操作语法】
	// ch <- value  // 发送
	// value := <-ch  // 接收
	// close(ch)    // 关闭
	//
	// 【无缓冲 vs 带缓冲】
	// | 类型       | 发送阻塞条件           | 接收阻塞条件     |
	// |------------|------------------------|------------------|
	// | 无缓冲     | 没有接收者             | 没有发送者       |
	// | 带缓冲     | 缓冲区满               | 缓冲区空         |
	//
	// 【无缓冲 channel 的同步特性】
	// 发送和接收必须同时准备好，类似"握手"
	// ========================================================================
	fmt.Println("\n--- Channel 基础 ---")

	// 创建无缓冲 channel
	// 【无缓冲的特点】
	// - 发送操作会阻塞，直到有接收者
	// - 接收操作会阻塞，直到有发送者
	// - 强制同步，确保数据交换
	ch := make(chan string)

	// 发送数据
	// 必须在另一个 goroutine 中发送，否则会死锁
	go func() {
		ch <- "Hello from channel" // 发送
	}()

	// 接收数据
	// <- ch 会阻塞，直到收到数据
	msg := <-ch
	fmt.Printf("收到: %s\n", msg)

	// 带缓冲的 channel
	// 【带缓冲的特点】
	// - 发送不会立即阻塞，直到缓冲区满
	// - 适合"生产者-消费者"模式
	// - cap(ch) 返回容量，len(ch) 返回当前元素数
	bufferedCh := make(chan int, 3) // 容量为 3
	bufferedCh <- 1                 // 不阻塞
	bufferedCh <- 2                 // 不阻塞
	bufferedCh <- 3                 // 不阻塞
	// bufferedCh <- 4 // 会阻塞，因为缓冲区已满

	fmt.Printf("缓冲 channel 长度: %d, 容量: %d\n", len(bufferedCh), cap(bufferedCh))

	// ========================================================================
	// 【Channel 操作】
	// ========================================================================
	// channel 的关闭和遍历
	//
	// 【关闭 channel】
	// - close(ch) 关闭 channel
	// - 关闭后不能再发送，但可以继续接收
	// - 接收已关闭的 channel 返回零值
	// - 重复关闭会 panic
	//
	// 【检查是否关闭】
	// value, ok := <-ch
	// - ok=true: channel 未关闭或有数据
	// - ok=false: channel 已关闭且无数据
	//
	// 【range 遍历】
	// for v := range ch { }
	// - 会一直接收直到 channel 关闭
	// - 如果不 close(ch)，会永远阻塞（死锁）
	// ========================================================================
	fmt.Println("\n--- Channel 操作 ---")

	// 关闭 channel
	dataCh := make(chan int, 5)
	go func() {
		for i := 1; i <= 5; i++ {
			dataCh <- i
		}
		close(dataCh) // 关闭 channel，通知接收方数据已发完
	}()

	// 使用 range 遍历 channel
	// 【range 的工作原理】
	// - 持续从 channel 接收数据
	// - channel 关闭时自动退出循环
	fmt.Print("从 channel 接收: ")
	for v := range dataCh {
		fmt.Printf("%d ", v)
	}
	fmt.Println()

	// 检查 channel 是否关闭
	// 【comma ok 模式】
	closedCh := make(chan int)
	close(closedCh)
	val, ok := <-closedCh
	fmt.Printf("从关闭的 channel 接收: val=%d, ok=%v\n", val, ok)
	// val=0 (零值), ok=false (已关闭)

	// ========================================================================
	// 【Select 语句】
	// ========================================================================
	// select 用于同时等待多个 channel 操作
	//
	// 【语法】
	// select {
	// case v := <-ch1:
	//     // ch1 可读
	// case ch2 <- v:
	//     // ch2 可写
	// case <-time.After(timeout):
	//     // 超时
	// default:
	//     // 所有 case 都不满足时执行
	// }
	//
	// 【特点】
	// - 多个 case 同时满足时，随机选择一个
	// - 没有 default 时，select 会阻塞
	// - 有 default 时，如果没有 case 满足，执行 default
	//
	// 【常见用途】
	// - 多路复用（同时监听多个 channel）
	// - 超时控制
	// - 非阻塞操作
	// ========================================================================
	fmt.Println("\n--- Select 语句 ---")

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(50 * time.Millisecond)
		ch1 <- "来自 ch1"
	}()
	go func() {
		time.Sleep(30 * time.Millisecond)
		ch2 <- "来自 ch2"
	}()

	// select 等待多个 channel
	// 【工作原理】
	// - 同时监听 ch1 和 ch2
	// - 哪个先准备好就执行哪个 case
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Printf("收到: %s\n", msg1)
		case msg2 := <-ch2:
			fmt.Printf("收到: %s\n", msg2)
		}
	}

	// select 带超时
	// 【time.After】
	// 返回一个 channel，在指定时间后发送当前时间
	// 常用于实现超时控制
	fmt.Println("\nselect 带超时:")
	timeoutCh := make(chan string)
	select {
	case msg := <-timeoutCh:
		fmt.Printf("收到: %s\n", msg)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("超时!")
	}

	// select 非阻塞操作
	// 【default case】
	// 当所有 case 都不满足时，立即执行 default
	// 实现非阻塞的发送/接收
	fmt.Println("\nselect 非阻塞:")
	nonBlockingCh := make(chan int)
	select {
	case v := <-nonBlockingCh:
		fmt.Printf("收到: %d\n", v)
	default:
		fmt.Println("channel 为空，继续执行")
	}

	// ========================================================================
	// 【单向 Channel】
	// ========================================================================
	// 可以将双向 channel 转换为单向 channel
	//
	// 【类型】
	// chan T     // 双向 channel
	// <-chan T   // 只读 channel（只能接收）
	// chan<- T   // 只写 channel（只能发送）
	//
	// 【用途】
	// - 限制函数对 channel 的操作
	// - 函数签名更清晰
	// - 编译时检查错误使用
	//
	// 【转换规则】
	// - 双向可以转单向
	// - 单向不能转双向
	// ========================================================================
	fmt.Println("\n--- 单向 Channel ---")

	// 生产者-消费者模式
	// generate 返回 <-chan int (只读)
	// square 接收 <-chan int，返回 <-chan int
	numbers := generate(5)     // 返回 <-chan int (只读)
	squares := square(numbers) // 接收 <-chan，返回 <-chan
	for n := range squares {
		fmt.Printf("%d ", n)
	}
	fmt.Println()

	// ========================================================================
	// 【sync.WaitGroup】
	// ========================================================================
	// WaitGroup 用于等待一组 goroutine 完成
	//
	// 【方法】
	// wg.Add(n)   // 增加计数器
	// wg.Done()   // 减少计数器（等价于 Add(-1)）
	// wg.Wait()   // 阻塞直到计数器为 0
	//
	// 【使用模式】
	// var wg sync.WaitGroup
	// for ... {
	//     wg.Add(1)
	//     go func() {
	//         defer wg.Done()
	//         // 工作
	//     }()
	// }
	// wg.Wait()
	//
	// 【注意】
	// - Add 必须在 go 语句之前调用
	// - Done 通常用 defer 调用，确保一定会执行
	// - 不要复制 WaitGroup（传指针）
	// ========================================================================
	fmt.Println("\n--- sync.WaitGroup ---")

	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1) // 必须在 go 之前调用
		go func(n int) {
			defer wg.Done() // 确保完成时调用 Done
			fmt.Printf("Worker %d 开始\n", n)
			time.Sleep(50 * time.Millisecond)
			fmt.Printf("Worker %d 完成\n", n)
		}(i)
	}

	wg.Wait() // 等待所有 goroutine 完成
	fmt.Println("所有 worker 完成")

	// ========================================================================
	// 【sync.Mutex】
	// ========================================================================
	// Mutex 是互斥锁，用于保护共享资源
	//
	// 【方法】
	// mu.Lock()    // 加锁
	// mu.Unlock()  // 解锁
	//
	// 【使用模式】
	// mu.Lock()
	// defer mu.Unlock()
	// // 访问共享资源
	//
	// 【注意】
	// - 未加锁时调用 Unlock 会 panic
	// - 不要复制 Mutex（传指针）
	// - 锁的粒度要适当，太大影响性能
	// - 推荐用 defer 解锁，防止忘记解锁
	//
	// 【何时使用】
	// - 需要保护共享数据的读写
	// - 不能用 channel 简单实现时
	// ========================================================================
	fmt.Println("\n--- sync.Mutex ---")

	counter := &SafeCounter{}

	var wg2 sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			counter.Increment() // 安全地增加计数
		}()
	}
	wg2.Wait()
	fmt.Printf("安全计数器: %d\n", counter.Value())

	// ========================================================================
	// 【sync.RWMutex】
	// ========================================================================
	// RWMutex 是读写锁，允许多个读者或一个写者
	//
	// 【方法】
	// mu.RLock()   // 读锁定
	// mu.RUnlock() // 读解锁
	// mu.Lock()    // 写锁定
	// mu.Unlock()  // 写解锁
	//
	// 【特点】
	// - 多个读者可以同时持有读锁
	// - 写者持有写锁时，其他读写都被阻塞
	// - 适合读多写少的场景
	//
	// 【选择指南】
	// | 场景           | 推荐使用     |
	// |----------------|--------------|
	// | 读写都频繁     | Mutex        |
	// | 读多写少       | RWMutex      |
	// | 只读           | 不需要锁     |
	// ========================================================================
	fmt.Println("\n--- sync.RWMutex ---")

	cache := &Cache{data: make(map[string]string)}
	cache.Set("key1", "value1")                       // 使用写锁
	fmt.Printf("cache.Get(\"key1\") = %s\n", cache.Get("key1")) // 使用读锁

	// ========================================================================
	// 【sync.Once】
	// ========================================================================
	// Once 确保函数只执行一次，常用于单例初始化
	//
	// 【方法】
	// once.Do(func)  // 只有第一次调用会执行 func
	//
	// 【特点】
	// - 线程安全
	// - 即使有多个 goroutine 同时调用，func 也只执行一次
	// - 常用于延迟初始化
	//
	// 【使用场景】
	// - 单例模式
	// - 延迟初始化
	// - 一次性的配置加载
	// ========================================================================
	fmt.Println("\n--- sync.Once ---")

	var once sync.Once
	initialize := func() {
		fmt.Println("初始化（只执行一次）")
	}

	var wg3 sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg3.Add(1)
		go func() {
			defer wg3.Done()
			once.Do(initialize) // 只有第一个到达的 goroutine 会执行
		}()
	}
	wg3.Wait()

	// ========================================================================
	// 【sync/atomic】
	// ========================================================================
	// atomic 包提供原子操作，无需加锁
	//
	// 【常用函数】
	// atomic.AddInt64(&val, delta)  // 原子加
	// atomic.LoadInt64(&val)        // 原子读
	// atomic.StoreInt64(&val, new)  // 原子写
	// atomic.CompareAndSwapInt64(&val, old, new)  // CAS
	//
	// 【特点】
	// - 比 Mutex 更轻量
	// - 只能用于简单的数值操作
	// - 底层使用 CPU 指令保证原子性
	//
	// 【何时使用】
	// - 简单的计数器
	// - 标志位
	// - 单个变量的读写
	// - 复杂操作还是用 Mutex
	// ========================================================================
	fmt.Println("\n--- sync/atomic ---")

	var atomicCounter int64 = 0

	var wg4 sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg4.Add(1)
		go func() {
			defer wg4.Done()
			atomic.AddInt64(&atomicCounter, 1) // 原子加 1
		}()
	}
	wg4.Wait()
	fmt.Printf("原子计数器: %d\n", atomicCounter)

	// ========================================================================
	// 【Context】
	// ========================================================================
	// context 用于控制 goroutine 的生命周期
	//
	// 【创建 Context】
	// context.Background()   // 根 context，通常在 main 中使用
	// context.TODO()         // 占位，表示未来会替换
	// context.WithCancel(parent)    // 带取消功能
	// context.WithTimeout(parent, d) // 带超时
	// context.WithDeadline(parent, t) // 带截止时间
	// context.WithValue(parent, k, v) // 携带值
	//
	// 【使用规则】
	// 1. Context 应该作为函数的第一个参数
	// 2. 不要存储 Context，而是传递它
	// 3. 不要传递 nil Context，用 context.TODO()
	//
	// 【监听取消】
	// select {
	// case <-ctx.Done():
	//     return ctx.Err()
	// default:
	//     // 继续工作
	// }
	// ========================================================================
	fmt.Println("\n--- Context ---")

	// 带取消的 context
	// 【WithCancel】
	// - 返回 ctx 和 cancel 函数
	// - 调用 cancel() 会取消 ctx
	// - ctx.Done() 返回的 channel 会关闭
	ctx, cancel := context.WithCancel(context.Background())
	go worker(ctx, "worker-1")
	time.Sleep(100 * time.Millisecond)
	cancel() // 取消
	time.Sleep(50 * time.Millisecond)

	// 带超时的 context
	// 【WithTimeout】
	// - 指定时间后自动取消
	// - 仍然要 defer cancel() 释放资源
	ctx2, cancel2 := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel2()

	select {
	case <-time.After(100 * time.Millisecond):
		fmt.Println("操作完成")
	case <-ctx2.Done():
		fmt.Printf("Context 超时: %v\n", ctx2.Err())
	}

	// ========================================================================
	// 【常见并发模式】
	// ========================================================================
	// Go 社区总结了很多并发模式
	//
	// 【扇出/扇入】(Fan-out/Fan-in)
	// - 扇出：一个输入，多个 goroutine 处理
	// - 扇入：多个输出合并为一个
	//
	// 【Worker Pool】
	// - 固定数量的 worker 处理任务队列
	// - 控制并发数量
	//
	// 【管道】(Pipeline)
	// - 多个阶段串联
	// - 每个阶段是一个 goroutine
	//
	// 【超时控制】
	// - 使用 context.WithTimeout
	// - 使用 time.After
	// ========================================================================
	fmt.Println("\n--- 常见并发模式 ---")

	// 1. 扇出/扇入
	fmt.Println("扇出/扇入模式:")
	fanOutFanIn()

	// 2. Worker Pool
	fmt.Println("\nWorker Pool 模式:")
	workerPool(3, 10)

	// ========================================================================
	// 【运行时信息】
	// ========================================================================
	// runtime 包提供与 Go 运行时交互的功能
	//
	// 【常用函数】
	// runtime.GOMAXPROCS(n) // 设置/获取最大 P 数量
	// runtime.NumCPU()      // CPU 核心数
	// runtime.NumGoroutine() // 当前 goroutine 数量
	// runtime.Gosched()     // 让出 CPU
	// runtime.Goexit()      // 退出当前 goroutine
	//
	// 【GOMAXPROCS】
	// - P 是 Go 调度器的处理器
	// - 默认等于 CPU 核心数
	// - 通常不需要修改
	// ========================================================================
	fmt.Println("\n--- 运行时信息 ---")
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0)) // 0 表示查询当前值
	fmt.Printf("NumCPU: %d\n", runtime.NumCPU())
	fmt.Printf("NumGoroutine: %d\n", runtime.NumGoroutine())
}

// ============================================================================
// 【辅助函数和类型】
// ============================================================================

// sayHello: 简单的打印函数
func sayHello(name string) {
	fmt.Printf("Hello from %s!\n", name)
}

// ============================================================================
// 【生成器模式】
// ============================================================================
// 生成器是一种常见的并发模式
// 函数返回一个只读 channel，调用者从中接收数据
//
// 【特点】
// - 生产者和消费者解耦
// - 消费者按需接收
// - 生产者在 goroutine 中运行
// ============================================================================

// generate: 生成器，产生 1 到 max 的数字
// 返回只读 channel <-chan int
func generate(max int) <-chan int {
	out := make(chan int)
	go func() {
		for i := 1; i <= max; i++ {
			out <- i
		}
		close(out) // 关闭通知消费者数据已发完
	}()
	return out
}

// ============================================================================
// 【管道阶段】
// ============================================================================
// 管道的每个阶段：
// 1. 从输入 channel 接收数据
// 2. 处理数据
// 3. 发送到输出 channel
// ============================================================================

// square: 管道阶段，计算平方
// 接收只读 channel，返回只读 channel
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out) // 输入关闭后，关闭输出
	}()
	return out
}

// ============================================================================
// 【线程安全的计数器】
// ============================================================================
// 使用 Mutex 保护共享数据
// 【设计模式】
// - 将 mutex 和它保护的数据放在同一个结构体中
// - 提供方法封装锁的操作
// - 使用 defer 确保解锁
// ============================================================================

// SafeCounter: 线程安全的计数器
type SafeCounter struct {
	mu    sync.Mutex // 互斥锁
	count int        // 被保护的数据
}

// Increment: 安全地增加计数
func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

// Value: 安全地读取计数
func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// ============================================================================
// 【带读写锁的缓存】
// ============================================================================
// 使用 RWMutex 优化读多写少的场景
// - Get 用读锁，允许并发读
// - Set 用写锁，独占写入
// ============================================================================

// Cache: 带读写锁的缓存
type Cache struct {
	mu   sync.RWMutex      // 读写锁
	data map[string]string // 被保护的数据
}

// Get: 读取缓存（使用读锁）
func (c *Cache) Get(key string) string {
	c.mu.RLock() // 读锁
	defer c.mu.RUnlock()
	return c.data[key]
}

// Set: 写入缓存（使用写锁）
func (c *Cache) Set(key, value string) {
	c.mu.Lock() // 写锁
	defer c.mu.Unlock()
	c.data[key] = value
}

// ============================================================================
// 【Context Worker】
// ============================================================================
// 响应 Context 取消的 worker
// 【模式】
// select {
// case <-ctx.Done():
//     return  // 收到取消信号，退出
// default:
//     // 继续工作
// }
// ============================================================================

// worker: 响应 context 取消的 worker
func worker(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done(): // 监听取消信号
			fmt.Printf("%s: 收到取消信号\n", name)
			return
		default:
			fmt.Printf("%s: 工作中...\n", name)
			time.Sleep(30 * time.Millisecond)
		}
	}
}

// ============================================================================
// 【扇出/扇入模式】
// ============================================================================
// 扇出：将工作分发给多个 goroutine
// 扇入：将多个 goroutine 的结果合并
//
// 【用途】
// - 并行处理
// - 提高吞吐量
// - CPU 密集型任务
// ============================================================================

// fanOutFanIn: 演示扇出/扇入模式
func fanOutFanIn() {
	in := generate(5)

	// 扇出: 启动多个 goroutine 处理
	// 注意：这里两个 square 共享同一个 in channel
	// 每个数字只会被一个 square 处理
	c1 := square(in)
	c2 := square(in)

	// 扇入: 合并结果
	for n := range merge(c1, c2) {
		fmt.Printf("%d ", n)
	}
	fmt.Println()
}

// merge: 合并多个 channel 为一个
// 【实现要点】
// - 为每个输入 channel 启动一个 goroutine
// - 使用 WaitGroup 等待所有输入关闭
// - 所有输入关闭后，关闭输出
func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// 为每个输入 channel 启动一个 goroutine
	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			out <- n
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// 等待所有输入关闭，然后关闭输出
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// ============================================================================
// 【Worker Pool 模式】
// ============================================================================
// 固定数量的 worker 处理任务队列
//
// 【组件】
// - jobs channel: 任务队列
// - results channel: 结果队列
// - worker: 从 jobs 取任务，处理后发到 results
//
// 【优点】
// - 控制并发数量
// - 复用 goroutine
// - 避免创建过多 goroutine
// ============================================================================

// workerPool: 演示 worker pool 模式
func workerPool(numWorkers, numJobs int) {
	jobs := make(chan int, numJobs)    // 任务队列
	results := make(chan int, numJobs) // 结果队列

	// 启动固定数量的 worker
	for w := 1; w <= numWorkers; w++ {
		go poolWorker(w, jobs, results)
	}

	// 发送任务
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // 关闭任务队列，通知 worker 没有更多任务

	// 收集结果
	fmt.Print("结果: ")
	for r := 1; r <= numJobs; r++ {
		fmt.Printf("%d ", <-results)
	}
	fmt.Println()
}

// poolWorker: worker pool 中的 worker
// 【参数】
// - id: worker ID（用于调试）
// - jobs: 只读的任务 channel
// - results: 只写的结果 channel
func poolWorker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		// 处理任务
		results <- j * 2
	}
}
