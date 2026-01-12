// 06_structs.rs - Rust结构体详解

// 1. 结构体的定义
// 使用struct关键字定义结构体
struct User {
    active: bool,
    username: String,
    email: String,
    sign_in_count: u64,
}

// 2. 元组结构体（Tuple Structs）
// 元组结构体没有字段名，只有类型
struct Color(i32, i32, i32);
struct Point(i32, i32, i32);

// 3. 单元结构体（Unit Structs）
// 没有任何字段的结构体，类似于单元类型()
struct Unit;

// 4. 结构体的方法定义（使用impl块）
impl User {
    // 实例方法，第一个参数是self
    fn get_email(&self) -> &String {
        &self.email
    }
    
    // 可变实例方法，可以修改self
    fn update_email(&mut self, new_email: String) {
        self.email = new_email;
    }
    
    // 关联函数（静态方法），不需要self
    fn new_user(username: String, email: String) -> User {
        User {
            active: true,
            username,
            email,
            sign_in_count: 1,
        }
    }
    
    // 关联函数示例：创建不活跃用户
    fn new_inactive_user(username: String, email: String) -> User {
        User {
            active: false,
            username,
            email,
            sign_in_count: 0,
        }
    }
}

// 5. 结构体可见性示例
// 公开结构体（使用pub关键字）
pub struct PublicStruct {
    pub public_field: i32,  // 公开字段
    private_field: i32,     // 私有字段（默认）
}

impl PublicStruct {
    pub fn new() -> PublicStruct {
        PublicStruct {
            public_field: 0,
            private_field: 10,
        }
    }
    
    pub fn get_private(&self) -> i32 {
        self.private_field
    }
    
    pub fn set_private(&mut self, value: i32) {
        self.private_field = value;
    }
}

// 6. 结构体所有权示例
struct Article {
    title: String,
    content: String,
    author: String,
}

impl Article {
    // 获取所有权的方法
    fn take_title(self) -> String {
        self.title
    }
    
    // 借用的方法
    fn get_title(&self) -> &str {
        &self.title
    }
}

pub fn run_example() {
    println!("=== Rust学习示例 ===\n");
    // 7. 结构体的实例化
    let user1 = User {
        active: true,
        username: String::from("alice"),
        email: String::from("alice@example.com"),
        sign_in_count: 1,
    };
    
    // 8. 结构体字段访问
    println!("用户名: {}, 邮箱: {}", user1.username, user1.email);
    
    // 9. 结构体的可变实例
    // 整个结构体必须是可变的，Rust不允许只将部分字段标记为可变
    let mut user2 = User {
        active: true,
        username: String::from("bob"),
        email: String::from("bob@example.com"),
        sign_in_count: 2,
    };
    
    // 修改字段
    user2.email = String::from("bob_new@example.com");
    println!("修改后的邮箱: {}", user2.email);
    
    // 10. 结构体更新语法
    // 可以基于现有结构体创建新结构体，只修改需要的字段
    let user3 = User {
        email: String::from("charlie@example.com"),
        ..user1  // 使用user1的其他字段
    };
    
    println!("user3的用户名: {}, 邮箱: {}", user3.username, user3.email);
    // 注意：user1的username字段被移动到了user3，因为它是String类型
    // println!("user1的用户名: {}", user1.username); // 这会报错
    
    // 11. 元组结构体的实例化和访问
    let black = Color(0, 0, 0);
    let origin = Point(0, 0, 0);
    
    println!("黑色的RGB值: {}, {}, {}", black.0, black.1, black.2);
    println!("原点坐标: {}, {}, {}", origin.0, origin.1, origin.2);
    
    // 注意：虽然Color和Point有相同的结构，但它们是不同的类型
    // let mix: Color = origin; // 这会报错，类型不匹配
    
    // 12. 单元结构体的实例化
    let unit = Unit;
    println!("单元结构体: {:?}", unit); // 需要Debug trait才能打印
    
    // 13. 调用结构体方法
    let user4 = User::new_user(
        String::from("david"),
        String::from("david@example.com")
    );
    
    println!("user4的邮箱: {}", user4.get_email());
    
    let mut user5 = User::new_inactive_user(
        String::from("eve"),
        String::from("eve@example.com")
    );
    
    println!("user5是否活跃: {}", user5.active);
    user5.update_email(String::from("eve_new@example.com"));
    println!("user5修改后的邮箱: {}", user5.get_email());
    
    // 14. 结构体所有权示例
    let article = Article {
        title: String::from("Rust结构体教程"),
        content: String::from("这是一篇关于Rust结构体的教程..."),
        author: String::from("Rust爱好者"),
    };
    
    println!("文章标题: {}", article.get_title());
    
    let title = article.take_title(); // 获取标题的所有权
    println!("获取到的标题: {}", title);
    // println!("文章标题: {}", article.get_title()); // 这会报错，因为article已经失去了title的所有权
    
    // 15. 公开结构体示例
    let mut public_struct = PublicStruct::new();
    println!("公开字段: {}, 私有字段: {}", 
             public_struct.public_field, 
             public_struct.get_private());
    
    public_struct.public_field = 20;
    public_struct.set_private(30);
    println!("修改后 - 公开字段: {}, 私有字段: {}", 
             public_struct.public_field, 
             public_struct.get_private());
}

// 16. 结构体的示例应用：矩形
struct Rectangle {
    width: u32,
    height: u32,
}

impl Rectangle {
    // 计算面积
    fn area(&self) -> u32 {
        self.width * self.height
    }
    
    // 检查是否能容纳另一个矩形
    fn can_hold(&self, other: &Rectangle) -> bool {
        self.width > other.width && self.height > other.height
    }
    
    // 创建正方形（关联函数）
    fn square(size: u32) -> Rectangle {
        Rectangle {
            width: size,
            height: size,
        }
    }
}

fn rectangle_example() {
    let rect1 = Rectangle {
        width: 30,
        height: 50,
    };
    
    let rect2 = Rectangle {
        width: 10,
        height: 40,
    };
    
    let rect3 = Rectangle {
        width: 60,
        height: 45,
    };
    
    let square = Rectangle::square(20);
    
    println!("\n矩形示例:");
    println!("rect1的面积: {}", rect1.area());
    println!("rect1能否容纳rect2: {}", rect1.can_hold(&rect2));
    println!("rect1能否容纳rect3: {}", rect1.can_hold(&rect3));
    println!("正方形的面积: {}", square.area());
}
// 用于单独运行本文件的main函数
fn main() {
    run_example();
}
