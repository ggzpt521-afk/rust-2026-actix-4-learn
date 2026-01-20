#!/usr/bin/env python3
"""
15_type_hints.py - 类型注解

Python 类型注解（Type Hints）：基本类型、泛型、TypeVar、Protocol 等
"""

from typing import (
    List, Dict, Set, Tuple, Optional, Union, Any,
    Callable, TypeVar, Generic, Protocol,
    Final, Literal, TypedDict, NewType,
    Sequence, Mapping, Iterable, Iterator,
    ClassVar, overload
)
from dataclasses import dataclass
from abc import abstractmethod

# =============================================================================
# 1. 基本类型注解
# =============================================================================

print("=== 基本类型注解 ===")

# 变量注解
name: str = "Alice"
age: int = 25
height: float = 1.68
is_student: bool = True
scores: list = [85, 92, 78]

print(f"name: {name} (type hint: str)")
print(f"age: {age} (type hint: int)")


# 函数注解
def greet(name: str, age: int) -> str:
    """带类型注解的函数"""
    return f"Hello, {name}! You are {age} years old."


result = greet("Bob", 30)
print(f"greet 结果: {result}")

# 查看注解
print(f"函数注解: {greet.__annotations__}")

# =============================================================================
# 2. 容器类型
# =============================================================================

print("\n=== 容器类型 ===")

# Python 3.9+ 可以直接使用内置类型
numbers: list[int] = [1, 2, 3]
pairs: dict[str, int] = {"a": 1, "b": 2}
unique: set[str] = {"x", "y", "z"}
point: tuple[int, int] = (3, 4)
mixed: tuple[str, int, float] = ("test", 1, 3.14)

# 可变长度元组
values: tuple[int, ...] = (1, 2, 3, 4, 5)

# 使用 typing 模块（兼容旧版本）
from typing import List, Dict, Set, Tuple

old_style: List[int] = [1, 2, 3]
old_dict: Dict[str, int] = {"a": 1}

print(f"numbers: {numbers}")
print(f"pairs: {pairs}")

# =============================================================================
# 3. Optional 和 Union
# =============================================================================

print("\n=== Optional 和 Union ===")


# Optional[X] 等价于 Union[X, None]
def find_user(user_id: int) -> Optional[str]:
    """可能返回 None 的函数"""
    users = {1: "Alice", 2: "Bob"}
    return users.get(user_id)


# Union 表示多种可能的类型
def process_input(value: Union[int, str]) -> str:
    """接受 int 或 str 的函数"""
    if isinstance(value, int):
        return f"整数: {value}"
    return f"字符串: {value}"


# Python 3.10+ 可以使用 | 语法
def new_style(value: int | str | None) -> str:
    """Python 3.10+ 的联合类型语法"""
    return str(value)


print(f"find_user(1): {find_user(1)}")
print(f"find_user(999): {find_user(999)}")
print(f"process_input(42): {process_input(42)}")
print(f"process_input('hello'): {process_input('hello')}")

# =============================================================================
# 4. Callable 类型
# =============================================================================

print("\n=== Callable 类型 ===")


# Callable[[参数类型...], 返回类型]
def apply_operation(
    func: Callable[[int, int], int],
    a: int,
    b: int
) -> int:
    """接受函数作为参数"""
    return func(a, b)


def add(x: int, y: int) -> int:
    return x + y


def multiply(x: int, y: int) -> int:
    return x * y


print(f"apply_operation(add, 3, 5): {apply_operation(add, 3, 5)}")
print(f"apply_operation(multiply, 3, 5): {apply_operation(multiply, 3, 5)}")


# 任意参数的 Callable
def log_call(func: Callable[..., Any]) -> Callable[..., Any]:
    """接受任意函数"""
    def wrapper(*args, **kwargs):
        print(f"  调用 {func.__name__}")
        return func(*args, **kwargs)
    return wrapper

# =============================================================================
# 5. TypeVar 和泛型
# =============================================================================

print("\n=== TypeVar 和泛型 ===")

# 定义类型变量
T = TypeVar('T')
K = TypeVar('K')
V = TypeVar('V')


def first(items: list[T]) -> T:
    """返回列表的第一个元素"""
    return items[0]


def get_or_default(d: dict[K, V], key: K, default: V) -> V:
    """获取字典值或返回默认值"""
    return d.get(key, default)


# 有约束的 TypeVar
Number = TypeVar('Number', int, float)


def double(x: Number) -> Number:
    """只接受 int 或 float"""
    return x * 2


print(f"first([1, 2, 3]): {first([1, 2, 3])}")
print(f"first(['a', 'b']): {first(['a', 'b'])}")
print(f"double(5): {double(5)}")
print(f"double(2.5): {double(2.5)}")


# 泛型类
class Stack(Generic[T]):
    """泛型栈"""

    def __init__(self) -> None:
        self._items: list[T] = []

    def push(self, item: T) -> None:
        self._items.append(item)

    def pop(self) -> T:
        return self._items.pop()

    def __repr__(self) -> str:
        return f"Stack({self._items})"


int_stack: Stack[int] = Stack()
int_stack.push(1)
int_stack.push(2)
print(f"int_stack: {int_stack}")

str_stack: Stack[str] = Stack()
str_stack.push("hello")
str_stack.push("world")
print(f"str_stack: {str_stack}")

# =============================================================================
# 6. Protocol（结构化子类型）
# =============================================================================

print("\n=== Protocol ===")


class Drawable(Protocol):
    """定义接口"""

    def draw(self) -> str:
        ...


class Circle:
    """实现 Drawable 协议（无需显式继承）"""

    def __init__(self, radius: float):
        self.radius = radius

    def draw(self) -> str:
        return f"Circle(r={self.radius})"


class Square:
    """实现 Drawable 协议"""

    def __init__(self, side: float):
        self.side = side

    def draw(self) -> str:
        return f"Square(s={self.side})"


def render(shape: Drawable) -> None:
    """接受任何有 draw 方法的对象"""
    print(f"  渲染: {shape.draw()}")


print("Protocol 示例:")
render(Circle(5))
render(Square(3))


# 运行时可检查的 Protocol
from typing import runtime_checkable


@runtime_checkable
class Sized(Protocol):
    def __len__(self) -> int:
        ...


print(f"\nisinstance([1,2,3], Sized): {isinstance([1, 2, 3], Sized)}")
print(f"isinstance('hello', Sized): {isinstance('hello', Sized)}")

# =============================================================================
# 7. Literal 和 Final
# =============================================================================

print("\n=== Literal 和 Final ===")


# Literal - 字面量类型
def set_mode(mode: Literal["read", "write", "append"]) -> str:
    return f"模式设置为: {mode}"


print(f"set_mode('read'): {set_mode('read')}")


# Final - 不可重新赋值的常量
MAX_SIZE: Final[int] = 100
API_URL: Final = "https://api.example.com"

print(f"MAX_SIZE: {MAX_SIZE}")
print(f"API_URL: {API_URL}")


# 类中的 Final
class Config:
    DEBUG: Final[bool] = False  # 类常量

# =============================================================================
# 8. TypedDict
# =============================================================================

print("\n=== TypedDict ===")


class UserDict(TypedDict):
    """类型化的字典"""
    name: str
    age: int
    email: str


class PartialUser(TypedDict, total=False):
    """所有键都是可选的"""
    name: str
    age: int


user: UserDict = {
    "name": "Alice",
    "age": 25,
    "email": "alice@example.com"
}

partial: PartialUser = {"name": "Bob"}  # age 是可选的

print(f"user: {user}")
print(f"partial: {partial}")

# =============================================================================
# 9. NewType
# =============================================================================

print("\n=== NewType ===")

# 创建新类型（主要用于类型检查，运行时无区别）
UserId = NewType('UserId', int)
ProductId = NewType('ProductId', int)


def get_user(user_id: UserId) -> str:
    return f"User #{user_id}"


def get_product(product_id: ProductId) -> str:
    return f"Product #{product_id}"


# 创建实例
uid = UserId(123)
pid = ProductId(456)

print(f"get_user(uid): {get_user(uid)}")
print(f"get_product(pid): {get_product(pid)}")

# 运行时 uid 就是普通 int
print(f"type(uid): {type(uid)}")

# =============================================================================
# 10. 函数重载
# =============================================================================

print("\n=== 函数重载 ===")


# 使用 @overload 定义多个签名
@overload
def process(x: int) -> int: ...
@overload
def process(x: str) -> str: ...
@overload
def process(x: list[int]) -> list[int]: ...


def process(x: int | str | list[int]) -> int | str | list[int]:
    """实际实现"""
    if isinstance(x, int):
        return x * 2
    elif isinstance(x, str):
        return x.upper()
    else:
        return [i * 2 for i in x]


print(f"process(5): {process(5)}")
print(f"process('hello'): {process('hello')}")
print(f"process([1, 2, 3]): {process([1, 2, 3])}")

# =============================================================================
# 11. 类型别名
# =============================================================================

print("\n=== 类型别名 ===")

# 简单别名
Vector = list[float]
Matrix = list[list[float]]

# 复杂别名
Callback = Callable[[str, int], bool]
UserCache = dict[int, tuple[str, int, bool]]


def scale_vector(v: Vector, factor: float) -> Vector:
    return [x * factor for x in v]


vec: Vector = [1.0, 2.0, 3.0]
print(f"scale_vector({vec}, 2): {scale_vector(vec, 2)}")

# Python 3.10+ TypeAlias
from typing import TypeAlias

Point2D: TypeAlias = tuple[float, float]
Point3D: TypeAlias = tuple[float, float, float]

# =============================================================================
# 12. 类型注解最佳实践
# =============================================================================

print("\n=== 类型注解最佳实践 ===")


# 1. 使用抽象类型而非具体类型
def good_practice(items: Sequence[int]) -> int:
    """使用 Sequence 而不是 list"""
    return sum(items)


# 2. 使用 dataclass 进行数据类型定义
@dataclass
class Point:
    x: float
    y: float

    def distance_from_origin(self) -> float:
        return (self.x ** 2 + self.y ** 2) ** 0.5


p = Point(3, 4)
print(f"Point: {p}, distance: {p.distance_from_origin()}")


# 3. 类变量 vs 实例变量
class MyClass:
    class_var: ClassVar[int] = 0  # 类变量
    instance_var: int  # 实例变量

    def __init__(self, value: int) -> None:
        self.instance_var = value


# 4. Self 类型（Python 3.11+）
from typing import Self


class Builder:
    def __init__(self) -> None:
        self.value = 0

    def add(self, x: int) -> Self:
        self.value += x
        return self

    def multiply(self, x: int) -> Self:
        self.value *= x
        return self


result = Builder().add(5).multiply(3).value
print(f"Builder 链式调用: {result}")


# 5. 运行时类型检查
def strict_add(a: int, b: int) -> int:
    if not isinstance(a, int) or not isinstance(b, int):
        raise TypeError("Both arguments must be integers")
    return a + b


print(f"strict_add(3, 5): {strict_add(3, 5)}")


if __name__ == "__main__":
    print("\n" + "=" * 50)
    print("15_type_hints.py 运行完成！")
    print("=" * 50)
