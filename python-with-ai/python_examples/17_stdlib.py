#!/usr/bin/env python3
"""
17_stdlib.py - 常用标准库

Python 常用标准库：os、pathlib、json、collections、functools、itertools 等
"""

# =============================================================================
# 1. os 模块
# =============================================================================

print("=== os 模块 ===")

import os

# 环境变量
print(f"HOME: {os.environ.get('HOME')}")
print(f"PATH 条目数: {len(os.environ.get('PATH', '').split(os.pathsep))}")

# 系统信息
print(f"\n系统信息:")
print(f"  os.name: {os.name}")
print(f"  os.getcwd(): {os.getcwd()}")
print(f"  os.cpu_count(): {os.cpu_count()}")

# 路径操作
print(f"\nos.path 操作:")
path = "/usr/local/bin/python"
print(f"  basename: {os.path.basename(path)}")
print(f"  dirname: {os.path.dirname(path)}")
print(f"  split: {os.path.split(path)}")
print(f"  splitext: {os.path.splitext('file.txt')}")
print(f"  join: {os.path.join('a', 'b', 'c')}")
print(f"  exists: {os.path.exists('/tmp')}")
print(f"  isfile: {os.path.isfile(__file__)}")
print(f"  isdir: {os.path.isdir('/tmp')}")

# =============================================================================
# 2. pathlib 模块
# =============================================================================

print("\n=== pathlib 模块 ===")

from pathlib import Path

# 路径创建
current = Path.cwd()
home = Path.home()
file_path = Path(__file__)

print(f"当前目录: {current}")
print(f"home: {home}")
print(f"当前文件: {file_path}")

# 路径操作
print(f"\n路径属性:")
print(f"  name: {file_path.name}")
print(f"  stem: {file_path.stem}")
print(f"  suffix: {file_path.suffix}")
print(f"  parent: {file_path.parent}")
print(f"  parts: {file_path.parts[-3:]}")

# 路径拼接
new_path = current / "subdir" / "file.txt"
print(f"  拼接: {new_path}")

# glob 模式
print(f"\nglob 示例:")
for py_file in current.glob("*.py"):
    print(f"  {py_file.name}")

# =============================================================================
# 3. json 模块
# =============================================================================

print("\n=== json 模块 ===")

import json

# Python 对象转 JSON
data = {
    "name": "Alice",
    "age": 25,
    "scores": [85, 92, 78],
    "active": True,
    "address": None
}

# 序列化
json_str = json.dumps(data, indent=2, ensure_ascii=False)
print(f"json.dumps:\n{json_str}")

# 反序列化
parsed = json.loads(json_str)
print(f"\njson.loads: {parsed}")

# 自定义序列化
from datetime import datetime, date


class CustomEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, (datetime, date)):
            return obj.isoformat()
        return super().default(obj)


data_with_date = {"created": datetime.now()}
json_str = json.dumps(data_with_date, cls=CustomEncoder)
print(f"\n自定义编码: {json_str}")

# =============================================================================
# 4. collections 模块
# =============================================================================

print("\n=== collections 模块 ===")

from collections import (
    Counter, defaultdict, OrderedDict,
    namedtuple, deque, ChainMap
)

# Counter
print("Counter:")
text = "abracadabra"
counter = Counter(text)
print(f"  Counter('{text}'): {counter}")
print(f"  most_common(3): {counter.most_common(3)}")
print(f"  elements: {''.join(sorted(counter.elements()))}")

# 计数器运算
c1 = Counter(a=3, b=1)
c2 = Counter(a=1, b=2)
print(f"  c1 + c2: {c1 + c2}")
print(f"  c1 - c2: {c1 - c2}")

# defaultdict
print("\ndefaultdict:")
dd = defaultdict(list)
for k, v in [('a', 1), ('b', 2), ('a', 3)]:
    dd[k].append(v)
print(f"  {dict(dd)}")

dd_int = defaultdict(int)
for char in "hello":
    dd_int[char] += 1
print(f"  计数: {dict(dd_int)}")

# namedtuple
print("\nnamedtuple:")
Point = namedtuple('Point', ['x', 'y'])
p = Point(3, 4)
print(f"  Point(3, 4): {p}")
print(f"  p.x={p.x}, p.y={p.y}")
print(f"  _asdict(): {p._asdict()}")
print(f"  _replace(x=10): {p._replace(x=10)}")

# deque
print("\ndeque:")
dq = deque([1, 2, 3], maxlen=5)
dq.append(4)
dq.appendleft(0)
print(f"  deque: {list(dq)}")
dq.rotate(2)
print(f"  rotate(2): {list(dq)}")

# ChainMap
print("\nChainMap:")
defaults = {'color': 'red', 'size': 'medium'}
custom = {'color': 'blue'}
combined = ChainMap(custom, defaults)
print(f"  color: {combined['color']}")  # 从 custom 获取
print(f"  size: {combined['size']}")    # 从 defaults 获取

# =============================================================================
# 5. functools 模块
# =============================================================================

print("\n=== functools 模块 ===")

import functools

# partial - 偏函数
def power(base, exp):
    return base ** exp


square = functools.partial(power, exp=2)
cube = functools.partial(power, exp=3)
print(f"partial:")
print(f"  square(5): {square(5)}")
print(f"  cube(5): {cube(5)}")

# lru_cache - 缓存
@functools.lru_cache(maxsize=128)
def fibonacci(n):
    if n < 2:
        return n
    return fibonacci(n-1) + fibonacci(n-2)


print(f"\nlru_cache:")
print(f"  fibonacci(30): {fibonacci(30)}")
print(f"  cache_info: {fibonacci.cache_info()}")

# reduce
from functools import reduce

numbers = [1, 2, 3, 4, 5]
product = reduce(lambda x, y: x * y, numbers)
print(f"\nreduce:")
print(f"  乘积: {product}")

# singledispatch
@functools.singledispatch
def process(value):
    return f"默认: {value}"


@process.register(int)
def _(value):
    return f"整数: {value * 2}"


@process.register(str)
def _(value):
    return f"字符串: {value.upper()}"


print(f"\nsingledispatch:")
print(f"  process(5): {process(5)}")
print(f"  process('hello'): {process('hello')}")
print(f"  process(3.14): {process(3.14)}")

# =============================================================================
# 6. itertools 模块
# =============================================================================

print("\n=== itertools 模块 ===")

import itertools

# 无限迭代器
print("无限迭代器:")
print(f"  count(10): {list(itertools.islice(itertools.count(10), 5))}")
print(f"  cycle('AB'): {list(itertools.islice(itertools.cycle('AB'), 6))}")
print(f"  repeat(1, 3): {list(itertools.repeat(1, 3))}")

# 组合迭代器
print("\n组合迭代器:")
print(f"  chain([1,2], [3,4]): {list(itertools.chain([1,2], [3,4]))}")
print(f"  zip_longest: {list(itertools.zip_longest([1,2], [3,4,5], fillvalue=0))}")

# 排列组合
print("\n排列组合:")
print(f"  permutations('AB'): {list(itertools.permutations('AB'))}")
print(f"  combinations('ABC', 2): {list(itertools.combinations('ABC', 2))}")
print(f"  product('AB', repeat=2): {list(itertools.product('AB', repeat=2))}")

# 分组
print("\ngroupby:")
data = [('A', 1), ('A', 2), ('B', 3), ('B', 4)]
for key, group in itertools.groupby(data, key=lambda x: x[0]):
    print(f"  {key}: {list(group)}")

# =============================================================================
# 7. re 模块（正则表达式）
# =============================================================================

print("\n=== re 模块 ===")

import re

text = "Email: alice@example.com, bob@test.org"

# 基本匹配
pattern = r'\w+@\w+\.\w+'
matches = re.findall(pattern, text)
print(f"findall: {matches}")

# 分组
pattern = r'(\w+)@(\w+)\.(\w+)'
match = re.search(pattern, text)
if match:
    print(f"search groups: {match.groups()}")

# 替换
result = re.sub(r'\w+@\w+\.\w+', '[EMAIL]', text)
print(f"sub: {result}")

# 分割
result = re.split(r'[,\s]+', text)
print(f"split: {result}")

# 编译正则表达式
email_pattern = re.compile(r'\w+@\w+\.\w+', re.IGNORECASE)
matches = email_pattern.findall(text)
print(f"compiled pattern: {matches}")

# =============================================================================
# 8. hashlib 模块
# =============================================================================

print("\n=== hashlib 模块 ===")

import hashlib

text = "Hello, Python!"

# MD5
md5_hash = hashlib.md5(text.encode()).hexdigest()
print(f"MD5: {md5_hash}")

# SHA-256
sha256_hash = hashlib.sha256(text.encode()).hexdigest()
print(f"SHA-256: {sha256_hash}")

# 文件哈希
def file_hash(filepath, algorithm='sha256'):
    h = hashlib.new(algorithm)
    with open(filepath, 'rb') as f:
        for chunk in iter(lambda: f.read(4096), b''):
            h.update(chunk)
    return h.hexdigest()


print(f"当前文件 SHA-256: {file_hash(__file__)[:32]}...")

# =============================================================================
# 9. random 模块
# =============================================================================

print("\n=== random 模块 ===")

import random

# 基本随机数
print(f"random(): {random.random():.4f}")
print(f"randint(1, 10): {random.randint(1, 10)}")
print(f"randrange(0, 10, 2): {random.randrange(0, 10, 2)}")
print(f"uniform(1, 10): {random.uniform(1, 10):.4f}")

# 序列操作
items = [1, 2, 3, 4, 5]
print(f"choice({items}): {random.choice(items)}")
print(f"choices({items}, k=3): {random.choices(items, k=3)}")
print(f"sample({items}, k=3): {random.sample(items, k=3)}")

shuffled = items.copy()
random.shuffle(shuffled)
print(f"shuffle: {shuffled}")

# 设置种子（可重现）
random.seed(42)
print(f"seed(42) -> random(): {random.random():.4f}")

# =============================================================================
# 10. logging 模块
# =============================================================================

print("\n=== logging 模块 ===")

import logging

# 基本配置
logging.basicConfig(
    level=logging.DEBUG,
    format='%(asctime)s - %(levelname)s - %(message)s',
    datefmt='%H:%M:%S'
)

# 创建 logger
logger = logging.getLogger(__name__)

# 不同级别的日志
logger.debug("调试信息")
logger.info("一般信息")
logger.warning("警告信息")
logger.error("错误信息")

# =============================================================================
# 11. argparse 模块
# =============================================================================

print("\n=== argparse 模块 ===")

import argparse

# 创建解析器
parser = argparse.ArgumentParser(description='示例程序')
parser.add_argument('--name', type=str, default='World', help='名字')
parser.add_argument('--count', type=int, default=1, help='次数')
parser.add_argument('--verbose', '-v', action='store_true', help='详细模式')

# 解析参数（使用空列表避免解析实际命令行）
args = parser.parse_args(['--name', 'Python', '--count', '3', '-v'])
print(f"args: {args}")
print(f"  name: {args.name}")
print(f"  count: {args.count}")
print(f"  verbose: {args.verbose}")

# =============================================================================
# 12. copy 模块
# =============================================================================

print("\n=== copy 模块 ===")

import copy

original = [[1, 2], [3, 4]]

# 浅拷贝
shallow = copy.copy(original)
shallow[0].append(5)
print(f"浅拷贝后 original: {original}")  # 内层列表被修改

# 深拷贝
original = [[1, 2], [3, 4]]
deep = copy.deepcopy(original)
deep[0].append(5)
print(f"深拷贝后 original: {original}")  # 不受影响

# =============================================================================
# 13. dataclasses 模块
# =============================================================================

print("\n=== dataclasses 模块 ===")

from dataclasses import dataclass, field, asdict, astuple


@dataclass
class Product:
    name: str
    price: float
    quantity: int = 0
    tags: list = field(default_factory=list)

    @property
    def total_value(self):
        return self.price * self.quantity


product = Product("Widget", 9.99, 100, ["sale", "new"])
print(f"Product: {product}")
print(f"total_value: {product.total_value}")
print(f"asdict: {asdict(product)}")
print(f"astuple: {astuple(product)}")


if __name__ == "__main__":
    print("\n" + "=" * 50)
    print("17_stdlib.py 运行完成！")
    print("=" * 50)
