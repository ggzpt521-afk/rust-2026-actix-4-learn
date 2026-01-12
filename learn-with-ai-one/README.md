# Rust学习资源集合

这是一个系统化学习Rust编程语言的资源集合，包含了Rust的核心语法和功能的详细讲解与示例代码。所有文件按照学习顺序编号，建议从01开始逐步学习。

## 学习顺序建议

建议按照以下顺序学习，从基础概念逐步过渡到高级特性：

1. `01_variables.rs` - 变量系统
2. `02_data_types.rs` - 数据类型
3. `03_functions.rs` - 函数
4. `04_control_flow.rs` - 流程控制
5. `05_ownership.rs` - 所有权系统（Rust核心概念）
6. `06_structs.rs` - 结构体
7. `07_enums.rs` - 枚举和模式匹配
8. `08_collections.rs` - 常见集合及操作
9. `09_packages_modules.rs` - 包和模块系统
10. `10_error_handling.rs` - 错误处理机制
11. `11_generics.rs` - 泛型编程
12. `12_traits.rs` - Trait系统
13. `13_lifetimes.rs` - 生命周期概念
14. `14_std_lib_macros.rs` - 常用标准库函数与实用宏
15. `15_async_await.rs` - Rust异步编程（async/await）

## 文件内容概览

### 1. `01_variables.rs` - 变量系统
- 可变与不可变变量
- 变量声明与赋值
- 变量遮蔽（Shadowing）
- 常量（Constants）
- 作用域与变量生命周期
- 变量类型推导

### 2. `02_data_types.rs` - 数据类型
- 标量类型（整数、浮点数、布尔值、字符）
- 复合类型（元组、数组、切片）
- 类型转换（显式与隐式）
- 类型别名
- 复合类型的解构

### 3. `03_functions.rs` - 函数
- 函数定义与调用
- 参数与返回值
- 所有权与函数
- 嵌套函数
- 高阶函数
- 闭包（Closures）
- 递归函数
- 发散函数（Diverging Functions）
- 函数指针

### 4. `04_control_flow.rs` - 流程控制
- if条件表达式
- if let表达式
- match模式匹配
- for循环
- while循环
- loop循环
- break和continue
- 循环标签
- while let条件循环
- 控制流组合使用（如FizzBuzz游戏）

### 5. `05_ownership.rs` - 所有权系统
- 所有权规则
- 变量作用域
- 移动语义
- 克隆（Clone）
- 不可变与可变借用
- 借用规则
- 切片类型
- 悬垂引用
- 所有权与函数

### 6. `06_structs.rs` - 结构体
- 结构体定义
- 结构体实例化
- 结构体方法
- 关联函数
- 可见性
- 结构体更新语法
- 元组结构体
- 单元结构体
- 结构体与所有权

### 7. `07_enums.rs` - 枚举和模式匹配
- 枚举定义
- 带数据的枚举变体
- 枚举方法
- Option和Result枚举
- 自定义Result类型
- 嵌套枚举
- match表达式
- if let和while let
- 模式解构

### 8. `08_collections.rs` - 常见集合及操作
- Vec（动态数组）
- String（字符串）
- HashMap（哈希映射）
- HashSet（哈希集合）
- VecDeque（双端队列）
- LinkedList（链表）
- BTreeMap和BTreeSet
- 集合的创建与初始化
- 添加、删除和访问元素
- 集合的遍历
- 自定义类型在集合中的使用

### 9. `09_packages_modules.rs` - 包和模块系统
- 包（Package）、箱（Crate）和模块（Module）
- 模块定义
- 嵌套模块
- 路径（Path）
- 可见性控制
- use关键字
- super和self关键字
- 重导出（Re-export）
- 包结构
- 工作区（Workspace）

### 10. `10_error_handling.rs` - 错误处理机制
- 不可恢复错误（panic!）
- 可恢复错误（Result<T, E>）
- 错误传播（?运算符）
- 自定义错误类型
- 错误链
- 错误处理最佳实践
- 结合panic!和Result

### 11. `11_generics.rs` - 泛型编程
- 泛型函数
- 泛型结构体
- 泛型枚举
- 泛型约束
- 泛型与所有权
- 标准库中的泛型
- 泛型性能
- 关联类型
- 泛型与trait对象

### 12. `12_traits.rs` - Trait系统
- Trait定义
- Trait实现
- 默认方法
- Trait作为参数
- Trait作为返回值
- Trait约束
- 多Trait约束
- 泛型Trait
- 关联类型
- Trait对象
- 运算符重载
- 标准库中的Trait

### 13. `13_lifetimes.rs` - 生命周期概念
- 基本生命周期概念
- 函数中的生命周期
- 结构体中的生命周期
- 生命周期省略规则
- 静态生命周期
- 生命周期约束
- 生命周期子类型
- 生命周期与Trait对象
- 高级生命周期应用

### 14. `14_std_lib_macros.rs` - 常用标准库函数与实用宏
- Option和Result相关函数（map, and_then, unwrap_or等）
- 集合相关函数（sort, iter, filter, map, fold等）
- 字符串处理函数（trim, split, contains, replace等）
- 数学函数（三角函数, 指数, 对数等）
- 时间函数（Instant, Duration, SystemTime等）
- 实用宏（println!, format!, vec!, assert!等）

### 15. `15_async_await.rs` - Rust异步编程（async/await）
- 基本异步概念
- 异步函数和await表达式
- 异步块
- 异步错误处理
- 并发执行
- 超时处理
- 异步流（Stream）概念
- 异步编程实际应用示例

## 使用方法

1. 确保已安装Rust环境，可以通过[rustup](https://rustup.rs/)安装
2. 进入项目目录
3. 运行单个示例文件：
   ```bash
   rustc src/01_variables.rs && ./01_variables
   ```
   或使用Cargo运行（如果已初始化Cargo项目）：
   ```bash
   cargo run --bin 01_variables
   ```

4. 学习建议：
   - 仔细阅读代码中的注释
   - 运行示例并观察输出
   - 修改示例代码进行实验
   - 尝试自己实现类似功能

## 学习建议

1. **循序渐进**：从基础开始，逐步理解Rust的核心概念
2. **实践为主**：多写代码，尝试解决实际问题
3. **理解所有权**：这是Rust最独特的概念，需要花时间掌握
4. **利用文档**：Rust有非常优秀的官方文档，可以随时查阅
5. **加入社区**：参与Rust社区讨论，向他人学习
6. **不要气馁**：Rust有一定的学习曲线，遇到困难时多查阅资料和示例

## 参考资源

- [Rust官方文档](https://doc.rust-lang.org/)
- [Rust程序设计语言（中文版）](https://kaisery.github.io/trpl-zh-cn/)
- [Rust By Example](https://doc.rust-lang.org/rust-by-example/)

## 贡献

欢迎对本项目进行改进和扩展！如果你发现任何错误或有更好的示例，欢迎提交Pull Request。

---

祝学习愉快！Rust是一门强大而有趣的编程语言，掌握它将为你打开系统编程和安全编程的大门。