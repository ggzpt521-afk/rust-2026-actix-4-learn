#!/usr/bin/env python3
"""
01_variables.py - 变量与常量

Python 是动态类型语言，变量无需声明类型。
本文件演示变量定义、命名规范、常量约定等基础概念。
"""

# =============================================================================
# 1. 变量定义与赋值
# =============================================================================

# 简单赋值
name = "Python"
age = 30
price = 99.99
is_active = True

print("=== 基本变量 ===")
print(f"name = {name}, type: {type(name).__name__}")
print(f"age = {age}, type: {type(age).__name__}")
print(f"price = {price}, type: {type(price).__name__}")
print(f"is_active = {is_active}, type: {type(is_active).__name__}")

# =============================================================================
# 2. 多重赋值
# =============================================================================

print("\n=== 多重赋值 ===")

# 同时给多个变量赋相同的值
x = y = z = 0
print(f"x = y = z = {x}")

# 同时给多个变量赋不同的值（解包赋值）
a, b, c = 1, 2, 3
print(f"a, b, c = {a}, {b}, {c}")

# 交换变量值（Pythonic 方式）
a, b = b, a
print(f"交换后: a = {a}, b = {b}")

# =============================================================================
# 3. 变量命名规范（PEP 8）
# =============================================================================

print("\n=== 命名规范示例 ===")

# 普通变量：小写字母，下划线分隔（snake_case）
user_name = "Alice"
total_count = 100
max_retry_times = 3

# 私有变量（约定，非强制）：单下划线开头
_internal_value = "内部使用"

# 名称修饰（Name Mangling）：双下划线开头
# 在类中会被重命名为 _ClassName__var
__private_var = "强私有"

# 特殊变量：双下划线开头和结尾（魔术方法/属性）
# 如 __name__, __doc__ 等，不建议自定义

print(f"user_name = {user_name}")
print(f"_internal_value = {_internal_value}")
print(f"__private_var = {__private_var}")

# =============================================================================
# 4. 常量（约定俗成，Python 无真正的常量）
# =============================================================================

print("\n=== 常量约定 ===")

# 常量命名：全大写，下划线分隔
MAX_CONNECTIONS = 100
DEFAULT_TIMEOUT = 30
PI = 3.14159
DATABASE_URL = "postgresql://localhost/mydb"

print(f"MAX_CONNECTIONS = {MAX_CONNECTIONS}")
print(f"PI = {PI}")

# Python 3.8+ 可使用 typing.Final 进行类型提示（但不会真正阻止修改）
from typing import Final

API_VERSION: Final[str] = "v1.0"
print(f"API_VERSION = {API_VERSION}")

# =============================================================================
# 5. 变量作用域
# =============================================================================

print("\n=== 变量作用域 ===")

# 全局变量
global_var = "我是全局变量"


def scope_demo():
    """演示变量作用域"""
    global global_var  # 声明需要在使用之前
    # 局部变量
    local_var = "我是局部变量"
    print(f"  函数内访问 global_var: {global_var}")
    print(f"  函数内访问 local_var: {local_var}")

    # 修改全局变量
    global_var = "全局变量被修改了"


scope_demo()
print(f"函数外 global_var: {global_var}")

# nonlocal 用于嵌套函数
def outer():
    outer_var = "外层变量"

    def inner():
        nonlocal outer_var
        outer_var = "被内层函数修改"
        print(f"  inner: {outer_var}")

    inner()
    print(f"  outer: {outer_var}")


print("\nnonlocal 示例:")
outer()

# =============================================================================
# 6. 变量的内存管理
# =============================================================================

print("\n=== 内存与引用 ===")

# Python 变量是对象的引用
list1 = [1, 2, 3]
list2 = list1  # list2 指向同一个对象
list2.append(4)

print(f"list1 = {list1}")  # [1, 2, 3, 4]
print(f"list2 = {list2}")  # [1, 2, 3, 4]
print(f"list1 is list2: {list1 is list2}")  # True

# 创建副本
list3 = list1.copy()  # 或 list1[:]
list3.append(5)

print(f"list1 = {list1}")  # [1, 2, 3, 4]
print(f"list3 = {list3}")  # [1, 2, 3, 4, 5]
print(f"list1 is list3: {list1 is list3}")  # False

# id() 查看对象内存地址
print(f"\nid(list1) = {id(list1)}")
print(f"id(list2) = {id(list2)}")
print(f"id(list3) = {id(list3)}")

# =============================================================================
# 7. 删除变量
# =============================================================================

print("\n=== 删除变量 ===")

temp_var = "临时数据"
print(f"删除前: temp_var = {temp_var}")

del temp_var  # 删除变量

try:
    print(temp_var)
except NameError as e:
    print(f"删除后访问变量: {e}")

# =============================================================================
# 8. 特殊值
# =============================================================================

print("\n=== 特殊值 ===")

# None - 表示空值或无值
result = None
print(f"result = {result}, type: {type(result).__name__}")
print(f"result is None: {result is None}")  # 推荐用 is 比较 None

# 区分 None、0、空字符串、空列表
values = [None, 0, "", [], False]
for v in values:
    print(f"  {repr(v):10} -> bool: {bool(v)}, is None: {v is None}")


if __name__ == "__main__":
    print("\n" + "=" * 50)
    print("01_variables.py 运行完成！")
    print("=" * 50)
