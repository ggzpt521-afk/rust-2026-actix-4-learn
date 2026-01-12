// Rust生命周期详解
// 生命周期是Rust中管理引用有效性的核心机制，用于防止悬垂引用和确保内存安全

pub fn run_example() {
    println!("=== Rust生命周期示例 ===\n");

    // 基本生命周期概念
    basic_lifetime_example();
    
    // 函数中的生命周期
    function_lifetime_example();
    
    // 结构体中的生命周期
    struct_lifetime_example();
    
    // 生命周期省略规则
    lifetime_elision_example();
    
    // 静态生命周期
    static_lifetime_example();
    
    // 生命周期约束
    lifetime_bound_example();
    
    // 生命周期子类型
    lifetime_subtyping_example();
    
    // 生命周期与Trait对象
    lifetime_trait_object_example();
    
    // 高级生命周期示例
    advanced_lifetime_example();
    
    println!("\n=== 生命周期示例结束 ===");
}

fn main() {
    run_example();
}

// longest函数：比较两个字符串引用，返回较长的那个
fn longest<'a>(x: &'a str, y: &'a str) -> &'a str {
    if x.len() > y.len() {
        x
    } else {
        y
    }
}

// 1. 基本生命周期概念
fn basic_lifetime_example() {
    println!("1. 基本生命周期概念:");
    
    // 生命周期是引用有效的范围
    let string1 = String::from("hello"); // string1的生命周期开始
    
    {
        let string2 = String::from("world"); // string2的生命周期开始
        let result = longest(string1.as_str(), string2.as_str());
        println!("较长的字符串是: {}", result);
    } // string2的生命周期结束
    
    // string1的生命周期结束
    
    println!();
}

// 2. 函数中的生命周期标注
fn function_lifetime_example() {
    println!("2. 函数中的生命周期标注:");
    
    // 生命周期参数使用撇号(')开头
    // 这个函数返回两个参数中较长的那个字符串的引用
    // longest函数已在外部定义
    
    let string1 = String::from("long string is long");
    let result;
    
    {
        let string2 = String::from("xyz");
        result = longest(string1.as_str(), string2.as_str());
        println!("在作用域内，较长的字符串是: {}", result);
    }
    
    // 错误: string2已经离开作用域，result指向无效的内存
    // println!("在作用域外: {}", result);
    
    println!("string1仍然有效: {}", string1);
    println!();
}

// 3. 结构体中的生命周期
fn struct_lifetime_example() {
    println!("3. 结构体中的生命周期:");
    
    // 包含引用的结构体必须标注生命周期
    struct ImportantExcerpt<'a> {
        part: &'a str,
    }
    
    let novel = String::from("Call me Ishmael. Some years ago...");
    let first_sentence = novel.split('.')
        .next()
        .expect("Could not find a '.'");
    
    let i = ImportantExcerpt {
        part: first_sentence,
    };
    
    println!("重要摘录: {}", i.part);
    
    // 结构体更新语法与生命周期
    let i2 = ImportantExcerpt {
        ..i
    };
    
    println!("更新后的摘录: {}", i2.part);
    println!();
}

// 4. 生命周期省略规则
fn lifetime_elision_example() {
    println!("4. 生命周期省略规则:");
    
    // Rust编译器有一套生命周期省略规则，可以在某些情况下省略显式标注
    
    // 规则1: 每个引用参数获得自己的生命周期参数
    // 例如: fn first_word(s: &str) -> &str { ... }
    // 等同于: fn first_word<'a>(s: &'a str) -> &str { ... }
    
    // 规则2: 如果只有一个输入生命周期参数，它被赋予所有输出生命周期参数
    // 例如: fn first_word<'a>(s: &'a str) -> &str { ... }
    // 等同于: fn first_word<'a>(s: &'a str) -> &'a str { ... }
    
    // 规则3: 如果有多个输入生命周期参数，但其中一个是&self或&mut self（方法），
    //        self的生命周期被赋予所有输出生命周期参数
    
    // 示例: 省略生命周期标注的函数
    fn first_word(s: &str) -> &str {
        let bytes = s.as_bytes();
        
        for (i, &item) in bytes.iter().enumerate() {
            if item == b' ' {
                return &s[0..i];
            }
        }
        
        &s[..]
    }
    
    let s = String::from("hello world");
    let word = first_word(&s);
    println!("第一个单词: {}", word);
    
    // 示例: 方法中的生命周期省略
    struct StrSlice {
        s: String,
    }
    
    impl StrSlice {
        // 这里应用了规则3
        fn get_slice(&self, start: usize, end: usize) -> &str {
            &self.s[start..end]
        }
    }
    
    let slice = StrSlice { s: String::from("Rust programming") };
    println!("切片: {}", slice.get_slice(0, 4));
    println!();
}

// 5. 静态生命周期
fn static_lifetime_example() {
    println!("5. 静态生命周期:");
    
    // 'static生命周期表示引用在整个程序运行期间都有效
    // 字符串字面量具有'static生命周期
    let s: &'static str = "我是一个静态字符串";
    println!("静态字符串: {}", s);
    
    // 显式声明静态生命周期
    fn get_static_str() -> &'static str {
        "这也是一个静态字符串"
    }
    
    println!("函数返回的静态字符串: {}", get_static_str());
    
    // 静态生命周期与其他生命周期
    fn longest_with_static<'a>(x: &'a str, y: &'static str) -> &'a str {
        if x.len() > y.len() {
            x
        } else {
            y
        }
    }
    
    let my_string = String::from("hello");
    let result = longest_with_static(my_string.as_str(), "world");
    println!("较长的字符串: {}", result);
    println!();
}

// 6. 生命周期约束
fn lifetime_bound_example() {
    println!("6. 生命周期约束:");
    
    use std::fmt::Display;
    
    // 同时包含生命周期和Trait约束
    fn longest_with_an_announcement<'a, T>(
        x: &'a str,
        y: &'a str,
        ann: T,
    ) -> &'a str
    where
        T: Display,
    {
        println!("公告: {}", ann);
        if x.len() > y.len() {
            x
        } else {
            y
        }
    }
    
    let string1 = String::from("hello");
    let string2 = String::from("world");
    let announcement = "这是一个重要公告!";
    
    let result = longest_with_an_announcement(string1.as_str(), string2.as_str(), announcement);
    println!("较长的字符串: {}", result);
    
    // 使用结构体的生命周期约束
    struct WithDisplay<'a, T: Display + 'a> {
        item: &'a T,
    }
    
    let num = 42;
    let with_display = WithDisplay { item: &num };
    println!("带Display的结构体: {}", with_display.item);
    println!();
}

// 7. 生命周期子类型
fn lifetime_subtyping_example() {
    println!("7. 生命周期子类型:");
    
    // 生命周期子类型表示一个生命周期比另一个生命周期更长
    // 使用'sub: 'super表示'sub是'super的子类型，即'sub比'super更长或相等
    
    fn longer_lifetime<'a, 'b: 'a>(x: &'a str, y: &'b str) -> &'a str {
        if x.len() > y.len() {
            x
        } else {
            // y的生命周期至少与x一样长，所以可以安全返回
            y
        }
    }
    
    let string1 = String::from("long string");
    let result;
    
    {
        let string2 = String::from("short");
        // string1的生命周期比string2长
        result = longer_lifetime(string2.as_str(), string1.as_str());
        println!("较长的字符串: {}", result);
    }
    
    println!("string1仍然有效: {}", string1);
    println!();
}

// 8. 生命周期与Trait对象
fn lifetime_trait_object_example() {
    println!("8. 生命周期与Trait对象:");
    
    trait Drawable {
        fn draw(&self);
    }
    
    struct Circle {
        radius: f64,
    }
    
    impl Drawable for Circle {
        fn draw(&self) {
            println!("绘制圆形，半径: {}", self.radius);
        }
    }
    
    struct Square {
        side: f64,
    }
    
    impl Drawable for Square {
        fn draw(&self) {
            println!("绘制正方形，边长: {}", self.side);
        }
    }
    
    // Trait对象中的生命周期
    fn get_drawable<'a>(shape_type: &str) -> Box<dyn Drawable + 'a> {
        match shape_type {
            "circle" => Box::new(Circle { radius: 1.5 }),
            "square" => Box::new(Square { side: 2.0 }),
            _ => Box::new(Circle { radius: 1.0 }),
        }
    }
    
    // 静态Trait对象（没有生命周期限制）
    fn get_static_drawable(shape_type: &str) -> Box<dyn Drawable + 'static> {
        get_drawable(shape_type)
    }
    
    let drawable = get_static_drawable("circle");
    drawable.draw();
    
    // 带有生命周期的Trait对象集合
    let mut drawables: Vec<Box<dyn Drawable + 'static>> = Vec::new();
    drawables.push(Box::new(Circle { radius: 3.0 }));
    drawables.push(Box::new(Square { side: 4.0 }));
    
    println!("绘制所有图形:");
    for d in drawables {
        d.draw();
    }
    println!();
}

// 9. 高级生命周期示例
fn advanced_lifetime_example() {
    println!("9. 高级生命周期示例:");
    
    // 多个生命周期参数
    fn multiple_lifetimes<'a, 'b>(x: &'a str, y: &'b str) -> (&'a str, &'b str) {
        (x, y)
    }
    
    let string1 = String::from("first");
    let string2 = String::from("second");
    let (result1, result2) = multiple_lifetimes(string1.as_str(), string2.as_str());
    println!("多生命周期参数结果: {}, {}", result1, result2);
    
    // 生命周期与泛型结合
    struct Wrapper<'a, T> {
        value: &'a T,
    }
    
    impl<'a, T> Wrapper<'a, T> {
        fn new(value: &'a T) -> Self {
            Wrapper { value }
        }
        
        fn get_value(&self) -> &'a T {
            self.value
        }
    }
    
    let num = 42;
    let wrapper = Wrapper::new(&num);
    println!("包装器中的值: {}", wrapper.get_value());
    
    // 生命周期与可变引用
    fn longest_mut<'a>(x: &'a mut str, y: &'a mut str) -> &'a mut str {
        if x.len() > y.len() {
            x
        } else {
            y
        }
    }
    
    // 需要使用String的可变切片
    let mut s1 = String::from("hello");
    let mut s2 = String::from("world");
    
    // 注意：这里不能同时借用两个可变引用给同一个函数
    // 因为Rust的借用规则不允许同时有多个可变引用指向同一数据
    // let result = longest_mut(&mut s1, &mut s2);
    // println!("较长的可变字符串: {}", result);
    
    println!();
}

// 10. 生命周期与实际应用
fn practical_lifetime_example() {
    println!("10. 生命周期与实际应用:");
    
    // 实际应用中的生命周期使用
    struct Document {
        content: String,
        tags: Vec<String>,
    }
    
    struct DocumentView<'a> {
        document: &'a Document,
        start: usize,
        end: usize,
    }
    
    impl<'a> DocumentView<'a> {
        fn new(document: &'a Document, start: usize, end: usize) -> Self {
            DocumentView {
                document,
                start: start.min(document.content.len()),
                end: end.min(document.content.len()),
            }
        }
        
        fn get_content(&self) -> &'a str {
            &self.document.content[self.start..self.end]
        }
        
        fn get_tags(&self) -> &'a [String] {
            &self.document.tags
        }
    }
    
    let doc = Document {
        content: String::from("Rust是一种系统编程语言，专注于安全性、并发和性能。"),
        tags: vec!["编程".to_string(), "Rust".to_string(), "系统编程".to_string()],
    };
    
    let view = DocumentView::new(&doc, 0, 15);
    println!("文档视图内容: {}", view.get_content());
    println!("文档标签: {:?}", view.get_tags());
    
    // 扩展文档视图
    let extended_view = DocumentView::new(&doc, 15, 30);
    println!("扩展视图内容: {}", extended_view.get_content());
    println!();
}