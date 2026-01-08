// ========== 1. 引入标准库和第三方库 ==========
use chrono::NaiveDateTime;          // 日期时间类型（无时区）
use dotenv::dotenv;                 // 加载 .env 文件到环境变量
use sqlx::postgres::PgPoolOptions;  // PostgreSQL 连接池
use std::env;                       // 读取环境变量
use std::io;                        // main 函数返回 io::Result

// ========== 2. 定义领域模型 ==========
#[derive(Debug)]                    // 自动生成 Debug 打印格式
pub struct Course {
    pub id: i32,                    // 主键
    pub teacher_id: i32,            // 外键
    pub name: String,               // 课程名
    pub time: Option<NaiveDateTime>, // 时间戳可空（Option → 显式空值）
}

// ========== 3. 异步 main（钉在 tokio 上） ==========
#[actix_web::main]                  // 宏：把 async main 绑在 tokio 运行时
async fn main() -> io::Result<()> { // 返回 I/O 错误类型（main 能返回）

    // 3.1 把 .env 文件加载到进程环境变量（失败也不 panic）
    dotenv().ok();
    println!("Hello, world!");

    // 3.2 读数据库连接串；expect 在缺失时给出友好错误
    let database_url = env::var("DATABASE_URL").expect("Database not in .env");

    // 3.3 **连接池**：复用 TCP + 会话，**比每次新建连接快 10×**
    //     PgPoolOptions::new() → 默认 10 连接，**异步**  
    let db_pool = PgPoolOptions::new()
        .connect(&database_url)      // **&str** → 借用，不拷贝
        .await                        // 异步等待 TCP + TLS 握手
        .unwrap();                    // 简化错误处理（测试可接受）

    // 3.4 **编译期检查 SQL**（sqlx::query! 宏）
    //     **占位符 $1** → PostgreSQL 风格；**参数类型必须对**（i32）
    let course_rows = sqlx::query!(
        r#"select * from rust_test1.course where id=$1"#,
        1i32                             // **i32** 与 SQL **integer** 对应
    )
    .fetch_all(&db_pool)               // **&Pool** → 借用池，**不转移所有权**
    .await                             // 异步等待结果集
    .unwrap();                         // 简化错误（生产用 ?）

    // 3.5 **空 Vec** 准备装结构体
    let mut course_list = vec![];

    // 3.6 **for 循环** → 把 **sqlx 返回的行** 转成 **自己定义的 Course**
    for row in course_rows {
        // 3.7 **row.id** → 编译期已知类型（i32），**直接拿**
        //     **row.teacher_id.unwrap()** → SQL 允许 NULL，**Option<i32>** → 手动解包
        //     **&db_pool** vs **row.id** → **& 表示“借用”**，**不拷贝大对象**
        course_list.push(Course {
            id: row.id,
            teacher_id: row.teacher_id.unwrap(),   // NULL → panic（测试可接受）
            name: row.name.unwrap(),               // NULL → panic
            time: Some(chrono::NaiveDateTime::from(row.time.unwrap())), // NULL → panic
        });
    }

    // 3.8 **Debug 打印** → 宏自动生成格式
    println!("courses are ={:?}", course_list);

    // 3.9 **Ok(())** → main 返回成功，**io::Result<()>** 要求
    Ok(())
}