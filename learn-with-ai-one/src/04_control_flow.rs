// 04_control_flow.rs - Rust流程控制详解

pub fn run_example() {
    println!("=== Rust学习示例 ===\n");
    // 1. if条件表达式
    // if语句在Rust中是表达式，不是语句，所以可以返回值
    let number = 6;
    
    if number % 2 == 0 {
        println!("{}是偶数", number);
    } else {
        println!("{}是奇数", number);
    }
    
    // if作为表达式返回值
    let result = if number > 0 {
        "正数"
    } else if number < 0 {
        "负数"
    } else {
        "零"
    };
    
    println!("{}是{}", number, result);
    
    // 注意：if的所有分支必须返回相同类型的值
    // let invalid = if number > 0 {
    //     1
    // } else {
    //     "错误"
    // };
    
    // 2. if let表达式
    // if let用于简化模式匹配，处理只关心一种情况的场景
    let optional_number: Option<i32> = Some(7);
    
    // 使用match
    match optional_number {
        Some(n) => println!("使用match获取到数字: {}", n),
        _ => ()  // 忽略其他情况
    }
    
    // 使用if let简化
    if let Some(n) = optional_number {
        println!("使用if let获取到数字: {}", n);
    }
    
    // if let配合else
    let optional_string: Option<String> = None;
    if let Some(s) = optional_string {
        println!("获取到字符串: {}", s);
    } else {
        println!("没有字符串");
    }
    
    // 3. match表达式（模式匹配）
    // match是Rust中强大的模式匹配结构
    let day = 3;
    
    let day_name = match day {
        1 => "星期一",
        2 => "星期二",
        3 => "星期三",
        4 => "星期四",
        5 => "星期五",
        6 => "星期六",
        7 => "星期日",
        _ => "无效的天数"  // 必须处理所有可能的情况，_是通配符
    };
    
    println!("{}号是{}", day, day_name);
    
    // match用于枚举
    enum Direction {
        Up,
        Down,
        Left,
        Right
    }
    
    let dir = Direction::Right;
    
    match dir {
        Direction::Up => println!("向上移动"),
        Direction::Down => println!("向下移动"),
        Direction::Left => println!("向左移动"),
        Direction::Right => println!("向右移动")
        // 这里不需要_，因为已经覆盖了所有枚举变体
    }
    
    // 4. for循环
    // for循环用于遍历迭代器或范围
    
    // 遍历范围（左闭右开）
    println!("\nfor循环遍历范围1..5:");
    for i in 1..5 {
        print!("{}", i);
        if i < 4 { print!(" "); }
    }
    println!();
    
    // 遍历包含结束值的范围
    println!("for循环遍历范围1..=5:");
    for i in 1..=5 {
        print!("{}", i);
        if i < 5 { print!(" "); }
    }
    println!();
    
    // 遍历数组
    let numbers = [10, 20, 30, 40, 50];
    println!("\nfor循环遍历数组:");
    for number in numbers.iter() {
        print!("{}", number);
        if *number < 50 { print!(" "); }
    }
    println!();
    
    // 遍历数组并获取索引
    println!("for循环遍历数组并获取索引:");
    for (index, number) in numbers.iter().enumerate() {
        println!("索引{}: 值{}", index, number);
    }
    
    // 遍历字符串
    let message = "Hello";
    println!("\nfor循环遍历字符串字符:");
    for c in message.chars() {
        print!("'{}\' ", c);
    }
    println!();
    
    // 5. while循环
    // while循环在条件为真时执行
    let mut count = 0;
    println!("\nwhile循环示例:");
    while count < 5 {
        println!("count = {}", count);
        count += 1;
    }
    
    // 6. loop循环
    // loop是无限循环，必须使用break终止
    let mut loop_count = 0;
    println!("\nloop循环示例:");
    
    loop {
        println!("loop_count = {}", loop_count);
        loop_count += 1;
        
        if loop_count >= 3 {
            break;  // 终止循环
        }
    }
    
    // loop作为表达式返回值
    let mut counter = 0;
    
    let result = loop {
        counter += 1;
        
        if counter == 10 {
            break counter * 2;  // 返回值
        }
    };
    
    println!("loop返回值: {}", result);
    
    // 7. break和continue控制流
    println!("\nbreak和continue示例:");
    
    for i in 1..=10 {
        if i % 2 == 0 {
            continue;  // 跳过偶数
        }
        
        if i > 7 {
            break;  // 大于7时终止循环
        }
        
        print!("{}", i);
        if i < 7 { print!(" "); }
    }
    println!();
    
    // 8. 循环标签
    // 使用'标签名: loop'来标记循环，用于嵌套循环中的break和continue
    println!("\n循环标签示例:");
    
    'outer: loop {
        let mut inner_count = 0;
        
        'inner: loop {
            println!("outer loop, inner_count = {}", inner_count);
            inner_count += 1;
            
            if inner_count >= 3 {
                break 'inner;  // 终止内部循环
            }
            
            if inner_count == 1 {
                break 'outer;  // 直接终止外部循环
            }
        }
    }
    
    // 9. while let条件循环
    // while let用于简化包含模式匹配的循环
    let mut stack = Vec::new();
    stack.push(1);
    stack.push(2);
    stack.push(3);
    
    println!("\nwhile let循环示例:");
    while let Some(top) = stack.pop() {
        println!("从栈中弹出: {}", top);
    }
    
    // 10. 流程控制的组合使用
    println!("\n流程控制组合使用示例:");
    
    let numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
    let mut even_sum = 0;
    let mut odd_sum = 0;
    
    for number in numbers.iter() {
        match number % 2 {
            0 => even_sum += number,
            1 => odd_sum += number,
            _ => unreachable!()  // 不可能的情况，用于调试
        }
    }
    
    println!("偶数和: {}, 奇数和: {}", even_sum, odd_sum);
}

// 11. 示例：使用流程控制实现FizzBuzz游戏
fn fizzbuzz() {
    println!("\nFizzBuzz游戏实现:");
    
    for n in 1..=20 {
        match (n % 3 == 0, n % 5 == 0) {
            (true, true) => println!("FizzBuzz"),
            (true, false) => println!("Fizz"),
            (false, true) => println!("Buzz"),
            (false, false) => println!("{}", n)
        }
    }
}
// 用于单独运行本文件的main函数
fn main() {
    run_example();
}
