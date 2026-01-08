// ========== 1. 依赖与类型 ==========
use super::models::*;               // 引入本地定义的 Course 结构体
use sqlx::postgres::PgPool;         // PostgreSQL 异步连接池（比单连接快 10×）


// ========== 2. 根据老师 ID 查所有课程 ==========
pub async fn get_courses_for_teacher_db(
    pool: &PgPool,                   // 2.1 **借用连接池** → 不转移所有权，**零成本**
    teacher_id: i32,                // 2.2 **i32** ↔ SQL **integer**，**类型必须对**
) -> Vec<Course> {                  // 2.3 返回 **Vec<Course>** → **零成本返回**（只是指针移动）

    // 2.4 **编译期检查 SQL**（sqlx::query! 宏）
    //     **占位符 $1** → PostgreSQL 风格；**参数类型必须对**（i32）
    let rows = sqlx::query!(
        r#"SELECT * FROM rust_test1.course WHERE teacher_id = $1"#,
        teacher_id
    )
    .fetch_all(pool)                 // 2.5 **异步取全部行** → **返回 Vec<PgRow>**
    .await                            // 2.6 **等待 IO 完成** → **不会阻塞线程**
    .unwrap();                        // 2.7 **简化错误**（测试可接受，生产用 ?）

    // 2.8 **Vec<Course>** 准备装结构体（零成本，只是指针数组）
    rows.iter()
        .map(|r| Course {             // 2.9 **逐行映射** → **零成本迭代**
            id: r.id,                                      // i32 ↔ INTEGER
            teacher_id: r.teacher_id.unwrap_or(0),         // Option<i32> → i32
            name: r.name.clone().unwrap_or_default(),      // Option<String> → String
            time: r.time,                                  // Option<NaiveDateTime> 直接用
        })
        .collect()                     // 2.14 **Vec<Course>** → **零成本收集**
}

// ========== 3. 根据老师 ID + 课程 ID 查单条课程 ==========
pub async fn get_course_detail_db(
    pool: &PgPool,                   // 3.1 **借用连接池** → **零成本**
    teacher_id: i32,                // 3.2 **i32 ↔ integer**
    course_id: i32,                // 3.3 **i32 ↔ integer**
) -> Course {                      // 3.4 返回 **单个 Course** → **零成本返回**

    // 3.5 **编译期检查 SQL** → **双条件查询**
    let row = sqlx::query!(
        r#"SELECT * FROM rust_test1.course WHERE teacher_id = $1 AND id = $2"#,
        teacher_id,
        course_id
    )
    .fetch_one(pool)                 // 3.6 **异步取一行** → **返回 PgRow**
    .await
    .unwrap();                       // 3.7 **unwrap()** → **测试可接受，生产用 ?**

    // 3.8 **直接构造 Course** → **零成本映射**
    Course {
        id: row.id,
        teacher_id: row.teacher_id.unwrap_or(0),
        name: row.name.clone().unwrap_or_default(),
        time: row.time,
    }
}

// ========== 4. 插入新课程并返回刚插入的行 ==========
pub async fn post_new_course_db(
    pool: &PgPool,                   // 4.1 **借用连接池** → **零成本**
    new_course: Course,              // 4.2 **Course 整体 move 进来** → **零成本（只是指针移动）**
) -> Course {                      // 4.3 返回 **刚插入的完整行** → **零成本返回**

    // 4.4 **编译期检查 SQL** → **INSERT … VALUES ($1,$2)**
    //     **不插入 id**：id 是 GENERATED ALWAYS（自增列），由数据库生成
    //     **fetch_one()** → **PostgreSQL 支持 RETURNING** → **返回刚插入的行**
    let row = sqlx::query!(
        r#"INSERT INTO rust_test1.course (teacher_id, name) VALUES ($1, $2) RETURNING *"#,
        new_course.teacher_id,
        new_course.name
    )
    .fetch_one(pool)                 // 4.5 **RETURNING * → 返回刚插入的行**
    .await
    .unwrap();                        // 4.6 **unwrap()** → **测试可接受**

    // 4.7 **直接构造返回的 Course** → **零成本映射**
    Course {
        id: row.id,
        teacher_id: row.teacher_id.unwrap_or(0),
        name: row.name.clone().unwrap_or_default(),
        time: row.time,
    }
}