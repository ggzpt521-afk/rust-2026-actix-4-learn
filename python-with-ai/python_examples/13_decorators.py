#!/usr/bin/env python3
"""
13_decorators.py - 装饰器

Python 装饰器：函数装饰器、类装饰器、带参数装饰器、functools 等
"""

import functools
import time
from typing import Callable, Any

# =============================================================================
# 1. 装饰器基础
# =============================================================================

print("=== 装饰器基础 ===")


def simple_decorator(func):
    """最简单的装饰器"""
    def wrapper():
        print("  函数调用前")
        result = func()
        print("  函数调用后")
        return result
    return wrapper


@simple_decorator
def say_hello():
    print("  Hello!")


# @decorator 语法等价于: say_hello = simple_decorator(say_hello)
say_hello()

# =============================================================================
# 2. 保留原函数信息
# =============================================================================

print("\n=== 保留原函数信息 ===")


def decorator_without_wraps(func):
    def wrapper(*args, **kwargs):
        return func(*args, **kwargs)
    return wrapper


def decorator_with_wraps(func):
    @functools.wraps(func)  # 保留原函数的元信息
    def wrapper(*args, **kwargs):
        return func(*args, **kwargs)
    return wrapper


@decorator_without_wraps
def func_a():
    """函数 A 的文档"""
    pass


@decorator_with_wraps
def func_b():
    """函数 B 的文档"""
    pass


print(f"无 @wraps: func_a.__name__ = {func_a.__name__}, __doc__ = {func_a.__doc__}")
print(f"有 @wraps: func_b.__name__ = {func_b.__name__}, __doc__ = {func_b.__doc__}")

# =============================================================================
# 3. 带参数的函数装饰器
# =============================================================================

print("\n=== 带参数的函数装饰器 ===")


def my_decorator(func):
    """接受任意参数的装饰器"""
    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        print(f"  调用 {func.__name__}，参数: args={args}, kwargs={kwargs}")
        result = func(*args, **kwargs)
        print(f"  返回值: {result}")
        return result
    return wrapper


@my_decorator
def add(a, b):
    """加法函数"""
    return a + b


@my_decorator
def greet(name, greeting="Hello"):
    """问候函数"""
    return f"{greeting}, {name}!"


add(3, 5)
greet("Alice", greeting="Hi")

# =============================================================================
# 4. 带参数的装饰器
# =============================================================================

print("\n=== 带参数的装饰器 ===")


def repeat(times):
    """重复执行装饰器"""
    def decorator(func):
        @functools.wraps(func)
        def wrapper(*args, **kwargs):
            results = []
            for _ in range(times):
                result = func(*args, **kwargs)
                results.append(result)
            return results
        return wrapper
    return decorator


@repeat(times=3)
def say_hi():
    print("  Hi!")
    return "done"


results = say_hi()
print(f"结果: {results}")


# 带可选参数的装饰器
def log(func=None, *, level="INFO"):
    """可带参数或不带参数使用的装饰器"""
    def decorator(fn):
        @functools.wraps(fn)
        def wrapper(*args, **kwargs):
            print(f"  [{level}] 调用 {fn.__name__}")
            return fn(*args, **kwargs)
        return wrapper

    if func is not None:
        # 不带参数调用: @log
        return decorator(func)
    # 带参数调用: @log(level="DEBUG")
    return decorator


@log
def func1():
    pass


@log(level="DEBUG")
def func2():
    pass


func1()
func2()

# =============================================================================
# 5. 实用装饰器示例
# =============================================================================

print("\n=== 实用装饰器示例 ===")


# 计时装饰器
def timer(func):
    """测量函数执行时间"""
    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        start = time.perf_counter()
        result = func(*args, **kwargs)
        end = time.perf_counter()
        print(f"  {func.__name__} 耗时: {end - start:.6f}秒")
        return result
    return wrapper


@timer
def slow_function():
    time.sleep(0.1)
    return "done"


slow_function()


# 缓存装饰器
def memoize(func):
    """简单的缓存装饰器"""
    cache = {}

    @functools.wraps(func)
    def wrapper(*args):
        if args not in cache:
            cache[args] = func(*args)
        return cache[args]
    return wrapper


@memoize
def fibonacci(n):
    """斐波那契数列（带缓存）"""
    if n < 2:
        return n
    return fibonacci(n - 1) + fibonacci(n - 2)


print(f"\nfibonacci(30) = {fibonacci(30)}")


# 使用内置的 lru_cache（推荐）
@functools.lru_cache(maxsize=128)
def fibonacci_builtin(n):
    """使用内置缓存的斐波那契"""
    if n < 2:
        return n
    return fibonacci_builtin(n - 1) + fibonacci_builtin(n - 2)


print(f"fibonacci_builtin(30) = {fibonacci_builtin(30)}")
print(f"缓存信息: {fibonacci_builtin.cache_info()}")


# 重试装饰器
def retry(max_attempts=3, delay=1):
    """重试装饰器"""
    def decorator(func):
        @functools.wraps(func)
        def wrapper(*args, **kwargs):
            last_exception = None
            for attempt in range(1, max_attempts + 1):
                try:
                    return func(*args, **kwargs)
                except Exception as e:
                    last_exception = e
                    print(f"  尝试 {attempt}/{max_attempts} 失败: {e}")
                    if attempt < max_attempts:
                        time.sleep(delay)
            raise last_exception
        return wrapper
    return decorator


attempt_count = 0


@retry(max_attempts=3, delay=0.1)
def unreliable_function():
    """模拟不可靠的函数"""
    global attempt_count
    attempt_count += 1
    if attempt_count < 3:
        raise ValueError("随机失败")
    return "成功！"


print("\n重试装饰器:")
try:
    result = unreliable_function()
    print(f"  最终结果: {result}")
except ValueError as e:
    print(f"  所有尝试都失败: {e}")


# 类型检查装饰器
def type_check(func):
    """基于类型注解进行类型检查"""
    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        hints = func.__annotations__
        # 检查位置参数
        for i, (arg, (name, expected_type)) in enumerate(
            zip(args, list(hints.items())[:-1] if 'return' in hints else list(hints.items()))
        ):
            if not isinstance(arg, expected_type):
                raise TypeError(f"参数 {name} 期望 {expected_type.__name__}，得到 {type(arg).__name__}")
        return func(*args, **kwargs)
    return wrapper


@type_check
def typed_add(a: int, b: int) -> int:
    return a + b


print(f"\ntyped_add(3, 5) = {typed_add(3, 5)}")
try:
    typed_add("3", 5)
except TypeError as e:
    print(f"类型检查失败: {e}")

# =============================================================================
# 6. 类装饰器
# =============================================================================

print("\n=== 类装饰器 ===")


# 用类实现装饰器
class CountCalls:
    """统计函数调用次数的装饰器类"""

    def __init__(self, func):
        functools.update_wrapper(self, func)
        self.func = func
        self.calls = 0

    def __call__(self, *args, **kwargs):
        self.calls += 1
        print(f"  {self.func.__name__} 被调用了 {self.calls} 次")
        return self.func(*args, **kwargs)


@CountCalls
def say_something():
    print("  Something!")


say_something()
say_something()
say_something()


# 装饰类的装饰器
def singleton(cls):
    """单例装饰器"""
    instances = {}

    @functools.wraps(cls)
    def get_instance(*args, **kwargs):
        if cls not in instances:
            instances[cls] = cls(*args, **kwargs)
        return instances[cls]

    return get_instance


@singleton
class Database:
    def __init__(self):
        print("  数据库初始化")
        self.connected = True


print("\n单例装饰器:")
db1 = Database()
db2 = Database()
print(f"db1 is db2: {db1 is db2}")


# 添加方法的类装饰器
def add_repr(cls):
    """自动添加 __repr__ 方法"""
    def __repr__(self):
        attrs = ', '.join(f'{k}={v!r}' for k, v in self.__dict__.items())
        return f'{cls.__name__}({attrs})'
    cls.__repr__ = __repr__
    return cls


@add_repr
class Point:
    def __init__(self, x, y):
        self.x = x
        self.y = y


print(f"\n@add_repr: {Point(3, 4)}")

# =============================================================================
# 7. 装饰器堆叠
# =============================================================================

print("\n=== 装饰器堆叠 ===")


def decorator_a(func):
    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        print("  A 前")
        result = func(*args, **kwargs)
        print("  A 后")
        return result
    return wrapper


def decorator_b(func):
    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        print("  B 前")
        result = func(*args, **kwargs)
        print("  B 后")
        return result
    return wrapper


@decorator_a
@decorator_b
def stacked_function():
    print("  函数执行")


# 等价于: decorator_a(decorator_b(stacked_function))
stacked_function()

# =============================================================================
# 8. 内置装饰器
# =============================================================================

print("\n=== 内置装饰器 ===")


class MyClass:
    _instances = 0

    def __init__(self, value):
        self._value = value
        MyClass._instances += 1

    @property
    def value(self):
        """属性 getter"""
        return self._value

    @value.setter
    def value(self, new_value):
        """属性 setter"""
        self._value = new_value

    @classmethod
    def get_instance_count(cls):
        """类方法"""
        return cls._instances

    @staticmethod
    def helper_function(x):
        """静态方法"""
        return x * 2


obj = MyClass(10)
print(f"property: obj.value = {obj.value}")
obj.value = 20
print(f"setter: obj.value = {obj.value}")
print(f"classmethod: MyClass.get_instance_count() = {MyClass.get_instance_count()}")
print(f"staticmethod: MyClass.helper_function(5) = {MyClass.helper_function(5)}")

# =============================================================================
# 9. functools 模块装饰器
# =============================================================================

print("\n=== functools 模块装饰器 ===")

# @functools.lru_cache - 带 LRU 缓存
# 已在前面演示

# @functools.cache - 无限缓存（Python 3.9+）
@functools.cache
def factorial(n):
    return n * factorial(n - 1) if n else 1


print(f"factorial(10) = {factorial(10)}")


# @functools.total_ordering - 自动生成比较方法
@functools.total_ordering
class Student:
    def __init__(self, name, score):
        self.name = name
        self.score = score

    def __eq__(self, other):
        return self.score == other.score

    def __lt__(self, other):
        return self.score < other.score


s1 = Student("Alice", 85)
s2 = Student("Bob", 90)
print(f"\n@total_ordering:")
print(f"  s1 < s2: {s1 < s2}")
print(f"  s1 <= s2: {s1 <= s2}")
print(f"  s1 > s2: {s1 > s2}")
print(f"  s1 >= s2: {s1 >= s2}")


# @functools.singledispatch - 单分派泛型函数
@functools.singledispatch
def process(value):
    print(f"  默认处理: {value}")


@process.register(int)
def _(value):
    print(f"  处理整数: {value}")


@process.register(list)
def _(value):
    print(f"  处理列表: {value}, 长度: {len(value)}")


@process.register(str)
def _(value):
    print(f"  处理字符串: '{value}'")


print("\n@singledispatch:")
process(42)
process([1, 2, 3])
process("hello")
process(3.14)  # 使用默认处理


if __name__ == "__main__":
    print("\n" + "=" * 50)
    print("13_decorators.py 运行完成！")
    print("=" * 50)
