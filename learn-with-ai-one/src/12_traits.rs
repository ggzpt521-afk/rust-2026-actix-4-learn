// Rust Trait系统详解
// Trait是Rust中实现代码复用和多态的核心机制，类似于其他语言中的接口但功能更强大

pub fn run_example() {
    println!("=== Rust学习示例 ===\n");
    println!("=== Rust Trait系统示例 ===\n");

    // 基本Trait实现和使用
    trait_example();
    
    // 默认方法
    default_method_example();
    
    // Trait作为参数和返回值
    trait_as_param_return_example();
    
    // Trait约束
    trait_bound_example();
    
    // 多Trait约束
    multiple_trait_bounds_example();
    
    // 泛型Trait
    generic_trait_example();
    
    // 关联类型
    associated_type_example();
    
    // Trait对象
    trait_object_example();
    
    // 运算符重载
    operator_overloading_example();
    
    // 标准库中的Trait
    std_trait_example();
    
    println!("\n=== Trait系统示例结束 ===");
}

// 1. 基本Trait定义和实现
fn trait_example() {
    println!("1. 基本Trait定义和实现:");
    
    // Trait定义
    trait Shape {
        fn area(&self) -> f64;
        fn perimeter(&self) -> f64;
    }
    
    // 实现Trait的结构体
    struct Rectangle {
        width: f64,
        height: f64,
    }
    
    struct Circle {
        radius: f64,
    }
    
    // 为Rectangle实现Shape Trait
    impl Shape for Rectangle {
        fn area(&self) -> f64 {
            self.width * self.height
        }
        
        fn perimeter(&self) -> f64 {
            2.0 * (self.width + self.height)
        }
    }
    
    // 为Circle实现Shape Trait
    impl Shape for Circle {
        fn area(&self) -> f64 {
            std::f64::consts::PI * self.radius * self.radius
        }
        
        fn perimeter(&self) -> f64 {
            2.0 * std::f64::consts::PI * self.radius
        }
    }
    
    // 使用实现了Trait的类型
    let rect = Rectangle { width: 5.0, height: 3.0 };
    let circle = Circle { radius: 2.5 };
    
    println!("矩形面积: {}", rect.area());
    println!("矩形周长: {}", rect.perimeter());
    println!("圆形面积: {}", circle.area());
    println!("圆形周长: {}", circle.perimeter());
    println!();
}

// 2. 默认方法
fn default_method_example() {
    println!("2. 默认方法:");
    
    trait Animal {
        fn speak(&self);
        
        // 默认方法
        fn sleep(&self) {
            println!("正在睡觉...");
        }
    }
    
    struct Dog;
    struct Cat;
    
    impl Animal for Dog {
        fn speak(&self) {
            println!("汪汪汪!");
        }
        
        // 可以重写默认方法
        fn sleep(&self) {
            println!("狗在睡觉...");
        }
    }
    
    impl Animal for Cat {
        fn speak(&self) {
            println!("喵喵喵!");
        }
        // 使用默认的sleep方法
    }
    
    let dog = Dog;
    let cat = Cat;
    
    dog.speak();
    dog.sleep();
    cat.speak();
    cat.sleep();
    println!();
}

// 3. Trait作为参数和返回值
fn trait_as_param_return_example() {
    println!("3. Trait作为参数和返回值:");
    
    trait Drawable {
        fn draw(&self);
    }
    
    struct Circle { radius: f64 }
    struct Square { side: f64 }
    
    impl Drawable for Circle {
        fn draw(&self) {
            println!("绘制一个半径为{}的圆形", self.radius);
        }
    }
    
    impl Drawable for Square {
        fn draw(&self) {
            println!("绘制一个边长为{}的正方形", self.side);
        }
    }
    
    // Trait作为参数（静态分发）
    fn draw_shape<T: Drawable>(shape: T) {
        shape.draw();
    }
    
    // Trait作为返回值（静态分发，只能返回一种类型）
    fn create_shape(shape_type: &str) -> impl Drawable {
        match shape_type {
            "circle" => Circle { radius: 2.0 },
            "square" => Square { side: 3.0 },
            _ => Circle { radius: 1.0 },
        }
    }
    
    let circle = Circle { radius: 1.5 };
    let square = Square { side: 2.5 };
    
    draw_shape(circle);
    draw_shape(square);
    
    let shape = create_shape("circle");
    shape.draw();
    println!();
}

// 4. Trait约束
fn trait_bound_example() {
    println!("4. Trait约束:");
    
    trait Printable {
        fn print(&self);
    }
    
    struct Person {
        name: String,
        age: u32,
    }
    
    impl Printable for Person {
        fn print(&self) {
            println!("姓名: {}, 年龄: {}", self.name, self.age);
        }
    }
    
    struct Car {
        brand: String,
        model: String,
    }
    
    impl Printable for Car {
        fn print(&self) {
            println!("汽车品牌: {}, 型号: {}", self.brand, self.model);
        }
    }
    
    // 带Trait约束的泛型函数
    fn print_item<T: Printable>(item: T) {
        item.print();
    }
    
    // 使用where子句的Trait约束（更易读的方式）
    fn print_items<T>(items: Vec<T>) where T: Printable {
        for item in items {
            item.print();
        }
    }
    
    let person = Person { name: "张三".to_string(), age: 30 };
    let car = Car { brand: "丰田".to_string(), model: "卡罗拉".to_string() };
    
    print_item(person);
    print_item(car);
    
    println!("\n打印多个项目:");
    let people = vec![
        Person { name: "李四".to_string(), age: 25 },
        Person { name: "王五".to_string(), age: 35 },
    ];
    print_items(people);
    println!();
}

// 5. 多Trait约束
fn multiple_trait_bounds_example() {
    println!("5. 多Trait约束:");
    
    trait Runnable {
        fn run(&self);
    }
    
    trait Swimmable {
        fn swim(&self);
    }
    
    struct Human {
        name: String,
    }
    
    // 实现多个Trait
    impl Runnable for Human {
        fn run(&self) {
            println!("{}在跑步", self.name);
        }
    }
    
    impl Swimmable for Human {
        fn swim(&self) {
            println!("{}在游泳", self.name);
        }
    }
    
    // 为了使用Display Trait，需要导入
    use std::fmt;
    
    impl fmt::Display for Human {
        fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
            write!(f, "{}先生", self.name)
        }
    }
    
    // 多Trait约束
    fn perform_activities<T>(entity: T) where 
        T: Runnable + Swimmable + fmt::Display {
        
        println!("实体: {}", entity);
        entity.run();
        entity.swim();
    }
    
    let person = Human { name: "赵六".to_string() };
    perform_activities(person);
    println!();
}

// 6. 泛型Trait
fn generic_trait_example() {
    println!("6. 泛型Trait:");
    
    // 泛型Trait定义
    trait Container<T> {
        fn add(&mut self, item: T);
        fn remove(&mut self) -> Option<T>;
        fn is_empty(&self) -> bool;
    }
    
    // 实现泛型Trait
    struct SimpleStack<T> {
        items: Vec<T>,
    }
    
    impl<T> SimpleStack<T> {
        fn new() -> Self {
            SimpleStack { items: Vec::new() }
        }
    }
    
    impl<T> Container<T> for SimpleStack<T> {
        fn add(&mut self, item: T) {
            self.items.push(item);
        }
        
        fn remove(&mut self) -> Option<T> {
            self.items.pop()
        }
        
        fn is_empty(&self) -> bool {
            self.items.is_empty()
        }
    }
    
    // 使用泛型Trait
    let mut stack = SimpleStack::new();
    println!("栈是否为空: {}", stack.is_empty());
    
    stack.add(10);
    stack.add(20);
    stack.add(30);
    
    println!("栈是否为空: {}", stack.is_empty());
    
    while let Some(item) = stack.remove() {
        println!("弹出元素: {}", item);
    }
    
    println!("栈是否为空: {}", stack.is_empty());
    println!();
}

// 7. 关联类型
fn associated_type_example() {
    println!("7. 关联类型:");
    
    // 带有关联类型的Trait
    trait Iterator {
        type Item;
        
        fn next(&mut self) -> Option<Self::Item>;
        
        // 使用关联类型的默认方法
        fn count(self) -> usize
        where
            Self: Sized,
        {
            let mut count = 0;
            while let Some(_) = self.next() {
                count += 1;
            }
            count
        }
    }
    
    // 实现带有关联类型的Trait
    struct Counter {
        current: u32,
        max: u32,
    }
    
    impl Counter {
        fn new(max: u32) -> Self {
            Counter { current: 0, max }
        }
    }
    
    impl Iterator for Counter {
        // 定义关联类型
        type Item = u32;
        
        fn next(&mut self) -> Option<Self::Item> {
            if self.current < self.max {
                let value = self.current;
                self.current += 1;
                Some(value)
            } else {
                None
            }
        }
    }
    
    let mut counter = Counter::new(5);
    println!("计数器元素:");
    while let Some(num) = counter.next() {
        println!("{}", num);
    }
    
    let counter2 = Counter::new(10);
    println!("计数器元素数量: {}", counter2.count());
    println!();
}

// 8. Trait对象（动态分发）
fn trait_object_example() {
    println!("8. Trait对象:");
    
    trait Drawable {
        fn draw(&self);
    }
    
    struct Circle { radius: f64 }
    struct Square { side: f64 }
    struct Triangle { base: f64, height: f64 }
    
    impl Drawable for Circle {
        fn draw(&self) {
            println!("绘制圆形，半径: {}", self.radius);
        }
    }
    
    impl Drawable for Square {
        fn draw(&self) {
            println!("绘制正方形，边长: {}", self.side);
        }
    }
    
    impl Drawable for Triangle {
        fn draw(&self) {
            println!("绘制三角形，底: {}, 高: {}", self.base, self.height);
        }
    }
    
    // 创建Trait对象向量（动态分发）
    let shapes: Vec<Box<dyn Drawable>> = vec![
        Box::new(Circle { radius: 1.5 }),
        Box::new(Square { side: 2.5 }),
        Box::new(Triangle { base: 3.0, height: 2.0 }),
    ];
    
    // 遍历并调用Trait方法
    println!("绘制所有图形:");
    for shape in shapes {
        shape.draw();
    }
    println!();
}

// 9. 运算符重载
fn operator_overloading_example() {
    println!("9. 运算符重载:");
    
    use std::ops::{Add, Sub, Mul, Div};
    
    struct Vector2D {
        x: f64,
        y: f64,
    }
    
    // 实现Add Trait重载+运算符
    impl Add for Vector2D {
        type Output = Vector2D;
        
        fn add(self, other: Vector2D) -> Vector2D {
            Vector2D {
                x: self.x + other.x,
                y: self.y + other.y,
            }
        }
    }
    
    // 实现Sub Trait重载-运算符
    impl Sub for Vector2D {
        type Output = Vector2D;
        
        fn sub(self, other: Vector2D) -> Vector2D {
            Vector2D {
                x: self.x - other.x,
                y: self.y - other.y,
            }
        }
    }
    
    // 实现Mul Trait重载*运算符（标量乘法）
    impl Mul<f64> for Vector2D {
        type Output = Vector2D;
        
        fn mul(self, scalar: f64) -> Vector2D {
            Vector2D {
                x: self.x * scalar,
                y: self.y * scalar,
            }
        }
    }
    
    let v1 = Vector2D { x: 1.0, y: 2.0 };
    let v2 = Vector2D { x: 3.0, y: 4.0 };
    
    let v3 = v1 + v2;
    println!("向量加法: ({}, {})
", v3.x, v3.y);
    
    let v4 = v3 - Vector2D { x: 2.0, y: 3.0 };
    println!("向量减法: ({}, {})
", v4.x, v4.y);
    
    let v5 = v4 * 2.0;
    println!("向量标量乘法: ({}, {})
", v5.x, v5.y);
    
    // 显示使用Add Trait的add方法
    let v6 = Vector2D::add(Vector2D { x: 1.0, y: 1.0 }, Vector2D { x: 2.0, y: 2.0 });
    println!("直接使用add方法: ({}, {})
", v6.x, v6.y);
    println!();
}

// 10. 标准库中的Trait
fn std_trait_example() {
    println!("10. 标准库中的Trait:");
    
    use std::fmt;
    use std::cmp::PartialEq;
    
    struct Person {
        name: String,
        age: u32,
    }
    
    // 实现Debug Trait用于调试输出
    impl fmt::Debug for Person {
        fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
            write!(f, "Person {{ name: {}, age: {} }}", self.name, self.age)
        }
    }
    
    // 实现Display Trait用于友好输出
    impl fmt::Display for Person {
        fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
            write!(f, "{} ({})", self.name, self.age)
        }
    }
    
    // 实现PartialEq Trait用于比较
    impl PartialEq for Person {
        fn eq(&self, other: &Self) -> bool {
            self.name == other.name && self.age == other.age
        }
    }
    
    let person1 = Person { name: "张三".to_string(), age: 30 };
    let person2 = Person { name: "李四".to_string(), age: 25 };
    let person3 = Person { name: "张三".to_string(), age: 30 };
    
    // 使用Debug Trait
    println!("Debug输出: {:?}", person1);
    
    // 使用Display Trait
    println!("Display输出: {}", person1);
    
    // 使用PartialEq Trait
    println!("person1 == person2: {}", person1 == person2);
    println!("person1 == person3: {}", person1 == person3);
    
    // 实现From Trait
    impl From<&str> for Person {
        fn from(s: &str) -> Self {
            let parts: Vec<&str> = s.split(", ").collect();
            Person {
                name: parts[0].to_string(),
                age: parts[1].parse().unwrap_or(0),
            }
        }
    }
    
    // 使用From Trait
    let person4: Person = "王五, 35".into();
    println!("从字符串创建的Person: {}", person4);
    println!();
}

// 11. 高级Trait概念 - 继承Trait
fn trait_inheritance_example() {
    println!("11. Trait继承:");
    
    trait Animal {
        fn eat(&self);
    }
    
    // 继承自Animal的Trait
    trait Mammal: Animal {
        fn sleep(&self);
        fn give_birth(&self);
    }
    
    struct Human {
        name: String,
    }
    
    // 必须先实现父Trait
    impl Animal for Human {
        fn eat(&self) {
            println!("{}正在吃饭", self.name);
        }
    }
    
    // 然后实现子Trait
    impl Mammal for Human {
        fn sleep(&self) {
            println!("{}正在睡觉", self.name);
        }
        
        fn give_birth(&self) {
            println!("{}是胎生动物", self.name);
        }
    }
    
    let person = Human { name: "张三".to_string() };
    person.eat();
    person.sleep();
    person.give_birth();
    println!();
}
// 用于单独运行本文件的main函数
fn main() {
    run_example();
}
