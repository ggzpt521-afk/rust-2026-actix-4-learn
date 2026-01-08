//Rust 把数据放在「进程内存」里，PHP 每个请求都「重新 new 一份」；
//所以 Rust 重启服务前数据一直在，PHP 请求结束就清空。
//Rust 全局状态 = 整个进程一块内存 → 请求之间天然共享
//PHP 默认 = 每个请求一份新内存 → 请求结束就清空
//想持久 → 都落盘（DB/Redis）；想共享 → 用进程外存储。
// ========== 1. 依赖与模块导入 ==========
use super::db_access::*;
use super::errors::MyErrorNew;
use super::state::AppState; // 全局共享状态（带锁的容器）
use crate::{ models::Course}; // 我们自己的课程结构体
use actix_web::body::MessageBody; //try_into_bytes 是 MessageBody 的方法 → 先 use actix_web::body::MessageBody; 再 .into_body().try_into_bytes()”
use actix_web::{HttpResponse, web}; // Web 框架核心类型
use chrono::Utc; // 时间戳生成器（UTC 时间）

// ========== 2. 健康检查 ==========
pub async fn health_check_handler(app_state: web::Data<AppState>) -> HttpResponse {
    // 2.1 只读字段无需加锁，直接引用
    let health_check_response = &app_state.health_check_response;

    // 2.2 计数器是 Mutex，必须加锁才能改；lock() 返回 MutexGuard<u32>
    //      unwrap() 在 poison 时 panic（测试可接受，生产建议 match）
    let mut visit_count = app_state.visit_count.lock().unwrap();

    // 2.3 拼接响应文本；format! 不会阻塞，因为只读字段无锁
    let response = format!("{}{} times", health_check_response, *visit_count);

    // 2.4 自增必须在 guard 作用域里，否则编译器不让改
    *visit_count += 1;
    // 2.5 guard 离开作用域 → 自动解锁，其他线程可继续读

    // 2.6 返回 JSON；&String 自动序列化成 JSON 字符串
    HttpResponse::Ok().json(&response)
}

// ========== 3. 新建课程 ==========
pub async fn new_course(
    new_course: web::Json<Course>,  // 3.1 请求体自动反序列化成 Course
    app_state: web::Data<AppState>, // 3.2 共享状态，内部是 Arc<AppState>
) -> HttpResponse {
    println!("Received new course");

    // 3.3 计算同一老师的已有课程数（用于生成自增 ID）
    //     clone() 会把整表复制一份 → O(n) 内存，测试可接受；
    //     生产环境建议 iter() + count()，避免整表克隆
    let course_count = app_state
        .courses
        .lock()
        .unwrap()
        .iter() // 只读迭代，无克隆
        .filter(|course| course.teacher_id == new_course.teacher_id)
        .count();

    // 3.4 构建新 Course；id 用 count+1 模拟自增，time 用当前 UTC
    let new_course = Course {
        teacher_id: new_course.teacher_id,
        id: 2,                              // 自增 ID
        name: new_course.name.clone(),      // 克隆字段，避免 move
        time: Some(Utc::now().naive_utc()), // 时间戳
    };

    // 3.5 再次加锁，把新课程 push 进 Vec
    app_state.courses.lock().unwrap().push(new_course);

    // 3.6 返回简单文本
    HttpResponse::Ok().body("course add")
}

pub async fn new_course_handle_db(
    new_course: web::Json<Course>,  // 3.1 请求体自动反序列化成 Course
    app_state: web::Data<AppState>, // 3.2 共享状态，内部是 Arc<AppState>
) -> HttpResponse {
    println!("Received new course");

    let course = post_new_course_db(&app_state.db, new_course.into()).await;
    HttpResponse::Ok().json(course)
}
// ========== 4. 根据老师 ID 查课程 ==========
pub async fn get_courses_for_teacher(
    app_state: web::Data<AppState>,
    params: web::Path<(i32, String)>, // 4.1 路径参数：/courses/{teacher_id}/{name}
) -> HttpResponse {
    // 4.2 解压元组 → (usize, String)
    let (teacher_id, _name) = params.into_inner();

    // 4.3 只读过滤：iter() 不克隆，filter 后 cloned() 把匹配项复制出来
    let filtered_courses = app_state
        .courses
        .lock()
        .unwrap()
        .iter()
        .filter(|course| course.teacher_id == teacher_id)
        .cloned() // Course 需实现 Clone
        .collect::<Vec<Course>>();

    // 4.4 REST 风格：空列表给 200 + []，前端不用判字符串
    if !filtered_courses.is_empty() {
        HttpResponse::Ok().json(filtered_courses)
    } else {
        HttpResponse::Ok().json(Vec::<Course>::new()) // 空数组
    }
}

pub async fn get_courses_for_teacher_handle_db(
    app_state: web::Data<AppState>,                   // 1.1 **共享状态** → **Arc<AppState>**，零成本借用
    params: web::Path<(usize, String)>,              // 1.2 **路径参数** → `/courses/{teacher_id}/{name}` → **零成本借用**
) -> Result<HttpResponse, MyErrorNew> {              // 1.3 **返回 Result** → **Ok(Json) 或 Err(MyErrorNew)****

    // 2.1 **解压元组** → (usize, String)
    let teacher_id = i32::try_from(params.0).unwrap(); // 2.2 **usize → i32** → **数据库 integer 对齐**

    // 3.1 **调用数据库函数** → **&Pool → 零成本借用**
    // 3.2 **.await** → **异步等待数据库 IO**，**不阻塞线程**
    // 3.3 **.map(|courses| …)** → **Ok 路径 → 把 Vec<Course> 转成 JSON**
    get_courses_for_teacher_db(&app_state.db, teacher_id)
        .await
        .map(|courses| HttpResponse::Ok().json(courses))   // 3.4 **Ok → JSON 响应**
}

pub async fn get_course_detail_handle_db(
    app_state: web::Data<AppState>,
    params: web::Path<(usize, usize)>,
) -> HttpResponse {
    let teacher_id = i32::try_from(params.0).unwrap();
    let course_id = i32::try_from(params.1).unwrap();
    let course = get_course_detail_db(&app_state.db, teacher_id, course_id).await;
    HttpResponse::Ok().json(course)
}

// ========== 5. 单元测试 ==========
#[cfg(test)]
mod tests {
    use super::*;
    use actix_web::{App, http::StatusCode};
    use dotenv::dotenv; // test里面新增
    use sqlx::postgres::PgPoolOptions;
    use std::env;
    use std::sync::Mutex;

    // 5.1 测试：POST /courses 成功创建
    #[actix_web::test]
    async fn post_course_test() {
        dotenv().ok();
        let database_url = env::var("DATABASE_URL").expect("DatabaseUrl not found");
        let db_pool = PgPoolOptions::new().connect(&database_url).await.unwrap();

        // 5.2 造请求体
        let course = web::Json(Course {
            teacher_id: 1,
            name: "test course".into(),
            id: 3,      // 由服务器生成
            time: None, // 由服务器生成
        });

        // 5.3 造空全局状态
        let app_state = web::Data::new(AppState {
            health_check_response: "OK".to_string(),
            visit_count: Mutex::new(0),
            courses: Mutex::new(vec![]),
            db: db_pool,
        });

        // 5.4 直接调处理器（绕过 HTTP 层，速度最快）
        let resp = new_course(course, app_state).await;

        // 5.5 断言
        assert_eq!(resp.status(), StatusCode::OK);

        // 2. 取出 body → 读成字节 → 再当 &str 用
        let bytes = resp.into_body().try_into_bytes().unwrap(); // Vec<u8>
        let body = std::str::from_utf8(&bytes).unwrap(); // &str
        // 3. 断言
        assert_eq!(body, "course add");
    }

    #[actix_web::test]
    async fn post_course_test_db() {
        dotenv().ok();
        let database_url = env::var("DATABASE_URL").expect("DatabaseUrl not found");
        let db_pool = PgPoolOptions::new().connect(&database_url).await.unwrap();

        // 5.2 造请求体
        let course = web::Json(Course {
            teacher_id: 1,
            name: "test course".into(),
            id: 4,      // 填写None 报错
            time: None, // 由服务器生成
        });

        // 5.3 造空全局状态
        let app_state = web::Data::new(AppState {
            health_check_response: "OK".to_string(),
            visit_count: Mutex::new(0),
            courses: Mutex::new(vec![]),
            db: db_pool,
        });

        // 5.4 直接调处理器（绕过 HTTP 层，速度最快）
        let resp = new_course_handle_db(course, app_state).await;

        // 5.5 断言
        assert_eq!(resp.status(), StatusCode::OK);

        // 2. 取出 body → 读成字节 → 再当 &str 用
        let bytes = resp.into_body().try_into_bytes().unwrap(); // Vec<u8>
        let body = std::str::from_utf8(&bytes).unwrap(); // &str
        let returned: Course = serde_json::from_slice(&bytes).unwrap(); // 反序列化

        // 3. 断言
        assert_eq!(returned.teacher_id, 1);
    }

    // 5.6 测试：GET /courses/{teacher_id}/{name} 空结果
    #[actix_web::test]
    async fn get_course_test() {
        dotenv().ok();
        let database_url = env::var("DATABASE_URL").expect("DatabaseUrl not found");
        let db_pool = PgPoolOptions::new().connect(&database_url).await.unwrap();
        let app_state = web::Data::new(AppState {
            health_check_response: "OK".to_string(),
            visit_count: Mutex::new(0),
            courses: Mutex::new(vec![]), // 空表 → 应返回 []
            db: db_pool,
        });

        // 5.7 构造双段路径
        let params = web::Path::from((1, "asdf".to_string()));
        let response = get_courses_for_teacher(app_state, params).await;

        assert_eq!(response.status(), StatusCode::OK);

        let bytes = response.into_body().try_into_bytes().unwrap(); // Vec<u8>
        let body: Vec<Course> = serde_json::from_slice(&bytes).unwrap();
        assert!(body.is_empty());
    }
}
