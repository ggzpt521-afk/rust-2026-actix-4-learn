// ============================================================================
// 07_layout_styling.rs - Slint 布局与样式示例
// ============================================================================
//
// 【核心概念】
// 布局和样式是 UI 开发的基础：
// 1. 布局：控制元素的位置和大小
// 2. 样式：控制元素的外观（颜色、边框、阴影等）
// 3. 响应式：适应不同屏幕尺寸
//
// 【原理说明】
// Slint 提供了多种布局方式：
// - VerticalLayout: 垂直排列（从上到下）
// - HorizontalLayout: 水平排列（从左到右）
// - GridLayout: 网格布局
// - 绝对定位：使用 x, y 属性
//
// 布局系统基于约束求解：
// - 每个元素有首选大小、最小/最大大小
// - 布局容器根据约束分配空间
// - stretch 属性控制空间分配比例
// ============================================================================

slint::slint! {
    import { Button } from "std-widgets.slint";

    export component LayoutStyling inherits Window {
        width: 500px;
        height: 450px;
        title: "布局样式示例";

        // 选项卡状态，用于切换不同的演示
        in-out property <int> active-tab: 0;

        VerticalLayout {
            padding: 20px;
            spacing: 15px;

            Text {
                text: "布局与样式示例";
                font-size: 24px;
                color: #333;
                // horizontal-alignment: 文本水平对齐方式
                // - left: 左对齐
                // - center: 居中
                // - right: 右对齐
                horizontal-alignment: center;
            }

            // ================================================================
            // 选项卡导航
            // ================================================================
            HorizontalLayout {
                spacing: 8px;
                // alignment: 子元素在布局中的对齐方式
                // - start: 起始位置
                // - center: 居中
                // - end: 结束位置
                // - stretch: 拉伸填满
                alignment: center;

                Button {
                    text: "水平布局";
                    clicked => { active-tab = 0; }
                }
                Button {
                    text: "垂直布局";
                    clicked => { active-tab = 1; }
                }
                Button {
                    text: "样式演示";
                    clicked => { active-tab = 2; }
                }
            }

            // 内容区域
            Rectangle {
                width: 100%;
                height: 280px;
                background: #f5f5f5;
                border-radius: 8px;

                // ============================================================
                // 条件渲染 - if 语法
                // ============================================================
                // if condition: Element { ... }
                // 当条件为真时渲染该元素

                // ============================================================
                // 水平布局演示
                // ============================================================
                if active-tab == 0: HorizontalLayout {
                    padding: 20px;
                    spacing: 15px;

                    // 固定宽度元素
                    Rectangle {
                        width: 80px;
                        height: 80px;
                        background: #0066cc;
                        border-radius: 8px;

                        Text {
                            text: "固定";
                            color: white;
                            vertical-alignment: center;
                            horizontal-alignment: center;
                        }
                    }

                    // --------------------------------------------------------
                    // horizontal-stretch - 水平拉伸因子
                    // --------------------------------------------------------
                    // 控制元素在水平方向上占用多少剩余空间
                    // stretch: 1 表示占用 1 份
                    // 多个元素的 stretch 值决定了它们的比例
                    Rectangle {
                        height: 80px;
                        background: #009900;
                        border-radius: 8px;
                        // 拉伸因子 1：占用剩余空间的 1/(1+2) = 1/3
                        horizontal-stretch: 1;

                        Text {
                            text: "拉伸 1";
                            color: white;
                            vertical-alignment: center;
                            horizontal-alignment: center;
                        }
                    }

                    Rectangle {
                        height: 80px;
                        background: #ff9900;
                        border-radius: 8px;
                        // 拉伸因子 2：占用剩余空间的 2/(1+2) = 2/3
                        horizontal-stretch: 2;

                        Text {
                            text: "拉伸 2";
                            color: white;
                            vertical-alignment: center;
                            horizontal-alignment: center;
                        }
                    }
                }

                // ============================================================
                // 垂直布局演示
                // ============================================================
                if active-tab == 1: VerticalLayout {
                    padding: 20px;
                    spacing: 15px;

                    // 固定高度元素
                    Rectangle {
                        width: 120px;
                        height: 60px;
                        background: #0066cc;
                        border-radius: 8px;

                        Text {
                            text: "固定高度";
                            color: white;
                            vertical-alignment: center;
                            horizontal-alignment: center;
                        }
                    }

                    // vertical-stretch: 垂直拉伸因子
                    // 与 horizontal-stretch 类似，但用于垂直方向
                    Rectangle {
                        width: 120px;
                        background: #009900;
                        border-radius: 8px;
                        vertical-stretch: 1;

                        Text {
                            text: "拉伸";
                            color: white;
                            vertical-alignment: center;
                            horizontal-alignment: center;
                        }
                    }
                }

                // ============================================================
                // 样式演示
                // ============================================================
                if active-tab == 2: VerticalLayout {
                    padding: 20px;
                    spacing: 15px;

                    // --------------------------------------------------------
                    // 边框样式
                    // --------------------------------------------------------
                    // border-width: 边框宽度
                    // border-color: 边框颜色
                    // border-radius: 圆角半径
                    Rectangle {
                        width: 120px;
                        height: 50px;
                        background: white;
                        border-width: 2px;
                        border-color: #0066cc;
                        border-radius: 8px;

                        Text {
                            text: "边框样式";
                            color: #333;
                            vertical-alignment: center;
                            horizontal-alignment: center;
                        }
                    }

                    // --------------------------------------------------------
                    // 阴影效果
                    // --------------------------------------------------------
                    // drop-shadow-offset-x/y: 阴影偏移
                    // drop-shadow-blur: 阴影模糊半径
                    // drop-shadow-color: 阴影颜色
                    Rectangle {
                        width: 120px;
                        height: 50px;
                        background: white;
                        border-radius: 8px;
                        drop-shadow-offset-y: 4px;
                        drop-shadow-blur: 8px;
                        // 颜色格式：#RRGGBBAA (带透明度)
                        drop-shadow-color: #00000033;

                        Text {
                            text: "阴影效果";
                            color: #333;
                            vertical-alignment: center;
                            horizontal-alignment: center;
                        }
                    }

                    // --------------------------------------------------------
                    // 圆形
                    // --------------------------------------------------------
                    // border-radius 设为宽高的一半，形成圆形
                    Rectangle {
                        width: 80px;
                        height: 80px;
                        background: #0066cc;
                        // 圆角 = 宽度/2 = 圆形
                        border-radius: 40px;

                        Text {
                            text: "圆形";
                            color: white;
                            vertical-alignment: center;
                            horizontal-alignment: center;
                        }
                    }
                }
            }

            Text {
                text: "使用拉伸因子和样式属性实现自适应布局";
                font-size: 12px;
                color: #666;
                horizontal-alignment: center;
            }
        }
    }
}

// ============================================================================
// main 函数
// ============================================================================
fn main() {
    let app = LayoutStyling::new().unwrap();
    app.run().unwrap();
}

// ============================================================================
// 【知识点总结】
// ============================================================================
//
// 1. 布局容器
//    - VerticalLayout: 垂直排列
//    - HorizontalLayout: 水平排列
//    - GridLayout: 网格布局（columns 属性设置列数）
//
// 2. 布局属性
//    - padding: 内边距
//    - spacing: 子元素间距
//    - alignment: 子元素对齐方式
//
// 3. 拉伸因子
//    - horizontal-stretch: 水平方向拉伸比例
//    - vertical-stretch: 垂直方向拉伸比例
//    - 用于响应式布局，分配剩余空间
//
// 4. 尺寸约束
//    - width/height: 固定尺寸
//    - min-width/min-height: 最小尺寸
//    - max-width/max-height: 最大尺寸
//    - preferred-width/height: 首选尺寸
//
// 5. 边框与圆角
//    - border-width: 边框宽度
//    - border-color: 边框颜色
//    - border-radius: 圆角半径
//
// 6. 阴影
//    - drop-shadow-offset-x/y: 阴影偏移
//    - drop-shadow-blur: 模糊半径
//    - drop-shadow-color: 阴影颜色
//
// 7. 条件渲染
//    - if condition: Element { ... }
//    - 条件为真时渲染元素
//    - 用于选项卡、折叠面板等
//
// 8. 响应式设计技巧
//    - 使用 % 单位相对于父元素
//    - 使用拉伸因子分配空间
//    - 设置最小/最大尺寸限制
// ============================================================================
