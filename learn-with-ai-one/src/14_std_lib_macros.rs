// Rust常用标准库函数与实用宏详解
// 本文件介绍Rust标准库中最常用的函数和宏，帮助开发者提高编程效率

pub fn run_example() {
    println!("=== Rust常用标准库函数与实用宏 ===\n");

    // Option和Result相关函数
    option_result_functions();
    
    // 集合相关函数
    collection_functions();
    
    // 字符串处理函数
    string_functions();
    
    // 数学函数
    math_functions();
    
    // 时间函数
    time_functions();
    
    // 实用宏
    useful_macros();
    
    println!("\n=== 示例结束 ===");
}

fn main() {
    run_example();
}

// 1. Option和Result相关函数
fn option_result_functions() {
    println!("1. Option和Result相关函数:");
    
    // Option类型示例
    let some_value: Option<i32> = Some(42);
    let none_value: Option<i32> = None;
    
    // map: 转换Option中的值
    let mapped = some_value.map(|x| x * 2);
    println!("map转换结果: {:?}", mapped);
    
    // and_then: 链式处理Option
    let result = some_value.and_then(|x| if x > 10 { Some(x + 5) } else { None });
    println!("and_then结果: {:?}", result);
    
    // unwrap_or: 获取值或默认值
    let unwrapped1 = some_value.unwrap_or(0);
    let unwrapped2 = none_value.unwrap_or(0);
    println!("unwrap_or结果: {}, {}", unwrapped1, unwrapped2);
    
    // unwrap_or_else: 获取结果或通过闭包生成默认值
    let unwrapped_else = none_value.unwrap_or_else(|| {
        println!("使用闭包生成默认值");
        100
    });
    println!("unwrap_or_else结果: {}", unwrapped_else);
    
    // ok_or: 将Option转换为Result
    let result_from_option = some_value.ok_or("值不存在");
    println!("ok_or结果: {:?}", result_from_option);
    
    // Result类型示例
    let ok_value: Result<i32, &str> = Ok(100);
    let err_value: Result<i32, &str> = Err("发生错误");
    
    // map_err: 转换错误类型
    let mapped_err = err_value.map_err(|e| format!("错误: {}", e));
    println!("map_err结果: {:?}", mapped_err);
    
    // and_then: 链式处理Result
    let chained_result = ok_value.and_then(|x| {
        if x > 50 {
            Ok(x * 2)
        } else {
            Err("值太小")
        }
    });
    println!("Result and_then结果: {:?}", chained_result);
    
    // unwrap_or_default: 获取值或类型默认值
    let default_ok = ok_value.unwrap_or_default();
    let default_err = err_value.unwrap_or_default();
    println!("unwrap_or_default结果: {}, {}", default_ok, default_err);
    
    println!();
}

// 2. 集合相关函数
fn collection_functions() {
    println!("2. 集合相关函数:");
    
    // Vec相关函数
    let mut vec = vec![3, 1, 4, 1, 5, 9];
    
    // sort: 排序
    vec.sort();
    println!("排序后的Vec: {:?}", vec);
    
    // iter: 创建不可变迭代器
    let sum: i32 = vec.iter().sum();
    println!("Vec元素和: {}", sum);
    
    // iter_mut: 创建可变迭代器
    vec.iter_mut().for_each(|x| *x *= 2);
    println!("每个元素乘以2: {:?}", vec);
    
    // filter: 过滤元素
    let even: Vec<_> = vec.iter().filter(|&&x| x % 2 == 0).collect();
    println!("偶数元素: {:?}", even);
    
    // map: 转换元素
    let squares: Vec<_> = vec![1, 2, 3, 4].iter().map(|&x| x * x).collect();
    println!("平方结果: {:?}", squares);
    
    // fold: 折叠计算
    let product: i32 = vec![1, 2, 3, 4].iter().fold(1, |acc, &x| acc * x);
    println!("乘积结果: {}", product);
    
    // find: 查找元素
    let found = vec.iter().find(|&&x| x > 10);
    println!("大于10的第一个元素: {:?}", found);
    
    // contains: 检查元素是否存在
    let has_element = vec.contains(&8);
    println!("Vec包含8: {}", has_element);
    
    // HashMap相关函数
    use std::collections::HashMap;
    
    let mut map = HashMap::new();
    map.insert("apple", 3);
    map.insert("banana", 5);
    map.insert("orange", 2);
    
    // get: 获取元素
    let apple_count = map.get("apple");
    println!("苹果数量: {:?}", apple_count);
    
    // entry: 安全插入或更新
    map.entry("grape").or_insert(4);
    map.entry("apple").and_modify(|count| *count += 2);
    println!("更新后的HashMap: {:?}", map);
    
    // keys和values: 获取键和值的迭代器
    let keys: Vec<_> = map.keys().collect();
    let values: Vec<_> = map.values().collect();
    println!("键: {:?}, 值: {:?}", keys, values);
    
    println!();
}

// 3. 字符串处理函数
fn string_functions() {
    println!("3. 字符串处理函数:");
    
    let s = String::from("  Rust Programming Language  ");
    let s2 = "Hello, World!";
    
    // trim, trim_start, trim_end: 去除空白字符
    println!("原始字符串: '{}'", s);
    println!("trim结果: '{}'", s.trim());
    println!("trim_start结果: '{}'", s.trim_start());
    println!("trim_end结果: '{}'", s.trim_end());
    
    // split, split_whitespace: 分割字符串
    let parts: Vec<&str> = s.split(" ").collect();
    let words: Vec<&str> = s.split_whitespace().collect();
    println!("split按空格分割: {:?}", parts);
    println!("split_whitespace分割: {:?}", words);
    
    // contains, starts_with, ends_with: 检查字符串内容
    println!("包含'Rust': {}", s.contains("Rust"));
    println!("以'  Ru'开头: {}", s.starts_with("  Ru"));
    println!("以'ge  '结尾: {}", s.ends_with("ge  "));
    
    // replace: 替换字符串
    let replaced = s.replace("Rust", "RUST");
    println!("替换'Rust'为'RUST': '{}'", replaced);
    
    // to_uppercase, to_lowercase: 转换大小写
    println!("大写: {}", s.to_uppercase());
    println!("小写: {}", s.to_lowercase());
    
    // chars: 字符迭代器
    let first_char = s.chars().next();
    println!("第一个字符: {:?}", first_char);
    
    // bytes: 字节迭代器
    let byte_count = s.bytes().count();
    println!("字节数: {}", byte_count);
    
    // parse: 字符串解析
    let num_str = "42";
    let num: Result<i32, _> = num_str.parse();
    println!("字符串解析为数字: {:?}", num);
    
    println!();
}

// 4. 数学函数
fn math_functions() {
    println!("4. 数学函数:");
    
    use std::f64;
    
    // 常量
    println!("数学常量:");
    println!("π: {}", f64::consts::PI);
    println!("e: {}", f64::consts::E);
    println!("√2: {}", f64::consts::SQRT_2);
    
    // 三角函数
    let angle = f64::consts::PI / 2.0; // 90度
    println!("\n三角函数:");
    println!("sin(π/2): {}", angle.sin());
    println!("cos(π/2): {}", angle.cos());
    println!("tan(π/4): {}", (f64::consts::PI / 4.0).tan());
    
    // 反三角函数
    println!("\n反三角函数:");
    println!("asin(1): {}", 1.0_f64.asin());
    println!("acos(0): {}", 0.0_f64.acos());
    println!("atan(1): {}", 1.0_f64.atan());
    
    // 指数和对数
    println!("\n指数和对数:");
    println!("e^2: {}", 2.0_f64.exp());
    println!("ln(e^2): {}", 2.0_f64.exp().ln());
    println!("log10(100): {}", 100.0_f64.log10());
    println!("2^8: {}", 2.0_f64.powf(8.0));
    
    // 其他数学函数
    println!("\n其他数学函数:");
    println!("√16: {}", 16.0_f64.sqrt());
    println!("立方根27: {}", 27.0_f64.cbrt());
    println!("绝对值-10: {}", (-10.0_f64).abs());
    println!("最大值5和10: {}", 5.0_f64.max(10.0_f64));
    println!("最小值5和10: {}", 5.0_f64.min(10.0_f64));
    println!("取整3.7: {}", 3.7_f64.floor());
    println!("向上取整3.2: {}", 3.2_f64.ceil());
    println!("四舍五入3.5: {}", 3.5_f64.round());
    
    println!();
}

// 5. 时间函数
fn time_functions() {
    println!("5. 时间函数:");
    
    use std::time::{Instant, Duration};
    
    // Instant: 测量时间点
    let start = Instant::now();
    
    // 模拟耗时操作
    std::thread::sleep(Duration::from_millis(100));
    
    let elapsed = start.elapsed();
    println!("耗时: {:?}", elapsed);
    
    // Duration: 表示时间段
    let duration1 = Duration::from_secs(5);
    let duration2 = Duration::from_millis(500);
    let total = duration1 + duration2;
    
    println!("Duration1: {:?}", duration1);
    println!("Duration2: {:?}", duration2);
    println!("总时长: {:?}", total);
    println!("总秒数: {}", total.as_secs_f64());
    
    // SystemTime: 系统时间
    use std::time::SystemTime;
    
    let now = SystemTime::now();
    
    // 转换为UNIX时间戳
    if let Ok(unix_time) = now.duration_since(SystemTime::UNIX_EPOCH) {
        println!("当前UNIX时间戳: {}", unix_time.as_secs());
    }
    
    println!();
}

// 6. 实用宏
fn useful_macros() {
    println!("6. 实用宏:");
    
    // println!: 打印输出
    println!("println!宏: 这是一个示例");
    println!("格式化输出: 数字{}，字符串'{}", 42, "hello");
    
    // format!: 格式化字符串
    let formatted = format!("数字{}，字符串'{}", 42, "world");
    println!("format!宏结果: {}", formatted);
    
    // vec!: 创建Vec
    let v = vec![1, 2, 3, 4, 5];
    println!("vec!宏创建的向量: {:?}", v);
    
    // Box::new和box!: 创建堆分配
    let boxed = Box::new(42);
    println!("Box::new创建的堆分配: {:?}", boxed);
    
    // Ok!和Err!: 创建Result
    let ok_result: Result<i32, &str> = Ok(42);
    let err_result: Result<i32, &str> = Err("错误");
    println!("Result宏结果: {:?}, {:?}", ok_result, err_result);
    
    // Some!: 创建Option
    let some_value = Some(42);
    println!("Some!宏结果: {:?}", some_value);
    
    // panic!: 程序崩溃
    // panic!("这会导致程序崩溃");
    
    // unwrap!: 解包Option或Result，如果失败则崩溃
    let unwrapped = Ok::<i32, &str>(42).unwrap();
    println!("unwrap!宏结果: {}", unwrapped);
    
    // assert!: 断言，失败则崩溃
    assert!(1 + 1 == 2, "1 + 1 应该等于2");
    println!("assert!宏通过");
    
    // assert_eq!: 断言相等
    assert_eq!(2 * 2, 4, "2 * 2 应该等于4");
    println!("assert_eq!宏通过");
    
    // unreachable!: 标记不可达代码
    fn example(x: bool) {
        if x {
            println!("x为true");
        } else {
            println!("x为false");
        }
        // unreachable!(); // 如果执行到这里就崩溃
    }
    example(true);
    
    // todo!: 标记待实现的代码
    fn todo_example() {
        // todo!("这个函数还没有实现");
    }
    
    // dbg!: 调试宏，打印表达式和值
    let a = 2;
    let b = dbg!(a * 2) + 1;
    println!("dbg!宏结果: b = {}", b);
    
    println!();
}