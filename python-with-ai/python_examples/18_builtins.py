#!/usr/bin/env python3
"""
18_builtins.py - 内置函数

Python 常用内置函数：len、map、filter、zip、sorted、enumerate 等
"""

# =============================================================================
# 1. 类型转换函数
# =============================================================================

print("=== 类型转换函数 ===")

# int()
print(f"int('42'): {int('42')}")
print(f"int(3.7): {int(3.7)}")
print(f"int('ff', 16): {int('ff', 16)}")
print(f"int('1010', 2): {int('1010', 2)}")

# float()
print(f"\nfloat('3.14'): {float('3.14')}")
print(f"float(42): {float(42)}")

# str()
print(f"\nstr(42): {str(42)}")
print(f"str([1, 2, 3]): {str([1, 2, 3])}")

# bool()
print(f"\nbool(1): {bool(1)}")
print(f"bool(0): {bool(0)}")
print(f"bool(''): {bool('')}")
print(f"bool([1]): {bool([1])}")

# list(), tuple(), set(), dict()
print(f"\nlist('abc'): {list('abc')}")
print(f"tuple([1, 2, 3]): {tuple([1, 2, 3])}")
print(f"set([1, 2, 2, 3]): {set([1, 2, 2, 3])}")
print(f"dict([('a', 1), ('b', 2)]): {dict([('a', 1), ('b', 2)])}")

# =============================================================================
# 2. 数学函数
# =============================================================================

print("\n=== 数学函数 ===")

numbers = [3, 1, 4, 1, 5, 9, 2, 6]

print(f"numbers: {numbers}")
print(f"len(): {len(numbers)}")
print(f"sum(): {sum(numbers)}")
print(f"min(): {min(numbers)}")
print(f"max(): {max(numbers)}")

# abs()
print(f"\nabs(-5): {abs(-5)}")
print(f"abs(3+4j): {abs(3+4j)}")  # 复数的模

# round()
print(f"\nround(3.7): {round(3.7)}")
print(f"round(3.5): {round(3.5)}")  # 银行家舍入
print(f"round(3.14159, 2): {round(3.14159, 2)}")

# pow()
print(f"\npow(2, 10): {pow(2, 10)}")
print(f"pow(2, 10, 1000): {pow(2, 10, 1000)}")  # 模幂运算

# divmod()
print(f"\ndivmod(17, 5): {divmod(17, 5)}")

# =============================================================================
# 3. 序列操作函数
# =============================================================================

print("\n=== 序列操作函数 ===")

numbers = [3, 1, 4, 1, 5, 9, 2, 6]

# sorted()
print(f"sorted({numbers}): {sorted(numbers)}")
print(f"sorted(reverse=True): {sorted(numbers, reverse=True)}")

words = ["apple", "Banana", "cherry"]
print(f"sorted 自定义: {sorted(words, key=str.lower)}")

# reversed()
print(f"\nlist(reversed({numbers})): {list(reversed(numbers))}")

# enumerate()
print("\nenumerate():")
for i, v in enumerate(['a', 'b', 'c'], start=1):
    print(f"  {i}: {v}")

# zip()
names = ["Alice", "Bob", "Charlie"]
scores = [85, 92, 78]
print(f"\nzip():")
for name, score in zip(names, scores):
    print(f"  {name}: {score}")

# 不等长序列
print(f"zip 不等长: {list(zip([1, 2], [3, 4, 5]))}")

# zip 解包
zipped = [(1, 'a'), (2, 'b'), (3, 'c')]
nums, chars = zip(*zipped)
print(f"zip 解包: nums={nums}, chars={chars}")

# =============================================================================
# 4. 函数式编程函数
# =============================================================================

print("\n=== 函数式编程函数 ===")

numbers = [1, 2, 3, 4, 5]

# map()
squared = list(map(lambda x: x ** 2, numbers))
print(f"map(lambda x: x**2, {numbers}): {squared}")

# 多参数 map
a = [1, 2, 3]
b = [10, 20, 30]
print(f"map(add, {a}, {b}): {list(map(lambda x, y: x + y, a, b))}")

# filter()
evens = list(filter(lambda x: x % 2 == 0, numbers))
print(f"\nfilter(偶数, {numbers}): {evens}")

# filter(None, ...) 过滤假值
mixed = [0, 1, '', 'hello', None, [], [1]]
truthy = list(filter(None, mixed))
print(f"filter(None, {mixed}): {truthy}")

# reduce() (在 functools 中)
from functools import reduce

product = reduce(lambda x, y: x * y, numbers)
print(f"\nreduce(乘法, {numbers}): {product}")

# =============================================================================
# 5. any() 和 all()
# =============================================================================

print("\n=== any() 和 all() ===")

values = [True, False, True]
print(f"values: {values}")
print(f"any(): {any(values)}")  # 任一为 True
print(f"all(): {all(values)}")  # 全部为 True

# 实际应用
numbers = [2, 4, 6, 8]
print(f"\n检查 {numbers}:")
print(f"  所有都是偶数: {all(n % 2 == 0 for n in numbers)}")
print(f"  存在大于5的数: {any(n > 5 for n in numbers)}")

# 空序列
print(f"\n空序列: any([])={any([])}, all([])={all([])}")

# =============================================================================
# 6. iter() 和 next()
# =============================================================================

print("\n=== iter() 和 next() ===")

# 基本用法
my_list = [1, 2, 3]
it = iter(my_list)
print(f"next(): {next(it)}")
print(f"next(): {next(it)}")
print(f"next(): {next(it)}")

# 默认值
it = iter([1, 2])
print(f"\nnext with default: {next(it, 'end')}")
print(f"next with default: {next(it, 'end')}")
print(f"next with default: {next(it, 'end')}")

# 哨兵模式
import random

random.seed(42)


def get_random():
    return random.randint(1, 10)


# iter(callable, sentinel) - 调用直到返回 sentinel
print(f"\n哨兵模式（直到返回 5）:")
for num in iter(get_random, 5):
    print(f"  {num}", end=" ")
print()

# =============================================================================
# 7. range()
# =============================================================================

print("\n=== range() ===")

print(f"range(5): {list(range(5))}")
print(f"range(2, 8): {list(range(2, 8))}")
print(f"range(0, 10, 2): {list(range(0, 10, 2))}")
print(f"range(10, 0, -2): {list(range(10, 0, -2))}")

# range 是惰性的
r = range(1000000)
print(f"\nrange(1000000) 内存高效")
print(f"  999999 in r: {999999 in r}")
print(f"  len(r): {len(r)}")

# =============================================================================
# 8. 对象属性函数
# =============================================================================

print("\n=== 对象属性函数 ===")


class MyClass:
    x = 10

    def __init__(self, y):
        self.y = y


obj = MyClass(20)

# getattr(), setattr(), delattr(), hasattr()
print(f"getattr(obj, 'x'): {getattr(obj, 'x')}")
print(f"getattr(obj, 'z', 'default'): {getattr(obj, 'z', 'default')}")
print(f"hasattr(obj, 'y'): {hasattr(obj, 'y')}")

setattr(obj, 'z', 30)
print(f"setattr 后 obj.z: {obj.z}")

# dir()
print(f"\ndir(obj) 部分: {[x for x in dir(obj) if not x.startswith('_')]}")

# vars()
print(f"vars(obj): {vars(obj)}")

# isinstance() 和 issubclass()
print(f"\nisinstance(obj, MyClass): {isinstance(obj, MyClass)}")
print(f"isinstance(1, (int, float)): {isinstance(1, (int, float))}")
print(f"issubclass(bool, int): {issubclass(bool, int)}")

# type()
print(f"\ntype(obj): {type(obj)}")
print(f"type(42): {type(42)}")

# =============================================================================
# 9. 输入/输出函数
# =============================================================================

print("\n=== 输入/输出函数 ===")

# print() 高级用法
print("普通输出")
print("分隔符", "示例", sep=" | ")
print("不换行", end="")
print(" -> 继续")

# 输出到文件
import io

buffer = io.StringIO()
print("写入buffer", file=buffer)
print(f"buffer内容: {buffer.getvalue()!r}")

# repr() vs str()
s = "Hello\nWorld"
print(f"\nstr(s): {str(s)}")
print(f"repr(s): {repr(s)}")

# ascii()
print(f"ascii('中文'): {ascii('中文')}")

# format()
print(f"\nformat(3.14159, '.2f'): {format(3.14159, '.2f')}")
print(f"format(255, 'x'): {format(255, 'x')}")
print(f"format(255, 'b'): {format(255, 'b')}")

# =============================================================================
# 10. 反射和元编程函数
# =============================================================================

print("\n=== 反射和元编程函数 ===")

# id() - 对象标识
a = [1, 2, 3]
b = a
c = [1, 2, 3]
print(f"id(a): {id(a)}")
print(f"a is b: {a is b}")
print(f"a is c: {a is c}")

# callable()
print(f"\ncallable(print): {callable(print)}")
print(f"callable(42): {callable(42)}")


# exec() 和 eval()
print("\neval():")
result = eval("2 + 3 * 4")
print(f"  eval('2 + 3 * 4'): {result}")

print("\nexec():")
exec("x = 10\nprint(f'  exec 中 x = {x}')")

# compile()
code = compile("print('编译的代码')", "<string>", "exec")
exec(code)

# globals() 和 locals()
print(f"\n部分 globals: {list(globals().keys())[:5]}")


def show_locals():
    local_var = 42
    print(f"locals(): {locals()}")


show_locals()

# =============================================================================
# 11. 其他常用函数
# =============================================================================

print("\n=== 其他常用函数 ===")

# bin(), oct(), hex()
n = 255
print(f"bin({n}): {bin(n)}")
print(f"oct({n}): {oct(n)}")
print(f"hex({n}): {hex(n)}")

# chr() 和 ord()
print(f"\nchr(65): {chr(65)}")
print(f"ord('A'): {ord('A')}")
print(f"chr(0x4e2d): {chr(0x4e2d)}")

# hash()
print(f"\nhash('hello'): {hash('hello')}")
print(f"hash((1, 2, 3)): {hash((1, 2, 3))}")

# slice()
s = slice(1, 5, 2)
lst = [0, 1, 2, 3, 4, 5, 6]
print(f"\nslice(1, 5, 2) on {lst}: {lst[s]}")

# memoryview()
data = bytearray(b"hello")
view = memoryview(data)
print(f"\nmemoryview: {view[0]} = {chr(view[0])}")

# __import__() - 动态导入
math_module = __import__('math')
print(f"\n动态导入 math.pi: {math_module.pi}")

# =============================================================================
# 12. 内置函数速查表
# =============================================================================

print("\n=== 内置函数速查表 ===")

cheatsheet = """
类型转换:    int, float, str, bool, list, tuple, set, dict, bytes
数学运算:    abs, round, pow, divmod, sum, min, max, len
序列操作:    sorted, reversed, enumerate, zip, range, slice
函数式:      map, filter, reduce(functools), any, all
迭代器:      iter, next
I/O:        print, input, open, format, repr, ascii
对象:       type, isinstance, issubclass, id, hash, callable
属性:       getattr, setattr, delattr, hasattr, dir, vars
元编程:     eval, exec, compile, globals, locals, __import__
进制转换:    bin, oct, hex, chr, ord
帮助:       help, __doc__
"""
print(cheatsheet)


if __name__ == "__main__":
    print("\n" + "=" * 50)
    print("18_builtins.py 运行完成！")
    print("=" * 50)
