// 10_error_handling.rs - Rust错误处理机制详解

// Rust的错误处理系统分为两类：
// 1. 不可恢复错误（Unrecoverable Errors）：使用panic!宏
// 2. 可恢复错误（Recoverable Errors）：使用Result<T, E>枚举

use std::error::Error;
use std::fmt;
use std::fs::File;
use std::io::{self, Read, Write};
use std::num::ParseIntError;
use std::string::FromUtf8Error;

// 1. 不可恢复错误：panic!
fn panic_example() {
    println!("=== 不可恢复错误：panic! ===");
    
    // panic!宏会导致程序崩溃并打印错误信息
    // panic!("这是一个不可恢复的错误！");
    
    // 使用panic!的调试信息
    let v = vec![1, 2, 3];
    // v[100]; // 索引越界，会自动panic
    
    println!("panic!演示完成");
}

// 2. 可恢复错误：Result<T, E>枚举
// Result<T, E>定义：enum Result<T, E> { Ok(T), Err(E) }

fn result_example() -> Result<(), std::io::Error> {
    println!("\n=== 可恢复错误：Result<T, E> ===");
    
    // 尝试打开文件
    let f = File::open("hello.txt");
    
    let mut f = match f {
        Ok(file) => file,
        Err(error) => {
            println!("打开文件失败：{}", error);
            return Err(error);
        }
    };
    
    let mut content = String::new();
    match f.read_to_string(&mut content) {
        Ok(_) => println!("文件内容：{}", content),
        Err(error) => {
            println!("读取文件失败：{}", error);
            return Err(error);
        }
    }
    
    Ok(())
}

// 3. 错误传播：?运算符
// ?运算符用于简化错误处理，自动将错误返回给调用者

fn read_username_from_file() -> Result<String, std::io::Error> {
    let mut f = File::open("hello.txt")?; // 如果失败，自动返回错误
    let mut s = String::new();
    f.read_to_string(&mut s)?; // 如果失败，自动返回错误
    Ok(s)
}

// 进一步简化
fn read_username_simple() -> Result<String, std::io::Error> {
    let mut s = String::new();
    File::open("hello.txt")?.read_to_string(&mut s)?;
    Ok(s)
}

// 更进一步简化（使用fs::read_to_string）
fn read_username_simplest() -> Result<String, std::io::Error> {
    std::fs::read_to_string("hello.txt")
}

// 4. ?运算符与From trait
// ?运算符会自动调用From trait将一种错误类型转换为另一种错误类型
// 这使得函数可以返回不同类型的错误

fn parse_number(s: &str) -> Result<i32, ParseIntError> {
    s.parse::<i32>()
}

// 5. 自定义错误类型

// 使用结构体定义自定义错误
#[derive(Debug)]
struct MyError {
    message: String,
    source: Option<Box<dyn Error + Send + Sync>>,
}

// 为MyError实现Display trait
impl fmt::Display for MyError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{}", self.message)
    }
}

// 为MyError实现Error trait
impl Error for MyError {
    fn source(&self) -> Option<&(dyn Error + 'static)> {
        self.source.as_ref().map(|e| e.as_ref())
    }
}

// 实现From trait，用于错误转换
impl From<ParseIntError> for MyError {
    fn from(error: ParseIntError) -> Self {
        Self {
            message: "解析整数失败".to_string(),
            source: Some(Box::new(error)),
        }
    }
}

impl From<std::io::Error> for MyError {
    fn from(error: std::io::Error) -> Self {
        Self {
            message: "IO操作失败".to_string(),
            source: Some(Box::new(error)),
        }
    }
}

// 使用自定义错误类型的函数
fn process_data() -> Result<(), MyError> {
    let _number = parse_number("not_a_number")?; // 自动转换为MyError
    let _content = std::fs::read_to_string("nonexistent.txt")?; // 自动转换为MyError
    Ok(())
}

// 6. 使用thiserror crate简化自定义错误（示例，实际需要添加依赖）
// 这是一个非常流行的错误处理库，简化了自定义错误的定义
/*
use thiserror::Error;

#[derive(Error, Debug)]
enum ThisErrorExample {
    #[error("解析错误: {0}")]
    ParseError(#[from] ParseIntError),
    
    #[error("IO错误: {0}")]
    IoError(#[from] std::io::Error),
    
    #[error("自定义错误: {0}")]
    CustomError(String),
}
*/

// 7. 使用anyhow crate进行错误处理（示例，实际需要添加依赖）
// anyhow提供了一种更灵活的错误处理方式，适合应用程序
/*
use anyhow::Result;

fn anyhow_example() -> Result<()> {
    let _number = parse_number("not_a_number")?;
    let _content = std::fs::read_to_string("nonexistent.txt")?;
    Ok(())
}
*/

// 8. 错误处理的组合

// 组合多个可能失败的操作
fn complex_operation() -> Result<i32, Box<dyn Error>> {
    // 读取文件内容
    let content = std::fs::read_to_string("number.txt")?;
    
    // 解析为整数
    let number: i32 = content.trim().parse()?;
    
    // 执行操作
    if number < 0 {
        return Err(Box::new(MyError {
            message: "数字不能为负数".to_string(),
            source: None,
        }));
    }
    
    Ok(number * 2)
}

// 9. 传播多个不同类型的错误

// 使用Box<dyn Error>作为返回类型，可以返回任何实现了Error trait的错误
fn multi_error_operation() -> Result<(), Box<dyn Error>> {
    // 尝试解析整数
    let number = "42".parse::<i32>()?;
    println!("解析的数字: {}", number);
    
    // 尝试打开文件
    let mut file = File::create("output.txt")?;
    
    // 尝试写入文件
    write!(file, "数字: {}", number)?;
    println!("文件写入成功");
    
    // 尝试读取不存在的文件
    // let _content = std::fs::read_to_string("nonexistent.txt")?;
    
    Ok(())
}

// 10. 错误处理的最佳实践

// 尽早返回错误
fn best_practice_early_return() -> Result<(), Box<dyn Error>> {
    let file = match File::open("data.txt") {
        Ok(f) => f,
        Err(e) => return Err(Box::new(e)),
    };
    
    // 继续处理文件
    let mut content = String::new();
    file.read_to_string(&mut content)?;
    
    Ok(())
}

// 提供有意义的错误信息
fn best_practice_meaningful_error() -> Result<(), MyError> {
    let file = File::open("config.toml")
        .map_err(|e| MyError {
            message: "无法打开配置文件config.toml".to_string(),
            source: Some(Box::new(e)),
        })?;
    
    // 处理文件
    Ok(())
}

// 11. 错误链

// Error trait的source()方法可以获取原始错误，形成错误链
fn error_chain_example() -> Result<(), MyError> {
    // 此函数会失败并返回包含错误链的MyError
    process_data()
}

// 12. 处理Option<T>和Result<T, E>的组合

fn option_result_combination() -> Result<Option<i32>, ParseIntError> {
    let numbers = vec!["1", "2", "three", "4"];
    
    for num_str in numbers {
        match num_str.parse::<i32>() {
            Ok(num) => return Ok(Some(num)),
            Err(e) => {
                println!("解析失败: {}, 继续尝试", e);
                continue;
            }
        }
    }
    
    Ok(None)
}

// 13. 使用unwrap()和expect()简化错误处理（仅适用于确定不会失败的情况）

fn unwrap_expect_example() {
    println!("\n=== 使用unwrap()和expect() ===");
    
    // unwrap()：如果Result是Ok则返回值，否则panic
    let f = File::open("hello.txt").unwrap();
    
    // expect()：类似unwrap()，但可以自定义panic信息
    let f = File::open("hello.txt").expect("无法打开hello.txt文件");
    
    println!("unwrap()和expect()演示完成");
}

// 14. 自定义Result类型

type MyResult<T> = Result<T, Box<dyn Error>>;

trait MyResultExt<T> {
    fn custom_error(self, message: &str) -> MyResult<T>;
}

impl<T, E: Error + Send + Sync + 'static> MyResultExt<T> for Result<T, E> {
    fn custom_error(self, message: &str) -> MyResult<T> {
        self.map_err(|e| {
            let custom_error = MyError {
                message: message.to_string(),
                source: Some(Box::new(e)),
            };
            Box::new(custom_error) as Box<dyn Error>
        })
    }
}

fn custom_result_type_example() -> MyResult<()> {
    let content = std::fs::read_to_string("data.txt")
        .custom_error("读取数据文件失败")?;
    
    let number: i32 = content.trim().parse()
        .custom_error("解析数据文件中的数字失败")?;
    
    println!("读取并解析的数字: {}", number);
    Ok(())
}

// 15. 错误处理与测试

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_parse_number() {
        assert_eq!(parse_number("42"), Ok(42));
        assert!(parse_number("not_a_number").is_err());
    }
    
    #[test]
    #[should_panic(expected = "索引越界")]
    fn test_panic() {
        let v = vec![1, 2, 3];
        v[100];
    }
}

pub fn run_example() {
    println!("=== Rust学习示例 ===\n");
    // 运行各个示例
    panic_example();
    
    if let Err(e) = result_example() {
        println!("result_example失败: {}", e);
    }
    
    match read_username_from_file() {
        Ok(username) => println!("\n读取的用户名: {}", username),
        Err(e) => println!("\n读取用户名失败: {}", e),
    }
    
    match process_data() {
        Ok(_) => println!("\nprocess_data成功"),
        Err(e) => {
            println!("\nprocess_data失败: {}", e);
            // 打印错误链
            if let Some(source) = e.source() {
                println!("  原始错误: {}", source);
            }
        }
    }
    
    match complex_operation() {
        Ok(result) => println!("\ncomplex_operation结果: {}", result),
        Err(e) => println!("\ncomplex_operation失败: {}", e),
    }
    
    match multi_error_operation() {
        Ok(_) => println!("\nmulti_error_operation成功"),
        Err(e) => println!("\nmulti_error_operation失败: {}", e),
    }
    
    match error_chain_example() {
        Ok(_) => println!("\nerror_chain_example成功"),
        Err(e) => {
            println!("\nerror_chain_example失败: {}", e);
            // 遍历错误链
            let mut current = Some(e.as_ref());
            let mut index = 0;
            while let Some(err) = current {
                println!("  错误{}: {}", index, err);
                current = err.source();
                index += 1;
            }
        }
    }
    
    match option_result_combination() {
        Ok(Some(num)) => println!("\noption_result_combination找到数字: {}", num),
        Ok(None) => println!("\noption_result_combination没有找到有效数字"),
        Err(e) => println!("\noption_result_combination失败: {}", e),
    }
    
    // unwrap_expect_example(); // 取消注释查看效果
    
    match custom_result_type_example() {
        Ok(_) => println!("\ncustom_result_type_example成功"),
        Err(e) => println!("\ncustom_result_type_example失败: {}", e),
    }
    
    // 16. 错误处理总结
    println!("\n=== 错误处理总结 ===");
    println!("1. 不可恢复错误：使用panic!宏，导致程序崩溃");
    println!("2. 可恢复错误：使用Result<T, E>枚举，允许优雅处理");
    println!("3. 错误传播：使用?运算符，简化错误传递");
    println!("4. 自定义错误：实现Error trait，提供有意义的错误信息");
    println!("5. 错误链：使用source()方法，追踪原始错误");
    println!("6. 最佳实践：");
    println!("   - 优先使用Result处理可恢复错误");
    println!("   - 仅在确实不可恢复的情况下使用panic!");
    println!("   - 提供清晰、有意义的错误信息");
    println!("   - 考虑使用错误处理库（thiserror, anyhow）简化开发");
    println!("   - 避免过度使用unwrap()和expect()");
}

// 17. 更多错误处理场景

// 处理多个错误类型的转换
fn handle_multiple_errors() -> Result<(), Box<dyn Error>> {
    // 模拟一个可能返回不同错误的操作
    let result1: Result<i32, ParseIntError> = "123".parse();
    let result2: Result<String, FromUtf8Error> = String::from_utf8(vec![255]);
    
    let num = result1?;
    println!("解析的数字: {}", num);
    
    let _text = result2?;
    Ok(())
}

// 异步错误处理（示例，实际需要使用async/await）
/*
async fn async_error_example() -> Result<(), Box<dyn Error + Send + Sync>> {
    // 异步操作中的错误处理
    let content = tokio::fs::read_to_string("file.txt").await?;
    let number: i32 = content.trim().parse()?;
    Ok(())
}
*/

// 错误恢复策略
fn error_recovery() -> Result<i32, Box<dyn Error>> {
    // 尝试从文件读取，失败则使用默认值
    let content = match std::fs::read_to_string("config.txt") {
        Ok(c) => c,
        Err(_) => {
            println!("使用默认配置");
            "42".to_string()
        }
    };
    
    let number: i32 = content.trim().parse()?;
    Ok(number)
}
// 用于单独运行本文件的main函数
fn main() {
    run_example();
}
