// ============================================================================
// 06_list_rendering.rs - Slint 列表渲染示例
// ============================================================================
//
// 【核心概念】
// 列表渲染是动态 UI 的关键特性：
// 1. 根据数据数组动态生成 UI 元素
// 2. 数据变化时自动更新列表
// 3. 支持索引访问和条件渲染
//
// 【原理说明】
// Slint 的列表渲染使用 for-in 语法：
// - for item in array: Element { ... }
// - 为数组的每个元素创建一个 Element 实例
// - 数组变化时，Slint 会智能地增删元素（diff 算法）
// - 支持获取索引：for item[index] in array
// ============================================================================

slint::slint! {
    // 导入标准组件
    import { Button, ScrollView } from "std-widgets.slint";

    export component ListRendering inherits Window {
        width: 400px;
        height: 400px;
        title: "列表渲染示例";

        // ====================================================================
        // 数组类型属性
        // ====================================================================
        // <[T]>: 数组类型，T 是元素类型
        // 支持的数组元素类型：
        // - 基础类型：string, int, float, bool
        // - 结构体类型：struct { field1: type1, ... }
        // - 其他数组：[[int]] (二维数组)
        //
        // 数组字面量语法：[item1, item2, item3]
        in-out property <[string]> items: ["项目 1", "项目 2", "项目 3"];

        // 用于显示计数（数组长度的替代方案）
        in-out property <int> item-count: 3;

        VerticalLayout {
            padding: 20px;
            spacing: 15px;

            Text {
                text: "列表渲染示例";
                font-size: 24px;
                color: #333;
                horizontal-alignment: center;
            }

            // 显示项目数量
            Text {
                text: "共 " + item-count + " 个项目";
                font-size: 14px;
                color: #666;
                horizontal-alignment: center;
            }

            // 添加按钮
            Button {
                text: "添加项目";
                clicked => {
                    // 更新计数
                    // 注意：在 Slint UI 层面修改数组比较复杂
                    // 通常需要通过 Rust 端操作
                    item-count = item-count + 1;
                }
            }

            // ================================================================
            // ScrollView - 可滚动容器
            // ================================================================
            // 当列表内容超出可视区域时，提供滚动功能
            ScrollView {
                width: 100%;
                height: 200px;

                VerticalLayout {
                    spacing: 8px;
                    padding: 10px;

                    // ========================================================
                    // for-in 循环 - 列表渲染核心语法
                    // ========================================================
                    // 语法：for variable[index] in array: Element { ... }
                    //
                    // 解释：
                    // - item: 当前元素的值
                    // - [index]: 可选，当前元素的索引（从 0 开始）
                    // - items: 要遍历的数组
                    // - Rectangle { ... }: 为每个元素创建的 UI
                    //
                    // 【原理】
                    // 1. Slint 遍历 items 数组
                    // 2. 为每个元素创建一个 Rectangle 实例
                    // 3. item 和 index 在 Element 内可用
                    // 4. 当 items 变化时，自动更新 UI
                    for item[index] in items: Rectangle {
                        width: 100%;
                        height: 40px;

                        // 根据索引设置不同的背景色
                        // 展示条件表达式与索引的结合使用
                        background: index == 0 ? #e3f2fd :
                                   (index == 1 ? #f3e5f5 : #e8f5e9);

                        border-radius: 8px;

                        HorizontalLayout {
                            padding-left: 15px;
                            padding-right: 15px;

                            Text {
                                // 显示当前项的文本
                                // item 是 string 类型
                                text: item;
                                font-size: 16px;
                                color: #333;
                                vertical-alignment: center;
                            }
                        }
                    }
                }
            }

            // 说明文字
            Text {
                text: "提示: 列表使用 for...in 语法渲染";
                font-size: 12px;
                color: #999;
                horizontal-alignment: center;
            }
        }
    }
}

// ============================================================================
// main 函数
// ============================================================================
fn main() {
    let app = ListRendering::new().unwrap();

    // ------------------------------------------------------------------------
    // 从 Rust 操作列表数据
    // ------------------------------------------------------------------------
    // Slint 数组在 Rust 中对应 ModelRc<T> 类型
    // 常用操作：
    //
    // 1. 获取数组：
    //    let model = app.get_items();
    //
    // 2. 设置数组：
    //    use slint::VecModel;
    //    use std::rc::Rc;
    //    let items = Rc::new(VecModel::from(vec!["A".into(), "B".into()]));
    //    app.set_items(items.into());
    //
    // 3. 修改数组元素：
    //    如果使用 VecModel，可以调用 push(), remove(), set_row_data() 等
    //
    // 4. 监听数组变化：
    //    VecModel 实现了 Model trait，支持 row_count(), row_data() 等

    app.run().unwrap();
}

// ============================================================================
// 【知识点总结】
// ============================================================================
//
// 1. 数组类型
//    - <[T]>: 声明数组类型
//    - [item1, item2]: 数组字面量
//    - 支持基础类型和结构体类型
//
// 2. for-in 语法
//    - for item in array: Element { ... }
//    - for item[index] in array: Element { ... }
//    - 为每个元素创建 UI 实例
//
// 3. 在循环中使用变量
//    - item: 当前元素值
//    - index: 当前元素索引
//    - 可在 Element 的任何属性中使用
//
// 4. 条件与索引结合
//    - background: index == 0 ? color1 : color2
//    - 根据索引实现交替样式等效果
//
// 5. Rust 端数组操作
//    - ModelRc<T>: Slint 数组的 Rust 类型
//    - VecModel<T>: 可变数组实现
//    - 支持 push, remove, set_row_data 等操作
//
// 6. 性能考虑
//    - Slint 使用 diff 算法优化更新
//    - 只更新变化的元素
//    - 大列表建议使用虚拟化 (ListView)
//
// 7. 与其他框架对比
//    - Slint for-in ≈ React array.map() / Vue v-for
//    - Slint index ≈ React key / Vue :key
// ============================================================================
