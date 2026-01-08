# Rust 学习笔记

本项目是 Rust 语言的学习实践，涵盖基础语法、Web 服务开发、数据库操作、错误处理等核心知识点。

---

## 目录

- [环境安装](#环境安装)
- [项目结构](#项目结构)
- [基础语法](#基础语法)
- [所有权系统](#所有权系统)
- [结构体与方法](#结构体与方法)
- [枚举与模式匹配](#枚举与模式匹配)
- [常见集合](#常见集合)
- [泛型与 Trait](#泛型与-trait)
- [错误处理](#错误处理)
- [Actix-Web 框架](#actix-web-框架)
- [数据库操作](#数据库操作-sqlx)
- [常用类型速查](#常用类型速查)
- [常用库发音](#常用库发音)

---

## 环境安装

```bash
# 1. 安装 Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# 2. 配置环境变量
. "$HOME/.cargo/env"            # For sh/bash/zsh/ash/dash/pdksh
source "$HOME/.cargo/env.fish"  # For fish

# 3. 直接运行 rs 文件
rustc main.rs
./main

# 4. 使用 Cargo 创建项目（推荐）
cargo new project_name
cd project_name
cargo run
```

### SQLx CLI 安装（数据库开发必需）

```bash
# 安装全局命令
cargo install sqlx-cli

# 生成类型缓存（让 IDE 和编译器知道字段类型）
cargo sqlx prepare
```

> `cargo sqlx prepare` 就是 "提前连数据库，把字段类型写进 `.sqlx/` 缓存"——让 IDE 和编译器在没有 `DATABASE_URL` 时也能知道 `r.id`、`r.name` 的类型。
>
> **工作流**："改表 → `cargo sqlx prepare` → 改 Rust 结构体 → `cargo check` 全绿 → 完成"

---

## 项目结构

```
learn-one/
├── ws/                    # 基础 Web 服务（内存存储）
├── db/                    # 数据库连接学习
├── ws-db/                 # Web 服务 + 数据库集成
├── ws-db-with-error/      # Web 服务 + 数据库 + 错误处理
└── error/                 # 错误处理基础
```

---

## 基础语法

### 变量与可变性

```rust
let x = 5;           // 不可变变量
let mut x = 5;       // 可变变量（加 mut）
x = 6;               // 可以修改
```

### 标量类型

| 类型 | 说明 | 示例 |
|------|------|------|
| 整数 | i8/i16/i32/i64/i128/isize, u8/u16/u32/u64/u128/usize | `let x: i32 = 42;` |
| 浮点 | f32, f64 | `let y: f64 = 3.14;` |
| 布尔 | bool | `let t: bool = true;` |
| 字符 | char（4字节 Unicode） | `let c: char = '中';` |

### 复合类型

```rust
// 元组 tuple - 固定长度，可不同类型
let tup: (i32, f64, u8) = (500, 6.4, 1);
let (x, y, z) = tup;           // 解构
let first = tup.0;             // 索引访问

// 数组 array - 固定长度，相同类型
let arr = [1, 2, 3, 4, 5];
let first = arr[0];
```

### 代码块（表达式）

```rust
let y = {
    let x = 1;
    x + 1    // 注意：没有分号 → 返回值
};
println!("y = {}", y);  // y = 2
```

### 函数

```rust
// 带返回值的函数
fn five() -> i32 {
    5    // 没有分号 → 返回值
}

fn add_one(x: i32) -> i32 {
    x + 1    // 没有分号 → 返回值
}
```

### 分号规则（重要）

```
Rust 的分号就是 "表达式尾巴开关"：
- 尾行不要分号 → 把值返回出去（相当于 return）
- 其他位置要分号 → 丢掉值，继续下一句
```

### 双冒号 :: 说明

```
双冒号 :: 在 Rust 里不是"调用方法"，而是路径（namespace）分隔符——
它告诉编译器："后面这个东西位于哪个模块/结构体/枚举/trait 里"

示例：
- String::from("hello")  → String 类型的 from 关联函数
- std::io::Result        → std 模块下 io 子模块的 Result 类型
```

---

## 控制流

### if 表达式

```rust
let z = 15;
if z > 10 {           // 注意：没有括号
    println!("大");
} else {
    println!("小");
}
```

### 循环

```rust
// loop - 无限循环
loop {
    println!("again!");
    break;  // 用 break 退出
}

// loop 返回值
let result = loop {
    count += 1;
    if count == 10 {
        break count * 2;  // break 可以返回值
    }
};

// while 循环
let mut n = 3;
while n != 0 {
    println!("{}!", n);
    n -= 1;
}

// for 循环（推荐）
let a = [1, 2, 3, 4, 5];
for i in a {
    println!("i = {}", i);
}

// 范围迭代 + 反转
for number in (1..4).rev() {
    println!("{}!", number);  // 3! 2! 1!
}
```

---

## 所有权系统

### 核心规则

```rust
// & = 借来看一眼，所有权仍在原变量
// 没有 & = 拿走所有权，原变量作废

// Copy 类型 → 赋值=按位拷贝 → 原变量仍有效 → 不需要 &
// 非 Copy 类型 → 赋值=移动所有权 → 原变量作废 → 只读场景必须用 &
```

| 类型 | 内存布局 | 赋值行为 | 原变量 |
|------|----------|----------|--------|
| `i32/f64/bool/char` | 值在栈上 | memcpy（1指令） | 仍有效 |
| `String/Vec/HashMap` | 栈指针+堆数据 | 移动指针 | 作废 |

**示例：**
```rust
let a = String::from("hello");  // a 拥有字符串
let b = a;                      // 移动所有权 → b 拥有，a 作废
// println!("{}", a);           // ❌ 编译错误：use of moved value
println!("{}", b);              // ✅ 只有 b 能用
```

### 作用域与自动回收

```rust
{
    let s = String::from("hello");  // s 开始有效
    // 使用 s
}   // 作用域结束，s 自动回收（调用 drop）
```

### 切片 Slice

```rust
let s = String::from("hello world");
let hello = &s[0..5];   // "hello"
let world = &s[6..11];  // "world"
```

---

## 结构体与方法

### 定义结构体

```rust
struct User {
    email: String,
    username: String,
    active: bool,
    sign_in_count: u64,
}

let user1 = User {
    email: String::from("test@example.com"),
    username: String::from("test"),
    active: true,
    sign_in_count: 1,
};
```

### Debug 打印

```rust
#[derive(Debug)]  // 必须加这个
struct Rectangle {
    width: u32,
    height: u32,
}

let rect = Rectangle { width: 30, height: 50 };
println!("rect is {:?}", rect);   // 必须用 {:?}
```

### 定义方法（impl）

```rust
#[derive(Debug)]
struct Rectangle {
    width: u32,
    height: u32,
}

impl Rectangle {
    // 方法：第一个参数是 &self
    fn area(&self) -> u32 {
        self.width * self.height
    }

    // 关联函数：没有 self 参数（类似静态方法）
    fn square(size: u32) -> Rectangle {
        Rectangle { width: size, height: size }
    }
}

// 使用
let rect = Rectangle { width: 30, height: 50 };
println!("area = {}", rect.area());

let sq = Rectangle::square(10);  // 关联函数用 :: 调用
```

---

## 枚举与模式匹配

### 定义枚举

```rust
enum IpAddrKind {
    V4,
    V6,
}

// 带数据的枚举
enum Coin {
    Penny,
    Nickel,
    Dime,
    Quarter(String),  // 变体可以带数据
}
```

### match 表达式

```rust
fn value_in_cents(coin: Coin) -> u8 {
    match coin {
        Coin::Penny => {
            println!("Lucky penny!");
            1
        }
        Coin::Nickel => 5,
        Coin::Dime => 10,
        Coin::Quarter(state) => {
            println!("State: {}", state);
            25
        }
    }
}
```

### if let 简化

```rust
let coin = Coin::Quarter("Alaska".to_string());
let mut count = 0;

// 只关心一种情况时用 if let
if let Coin::Quarter(state) = coin {
    println!("State quarter from {}!", state);
} else {
    count += 1;
}
```

---

## 常见集合

### Vector

```rust
// 创建
let v: Vec<i32> = Vec::new();
let v = vec![1, 2, 3];  // 宏创建

// 操作
let mut v = Vec::new();
v.push(5);
v.push(6);

// 访问
let third = &v[2];           // 可能 panic
let third = v.get(2);        // 返回 Option<&T>
```

### HashMap

```rust
use std::collections::HashMap;

let mut map = HashMap::new();
map.insert("key1", 10);
map.insert("key2", 20);

// 统计单词频率
let text = "hello world wonderful world";
let mut map = HashMap::new();
for word in text.split_whitespace() {
    let count = map.entry(word).or_insert(0);
    *count += 1;
}
println!("{:?}", map);  // {"hello": 1, "world": 2, "wonderful": 1}
```

---

## 泛型与 Trait

### 泛型

```rust
// 泛型结构体
struct Point<T> {
    x: T,
    y: T,
}

// 泛型方法
impl<T> Point<T> {
    fn x(&self) -> &T {
        &self.x
    }
}

let p = Point { x: 5, y: 10 };
println!("p.x = {}", p.x());
```

### Trait（类似接口）

```rust
// 定义 trait
pub trait Summary {
    fn summarize(&self) -> String;
}

// 为类型实现 trait
pub struct NewsArticle {
    pub headline: String,
    pub author: String,
}

impl Summary for NewsArticle {
    fn summarize(&self) -> String {
        format!("{}, by {}", self.headline, self.author)
    }
}

pub struct Tweet {
    pub username: String,
    pub content: String,
}

impl Summary for Tweet {
    fn summarize(&self) -> String {
        format!("{}: {}", self.username, self.content)
    }
}
```

---

## 错误处理

### panic! 宏

```rust
fn main() {
    panic!("crash and burn");  // 程序立即终止
}
```

### Result 类型

```rust
use std::fs::File;
use std::io::{self, Read};

// 方式一：match 手动处理
fn read_file() -> Result<String, io::Error> {
    let mut f = match File::open("hello.txt") {
        Ok(file) => file,
        Err(e) => return Err(e),
    };
    let mut s = String::new();
    match f.read_to_string(&mut s) {
        Ok(_) => Ok(s),
        Err(e) => Err(e),
    }
}

// 方式二：? 运算符（推荐）
fn read_file_short() -> Result<String, io::Error> {
    let mut f = File::open("hello.txt")?;  // 失败自动返回 Err
    let mut s = String::new();
    f.read_to_string(&mut s)?;
    Ok(s)
}
```

### unwrap 和 expect

```rust
// unwrap：成功返回值，失败 panic
let f = File::open("hello.txt").unwrap();

// expect：失败时带自定义错误信息
let f = File::open("hello.txt").expect("Failed to open hello.txt");
```

### 自定义错误类型

```rust
use actix_web::{HttpResponse, http::StatusCode};
use serde::{Serialize, Deserialize};
use std::fmt;

#[derive(Debug, Serialize, Deserialize)]
pub enum MyError {
    DbError(String),
    ActixError(String),
    NotFound(String),
}

// 实现 Display
impl fmt::Display for MyError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{:?}", self)
    }
}

// 实现 ResponseError → Actix-Web 自动转换
impl actix_web::error::ResponseError for MyError {
    fn status_code(&self) -> StatusCode {
        match self {
            MyError::DbError(_) => StatusCode::INTERNAL_SERVER_ERROR,
            MyError::NotFound(_) => StatusCode::NOT_FOUND,
            _ => StatusCode::INTERNAL_SERVER_ERROR,
        }
    }
}

// From trait：自动类型转换
impl From<sqlx::Error> for MyError {
    fn from(err: sqlx::Error) -> Self {
        MyError::DbError(err.to_string())
    }
}
```

> **From trait 就是 "零成本类型转换器"**——输入 A，输出 B，编译器自动调用。
>
> **fn = 造函数；impl = 把函数（或 trait）装到类型上。**

---

## Actix-Web 框架

### Tokio 异步运行时

```
Tokio 就是 "Rust 的异步运行时"——
把 async/await 变成真正能跑的多线程、零成本、无阻塞的并发程序；
一句话：没有 Tokio，Rust 的 async 只是语法糖。
```

### io::Result<()> 说明

```
() 是 Rust 的单元类型（unit type），代表"成功但无数据"
io::Result<()> 就是 "I/O 操作成功，没有额外返回值" 的标准写法
相当于其他语言的 void，但是一个真正的类型，可以放在泛型里
```

### 项目启动

```rust
use actix_web::{web, App, HttpServer, HttpResponse, Responder};
use std::io;
use std::sync::Mutex;

#[actix_web::main]  // 把 async fn main 钉在 tokio 运行时上
async fn main() -> io::Result<()> {
    let share_data = web::Data::new(AppState {
        health_check_response: "I'm OK".to_string(),
        visit_count: Mutex::new(0),
        courses: Mutex::new(vec![]),
    });

    // move 捕获环境，保证闭包 Send + 'static
    let app = move || {
        App::new()
            .app_data(share_data.clone())
            .configure(general_routes)
            .configure(course_routes)
    };

    HttpServer::new(app)
        .bind("127.0.0.1:3339")?
        .run()
        .await
}
```

### 路由配置

```rust
pub fn general_routes(cfg: &mut web::ServiceConfig) {
    cfg.route("/health", web::get().to(health_check_handler));
}

pub fn course_routes(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/courses")
            .route("/", web::post().to(new_course))
            .route("/{user_id}/{name}", web::get().to(get_courses_for_teacher)),
    );
}
```

### Handler 函数

```rust
// 健康检查
pub async fn health_check_handler() -> impl Responder {
    HttpResponse::Ok().json("httpserver is running")
}

// 带参数的 handler
pub async fn new_course(
    new_course: web::Json<Course>,   // 请求体自动反序列化
    app_state: web::Data<AppState>,  // 共享状态（Arc<AppState>）
) -> HttpResponse {
    // 处理逻辑...
    HttpResponse::Ok().json(course)
}
```

### 共享状态与线程安全

```rust
use std::sync::Mutex;
use sqlx::postgres::PgPool;

pub struct AppState {
    pub health_check_response: String,    // 只读字段无需加锁
    pub visit_count: Mutex<u32>,          // 写字段需要 Mutex
    pub courses: Mutex<Vec<Course>>,      // 动态数组需要 Mutex
    pub db: PgPool,                       // 数据库连接池
}

// Mutex 使用
let mut count = app_state.visit_count.lock().unwrap();
*count += 1;
// guard 离开作用域 → 自动解锁
```

**Rust vs PHP 内存模型对比：**
- **Rust**：进程内存，重启前数据一直在，请求间天然共享
- **PHP**：每个请求重新 new 一份，请求结束就清空
- **持久化方案**：都落盘（DB/Redis）

### Actix 并发模型

1. **异步 I/O**：单线程处理多连接
2. **多线程并行**：worker 数等于系统 CPU 数

---

## 数据库操作 (SQLx)

### 连接池

```rust
use sqlx::postgres::PgPoolOptions;
use dotenv::dotenv;
use std::env;

dotenv().ok();  // 加载 .env 文件
let database_url = env::var("DATABASE_URL").expect("DATABASE_URL not set");

let db_pool = PgPoolOptions::new()
    .connect(&database_url)  // &str → 借用，不拷贝
    .await                   // 异步等待 TCP + TLS 握手
    .unwrap();
```

### 查询操作

```rust
// 编译期检查 SQL（sqlx::query! 宏）
let rows = sqlx::query!(
    r#"SELECT * FROM rust_test1.course WHERE teacher_id = $1"#,
    teacher_id  // 占位符 $1 → PostgreSQL 风格
)
.fetch_all(&pool)  // &Pool → 借用池，不转移所有权
.await?;

// 映射到结构体
let courses: Vec<Course> = rows.iter()
    .map(|r| Course {
        id: r.id,
        teacher_id: r.teacher_id.unwrap_or(0),
        name: r.name.clone().unwrap_or_default(),
        time: r.time,
    })
    .collect();
```

### 插入并返回

```rust
let row = sqlx::query!(
    r#"INSERT INTO rust_test1.course (teacher_id, name) VALUES ($1, $2) RETURNING *"#,
    new_course.teacher_id,
    new_course.name
)
.fetch_one(pool)  // RETURNING * → 返回刚插入的行
.await?;
```

### 环境配置

在项目根目录创建 `.env` 文件：
```
DATABASE_URL=postgres://user:password@localhost:5432/database
TZ=Asia/Shanghai
```

---

## 属性宏 (#[...])

```
#[] = "给代码贴便签"
便签内容 = 指令/配置/标记，谁看到谁处理
```

| 场景 | 属性 | 效果 |
|------|------|------|
| 入口函数 | `#[actix_web::main]` | 把 async fn main 钉在 tokio 运行时 |
| 测试 | `#[test]` | 告诉 cargo test 这是单元测试 |
| 派生 trait | `#[derive(Debug, Clone)]` | 自动生成实现 |
| 条件编译 | `#[cfg(target_os = "macos")]` | 只在 macOS 下编译 |
| 序列化 | `#[serde(rename = "user_name")]` | 重命名字段 |
| 文档 | `#[doc = "说明"]` | 给函数加文档注释 |
| 性能 | `#[inline]` | 建议编译器做内联展开 |
| FFI | `#[no_mangle]` | 禁止改名，方便 C 语言链接 |
| 手动指定目录 | `#[path = "../handlers.rs"]` | 指定模块文件位置 |

---

## 序列化与反序列化 (Serde)

```rust
use serde::{Deserialize, Serialize};
use chrono::NaiveDateTime;

#[derive(Deserialize, Serialize, Debug, Clone)]
pub struct Course {
    pub teacher_id: i32,
    pub id: Option<i32>,           // Option → 兼容数据库 NULL
    pub name: String,
    pub time: Option<NaiveDateTime>,
}
```

**web::Json 使用：**
```rust
async fn handler(payload: web::Json<Course>) -> HttpResponse {
    let course: Course = payload.into_inner();  // 取出内部值
    HttpResponse::Ok().json(course)             // 自动序列化响应
}
```

---

## 单元测试

```rust
#[cfg(test)]
mod tests {
    use super::*;
    use actix_web::http::StatusCode;

    #[actix_web::test]
    async fn post_course_test() {
        // 构造请求体
        let course = web::Json(Course { ... });

        // 构造状态
        let app_state = web::Data::new(AppState { ... });

        // 直接调 handler（绕过 HTTP 层，速度最快）
        let resp = new_course(course, app_state).await;

        // 断言
        assert_eq!(resp.status(), StatusCode::OK);

        // 取出 body 验证
        let bytes = resp.into_body().try_into_bytes().unwrap();
        let body = std::str::from_utf8(&bytes).unwrap();
        assert_eq!(body, "course add");
    }
}
```

---

## 常用类型速查

| 类型 | 说明 |
|------|------|
| `String` | 堆上可变字符串 |
| `&str` | 字符串切片（借用） |
| `Vec<T>` | 动态数组 |
| `Option<T>` | 可空值（Some/None） |
| `Result<T, E>` | 错误处理（Ok/Err） |
| `Mutex<T>` | 互斥锁 |
| `Arc<T>` | 原子引用计数（线程安全共享） |
| `web::Data<T>` | Actix-Web 共享状态（内部是 Arc） |
| `web::Json<T>` | JSON 请求/响应包装器 |
| `web::Path<T>` | 路径参数提取器 |

---

## 关键概念总结

1. **零成本抽象**：Rust 的抽象不会带来运行时开销
2. **借用规则**：
   - 同一时刻，只能有一个可变引用 **或** 多个不可变引用
   - 引用必须始终有效
3. **生命周期**：编译器自动推断，复杂场景需要显式标注
4. **异步编程**：`.await` 挂起当前任务，不阻塞线程
5. **错误传播**：`?` 运算符 + `From` trait 实现优雅错误处理

---

## 常用库发音

| 库名 | 发音 | 含义 |
|------|------|------|
| actix | /ˈæktɪks/  | Web 框架 |
| chrono | /ˈkrɑːnoʊ/  | 时间日期 |
| dotenv | /dɑːt env/ | 环境变量 |
| serde | /ˈsɜːrdi/  | 序列化 |
| sqlx | /ˌes kjuː el ˈeks/ | 数据库 |
| tokio | /ˈtoʊkioʊ/ | 异步运行时 |

| 关键词 | 发音 | 含义 |
|--------|------|------|
| features | /ˈfiːtʃərz/ | 特征 |
| derive | /dɪˈraɪv/ | 派生 |
| schemas | /ˈskiːməz/ | 模式 |
| identity | /aɪˈdentəti/ | 自增 |
| naive | /naɪˈiːv/ | 原始的 |
| async | /əˈsɪŋk/ | 异步 |
| sync | /sɪŋk/ | 同步 |
| mutate | /ˈmjuːteɪt/ | 改变 |
| mutex | /ˈmjuːteks/ | 互斥锁 |

---

## 运行项目

```bash
# 运行 Web 服务
cd ws/webservice && cargo run --bin teacher-service

# 运行数据库测试
cd db && cargo run

# 运行单元测试
cargo test
```

---

## 学习资源

- [The Rust Book](https://rustwiki.org/zh-CN/book/title-page.html)
- [Actix-Web 文档](https://actix.rs/docs/)
- [SQLx 文档](https://docs.rs/sqlx/latest/sqlx/)
- [Serde 文档](https://serde.rs/)
- [Rust 中文学习文档](https://rustwiki.org/)
