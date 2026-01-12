// 03_functions.rs - Rust函数详解

// 1. 函数定义的基本语法
// 使用fn关键字定义函数，函数名使用snake_case命名规范
fn greet() {
    println!("Hello, Rust!");
}

// 2. 带参数的函数
// 参数需要指定类型，多个参数用逗号分隔
fn add(a: i32, b: i32) {
    println!("{} + {} = {}", a, b, a + b);
}

// 3. 带返回值的函数
// 返回值类型使用->指定
fn multiply(a: i32, b: i32) -> i32 {
    a * b // 表达式作为返回值，不需要分号
}

// 4. 提前返回（使用return关键字）
fn divide(a: i32, b: i32) -> Option<f64> {
    if b == 0 {
        return None; // 提前返回，需要分号
    }
    Some(a as f64 / b as f64) // 正常返回
}

// 5. 无返回值的函数（隐式返回()）
fn print_result(result: i32) -> () {
    println!("结果是: {}", result);
    // 隐式返回()，可以省略
}

// 6. 函数参数的所有权示例
fn take_ownership(s: String) {
    println!("获取到的字符串是: {}", s);
    // s在这里被销毁
}

fn borrow_reference(s: &String) {
    println!("借用的字符串是: {}", s);
    // 引用离开作用域，不影响原字符串
}

// 7. 可变引用参数
fn modify_string(s: &mut String) {
    s.push_str(", World!");
}

// 8. 嵌套函数
fn outer_function(x: i32) {
    println!("外部函数，x = {}", x);
    
    // 在函数内部定义嵌套函数
    fn inner_function(y: i32) {
        println!("内部函数，y = {}", y);
    }
    
    inner_function(x * 2);
}

// 9. 高阶函数（函数作为参数）
// 定义一个接受函数作为参数的函数
fn apply_function<F>(x: i32, f: F) -> i32 
where F: Fn(i32) -> i32 {
    f(x)
}

// 用于测试的函数
fn double(x: i32) -> i32 {
    x * 2
}

fn square(x: i32) -> i32 {
    x * x
}

// 10. 闭包（Closures）
// 闭包是可以捕获其环境变量的匿名函数
fn closure_example() {
    let factor = 3;
    
    // 闭包定义，使用||代替参数列表
    let triple = |x| x * factor;
    
    println!("闭包示例：3 * {} = {}", 5, triple(5));
    
    // 带类型注解的闭包
    let add_one: fn(i32) -> i32 = |x: i32| -> i32 { x + 1 };
    println!("带类型注解的闭包：{} + 1 = {}", 10, add_one(10));
    
    // 多参数闭包
    let sum = |x, y| x + y;
    println!("多参数闭包：{} + {} = {}", 3, 4, sum(3, 4));
}

// 11. 递归函数
// 递归函数必须指定返回类型
fn factorial(n: u32) -> u32 {
    if n == 0 || n == 1 {
        1
    } else {
        n * factorial(n - 1)
    }
}

// 12. 发散函数（Diverging Functions）
// 发散函数永远不会返回，使用!作为返回类型
fn diverging_function() -> ! {
    panic!("这个函数永远不会返回");
}

pub fn run_example() {
    println!("=== Rust学习示例 ===\n");
    // 调用基本函数
    greet();
    
    // 调用带参数的函数
    add(3, 5);
    
    // 调用带返回值的函数
    let result = multiply(4, 6);
    println!("4 * 6 = {}", result);
    
    // 调用带条件返回的函数
    match divide(10, 2) {
        Some(value) => println!("10 / 2 = {}", value),
        None => println!("除数不能为0")
    }
    
    // 调用无返回值的函数
    print_result(result);
    
    // 演示所有权传递
    let s = String::from("Hello");
    take_ownership(s);
    // println!("s: {}", s); // 这会报错，因为s的所有权已经被转移
    
    let s2 = String::from("Hello");
    borrow_reference(&s2);
    println!("原字符串s2: {}", s2); // 可以正常访问
    
    // 演示可变引用
    let mut s3 = String::from("Hello");
    modify_string(&mut s3);
    println!("修改后的字符串s3: {}", s3);
    
    // 调用嵌套函数
    outer_function(5);
    
    // 调用高阶函数
    let double_result = apply_function(10, double);
    let square_result = apply_function(10, square);
    println!("高阶函数示例：");
    println!("double(10) = {}", double_result);
    println!("square(10) = {}", square_result);
    
    // 调用闭包示例
    closure_example();
    
    // 调用递归函数
    let n = 5;
    println!("{}的阶乘是：{}", n, factorial(n));
    
    // 调用发散函数（会导致程序崩溃，演示用）
    // diverging_function();
}

// 13. 函数指针
fn function_pointer_example() {
    // 定义函数指针类型
    type Operation = fn(i32, i32) -> i32;
    
    // 将函数赋值给函数指针变量
    let add_func: Operation = |a, b| a + b;
    let subtract_func: Operation = |a, b| a - b;
    
    println!("函数指针示例：");
    println!("add_func(10, 5) = {}", add_func(10, 5));
    println!("subtract_func(10, 5) = {}", subtract_func(10, 5));
}
// 用于单独运行本文件的main函数
fn main() {
    run_example();
}
