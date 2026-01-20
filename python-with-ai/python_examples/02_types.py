#!/usr/bin/env python3
"""
02_types.py - 基本数据类型

Python 内置的基本数据类型包括：
- 数值类型：int, float, complex
- 布尔类型：bool
- 字符串类型：str
- 字节类型：bytes, bytearray
"""

# =============================================================================
# 1. 整数（int）
# =============================================================================

print("=== 整数类型 (int) ===")

# Python 3 的整数没有大小限制
small_int = 42
big_int = 123456789012345678901234567890

print(f"小整数: {small_int}")
print(f"大整数: {big_int}")

# 不同进制表示
decimal = 255       # 十进制
binary = 0b11111111  # 二进制（0b 前缀）
octal = 0o377       # 八进制（0o 前缀）
hexadecimal = 0xFF  # 十六进制（0x 前缀）

print(f"\n不同进制表示 255:")
print(f"  十进制: {decimal}")
print(f"  二进制: {bin(binary)} = {binary}")
print(f"  八进制: {oct(octal)} = {octal}")
print(f"  十六进制: {hex(hexadecimal)} = {hexadecimal}")

# 数字分隔符（提高可读性，Python 3.6+）
million = 1_000_000
binary_readable = 0b1111_0000_1111_0000

print(f"\n使用下划线分隔: {million:,}")

# =============================================================================
# 2. 浮点数（float）
# =============================================================================

print("\n=== 浮点数类型 (float) ===")

pi = 3.14159
negative = -0.001
scientific = 1.5e-10  # 科学计数法

print(f"pi = {pi}")
print(f"negative = {negative}")
print(f"scientific = {scientific}")

# 浮点数精度问题
result = 0.1 + 0.2
print(f"\n0.1 + 0.2 = {result}")  # 不等于 0.3！
print(f"0.1 + 0.2 == 0.3: {result == 0.3}")

# 使用 decimal 模块获得精确计算
from decimal import Decimal, getcontext

getcontext().prec = 6  # 设置精度
d1 = Decimal("0.1")
d2 = Decimal("0.2")
print(f"Decimal: 0.1 + 0.2 = {d1 + d2}")

# 特殊浮点值
import math

print(f"\n特殊浮点值:")
print(f"  正无穷: {float('inf')}")
print(f"  负无穷: {float('-inf')}")
print(f"  NaN: {float('nan')}")
print(f"  math.isnan(float('nan')): {math.isnan(float('nan'))}")
print(f"  math.isinf(float('inf')): {math.isinf(float('inf'))}")

# =============================================================================
# 3. 复数（complex）
# =============================================================================

print("\n=== 复数类型 (complex) ===")

c1 = 3 + 4j
c2 = complex(1, 2)

print(f"c1 = {c1}")
print(f"c2 = {c2}")
print(f"c1 的实部: {c1.real}")
print(f"c1 的虚部: {c1.imag}")
print(f"c1 的共轭: {c1.conjugate()}")
print(f"c1 + c2 = {c1 + c2}")
print(f"|c1| (模) = {abs(c1)}")

# =============================================================================
# 4. 布尔类型（bool）
# =============================================================================

print("\n=== 布尔类型 (bool) ===")

is_valid = True
is_empty = False

print(f"is_valid = {is_valid}, type: {type(is_valid).__name__}")
print(f"is_empty = {is_empty}")

# bool 是 int 的子类
print(f"\nbool 是 int 的子类: {issubclass(bool, int)}")
print(f"True + True = {True + True}")   # 2
print(f"True * 10 = {True * 10}")       # 10

# 真值测试（Truthy / Falsy）
print("\n假值（Falsy）示例:")
falsy_values = [False, None, 0, 0.0, 0j, "", [], {}, set(), frozenset()]
for v in falsy_values:
    print(f"  bool({repr(v):15}) = {bool(v)}")

print("\n真值（Truthy）示例:")
truthy_values = [True, 1, -1, 0.1, "hello", [1], {"a": 1}]
for v in truthy_values:
    print(f"  bool({repr(v):15}) = {bool(v)}")

# =============================================================================
# 5. 字符串（str）- 基础部分
# =============================================================================

print("\n=== 字符串类型 (str) ===")

# 多种定义方式
s1 = 'single quotes'
s2 = "double quotes"
s3 = """多行字符串
可以跨越
多行"""
s4 = '''也可以用
单引号'''

print(f"s1: {s1}")
print(f"s3:\n{s3}")

# 原始字符串（Raw String）
path = r"C:\Users\name\documents"  # 反斜杠不转义
print(f"\n原始字符串: {path}")

# Unicode 字符串
chinese = "你好，世界！"
emoji = "Python 🐍"
print(f"中文: {chinese}")
print(f"Emoji: {emoji}")

# 字符串编码
encoded = chinese.encode("utf-8")
decoded = encoded.decode("utf-8")
print(f"\n编码: {encoded}")
print(f"解码: {decoded}")

# =============================================================================
# 6. 字节类型（bytes 和 bytearray）
# =============================================================================

print("\n=== 字节类型 (bytes/bytearray) ===")

# bytes - 不可变字节序列
b1 = b"hello"
b2 = bytes([72, 101, 108, 108, 111])  # ASCII 码

print(f"b1 = {b1}")
print(f"b2 = {b2}")
print(f"b1 == b2: {b1 == b2}")

# bytearray - 可变字节序列
ba = bytearray(b"hello")
ba[0] = 72  # 修改第一个字节
ba.append(33)  # 添加 '!'

print(f"bytearray: {ba}")
print(f"转换为 bytes: {bytes(ba)}")

# =============================================================================
# 7. 类型转换
# =============================================================================

print("\n=== 类型转换 ===")

# 显式类型转换
num_str = "123"
num_int = int(num_str)
num_float = float(num_str)

print(f"str -> int: '{num_str}' -> {num_int}")
print(f"str -> float: '{num_str}' -> {num_float}")
print(f"int -> str: {num_int} -> '{str(num_int)}'")
print(f"float -> int: {3.7} -> {int(3.7)}")  # 截断，不是四舍五入

# 四舍五入
print(f"round(3.7) = {round(3.7)}")
print(f"round(3.5) = {round(3.5)}")  # 银行家舍入法

# 进制转换
print(f"\n进制转换:")
print(f"  int('ff', 16) = {int('ff', 16)}")
print(f"  int('1010', 2) = {int('1010', 2)}")

# =============================================================================
# 8. 类型检查
# =============================================================================

print("\n=== 类型检查 ===")

value = 42

# type() 精确匹配
print(f"type({value}) == int: {type(value) == int}")

# isinstance() 支持继承检查（推荐）
print(f"isinstance({value}, int): {isinstance(value, int)}")
print(f"isinstance({value}, (int, float)): {isinstance(value, (int, float))}")

# 检查是否为数值类型
from numbers import Number

print(f"isinstance(3.14, Number): {isinstance(3.14, Number)}")
print(f"isinstance(3+4j, Number): {isinstance(3+4j, Number)}")

# =============================================================================
# 9. 数值运算
# =============================================================================

print("\n=== 数值运算 ===")

a, b = 17, 5

print(f"a = {a}, b = {b}")
print(f"  a + b = {a + b}")      # 加法
print(f"  a - b = {a - b}")      # 减法
print(f"  a * b = {a * b}")      # 乘法
print(f"  a / b = {a / b}")      # 除法（总是返回 float）
print(f"  a // b = {a // b}")    # 整除（向下取整）
print(f"  a % b = {a % b}")      # 取模
print(f"  a ** b = {a ** b}")    # 幂运算
print(f"  -a // b = {-a // b}")  # 注意负数整除

# divmod 同时返回商和余数
quotient, remainder = divmod(a, b)
print(f"  divmod({a}, {b}) = ({quotient}, {remainder})")

# 位运算
print(f"\n位运算 (a=17=0b10001, b=5=0b00101):")
print(f"  a & b (AND) = {a & b} = {bin(a & b)}")
print(f"  a | b (OR)  = {a | b} = {bin(a | b)}")
print(f"  a ^ b (XOR) = {a ^ b} = {bin(a ^ b)}")
print(f"  ~a (NOT)    = {~a}")
print(f"  a << 2      = {a << 2}")  # 左移
print(f"  a >> 2      = {a >> 2}")  # 右移


if __name__ == "__main__":
    print("\n" + "=" * 50)
    print("02_types.py 运行完成！")
    print("=" * 50)
