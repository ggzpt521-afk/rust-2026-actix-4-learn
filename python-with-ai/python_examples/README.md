# Python 学习示例集

一套系统覆盖 Python 核心概念的可运行示例文件，适合从入门到进阶的学习者。

## 学习顺序

建议按照以下顺序学习，每个文件都可以直接运行：

```bash
python 01_variables.py
```

### 第一阶段：基础语法

| 序号 | 文件 | 内容概要 |
|------|------|----------|
| 01 | `01_variables.py` | 变量定义、命名规范、常量约定、作用域、内存引用 |
| 02 | `02_types.py` | int/float/complex/bool/str/bytes、类型转换、数值运算 |
| 03 | `03_functions.py` | 函数定义、参数类型、多返回值、Lambda、闭包、递归 |
| 04 | `04_control_flow.py` | if/elif/else、for/while 循环、match-case（Python 3.10+）|

### 第二阶段：数据结构

| 序号 | 文件 | 内容概要 |
|------|------|----------|
| 05 | `05_collections.py` | list/tuple/set/dict、推导式、collections 模块 |
| 06 | `06_strings.py` | 字符串操作、格式化（f-string）、正则表达式 |

### 第三阶段：模块化编程

| 序号 | 文件 | 内容概要 |
|------|------|----------|
| 07 | `07_modules.py` | 模块导入、包结构、虚拟环境、pip 使用 |
| 08 | `08_classes.py` | 类定义、继承、多态、特殊方法、dataclass、抽象类 |
| 09 | `09_exceptions.py` | try/except/finally、自定义异常、异常链 |

### 第四阶段：文件与系统

| 序号 | 文件 | 内容概要 |
|------|------|----------|
| 10 | `10_files.py` | 文件读写、pathlib、JSON/CSV、临时文件 |
| 11 | `11_datetime.py` | date/time/datetime、时区处理、日期计算 |

### 第五阶段：高级特性

| 序号 | 文件 | 内容概要 |
|------|------|----------|
| 12 | `12_generators.py` | 迭代器协议、生成器函数/表达式、itertools |
| 13 | `13_decorators.py` | 函数装饰器、类装饰器、functools |
| 14 | `14_context_managers.py` | with 语句、contextlib、异步上下文管理器 |
| 15 | `15_type_hints.py` | 类型注解、泛型、Protocol、TypedDict |

### 第六阶段：并发编程

| 序号 | 文件 | 内容概要 |
|------|------|----------|
| 16 | `16_concurrency.py` | threading、multiprocessing、asyncio、concurrent.futures |

### 第七阶段：标准库与测试

| 序号 | 文件 | 内容概要 |
|------|------|----------|
| 17 | `17_stdlib.py` | os/pathlib/json/collections/functools/re/logging 等 |
| 18 | `18_builtins.py` | len/map/filter/zip/sorted/enumerate 等内置函数 |
| 19 | `19_testing.py` | unittest、pytest、mock、doctest、测试覆盖率 |

## 运行环境

- **Python 版本**: 3.10+ （推荐 3.11+，部分 match-case 示例需要）
- **操作系统**: macOS / Linux / Windows

## 快速开始

```bash
# 克隆或进入目录
cd python_examples

# 运行单个文件
python 01_variables.py

# 运行所有文件（检查是否有语法错误）
for f in *.py; do echo "=== $f ===" && python "$f" && echo; done

# 运行测试示例
python 19_testing.py
# 或使用 pytest
pip install pytest
pytest 19_testing.py -v
```

## 文件详细说明

### 01_variables.py - 变量与常量
- 变量定义与多重赋值
- 命名规范（PEP 8）
- 常量约定（`Final`）
- 变量作用域（`global`、`nonlocal`）
- 内存管理与引用

### 02_types.py - 基本数据类型
- 整数（无限精度、进制表示）
- 浮点数（精度问题、`Decimal`）
- 复数
- 布尔类型（真值测试）
- 字符串与字节类型

### 03_functions.py - 函数
- 位置参数、默认参数、关键字参数
- `*args` 和 `**kwargs`
- 多返回值（元组解包）
- Lambda 表达式
- 闭包与高阶函数
- 函数注解

### 04_control_flow.py - 流程控制
- 条件语句与三元表达式
- `for` 循环（`enumerate`、`zip`）
- `while` 循环
- `break`、`continue`、`else` 子句
- `match-case` 模式匹配（Python 3.10+）
- 海象运算符 `:=`（Python 3.8+）

### 05_collections.py - 组合数据类型
- 列表（方法、切片、推导式）
- 元组（不可变、命名元组）
- 集合（运算、frozenset）
- 字典（方法、合并）
- `collections` 模块（Counter、defaultdict、deque）

### 06_strings.py - 字符串与格式化
- 字符串方法（查找、替换、分割）
- f-string 格式化
- `format()` 方法
- 正则表达式（`re` 模块）
- Unicode 与编码

### 07_modules.py - 模块与包管理
- 导入方式（`import`、`from ... import`）
- 模块搜索路径
- 包结构与 `__init__.py`
- 相对导入与绝对导入
- 虚拟环境与 pip

### 08_classes.py - 类与对象
- 类定义、实例变量、类变量
- 访问控制（`_` 和 `__`）
- `@property` 装饰器
- `@classmethod` 和 `@staticmethod`
- 继承与多重继承
- 特殊方法（`__init__`、`__repr__`、`__add__` 等）
- `dataclass`
- 抽象基类（`ABC`）

### 09_exceptions.py - 异常处理
- `try`/`except`/`else`/`finally`
- 多异常处理
- 自定义异常
- 异常链（`raise ... from`）
- 断言（`assert`）
- 警告（`warnings`）

### 10_files.py - 文件与操作系统交互
- 文件读写模式
- `pathlib` 路径操作
- `os` 模块
- JSON/CSV 文件处理
- 临时文件
- 文件锁定

### 11_datetime.py - 日期时间处理
- `date`、`time`、`datetime` 对象
- `timedelta` 时间差
- 时区处理（`zoneinfo`）
- 日期格式化与解析
- `calendar` 模块

### 12_generators.py - 生成器与迭代器
- 迭代器协议（`__iter__`、`__next__`）
- 生成器函数（`yield`）
- 生成器表达式
- `yield from`
- `itertools` 模块

### 13_decorators.py - 装饰器
- 函数装饰器基础
- `@functools.wraps`
- 带参数的装饰器
- 类装饰器
- 实用装饰器（计时、缓存、重试）
- 内置装饰器（`@property`、`@staticmethod`）

### 14_context_managers.py - 上下文管理器
- `with` 语句
- `__enter__` 和 `__exit__`
- `@contextlib.contextmanager`
- `contextlib` 工具
- 异步上下文管理器
- `ExitStack`

### 15_type_hints.py - 类型注解
- 基本类型注解
- 容器类型（`list[int]`）
- `Optional` 和 `Union`
- `Callable`
- `TypeVar` 和泛型
- `Protocol`
- `Literal`、`Final`、`TypedDict`

### 16_concurrency.py - 并发与异步编程
- `threading` 多线程
- 线程同步（Lock、Semaphore、Event）
- `concurrent.futures`
- `asyncio` 基础
- 异步迭代器和生成器
- `multiprocessing` 多进程

### 17_stdlib.py - 常用标准库
- `os` 和 `pathlib`
- `json`
- `collections`
- `functools`
- `itertools`
- `re`（正则表达式）
- `hashlib`
- `random`
- `logging`
- `argparse`
- `dataclasses`

### 18_builtins.py - 内置函数
- 类型转换（`int`、`str`、`list` 等）
- 数学函数（`sum`、`min`、`max`、`abs`）
- 序列操作（`sorted`、`reversed`、`enumerate`、`zip`）
- 函数式（`map`、`filter`、`any`、`all`）
- 对象属性（`getattr`、`hasattr`、`isinstance`）

### 19_testing.py - 单元测试
- `unittest` 框架
- 断言方法
- `setUp` 和 `tearDown`
- Mock 对象
- `doctest`
- pytest 风格
- 测试覆盖率

## 代码风格

所有示例代码遵循：
- **PEP 8** 代码风格指南
- **PEP 257** 文档字符串规范
- 使用 **中文注释** 便于理解
- **Pythonic** 写法

## 扩展学习

完成本教程后，建议继续学习：

1. **Web 开发**: Flask / Django / FastAPI
2. **数据科学**: NumPy / Pandas / Matplotlib
3. **机器学习**: Scikit-learn / PyTorch / TensorFlow
4. **自动化**: Selenium / Scrapy / Airflow
5. **深入理解**: CPython 源码 / Python 数据模型

## 许可证

本示例代码仅供学习使用，可自由修改和分发。
