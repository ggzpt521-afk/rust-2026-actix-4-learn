# Slint UI 框架学习教程

本项目是 Slint UI 框架的系统学习教程，包含 9 个循序渐进的示例，涵盖从基础组件到异步数据交互的完整知识体系。

> **版本说明**：本项目已升级至 Slint 1.x 版本，所有示例都经过更新以兼容最新的 Slint API。

---

## 目录

- [快速开始](#快速开始)
- [核心概念](#核心概念)
- [示例列表](#示例列表)
- [核心原理详解](#核心原理详解)
- [常用 API 参考](#常用-api-参考)
- [最佳实践](#最佳实践)

---

## 快速开始

### 环境要求

- Rust 1.70+
- Cargo

### 运行示例

```bash
# 在项目根目录执行

# 运行基础组件示例
cargo run --bin 01_basic_components

# 运行数据绑定示例
cargo run --bin 02_data_binding

# 运行事件处理示例
cargo run --bin 03_event_handling

# 运行状态管理示例
cargo run --bin 04_state_management

# 运行自定义组件示例
cargo run --bin 05_custom_components

# 运行列表渲染示例
cargo run --bin 06_list_rendering

# 运行布局样式示例
cargo run --bin 07_layout_styling

# 运行跨平台构建示例
cargo run --bin 08_cross_platform

# 运行异步数据交互示例
cargo run --bin 09_async_data
```

---

## 核心概念

### 1. 什么是 Slint？

Slint 是一个现代化的声明式 GUI 框架，专为桌面和嵌入式应用设计。

**核心特点：**
- **声明式 UI**：用 DSL 描述"是什么"而非"怎么做"
- **编译时安全**：UI 定义在编译时检查，避免运行时错误
- **原生性能**：编译为原生代码，无虚拟机开销
- **跨平台**：Windows、macOS、Linux、嵌入式设备

### 2. Slint 架构

```
┌─────────────────────────────────────────┐
│              应用程序                    │
├─────────────────────────────────────────┤
│   Slint DSL (UI 定义)  │  Rust (业务逻辑) │
├─────────────────────────────────────────┤
│           slint::slint! 宏              │
│         (编译时代码生成)                 │
├─────────────────────────────────────────┤
│           Slint 运行时                   │
│    (事件循环、渲染、属性系统)             │
├─────────────────────────────────────────┤
│         渲染后端 (可选)                  │
│   FemtoVG (OpenGL) / Skia / Software    │
├─────────────────────────────────────────┤
│            操作系统                      │
│     Windows / macOS / Linux             │
└─────────────────────────────────────────┘
```

### 3. slint::slint! 宏的工作原理

```rust
slint::slint! {
    export component MyApp inherits Window {
        // UI 定义
    }
}
```

**编译时处理流程：**
1. 宏解析 Slint DSL 代码
2. 生成对应的 Rust 结构体 (`MyApp`)
3. 生成 `new()`, `run()` 等方法
4. 生成属性的 getter/setter (`get_xxx()`, `set_xxx()`)
5. 生成回调处理器 (`on_xxx()`)

---

## 示例列表

| 示例 | 文件 | 核心知识点 |
|------|------|-----------|
| 01 | basic_components.rs | 基础组件、Window、Rectangle、Text |
| 02 | data_binding.rs | 属性系统、数据绑定、双向绑定 |
| 03 | event_handling.rs | 事件回调、clicked、edited |
| 04 | state_management.rs | 状态管理、条件表达式 |
| 05 | custom_components.rs | 自定义组件、组件复用 |
| 06 | list_rendering.rs | 列表渲染、for-in 循环 |
| 07 | layout_styling.rs | 布局容器、样式属性 |
| 08 | cross_platform.rs | 跨平台、条件编译 |
| 09 | async_data.rs | 异步处理、线程安全 |

---

## 核心原理详解

### 1. 属性系统 (Property System)

属性是 Slint 数据绑定的基础。每个属性都是响应式的——当属性值变化时，所有依赖它的 UI 元素会自动更新。

**属性修饰符对照表：**

| 修饰符 | 内部读 | 内部写 | 外部读 | 外部写 | 典型用途 |
|--------|--------|--------|--------|--------|----------|
| `in` | ✓ | ✗ | ✓ | ✓ | 接收父组件配置 |
| `out` | ✓ | ✓ | ✓ | ✗ | 向父组件暴露状态 |
| `in-out` | ✓ | ✓ | ✓ | ✓ | 双向数据同步 |
| `private` | ✓ | ✓ | ✗ | ✗ | 组件内部状态 |

**声明语法：**
```slint
// 输入属性：父组件设置，子组件只读
in property <string> title: "默认标题";

// 输出属性：子组件设置，父组件只读
out property <int> count: 0;

// 双向属性：父子都可读写
in-out property <bool> enabled: true;

// 私有属性：仅组件内部使用
property <color> bg-color: #ffffff;
```

**支持的数据类型：**

| 类型 | 说明 | 示例 |
|------|------|------|
| `string` | 字符串 | `"Hello"` |
| `int` | 整数 | `42` |
| `float` | 浮点数 | `3.14` |
| `bool` | 布尔值 | `true` / `false` |
| `color` | 颜色 | `#ff0000` / `red` |
| `length` | 长度 | `100px` / `50%` |
| `duration` | 时间 | `250ms` |
| `image` | 图像 | `@image-url("path")` |
| `[T]` | 数组 | `["a", "b", "c"]` |
| `struct` | 结构体 | 自定义结构 |

### 2. 数据绑定原理

**响应式更新机制：**

```
┌──────────────┐      修改       ┌──────────────┐
│   属性值     │ ─────────────→ │   依赖追踪   │
└──────────────┘                └──────────────┘
                                      │
                                      ↓ 通知
                                ┌──────────────┐
                                │  UI 重新计算  │
                                └──────────────┘
                                      │
                                      ↓ 更新
                                ┌──────────────┐
                                │   渲染更新   │
                                └──────────────┘
```

**单向绑定（属性 → UI）：**
```slint
Text {
    text: my-property;  // 属性变化时 UI 自动更新
}
```

**双向绑定（属性 ↔ UI）：**
```slint
LineEdit {
    text: user-input;                     // 属性 → UI
    edited => { user-input = self.text; } // UI → 属性
}
```

**计算属性：**
```slint
// 表达式绑定：依赖任一属性变化时重新计算
Text {
    text: "总计: " + (price * quantity);
    color: is-valid ? green : red;
}
```

### 3. 事件处理机制

**事件回调语法：**
```slint
// 基础语法
clicked => { 语句; }

// 多语句
clicked => {
    count = count + 1;
    update-display();
}

// 带参数（某些事件）
key-pressed(event) => {
    if event.text == "Enter" { submit(); }
}
```

**常用事件一览：**

| 组件 | 事件 | 触发时机 | 参数 |
|------|------|----------|------|
| Button | `clicked` | 按钮点击 | 无 |
| LineEdit | `edited` | 文本编辑完成 | 无 |
| LineEdit | `accepted` | 按下回车 | 无 |
| CheckBox | `toggled` | 选中状态变化 | 无 |
| TouchArea | `clicked` | 区域点击 | 无 |
| TouchArea | `moved` | 指针移动 | 无 |
| FocusScope | `key-pressed` | 键盘按下 | KeyEvent |
| FocusScope | `key-released` | 键盘释放 | KeyEvent |

### 4. 布局系统

**布局容器类型：**

| 容器 | 方向 | 主要属性 |
|------|------|----------|
| `VerticalLayout` | 垂直 | padding, spacing, alignment |
| `HorizontalLayout` | 水平 | padding, spacing, alignment |
| `GridLayout` | 网格 | spacing, columns |

**布局属性详解：**

```slint
VerticalLayout {
    // 内边距：容器边缘到内容的距离
    padding: 20px;
    padding-left: 10px;    // 单独设置某一边

    // 间距：子元素之间的距离
    spacing: 10px;

    // 对齐：子元素在容器中的位置
    alignment: center;     // start, center, end, stretch
}
```

**拉伸因子原理：**

```slint
HorizontalLayout {
    Rectangle { width: 100px; }            // 固定 100px
    Rectangle { horizontal-stretch: 1; }   // 占 1 份
    Rectangle { horizontal-stretch: 2; }   // 占 2 份
}
// 假设总宽度 400px，固定元素占 100px
// 剩余 300px 按 1:2 分配 → 100px 和 200px
```

### 5. 条件渲染

**条件表达式（用于属性值）：**
```slint
// 三元运算符
background: is-dark ? #333333 : #ffffff;

// 嵌套条件
color: level == "error" ? red :
       (level == "warn" ? orange : green);
```

**条件元素（用于渲染控制）：**
```slint
// 条件显示元素
if show-details: Text { text: "详情内容"; }

// 多条件分支
if tab == 0: Panel1 { }
if tab == 1: Panel2 { }
if tab == 2: Panel3 { }
```

### 6. 列表渲染

**基础语法：**
```slint
for item in items: Element {
    text: item;
}
```

**带索引：**
```slint
for item[index] in items: Rectangle {
    background: index % 2 == 0 ? #f0f0f0 : #ffffff;
    Text { text: index + ": " + item; }
}
```

**渲染原理：**
```
┌──────────────┐
│  数组数据    │ items = ["A", "B", "C"]
└──────────────┘
       │
       ↓ for-in 遍历
┌──────────────┐
│   创建实例   │ 为每个元素创建 UI 实例
└──────────────┘
       │
       ↓
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│  实例 0      │    │  实例 1      │    │  实例 2      │
│  item="A"    │    │  item="B"    │    │  item="C"    │
│  index=0     │    │  index=1     │    │  index=2     │
└──────────────┘    └──────────────┘    └──────────────┘
```

### 7. 组件系统

**组件定义：**
```slint
// 内部组件（不能从 Rust 访问）
component MyButton {
    in property <string> label: "按钮";
    in-out property <bool> pressed: false;
    callback clicked();

    Rectangle {
        // UI 实现
    }
}

// 导出组件（可从 Rust 访问）
export component MyApp inherits Window {
    // ...
}
```

**组件使用：**
```slint
// 实例化组件
MyButton {
    label: "提交";
    clicked => { handle-submit(); }
}

// 多实例
MyButton { label: "确定"; }
MyButton { label: "取消"; }
```

### 8. Rust 与 Slint 交互

**属性访问：**
```rust
// 获取属性值
let value = app.get_my_property();

// 设置属性值
app.set_my_property("新值".into());

// 注意：属性名中的 - 在 Rust 中变成 _
// Slint: my-property → Rust: my_property
```

**回调注册：**
```rust
// 注册回调处理器
app.on_my_callback(|arg1, arg2| {
    println!("回调被调用: {}, {}", arg1, arg2);
});

// 带返回值的回调
app.on_calculate(|a, b| {
    a + b  // 返回值
});
```

**线程安全更新：**
```rust
// 创建弱引用（避免循环引用）
let app_weak = app.as_weak();

// 在后台线程中
thread::spawn(move || {
    // 执行耗时操作...
    let result = fetch_data();

    // 回到主线程更新 UI
    slint::invoke_from_event_loop(move || {
        if let Some(app) = app_weak.upgrade() {
            app.set_result(result.into());
        }
    }).unwrap();
});
```

**为什么需要 `invoke_from_event_loop`：**
- Slint UI 不是线程安全的
- 只能在主线程修改属性
- 该函数将闭包调度到主线程事件循环执行

---

## 常用 API 参考

### Slint DSL 语法

```slint
// ============ 组件声明 ============
export component Name inherits Parent { }  // 导出组件
component Name { }                          // 内部组件

// ============ 属性声明 ============
in property <type> name: default;           // 输入属性
out property <type> name: default;          // 输出属性
in-out property <type> name: default;       // 双向属性
property <type> name: default;              // 私有属性

// ============ 回调声明 ============
callback name(param-types);                 // 无返回值
callback name(param-types) -> return-type;  // 有返回值

// ============ 事件处理 ============
event-name => { statements; }               // 处理事件

// ============ 条件渲染 ============
if condition: Element { }                   // 条件显示

// ============ 列表渲染 ============
for item in array: Element { }              // 基础循环
for item[index] in array: Element { }       // 带索引

// ============ 导入 ============
import { Widget } from "std-widgets.slint"; // 导入组件
```

### Rust API

```rust
// ============ 生命周期 ============
let app = MyComponent::new().unwrap();      // 创建实例
app.run().unwrap();                         // 运行应用

// ============ 属性操作 ============
app.get_property_name();                    // 获取属性
app.set_property_name(value);               // 设置属性

// ============ 回调操作 ============
app.on_callback_name(|args| { });           // 注册回调
app.invoke_callback_name(args);             // 调用回调

// ============ 引用管理 ============
let weak = app.as_weak();                   // 创建弱引用
weak.upgrade();                             // 升级为强引用

// ============ 线程操作 ============
slint::invoke_from_event_loop(|| { });      // 主线程执行
```

---

## 最佳实践

### 1. 组件设计原则

```
✓ 单一职责：每个组件只做一件事
✓ 可配置性：通过属性配置行为
✓ 可复用性：设计通用的组件接口
✓ 封装性：隐藏内部实现细节
```

### 2. 状态管理原则

```
✓ 集中管理：在组件顶层定义状态
✓ 最小化：只保存必要的状态
✓ 派生计算：能计算的不要存储
✓ 单一数据源：避免重复状态
```

### 3. 性能优化

```
✓ 避免重复计算：使用属性缓存结果
✓ 虚拟化列表：大列表使用 ListView
✓ 条件渲染：不需要的 UI 不渲染
✓ 减少重绘：批量更新属性
```

### 4. 异步处理

```
✓ 后台执行：耗时操作放后台线程
✓ 状态反馈：显示加载状态和进度
✓ 错误处理：处理所有可能的错误
✓ 取消支持：支持取消长时间操作
```

### 5. 跨平台开发

```
✓ 避免硬编码：不硬编码路径分隔符
✓ 条件编译：使用 #[cfg()] 处理平台差异
✓ 测试覆盖：在所有目标平台测试
✓ 响应式设计：适应不同屏幕尺寸
```

---

## 项目结构

```
slint-ai-one/
├── src/
│   ├── 01_basic_components.rs    # 基础组件示例（Rust + 内联 Slint）
│   ├── 02_data_binding.rs        # 数据绑定示例
│   ├── 03_event_handling.rs      # 事件处理示例
│   ├── 04_state_management.rs    # 状态管理示例
│   ├── 05_custom_components.rs   # 自定义组件示例
│   ├── 06_list_rendering.rs      # 列表渲染示例
│   ├── 07_layout_styling.rs      # 布局样式示例
│   ├── 08_cross_platform.rs      # 跨平台构建示例
│   ├── 09_async_data.rs          # 异步数据示例
│   ├── 01_basic_components.slint # 独立 UI 定义文件
│   ├── 02_data_binding.slint     # 数据绑定 UI
│   ├── 03_event_handling.slint   # 事件处理 UI
│   ├── 04_state_management.slint # 状态管理 UI
│   ├── 05_custom_components.slint# 自定义组件 UI
│   ├── 06_list_rendering.slint   # 列表渲染 UI
│   ├── 07_layout_styling.slint   # 布局样式 UI
│   ├── 08_cross_platform.slint   # 跨平台 UI
│   └── 09_async_data.slint       # 异步数据 UI
├── Cargo.toml                    # 项目配置
└── README.md                     # 本文档
```

---

## .slint 文件详解

### 1. .slint 文件是什么？

`.slint` 文件是 Slint UI 框架的声明式界面定义文件，使用 Slint DSL（领域特定语言）编写。

**核心特点：**
- **声明式**：描述"UI 是什么"而不是"如何创建 UI"
- **类型安全**：编译时检查，避免运行时错误
- **独立文件**：可以与 Rust 代码分离，便于 UI 设计师协作
- **可复用**：定义的组件可以在多处使用

### 2. .slint 与 .rs 文件的关系

```
┌─────────────────────────────────────────────────────────────┐
│                        应用程序                              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   ┌─────────────────┐          ┌─────────────────┐         │
│   │   .slint 文件    │◄────────►│    .rs 文件     │         │
│   │   (UI 定义)      │          │   (业务逻辑)    │         │
│   └─────────────────┘          └─────────────────┘         │
│          │                            │                     │
│          │  property                  │  get/set            │
│          │  callback                  │  on_xxx()           │
│          │  struct                    │  类型映射            │
│          │                            │                     │
│          ▼                            ▼                     │
│   ┌─────────────────────────────────────────────┐          │
│   │          Slint 编译器/运行时                  │          │
│   │   • 解析 .slint 生成 Rust 代码                │          │
│   │   • 管理属性响应式系统                        │          │
│   │   • 处理事件分发                              │          │
│   │   • 执行 UI 渲染                              │          │
│   └─────────────────────────────────────────────┘          │
└─────────────────────────────────────────────────────────────┘
```

### 3. 两种使用 .slint 的方式

**方式 1：内联 slint! 宏（本项目 .rs 文件采用）**

```rust
// 在 .rs 文件中直接写 Slint 代码
slint::slint! {
    export component MyApp inherits Window {
        Text { text: "Hello"; }
    }
}

fn main() {
    let app = MyApp::new().unwrap();
    app.run().unwrap();
}
```

优点：简单直接，适合小型项目和示例

**方式 2：独立 .slint 文件 + build.rs（推荐用于大型项目）**

```rust
// build.rs
fn main() {
    slint_build::compile("src/ui/main.slint").unwrap();
}

// main.rs
slint::include_modules!();

fn main() {
    let app = MyApp::new().unwrap();
    app.run().unwrap();
}
```

优点：UI 与逻辑分离，便于团队协作，IDE 支持更好

### 4. .slint 文件的核心语法

#### 组件声明
```slint
// 导出组件（可从 Rust 访问）
export component MyApp inherits Window { }

// 内部组件（仅 .slint 内部使用）
component MyButton { }

// 结构体定义
export struct TodoItem {
    id: int,
    title: string,
    completed: bool,
}
```

#### 属性声明
```slint
// 输入属性（Rust 可写，组件内只读）
in property <string> title: "默认";

// 输出属性（组件内可写，Rust 只读）
out property <int> count: 0;

// 双向属性（都可读写）
in-out property <bool> enabled: true;

// 私有属性（仅组件内部）
property <color> bg: #fff;
```

#### 回调声明
```slint
// 无返回值回调
callback clicked();
callback submit(string);

// 有返回值回调
callback calculate(int, int) -> int;
```

#### 函数定义
```slint
// 组件内部函数
function get-color() -> color {
    if count > 0 { return green; }
    else { return red; }
}
```

### 5. .slint 与 .rs 的通信桥梁

#### property（属性）—— 数据传递

```
┌─────────────────┐                    ┌─────────────────┐
│     .slint      │                    │      .rs        │
├─────────────────┤                    ├─────────────────┤
│                 │   get_xxx()        │                 │
│  in property    │ ◄──────────────────│  读取属性       │
│                 │   set_xxx()        │  设置属性       │
│                 │ ◄──────────────────│                 │
├─────────────────┤                    ├─────────────────┤
│                 │   get_xxx()        │                 │
│  out property   │ ──────────────────►│  读取状态       │
│                 │                    │                 │
├─────────────────┤                    ├─────────────────┤
│                 │   get/set_xxx()    │                 │
│  in-out         │ ◄────────────────►│  双向同步       │
│  property       │                    │                 │
└─────────────────┘                    └─────────────────┘
```

**Rust 代码示例：**
```rust
// .slint: in-out property <string> user-name;
// Rust 中属性名 - 变为 _

// 读取
let name = app.get_user_name();

// 设置
app.set_user_name("新用户".into());
```

#### callback（回调）—— 事件通知

```
┌─────────────────┐                    ┌─────────────────┐
│     .slint      │                    │      .rs        │
├─────────────────┤                    ├─────────────────┤
│                 │   on_xxx()         │                 │
│  callback xxx;  │ ──────────────────►│  注册处理器     │
│                 │                    │                 │
│  clicked => {   │                    │  app.on_xxx(    │
│    xxx();       │   触发             │    |args| {     │
│  }              │ ──────────────────►│      // 处理    │
│                 │                    │    }            │
│                 │                    │  );             │
└─────────────────┘                    └─────────────────┘
```

**Rust 代码示例：**
```rust
// .slint: callback fetch-data(string);

// 注册回调处理器
app.on_fetch_data(|query| {
    println!("收到查询: {}", query);
    // 执行异步操作...
});
```

#### struct（结构体）—— 复杂数据

```slint
// .slint 中定义
export struct WeatherData {
    city: string,
    temperature: int,
    description: string,
}

in-out property <WeatherData> weather;
```

```rust
// Rust 中使用
let weather = WeatherData {
    city: "北京".into(),
    temperature: 22,
    description: "晴朗".into(),
};
app.set_weather(weather);
```

### 6. .slint 文件示例对照表

| .slint 文件 | 对应 .rs 文件 | 核心演示内容 |
|------------|--------------|-------------|
| 01_basic_components.slint | 01_basic_components.rs | 基础组件：Text, Button, LineEdit, CheckBox |
| 02_data_binding.slint | 02_data_binding.rs | 属性绑定、双向绑定、表达式绑定 |
| 03_event_handling.slint | 03_event_handling.rs | 事件处理、callback 机制 |
| 04_state_management.slint | 04_state_management.rs | struct、状态管理、函数 |
| 05_custom_components.slint | 05_custom_components.rs | 自定义组件、属性接口、回调 |
| 06_list_rendering.slint | 06_list_rendering.rs | for-in 循环、数组操作 |
| 07_layout_styling.slint | 07_layout_styling.rs | 布局容器、样式属性、变量 |
| 08_cross_platform.slint | 08_cross_platform.rs | 平台信息传递、条件编译 |
| 09_async_data.slint | 09_async_data.rs | 异步状态、加载指示器 |

### 7. Slint 版本差异对照

| 功能 | Slint 0.x 语法 | Slint 1.x 语法 |
|------|---------------|---------------|
| 事件处理 | `on-clicked: { }` | `clicked => { }` |
| 文本编辑 | `on-text-changed: { }` | `edited => { }` |
| 条件表达式 | `if a then b else c` | `a ? b : c` |
| 减法运算 | `counter -= 1` | `counter = counter - 1` |
| 顶级组件 | `export component Name { }` | `export component Name inherits Window { }` |
| 布局对齐 | `horizontal-alignment: center` | `alignment: center`（在布局中）|
| 拉伸因子 | `stretch: 1` | `horizontal-stretch: 1` / `vertical-stretch: 1` |

> **注意**：本项目的 .slint 文件使用 Slint 0.x 语法作为参考示例，而 .rs 文件中的内联代码已更新为 Slint 1.x 语法。

---

## 参考资源

- [Slint 官方文档](https://slint.dev/docs)
- [Slint GitHub 仓库](https://github.com/slint-ui/slint)
- [Slint 示例集合](https://github.com/slint-ui/slint/tree/master/examples)
- [Rust 官方文档](https://doc.rust-lang.org/)

---

## 许可证

MIT License
