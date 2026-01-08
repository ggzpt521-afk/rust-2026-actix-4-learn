use actix_web::{HttpResponse, Result, error, http::StatusCode};
use serde::{Deserialize, Serialize};
use std::fmt;

// ========== 1. 自定义错误枚举（可序列化 + Debug） ==========
#[derive(Debug, Serialize, Deserialize)]
pub enum MyErrorNew {
    DbError(String),    // 数据库错误
    ActixError(String), // 框架错误
    NotFound(String),   // 资源未找到
}

// ========== 2. HTTP 响应结构体（可序列化） ==========
#[derive(Debug, Serialize, Deserialize)]
pub struct MyErrorNewResponse {
    error_message: String, // 人类可读的错误信息
}

// ========== 3. impl MyErrorNew → 自定义方法 ==========
impl MyErrorNew {
    // 3.1 &self = “把当前错误借给你看一眼” → 零成本只读借用
    fn error_response(&self) -> String {
        match self {
            MyErrorNew::DbError(msg) => format!("数据库错误: {}", msg),
            MyErrorNew::ActixError(msg) => format!("框架错误: {}", msg),
            MyErrorNew::NotFound(msg) => format!("资源未找到: {}", msg),
        }
    }
}

// ========== 4. impl Display → 人类可读字符串 ==========
impl fmt::Display for MyErrorNew {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        // 直接打印枚举本身（Debug 已足够）
        write!(f, "{:?}", self)
    }
}

// ========== 5. impl ResponseError → Actix-Web 自动转换 ==========
impl actix_web::error::ResponseError for MyErrorNew {
    // 5.1 &self = “借当前错误看一眼” → 返回 HTTP 状态码
    fn status_code(&self) -> StatusCode {
        match self {
            MyErrorNew::DbError(_) => StatusCode::INTERNAL_SERVER_ERROR, // 500
            MyErrorNew::ActixError(_) => StatusCode::INTERNAL_SERVER_ERROR, // 500
            MyErrorNew::NotFound(_) => StatusCode::NOT_FOUND,            // 404
        }
    }

    // 5.2 &self = “借当前错误看一眼” → 返回 JSON 响应体
    fn error_response(&self) -> HttpResponse {
        let resp = MyErrorNewResponse {
            error_message: self.error_response(), // 调用 3.1 的人类可读信息
        };
        // build(status_code()) + json() → 返回 JSON + 状态码
        HttpResponse::build(self.status_code()).json(resp)
    }
}

// ========== 6. 把 Actix 错误自动转成 MyErrorNew ==========
impl From<actix_web::error::Error> for MyErrorNew {
    // 1.1 &self 不存在 → 因为 From 是 **关联函数**，不是方法
    // 1.2 from(err) → 输入一个 Actix 错误，输出一个 MyErrorNew
    fn from(err: actix_web::error::Error) -> Self {
        // 1.3 **零成本转换** → 只拷字符串，不移动原错误
        MyErrorNew::ActixError(err.to_string())
    }
}

// ========== 7. 把 SQLx 错误自动转成 MyErrorNew ==========
impl From<sqlx::Error> for MyErrorNew {
    // 2.1 from(err) → 输入一个 SQLx 错误，输出一个 MyErrorNew
    fn from(err: sqlx::Error) -> Self {
        // 2.2 **零成本转换** → 只拷字符串，不移动原错误
        MyErrorNew::DbError(err.to_string())
    }
}

// ========== 8. 一键使用（? 运算符自动转换） ==========
// pub async fn demo() -> Result<String, MyErrorNew> {
//     // 6.1 ? 运算符：如果 Err → 自动转成 MyErrorNew 并提前返回
//     // 这里模拟数据库错误
//     Err(MyErrorNew::DbError("连接超时".to_string()))
// }

//“fn = 造函数；impl = 把函数（或 trait）装到类型上。”
//From trait 就是 “零成本类型转换器”——
//输入 A，输出 B，不移动原对象，编译器自动调用。
