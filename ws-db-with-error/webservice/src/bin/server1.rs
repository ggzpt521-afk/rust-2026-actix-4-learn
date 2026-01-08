//双冒号 :: 在 Rust 里 不是“调用方法”，而是 路径（namespace）分隔符—— “后面这个东西位于哪个模块/结构体/枚举/ trait 里
// 引入 actix-web 核心部件；Responder 让异步函数可以直接当 HTTP 响应
use actix_web::{App, HttpResponse, HttpServer, Responder, web};
// Rust 标准库 I/O 错误类型，main 函数用它做错误载体
use std::io;

// ====== 路由模块 ======
// 把所有跟“通用/健康”相关的路由注册到 ServiceConfig
// 原理：actix 启动时会回调这个函数，把路径 + 处理器装进路由表
// 有分号 相当于 不需要返回值
pub fn general_routes(cfg: &mut web::ServiceConfig) {
    // 注册 GET /health  →  由 health_check_handler 异步函数处理
    cfg.route("/health", web::get().to(health_check_handler));
    cfg.route("/healthy", web::get().to(health_check_handler));
}

// ====== 处理器（Controller） ======
// 异步函数签名：返回 impl Responder → actix 能把它变成 HTTP 响应
// 功能：返回 200 OK + JSON 字符串，供负载均衡/监控探活
// 尾行不要分号 → 把值返回出去；    		相当于 return a;
pub async fn health_check_handler() -> impl Responder {
    // HttpResponse::Ok() 生成 200 状态；json() 自动序列化并设置 Content-Type: application/json
    HttpResponse::Ok().json("httpserver is running")
}

// ====== 入口：main ======
//顶级目录 执行 cargo run -p webservice --bin=server1 
//平级目录 webservice目录执行 cargo run --bin=server1
//运行起来之后执行 http://localhost:9919/health

// #[actix_web::main] 是宏，把 async main 钉在 tokio 运行时上（Rust 原生 main 不能 async）
//io::Result<()> 里的 () 不是“空指针”也不是“无返回值”，而是 Rust 的单元类型（unit type），只有一个值，也叫 ()。
//它代表 “业务成功，但没有额外数据要带回来”，相当于 C/Java 里的 void，但是一个真正的类型，可以放在泛型里、可以 Ok(())、可以 match，编译器能检查它，所以比 void 更安全
//() 是 Rust 的“成功但无数据”类型，0 字节、可泛型、编译期安全，
//io::Result<()> 就是 “I/O 操作成功，没有额外返回值” 的标准写法。
#[actix_web::main]
async fn main() -> io::Result<()> {
    // 构造 **应用工厂**：每次新连接，actix 会调用这个闭包生成独立的 App 实例
    // move 捕获空环境，保证闭包 Send + 'static
    let app = move || App::new().configure(general_routes);

    // HttpServer 是 tokio 上的异步 TCP 服务器；new(app) 把工厂传进去
    // .bind() 返回 Result，? 把绑定失败（端口被占等）向上抛
    // .run() 生成 *Server* future；await 让它一直监听，直到进程被杀
    HttpServer::new(app)
        .bind("127.0.0.1:9919")? // 监听本地 9919，成功则返回 Server，失败早抛
        .run() // 启动 tokio task，开始 accept 连接
        .await // 阻塞在这里，永不返回（除非出错）
}
