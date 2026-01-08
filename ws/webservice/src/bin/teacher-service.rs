// 引入 Actix Web 框架的核心组件：
// - `web`：用于处理请求参数、共享状态（Data）、路径配置等；
// - `App`：代表一个 Web 应用实例；
// - `HttpServer`：用于创建并运行 HTTP 服务器。
use actix_web::{web, App, HttpServer};

// 引入标准库的 I/O 模块，用于处理如端口绑定失败等 I/O 错误。
use std::io;

// 引入标准库的互斥锁 Mutex，用于在多线程环境中安全地修改共享数据（如访问计数）。
use std::sync::Mutex;

// 手动指定模块文件路径（不推荐常规使用，但可用于特殊项目结构）：
// 将上一级目录中的 `handlers.rs` 文件作为本地模块 `handlers` 引入。
//一句话记忆
//#[] = “给代码贴便签”
//便签内容 = 指令/配置/标记，谁看到谁处理。
//| 场景           | 属性                                | 效果                                  |
//| ------------ | --------------------------------- | ----------------------------------- |
//| **入口函数**     | `#[actix_web::main]`              | 把普通 `async fn main()` 钉在 tokio 运行时上 |
//| **测试**       | `#[test]`                         | 告诉 `cargo test` 这是单元测试              |
//| **派生 trait** | `#[derive(Debug, Clone)]`         | 自动生成 `Debug` 和 `Clone` 实现           |
//| **条件编译**     | `#[cfg(target_os = "macos")]`     | 只在 macOS 下编译这段代码                    |
//| **文档**       | `#[doc = "health check handler"]` | 给函数加文档注释                            |
//| **序列化**      | `#[serde(rename = "user_name")]`  | 让 serde 把字段序列化成 `user_name`         |
//| **性能**       | `#[inline]`                       | 建议编译器做内联展开                          |
//| **FFI**      | `#[no_mangle]`                    | 禁止改名，方便 C 语言链接                      |
//| **递归宏**      | `#[recursion_limit = "256"]`      | 提高宏展开深度上限                           |
//| **手动指定目录**  | `#[path = "..."]`                | 手动指定模块文件位置                         |


#[path = "../handlers.rs"]
mod handlers;

// 将上一级目录中的 `routers.rs` 文件作为本地模块 `routers` 引入。
#[path = "../routers.rs"]
mod routers;

// 将上一级目录中的 `state.rs` 文件作为本地模块 `state` 引入，
// 该文件中应定义了应用的全局状态结构体 `AppState`。
#[path = "../state.rs"]
mod state;

#[path = "../models.rs"]
mod models;

// 从 `routers` 模块中导入所有公开项（通常是路由配置函数，如 `general_routes`）。
use routers::*;

// 从 `state` 模块中导入 `AppState` 类型，用于构建应用的共享状态。
use state::AppState;

// `#[actix_web::main]` 是 Actix Web 提供的宏，用于将 `async fn main` 转换为
// 基于 Tokio 异步运行时的入口点。没有它，Rust 不允许 `main` 函数是异步的。
#[actix_web::main]
async fn main() -> io::Result<()> {
    // 创建应用的全局共享状态实例，并用 `web::Data::new()` 包装。
    // `web::Data<T>` 是 Actix Web 提供的线程安全共享容器（内部基于 Arc），
    // 允许多个 handler 安全地读取或修改该状态。
    let share_data = web::Data::new(
        AppState {
            // 初始化健康检查响应内容为字符串 "I'm OK"
            health_check_response: "I'm OK".to_string(),
            // 初始化访问计数器为 0，并用 Mutex 包裹以支持多线程安全修改
            // ⚠️ 注意：此处字段名必须与 `state.rs` 中定义的完全一致（建议拼写为 visit_count）
            visit_count: Mutex::new(0),
            //let v1 = vec![];        // 宏展开 = Vec::new() 一样快
            //let v2 = Vec::new();    // 直接空 Vec
            //Rust 里根本没有 vec[] 这种写法，只有vec![] 和 Vec::new()
            courses: Mutex::new(vec![])
        }
    );

    // 定义一个闭包 `app`，用于生成新的 `App` 实例。
    // 使用 `move ||` 表示该闭包“获取”外部变量 `share_data` 的所有权。
    // 因为服务器可能启动多个 worker 线程，每个线程都会调用此闭包一次，
    // 所以需要能多次克隆 `share_data`（`web::Data` 实现了 Clone）。
    let app = move || {
        App::new()
            // 将共享状态 `share_data` 注册到应用中，使所有 handler 都能通过参数注入访问它
            .app_data(share_data.clone())
            // 调用 `general_routes` 函数来批量注册路由（该函数应在 `routers.rs` 中定义）
            .configure(general_routes)
            .configure(course_routes)
    };

    // 启动 HTTP 服务器：
    // 1. `HttpServer::new(app)`：传入上面定义的应用工厂闭包；
    // 2. `.bind("127.0.0.1:3339")?`：尝试绑定到本地 3339 端口，若失败则返回错误（`?` 传播）；
    // 3. `.run().await`：异步启动服务器并阻塞等待其结束（通常直到 Ctrl+C 终止）。
    HttpServer::new(app).bind("127.0.0.1:3339")?.run().await
}