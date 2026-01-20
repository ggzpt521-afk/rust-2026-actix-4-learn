#!/usr/bin/env python3
"""
12_generators.py - 生成器与迭代器

Python 迭代协议、生成器函数、生成器表达式、itertools 等
"""

# =============================================================================
# 1. 迭代器协议
# =============================================================================

print("=== 迭代器协议 ===")


class CountDown:
    """自定义迭代器 - 倒计时"""

    def __init__(self, start):
        self.start = start

    def __iter__(self):
        """返回迭代器对象（自身）"""
        return self

    def __next__(self):
        """返回下一个值"""
        if self.start <= 0:
            raise StopIteration
        self.start -= 1
        return self.start + 1


# 使用迭代器
print("倒计时:")
for num in CountDown(5):
    print(f"  {num}")

# 手动迭代
print("\n手动迭代:")
counter = CountDown(3)
print(f"  next: {next(counter)}")
print(f"  next: {next(counter)}")
print(f"  next: {next(counter)}")
try:
    next(counter)
except StopIteration:
    print("  迭代结束")

# =============================================================================
# 2. 可迭代对象 vs 迭代器
# =============================================================================

print("\n=== 可迭代对象 vs 迭代器 ===")


class NumberRange:
    """可迭代对象（每次迭代返回新的迭代器）"""

    def __init__(self, start, end):
        self.start = start
        self.end = end

    def __iter__(self):
        """返回新的迭代器"""
        return NumberRangeIterator(self.start, self.end)


class NumberRangeIterator:
    """迭代器"""

    def __init__(self, start, end):
        self.current = start
        self.end = end

    def __iter__(self):
        return self

    def __next__(self):
        if self.current >= self.end:
            raise StopIteration
        value = self.current
        self.current += 1
        return value


numbers = NumberRange(1, 4)
print("第一次迭代:", list(numbers))
print("第二次迭代:", list(numbers))  # 可以重复迭代

# 检查是否可迭代
from collections.abc import Iterable, Iterator

print(f"\nNumberRange 是 Iterable: {isinstance(numbers, Iterable)}")
print(f"NumberRange 是 Iterator: {isinstance(numbers, Iterator)}")

# =============================================================================
# 3. 生成器函数
# =============================================================================

print("\n=== 生成器函数 ===")


def countdown(n):
    """生成器函数 - 使用 yield"""
    print("  开始倒计时...")
    while n > 0:
        yield n
        n -= 1
    print("  倒计时结束!")


# 调用生成器函数返回生成器对象
gen = countdown(5)
print(f"生成器类型: {type(gen)}")

# 迭代生成器
print("\n迭代生成器:")
for num in countdown(3):
    print(f"  {num}")


# 无限生成器
def infinite_sequence():
    """无限序列生成器"""
    num = 0
    while True:
        yield num
        num += 1


print("\n无限生成器（取前5个）:")
gen = infinite_sequence()
for _ in range(5):
    print(f"  {next(gen)}")


# 生成器处理大文件
def read_large_file(file_path):
    """生成器读取大文件（内存友好）"""
    with open(file_path, 'r') as f:
        for line in f:
            yield line.strip()


# =============================================================================
# 4. yield 表达式
# =============================================================================

print("\n=== yield 表达式 ===")


def fibonacci(limit):
    """斐波那契数列生成器"""
    a, b = 0, 1
    while a < limit:
        yield a
        a, b = b, a + b


print("斐波那契数列 (< 100):")
print(list(fibonacci(100)))


# yield from - 委托生成器
def chain(*iterables):
    """使用 yield from 链接多个可迭代对象"""
    for it in iterables:
        yield from it


print("\nyield from 链接:")
result = list(chain([1, 2], [3, 4], [5, 6]))
print(f"  {result}")


# 嵌套生成器
def flatten(nested_list):
    """递归展平嵌套列表"""
    for item in nested_list:
        if isinstance(item, list):
            yield from flatten(item)
        else:
            yield item


nested = [1, [2, 3, [4, 5]], 6, [7, [8, 9]]]
print(f"\n展平 {nested}:")
print(f"  {list(flatten(nested))}")

# =============================================================================
# 5. 生成器表达式
# =============================================================================

print("\n=== 生成器表达式 ===")

# 列表推导式（立即求值，占用内存）
list_comp = [x ** 2 for x in range(10)]
print(f"列表推导式: {list_comp}")
print(f"  类型: {type(list_comp)}")

# 生成器表达式（惰性求值，节省内存）
gen_exp = (x ** 2 for x in range(10))
print(f"生成器表达式: {gen_exp}")
print(f"  类型: {type(gen_exp)}")
print(f"  转为列表: {list(gen_exp)}")

# 内存对比
import sys

list_size = sys.getsizeof([x ** 2 for x in range(1000)])
gen_size = sys.getsizeof(x ** 2 for x in range(1000))
print(f"\n内存占用对比（1000个元素）:")
print(f"  列表: {list_size} bytes")
print(f"  生成器: {gen_size} bytes")

# 生成器表达式作为函数参数（可省略括号）
total = sum(x ** 2 for x in range(10))
print(f"\nsum(x**2 for x in range(10)) = {total}")

# =============================================================================
# 6. 生成器方法
# =============================================================================

print("\n=== 生成器方法 ===")


def echo_generator():
    """可接收值的生成器"""
    print("  生成器启动")
    while True:
        received = yield
        print(f"  收到: {received}")


# send() 方法
gen = echo_generator()
next(gen)  # 启动生成器，运行到第一个 yield
gen.send("Hello")
gen.send("World")


# 带返回值的生成器
def accumulator():
    """累加器生成器"""
    total = 0
    while True:
        value = yield total
        if value is None:
            break
        total += value
    return total


gen = accumulator()
next(gen)  # 启动
gen.send(10)
gen.send(20)
gen.send(30)

try:
    gen.send(None)  # 结束
except StopIteration as e:
    print(f"\n累加器返回值: {e.value}")


# throw() 方法
def generator_with_exception():
    """处理异常的生成器"""
    try:
        yield 1
        yield 2
        yield 3
    except ValueError:
        print("  捕获到 ValueError")
        yield "recovered"


gen = generator_with_exception()
print(f"\nnext: {next(gen)}")
print(f"next: {next(gen)}")
print(f"throw: {gen.throw(ValueError)}")

# close() 方法
gen = countdown(10)
next(gen)
gen.close()  # 关闭生成器
print("生成器已关闭")

# =============================================================================
# 7. itertools 模块
# =============================================================================

print("\n=== itertools 模块 ===")

import itertools

# count - 无限计数器
print("count(10, 2):", list(itertools.islice(itertools.count(10, 2), 5)))

# cycle - 无限循环
print("cycle('ABC'):", list(itertools.islice(itertools.cycle('ABC'), 7)))

# repeat - 重复
print("repeat('X', 3):", list(itertools.repeat('X', 3)))

# chain - 链接
print("chain([1,2], [3,4]):", list(itertools.chain([1, 2], [3, 4])))

# chain.from_iterable
print("chain.from_iterable:", list(itertools.chain.from_iterable([[1, 2], [3, 4]])))

# compress - 过滤
print("compress('ABCDE', [1,0,1,0,1]):", list(itertools.compress('ABCDE', [1, 0, 1, 0, 1])))

# dropwhile / takewhile
print("dropwhile(x<5, [1,4,6,4,1]):", list(itertools.dropwhile(lambda x: x < 5, [1, 4, 6, 4, 1])))
print("takewhile(x<5, [1,4,6,4,1]):", list(itertools.takewhile(lambda x: x < 5, [1, 4, 6, 4, 1])))

# filterfalse
print("filterfalse(x%2, range(10)):", list(itertools.filterfalse(lambda x: x % 2, range(10))))

# groupby
print("\ngroupby 示例:")
data = [('A', 1), ('A', 2), ('B', 3), ('B', 4), ('A', 5)]
data.sort(key=lambda x: x[0])  # groupby 需要先排序
for key, group in itertools.groupby(data, key=lambda x: x[0]):
    print(f"  {key}: {list(group)}")

# islice - 切片
print("\nislice(range(10), 2, 8, 2):", list(itertools.islice(range(10), 2, 8, 2)))

# starmap
print("starmap(pow, [(2,3), (3,2)]):", list(itertools.starmap(pow, [(2, 3), (3, 2)])))

# zip_longest
print("zip_longest([1,2], [1,2,3]):", list(itertools.zip_longest([1, 2], [1, 2, 3], fillvalue=0)))

# 排列组合
print("\n排列组合:")
print(f"  permutations('ABC', 2): {list(itertools.permutations('ABC', 2))}")
print(f"  combinations('ABC', 2): {list(itertools.combinations('ABC', 2))}")
print(f"  combinations_with_replacement('AB', 2): {list(itertools.combinations_with_replacement('AB', 2))}")
print(f"  product('AB', repeat=2): {list(itertools.product('AB', repeat=2))}")

# =============================================================================
# 8. 实用生成器示例
# =============================================================================

print("\n=== 实用生成器示例 ===")


def sliding_window(iterable, n):
    """滑动窗口"""
    it = iter(iterable)
    window = []
    for _ in range(n):
        window.append(next(it))
    yield tuple(window)
    for item in it:
        window.pop(0)
        window.append(item)
        yield tuple(window)


print("滑动窗口 (size=3):")
print(f"  {list(sliding_window(range(6), 3))}")


def batch(iterable, size):
    """分批处理"""
    it = iter(iterable)
    while True:
        chunk = list(itertools.islice(it, size))
        if not chunk:
            break
        yield chunk


print("\n分批处理 (size=3):")
print(f"  {list(batch(range(10), 3))}")


def unique_justseen(iterable):
    """去除连续重复"""
    return map(next, map(lambda x: x[1], itertools.groupby(iterable)))


print("\n去除连续重复:")
print(f"  {list(unique_justseen('AAAABBBCCDAABB'))}")


if __name__ == "__main__":
    print("\n" + "=" * 50)
    print("12_generators.py 运行完成！")
    print("=" * 50)
