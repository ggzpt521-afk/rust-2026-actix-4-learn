// ============================================================================
// 09_async_data.rs - Slint 异步数据交互示例
// ============================================================================
//
// 【核心概念】
// 异步数据处理是现代应用的关键：
// 1. 网络请求：获取远程数据
// 2. 文件操作：读写大文件
// 3. 计算任务：CPU 密集型操作
// 4. 保持 UI 响应：不阻塞主线程
//
// 【原理说明】
// Slint 的异步模型：
// - UI 运行在主线程（事件循环）
// - 耗时操作应在后台线程执行
// - 使用 slint::invoke_from_event_loop() 从后台线程更新 UI
// - 这保证了线程安全，避免数据竞争
//
// 【为什么需要异步】
// 如果在主线程执行耗时操作（如网络请求）：
// - UI 会冻结，无法响应用户操作
// - 系统可能判定程序"无响应"
// - 用户体验极差
// ============================================================================

slint::slint! {
    import { Button, LineEdit } from "std-widgets.slint";

    export component AsyncData inherits Window {
        width: 400px;
        height: 350px;
        title: "异步数据示例";

        // ====================================================================
        // 应用状态属性
        // ====================================================================
        // 这些属性用于跟踪异步操作的状态

        // 加载状态：是否正在加载数据
        in-out property <bool> is-loading: false;

        // 状态消息：显示给用户的提示信息
        in-out property <string> status-message: "点击按钮获取数据";

        // 结果数据：异步操作返回的数据
        in-out property <string> result-data: "";

        // 用户输入
        in-out property <string> input-text: "北京";

        VerticalLayout {
            padding: 20px;
            spacing: 15px;

            Text {
                text: "异步数据交互示例";
                font-size: 24px;
                color: #333;
                horizontal-alignment: center;
            }

            // 输入区域
            HorizontalLayout {
                spacing: 10px;
                alignment: center;

                // 文本输入框
                LineEdit {
                    text: input-text;
                    width: 200px;
                    // 双向绑定：UI 变化 -> 属性更新
                    edited => { input-text = self.text; }
                }

                // 提交按钮
                Button {
                    // 条件表达式：根据加载状态显示不同文本
                    text: is-loading ? "加载中..." : "获取数据";

                    // enabled: 控制按钮是否可点击
                    // 加载中或输入为空时禁用按钮
                    enabled: !is-loading && input-text != "";

                    clicked => {
                        // 更新状态，表示开始加载
                        is-loading = true;
                        status-message = "正在获取数据...";

                        // 调用回调，触发 Rust 端的异步操作
                        fetch-data(input-text);
                    }
                }
            }

            // 状态消息显示
            Text {
                text: status-message;
                font-size: 14px;
                color: #666;
                horizontal-alignment: center;
            }

            // ================================================================
            // 条件渲染：只有当有结果时才显示结果卡片
            // ================================================================
            if result-data != "": Rectangle {
                width: 100%;
                height: 120px;
                background: white;
                border-width: 1px;
                border-color: #e0e0e0;
                border-radius: 8px;

                VerticalLayout {
                    padding: 15px;
                    spacing: 10px;

                    Text {
                        text: "查询结果:";
                        font-size: 14px;
                        color: #333;
                    }

                    Text {
                        text: result-data;
                        font-size: 16px;
                        color: #0066cc;
                    }
                }
            }

            // 说明文字
            Text {
                text: "本示例演示异步数据处理:\n• 后台线程执行\n• UI 保持响应\n• 显示加载状态";
                font-size: 12px;
                color: #999;
                horizontal-alignment: center;
            }
        }

        // ====================================================================
        // 回调声明
        // ====================================================================
        // callback: 声明一个可以从 Slint 调用、在 Rust 中实现的函数
        // 语法：callback 名称(参数类型, ...);
        //
        // 这是 Slint 和 Rust 之间的桥梁：
        // - Slint 端：调用 fetch-data(...)
        // - Rust 端：实现 on_fetch_data(|args| { ... })
        callback fetch-data(string);
    }
}

// 引入标准库的线程和时间模块
use std::thread;
use std::time::Duration;

// ============================================================================
// main 函数
// ============================================================================
fn main() {
    let app = AsyncData::new().unwrap();

    // ------------------------------------------------------------------------
    // 弱引用 (Weak Reference)
    // ------------------------------------------------------------------------
    // as_weak() 创建组件的弱引用
    // 弱引用不会阻止组件被销毁
    // 在回调中使用弱引用避免循环引用
    let app_weak = app.as_weak();

    // ------------------------------------------------------------------------
    // 注册回调处理函数
    // ------------------------------------------------------------------------
    // on_fetch_data: 自动生成的方法，用于注册 fetch-data 回调
    // 当 Slint 端调用 fetch-data(query) 时，这个闭包会被执行
    app.on_fetch_data(move |query| {
        // 克隆弱引用用于闭包
        // 需要在闭包外克隆，因为闭包会 move 进新线程
        let app_weak = app_weak.clone();

        // 转换为 Rust String
        // Slint 的 SharedString 需要转换为 String 才能 move 进新线程
        let query = query.to_string();

        // --------------------------------------------------------------------
        // 创建新线程执行异步操作
        // --------------------------------------------------------------------
        // thread::spawn: 创建新线程
        // move ||: 闭包获取外部变量的所有权
        //
        // 【为什么用线程而不是 async/await】
        // 1. 简单直接，不需要异步运行时
        // 2. Slint 主循环不是异步的
        // 3. 对于简单的异步任务足够了
        thread::spawn(move || {
            // 模拟网络延迟（实际应用中这里是真正的网络请求）
            thread::sleep(Duration::from_secs(1));

            // 模拟获取的数据
            let result = format!("城市: {}\n温度: 22°C\n天气: 晴朗", query);

            // ----------------------------------------------------------------
            // 从后台线程更新 UI
            // ----------------------------------------------------------------
            // slint::invoke_from_event_loop: 在主线程事件循环中执行代码
            //
            // 【为什么需要这个函数】
            // - Slint 的 UI 不是线程安全的
            // - 只能在主线程修改 UI 属性
            // - 这个函数将闭包调度到主线程执行
            //
            // 【原理】
            // 1. 将闭包发送到主线程的消息队列
            // 2. 主线程的事件循环取出闭包
            // 3. 在主线程上下文中执行闭包
            // 4. 此时修改属性是安全的
            slint::invoke_from_event_loop(move || {
                // upgrade(): 尝试将弱引用转换为强引用
                // 如果组件已销毁，返回 None
                if let Some(app) = app_weak.upgrade() {
                    // 更新 UI 属性
                    app.set_result_data(result.into());
                    app.set_status_message("数据获取成功!".into());
                    app.set_is_loading(false);
                }
                // 如果 upgrade() 返回 None，说明窗口已关闭
                // 此时什么也不做，安全退出
            }).unwrap();
        });
    });

    // 运行应用
    app.run().unwrap();
}

// ============================================================================
// 【知识点总结】
// ============================================================================
//
// 1. 异步模式
//    - UI 线程：运行事件循环，处理用户交互
//    - 工作线程：执行耗时操作（网络、IO、计算）
//    - 线程通信：通过 invoke_from_event_loop
//
// 2. 弱引用 (Weak)
//    - as_weak(): 创建弱引用
//    - upgrade(): 尝试转换为强引用
//    - 避免循环引用和内存泄漏
//
// 3. 回调机制
//    - Slint: callback name(types);
//    - Rust: app.on_name(|args| { ... });
//    - 双向通信的桥梁
//
// 4. 线程安全
//    - Slint UI 不是线程安全的
//    - 只能在主线程修改属性
//    - invoke_from_event_loop 保证线程安全
//
// 5. 状态管理模式
//    - is-loading: 加载状态
//    - status-message: 提示信息
//    - result-data: 结果数据
//    - 三个状态覆盖完整的异步生命周期
//
// 6. 替代方案
//    - tokio + slint: 使用 tokio 运行时
//    - async-std: 另一个异步运行时
//    - 信道 (channel): std::sync::mpsc 或 crossbeam
//
// 7. 错误处理
//    - 网络错误：设置 has-error 状态
//    - 超时：设置超时状态
//    - 取消：检查组件是否仍存在
// ============================================================================
