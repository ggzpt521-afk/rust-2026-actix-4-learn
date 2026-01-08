// 从父模块（通常是 main.rs 或 lib.rs 所在的上一级）导入所有公开的 handler 函数。
// 这里假设 `health_check_handler` 在 `handlers.rs` 中被定义并标记为 `pub`。
use super::handlers::*;

// 引入 Actix Web 的 `web` 模块，用于访问路由构建器（如 `web::get`, `web::post` 等）。
use actix_web::web;

// 定义一个公共函数 `general_routes`，用于集中配置应用的路由。
// 参数 `cfg: &mut web::ServiceConfig` 是 Actix Web 提供的路由配置上下文，
// 允许我们在其中注册多个路由。
pub fn general_routes(cfg: &mut web::ServiceConfig) {
    // 注册一个 GET 路由：
    // - 路径为 "/health"
    // - 使用 `web::get()` 创建一个 GET 请求处理器
    // - 通过 `.to(health_check_handler)` 绑定具体的处理函数
    // 注意：`health_check_handler` 必须是一个符合 Actix Web handler 签名的异步函数
    cfg.service(web::resource("/health").route(web::get().to(health_check_handler)));
}

// 引入 Actix Web 的 `web::ServiceConfig` 类型（通常已在上级模块引入，此处仅为上下文说明）
// 本函数用于集中注册与“课程（Course）”相关的所有 API 路由。

/// 注册课程相关路由的配置函数。
/// 
/// 原理说明：
/// - Actix Web 使用“服务配置（ServiceConfig）”模式来组织路由，支持模块化、嵌套和作用域隔离。
/// - `web::scope()` 允许将一组路由挂载到公共路径前缀下（如 `/courses`），避免重复书写前缀。
/// - 所有路由在应用启动时被收集并编译为高效的路由匹配表（基于 radix tree），性能极高。
//web::scope("/courses") —— 路由作用域
//作用：为一组路由添加公共路径前缀，提升代码组织性和可维护性。
//优势：
//避免在每个 .route() 中重复写 /courses
//可在 scope 上统一添加中间件（如认证、日志），例如：

//web::scope("/courses")
//    .wrap(AuthMiddleware) // 所有 /courses/* 路由都需认证
//    .route(...)
pub fn course_routes(cfg: &mut web::ServiceConfig) {
    
    // 向全局路由配置 `cfg` 中注册一个“作用域服务（scoped service）”。
    // `cfg.service(...)` 是注册子路由的标准方式，支持嵌套、中间件和生命周期管理。
    cfg.service(
        // 创建一个路由作用域（scope），所有子路由自动继承前缀 `/courses`
        // 例如：`.route("/", ...)` 实际对应完整路径 `/courses/`
        //       `.route("/{user_id}", ...)` 对应 `/courses/{user_id}`
        web::scope("/courses")                        
            
            // 注册 POST /courses 路由
            // - 路径：`/`（相对于 scope 前缀，即完整路径为 `/courses/`）
            // - HTTP 方法：POST（通过 `web::post()` 指定）
            // - 处理函数：`new_course`（必须是一个符合 Actix Web handler 签名的异步函数）
            //   通常用于创建新课程，请求体为 JSON 格式的 Course 数据
            .route("/", web::post().to(new_course))  
            
            // 注册 GET /courses/{user_id} 路由
            // - 路径：`/{user_id}`（完整路径为 `/courses/{user_id}`）
            // - HTTP 方法：GET（通过 `web::get()` 指定）
            // - 路径参数：`{user_id}` 会被自动提取，并传递给 handler（如通过 `web::Path<usize>`）
            // - 处理函数：`get_courses_for_teacher`，用于根据教师 ID 查询其所有课程
            .route("/{user_id}/{name}", web::get().to(get_courses_for_teacher)),
    );
}