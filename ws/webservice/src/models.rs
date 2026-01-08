// === 引入必要的依赖模块 ===

// Actix Web 的 `web` 模块：用于处理 HTTP 请求体（如 JSON）、路径参数等。
use actix_web::web;

// chrono 库中的 `NaiveDateTime`：表示不带时区的时间（格式如 2025-01-01 12:00:00），
// 常用于数据库存储或简单时间记录（注意：生产环境建议用带时区的 DateTime）。
use chrono::NaiveDateTime;

// serde 的核心 trait：
// - `Deserialize`：允许从 JSON 字符串反序列化为结构体（接收请求）
// - `Serialize`：允许将结构体序列化为 JSON 字符串（返回响应）
use serde::{Deserialize, Serialize};


// === 定义 Course 结构体 ===
//
// 使用 `#[derive(...)]` 自动实现多个常用 trait：
// - `Deserialize` / `Serialize`：支持与 JSON 互转（用于 API 输入/输出）
// - `Debug`：支持通过 `println!("{:?}", course)` 调试打印
// - `Clone`：允许复制整个结构体（因为所有字段都实现了 Clone）
//
// 设计说明：
// - `teacher_id` 是必填项（usize），表示所属教师
// - `id` 是可选项（Option<usize>），因为新建课程时数据库尚未分配 ID
// - `name` 是课程名称，必填（String）
// - `time` 是创建/更新时间，可为空（Option<NaiveDateTime>），兼容数据库 NULL
#[derive(Deserialize, Serialize, Debug, Clone)]
pub struct Course {
    pub teacher_id: usize,
    pub id: Option<usize>,
    pub name: String,
    pub time: Option<NaiveDateTime>,
}


// === 关于 From<web::Json<Course>> for Course 的说明 ===
//
// ❌ 原始错误写法（已注释掉）：
// impl From<web::Json<Course>> for Course {
//     fn from(course: web::Json<Course>) -> Self {
//         Course {
//             teacher_id: Course.teacher_id,   // ← 错误！Course 是类型名，不能这样访问字段
//             id: Course.id,                   // ← 同上，编译会报错
//             name: Course.name.clone(),       // ← 同上
//             time: Course.time                // ← 同上
//         }
//     }
// }
//
// 🔍 错误原因：
// - `Course.teacher_id` 不是合法表达式。Rust 中不能通过“类型名.字段”访问数据。
// - 正确做法是使用函数参数 `course`（小写变量名）来访问其内部字段。
// - 更重要的是：`web::Json<T>` 是一个包装器，已经提供了安全解包方法。
//
// ✅ 正确且推荐的做法：
// 实际上，**通常不需要手动实现这个 From trait**！
// 因为 Actix Web 的 `web::Json<T>` 已经内置了以下能力：
//   1. 在 handler 中直接作为参数：`fn handler(payload: web::Json<Course>)`
//   2. 通过 `.into_inner()` 方法获取内部的 `Course`
//   3. `web::Json<T>` 已经实现了 `Into<T>`，所以可以直接调用 `.into()`
//
// 因此，如果你真的需要 `From<web::Json<Course>> for Course`，
// 最简洁、安全的实现如下（但一般没必要写）：
impl From<web::Json<Course>> for Course {
    fn from(course: web::Json<Course>) -> Self {
        Course {
            teacher_id: course.teacher_id,   
            id: course.id,                   
            name: course.name.clone(),       
            time: course.time               
        }
    }
}

// === （可选）正确实现 From trait 的方式 ===
//
// 注意：此实现是冗余的，因为 `web::Json<Course>` 已经可以 `.into()` 转为 `Course`。
// 仅在特殊场景（如泛型约束要求必须有 From 实现）才需要。
// impl From<web::Json<Course>> for Course {
//     fn from(course: web::Json<Course>) -> Self {
//         // `web::Json<T>` 提供 `.into_inner()` 方法，安全地取出内部 T
//         // 这是最推荐的方式，语义清晰且零成本
//         course.into_inner()
//     }
// }


// === 最佳实践建议（无需额外代码）===
//
// 在你的 handler 函数中，直接这样使用即可：
//
// async fn create_course(payload: web::Json<Course>) -> impl Responder {
//     // 方式 1：使用 .into_inner()
//     let course: Course = payload.into_inner();
//
//     // 方式 2：利用已有的 Into 实现
//     // let course: Course = payload.into();
//
//     // 现在 course 是纯 Course 结构体，可存入数据库或处理
//     HttpResponse::Ok().json(course)
// }
//
// 因此，**本文件不需要任何 From 实现**，保持简洁即可。


// === 总结 ===
//
// - 结构体 `Course` 已正确配置 serde 和调试支持。
// - 字段设计合理，兼容数据库常见场景（ID 和时间可为空）。
// - 无需手动实现 `From<web::Json<Course>>`，Actix Web 已提供更优方案。
// - 避免重复造轮子，优先使用框架内置功能。