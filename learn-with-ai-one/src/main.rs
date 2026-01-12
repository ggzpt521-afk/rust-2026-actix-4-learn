// Rust学习示例主程序
// 该程序提供一个菜单，允许用户选择要运行的示例

use std::process::Command;
use std::path::Path;
use std::io::{self, BufRead};

// 导入13-15号文件，这些文件已经有run_example函数
#[path = "13_lifetimes.rs"] mod lifetimes;
#[path = "14_std_lib_macros.rs"] mod std_lib_macros;
#[path = "15_async_await.rs"] mod async_await;

fn main() {
    println!("=== Rust学习示例 ===\n");
    println!("请选择要运行的示例：");
    println!("1. 变量和可变性");
    println!("2. 数据类型");
    println!("3. 函数");
    println!("4. 控制流");
    println!("5. 所有权");
    println!("6. 结构体");
    println!("7. 枚举和模式匹配");
    println!("8. 集合");
    println!("9. 包、 crate 和模块");
    println!("10. 错误处理");
    println!("11. 泛型");
    println!("12. Trait");
    println!("13. 生命周期");
    println!("14. 常用标准库函数与实用宏");
    println!("15. 异步编程（async/await）");
    println!("0. 退出");
    println!();
    
    print!("请输入选择 (0-15): ");
    // 手动刷新输出缓冲区，确保提示信息先显示
    std::io::Write::flush(&mut std::io::stdout()).unwrap();
    
    // 从用户输入读取选择
    let mut choice: u8;
    let stdin = io::stdin();
    let mut input = String::new();
    
    loop {
        input.clear();
        if let Err(e) = stdin.lock().read_line(&mut input) {
            println!("读取输入错误: {}", e);
            continue;
        }
        
        // 去除输入中的换行符和空格
        let trimmed = input.trim();
        
        // 尝试解析为u8
        match trimmed.parse::<u8>() {
            Ok(num) if num <= 15 => {
                choice = num;
                break;
            }
            Ok(_) => println!("选择无效，请输入0-15之间的数字"),
            Err(_) => println!("输入格式错误，请输入数字"),
        }
    }
    
    println!("\n您选择了: {}", choice);
    
    match choice {
        1 => run_example_file("01_variables"),
        2 => run_example_file("02_data_types"),
        3 => run_example_file("03_functions"),
        4 => run_example_file("04_control_flow"),
        5 => run_example_file("05_ownership"),
        6 => run_example_file("06_structs"),
        7 => run_example_file("07_enums"),
        8 => run_example_file("08_collections"),
        9 => run_example_file("09_packages_modules"),
        10 => run_example_file("10_error_handling"),
        11 => run_example_file("11_generics"),
        12 => run_example_file("12_traits"),
        13 => lifetimes::run_example(),
        14 => std_lib_macros::run_example(),
        15 => async_await::run_example(),
        0 => println!("退出程序") ,
        _ => println!("无效选择") ,
    }
    
    println!("\n=== 程序结束 ===");
}

// 运行示例文件的函数
fn run_example_file(filename: &str) {
    let file_path = Path::new("src").join(format!("{}.rs", filename));
    
    if !file_path.exists() {
        println!("错误：文件 {} 不存在", file_path.display());
        return;
    }
    
    println!("\n正在运行示例: {}\n", filename);
    
    // 编译并运行示例文件
    let output = Command::new("rustc")
        .arg(&file_path)
        .arg("-o")
        .arg(filename)
        .output();
    
    match output {
        Ok(compilation) => {
            if compilation.status.success() {
                // 运行编译后的程序
                let run_output = Command::new(format!("./{}", filename))
                    .output();
                
                match run_output {
                    Ok(run) => {
                        if run.status.success() {
                            println!("{}", String::from_utf8_lossy(&run.stdout));
                        } else {
                            println!("运行错误:");
                            println!("{}", String::from_utf8_lossy(&run.stderr));
                        }
                    }
                    Err(e) => println!("运行失败: {}", e),
                }
                
                // 清理编译后的程序
                let _ = Command::new("rm").arg(filename).output();
            } else {
                println!("编译错误:");
                println!("{}", String::from_utf8_lossy(&compilation.stderr));
            }
        }
        Err(e) => println!("编译失败: {}", e),
    }
}
