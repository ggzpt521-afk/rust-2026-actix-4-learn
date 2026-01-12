// 09_packages_modules.rs - Rust包和模块系统详解

// Rust使用包(Package)、箱(Crate)和模块(Module)来组织代码
// 1. 包(Package)：一个项目，包含Cargo.toml文件，用于描述项目和依赖
// 2. 箱(Crate)：编译的基本单位，可以是二进制箱或库箱
// 3. 模块(Module)：用于组织代码，控制可见性
// 4. 路径(Path)：用于引用模块、类型、函数等

// 以下是一个模块系统的示例结构：
// my_package/
// ├── Cargo.toml          # 包配置文件
// └── src/                # 源代码目录
//     ├── main.rs         # 二进制箱入口（与包同名）
//     ├── lib.rs          # 库箱入口（与包同名）
//     ├── module1.rs      # 模块1
//     └── module2/
//         ├── mod.rs      # 模块2的入口
//         └── submodule.rs # 子模块

// 由于我们在单个文件中演示，将使用嵌套模块来模拟目录结构

// 1. 定义模块
mod my_module {
    // 模块中的代码
    pub fn public_function() {
        println!("这是一个公共函数");
        private_function(); // 模块内部可以访问私有函数
    }
    
    fn private_function() {
        println!("这是一个私有函数");
    }
    
    // 2. 嵌套模块
    pub mod nested_module {
        pub fn nested_public_function() {
            println!("这是嵌套模块中的公共函数");
        }
        
        fn nested_private_function() {
            println!("这是嵌套模块中的私有函数");
        }
    }
    
    // 3. 结构体的可见性
    pub struct PublicStruct {
        pub public_field: i32,  // 公共字段
        private_field: i32,     // 私有字段
    }
    
    struct PrivateStruct {
        field: i32,
    }
    
    impl PublicStruct {
        pub fn new(public: i32, private: i32) -> Self {
            Self {
                public_field: public,
                private_field: private,
            }
        }
        
        pub fn access_private_field(&self) -> i32 {
            self.private_field
        }
    }
    
    // 4. 枚举的可见性
    pub enum PublicEnum {
        Variant1,              // 枚举变体默认是公共的
        Variant2(i32),         // 带数据的公共变体
        PrivateVariant,        // 枚举变体默认是公共的，即使枚举是公共的
    }
    
    enum PrivateEnum {
        Variant1,
    }
}

// 5. 使用pub mod声明模块（与上面的mod不同，这会创建一个新的模块文件）
// pub mod module1; // 这会加载src/module1.rs文件中的模块

// 6. 引用路径
// 有两种方式引用路径：
// - 绝对路径：从箱根目录开始，使用crate关键字
// - 相对路径：从当前模块开始，使用self、super或模块名

// 7. use关键字：用于导入路径，简化代码
use crate::my_module::{public_function, PublicStruct, nested_module};
use crate::my_module::nested_module::nested_public_function as npf; // 使用as重命名

// 8. 导入整个模块
use crate::my_module; // 导入整个模块

// 9. 使用通配符导入所有公共项
use crate::my_module::*; // 不推荐在生产代码中使用，可能导致名称冲突

// 10. 从外部包导入
// use std::collections::HashMap; // 从标准库导入
// use serde::{Serialize, Deserialize}; // 从外部包导入

pub fn run_example() {
    println!("=== Rust学习示例 ===\n");
    println!("=== Rust包和模块系统 ===");
    
    // 11. 使用绝对路径调用函数
    crate::my_module::public_function();
    
    // 12. 使用导入的函数
    public_function();
    
    // 13. 使用结构体
    let mut s = PublicStruct::new(10, 20);
    println!("公共字段: {}", s.public_field);
    // println!("私有字段: {}", s.private_field); // 这会报错，因为私有字段不可访问
    println!("通过方法访问私有字段: {}", s.access_private_field());
    
    // 14. 使用嵌套模块
    my_module::nested_module::nested_public_function();
    nested_module::nested_public_function();
    npf(); // 使用重命名的函数
    
    // 15. 使用枚举
    let variant = my_module::PublicEnum::Variant2(42);
    match variant {
        my_module::PublicEnum::Variant1 => println!("Variant1"),
        my_module::PublicEnum::Variant2(n) => println!("Variant2: {}", n),
        my_module::PublicEnum::PrivateVariant => println!("PrivateVariant"),
    }
    
    // 16. 演示模块可见性规则
    println!("\n=== 模块可见性规则 ===");
    println!("1. 默认情况下，所有项（函数、结构体、枚举等）都是私有的");
    println!("2. 使用pub关键字使项变为公共的");
    println!("3. 公共结构体的字段默认是私有的，需要单独使用pub关键字");
    println!("4. 公共枚举的变体默认是公共的");
    println!("5. 模块本身默认是私有的，需要使用pub mod使其变为公共的");
    
    // 17. super关键字：用于引用父模块
    println!("\n=== super关键字的使用 ===");
    outer_module::inner_module::call_outer_function();
    
    // 18. self关键字：用于引用当前模块
    println!("\n=== self关键字的使用 ===");
    self::my_module::public_function(); // 等同于crate::my_module::public_function()
}

// 19. 演示super关键字
mod outer_module {
    pub fn outer_function() {
        println!("这是外部模块的函数");
    }
    
    pub mod inner_module {
        pub fn call_outer_function() {
            super::outer_function(); // 使用super引用父模块的函数
        }
    }
}

// 20. 模块文件系统
// Rust的模块系统与文件系统是对应的：
// - 单个文件模块：src/module.rs -> mod module;
// - 目录模块：src/module/mod.rs -> mod module;
// - 子模块：src/module/submodule.rs -> mod submodule;

// 以下是模块文件系统的示例：
// // src/lib.rs
// pub mod module1;
// pub mod module2;

// // src/module1.rs
// pub fn module1_function() {}

// // src/module2/mod.rs
// pub mod submodule;
// pub fn module2_function() {}

// // src/module2/submodule.rs
// pub fn submodule_function() {}

// 21. 可见性规则详细说明
// - private：默认，只能在当前模块内访问
// - pub：公共的，可以被任何模块访问
// - pub(crate)：在当前箱内可见
// - pub(super)：在父模块内可见
// - pub(in path)：在指定路径的模块内可见

mod visibility_demo {
    pub fn public_function() {}
    pub(crate) fn crate_function() {}
    pub(super) fn super_function() {}
    pub(in crate::visibility_demo) fn in_module_function() {}
    fn private_function() {}
    
    mod inner {
        use super::*;
        
        pub fn test_visibility() {
            public_function();      // 可访问
            crate_function();       // 可访问
            super_function();       // 可访问
            in_module_function();   // 可访问
            // private_function();  // 不可访问，因为是父模块的私有函数
        }
    }
}

// 22. 使用use导入多个项
use std::fs::File;
use std::io::{self, Read};

fn use_example() {
    // 使用导入的类型
    let file_result = File::open("test.txt");
    
    match file_result {
        Ok(mut file) => {
            let mut content = String::new();
            match file.read_to_string(&mut content) {
                Ok(_) => println!("文件内容: {}", content),
                Err(e) => println!("读取错误: {}", e),
            }
        },
        Err(e) => println!("打开错误: {}", e),
    }
}

// 23. 重导出(re-export)
// 使用pub use可以将内部实现重新导出，简化外部API
mod internal {
    pub fn internal_function() {
        println!("内部函数");
    }
}

pub use crate::internal::internal_function; // 重导出内部函数

// 24. 包的结构
// Rust包可以包含：
// - 零个或一个库箱（src/lib.rs）
// - 任意数量的二进制箱（src/bin/*.rs）
// - 测试目录（tests/）
// - 示例目录（examples/）
// - 基准目录（benches/）

// 25. 工作区(Workspace)
// 工作区用于管理多个相关的包
// [workspace]
// members = [
//     "package1",
//     "package2",
//     "package3",
// ]

// 26. 模块系统最佳实践
// 1. 使用模块组织相关代码
// 2. 遵循单一职责原则
// 3. 合理使用可见性控制
// 4. 使用use导入常用路径，避免重复
// 5. 为公共API提供清晰的文档
// 6. 使用重导出简化外部API

// 27. 示例：一个完整的模块系统
// 以下是一个简化的博客系统的模块结构：

mod blog {
    pub mod post {
        pub struct Post {
            pub title: String,
            pub content: String,
            pub author: String,
        }
        
        impl Post {
            pub fn new(title: &str, content: &str, author: &str) -> Self {
                Self {
                    title: title.to_string(),
                    content: content.to_string(),
                    author: author.to_string(),
                }
            }
        }
    }
    
    pub mod comment {
        pub struct Comment {
            pub content: String,
            pub author: String,
        }
        
        impl Comment {
            pub fn new(content: &str, author: &str) -> Self {
                Self {
                    content: content.to_string(),
                    author: author.to_string(),
                }
            }
        }
    }
    
    pub mod utils {
        pub fn format_post(post: &crate::blog::post::Post) -> String {
            format!("Title: {}\nAuthor: {}\n\n{}", post.title, post.author, post.content)
        }
    }
}

fn blog_example() {
    let post = blog::post::Post::new(
        "Rust模块系统",
        "Rust的模块系统非常强大...",
        "Rust开发者"
    );
    
    let comment = blog::comment::Comment::new(
        "这是一篇很好的文章！",
        "读者"
    );
    
    let formatted_post = blog::utils::format_post(&post);
    println!("\n=== 博客文章 ===");
    println!("{}", formatted_post);
    println!("\n评论: {}", comment.content);
    println!("评论作者: {}", comment.author);
}
// 用于单独运行本文件的main函数
fn main() {
    run_example();
}
