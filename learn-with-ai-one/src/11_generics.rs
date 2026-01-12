// 11_generics.rs - Rust泛型编程详解

// 泛型（Generics）是Rust中用于编写通用代码的特性，允许在不指定具体类型的情况下编写函数、结构体、枚举等
// 泛型的主要优点：
// 1. 代码复用：可以为多种类型编写相同的逻辑
// 2. 类型安全：编译时进行类型检查
// 3. 零成本抽象：运行时不会产生额外开销

// 1. 泛型函数

// 定义一个泛型函数，用于比较两个值是否相等
fn is_equal<T: PartialEq>(a: T, b: T) -> bool {
    a == b
}

// 泛型函数的使用
fn generic_function_example() {
    println!("=== 泛型函数 ===");
    
    // 比较整数
    println!("1 == 2: {}", is_equal(1, 2));
    println!("5 == 5: {}", is_equal(5, 5));
    
    // 比较字符串
    println!("\"hello\" == \"world\": {}", is_equal("hello", "world"));
    println!("\"rust\" == \"rust\": {}", is_equal("rust", "rust"));
    
    // 比较浮点数
    println!("3.14 == 2.71: {}", is_equal(3.14, 2.71));
    
    // 比较自定义类型（需要实现PartialEq trait）
    #[derive(PartialEq)]
    struct Point { x: i32, y: i32 }
    
    let p1 = Point { x: 1, y: 2 };
    let p2 = Point { x: 1, y: 2 };
    let p3 = Point { x: 3, y: 4 };
    
    println!("p1 == p2: {}", is_equal(p1, p2));
    println!("p1 == p3: {}", is_equal(p2, p3));
}

// 2. 泛型结构体

// 定义一个泛型结构体，用于存储一对值
struct Pair<T, U> {
    first: T,
    second: U,
}

// 为泛型结构体实现方法
impl<T, U> Pair<T, U> {
    // 创建新的Pair实例
    fn new(first: T, second: U) -> Self {
        Self { first, second }
    }
    
    // 获取first字段的引用
    fn first(&self) -> &T {
        &self.first
    }
    
    // 获取second字段的引用
    fn second(&self) -> &U {
        &self.second
    }
}

// 为特定类型的泛型结构体实现方法
impl Pair<i32, f64> {
    // 计算两个值的和（仅适用于i32和f64的组合）
    fn sum(&self) -> f64 {
        self.first as f64 + self.second
    }
}

// 泛型结构体的使用
fn generic_struct_example() {
    println!("\n=== 泛型结构体 ===");
    
    // 创建不同类型的Pair实例
    let pair1 = Pair::new(1, 2.5);
    let pair2 = Pair::new("hello", true);
    let pair3 = Pair::new(Pair::new(1, 2), Pair::new(3, 4));
    
    println!("pair1: ({}, {})", pair1.first(), pair1.second());
    println!("pair2: ({}, {})", pair2.first(), pair2.second());
    
    // 使用特定类型的方法
    println!("pair1的和: {}", pair1.sum());
}

// 3. 泛型枚举

// 定义一个泛型枚举，用于表示可能存在或不存在的值（类似Option<T>）
enum MyOption<T> {
    Some(T),
    None,
}

// 为泛型枚举实现方法
impl<T> MyOption<T> {
    fn is_some(&self) -> bool {
        match self {
            MyOption::Some(_) => true,
            MyOption::None => false,
        }
    }
    
    fn unwrap(self) -> T {
        match self {
            MyOption::Some(value) => value,
            MyOption::None => panic!("尝试从None中unwrap值"),
        }
    }
}

// 泛型枚举的使用
fn generic_enum_example() {
    println!("\n=== 泛型枚举 ===");
    
    let some_int = MyOption::Some(42);
    let some_string = MyOption::Some(String::from("hello"));
    let none_value: MyOption<i32> = MyOption::None;
    
    println!("some_int是Some吗: {}", some_int.is_some());
    println!("none_value是Some吗: {}", none_value.is_some());
    println!("some_int的值: {}", some_int.unwrap());
    
    // 使用标准库的Option<T>
    let std_some = Some(3.14);
    let std_none: Option<String> = None;
    
    println!("std_some的值: {:?}", std_some);
    println!("std_none的值: {:?}", std_none);
}

// 4. 泛型约束（Generic Constraints）

// 使用trait bounds来约束泛型类型
// trait bounds指定泛型类型必须实现的trait

// 定义一个泛型函数，要求类型T实现Debug trait
use std::fmt::Debug;

fn print_debug<T: Debug>(value: T) {
    println!("{:?}", value);
}

// 多重trait bounds
fn print_and_clone<T: Debug + Clone>(value: T) {
    println!("值: {:?}", value);
    let cloned = value.clone();
    println!("克隆的值: {:?}", cloned);
}

// 使用where子句简化trait bounds
fn complex_generic<T, U>(a: T, b: U) -> i32
where
    T: Into<i32> + Debug,
    U: Into<i32> + Debug,
{
    println!("a: {:?}, b: {:?}", a, b);
    a.into() + b.into()
}

// 泛型约束的使用
fn generic_constraints_example() {
    println!("\n=== 泛型约束 ===");
    
    print_debug(42);
    print_debug("hello");
    print_debug(Point { x: 1, y: 2 });
    
    print_and_clone(42);
    print_and_clone(String::from("rust"));
    
    let sum = complex_generic(10, 20);
    println!("complex_generic(10, 20) = {}", sum);
    
    let sum = complex_generic(5, 3.5);
    println!("complex_generic(5, 3.5) = {}", sum);
}

// 5. 泛型与所有权

// 泛型函数可以接受不同所有权类型的参数

// 接受值（移动所有权）
fn take_ownership<T>(value: T) -> T {
    println!("接受了所有权");
    value // 返回所有权
}

// 接受不可变引用
fn borrow_immutable<T>(value: &T) where T: Debug {
    println!("不可变引用: {:?}", value);
}

// 接受可变引用
fn borrow_mutable<T>(value: &mut T) where T: Debug + Default {
    *value = T::default();
    println!("修改后的值: {:?}", value);
}

// 泛型与所有权的使用
fn generics_and_ownership() {
    println!("\n=== 泛型与所有权 ===");
    
    let s = String::from("hello");
    let s = take_ownership(s); // 移动所有权
    println!("s: {}", s);
    
    borrow_immutable(&s); // 不可变借用
    println!("s: {}", s);
    
    let mut vec = vec![1, 2, 3];
    borrow_mutable(&mut vec); // 可变借用
    println!("vec: {:?}", vec);
}

// 6. 泛型在标准库中的应用

// Vec<T>：动态数组
// HashMap<K, V>：哈希映射
// Option<T>：可选值
// Result<T, E>：结果

fn std_lib_generics_example() {
    println!("\n=== 标准库中的泛型 ===");
    
    // Vec<T>
    let mut vec = Vec::new();
    vec.push(1);
    vec.push(2);
    vec.push(3);
    println!("Vec<T>: {:?}", vec);
    
    // HashMap<K, V>
    use std::collections::HashMap;
    
    let mut map = HashMap::new();
    map.insert("apple", 1);
    map.insert("banana", 2);
    map.insert("cherry", 3);
    println!("HashMap<K, V>: {:?}", map);
    
    // Option<T>
    let some_value = Some(42);
    let none_value: Option<i32> = None;
    println!("Option<T> - Some: {:?}, None: {:?}", some_value, none_value);
    
    // Result<T, E>
    let ok_result: Result<i32, &str> = Ok(100);
    let err_result: Result<i32, &str> = Err("错误信息");
    println!("Result<T, E> - Ok: {:?}, Err: {:?}", ok_result, err_result);
}

// 7. 泛型的性能考虑

// Rust的泛型是零成本抽象：
// 1. 编译时会为每种具体类型生成专门的代码（单态化）
// 2. 运行时没有额外开销
// 3. 与手动编写每种类型的代码性能相同

fn performance_example() {
    println!("\n=== 泛型的性能考虑 ===");
    println!("Rust的泛型是零成本抽象：");
    println!("1. 编译时单态化：为每种具体类型生成专门的代码");
    println!("2. 运行时无额外开销：与手动编写每种类型的代码性能相同");
    println!("3. 类型安全：编译时进行类型检查");
}

// 8. 高级泛型特性

// 关联类型（Associated Types）
// 关联类型允许trait定义一个类型占位符，由实现该trait的类型来指定具体类型

pub trait Iterator {
    type Item; // 关联类型
    
    fn next(&mut self) -> Option<Self::Item>;
}

// 示例：实现一个简单的迭代器
struct Counter {
    count: u32,
}

impl Counter {
    fn new() -> Counter {
        Counter { count: 0 }
    }
}

impl Iterator for Counter {
    type Item = u32;
    
    fn next(&mut self) -> Option<Self::Item> {
        if self.count < 5 {
            self.count += 1;
            Some(self.count)
        } else {
            None
        }
    }
}

// 关联类型的使用
fn associated_types_example() {
    println!("\n=== 关联类型 ===");
    
    let mut counter = Counter::new();
    
    println!("Counter迭代器的结果:");
    while let Some(value) = counter.next() {
        println!("{}", value);
    }
}

// 9. 泛型与trait对象

// 泛型用于在编译时确定类型（静态分发）
// trait对象用于在运行时确定类型（动态分发）

// 静态分发：编译时生成专门的代码
fn static_dispatch<T: Debug>(value: T) {
    println!("静态分发: {:?}", value);
}

// 动态分发：运行时通过虚表查找
fn dynamic_dispatch(value: &dyn Debug) {
    println!("动态分发: {:?}", value);
}

// 泛型与trait对象的比较
fn generics_vs_trait_objects() {
    println!("\n=== 泛型与trait对象 ===");
    
    let x = 42;
    let s = String::from("hello");
    
    // 静态分发
    static_dispatch(x);
    static_dispatch(s.clone());
    
    // 动态分发
    dynamic_dispatch(&x);
    dynamic_dispatch(&s);
    
    println!("\n静态分发 vs 动态分发:");
    println!("- 静态分发: 编译时生成专门的代码，性能更好");
    println!("- 动态分发: 运行时通过虚表查找，更灵活，支持异质集合");
}

// 10. 泛型的实际应用

// 实现一个泛型栈
struct Stack<T> {
    elements: Vec<T>,
}

impl<T> Stack<T> {
    // 创建新栈
    fn new() -> Self {
        Stack { elements: Vec::new() }
    }
    
    // 压入元素
    fn push(&mut self, element: T) {
        self.elements.push(element);
    }
    
    // 弹出元素
    fn pop(&mut self) -> Option<T> {
        self.elements.pop()
    }
    
    // 查看栈顶元素
    fn peek(&self) -> Option<&T> {
        self.elements.last()
    }
    
    // 检查栈是否为空
    fn is_empty(&self) -> bool {
        self.elements.is_empty()
    }
    
    // 获取栈的大小
    fn size(&self) -> usize {
        self.elements.len()
    }
}

// 泛型栈的使用
fn generic_stack_example() {
    println!("\n=== 泛型栈的实际应用 ===");
    
    // 创建一个存储i32的栈
    let mut int_stack = Stack::new();
    int_stack.push(1);
    int_stack.push(2);
    int_stack.push(3);
    
    println!("int_stack的大小: {}", int_stack.size());
    println!("int_stack的栈顶元素: {:?}", int_stack.peek());
    
    while let Some(value) = int_stack.pop() {
        println!("弹出: {}", value);
    }
    
    println!("int_stack是否为空: {}", int_stack.is_empty());
    
    // 创建一个存储String的栈
    let mut string_stack = Stack::new();
    string_stack.push(String::from("hello"));
    string_stack.push(String::from("world"));
    
    println!("\nstring_stack的大小: {}", string_stack.size());
    println!("string_stack的栈顶元素: {:?}", string_stack.peek());
}

// 11. 泛型约束的高级用法

// 使用std::ops模块的trait实现运算符重载
use std::ops::Add;

#[derive(Debug)]
struct Point<T> {
    x: T,
    y: T,
}

// 为Point<T>实现Add trait
impl<T: Add<Output = T>> Add for Point<T> {
    type Output = Point<T>;
    
    fn add(self, other: Point<T>) -> Point<T> {
        Point {
            x: self.x + other.x,
            y: self.y + other.y,
        }
    }
}

// 使用Into trait进行类型转换
fn sum<T: Into<i32>>(a: T, b: T) -> i32 {
    a.into() + b.into()
}

// 高级泛型约束的使用
fn advanced_constraints_example() {
    println!("\n=== 高级泛型约束 ===");
    
    // 使用Add trait
    let p1 = Point { x: 1, y: 2 };
    let p2 = Point { x: 3, y: 4 };
    let p3 = p1 + p2;
    
    println!("p1 + p2 = {:?}", p3);
    
    // 使用Into trait
    println!("sum(1, 2) = {}", sum(1, 2));
    println!("sum(1.5, 2.5) = {}", sum(1.5, 2.5));
    println!("sum('a' as u8, 'b' as u8) = {}", sum('a' as u8, 'b' as u8));
}

pub fn run_example() {
    println!("=== Rust学习示例 ===\n");
    // 运行所有示例
    generic_function_example();
    generic_struct_example();
    generic_enum_example();
    generic_constraints_example();
    generics_and_ownership();
    std_lib_generics_example();
    performance_example();
    associated_types_example();
    generics_vs_trait_objects();
    generic_stack_example();
    advanced_constraints_example();
    
    // 12. 泛型总结
    println!("\n=== 泛型总结 ===");
    println!("1. 泛型允许编写通用代码，支持多种数据类型");
    println!("2. 泛型包括：泛型函数、泛型结构体、泛型枚举");
    println!("3. 泛型约束使用trait bounds限制泛型类型");
    println!("4. where子句可以简化复杂的trait bounds");
    println!("5. 泛型是零成本抽象：编译时单态化，运行时无额外开销");
    println!("6. 泛型与trait对象的选择：");
    println!("   - 静态分发（泛型）：性能优先，编译时确定类型");
    println!("   - 动态分发（trait对象）：灵活性优先，运行时确定类型");
    println!("7. 标准库广泛使用泛型：Vec<T>, HashMap<K, V>, Option<T>, Result<T, E>等");
}

// 13. 泛型的最佳实践

// 1. 保持泛型参数尽可能简单
// 2. 合理使用trait bounds，避免过度约束
// 3. 考虑使用where子句提高可读性
// 4. 对于性能敏感的代码，优先使用泛型（静态分发）
// 5. 对于需要运行时多态的场景，使用trait对象（动态分发）
// 6. 为泛型类型提供有用的文档，说明泛型约束的用途
// 7. 考虑实现常见的trait（如Debug, Clone, PartialEq等）以提高可用性

// 14. 泛型与生命周期

// 泛型可以与生命周期参数结合使用
// 这将在13_lifetimes.rs中详细介绍

// 示例：带有生命周期参数的泛型函数
/*
fn longest<'a, T>(x: &'a T, y: &'a T) -> &'a T
where
    T: PartialOrd,
{
    if x > y {
        x
    } else {
        y
    }
}
*/

// 15. 泛型与宏

// 宏也可以实现类似泛型的功能，但宏是在语法层面进行展开
// 泛型是在类型层面进行抽象

// 示例：使用宏实现一个简单的泛型打印函数
/*
macro_rules! print_generic {
    ($x:expr) => {
        println!("值: {}", $x);
    };
}
*/

// 泛型与宏的比较：
// - 泛型：类型安全，编译时检查，性能好
// - 宏：更灵活，可以处理任意语法，但可能导致复杂的错误信息
// 用于单独运行本文件的main函数
fn main() {
    run_example();
}
