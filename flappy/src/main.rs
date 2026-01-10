//! # Flappy Dragon 游戏
//!
//! 这是一个使用 Rust 和 bracket-lib 库实现的 Flappy Bird 风格游戏。
//!
//! ## 游戏原理
//!
//! 玩家控制一只"龙"（用 @ 符号表示），通过按空格键来"拍打翅膀"向上飞行，
//! 同时受到重力影响会自动下落。玩家需要穿过不断出现的障碍物（管道），
//! 每成功穿过一个障碍物得1分。碰到障碍物或掉出屏幕则游戏结束。
//!
//! ## 核心机制
//!
//! 1. **重力系统**: 玩家持续受到向下的加速度影响
//! 2. **拍打机制**: 按空格键给予玩家向上的速度
//! 3. **障碍物生成**: 随机生成带有缺口的管道，缺口大小随分数增加而减小
//! 4. **碰撞检测**: 检测玩家是否撞到障碍物或超出屏幕边界
//! 5. **状态机**: 游戏在菜单、游戏中、结束三种状态间切换

use bracket_lib::prelude::*;

// ============================================================================
// 游戏常量配置
// ============================================================================

/// 屏幕宽度（字符单位）
/// 游戏窗口横向可显示80个字符
const SCREEN_WIDTH: i32 = 80;

/// 屏幕高度（字符单位）
/// 游戏窗口纵向可显示50个字符
const SCREEN_HEIGHT: i32 = 50;

/// 帧持续时间（毫秒）
/// 控制游戏更新频率，值越大游戏越慢
/// 75ms 约等于 13 FPS 的游戏逻辑更新速度
const FRAME_DURATION: f32 = 75.0;

// ============================================================================
// 游戏状态枚举
// ============================================================================

/// 游戏模式枚举
///
/// 使用状态机模式管理游戏的不同阶段：
/// - Menu: 主菜单界面，等待玩家开始游戏
/// - Playing: 游戏进行中，处理玩家输入和游戏逻辑
/// - End: 游戏结束界面，显示分数并等待重新开始
enum GameMode {
    /// 主菜单状态
    Menu,
    /// 游戏进行中状态
    Playing,
    /// 游戏结束状态
    End,
}

// ============================================================================
// 游戏主状态结构体
// ============================================================================

/// 游戏主状态结构体
///
/// 存储游戏运行所需的所有状态数据：
/// - player: 玩家对象，包含位置和速度信息
/// - frame_time: 帧时间累加器，用于控制游戏逻辑更新频率
/// - mode: 当前游戏模式
/// - obstacle: 当前障碍物对象
/// - score: 玩家得分
struct State {
    /// 玩家对象
    player: Player,
    /// 帧时间累加器（毫秒）
    /// 用于实现固定时间步长的游戏循环
    frame_time: f32,
    /// 当前游戏模式
    mode: GameMode,
    /// 当前障碍物
    obstacle: Obstacle,
    /// 玩家得分
    score: i32,
}

// ============================================================================
// 障碍物结构体及实现
// ============================================================================

/// 障碍物结构体
///
/// 表示游戏中的管道障碍物，由上下两部分组成，中间有一个缺口供玩家通过。
///
/// ## 设计原理
///
/// 障碍物使用世界坐标系统（x 随玩家移动而相对变化），
/// 渲染时转换为屏幕坐标。缺口位置随机生成，
/// 缺口大小随游戏进行（分数增加）而逐渐减小，增加难度。
struct Obstacle {
    /// 障碍物的世界 x 坐标
    x: i32,
    /// 缺口中心的 y 坐标
    gap_y: i32,
    /// 缺口大小（半径的2倍）
    size: i32,
}

impl Obstacle {
    /// 创建新的障碍物
    ///
    /// # 参数
    ///
    /// * `x` - 障碍物的初始 x 坐标（世界坐标）
    /// * `score` - 当前分数，用于计算缺口大小
    ///
    /// # 返回值
    ///
    /// 返回一个新的 Obstacle 实例
    ///
    /// # 算法说明
    ///
    /// - 缺口 y 位置：在 10-50 范围内随机生成
    /// - 缺口大小：max(2, 20 - score)，最小为2，随分数增加而减小
    fn new(x: i32, score: i32) -> Self {
        let mut random = RandomNumberGenerator::new();
        Obstacle {
            x,
            gap_y: random.range(10, 50),
            size: i32::max(2, 20 - score),
        }
    }

    /// 渲染障碍物到屏幕
    ///
    /// # 参数
    ///
    /// * `ctx` - BTerm 上下文，用于绘制
    /// * `player_x` - 玩家的 x 坐标，用于计算屏幕坐标
    ///
    /// # 渲染原理
    ///
    /// 1. 计算屏幕坐标：screen_x = obstacle.x - player_x
    /// 2. 绘制上半部分管道：从 y=0 到 gap_y - half_size
    /// 3. 绘制下半部分管道：从 gap_y + half_size 到屏幕底部
    /// 4. 使用红色 '|' 字符表示管道
    fn render(&mut self, ctx: &mut BTerm, player_x: i32) {
        // 将世界坐标转换为屏幕坐标
        let screen_x = self.x - player_x;
        let half_size = self.size / 2;

        // 绘制上半部分管道（从顶部到缺口上边缘）
        for y in 0..self.gap_y - half_size {
            ctx.set(screen_x, y, RED, BLACK, to_cp437('|'));
        }

        // 绘制下半部分管道（从缺口下边缘到底部）
        for y in self.gap_y + half_size..SCREEN_HEIGHT {
            ctx.set(screen_x, y, RED, BLACK, to_cp437('|'));
        }
    }

    /// 检测玩家是否撞到障碍物
    ///
    /// # 参数
    ///
    /// * `player` - 玩家对象引用
    ///
    /// # 返回值
    ///
    /// 如果玩家与障碍物碰撞返回 true，否则返回 false
    ///
    /// # 碰撞检测原理
    ///
    /// 碰撞发生的条件（必须同时满足）：
    /// 1. 玩家 x 坐标等于障碍物 x 坐标（水平重叠）
    /// 2. 玩家 y 坐标在缺口范围之外（在缺口上方或下方）
    fn hit_obstacle(&self, player: &Player) -> bool {
        let half_size = self.size / 2;
        // 检查 x 坐标是否重叠
        let does_x_match = player.x == self.x;
        // 检查玩家是否在缺口上方
        let player_above_gap = player.y < self.gap_y - half_size;
        // 检查玩家是否在缺口下方
        let player_below_gap = player.y > self.gap_y + half_size;

        does_x_match && (player_above_gap || player_below_gap)
    }
}

// ============================================================================
// 玩家结构体及实现
// ============================================================================

/// 玩家结构体
///
/// 表示游戏中玩家控制的角色（龙/小鸟）。
///
/// ## 物理模型
///
/// 使用简化的物理模型：
/// - 位置 (x, y)：整数坐标，x 表示前进距离，y 表示高度
/// - 速度 (velocity)：浮点数，表示垂直方向速度
/// - 重力：每帧增加 0.2 的向下速度
/// - 拍打：将速度设为 -2.0（向上）
struct Player {
    /// 玩家世界 x 坐标（表示前进的距离）
    x: i32,
    /// 玩家 y 坐标（垂直位置，0 为顶部）
    y: i32,
    /// 垂直速度（正值向下，负值向上）
    velocity: f32,
}

impl Player {
    /// 创建新玩家
    ///
    /// # 参数
    ///
    /// * `x` - 初始 x 坐标
    /// * `y` - 初始 y 坐标
    ///
    /// # 返回值
    ///
    /// 返回一个新的 Player 实例，初始速度为 0
    fn new(x: i32, y: i32) -> Self {
        Player {
            x,
            y,
            velocity: 0.0,
        }
    }

    /// 渲染玩家到屏幕
    ///
    /// # 参数
    ///
    /// * `ctx` - BTerm 上下文
    ///
    /// # 说明
    ///
    /// 玩家始终显示在屏幕左侧 x=0 的位置，
    /// 使用黄色 '@' 字符表示
    fn render(&mut self, ctx: &mut BTerm) {
        ctx.set(0, self.y, YELLOW, BLACK, to_cp437('@'));
    }

    /// 应用重力并移动玩家
    ///
    /// # 物理计算原理
    ///
    /// 每次调用时执行以下操作：
    /// 1. 增加向下的速度（重力加速度 0.2），最大速度限制为 2.0
    /// 2. 将速度应用到 y 坐标（向下移动）
    /// 3. x 坐标增加 1（自动前进）
    /// 4. 如果 y < 0，将 y 设为 0（防止飞出屏幕顶部）
    ///
    /// 这实现了简单的抛物线运动效果
    fn gravity_and_move(&mut self) {
        // 应用重力加速度，但限制最大下落速度
        if self.velocity < 2.0 {
            self.velocity += 0.2;
        }
        // 将速度应用到位置
        self.y += self.velocity as i32;

        // 自动向前移动
        self.x += 1;

        // 防止飞出屏幕顶部
        if self.y < 0 {
            self.y = 0;
        }
    }

    /// 拍打翅膀（向上飞）
    ///
    /// # 说明
    ///
    /// 将垂直速度设为 -2.0，使玩家向上移动。
    /// 这会立即改变速度方向，模拟拍打翅膀的效果。
    fn flap(&mut self) {
        self.velocity = -2.0;
    }
}

// ============================================================================
// 游戏状态实现
// ============================================================================

impl State {
    /// 创建新的游戏状态
    ///
    /// # 返回值
    ///
    /// 返回初始化的游戏状态：
    /// - 玩家位于 (5, 25)
    /// - 第一个障碍物在屏幕右边缘
    /// - 游戏模式为菜单
    /// - 分数为 0
    fn new() -> Self {
        State {
            player: Player::new(5, 25),
            frame_time: 0.0,
            mode: GameMode::Menu,
            obstacle: Obstacle::new(SCREEN_WIDTH, 0),
            score: 0,
        }
    }

    /// 游戏主循环逻辑
    ///
    /// # 参数
    ///
    /// * `ctx` - BTerm 上下文
    ///
    /// # 游戏循环原理
    ///
    /// 1. **清屏**: 使用深蓝色背景
    /// 2. **时间控制**: 累积帧时间，达到阈值时更新游戏逻辑
    /// 3. **输入处理**: 检测空格键，触发拍打
    /// 4. **渲染**: 绘制玩家、障碍物、UI
    /// 5. **得分**: 玩家通过障碍物时加分
    /// 6. **碰撞检测**: 检测死亡条件
    fn play(&mut self, ctx: &mut BTerm) {
        // 清屏并设置背景色为深蓝色
        ctx.cls_bg(NAVY);

        // 累积帧时间
        self.frame_time += ctx.frame_time_ms;

        // 固定时间步长更新游戏逻辑
        // 只有当累积时间超过 FRAME_DURATION 时才更新
        if self.frame_time > FRAME_DURATION {
            self.frame_time = 0.0;
            self.player.gravity_and_move();
        }

        // 处理空格键输入 - 拍打翅膀
        if let Some(VirtualKeyCode::Space) = ctx.key {
            self.player.flap();
        }

        // 渲染玩家
        self.player.render(ctx);

        // 显示 UI 信息
        ctx.print(0, 0, "Press space to flap");
        ctx.print(0, 1, &format!("Score {}", self.score));

        // 渲染障碍物
        self.obstacle.render(ctx, self.player.x);

        // 检测是否通过障碍物并计分
        // 当玩家 x 坐标超过障碍物 x 坐标时，表示成功通过
        if self.player.x > self.obstacle.x {
            self.score += 1;
            // 生成新障碍物，位置在当前位置 + 屏幕宽度处
            self.obstacle = Obstacle::new(self.player.x + SCREEN_WIDTH, self.score);
        }

        // 死亡检测：
        // 1. 玩家掉出屏幕底部
        // 2. 玩家撞到障碍物
        if self.player.y > SCREEN_HEIGHT || self.obstacle.hit_obstacle(&self.player) {
            self.mode = GameMode::End;
        }
    }

    /// 重新开始游戏
    ///
    /// # 说明
    ///
    /// 重置所有游戏状态到初始值：
    /// - 切换到游戏模式
    /// - 重置帧时间
    /// - 重新创建玩家
    /// - 重新创建障碍物
    /// - 重置分数
    fn restart(&mut self) {
        self.mode = GameMode::Playing;
        self.frame_time = 0.0;
        self.player = Player::new(5, 25);
        self.obstacle = Obstacle::new(SCREEN_WIDTH, 0);
        self.score = 0;
    }

    /// 显示主菜单
    ///
    /// # 参数
    ///
    /// * `ctx` - BTerm 上下文
    ///
    /// # 说明
    ///
    /// 显示欢迎信息和操作提示：
    /// - P 键开始游戏
    /// - Q 键退出
    fn main_menu(&mut self, ctx: &mut BTerm) {
        ctx.cls();
        ctx.print_centered(5, "welcome here");
        ctx.print_centered(8, "(P) Play");
        ctx.print_centered(9, "(Q) Quit");

        // 处理菜单输入
        if let Some(key) = ctx.key {
            match key {
                VirtualKeyCode::P => self.restart(),
                VirtualKeyCode::Q => ctx.quitting = true,
                _ => {}
            }
        }
    }

    /// 显示死亡/游戏结束界面
    ///
    /// # 参数
    ///
    /// * `ctx` - BTerm 上下文
    ///
    /// # 说明
    ///
    /// 显示游戏结束信息、最终得分和操作提示
    fn dead(&mut self, ctx: &mut BTerm) {
        ctx.cls();
        ctx.print_centered(5, "You are dead");
        ctx.print_centered(6, &format!("you earned {} point", self.score));
        ctx.print_centered(8, "(P) Play");
        ctx.print_centered(9, "(Q) Quit");

        // 处理结束界面输入
        if let Some(key) = ctx.key {
            match key {
                VirtualKeyCode::P => self.restart(),
                VirtualKeyCode::Q => ctx.quitting = true,
                _ => {}
            }
        }
    }
}

// ============================================================================
// GameState trait 实现
// ============================================================================

/// 为 State 实现 bracket-lib 的 GameState trait
///
/// # 原理
///
/// bracket-lib 使用 GameState trait 来驱动游戏循环。
/// tick 方法在每一帧被调用，我们在其中根据当前游戏模式
/// 分发到对应的处理函数。
impl GameState for State {
    /// 游戏主循环入口
    ///
    /// # 参数
    ///
    /// * `ctx` - BTerm 上下文，提供输入处理和渲染功能
    ///
    /// # 状态机模式
    ///
    /// 使用 match 表达式根据当前游戏模式分发到对应处理函数：
    /// - Menu -> main_menu(): 显示主菜单
    /// - Playing -> play(): 执行游戏逻辑
    /// - End -> dead(): 显示结束界面
    fn tick(&mut self, ctx: &mut BTerm) {
        match self.mode {
            GameMode::Menu => self.main_menu(ctx),
            GameMode::End => self.dead(ctx),
            GameMode::Playing => self.play(ctx),
        }
    }
}

// ============================================================================
// 程序入口
// ============================================================================

/// 程序主入口
///
/// # 返回值
///
/// 返回 BError，bracket-lib 的错误类型
///
/// # 初始化流程
///
/// 1. 使用 BTermBuilder 创建 80x50 的终端窗口
/// 2. 设置窗口标题为 "flappy dragon"
/// 3. 调用 main_loop 启动游戏循环，传入初始游戏状态
///
/// # bracket-lib 游戏循环
///
/// main_loop 函数会：
/// 1. 持续调用 State::tick() 方法
/// 2. 处理窗口事件（关闭、调整大小等）
/// 3. 管理渲染和输入
fn main() -> BError {
    println!("Hello, world!");

    // 创建游戏窗口
    let context = BTermBuilder::simple80x50()
        .with_title("flappy dragon")
        .build()?;

    // 启动游戏主循环
    main_loop(context, State::new())
}
