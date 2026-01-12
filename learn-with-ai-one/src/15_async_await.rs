// Rust异步编程（async/await）详解
// 本文件介绍Rust中的异步编程模型，包括async函数、await表达式、Future等概念

use std::time::Duration;
use std::thread;

// 注意：实际运行异步代码需要使用tokio或async-std等运行时
// 这里我们使用模拟的方式展示异步编程的概念

pub fn run_example() {
    println!("=== Rust异步编程（async/await）===\n");

    // 基本异步概念
    basic_async_concepts();
    
    // 异步函数和await
    async_functions_and_await();
    
    // 异步块
    async_blocks();
    
    // 错误处理
    async_error_handling();
    
    // 并发执行
    async_concurrency();
    
    // 超时处理
    async_timeout();
    
    // 异步流（Stream）概念
    async_stream_concepts();
    
    // 实际应用示例
    async_practical_example();
    
    println!("\n=== 示例结束 ===");
}

fn main() {
    run_example();
}

// 1. 基本异步概念
fn basic_async_concepts() {
    println!("1. 基本异步概念:");
    
    println!("- 异步编程允许程序在等待某个操作完成时执行其他任务");
    println!("- Rust使用Future trait表示异步操作的结果");
    println!("- async/await语法提供了更简洁的异步编程方式");
    println!("- 需要运行时（如tokio、async-std）来执行异步代码");
    
    // 模拟异步操作
    simulate_async_operation();
    
    println!();
}

// 模拟异步操作
fn simulate_async_operation() {
    println!("\n模拟异步操作:");
    
    // 主线程继续执行
    println!("主线程继续执行...");
    
    // 创建新线程模拟异步操作
    thread::spawn(|| {
        thread::sleep(Duration::from_millis(200));
        println!("异步操作完成!");
    });
    
    // 主线程继续执行
    println!("主线程继续做其他事情...");
    
    // 等待足够长的时间让异步操作完成
    thread::sleep(Duration::from_millis(300));
}

// 2. 异步函数和await
fn async_functions_and_await() {
    println!("2. 异步函数和await:");
    
    println!("- async关键字用于定义异步函数");
    println!("- 异步函数返回Future trait的实现");
    println!("- await关键字用于等待异步操作完成");
    println!("- await只能在async函数或async块中使用");
    
    // 模拟异步函数调用
    simulate_async_function();
    
    println!();
}

// 模拟异步函数
fn simulate_async_function() {
    println!("\n模拟异步函数调用:");
    
    // 异步函数定义（这里使用普通函数模拟）
    fn async_function(name: &str) -> String {
        thread::sleep(Duration::from_millis(150));
        format!("{name} 完成")
    }
    
    // 模拟异步调用链
    let result1 = async_function("任务1");
    println!("任务1结果: {result1}");
    
    let result2 = async_function("任务2");
    println!("任务2结果: {result2}");
    
    let combined = format!("{result1}，{result2}");
    println!("组合结果: {combined}");
}

// 3. 异步块
fn async_blocks() {
    println!("3. 异步块:");
    
    println!("- async块创建一个Future实例");
    println!("- 语法: async {{ /* 异步代码 */ }}");
    println!("- 可以在任何地方使用，不局限于函数内部");
    
    // 模拟异步块
    simulate_async_block();
    
    println!();
}

// 模拟异步块
fn simulate_async_block() {
    println!("\n模拟异步块:");
    
    // 定义一个模拟的异步块
    let async_block = || {
        thread::sleep(Duration::from_millis(100));
        "异步块执行完成"
    };
    
    // 执行异步块
    let result = async_block();
    println!("异步块结果: {result}");
    
    // 嵌套模拟
    let nested_async = || {
        thread::sleep(Duration::from_millis(50));
        let inner_result = async_block();
        format!("外层结果 + {inner_result}")
    };
    
    let nested_result = nested_async();
    println!("嵌套异步块结果: {nested_result}");
}

// 4. 错误处理
fn async_error_handling() {
    println!("4. 异步错误处理:");
    
    println!("- 异步函数可以返回Result<T, E>类型");
    println!("- 使用?操作符传播错误");
    println!("- 需要处理Future中的错误");
    
    // 模拟异步错误处理
    simulate_async_error_handling();
    
    println!();
}

// 模拟异步错误处理
fn simulate_async_error_handling() {
    println!("\n模拟异步错误处理:");
    
    // 模拟可能失败的异步操作
    fn async_operation_with_error(success: bool) -> Result<String, String> {
        thread::sleep(Duration::from_millis(100));
        
        if success {
            Ok("操作成功".to_string())
        } else {
            Err("操作失败".to_string())
        }
    }
    
    // 成功的情况
    match async_operation_with_error(true) {
        Ok(result) => println!("成功情况: {result}"),
        Err(e) => println!("成功情况错误: {e}"),
    }
    
    // 失败的情况
    match async_operation_with_error(false) {
        Ok(result) => println!("失败情况: {result}"),
        Err(e) => println!("失败情况错误: {e}"),
    }
    
    // 模拟链式调用的错误传播
    fn chained_async_operations() -> Result<String, String> {
        let result1 = async_operation_with_error(true)?;
        let result2 = async_operation_with_error(false)?; // 这里会失败
        Ok(format!("{result1}, {result2}"))
    }
    
    match chained_async_operations() {
        Ok(result) => println!("链式调用结果: {result}"),
        Err(e) => println!("链式调用错误: {e}"),
    }
}

// 5. 并发执行
fn async_concurrency() {
    println!("5. 并发执行:");
    
    println!("- 使用join!宏并发执行多个异步任务");
    println!("- 所有任务完成后才继续执行");
    println!("- 任务之间可以共享数据（需要适当的同步机制）");
    
    // 模拟并发执行
    simulate_async_concurrency();
    
    println!();
}

// 模拟并发执行
fn simulate_async_concurrency() {
    println!("\n模拟并发执行:");
    
    // 定义三个异步任务
    let task1 = || {
        thread::sleep(Duration::from_millis(200));
        "任务1完成"
    };
    
    let task2 = || {
        thread::sleep(Duration::from_millis(150));
        "任务2完成"
    };
    
    let task3 = || {
        thread::sleep(Duration::from_millis(250));
        "任务3完成"
    };
    
    // 并发执行（使用线程模拟）
    let start = std::time::Instant::now();
    
    let handle1 = thread::spawn(task1);
    let handle2 = thread::spawn(task2);
    let handle3 = thread::spawn(task3);
    
    // 等待所有任务完成
    let result1 = handle1.join().unwrap();
    let result2 = handle2.join().unwrap();
    let result3 = handle3.join().unwrap();
    
    let elapsed = start.elapsed();
    
    println!("任务1结果: {result1}");
    println!("任务2结果: {result2}");
    println!("任务3结果: {result3}");
    println!("总耗时: {:?}", elapsed);
    println!("如果串行执行，预计耗时约600ms");
}

// 6. 超时处理
fn async_timeout() {
    println!("6. 超时处理:");
    
    println!("- 使用timeout!宏或类似机制设置异步操作的超时");
    println!("- 避免长时间等待导致的资源浪费");
    println!("- 提高系统的响应性和稳定性");
    
    // 模拟超时处理
    simulate_async_timeout();
    
    println!();
}

// 模拟超时处理
fn simulate_async_timeout() {
    println!("\n模拟超时处理:");
    
    // 模拟可能超时的异步操作
    fn async_operation_with_timeout(duration: Duration) -> Result<String, String> {
        thread::sleep(duration);
        Ok("操作完成".to_string())
    }
    
    // 模拟超时机制
    fn with_timeout<F, T>(f: F, timeout: Duration) -> Result<T, String>
    where
        F: FnOnce() -> Result<T, String> + Send + 'static,
        T: Send + 'static,
    {
        let handle = thread::spawn(f);
        
        match handle.join() {
            Ok(result) => result,
            Err(_) => Err("线程执行错误".to_string()),
        }
    }
    
    // 不超时的情况
    let start1 = std::time::Instant::now();
    match with_timeout(|| async_operation_with_timeout(Duration::from_millis(100)), Duration::from_millis(200)) {
        Ok(result) => println!("不超时情况: {result}"),
        Err(e) => println!("不超时情况错误: {e}"),
    }
    let elapsed1 = start1.elapsed();
    println!("不超时耗时: {:?}", elapsed1);
    
    println!();
}

// 7. 异步流（Stream）概念
fn async_stream_concepts() {
    println!("7. 异步流（Stream）概念:");
    
    println!("- Stream表示异步产生的一系列值");
    println!("- 类似于迭代器，但值是异步产生的");
    println!("- 使用next().await获取下一个值");
    println!("- 可以与map、filter等操作符一起使用");
    
    // 模拟异步流
    simulate_async_stream();
    
    println!();
}

// 模拟异步流
fn simulate_async_stream() {
    println!("\n模拟异步流:");
    
    // 模拟Stream的迭代器
    struct MockStream {
        current: u32,
        max: u32,
    }
    
    impl MockStream {
        fn new(max: u32) -> Self {
            MockStream { current: 0, max }
        }
        
        // 模拟next().await
        fn next(&mut self) -> Option<u32> {
            thread::sleep(Duration::from_millis(50));
            
            if self.current < self.max {
                let value = self.current;
                self.current += 1;
                Some(value)
            } else {
                None
            }
        }
    }
    
    // 使用模拟流
    let mut stream = MockStream::new(5);
    
    println!("从流中获取值:");
    while let Some(value) = stream.next() {
        println!("获取到值: {}", value);
    }
    
    println!("流结束");
}

// 8. 实际应用示例
fn async_practical_example() {
    println!("8. 实际应用示例:");
    
    println!("异步编程在实际应用中的常见场景:");
    println!("- 网络请求和API调用");
    println!("- 文件I/O操作");
    println!("- 数据库查询");
    println!("- 并发任务处理");
    println!("- Web服务器和客户端");
    
    // 模拟异步Web请求
    simulate_async_web_request();
    
    println!();
}

// 模拟异步Web请求
fn simulate_async_web_request() {
    println!("\n模拟异步Web请求:");
    
    // 模拟HTTP客户端
    struct MockHttpClient;
    
    impl MockHttpClient {
        fn new() -> Self {
            MockHttpClient
        }
        
        // 模拟异步GET请求
        fn get(&self, url: &str) -> Result<String, String> {
            println!("发送GET请求到: {}", url);
            thread::sleep(Duration::from_millis(150));
            
            Ok(format!("{url} 的响应内容"))
        }
        
        // 模拟异步POST请求
        fn post(&self, url: &str, data: &str) -> Result<String, String> {
            println!("发送POST请求到: {}，数据: {}", url, data);
            thread::sleep(Duration::from_millis(200));
            
            Ok(format!("POST请求成功，响应: {data}"))
        }
    }
    
    // 模拟异步Web服务客户端
    async fn simulate_web_client() {
        // 并发发送多个请求（模拟）
        // 为每个线程创建一个新的客户端实例
        let handle1 = thread::spawn(|| MockHttpClient::new().get("https://api.example.com/users"));
        let handle2 = thread::spawn(|| MockHttpClient::new().get("https://api.example.com/products"));
        let handle3 = thread::spawn(|| MockHttpClient::new().post("https://api.example.com/orders", "{\"item\": \"book\"}"));
        
        // 等待所有请求完成
        let result1 = handle1.join().unwrap();
        let result2 = handle2.join().unwrap();
        let result3 = handle3.join().unwrap();
        
        // 处理结果
        println!("\n处理请求结果:");
        match result1 {
            Ok(response) => println!("用户API响应: {}", response),
            Err(e) => println!("用户API错误: {}", e),
        }
        
        match result2 {
            Ok(response) => println!("产品API响应: {}", response),
            Err(e) => println!("产品API错误: {}", e),
        }
        
        match result3 {
            Ok(response) => println!("订单API响应: {}", response),
            Err(e) => println!("订单API错误: {}", e),
        }
    }
    
    // 运行模拟的Web客户端
    simulate_web_client();
    
    println!("\n异步Web客户端模拟完成");
}