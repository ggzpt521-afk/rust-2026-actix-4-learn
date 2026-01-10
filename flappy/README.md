# Flappy Dragon

一个使用 Rust 和 bracket-lib 库实现的 Flappy Bird 风格游戏。

## 项目结构

```
flappy/
├── Cargo.toml          # 项目配置文件
├── Cargo.lock          # 依赖锁定文件
├── README.md           # 项目文档
├── src/
│   └── main.rs         # 游戏主代码
└── target/             # 编译输出目录
```

## 依赖

- **bracket-lib** (0.8.7): 一个用于创建 Roguelike 和终端风格游戏的 Rust 库

## bracket-lib 库介绍

### 什么是 bracket-lib？

[bracket-lib](https://github.com/amethyst/bracket-lib) 是一个 Rust 游戏开发库，最初名为 **RLTK (Roguelike Toolkit)**，专为 Roguelike 游戏开发设计，后来扩展为通用游戏开发库。它是《Hands-on Rust》一书的主要支持库。

### 设计理念

bracket-lib 的核心设计理念是**用户友好优先**。当开发者需要在性能和易用性之间做选择时，99% 的情况下会选择易用性。这使得它非常适合：
- 游戏开发教学
- 快速原型开发
- 初学者学习游戏编程

### 核心功能与原理

#### 1. 终端模拟 (Terminal Emulation)

```
┌─────────────────────────────────────────┐
│  bracket-lib 提供虚拟 ASCII 终端        │
│                                         │
│  @ . . . . . . . . . . . . . . . . │   │
│  . . . . . . . . . . . . . . . . . │   │
│  . . . . . . . . . . . . . . . . . │   │
│                                         │
│  使用 Codepage-437 字符集               │
│  支持瓦片图形和多图层                    │
└─────────────────────────────────────────┘
```

- 提供虚拟的 ASCII/Codepage-437 终端
- 支持可选的瓦片图形 (Tile Graphics)
- 支持多图层渲染
- 跨平台外观一致（使用精灵渲染字符）

#### 2. 游戏循环 (Game Loop)

```rust
// bracket-lib 的游戏循环模式
trait GameState {
    fn tick(&mut self, ctx: &mut BTerm);
}

// main_loop 会持续调用 tick() 方法
main_loop(context, State::new())
```

**工作原理**：
1. `main_loop` 函数接管程序控制权
2. 每一帧调用实现了 `GameState` trait 的 `tick()` 方法
3. 自动处理窗口事件（关闭、调整大小等）
4. 管理渲染和输入

#### 3. 渲染后端 (Rendering Backends)

bracket-lib 支持多种渲染后端，通过 feature flags 切换：

| 后端 | 说明 |
|------|------|
| OpenGL (默认) | 桌面平台默认使用 |
| WebGL | 编译为 WASM 时自动使用 |
| WebGPU | Vulkan/Metal/WGPU 支持 |
| Crossterm | 纯终端模式，无图形窗口 |

#### 4. 模块化架构

bracket-lib 由多个子 crate 组成，可按需选择：

| 子 crate | 功能 |
|----------|------|
| **bracket-terminal** | 终端模拟和渲染 |
| **bracket-random** | 随机数生成，支持 RPG 骰子语法（如 `3d6+12`） |
| **bracket-pathfinding** | A* 寻路算法和 Dijkstra 地图 |
| **bracket-geometry** | 几何图元（Point, Rect）、距离计算、Bresenham 算法 |
| **bracket-noise** | 噪声生成（Perlin、Simplex 等） |

### 本项目使用的 bracket-lib 功能

```rust
use bracket_lib::prelude::*;

// 1. BTerm - 终端上下文，提供渲染和输入功能
ctx.cls()              // 清屏
ctx.cls_bg(color)      // 用指定颜色清屏
ctx.set(x, y, fg, bg, glyph)  // 在指定位置绘制字符
ctx.print(x, y, text)  // 打印文本
ctx.print_centered(y, text)   // 居中打印

// 2. 颜色常量
RED, YELLOW, BLACK, NAVY  // 预定义颜色

// 3. 字符转换
to_cp437('@')          // 将字符转换为 Codepage-437 编码

// 4. 输入处理
ctx.key                // Option<VirtualKeyCode>，当前按下的键

// 5. 随机数
RandomNumberGenerator::new()
random.range(min, max) // 生成范围内随机数

// 6. 窗口构建
BTermBuilder::simple80x50()
    .with_title("title")
    .build()?
```

### 为什么选择 bracket-lib？

| 优点 | 说明 |
|------|------|
| 简单易用 | API 设计友好，适合初学者 |
| 跨平台 | 支持 Windows、macOS、Linux、Web |
| 轻量级 | 专注于 2D 游戏，无复杂依赖 |
| 文档完善 | 配套教程和书籍 |
| 快速原型 | 几行代码即可创建游戏窗口 |

### 参考资源

- [GitHub 仓库](https://github.com/amethyst/bracket-lib)
- [官方文档](https://docs.rs/bracket-lib)
- [使用指南](https://bfnightly.bracketproductions.com/bracket-lib/what_is_it.html)
- [crates.io](https://crates.io/crates/bracket-lib)

## 游戏原理

玩家控制一只"龙"（用 `@` 符号表示），通过按空格键来"拍打翅膀"向上飞行，同时受到重力影响会自动下落。玩家需要穿过不断出现的障碍物（管道），每成功穿过一个障碍物得 1 分。碰到障碍物或掉出屏幕则游戏结束。

## 核心设计

### 1. 状态机模式 (State Machine Pattern)

游戏使用状态机管理不同的游戏阶段：

```
┌─────────┐     按 P 键     ┌─────────┐
│  Menu   │ ──────────────> │ Playing │
└─────────┘                 └─────────┘
     ^                           │
     │        按 P 键            │ 碰撞/掉落
     └──────────────────────────┐│
                                vv
                           ┌─────────┐
                           │   End   │
                           └─────────┘
```

- **Menu**: 主菜单界面，等待玩家开始游戏
- **Playing**: 游戏进行中，处理玩家输入和游戏逻辑
- **End**: 游戏结束界面，显示分数并等待重新开始

### 2. 物理系统

#### 重力模型

```rust
// 重力加速度：每帧增加 0.2 的向下速度
if self.velocity < 2.0 {
    self.velocity += 0.2;  // 限制最大下落速度
}
self.y += self.velocity as i32;
```

#### 拍打机制

```rust
// 按空格键将速度设为 -2.0（向上）
fn flap(&mut self) {
    self.velocity = -2.0;
}
```

### 3. 坐标系统

游戏使用两套坐标系统：

| 坐标系 | 用途 | 说明 |
|--------|------|------|
| 世界坐标 | 游戏逻辑 | 玩家和障碍物的实际位置 |
| 屏幕坐标 | 渲染显示 | 相对于屏幕左上角的位置 |

**坐标转换公式**：`screen_x = world_x - player_x`

### 4. 碰撞检测

```
障碍物结构：
   │ │  <- 上半部分管道
   │ │

   ├─┤  <- 缺口 (gap_y ± half_size)

   │ │  <- 下半部分管道
   │ │
```

碰撞条件（同时满足）：
1. 玩家 x 坐标等于障碍物 x 坐标（水平重叠）
2. 玩家 y 坐标在缺口范围之外

### 5. 难度递增

障碍物缺口大小随分数增加而减小：

```rust
size = max(2, 20 - score)
```

- 初始缺口大小：20
- 最小缺口大小：2
- 每得一分，缺口减小 1

## 核心数据结构

### GameMode (游戏模式枚举)

```rust
enum GameMode {
    Menu,     // 主菜单
    Playing,  // 游戏中
    End,      // 游戏结束
}
```

### State (游戏主状态)

```rust
struct State {
    player: Player,       // 玩家对象
    frame_time: f32,      // 帧时间累加器
    mode: GameMode,       // 当前游戏模式
    obstacle: Obstacle,   // 当前障碍物
    score: i32,           // 玩家得分
}
```

### Player (玩家)

```rust
struct Player {
    x: i32,           // 世界 x 坐标（前进距离）
    y: i32,           // y 坐标（垂直位置）
    velocity: f32,    // 垂直速度
}
```

### Obstacle (障碍物)

```rust
struct Obstacle {
    x: i32,       // 世界 x 坐标
    gap_y: i32,   // 缺口中心 y 坐标
    size: i32,    // 缺口大小
}
```

## 核心函数

| 函数 | 所属 | 功能 |
|------|------|------|
| `Player::new(x, y)` | Player | 创建玩家，初始化位置和速度 |
| `Player::gravity_and_move()` | Player | 应用重力，更新位置 |
| `Player::flap()` | Player | 拍打翅膀，设置向上速度 |
| `Player::render(ctx)` | Player | 渲染玩家到屏幕 |
| `Obstacle::new(x, score)` | Obstacle | 创建障碍物，随机生成缺口位置 |
| `Obstacle::render(ctx, player_x)` | Obstacle | 渲染障碍物（坐标转换） |
| `Obstacle::hit_obstacle(player)` | Obstacle | 碰撞检测 |
| `State::new()` | State | 初始化游戏状态 |
| `State::play(ctx)` | State | 游戏主循环逻辑 |
| `State::restart()` | State | 重置游戏状态 |
| `State::main_menu(ctx)` | State | 显示主菜单 |
| `State::dead(ctx)` | State | 显示游戏结束界面 |
| `State::tick(ctx)` | GameState trait | 游戏循环入口，状态分发 |

## 游戏循环

```
main_loop() 调用 State::tick()
        │
        v
┌───────────────────────────────────┐
│            tick()                 │
│   match self.mode {               │
│     Menu -> main_menu()           │
│     Playing -> play()             │
│     End -> dead()                 │
│   }                               │
└───────────────────────────────────┘
        │
        v (Playing 模式)
┌───────────────────────────────────┐
│            play()                 │
│  1. 清屏 (cls_bg)                 │
│  2. 累积帧时间                     │
│  3. 固定时间步长更新逻辑           │
│     - gravity_and_move()          │
│  4. 处理输入 (空格 -> flap)       │
│  5. 渲染玩家和障碍物              │
│  6. 检测得分和碰撞                │
└───────────────────────────────────┘
```

## 运行游戏

```bash
# 编译并运行
cargo run

# 仅编译
cargo build

# 编译发布版本
cargo build --release
```

## 操作说明

| 按键 | 功能 |
|------|------|
| P | 开始游戏 / 重新开始 |
| Q | 退出游戏 |
| Space | 拍打翅膀（向上飞） |

## 游戏常量

| 常量 | 值 | 说明 |
|------|-----|------|
| SCREEN_WIDTH | 80 | 屏幕宽度（字符） |
| SCREEN_HEIGHT | 50 | 屏幕高度（字符） |
| FRAME_DURATION | 75.0 | 帧持续时间（ms），约 13 FPS |

## 技术要点

1. **固定时间步长**: 使用 `frame_time` 累加器确保游戏逻辑以固定频率更新，与渲染帧率解耦

2. **GameState Trait**: 实现 bracket-lib 的 `GameState` trait，通过 `tick()` 方法驱动游戏循环

3. **状态机模式**: 使用枚举和 match 表达式实现清晰的状态管理

4. **世界/屏幕坐标分离**: 障碍物使用世界坐标存储，渲染时转换为屏幕坐标，实现"无限"滚动效果

## Rust 语法解释

### Self 与 self

在 Rust 中，`Self` 和 `self` 是两个不同但相关的概念：

#### Self（大写）- 类型别名

`Self` 是当前类型的别名，在 `impl` 块中代表正在实现的类型。

```rust
impl Player {
    // Self 等价于 Player
    fn new(x: i32, y: i32) -> Self {  // 返回类型 Self = Player
        Player {                       // 也可以写成 Self { ... }
            x,
            y,
            velocity: 0.0,
        }
    }
}

// 以下两种写法完全等价：
fn new() -> Self { ... }
fn new() -> Player { ... }
```

**为什么使用 Self？**
- 代码更简洁，避免重复类型名
- 重构时更方便（修改类型名只需改一处）
- 在 trait 实现中，Self 自动指向实现该 trait 的类型

#### self（小写）- 实例引用

`self` 是方法的第一个参数，代表调用该方法的实例本身。

```rust
impl Player {
    // self 的三种形式：

    // 1. self - 获取所有权（消耗实例）
    fn consume(self) { ... }

    // 2. &self - 不可变借用（只读访问）
    fn read_only(&self) -> i32 {
        self.x  // 可以读取字段
    }

    // 3. &mut self - 可变借用（可以修改）
    fn modify(&mut self) {
        self.x += 1;  // 可以修改字段
    }
}
```

### &self vs &mut self 详解

| 形式 | 含义 | 能做什么 | 使用场景 |
|------|------|----------|----------|
| `self` | 获取所有权 | 消耗实例，调用后无法再使用 | 转换、销毁操作 |
| `&self` | 不可变借用 | 只能读取，不能修改 | 查询、获取信息 |
| `&mut self` | 可变借用 | 可以读取和修改 | 更新状态、修改数据 |

#### 本项目中的实际例子

```rust
impl Player {
    // 使用 &self - 只需要读取数据，不修改
    // 注意：hit_obstacle 在 Obstacle 中，这里用 &Player 作为参数
    fn hit_obstacle(&self, player: &Player) -> bool {
        // 只读取 self.x, self.gap_y 等，不修改任何值
        let does_x_match = player.x == self.x;
        // ...
    }

    // 使用 &mut self - 需要修改玩家状态
    fn gravity_and_move(&mut self) {
        self.velocity += 0.2;  // 修改速度
        self.y += self.velocity as i32;  // 修改位置
        self.x += 1;  // 修改位置
    }

    // 使用 &mut self - 需要修改速度
    fn flap(&mut self) {
        self.velocity = -2.0;  // 修改速度
    }

    // 使用 &mut self - 虽然只是渲染，但 BTerm API 需要
    fn render(&mut self, ctx: &mut BTerm) {
        ctx.set(0, self.y, YELLOW, BLACK, to_cp437('@'));
    }
}
```

### 方法调用语法糖

Rust 会自动处理引用，以下调用是等价的：

```rust
let mut player = Player::new(5, 25);

// 这两种写法等价：
player.flap();           // Rust 自动借用为 &mut player
Player::flap(&mut player);  // 显式传递引用

// 这两种写法也等价：
player.x;                // 自动解引用
(&player).x;             // 显式引用
```

### 所有权规则回顾

```
┌─────────────────────────────────────────────────────────┐
│                    Rust 所有权规则                       │
├─────────────────────────────────────────────────────────┤
│ 1. 每个值都有一个所有者（owner）                         │
│ 2. 同一时间只能有一个所有者                              │
│ 3. 当所有者离开作用域，值被丢弃（drop）                  │
├─────────────────────────────────────────────────────────┤
│                      借用规则                            │
├─────────────────────────────────────────────────────────┤
│ • 可以有多个不可变引用（&T）                             │
│ • 或者只有一个可变引用（&mut T）                         │
│ • 不能同时存在可变和不可变引用                           │
└─────────────────────────────────────────────────────────┘
```

### 为什么这样设计？

1. **内存安全**: 编译器在编译时检查所有借用，防止数据竞争和悬垂指针
2. **零成本抽象**: 借用检查在编译时完成，运行时无开销
3. **明确意图**: 从函数签名就能看出是否会修改数据

### 本项目中使用的其他 Rust 语法

#### 1. use 和 prelude 模式

```rust
use bracket_lib::prelude::*;
```

- `use` 关键字用于导入模块中的项
- `prelude` 是一种约定，包含库最常用的类型和函数
- `*` 表示导入模块中所有公开的项（glob import）

#### 2. const 常量

```rust
const SCREEN_WIDTH: i32 = 80;
```

- `const` 定义编译时常量，必须指定类型
- 命名约定：`UPPER_SNAKE_CASE`
- 常量会被内联到使用的地方，无运行时开销

#### 3. enum 枚举

```rust
enum GameMode {
    Menu,
    Playing,
    End,
}
```

- 枚举定义一组互斥的变体（variants）
- 每个变体可以携带不同类型的数据（本项目中未使用）
- 配合 `match` 使用实现状态机模式

#### 4. struct 结构体

```rust
struct Player {
    x: i32,
    y: i32,
    velocity: f32,
}
```

- 结构体将相关数据组合在一起
- 字段默认私有，通过方法访问

**字段简写语法**：
```rust
// 当变量名与字段名相同时可以简写
Player { x, y, velocity: 0.0 }
// 等价于
Player { x: x, y: y, velocity: 0.0 }
```

#### 5. impl 实现块

```rust
impl Player {
    fn new(x: i32, y: i32) -> Self { ... }  // 关联函数
    fn flap(&mut self) { ... }              // 方法
}
```

- `impl` 为类型添加方法和关联函数
- **关联函数**：没有 `self` 参数，通过 `Type::function()` 调用
- **方法**：第一个参数是 `self`，通过 `instance.method()` 调用

#### 6. trait 实现

```rust
impl GameState for State {
    fn tick(&mut self, ctx: &mut BTerm) { ... }
}
```

- `impl Trait for Type` 为类型实现特定接口
- 必须实现 trait 定义的所有必需方法
- 类似其他语言的接口（interface）

#### 7. match 模式匹配

```rust
match self.mode {
    GameMode::Menu => self.main_menu(ctx),
    GameMode::Playing => self.play(ctx),
    GameMode::End => self.dead(ctx),
}

match key {
    VirtualKeyCode::P => self.restart(),
    VirtualKeyCode::Q => ctx.quitting = true,
    _ => {}  // _ 匹配所有其他情况
}
```

- `match` 必须穷尽所有可能的情况
- `_` 是通配符，匹配任意值
- `=>` 后面是对应的执行代码
- `{}` 表示空操作（什么都不做）

#### 8. if let 语法糖

```rust
if let Some(VirtualKeyCode::Space) = ctx.key {
    self.player.flap();
}

// 等价于
match ctx.key {
    Some(VirtualKeyCode::Space) => self.player.flap(),
    _ => {}
}
```

- `if let` 是只关心一种模式时的简化写法
- 比完整的 `match` 更简洁
- 适合只处理 `Option` 的 `Some` 情况

#### 9. Option 类型

```rust
ctx.key  // 类型是 Option<VirtualKeyCode>
```

```rust
enum Option<T> {
    Some(T),  // 有值
    None,     // 无值
}
```

- Rust 没有 null，用 `Option` 表示可能缺失的值
- 强制处理"无值"的情况，避免空指针错误

#### 10. for 循环和 Range

```rust
for y in 0..self.gap_y - half_size {
    // y 从 0 到 (gap_y - half_size - 1)
}

for y in self.gap_y + half_size..SCREEN_HEIGHT {
    // y 从 (gap_y + half_size) 到 (SCREEN_HEIGHT - 1)
}
```

- `a..b` 创建一个左闭右开区间 [a, b)
- `a..=b` 创建一个闭区间 [a, b]（本项目未使用）
- `for` 循环自动处理迭代器

#### 11. as 类型转换

```rust
self.y += self.velocity as i32;
```

- `as` 用于基本类型之间的转换
- `f32 as i32` 会截断小数部分
- 注意：可能丢失精度或溢出

#### 12. format! 宏

```rust
&format!("Score {}", self.score)
```

- `format!` 返回格式化的 `String`
- `{}` 是占位符，按顺序替换
- `&format!(...)` 返回 `&String`，可自动转换为 `&str`

#### 13. ? 错误传播运算符

```rust
fn main() -> BError {
    let context = BTermBuilder::simple80x50()
        .with_title("flappy dragon")
        .build()?;  // 如果 build() 返回 Err，立即返回该错误

    main_loop(context, State::new())
}
```

- `?` 用于简化错误处理
- 如果是 `Ok(value)`，提取 value 继续执行
- 如果是 `Err(e)`，立即从函数返回该错误
- 要求函数返回类型为 `Result`

#### 14. 方法链（链式调用）

```rust
BTermBuilder::simple80x50()
    .with_title("flappy dragon")
    .build()?
```

- 每个方法返回 `self` 或新对象，支持连续调用
- 也叫建造者模式（Builder Pattern）
- 使配置代码更易读

#### 15. let 与 let mut

```rust
let screen_x = self.x - player_x;      // 不可变绑定
let mut random = RandomNumberGenerator::new();  // 可变绑定
```

- 默认 `let` 创建不可变绑定
- `let mut` 允许重新赋值
- Rust 鼓励默认不可变

#### 16. :: 路径分隔符

```rust
GameMode::Menu          // 枚举变体
Player::new(5, 25)      // 关联函数调用
VirtualKeyCode::Space   // 外部枚举变体
i32::max(2, 20 - score) // 类型的关联函数
```

- `::` 用于访问模块、类型、枚举的内部项
- `Type::function()` 调用关联函数
- `Enum::Variant` 访问枚举变体

#### 17. 逻辑运算符

```rust
// && 短路与：第一个为 false 则不计算第二个
does_x_match && (player_above_gap || player_below_gap)

// || 短路或：第一个为 true 则不计算第二个
self.player.y > SCREEN_HEIGHT || self.obstacle.hit_obstacle(&self.player)
```

#### 18. 文档注释

```rust
/// 这是文档注释，用于函数/结构体/模块
/// 支持 Markdown 格式
///
/// # 示例
/// ```
/// let player = Player::new(5, 25);
/// ```
fn new(x: i32, y: i32) -> Self { ... }

//! 这是模块级文档注释
//! 用于描述整个模块/文件
```

- `///` 文档注释，用于紧随其后的项
- `//!` 内部文档注释，用于当前模块
- 可用 `cargo doc` 生成 HTML 文档

## License

MIT
