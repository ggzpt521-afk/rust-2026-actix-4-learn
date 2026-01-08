// ========== 0. 引入标准库 ==========
use std::num::ParseIntError;   // 标准库提供的“字符串转整数失败”错误类型

// ========== 1. 自定义错误枚举（最小可运行版） ==========
#[derive(Debug)]                // Debug → 可以 {:?} 打印
pub enum MyError {
    NotFound,                   // ① 资源不存在
    BadInput(String),           // ② 输入无效，带描述
}

// ========== 2. 主函数（测试用） ==========
fn main() {
    // 2.1 正常输入
    let result = square("32");
    println!("Hello, world!");
    println!("test == {:?}", result);        // Ok(1024)

    // 2.2 异常输入（非数字）
    let result2 = square("RT");
    println!("error == {:?}", result2);      // Err(ParseIntError { … })

    // 2.3 用 ? 运算符（正常输入）
    let result3 = squareDealErr("32");
    println!("error == {:?}", result3);      // Ok(1024)

    // 2.4 用 ? 运算符（异常输入）
    let result5 = squareDealErr("RT");
    println!("error == {:?}", result5);      // Err(ParseIntError { … })
}

// ========== 3. 手工 match 版（显式处理错误） ==========
fn square(val: &str) -> Result<i32, ParseIntError> {
    // 3.1 尝试解析字符串 → Result<i32, ParseIntError>
    match val.parse::<i32>() {
        Ok(num) => Ok(num.pow(2)),        // 成功 → 返回平方
        Err(e) => Err(e),                 // 失败 → 返回错误
    }
}

// ========== 4. ? 运算符版（隐式处理错误） ==========
fn squareDealErr(val: &str) -> Result<i32, ParseIntError> {
    // 4.1 ? 运算符：如果 parse 成功 → 返回 i32；如果失败 → 提前返回 Err(e)
    let num = val.parse::<i32>()?;       // **? = 自动解包 + 提前返回**
    Ok(num)                              // 成功路径
}