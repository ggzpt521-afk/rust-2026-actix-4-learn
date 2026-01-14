// ============================================================================
// 03_event_handling.rs - Slint 事件处理示例
// ============================================================================
//
// 【核心概念】
// 事件处理是用户交互的基础，Slint 提供了声明式的事件处理机制：
// 1. 内置事件：clicked, edited, changed, pointer-event 等
// 2. 自定义回调：callback 声明
// 3. 事件冒泡：事件从子元素向父元素传播
//
// 【原理说明】
// Slint 的事件系统工作流程：
// 1. 用户操作 (点击、输入等) 产生原生事件
// 2. Slint 运行时捕获事件并确定目标元素
// 3. 调用元素上注册的事件回调
// 4. 回调中可以修改属性，触发 UI 更新
// ============================================================================

slint::slint! {
    // 导入标准组件
    import { Button, ScrollView } from "std-widgets.slint";

    export component EventHandling inherits Window {
        width: 400px;
        height: 350px;
        title: "事件处理示例";

        // ====================================================================
        // 组件状态属性
        // ====================================================================
        // 这些属性用于跟踪和显示事件信息
        in-out property <int> click-count: 0;      // 点击计数
        in-out property <string> event-log: "事件日志：\n";  // 事件日志

        Rectangle {
            width: 100%;
            height: 100%;
            background: #f0f0f0;

            VerticalLayout {
                padding: 20px;
                spacing: 10px;

                Text {
                    text: "事件处理示例";
                    font-size: 24px;
                    color: #333;
                }

                // ============================================================
                // 点击事件 (clicked)
                // ============================================================
                // clicked 是最常用的事件，当用户点击元素时触发
                //
                // 事件回调语法：
                // event-name => { 语句1; 语句2; ... }
                //
                // 回调体内可以：
                // - 修改属性值
                // - 调用其他回调
                // - 执行条件判断
                Button {
                    // 动态文本：显示当前点击次数
                    text: "点击我 (计数: " + click-count + ")";

                    // clicked 事件回调
                    clicked => {
                        // 增加点击计数
                        // 注意：click-count = click-count + 1 而不是 click-count++
                        // Slint 不支持 ++ 运算符
                        click-count = click-count + 1;

                        // 更新事件日志
                        // 字符串拼接使用 + 运算符
                        event-log = event-log + "按钮点击，计数: " + click-count + "\n";
                    }
                }

                // ============================================================
                // 鼠标悬停区域
                // ============================================================
                // Rectangle 支持鼠标相关事件：
                // - TouchArea: 专门用于处理触摸/鼠标事件的组件
                // - pointer-event: 低级指针事件
                //
                // 【注意】普通 Rectangle 默认不处理鼠标事件
                // 需要使用 TouchArea 或设置特定属性
                Rectangle {
                    width: 150px;
                    height: 60px;
                    background: #0066cc;

                    Text {
                        text: "鼠标悬停此处";
                        color: white;
                        vertical-alignment: center;
                        horizontal-alignment: center;
                    }

                    // 【扩展知识】如需处理鼠标悬停，可以使用 TouchArea：
                    // TouchArea {
                    //     width: 100%;
                    //     height: 100%;
                    //     // has-hover 属性：鼠标是否在区域内
                    //     // clicked => { ... }
                    //     // moved => { ... }
                    // }
                }

                // ============================================================
                // ScrollView - 滚动视图
                // ============================================================
                // ScrollView 是一个可滚动的容器
                // 当内容超出可视区域时，会显示滚动条
                //
                // 常用属性：
                // - viewport-width/height: 内容区域大小
                // - visible-width/height: 可视区域大小
                // - horizontal-scrollbar-policy: 水平滚动条策略
                // - vertical-scrollbar-policy: 垂直滚动条策略
                ScrollView {
                    width: 100%;
                    height: 120px;

                    Text {
                        // 绑定到 event-log 属性
                        // 每次事件日志更新，这里自动显示新内容
                        text: event-log;
                        font-size: 12px;
                        color: #333;
                    }
                }
            }
        }
    }
}

// ============================================================================
// main 函数
// ============================================================================
fn main() {
    let app = EventHandling::new().unwrap();

    // 设置初始事件日志
    // .into() 将 Rust 字符串转换为 Slint 的 SharedString
    app.set_event_log("事件日志：\n".into());

    app.run().unwrap();
}

// ============================================================================
// 【知识点总结】
// ============================================================================
//
// 1. 内置事件
//    - clicked: 点击事件 (Button, TouchArea 等)
//    - edited: 文本编辑完成 (LineEdit)
//    - changed: 值改变 (CheckBox, Slider 等)
//    - key-pressed: 键盘按下
//    - key-released: 键盘释放
//
// 2. 事件回调语法
//    event-name => { statements; }
//
// 3. TouchArea 组件
//    专门处理触摸/鼠标事件的组件：
//    - clicked: 点击
//    - moved: 移动
//    - pressed/released: 按下/释放
//    - has-hover: 悬停状态属性
//    - pressed-x/pressed-y: 按下位置
//    - mouse-x/mouse-y: 当前位置
//
// 4. 事件传播
//    - 默认情况下，事件不会冒泡
//    - TouchArea 可以设置 clicked-at 等来获取点击位置
//    - 可以通过父子组件的回调链实现类似冒泡的效果
//
// 5. 常用事件组件
//    - Button: clicked
//    - LineEdit: edited, accepted
//    - TouchArea: clicked, moved, scroll
//    - FocusScope: key-pressed, key-released
// ============================================================================
