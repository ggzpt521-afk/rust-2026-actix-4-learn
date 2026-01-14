// ============================================================================
// 08_cross_platform.rs - Slint 跨平台构建示例
// ============================================================================
//
// 【核心概念】
// Slint 的跨平台特性：
// 1. 一套代码运行在多个平台（Windows, macOS, Linux, 嵌入式）
// 2. 原生性能：编译为原生代码，不需要虚拟机
// 3. 统一渲染：在所有平台上外观一致
// 4. 平台集成：可以访问平台特定的功能
//
// 【原理说明】
// Slint 的跨平台实现：
// - 使用 Rust 作为后端，Rust 本身支持跨平台
// - 渲染后端可选：OpenGL, Skia, 软件渲染
// - 窗口管理使用 winit 库（跨平台窗口抽象）
// - 条件编译：#[cfg(target_os = "...")] 处理平台差异
// ============================================================================

slint::slint! {
    export component CrossPlatform inherits Window {
        width: 400px;
        height: 350px;
        title: "跨平台构建示例";

        // ====================================================================
        // 平台信息属性
        // ====================================================================
        // 这些属性由 Rust 代码在运行时设置
        // 展示了如何将系统信息传递给 UI
        in property <string> platform-name: "";   // 平台名称
        in property <string> platform-icon: "";   // 平台图标（emoji）
        in property <string> os-version: "";      // 系统版本
        in property <string> architecture: "";    // CPU 架构

        VerticalLayout {
            padding: 20px;
            spacing: 15px;

            // 标题区域，显示平台图标和标题
            HorizontalLayout {
                spacing: 10px;
                alignment: center;

                // 平台图标（emoji）
                Text {
                    text: platform-icon;
                    font-size: 32px;
                }

                Text {
                    text: "跨平台构建示例";
                    font-size: 24px;
                    color: #333;
                }
            }

            // ================================================================
            // 平台信息卡片
            // ================================================================
            // 显示从 Rust 代码获取的系统信息
            Rectangle {
                width: 100%;
                background: white;
                border-width: 1px;
                border-color: #e0e0e0;
                border-radius: 8px;

                VerticalLayout {
                    padding: 20px;
                    spacing: 10px;

                    // 平台名称
                    HorizontalLayout {
                        spacing: 10px;
                        Text {
                            text: "平台:";
                            color: #333;
                            width: 80px;
                        }
                        Text {
                            text: platform-name;
                            color: #0066cc;
                        }
                    }

                    // 系统版本
                    HorizontalLayout {
                        spacing: 10px;
                        Text {
                            text: "版本:";
                            color: #333;
                            width: 80px;
                        }
                        Text {
                            text: os-version;
                            color: #333;
                        }
                    }

                    // CPU 架构
                    HorizontalLayout {
                        spacing: 10px;
                        Text {
                            text: "架构:";
                            color: #333;
                            width: 80px;
                        }
                        Text {
                            text: architecture;
                            color: #333;
                        }
                    }
                }
            }

            // ================================================================
            // 跨平台功能说明
            // ================================================================
            Rectangle {
                width: 100%;
                background: white;
                border-width: 1px;
                border-color: #e0e0e0;
                border-radius: 8px;

                VerticalLayout {
                    padding: 15px;
                    spacing: 5px;

                    Text { text: "跨平台特性:"; font-size: 14px; color: #333; }
                    Text { text: "• 自动适应不同操作系统"; font-size: 12px; color: #666; }
                    Text { text: "• 统一的代码库"; font-size: 12px; color: #666; }
                    Text { text: "• 原生性能体验"; font-size: 12px; color: #666; }
                    Text { text: "• 支持 Windows、macOS、Linux"; font-size: 12px; color: #666; }
                }
            }
        }
    }
}

// ============================================================================
// main 函数
// ============================================================================
fn main() {
    let app = CrossPlatform::new().unwrap();

    // ------------------------------------------------------------------------
    // 获取平台信息
    // ------------------------------------------------------------------------
    // 使用 sys-info crate 获取系统信息
    // 这些信息在编译时无法确定，必须在运行时获取

    // os_type(): 返回操作系统类型
    // - "Darwin": macOS
    // - "Linux": Linux
    // - "Windows_NT": Windows
    let os_type = sys_info::os_type().unwrap_or("未知".into());

    // os_release(): 返回系统版本号
    let os_release = sys_info::os_release().unwrap_or("未知".into());

    // ------------------------------------------------------------------------
    // 根据操作系统类型设置图标和名称
    // ------------------------------------------------------------------------
    // 使用 match 表达式进行模式匹配
    let (name, icon) = match os_type.as_str() {
        "Darwin" => ("macOS", "🍎"),      // Apple macOS
        "Linux" => ("Linux", "🐧"),        // Linux (企鹅)
        "Windows_NT" => ("Windows", "🪟"), // Windows (窗户)
        _ => (os_type.as_str(), "📱"),     // 其他/未知
    };

    // 获取 CPU 架构
    let arch = get_arch();

    // ------------------------------------------------------------------------
    // 设置平台信息到组件
    // ------------------------------------------------------------------------
    // 使用自动生成的 set_xxx 方法设置属性
    // .into() 将 &str 转换为 SharedString
    app.set_platform_name(name.into());
    app.set_platform_icon(icon.into());
    app.set_os_version(os_release.into());
    app.set_architecture(arch.into());

    app.run().unwrap();
}

// ============================================================================
// 获取 CPU 架构信息
// ============================================================================
// 使用条件编译 (#[cfg(...)]) 在编译时确定 CPU 架构
// 这是 Rust 的编译时特性，不是运行时检测
fn get_arch() -> String {
    // #[cfg(target_arch = "x86")]: 当目标架构为 x86 时编译此代码
    #[cfg(target_arch = "x86")]
    return "x86".into();

    // x86_64: 64 位 Intel/AMD 处理器
    #[cfg(target_arch = "x86_64")]
    return "x86_64".into();

    // arm: 32 位 ARM 处理器
    #[cfg(target_arch = "arm")]
    return "ARM".into();

    // aarch64: 64 位 ARM 处理器 (Apple M1/M2, 新款手机等)
    #[cfg(target_arch = "aarch64")]
    return "ARM64".into();

    // 其他架构的兜底处理
    #[cfg(not(any(
        target_arch = "x86",
        target_arch = "x86_64",
        target_arch = "arm",
        target_arch = "aarch64"
    )))]
    return "其他架构".into();
}

// ============================================================================
// 【知识点总结】
// ============================================================================
//
// 1. 跨平台构建
//    - cargo build: 构建当前平台
//    - cargo build --target xxx: 交叉编译到其他平台
//    - 常用目标：
//      - x86_64-pc-windows-msvc (Windows)
//      - x86_64-unknown-linux-gnu (Linux)
//      - x86_64-apple-darwin (macOS Intel)
//      - aarch64-apple-darwin (macOS Apple Silicon)
//
// 2. 条件编译
//    - #[cfg(target_os = "windows")]: 按操作系统
//    - #[cfg(target_arch = "x86_64")]: 按 CPU 架构
//    - #[cfg(feature = "xxx")]: 按功能特性
//
// 3. 平台信息获取
//    - sys-info crate: 获取系统信息
//    - std::env::consts: 标准库常量
//
// 4. Slint 渲染后端
//    - femtovg: 基于 OpenGL 的矢量渲染
//    - skia: Google 的 2D 图形库
//    - software: 纯软件渲染（嵌入式）
//
// 5. 跨平台最佳实践
//    - 避免硬编码路径分隔符
//    - 使用标准库的跨平台 API
//    - 条件编译处理平台差异
//    - 测试所有目标平台
//
// 6. 打包发布
//    - Windows: .exe 文件或 MSIX
//    - macOS: .app 包或 DMG
//    - Linux: AppImage, Flatpak, 或 DEB/RPM
// ============================================================================
