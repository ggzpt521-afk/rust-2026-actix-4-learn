// 引入标准库中的 `Mutex` 类型。 /mju:teks/
// `Mutex`（互斥锁）是一种用于在多线程环境中安全地共享和修改数据的同步原语。
// 它确保同一时间只有一个线程可以访问被它保护的数据，从而避免数据竞争（data race）。
use std::sync::Mutex;
use super::models::Course;  //需要在 teacher-service.rs 声明下mod 这里才能调用 否则报错
use sqlx::postgres::PgPool;

// 使用 `pub` 关键字声明一个公共的结构体 `AppState`。
// `pub` 表示这个结构体可以在当前模块之外被其他模块或 crate 访问。
// `AppState` 通常用于 Web 应用（如使用 Axum、Actix-web 等框架）中存储应用的全局状态。
pub struct AppState {
    // 字段 `health_check_response` 是一个 `String` 类型，且是公共的（`pub`）。
    // 它可能用于存储健康检查接口返回的固定响应内容，例如 "OK"。
    // 因为 `String` 是不可变的（除非显式可变），如果这个字段不需要在运行时修改，
    // 就不需要用 `Mutex` 包裹；但如果将来需要修改，可能也需要加锁。
    pub health_check_response: String,

    // 字段 `visit_count`（注意：拼写应为 `visit_count`，可能是笔误）是一个 `Mutex<u32>` 类型。
    // `Mutex<u32>` 表示一个被互斥锁保护的 32 位无符号整数。
    // 这个字段用于记录访问次数（比如网页被访问了多少次）。
    // 由于多个线程（例如处理 HTTP 请求的线程）可能会同时读写这个计数器，
    // 必须使用 `Mutex` 来保证线程安全。
    //
    // 注意：
    // - `Mutex<T>` 本身不是 `Send` 或 `Sync` 的，但 `std::sync::Mutex<T>` 是 `Send + Sync` 的，
    //   所以它可以安全地在线程间共享（前提是 T 也是 Send + Sync）。
    // - `u32` 是基本类型，满足这些要求，因此 `Mutex<u32>` 可以安全地放在共享状态中。
    pub visit_count: Mutex<u32>,

    // 就是 “一个带锁的公共课程列表”——
    // Vec<Course> 是 真正的数据；Mutex 是 看门的大锁；pub 表示 谁都看得见；
    //| 片段            | 含义                            |
    //| ------------- | ----------------------------- |
    //| `pub`         | 字段公开，**外部模块可读可写**             |
    //| `courses`      | 字段名，**课程列表**                  |
    //| `Mutex<...>`  | **互斥锁**，**同一时刻只允许一个线程访问内部数据** |
    //| `Vec<Course>` | **动态数组**，里面存 **Course 结构体实例** |
    pub courses: Mutex<Vec<Course>>,

    pub db: PgPool
}
