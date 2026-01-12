// 01_variables.rs - Rust变量系统详解

pub fn run_example() {
    println!("=== Rust变量系统详解 ===\n");
    // 1. 变量声明与不可变性（默认）
    // 在Rust中，变量默认是不可变的（immutable）
    let x = 5; // 声明一个不可变变量x，值为5
    println!("x的值是: {}", x);
    
    // 尝试修改不可变变量会导致编译错误
    // x = 6; // 取消注释这行会报错：cannot assign twice to immutable variable
    
    // 2. 可变变量（使用mut关键字）
    // 使用mut关键字可以使变量可变
    let mut y = 5; // 声明一个可变变量y，值为5
    println!("y的初始值是: {}", y);
    
    y = 6; // 可以修改可变变量
    println!("y修改后的值是: {}", y);
    
    // 3. 变量遮蔽（Shadowing）
    // 可以使用相同的名称声明新变量，新变量会遮蔽旧变量
    let z = 5;
    println!("z的初始值是: {}", z);
    
    let z = z + 1; // 遮蔽旧变量z，新z的值是6
    println!("z遮蔽后的值是: {}", z);
    
    let z = z * 2; // 再次遮蔽，新z的值是12
    println!("z再次遮蔽后的值是: {}", z);
    
    // 注意：变量遮蔽与mut的区别
    // - mut允许在同一变量上修改值
    // - 遮蔽创建了一个新变量，可以改变类型而保持名称不变
    let spaces = "   ";
    println!("spaces的类型是字符串，内容是: '{spaces}'");
    
    let spaces = spaces.len(); // 遮蔽成一个数字类型
    println!("遮蔽后spaces的类型是usize，值是: {spaces}");
    
    // 4. 常量（Constants）
    // 使用const关键字声明常量，必须注明类型
    const MAX_POINTS: u32 = 100_000;
    println!("常量MAX_POINTS的值是: {}", MAX_POINTS);
    
    // 常量与不可变变量的区别：
    // - 常量必须使用const关键字，而不是let
    // - 常量必须标注类型
    // - 常量可以在任何作用域内声明，包括全局作用域
    // - 常量只能被常量表达式初始化，不能是函数调用或其他在运行时计算的值
    
    // 5. 静态变量（Static Variables）
    // 静态变量在程序的整个生命周期内都存在
    static GLOBAL_COUNT: i32 = 0;
    println!("静态变量GLOBAL_COUNT的值是: {}", GLOBAL_COUNT);
    
    // 注意：修改静态变量需要使用unsafe块
    // unsafe {
    //     GLOBAL_COUNT += 1; // 这会导致未定义行为
    // }
    
    // 6. 类型注解
    // 当Rust无法推断类型时，需要显式标注类型
    let explicit_type: i32 = 42; // 显式标注为i32类型
    println!("显式类型标注的变量值是: {}", explicit_type);
    
    // 没有初始值时必须标注类型
    let mut empty_variable: String;
    empty_variable = String::from("Hello");
    println!("先声明后初始化的变量值是: {}", empty_variable);
    
    // 调用变量作用域示例函数
    variable_scope();
    
    println!("\n=== 示例结束 ===");
}

// 7. 变量作用域示例
fn variable_scope() {
    // 作用域是变量有效的范围
    let outer = 10;
    
    { // 进入新的作用域
        let inner = 20;
        println!("内部作用域：outer = {}, inner = {}", outer, inner);
        
        let outer = 30; // 遮蔽外部作用域的outer变量
        println!("内部作用域遮蔽后：outer = {}, inner = {}", outer, inner);
    } // 离开内部作用域，inner变量被销毁
    
    println!("外部作用域：outer = {}", outer);
    // println!("尝试访问inner: {}", inner); // 这会报错，因为inner不在作用域内
}

// 用于单独运行本文件的main函数
fn main() {
    run_example();
}