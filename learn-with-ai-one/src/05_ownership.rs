// 05_ownership.rs - Rust所有权系统详解

pub fn run_example() {
    println!("=== Rust学习示例 ===\n");
    // 1. 所有权的三条规则
    // - Rust中的每个值都有一个所有者变量
    // - 同一时刻只能有一个所有者
    // - 当所有者离开作用域时，值会被自动销毁
    
    // 2. 变量作用域与String类型
    // 基本数据类型（如i32, bool等）是在栈上存储的，而String类型是在堆上存储的
    {
        let s = String::from("hello"); // s进入作用域
        println!("s的值: {}", s);
    } // s离开作用域，Rust自动调用drop函数释放堆内存
    
    // 3. 移动（Move）语义
    // Rust中，当一个非基本类型的变量被赋值给另一个变量时，所有权会发生转移
    let s1 = String::from("hello");
    let s2 = s1; // s1的所有权转移给s2，s1不再有效
    
    println!("s2的值: {}", s2);
    // println!("s1的值: {}", s1); // 这会报错，因为s1已经不再拥有该值的所有权
    
    // 基本数据类型的情况（Copy trait）
    // 对于实现了Copy trait的类型，赋值时会进行复制而不是移动
    let x = 5;
    let y = x;
    
    println!("x的值: {}, y的值: {}", x, y); // 两者都有效
    
    // 4. 克隆（Clone）操作
    // 如果需要深拷贝堆上的数据，可以使用clone方法
    let s1 = String::from("hello");
    let s2 = s1.clone();
    
    println!("s1的值: {}, s2的值: {}", s1, s2); // 两者都有效
    
    // 5. 不可变借用（Borrowing）
    // 可以通过引用（&）来借用值的所有权，而不是获取所有权
    let s1 = String::from("hello");
    let len = calculate_length(&s1); // 传递s1的引用
    
    println!("'{}'的长度是: {}", s1, len); // s1仍然有效
    
    // 6. 可变借用
    // 使用&mut来创建可变引用，可以修改引用指向的值
    let mut s = String::from("hello");
    change_string(&mut s); // 传递可变引用
    
    println!("修改后的字符串: {}", s);
    
    // 注意：可变借用的规则
    // - 在同一作用域内，一个值只能有一个可变引用
    // - 可变引用和不可变引用不能同时存在
    
    // 7. 借用规则示例
    let mut s = String::from("hello");
    
    // 同一作用域内只能有一个可变引用
    // let r1 = &mut s;
    // let r2 = &mut s; // 这会报错
    
    // 可变引用和不可变引用不能同时存在
    // let r1 = &s;
    // let r2 = &mut s; // 这会报错
    
    // 8. 切片（Slices）类型
    // 切片是对集合中一段连续元素的引用，没有所有权
    
    // 字符串切片
    let s = String::from("hello world");
    let hello = &s[0..5]; // 从索引0到5（不包含）的切片
    let world = &s[6..11]; // 从索引6到11的切片
    
    println!("字符串切片: '{}' 和 '{}'", hello, world);
    
    // 简化的切片语法
    let slice1 = &s[..5]; // 等同于 &s[0..5]
    let slice2 = &s[6..]; // 等同于 &s[6..s.len()]
    let slice3 = &s[..]; // 等同于 &s[0..s.len()]
    
    println!("简化切片语法: '{}', '{}', '{}'", slice1, slice2, slice3);
    
    // 9. 数组切片
    let a = [1, 2, 3, 4, 5];
    let slice = &a[1..3]; // 数组切片
    
    println!("数组切片: {:?}", slice);
    
    // 10. 字符串字面量是切片
    let s = "Hello, world!"; // s的类型是 &str，它是一个指向二进制程序特定位置的切片
    
    // 11. 悬垂引用（Dangling References）
    // Rust编译器会防止悬垂引用的产生
    // 下面的代码会编译失败，因为返回的引用指向了一个即将销毁的变量
    // let r = dangle();
    
    // 12. 所有权与函数参数
    let s = String::from("hello");
    
    // 传递所有权
    take_ownership(s);
    // println!("s: {}", s); // 这会报错，因为s的所有权已被转移
    
    // 传递不可变引用
    let s = String::from("hello");
    borrow_immutable(&s);
    println!("s仍然有效: {}", s);
    
    // 传递可变引用
    let mut s = String::from("hello");
    borrow_mutable(&mut s);
    println!("修改后的s: {}", s);
    
    // 13. 所有权与函数返回值
    let s1 = give_ownership(); // 获取返回值的所有权
    println!("从函数获取的字符串: {}", s1);
    
    let s2 = String::from("hello");
    let s3 = take_and_give_back(s2); // s2的所有权转移到函数，然后返回给s3
    println!("s3的值: {}", s3);
    // println!("s2的值: {}", s2); // 这会报错，因为s2的所有权已被转移
}

// 计算字符串长度（不可变借用）
fn calculate_length(s: &String) -> usize {
    s.len() // 返回字符串长度，不获取所有权
}

// 修改字符串（可变借用）
fn change_string(s: &mut String) {
    s.push_str(", world"); // 修改借用的字符串
}

// 获取字符串所有权
fn take_ownership(s: String) {
    println!("获取到的字符串: {}", s);
    // s离开作用域，被自动销毁
}

// 不可变借用
fn borrow_immutable(s: &String) {
    println!("不可变借用的字符串: {}", s);
    // 引用离开作用域，不影响原字符串
}

// 可变借用
fn borrow_mutable(s: &mut String) {
    s.push_str(", Rust");
    println!("可变借用并修改的字符串: {}", s);
}

// 返回一个String，调用者获得所有权
fn give_ownership() -> String {
    let s = String::from("hello from give_ownership");
    s // 返回s，所有权转移给调用者
}

// 获取一个String，然后返回它
fn take_and_give_back(s: String) -> String {
    println!("获取到并即将返回的字符串: {}", s);
    s // 返回s，所有权转移给调用者
}

// 尝试返回悬垂引用（编译失败）
// fn dangle() -> &String {
//     let s = String::from("hello"); // s进入作用域
//     &s // 返回s的引用
// } // s离开作用域，被销毁，返回的引用指向了无效内存

// 14. 切片作为函数参数
// 更好的方式是接受&str类型，这样既可以接受字符串切片，也可以接受String的引用
fn first_word(s: &str) -> &str {
    let bytes = s.as_bytes();
    
    for (i, &item) in bytes.iter().enumerate() {
        if item == b' ' {
            return &s[0..i];
        }
    }
    
    &s[..] // 返回整个字符串切片
}

// 切片函数示例
fn slice_example() {
    let my_string = String::from("hello world");
    let word = first_word(&my_string[..]); // 传递String的切片
    
    let my_string_literal = "hello world";
    let word2 = first_word(my_string_literal); // 传递字符串字面量（本身就是切片）
    
    println!("第一个单词: '{}' 和 '{}'", word, word2);
}
// 用于单独运行本文件的main函数
fn main() {
    run_example();
}
