// ============================================================================
// 04_state_management.rs - Slint 状态管理示例
// ============================================================================
//
// 【核心概念】
// 状态管理是 UI 应用的核心，涉及：
// 1. 应用状态的定义和存储
// 2. 状态变化如何反映到 UI
// 3. UI 操作如何改变状态
// 4. 复杂状态逻辑的组织
//
// 【原理说明】
// Slint 的状态管理基于属性系统：
// - 每个属性都是状态的一部分
// - 属性变化自动触发依赖它的 UI 元素更新
// - 条件表达式 (? :) 实现状态驱动的 UI 变化
// - 这种模式称为"单向数据流"或"响应式状态"
// ============================================================================

slint::slint! {
    // 导入按钮组件
    import { Button } from "std-widgets.slint";

    export component StateManagement inherits Window {
        width: 400px;
        height: 350px;
        title: "状态管理示例";

        // ====================================================================
        // 应用状态定义
        // ====================================================================
        // 所有状态都通过属性来定义和管理
        // 这些属性构成了应用的"状态模型"

        // 主题状态：控制应用的视觉外观
        // "light" 或 "dark"
        in-out property <string> theme: "light";

        // 语言状态：控制应用的语言设置
        in-out property <string> language: "zh-CN";

        // 通知状态：控制是否显示通知
        in-out property <bool> notifications: true;

        // 计数器状态：一个简单的数值状态
        in-out property <int> counter: 0;

        Rectangle {
            width: 100%;
            height: 100%;

            // ================================================================
            // 条件表达式 - 状态驱动的样式
            // ================================================================
            // 语法：condition ? value_if_true : value_if_false
            // 这是三元运算符，类似于大多数编程语言
            //
            // 【原理】当 theme 属性变化时：
            // 1. Slint 检测到 theme 被修改
            // 2. 重新计算所有依赖 theme 的表达式
            // 3. 如果结果不同，更新对应的 UI 属性
            // 4. 触发重绘
            background: theme == "light" ? #ffffff : #333333;

            VerticalLayout {
                padding: 20px;
                spacing: 10px;

                Text {
                    text: "状态管理示例";
                    font-size: 24px;
                    // 文字颜色也随主题变化
                    color: theme == "light" ? #333333 : #ffffff;
                }

                // ============================================================
                // 主题切换 - 状态修改示例
                // ============================================================
                HorizontalLayout {
                    spacing: 10px;

                    Text {
                        text: "主题:";
                        color: theme == "light" ? #333333 : #ffffff;
                        vertical-alignment: center;
                    }

                    // 切换到浅色主题
                    Button {
                        text: "浅色";
                        // 点击时修改状态
                        // theme = "light" 会触发所有依赖 theme 的 UI 更新
                        clicked => { theme = "light"; }
                    }

                    // 切换到深色主题
                    Button {
                        text: "深色";
                        clicked => { theme = "dark"; }
                    }
                }

                // ============================================================
                // 语言选择 - 另一个状态示例
                // ============================================================
                HorizontalLayout {
                    spacing: 10px;

                    Text {
                        text: "语言:";
                        color: theme == "light" ? #333333 : #ffffff;
                        vertical-alignment: center;
                    }

                    Button {
                        text: "中文";
                        clicked => { language = "zh-CN"; }
                    }

                    Button {
                        text: "English";
                        clicked => { language = "en-US"; }
                    }
                }

                // ============================================================
                // 计数器 - 数值状态示例
                // ============================================================
                HorizontalLayout {
                    spacing: 10px;

                    Button {
                        text: "-";
                        // 修改数值状态
                        // 注意空格：counter - 1 而非 counter-1
                        clicked => { counter = counter - 1 ; }
                    }

                    Text {
                        text: counter;
                        font-size: 20px;
                        width: 60px;
                        // 多个属性都可以依赖同一个状态
                        color: theme == "light" ? #333333 : #ffffff;
                        horizontal-alignment: center;
                        vertical-alignment: center;
                    }

                    Button {
                        text: "+";
                        clicked => { counter = counter + 1; }
                    }
                }

                // ============================================================
                // 状态组合显示
                // ============================================================
                // 多个状态可以组合在一个表达式中
                Text {
                    // 字符串拼接显示多个状态
                    text: "当前状态: " + theme + " / " + language;
                    font-size: 14px;
                    // 嵌套的条件表达式
                    color: theme == "light" ? #666666 : #cccccc;
                }
            }
        }
    }
}

// ============================================================================
// main 函数
// ============================================================================
fn main() {
    let app = StateManagement::new().unwrap();

    // ------------------------------------------------------------------------
    // 从 Rust 端读取和修改状态
    // ------------------------------------------------------------------------
    // 虽然这个示例主要在 Slint 端管理状态
    // 但 Rust 端也可以完全访问和控制这些状态

    // 读取状态示例 (已注释，仅作说明)：
    // let current_theme = app.get_theme();
    // let current_counter = app.get_counter();

    // 设置状态示例 (已注释，仅作说明)：
    // app.set_theme("dark".into());
    // app.set_counter(100);

    app.run().unwrap();
}

// ============================================================================
// 【知识点总结】
// ============================================================================
//
// 1. 状态定义
//    - 使用 property 声明状态
//    - 支持各种类型：string, int, bool, color 等
//    - 状态是响应式的，变化会自动触发 UI 更新
//
// 2. 条件表达式
//    - 语法：condition ? true_value : false_value
//    - 可用于任何属性值
//    - 支持嵌套：a ? b : (c ? d : e)
//
// 3. 状态驱动 UI
//    - UI 属性可以绑定到状态表达式
//    - 状态变化 -> 表达式重新计算 -> UI 更新
//    - 多个 UI 元素可以依赖同一状态
//
// 4. 状态修改
//    - 在事件回调中：property = new_value
//    - 从 Rust 端：app.set_property(value)
//    - 修改会自动触发所有依赖项更新
//
// 5. 最佳实践
//    - 状态集中管理（在组件顶层定义）
//    - 状态命名清晰（描述含义而非用途）
//    - 避免过多状态（考虑是否可以派生计算）
//    - 使用条件表达式而非重复代码
//
// 6. 与 React/Vue 的对比
//    - Slint 属性 ≈ React state / Vue ref
//    - Slint 条件表达式 ≈ React 条件渲染 / Vue v-if
//    - Slint 自动更新 ≈ React re-render / Vue 响应式
// ============================================================================
