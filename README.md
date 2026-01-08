# Rust 学习笔记

本项目是 Rust 语言的学习实践，涵盖 Web 服务开发、数据库操作、错误处理等核心知识点。

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

## 核心知识点

### 1. 所有权系统 (Ownership)

```rust
// & = 借来看一眼，所有权仍在原变量
// 没有 & = 拿走所有权，原变量作废

// Copy 类型 → 赋值=按位拷贝 → 原变量仍有效 → 不需要 &
// 非 Copy 类型 → 赋值=移动所有权 → 原变量作废 → 只读场景必须用 &
```

| 类型 | 直觉 | 要不要 & |
|------|------|----------|
| `i32/f64/bool/char/数组/元组` | 小整数类 | 不需要 |
| `String/Vec/HashMap/Box/自定义 struct` | 大对象类 | 必须 & |

**示例：**
```rust
let a = String::from("hello");  // a 拥有字符串
let b = a;                      // 移动所有权 → b 拥有，a 作废
// println!("{}", a);           // ❌ 编译错误：use of moved value
println!("{}", b);              // ✅ 只有 b 能用
```

---

### 2. Actix-Web 框架

#### 2.1 项目启动
```rust
#[actix_web::main]  // 把 async fn main 钉在 tokio 运行时上
async fn main() -> io::Result<()> {
    let share_data = web::Data::new(AppState { ... });

    let app = move || {
        App::new()
            .app_data(share_data.clone())
            .configure(general_routes)
            .configure(course_routes)
    };

    HttpServer::new(app).bind("127.0.0.1:3339")?.run().await
}
```

#### 2.2 路由配置
```rust
pub fn course_routes(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/courses")
            .route("/", web::post().to(new_course))
            .route("/{user_id}/{name}", web::get().to(get_courses_for_teacher)),
    );
}
```

#### 2.3 Handler 函数
```rust
pub async fn new_course(
    new_course: web::Json<Course>,   // 请求体自动反序列化
    app_state: web::Data<AppState>,  // 共享状态（Arc<AppState>）
) -> HttpResponse {
    // 处理逻辑...
    HttpResponse::Ok().json(course)
}
```

---

### 3. 共享状态与线程安全

```rust
use std::sync::Mutex;

pub struct AppState {
    pub health_check_response: String,    // 只读字段无需加锁
    pub visit_count: Mutex<u32>,          // 写字段需要 Mutex
    pub courses: Mutex<Vec<Course>>,      // 动态数组需要 Mutex
    pub db: PgPool,                       // 数据库连接池
}
```

**Mutex 使用：**
```rust
// 加锁 → 修改 → 自动解锁
let mut count = app_state.visit_count.lock().unwrap();
*count += 1;
// guard 离开作用域 → 自动解锁
```

**Rust vs PHP 内存模型对比：**
- Rust：进程内存，重启前数据一直在，请求间天然共享
- PHP：每个请求重新 new 一份，请求结束就清空
- 持久化方案：都落盘（DB/Redis）

---

### 4. 数据库操作 (SQLx)

#### 4.1 连接池
```rust
use sqlx::postgres::PgPoolOptions;

let db_pool = PgPoolOptions::new()
    .connect(&database_url)  // &str → 借用，不拷贝
    .await                   // 异步等待 TCP + TLS 握手
    .unwrap();
```

#### 4.2 查询操作
```rust
// 编译期检查 SQL（sqlx::query! 宏）
let rows = sqlx::query!(
    r#"SELECT * FROM rust_test1.course WHERE teacher_id = $1"#,
    teacher_id  // 占位符 $1 → PostgreSQL 风格
)
.fetch_all(&pool)  // &Pool → 借用池，不转移所有权
.await
.unwrap();
```

#### 4.3 插入并返回
```rust
let row = sqlx::query!(
    r#"INSERT INTO rust_test1.course (teacher_id, name) VALUES ($1, $2) RETURNING *"#,
    new_course.teacher_id,
    new_course.name
)
.fetch_one(pool)  // RETURNING * → 返回刚插入的行
.await
.unwrap();
```

---

### 5. 错误处理

#### 5.1 基础错误处理
```rust
// 方式一：match 手动处理
fn square(val: &str) -> Result<i32, ParseIntError> {
    match val.parse::<i32>() {
        Ok(num) => Ok(num.pow(2)),
        Err(e) => Err(e),
    }
}

// 方式二：? 运算符（推荐）
fn square(val: &str) -> Result<i32, ParseIntError> {
    let num = val.parse::<i32>()?;  // 失败自动返回 Err
    Ok(num.pow(2))
}
```

#### 5.2 自定义错误类型
```rust
#[derive(Debug, Serialize, Deserialize)]
pub enum MyError {
    DbError(String),      // 数据库错误
    ActixError(String),   // 框架错误
    NotFound(String),     // 资源未找到
}

// 实现 Display trait
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
```

#### 5.3 From trait 自动转换
```rust
// 把 SQLx 错误自动转成 MyError
impl From<sqlx::Error> for MyError {
    fn from(err: sqlx::Error) -> Self {
        MyError::DbError(err.to_string())
    }
}

// 使用时 ? 自动调用 From
pub async fn get_courses() -> Result<Vec<Course>, MyError> {
    let rows = sqlx::query!(...).fetch_all(pool).await?;  // 自动转换
    Ok(courses)
}
```

---

### 6. 序列化与反序列化 (Serde)

```rust
use serde::{Deserialize, Serialize};

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
// 自动反序列化请求体
async fn handler(payload: web::Json<Course>) -> HttpResponse {
    let course: Course = payload.into_inner();  // 取出内部值
    HttpResponse::Ok().json(course)             // 自动序列化响应
}
```

---

### 7. 属性宏 (#[...])

| 场景 | 属性 | 效果 |
|------|------|------|
| 入口函数 | `#[actix_web::main]` | 把 async fn main 钉在 tokio 运行时 |
| 测试 | `#[test]` | 告诉 cargo test 这是单元测试 |
| 派生 trait | `#[derive(Debug, Clone)]` | 自动生成实现 |
| 条件编译 | `#[cfg(target_os = "macos")]` | 只在 macOS 下编译 |
| 序列化 | `#[serde(rename = "user_name")]` | 重命名字段 |
| 手动指定目录 | `#[path = "../handlers.rs"]` | 指定模块文件位置 |

---

### 8. 单元测试

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

### 9. 常用类型速查

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

### 10. 关键概念总结

1. **零成本抽象**：Rust 的抽象不会带来运行时开销
2. **借用规则**：
   - 同一时刻，只能有一个可变引用 **或** 多个不可变引用
   - 引用必须始终有效
3. **生命周期**：编译器自动推断，复杂场景需要显式标注
4. **异步编程**：`.await` 挂起当前任务，不阻塞线程
5. **错误传播**：`?` 运算符 + `From` trait 实现优雅错误处理

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

## 环境配置

在项目根目录创建 `.env` 文件：
```
DATABASE_URL=postgres://user:password@localhost:5432/database
```

---

## 学习资源

- [The Rust Book](https://doc.rust-lang.org/book/)
- [Actix-Web 文档](https://actix.rs/docs/)
- [SQLx 文档](https://docs.rs/sqlx/latest/sqlx/)
- [Serde 文档](https://serde.rs/)
