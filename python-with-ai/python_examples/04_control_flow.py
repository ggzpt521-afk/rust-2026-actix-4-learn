#!/usr/bin/env python3
"""
04_control_flow.py - 流程控制

Python 流程控制语句：if/elif/else, for, while, match-case（Python 3.10+）
"""

# =============================================================================
# 1. 条件语句（if/elif/else）
# =============================================================================

print("=== 条件语句 ===")

score = 85

if score >= 90:
    grade = "A"
elif score >= 80:
    grade = "B"
elif score >= 70:
    grade = "C"
elif score >= 60:
    grade = "D"
else:
    grade = "F"

print(f"分数: {score}, 等级: {grade}")

# 单行条件表达式（三元运算符）
status = "及格" if score >= 60 else "不及格"
print(f"状态: {status}")

# 条件表达式链
age = 25
category = "儿童" if age < 12 else "青少年" if age < 18 else "成人" if age < 60 else "老年"
print(f"年龄 {age}, 类别: {category}")

# =============================================================================
# 2. 逻辑运算符
# =============================================================================

print("\n=== 逻辑运算符 ===")

x, y = 10, 20

# and, or, not
print(f"x > 5 and y > 15: {x > 5 and y > 15}")
print(f"x > 15 or y > 15: {x > 15 or y > 15}")
print(f"not x > 15: {not x > 15}")

# 短路求值
print(f"\n短路求值:")
print(f"True or print('不会执行'): {True or print('不会执行')}")
print(f"False and print('不会执行'): {False and print('不会执行')}")

# and/or 返回决定结果的操作数（不一定是 bool）
print(f"\nand/or 返回值:")
print(f"'hello' and 'world' = {'hello' and 'world'}")  # 'world'
print(f"'' and 'world' = {repr('' and 'world')}")  # ''
print(f"'hello' or 'world' = {'hello' or 'world'}")  # 'hello'
print(f"'' or 'world' = {'' or 'world'}")  # 'world'

# 常用技巧：设置默认值
name = ""
display_name = name or "Anonymous"
print(f"display_name = {display_name}")

# =============================================================================
# 3. 比较运算符
# =============================================================================

print("\n=== 比较运算符 ===")

a, b, c = 5, 10, 15

# 链式比较
print(f"a < b < c: {a < b < c}")  # 等价于 a < b and b < c
print(f"a < b > c: {a < b > c}")  # False

# is 与 == 的区别
list1 = [1, 2, 3]
list2 = [1, 2, 3]
list3 = list1

print(f"\nlist1 == list2: {list1 == list2}")  # True（值相等）
print(f"list1 is list2: {list1 is list2}")  # False（不同对象）
print(f"list1 is list3: {list1 is list3}")  # True（同一对象）

# None 比较应该用 is
value = None
print(f"value is None: {value is None}")  # 推荐
print(f"value == None: {value == None}")  # 不推荐

# =============================================================================
# 4. for 循环
# =============================================================================

print("\n=== for 循环 ===")

# 遍历列表
fruits = ["apple", "banana", "cherry"]
for fruit in fruits:
    print(f"  {fruit}")

# range() 函数
print("\nrange(5):")
for i in range(5):
    print(f"  i = {i}")

print("\nrange(2, 8, 2):")  # start, stop, step
for i in range(2, 8, 2):
    print(f"  i = {i}")

# enumerate() - 同时获取索引和值
print("\nenumerate():")
for index, fruit in enumerate(fruits):
    print(f"  {index}: {fruit}")

# 指定起始索引
print("\nenumerate(start=1):")
for index, fruit in enumerate(fruits, start=1):
    print(f"  {index}: {fruit}")

# zip() - 并行遍历多个序列
names = ["Alice", "Bob", "Charlie"]
scores = [85, 92, 78]
cities = ["Beijing", "Shanghai", "Shenzhen"]

print("\nzip():")
for name, score, city in zip(names, scores, cities):
    print(f"  {name}: {score}分, {city}")

# 遍历字典
print("\n遍历字典:")
person = {"name": "Alice", "age": 25, "city": "Beijing"}

for key in person:  # 遍历键
    print(f"  {key}: {person[key]}")

print("\n.items():")
for key, value in person.items():  # 遍历键值对
    print(f"  {key} = {value}")

# =============================================================================
# 5. while 循环
# =============================================================================

print("\n=== while 循环 ===")

count = 0
while count < 5:
    print(f"  count = {count}")
    count += 1

# while-else（循环正常结束时执行 else）
print("\nwhile-else:")
n = 0
while n < 3:
    print(f"  n = {n}")
    n += 1
else:
    print("  循环正常结束")

# =============================================================================
# 6. break 和 continue
# =============================================================================

print("\n=== break 和 continue ===")

# break - 跳出循环
print("break 示例:")
for i in range(10):
    if i == 5:
        print(f"  在 i={i} 处 break")
        break
    print(f"  i = {i}")

# continue - 跳过当前迭代
print("\ncontinue 示例:")
for i in range(6):
    if i == 3:
        print(f"  跳过 i={i}")
        continue
    print(f"  i = {i}")

# for-else（没有 break 时执行 else）
print("\nfor-else 示例:")
for n in range(2, 10):
    for x in range(2, n):
        if n % x == 0:
            break
    else:
        # 没有找到因子，是素数
        print(f"  {n} 是素数")

# =============================================================================
# 7. match-case（Python 3.10+）
# =============================================================================

print("\n=== match-case (Python 3.10+) ===")


def http_status(status):
    """根据 HTTP 状态码返回描述"""
    match status:
        case 200:
            return "OK"
        case 201:
            return "Created"
        case 400:
            return "Bad Request"
        case 404:
            return "Not Found"
        case 500:
            return "Internal Server Error"
        case _:  # 默认情况（通配符）
            return "Unknown Status"


print(f"HTTP 200: {http_status(200)}")
print(f"HTTP 404: {http_status(404)}")
print(f"HTTP 999: {http_status(999)}")


# 模式匹配 - 多个值
def day_type(day):
    match day:
        case "Saturday" | "Sunday":  # OR 模式
            return "Weekend"
        case "Monday" | "Tuesday" | "Wednesday" | "Thursday" | "Friday":
            return "Weekday"
        case _:
            return "Invalid"


print(f"\nSaturday: {day_type('Saturday')}")
print(f"Monday: {day_type('Monday')}")


# 模式匹配 - 序列解构
def process_point(point):
    match point:
        case (0, 0):
            return "原点"
        case (0, y):
            return f"Y轴上，y={y}"
        case (x, 0):
            return f"X轴上，x={x}"
        case (x, y):
            return f"普通点 ({x}, {y})"
        case _:
            return "不是有效的点"


print(f"\n(0, 0): {process_point((0, 0))}")
print(f"(0, 5): {process_point((0, 5))}")
print(f"(3, 4): {process_point((3, 4))}")


# 模式匹配 - 字典解构
def process_command(command):
    match command:
        case {"action": "create", "name": name}:
            return f"创建: {name}"
        case {"action": "delete", "id": id}:
            return f"删除 ID: {id}"
        case {"action": "update", "id": id, "data": data}:
            return f"更新 ID {id}: {data}"
        case _:
            return "未知命令"


print(f"\n{process_command({'action': 'create', 'name': 'test'})}")
print(f"{process_command({'action': 'delete', 'id': 123})}")


# 模式匹配 - 带守卫条件
def classify_number(n):
    match n:
        case n if n < 0:
            return "负数"
        case 0:
            return "零"
        case n if n % 2 == 0:
            return "正偶数"
        case _:
            return "正奇数"


print(f"\n-5: {classify_number(-5)}")
print(f"0: {classify_number(0)}")
print(f"4: {classify_number(4)}")
print(f"7: {classify_number(7)}")


# 模式匹配 - 类实例
class Point:
    def __init__(self, x, y):
        self.x = x
        self.y = y


def describe_point(point):
    match point:
        case Point(x=0, y=0):
            return "原点"
        case Point(x=x, y=y) if x == y:
            return f"对角线上的点 ({x}, {y})"
        case Point(x=x, y=y):
            return f"点 ({x}, {y})"
        case _:
            return "不是 Point 对象"


print(f"\nPoint(0, 0): {describe_point(Point(0, 0))}")
print(f"Point(5, 5): {describe_point(Point(5, 5))}")
print(f"Point(3, 4): {describe_point(Point(3, 4))}")

# =============================================================================
# 8. 列表推导式与条件
# =============================================================================

print("\n=== 列表推导式与条件 ===")

# 带条件的列表推导式
numbers = range(1, 11)
evens = [x for x in numbers if x % 2 == 0]
print(f"偶数: {evens}")

# 带 if-else 的列表推导式
labels = ["偶数" if x % 2 == 0 else "奇数" for x in range(1, 6)]
print(f"标签: {labels}")

# 嵌套循环
matrix = [[1, 2, 3], [4, 5, 6], [7, 8, 9]]
flattened = [num for row in matrix for num in row]
print(f"扁平化: {flattened}")

# =============================================================================
# 9. 海象运算符（Python 3.8+）
# =============================================================================

print("\n=== 海象运算符 := (Python 3.8+) ===")

# 在条件中赋值并使用
data = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

# 传统方式
# n = len(data)
# if n > 5:
#     print(f"列表有 {n} 个元素")

# 使用海象运算符
if (n := len(data)) > 5:
    print(f"列表有 {n} 个元素")

# 在 while 循环中使用
print("\n从列表中弹出元素:")
stack = [1, 2, 3]
while (item := stack.pop() if stack else None) is not None:
    print(f"  弹出: {item}")

# 在列表推导式中使用（避免重复计算）
results = [y for x in range(5) if (y := x ** 2) > 5]
print(f"平方大于5的结果: {results}")


if __name__ == "__main__":
    print("\n" + "=" * 50)
    print("04_control_flow.py 运行完成！")
    print("=" * 50)
