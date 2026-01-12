// 07_enums.rs - Rust枚举与模式匹配详解

// 1. 枚举的基本定义
// 使用enum关键字定义枚举
enum Direction {
    Up,
    Down,
    Left,
    Right,
}

// 2. 带数据的枚举变体
enum Message {
    Quit,                   // 无数据
    Move { x: i32, y: i32 }, // 包含结构体
    Write(String),          // 包含单个String
    ChangeColor(i32, i32, i32), // 包含三个i32值
}

// 3. 枚举的方法定义（使用impl块）
impl Message {
    fn call(&self) {
        // 根据枚举变体执行不同的操作
        match self {
            Message::Quit => println!("退出消息"),
            Message::Move { x, y } => println!("移动到 x: {}, y: {}", x, y),
            Message::Write(text) => println!("写入消息: {}", text),
            Message::ChangeColor(r, g, b) => println!("改变颜色为 R: {}, G: {}, B: {}", r, g, b),
        }
    }
}

// 4. Option枚举（Rust标准库提供）
// Option<T>表示一个值可能存在（Some(T)）或不存在（None）
// 定义类似于：enum Option<T> { Some(T), None }

// 5. Result枚举（Rust标准库提供）
// Result<T, E>表示操作可能成功（Ok(T)）或失败（Err(E)）
// 定义类似于：enum Result<T, E> { Ok(T), Err(E) }

// 6. 自定义Result类型
type MyResult<T> = Result<T, MyError>;

enum MyError {
    NotFound,
    PermissionDenied,
    InvalidInput(String),
}

// 7. 嵌套枚举
enum OuterEnum {
    Variant1,
    Variant2(InnerEnum),
}

enum InnerEnum {
    Inner1,
    Inner2(i32),
}

pub fn run_example() {
    println!("=== Rust学习示例 ===\n");
    // 8. 枚举变体的使用
    let dir = Direction::Up;
    println!("方向: {:?}", dir); // 需要Debug trait
    
    // 9. 带数据的枚举变体实例化
    let msg1 = Message::Quit;
    let msg2 = Message::Move { x: 10, y: 20 };
    let msg3 = Message::Write(String::from("Hello, Rust!"));
    let msg4 = Message::ChangeColor(255, 0, 0);
    
    // 调用枚举方法
    msg1.call();
    msg2.call();
    msg3.call();
    msg4.call();
    
    // 10. 模式匹配（match表达式）
    let number = Some(7);
    
    let result = match number {
        Some(5) => "五",
        Some(n) if n % 2 == 0 => "偶数",
        Some(n) => format!("奇数: {}", n),
        None => "没有数字",
    };
    
    println!("模式匹配结果: {}", result);
    
    // 11. Option枚举的使用示例
    let some_number = Some(5);
    let absent_number: Option<i32> = None;
    
    // Option<T>和T是不同的类型，不能直接运算
    // let sum = some_number + 1; // 这会报错
    
    // 使用match处理Option
    match some_number {
        Some(n) => println!("some_number的值是: {}", n),
        None => println!("some_number没有值"),
    }
    
    // 12. Result枚举的使用示例
    let success_result: Result<i32, &str> = Ok(10);
    let error_result: Result<i32, &str> = Err("发生错误");
    
    match success_result {
        Ok(value) => println!("成功: {}", value),
        Err(e) => println!("失败: {}", e),
    }
    
    match error_result {
        Ok(value) => println!("成功: {}", value),
        Err(e) => println!("失败: {}", e),
    }
    
    // 13. if let表达式（简化模式匹配）
    let some_u8_value = Some(0u8);
    
    // 使用match
    match some_u8_value {
        Some(3) => println!("数字3"),
        _ => (),
    }
    
    // 使用if let简化
    if let Some(3) = some_u8_value {
        println!("数字3");
    }
    
    // if let配合else
    if let Some(n) = some_u8_value {
        println!("获取到数字: {}", n);
    } else {
        println!("没有数字");
    }
    
    // 14. while let表达式
    let mut stack = Vec::new();
    stack.push(1);
    stack.push(2);
    stack.push(3);
    
    println!("\n使用while let遍历栈:");
    while let Some(top) = stack.pop() {
        println!("弹出: {}", top);
    }
    
    // 15. 模式匹配中的高级模式
    // 解构结构体
    struct Point { x: i32, y: i32 }
    let p = Point { x: 10, y: 20 };
    
    match p {
        Point { x, y: 0 } => println!("在x轴上，x = {}", x),
        Point { x: 0, y } => println!("在y轴上，y = {}", y),
        Point { x, y } => println!("在点({}, {})", x, y) 
    }
    
    // 解构枚举中的结构体
    let move_msg = Message::Move { x: 5, y: 10 };
    
    if let Message::Move { x, y } = move_msg {
        println!("移动消息: x = {}, y = {}", x, y);
    }
    
    // 解构元组
    let tuple = (1, 2, (3, 4));
    
    match tuple {
        (a, b, (c, d)) if a + b > 2 => println!("a + b > 2: {}, {}, {}, {}", a, b, c, d),
        _ => println!("其他情况"),
    }
    
    // 16. 枚举的实际应用示例：计算表达式
enum Expr {
    Literal(i32),
    Add(Box<Expr>, Box<Expr>),
    Subtract(Box<Expr>, Box<Expr>),
    Multiply(Box<Expr>, Box<Expr>),
    Divide(Box<Expr>, Box<Expr>),
}

impl Expr {
    fn evaluate(&self) -> i32 {
        match self {
            Expr::Literal(n) => *n,
            Expr::Add(left, right) => left.evaluate() + right.evaluate(),
            Expr::Subtract(left, right) => left.evaluate() - right.evaluate(),
            Expr::Multiply(left, right) => left.evaluate() * right.evaluate(),
            Expr::Divide(left, right) => left.evaluate() / right.evaluate(),
        }
    }
}
    
    // 创建表达式：1 + 2 * 3
    let expr = Expr::Add(
        Box::new(Expr::Literal(1)),
        Box::new(Expr::Multiply(
            Box::new(Expr::Literal(2)),
            Box::new(Expr::Literal(3)),
        )),
    );
    
    println!("\n表达式计算结果: {}", expr.evaluate());
    
    // 17. 自定义Result类型的使用
    fn divide(a: i32, b: i32) -> MyResult<i32> {
        if b == 0 {
            Err(MyError::InvalidInput(String::from("除数不能为0")))
        } else {
            Ok(a / b)
        }
    }
    
    match divide(10, 2) {
        Ok(result) => println!("10 / 2 = {}", result),
        Err(e) => match e {
            MyError::NotFound => println!("未找到"),
            MyError::PermissionDenied => println!("权限不足"),
            MyError::InvalidInput(msg) => println!("无效输入: {}", msg),
        },
    }
    
    match divide(10, 0) {
        Ok(result) => println!("10 / 0 = {}", result),
        Err(e) => match e {
            MyError::NotFound => println!("未找到"),
            MyError::PermissionDenied => println!("权限不足"),
            MyError::InvalidInput(msg) => println!("无效输入: {}", msg),
        },
    }
    
    // 18. 使用_通配符和..省略符
    let colors = (255, 0, 130);
    
    match colors {
        (r, _, b) => println!("红色: {}, 蓝色: {}", r, b), // 忽略绿色
    }
    
    struct Person { name: String, age: u32, city: String }
    let person = Person { name: String::from("Alice"), age: 30, city: String::from("Beijing") };
    
    match person {
        Person { name, .. } => println!("姓名: {}", name), // 只关心name字段
    }
}

// 19. 为枚举实现Debug trait以便打印
#[derive(Debug)]
enum DebugEnum {
    Variant1,
    Variant2(i32, String),
    Variant3 { x: f64, y: f64 },
}

fn debug_example() {
    let enum1 = DebugEnum::Variant1;
    let enum2 = DebugEnum::Variant2(42, String::from("hello"));
    let enum3 = DebugEnum::Variant3 { x: 3.14, y: 2.71 };
    
    println!("\n调试示例:");
    println!("enum1: {:?}", enum1);
    println!("enum2: {:?}", enum2);
    println!("enum3: {:?}", enum3);
}
// 用于单独运行本文件的main函数
fn main() {
    run_example();
}
