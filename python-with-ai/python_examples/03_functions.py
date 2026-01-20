#!/usr/bin/env python3
"""
03_functions.py - 函数

Python 函数定义、参数类型、返回值、作用域等核心概念。
包括多返回值、可变参数、关键字参数、Lambda 表达式等。
"""

# =============================================================================
# 1. 基本函数定义
# =============================================================================

print("=== 基本函数定义 ===")


def greet(name):
    """
    简单的问候函数。

    Args:
        name: 要问候的名字

    Returns:
        问候字符串
    """
    return f"Hello, {name}!"


result = greet("Python")
print(result)
print(f"函数文档: {greet.__doc__.strip()}")

# =============================================================================
# 2. 参数类型
# =============================================================================

print("\n=== 参数类型 ===")


# 位置参数
def add(a, b):
    return a + b


print(f"位置参数: add(3, 5) = {add(3, 5)}")


# 默认参数（注意：默认值只计算一次）
def power(base, exp=2):
    return base ** exp


print(f"默认参数: power(3) = {power(3)}")
print(f"默认参数: power(2, 10) = {power(2, 10)}")


# 默认参数的陷阱（可变对象作为默认值）
def append_item_wrong(item, lst=[]):  # 不推荐！
    lst.append(item)
    return lst


def append_item_correct(item, lst=None):  # 正确做法
    if lst is None:
        lst = []
    lst.append(item)
    return lst


print("\n默认参数陷阱:")
print(f"  错误方式第一次: {append_item_wrong(1)}")
print(f"  错误方式第二次: {append_item_wrong(2)}")  # [1, 2] 而不是 [2]！
print(f"  正确方式第一次: {append_item_correct(1)}")
print(f"  正确方式第二次: {append_item_correct(2)}")  # [2]


# 关键字参数
def describe_person(name, age, city):
    return f"{name} is {age} years old, lives in {city}"


print(f"\n关键字参数: {describe_person(age=25, city='Beijing', name='Alice')}")


# 强制关键字参数（* 后面的参数必须用关键字指定）
def make_request(url, *, method="GET", timeout=30):
    return f"{method} {url} (timeout={timeout}s)"


print(f"强制关键字: {make_request('https://api.example.com', method='POST')}")


# 仅位置参数（Python 3.8+，/ 前面的参数必须用位置方式传递）
def calculate(x, y, /, *, operation="add"):
    if operation == "add":
        return x + y
    return x - y


print(f"仅位置参数: {calculate(10, 5, operation='sub')}")

# =============================================================================
# 3. 可变参数
# =============================================================================

print("\n=== 可变参数 ===")


# *args - 收集位置参数为元组
def sum_all(*args):
    """接受任意数量的位置参数"""
    print(f"  args = {args}, type = {type(args).__name__}")
    return sum(args)


print(f"*args: sum_all(1, 2, 3, 4, 5) = {sum_all(1, 2, 3, 4, 5)}")


# **kwargs - 收集关键字参数为字典
def print_info(**kwargs):
    """接受任意数量的关键字参数"""
    print(f"  kwargs = {kwargs}, type = {type(kwargs).__name__}")
    for key, value in kwargs.items():
        print(f"    {key}: {value}")


print("\n**kwargs:")
print_info(name="Bob", age=30, language="Python")


# 组合使用
def universal_function(required, *args, default="default", **kwargs):
    """演示所有参数类型的组合"""
    print(f"  required = {required}")
    print(f"  args = {args}")
    print(f"  default = {default}")
    print(f"  kwargs = {kwargs}")


print("\n组合参数:")
universal_function("必需", 1, 2, 3, default="自定义", extra="额外")

# =============================================================================
# 4. 参数解包
# =============================================================================

print("\n=== 参数解包 ===")

# 列表/元组解包为位置参数
numbers = [1, 2, 3]
print(f"列表解包: add(*[3, 5]) = {add(*[3, 5])}")

# 字典解包为关键字参数
person = {"name": "Charlie", "age": 28, "city": "Shanghai"}
print(f"字典解包: {describe_person(**person)}")

# =============================================================================
# 5. 多返回值
# =============================================================================

print("\n=== 多返回值 ===")


def get_stats(numbers):
    """计算统计信息，返回多个值"""
    if not numbers:
        return None, None, None, 0

    total = sum(numbers)
    avg = total / len(numbers)
    min_val = min(numbers)
    max_val = max(numbers)

    return min_val, max_val, avg, len(numbers)  # 返回元组


data = [23, 45, 12, 67, 34, 89, 11]

# 解包多返回值
minimum, maximum, average, count = get_stats(data)
print(f"数据: {data}")
print(f"最小值: {minimum}, 最大值: {maximum}")
print(f"平均值: {average:.2f}, 数量: {count}")

# 也可以作为元组接收
result = get_stats(data)
print(f"作为元组: {result}")

# 忽略部分返回值
_, max_only, *rest = get_stats(data)
print(f"只要最大值: {max_only}")


# 使用命名元组返回多值（更清晰）
from collections import namedtuple

Stats = namedtuple("Stats", ["min", "max", "avg", "count"])


def get_stats_named(numbers):
    """使用命名元组返回统计信息"""
    if not numbers:
        return Stats(None, None, None, 0)
    return Stats(
        min=min(numbers),
        max=max(numbers),
        avg=sum(numbers) / len(numbers),
        count=len(numbers),
    )


stats = get_stats_named(data)
print(f"\n命名元组: {stats}")
print(f"访问属性: stats.avg = {stats.avg:.2f}")

# =============================================================================
# 6. Lambda 表达式
# =============================================================================

print("\n=== Lambda 表达式 ===")

# 匿名函数
square = lambda x: x ** 2
print(f"lambda square(5) = {square(5)}")

# 多参数
multiply = lambda x, y: x * y
print(f"lambda multiply(3, 4) = {multiply(3, 4)}")

# 常与高阶函数配合使用
numbers = [1, 2, 3, 4, 5]
squared = list(map(lambda x: x ** 2, numbers))
print(f"map with lambda: {squared}")

# 排序时使用
students = [("Alice", 85), ("Bob", 92), ("Charlie", 78)]
by_score = sorted(students, key=lambda s: s[1], reverse=True)
print(f"按分数排序: {by_score}")

# =============================================================================
# 7. 高阶函数
# =============================================================================

print("\n=== 高阶函数 ===")


# 函数作为参数
def apply_operation(func, x, y):
    """接受函数作为参数"""
    return func(x, y)


print(f"apply_operation(add, 10, 20) = {apply_operation(add, 10, 20)}")
print(f"apply_operation(lambda x, y: x * y, 10, 20) = {apply_operation(lambda x, y: x * y, 10, 20)}")


# 函数作为返回值
def make_multiplier(n):
    """返回一个乘法函数（闭包）"""
    def multiplier(x):
        return x * n
    return multiplier


double = make_multiplier(2)
triple = make_multiplier(3)

print(f"double(5) = {double(5)}")
print(f"triple(5) = {triple(5)}")

# =============================================================================
# 8. 闭包
# =============================================================================

print("\n=== 闭包 ===")


def counter():
    """创建一个计数器闭包"""
    count = 0

    def increment():
        nonlocal count  # 声明使用外层变量
        count += 1
        return count

    return increment


counter1 = counter()
counter2 = counter()

print(f"counter1: {counter1()}, {counter1()}, {counter1()}")
print(f"counter2: {counter2()}, {counter2()}")  # 独立的计数器

# =============================================================================
# 9. 递归函数
# =============================================================================

print("\n=== 递归函数 ===")


def factorial(n):
    """计算阶乘（递归实现）"""
    if n <= 1:
        return 1
    return n * factorial(n - 1)


print(f"factorial(5) = {factorial(5)}")


# 尾递归优化（Python 不支持，但可以改写）
def factorial_iter(n, acc=1):
    """尾递归形式（Python 中仍会栈溢出）"""
    if n <= 1:
        return acc
    return factorial_iter(n - 1, n * acc)


print(f"factorial_iter(5) = {factorial_iter(5)}")

# 获取递归限制
import sys

print(f"递归限制: {sys.getrecursionlimit()}")

# =============================================================================
# 10. 函数注解（Function Annotations）
# =============================================================================

print("\n=== 函数注解 ===")


def typed_add(x: int, y: int) -> int:
    """带类型注解的加法函数"""
    return x + y


print(f"typed_add(3, 5) = {typed_add(3, 5)}")
print(f"函数注解: {typed_add.__annotations__}")

# 注解不会强制类型检查，仍可传入其他类型
print(f"typed_add('a', 'b') = {typed_add('a', 'b')}")  # 字符串拼接

# =============================================================================
# 11. 文档字符串（Docstring）
# =============================================================================

print("\n=== 文档字符串 ===")


def well_documented_function(param1: str, param2: int = 0) -> dict:
    """
    这是一个文档完善的函数示例。

    这里是更详细的描述，解释函数的用途和行为。

    Args:
        param1: 第一个参数，字符串类型
        param2: 第二个参数，整数类型，默认为 0

    Returns:
        包含参数信息的字典

    Raises:
        ValueError: 当 param1 为空字符串时抛出

    Examples:
        >>> well_documented_function("test", 42)
        {'param1': 'test', 'param2': 42}
    """
    if not param1:
        raise ValueError("param1 cannot be empty")
    return {"param1": param1, "param2": param2}


# 访问文档
print(well_documented_function.__doc__)

# help() 函数显示文档
# help(well_documented_function)


if __name__ == "__main__":
    print("\n" + "=" * 50)
    print("03_functions.py 运行完成！")
    print("=" * 50)
