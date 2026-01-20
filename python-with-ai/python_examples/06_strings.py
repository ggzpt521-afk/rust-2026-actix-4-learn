#!/usr/bin/env python3
"""
06_strings.py - 字符串与格式化

Python 字符串操作、格式化方法、正则表达式基础
"""

# =============================================================================
# 1. 字符串基础操作
# =============================================================================

print("=== 字符串基础操作 ===")

s = "Hello, Python!"

# 索引和切片
print(f"s = '{s}'")
print(f"s[0] = '{s[0]}'")
print(f"s[-1] = '{s[-1]}'")
print(f"s[0:5] = '{s[0:5]}'")
print(f"s[7:] = '{s[7:]}'")
print(f"s[::-1] = '{s[::-1]}'")  # 反转

# 长度
print(f"len(s) = {len(s)}")

# 成员检测
print(f"'Python' in s: {'Python' in s}")
print(f"'Java' not in s: {'Java' not in s}")

# 字符串不可变
try:
    s[0] = "h"
except TypeError as e:
    print(f"字符串不可变: {e}")

# =============================================================================
# 2. 字符串方法
# =============================================================================

print("\n=== 字符串方法 ===")

text = "  Hello, World!  "

# 大小写转换
print(f"原字符串: '{text.strip()}'")
print(f"upper(): '{text.strip().upper()}'")
print(f"lower(): '{text.strip().lower()}'")
print(f"title(): '{text.strip().title()}'")
print(f"capitalize(): '{text.strip().capitalize()}'")
print(f"swapcase(): '{text.strip().swapcase()}'")

# 去除空白
print(f"\n去除空白:")
print(f"strip(): '{text.strip()}'")
print(f"lstrip(): '{text.lstrip()}'")
print(f"rstrip(): '{text.rstrip()}'")

# 查找和替换
print(f"\n查找和替换:")
s = "hello world, hello python"
print(f"s = '{s}'")
print(f"find('hello'): {s.find('hello')}")
print(f"rfind('hello'): {s.rfind('hello')}")
print(f"index('world'): {s.index('world')}")
print(f"count('hello'): {s.count('hello')}")
print(f"replace('hello', 'hi'): '{s.replace('hello', 'hi')}'")
print(f"replace('hello', 'hi', 1): '{s.replace('hello', 'hi', 1)}'")

# 判断方法
print(f"\n判断方法:")
print(f"'hello'.isalpha(): {'hello'.isalpha()}")
print(f"'123'.isdigit(): {'123'.isdigit()}")
print(f"'hello123'.isalnum(): {'hello123'.isalnum()}")
print(f"'   '.isspace(): {'   '.isspace()}")
print(f"'Hello'.istitle(): {'Hello'.istitle()}")
print(f"'HELLO'.isupper(): {'HELLO'.isupper()}")
print(f"'hello'.islower(): {'hello'.islower()}")

# 前缀和后缀
print(f"\n前缀和后缀:")
filename = "document.pdf"
print(f"'{filename}'.startswith('doc'): {filename.startswith('doc')}")
print(f"'{filename}'.endswith('.pdf'): {filename.endswith('.pdf')}")
print(f"'{filename}'.endswith(('.pdf', '.doc')): {filename.endswith(('.pdf', '.doc'))}")

# 分割和连接
print(f"\n分割和连接:")
csv_line = "apple,banana,cherry"
print(f"split(','): {csv_line.split(',')}")

path = "/usr/local/bin"
print(f"split('/'): {path.split('/')}")

text = "Line 1\nLine 2\nLine 3"
print(f"splitlines(): {text.splitlines()}")

words = ["Hello", "World"]
print(f"' '.join(words): {' '.join(words)}")
print(f"'-'.join(words): {'-'.join(words)}")

# 对齐和填充
print(f"\n对齐和填充:")
s = "Python"
print(f"center(20, '-'): '{s.center(20, '-')}'")
print(f"ljust(20, '.'): '{s.ljust(20, '.')}'")
print(f"rjust(20, '.'): '{s.rjust(20, '.')}'")
print(f"zfill(10): '{'42'.zfill(10)}'")

# =============================================================================
# 3. 字符串格式化
# =============================================================================

print("\n=== 字符串格式化 ===")

name = "Alice"
age = 25
height = 1.68
score = 95.5

# f-string（Python 3.6+，推荐）
print("f-string 格式化:")
print(f"  基本: {name} is {age} years old")
print(f"  表达式: {name.upper()} is {age * 12} months old")
print(f"  宽度: |{name:10}|{age:5}|")
print(f"  对齐: |{name:<10}|{name:>10}|{name:^10}|")
print(f"  填充: |{name:*^10}|")
print(f"  精度: {height:.1f}m, {score:.0f}分")
print(f"  千位分隔: {1234567:,}")
print(f"  百分比: {0.856:.2%}")
print(f"  进制: {255:#x}, {255:#b}, {255:#o}")

# 调试模式（Python 3.8+）
x = 10
y = 20
print(f"  调试: {x=}, {y=}, {x+y=}")

# format() 方法
print("\nformat() 方法:")
template = "{name} scored {score:.1f} points"
print(f"  {template.format(name='Bob', score=87.5)}")

# 位置参数
print("  {0} vs {1}".format("Python", "Java"))
print("  {1} vs {0}".format("Python", "Java"))

# 格式规格说明
print("  {:>10.2f}".format(3.14159))
print("  {:0>5d}".format(42))

# % 格式化（旧式，不推荐但需了解）
print("\n% 格式化（旧式）:")
print("  %s is %d years old" % (name, age))
print("  Pi = %.4f" % 3.14159)
print("  %(name)s: %(score).1f" % {"name": "Charlie", "score": 88.5})

# =============================================================================
# 4. 字符串模板
# =============================================================================

print("\n=== 字符串模板 ===")

from string import Template

template = Template("$name is $age years old")
result = template.substitute(name="David", age=30)
print(f"Template: {result}")

# 安全替换（缺少键不报错）
template = Template("$name, $title")
result = template.safe_substitute(name="Eve")
print(f"safe_substitute: {result}")

# =============================================================================
# 5. 原始字符串和转义
# =============================================================================

print("\n=== 原始字符串和转义 ===")

# 常见转义序列
print("转义序列:")
print("  换行: Hello\\nWorld -> Hello\nWorld")
print("  制表: Hello\\tWorld -> Hello\tWorld")
print("  反斜杠: \\\\ -> \\")
print("  引号: \\' -> '")

# 原始字符串
path = r"C:\Users\name\documents"
print(f"\n原始字符串: {path}")

regex_pattern = r"\d+\.\d+"
print(f"正则表达式: {regex_pattern}")

# =============================================================================
# 6. Unicode 和编码
# =============================================================================

print("\n=== Unicode 和编码 ===")

# Unicode 字符
chinese = "中文"
emoji = "Python 🐍 is fun! 🎉"

print(f"中文: {chinese}")
print(f"Emoji: {emoji}")

# Unicode 转义
print(f"Unicode 转义: {'\\u4e2d\\u6587'} = {'\u4e2d\u6587'}")
print(f"Unicode 名称: {'\\N{{SNAKE}}'} = {chr(0x1F40D)}")

# ord() 和 chr()
print(f"\nord('A') = {ord('A')}")
print(f"chr(65) = '{chr(65)}'")
print(f"ord('中') = {ord('中')}")
print(f"chr(20013) = '{chr(20013)}'")

# 编码和解码
text = "Hello, 世界!"
encoded = text.encode("utf-8")
print(f"\n原文: {text}")
print(f"UTF-8 编码: {encoded}")
print(f"解码: {encoded.decode('utf-8')}")

# 不同编码
print(f"GBK 编码: {text.encode('gbk')}")

# =============================================================================
# 7. 正则表达式
# =============================================================================

print("\n=== 正则表达式 ===")

import re

text = "Contact: alice@example.com, bob@test.org, invalid-email"

# 基本匹配
pattern = r"\w+@\w+\.\w+"
matches = re.findall(pattern, text)
print(f"findall 邮箱: {matches}")

# match vs search
print("\nmatch vs search:")
print(f"match('Contact', text): {re.match('Contact', text)}")
print(f"match('alice', text): {re.match('alice', text)}")  # None
print(f"search('alice', text): {re.search('alice', text)}")

# 分组
pattern = r"(\w+)@(\w+)\.(\w+)"
match = re.search(pattern, text)
if match:
    print(f"\n分组匹配:")
    print(f"  完整匹配: {match.group(0)}")
    print(f"  用户名: {match.group(1)}")
    print(f"  域名: {match.group(2)}")
    print(f"  后缀: {match.group(3)}")
    print(f"  所有分组: {match.groups()}")

# 命名分组
pattern = r"(?P<user>\w+)@(?P<domain>\w+)\.(?P<suffix>\w+)"
match = re.search(pattern, text)
if match:
    print(f"\n命名分组:")
    print(f"  user: {match.group('user')}")
    print(f"  domain: {match.group('domain')}")
    print(f"  groupdict: {match.groupdict()}")

# 替换
result = re.sub(r"\w+@\w+\.\w+", "[EMAIL]", text)
print(f"\n替换: {result}")

# 分割
text = "apple, banana; cherry  orange"
result = re.split(r"[,;\s]+", text)
print(f"正则分割: {result}")

# 编译正则表达式（提高效率）
email_pattern = re.compile(r"\w+@\w+\.\w+", re.IGNORECASE)
print(f"\n编译后的模式: {email_pattern.findall(text)}")

# 常用正则模式示例
print("\n常用正则模式:")
patterns = {
    "数字": r"\d+",
    "单词": r"\b\w+\b",
    "中文": r"[\u4e00-\u9fa5]+",
    "手机号": r"1[3-9]\d{9}",
    "日期": r"\d{4}-\d{2}-\d{2}",
}

test_text = "2024年01月15日，张三的手机是13912345678，订单号123456"
for name, pattern in patterns.items():
    matches = re.findall(pattern, test_text)
    print(f"  {name}: {matches}")

# =============================================================================
# 8. 字符串性能
# =============================================================================

print("\n=== 字符串性能提示 ===")

# 字符串拼接
# 不好的做法（每次创建新字符串）
# result = ""
# for i in range(1000):
#     result += str(i)

# 好的做法
parts = [str(i) for i in range(10)]
result = "".join(parts)
print(f"使用 join 拼接: {result}")

# 使用 io.StringIO 进行大量字符串操作
from io import StringIO

buffer = StringIO()
for i in range(5):
    buffer.write(f"Line {i}\n")
print(f"StringIO 结果:\n{buffer.getvalue()}")


if __name__ == "__main__":
    print("\n" + "=" * 50)
    print("06_strings.py 运行完成！")
    print("=" * 50)
