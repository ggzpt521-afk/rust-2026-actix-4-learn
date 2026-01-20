#!/usr/bin/env python3
"""
05_collections.py - 组合数据类型

Python 内置的组合数据类型：列表(list)、元组(tuple)、集合(set)、字典(dict)
"""

# =============================================================================
# 1. 列表（List）- 可变有序序列
# =============================================================================

print("=== 列表 (List) ===")

# 创建列表
empty_list = []
numbers = [1, 2, 3, 4, 5]
mixed = [1, "hello", 3.14, True, None]
nested = [[1, 2], [3, 4], [5, 6]]

print(f"numbers = {numbers}")
print(f"mixed = {mixed}")
print(f"nested = {nested}")

# 列表推导式
squares = [x ** 2 for x in range(1, 6)]
print(f"squares = {squares}")

# list() 构造函数
from_range = list(range(5))
from_string = list("hello")
print(f"from_range = {from_range}")
print(f"from_string = {from_string}")

# 索引和切片
print(f"\n索引和切片:")
print(f"  numbers[0] = {numbers[0]}")        # 第一个元素
print(f"  numbers[-1] = {numbers[-1]}")      # 最后一个元素
print(f"  numbers[1:4] = {numbers[1:4]}")    # 切片 [1, 4)
print(f"  numbers[::2] = {numbers[::2]}")    # 步长为 2
print(f"  numbers[::-1] = {numbers[::-1]}")  # 反转

# 列表方法
print(f"\n列表方法:")
lst = [3, 1, 4, 1, 5, 9, 2, 6]
print(f"原列表: {lst}")

lst.append(7)  # 末尾添加
print(f"append(7): {lst}")

lst.insert(0, 0)  # 指定位置插入
print(f"insert(0, 0): {lst}")

lst.extend([8, 9])  # 扩展列表
print(f"extend([8, 9]): {lst}")

popped = lst.pop()  # 弹出末尾元素
print(f"pop() -> {popped}: {lst}")

lst.remove(1)  # 删除第一个匹配的值
print(f"remove(1): {lst}")

# 排序
lst_copy = lst.copy()
lst_copy.sort()
print(f"sort(): {lst_copy}")

lst_copy.sort(reverse=True)
print(f"sort(reverse=True): {lst_copy}")

# 不修改原列表的排序
print(f"sorted(lst): {sorted(lst)}")
print(f"原列表不变: {lst}")

# 其他方法
print(f"\ncount(4): {lst.count(4)}")
print(f"index(5): {lst.index(5)}")
print(f"len(lst): {len(lst)}")

# 清空和反转
lst.reverse()
print(f"reverse(): {lst}")

# =============================================================================
# 2. 元组（Tuple）- 不可变有序序列
# =============================================================================

print("\n=== 元组 (Tuple) ===")

# 创建元组
empty_tuple = ()
single = (1,)  # 单元素元组必须加逗号
point = (3, 4)
rgb = (255, 128, 0)
mixed_tuple = (1, "hello", 3.14)

print(f"single = {single}, type = {type(single)}")
print(f"point = {point}")

# 不加括号也可以创建元组
coordinates = 10, 20, 30
print(f"coordinates = {coordinates}")

# 元组解包
x, y = point
print(f"解包: x = {x}, y = {y}")

# 带星号的解包
first, *middle, last = [1, 2, 3, 4, 5]
print(f"first = {first}, middle = {middle}, last = {last}")

# 元组是不可变的
try:
    point[0] = 10
except TypeError as e:
    print(f"元组不可变: {e}")

# 但元组中的可变对象可以修改
mutable_in_tuple = ([1, 2], [3, 4])
mutable_in_tuple[0].append(3)
print(f"元组中的列表可修改: {mutable_in_tuple}")

# 元组方法（只有 count 和 index）
print(f"\n元组方法:")
t = (1, 2, 2, 3, 2, 4)
print(f"t.count(2) = {t.count(2)}")
print(f"t.index(3) = {t.index(3)}")

# 命名元组
from collections import namedtuple

Point = namedtuple("Point", ["x", "y"])
p = Point(3, 4)
print(f"\n命名元组:")
print(f"Point(3, 4) = {p}")
print(f"p.x = {p.x}, p.y = {p.y}")
print(f"p[0] = {p[0]}, p[1] = {p[1]}")

# =============================================================================
# 3. 集合（Set）- 可变无序不重复集合
# =============================================================================

print("\n=== 集合 (Set) ===")

# 创建集合
empty_set = set()  # 不能用 {}，那是空字典
numbers_set = {1, 2, 3, 4, 5}
from_list = set([1, 2, 2, 3, 3, 3])  # 自动去重

print(f"numbers_set = {numbers_set}")
print(f"from_list = {from_list}")

# 集合推导式
squares_set = {x ** 2 for x in range(1, 6)}
print(f"squares_set = {squares_set}")

# 集合操作
a = {1, 2, 3, 4, 5}
b = {4, 5, 6, 7, 8}

print(f"\na = {a}")
print(f"b = {b}")

print(f"并集 a | b = {a | b}")
print(f"并集 a.union(b) = {a.union(b)}")

print(f"交集 a & b = {a & b}")
print(f"交集 a.intersection(b) = {a.intersection(b)}")

print(f"差集 a - b = {a - b}")
print(f"差集 a.difference(b) = {a.difference(b)}")

print(f"对称差集 a ^ b = {a ^ b}")
print(f"对称差集 a.symmetric_difference(b) = {a.symmetric_difference(b)}")

# 子集和超集
c = {1, 2, 3}
print(f"\nc = {c}")
print(f"c <= a (子集): {c <= a}")
print(f"c.issubset(a): {c.issubset(a)}")
print(f"a >= c (超集): {a >= c}")
print(f"a.issuperset(c): {a.issuperset(c)}")

# 集合方法
s = {1, 2, 3}
print(f"\n集合方法:")

s.add(4)
print(f"add(4): {s}")

s.update([5, 6])
print(f"update([5, 6]): {s}")

s.discard(10)  # 不存在不报错
print(f"discard(10): {s}")

s.remove(6)  # 不存在会报错
print(f"remove(6): {s}")

popped = s.pop()  # 随机弹出一个元素
print(f"pop() -> {popped}: {s}")

# 不可变集合（frozenset）
fs = frozenset([1, 2, 3])
print(f"\nfrozenset: {fs}")
print(f"可以作为字典的键: {{{fs}: 'value'}}")

# =============================================================================
# 4. 字典（Dict）- 可变键值对集合
# =============================================================================

print("\n=== 字典 (Dict) ===")

# 创建字典
empty_dict = {}
person = {"name": "Alice", "age": 25, "city": "Beijing"}
from_pairs = dict([("a", 1), ("b", 2)])
from_kwargs = dict(x=1, y=2, z=3)

print(f"person = {person}")
print(f"from_pairs = {from_pairs}")
print(f"from_kwargs = {from_kwargs}")

# 字典推导式
squares_dict = {x: x ** 2 for x in range(1, 6)}
print(f"squares_dict = {squares_dict}")

# 访问元素
print(f"\n访问元素:")
print(f"person['name'] = {person['name']}")
print(f"person.get('age') = {person.get('age')}")
print(f"person.get('salary', 0) = {person.get('salary', 0)}")  # 默认值

# 修改和添加
person["age"] = 26
person["email"] = "alice@example.com"
print(f"修改后: {person}")

# 删除元素
del person["email"]
print(f"del 后: {person}")

popped = person.pop("city", None)
print(f"pop('city') -> {popped}: {person}")

# 字典方法
d = {"a": 1, "b": 2, "c": 3}
print(f"\n字典方法:")
print(f"d.keys() = {list(d.keys())}")
print(f"d.values() = {list(d.values())}")
print(f"d.items() = {list(d.items())}")

# 合并字典
d1 = {"a": 1, "b": 2}
d2 = {"b": 3, "c": 4}

# Python 3.9+ 使用 | 运算符
merged = d1 | d2
print(f"\nd1 | d2 = {merged}")

# update 方法（就地修改）
d1.update(d2)
print(f"d1.update(d2): {d1}")

# setdefault - 获取值，不存在则设置默认值
d = {"a": 1}
value = d.setdefault("b", 2)
print(f"\nsetdefault('b', 2) -> {value}: {d}")
value = d.setdefault("a", 100)  # 已存在，不会修改
print(f"setdefault('a', 100) -> {value}: {d}")

# 字典视图
print(f"\n字典视图是动态的:")
d = {"a": 1}
keys = d.keys()
print(f"keys = {list(keys)}")
d["b"] = 2
print(f"添加 'b' 后 keys = {list(keys)}")

# =============================================================================
# 5. 序列操作通用方法
# =============================================================================

print("\n=== 序列通用操作 ===")

seq = [1, 2, 3, 4, 5]

print(f"len(seq) = {len(seq)}")
print(f"min(seq) = {min(seq)}")
print(f"max(seq) = {max(seq)}")
print(f"sum(seq) = {sum(seq)}")

print(f"3 in seq: {3 in seq}")
print(f"10 not in seq: {10 not in seq}")

# 序列拼接和重复
print(f"[1, 2] + [3, 4] = {[1, 2] + [3, 4]}")
print(f"[1, 2] * 3 = {[1, 2] * 3}")

# =============================================================================
# 6. 深拷贝与浅拷贝
# =============================================================================

print("\n=== 深拷贝与浅拷贝 ===")

import copy

original = [[1, 2], [3, 4]]

# 浅拷贝
shallow = original.copy()  # 或 list(original) 或 original[:]
shallow[0].append(5)

print(f"original = {original}")  # [[1, 2, 5], [3, 4]] - 内层列表被修改
print(f"shallow = {shallow}")

# 深拷贝
original = [[1, 2], [3, 4]]
deep = copy.deepcopy(original)
deep[0].append(5)

print(f"original = {original}")  # [[1, 2], [3, 4]] - 不受影响
print(f"deep = {deep}")

# =============================================================================
# 7. collections 模块常用容器
# =============================================================================

print("\n=== collections 模块 ===")

from collections import Counter, defaultdict, OrderedDict, deque

# Counter - 计数器
print("Counter:")
text = "abracadabra"
counter = Counter(text)
print(f"Counter('{text}') = {counter}")
print(f"most_common(3) = {counter.most_common(3)}")

# defaultdict - 带默认值的字典
print("\ndefaultdict:")
dd = defaultdict(list)
dd["fruits"].append("apple")
dd["fruits"].append("banana")
dd["vegetables"].append("carrot")
print(f"defaultdict(list): {dict(dd)}")

dd_int = defaultdict(int)
for char in "hello":
    dd_int[char] += 1
print(f"defaultdict(int): {dict(dd_int)}")

# deque - 双端队列
print("\ndeque:")
dq = deque([1, 2, 3])
dq.appendleft(0)
dq.append(4)
print(f"deque: {dq}")
dq.rotate(1)  # 右旋
print(f"rotate(1): {dq}")
dq.rotate(-2)  # 左旋
print(f"rotate(-2): {dq}")

# 限制长度的 deque
limited = deque(maxlen=3)
for i in range(5):
    limited.append(i)
    print(f"  append({i}): {list(limited)}")


if __name__ == "__main__":
    print("\n" + "=" * 50)
    print("05_collections.py 运行完成！")
    print("=" * 50)
