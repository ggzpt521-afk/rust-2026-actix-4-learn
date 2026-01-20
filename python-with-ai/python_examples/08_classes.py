#!/usr/bin/env python3
"""
08_classes.py - 类与对象

Python 面向对象编程：类定义、继承、多态、特殊方法、属性、抽象类等
"""

# =============================================================================
# 1. 类的基本定义
# =============================================================================

print("=== 类的基本定义 ===")


class Person:
    """人员类 - 演示基本类定义"""

    # 类变量（所有实例共享）
    species = "Homo sapiens"
    count = 0

    def __init__(self, name, age):
        """构造方法"""
        # 实例变量
        self.name = name
        self.age = age
        Person.count += 1

    def greet(self):
        """实例方法"""
        return f"Hello, I'm {self.name}, {self.age} years old."

    def have_birthday(self):
        """修改实例状态"""
        self.age += 1
        return f"Happy birthday! Now I'm {self.age}."


# 创建实例
alice = Person("Alice", 25)
bob = Person("Bob", 30)

print(f"alice.greet(): {alice.greet()}")
print(f"bob.greet(): {bob.greet()}")
print(f"alice.have_birthday(): {alice.have_birthday()}")

# 类变量访问
print(f"\n类变量:")
print(f"Person.species = {Person.species}")
print(f"Person.count = {Person.count}")
print(f"alice.species = {alice.species}")  # 也可以通过实例访问

# =============================================================================
# 2. 访问控制（约定）
# =============================================================================

print("\n=== 访问控制 ===")


class BankAccount:
    """银行账户类 - 演示访问控制"""

    def __init__(self, owner, balance=0):
        self.owner = owner          # 公开属性
        self._balance = balance      # 受保护属性（约定）
        self.__pin = "1234"          # 私有属性（名称修饰）

    def deposit(self, amount):
        if amount > 0:
            self._balance += amount
        return self._balance

    def get_balance(self):
        return self._balance

    def _internal_method(self):
        """受保护方法（约定）"""
        return "Internal method"

    def __private_method(self):
        """私有方法"""
        return "Private method"


account = BankAccount("Alice", 1000)
print(f"owner: {account.owner}")
print(f"_balance: {account._balance}")  # 可以访问，但不建议

# 名称修饰后的私有属性
print(f"__pin (通过名称修饰): {account._BankAccount__pin}")

# =============================================================================
# 3. 属性装饰器（Property）
# =============================================================================

print("\n=== 属性装饰器 ===")


class Circle:
    """圆形类 - 演示 property 装饰器"""

    def __init__(self, radius):
        self._radius = radius

    @property
    def radius(self):
        """getter"""
        return self._radius

    @radius.setter
    def radius(self, value):
        """setter"""
        if value < 0:
            raise ValueError("Radius cannot be negative")
        self._radius = value

    @radius.deleter
    def radius(self):
        """deleter"""
        print("Deleting radius")
        del self._radius

    @property
    def diameter(self):
        """只读属性"""
        return self._radius * 2

    @property
    def area(self):
        """计算属性"""
        import math
        return math.pi * self._radius ** 2


circle = Circle(5)
print(f"radius = {circle.radius}")
print(f"diameter = {circle.diameter}")
print(f"area = {circle.area:.2f}")

circle.radius = 10  # 使用 setter
print(f"修改后 radius = {circle.radius}")

try:
    circle.radius = -1
except ValueError as e:
    print(f"设置负值: {e}")

# =============================================================================
# 4. 类方法和静态方法
# =============================================================================

print("\n=== 类方法和静态方法 ===")


class Date:
    """日期类 - 演示类方法和静态方法"""

    def __init__(self, year, month, day):
        self.year = year
        self.month = month
        self.day = day

    def __repr__(self):
        return f"Date({self.year}, {self.month}, {self.day})"

    @classmethod
    def from_string(cls, date_string):
        """类方法 - 工厂方法模式"""
        year, month, day = map(int, date_string.split("-"))
        return cls(year, month, day)

    @classmethod
    def today(cls):
        """类方法 - 返回今天的日期"""
        import datetime
        t = datetime.date.today()
        return cls(t.year, t.month, t.day)

    @staticmethod
    def is_valid_date(year, month, day):
        """静态方法 - 不需要访问实例或类"""
        if month < 1 or month > 12:
            return False
        if day < 1 or day > 31:
            return False
        return True


# 类方法调用
d1 = Date(2024, 1, 15)
d2 = Date.from_string("2024-06-20")
d3 = Date.today()

print(f"普通创建: {d1}")
print(f"from_string: {d2}")
print(f"today: {d3}")

# 静态方法调用
print(f"is_valid_date(2024, 2, 30): {Date.is_valid_date(2024, 2, 30)}")

# =============================================================================
# 5. 继承
# =============================================================================

print("\n=== 继承 ===")


class Animal:
    """动物基类"""

    def __init__(self, name):
        self.name = name

    def speak(self):
        raise NotImplementedError("Subclass must implement")

    def describe(self):
        return f"I am {self.name}"


class Dog(Animal):
    """狗类 - 继承自 Animal"""

    def __init__(self, name, breed):
        super().__init__(name)  # 调用父类构造方法
        self.breed = breed

    def speak(self):
        return f"{self.name} says: Woof!"

    def fetch(self):
        return f"{self.name} is fetching the ball"


class Cat(Animal):
    """猫类 - 继承自 Animal"""

    def speak(self):
        return f"{self.name} says: Meow!"

    def climb(self):
        return f"{self.name} is climbing"


dog = Dog("Buddy", "Golden Retriever")
cat = Cat("Whiskers")

print(f"dog.describe(): {dog.describe()}")
print(f"dog.speak(): {dog.speak()}")
print(f"dog.fetch(): {dog.fetch()}")
print(f"dog.breed: {dog.breed}")

print(f"\ncat.describe(): {cat.describe()}")
print(f"cat.speak(): {cat.speak()}")

# isinstance 和 issubclass
print(f"\nisinstance(dog, Dog): {isinstance(dog, Dog)}")
print(f"isinstance(dog, Animal): {isinstance(dog, Animal)}")
print(f"issubclass(Dog, Animal): {issubclass(Dog, Animal)}")

# =============================================================================
# 6. 多重继承
# =============================================================================

print("\n=== 多重继承 ===")


class Flyable:
    """可飞行的混入类"""

    def fly(self):
        return f"{self.name} is flying"


class Swimmable:
    """可游泳的混入类"""

    def swim(self):
        return f"{self.name} is swimming"


class Duck(Animal, Flyable, Swimmable):
    """鸭子类 - 多重继承"""

    def speak(self):
        return f"{self.name} says: Quack!"


duck = Duck("Donald")
print(f"duck.speak(): {duck.speak()}")
print(f"duck.fly(): {duck.fly()}")
print(f"duck.swim(): {duck.swim()}")

# MRO（方法解析顺序）
print(f"\nDuck MRO: {[c.__name__ for c in Duck.__mro__]}")

# =============================================================================
# 7. 多态
# =============================================================================

print("\n=== 多态 ===")


def animal_sound(animal):
    """多态示例 - 接受任何有 speak 方法的对象"""
    print(animal.speak())


animals = [Dog("Buddy", "Labrador"), Cat("Whiskers"), Duck("Donald")]
for animal in animals:
    animal_sound(animal)

# =============================================================================
# 8. 特殊方法（魔术方法）
# =============================================================================

print("\n=== 特殊方法 ===")


class Vector:
    """向量类 - 演示特殊方法"""

    def __init__(self, x, y):
        self.x = x
        self.y = y

    def __repr__(self):
        """开发者友好的字符串表示"""
        return f"Vector({self.x}, {self.y})"

    def __str__(self):
        """用户友好的字符串表示"""
        return f"({self.x}, {self.y})"

    def __eq__(self, other):
        """相等比较"""
        if isinstance(other, Vector):
            return self.x == other.x and self.y == other.y
        return False

    def __add__(self, other):
        """加法运算"""
        return Vector(self.x + other.x, self.y + other.y)

    def __sub__(self, other):
        """减法运算"""
        return Vector(self.x - other.x, self.y - other.y)

    def __mul__(self, scalar):
        """标量乘法"""
        return Vector(self.x * scalar, self.y * scalar)

    def __rmul__(self, scalar):
        """反向乘法（scalar * vector）"""
        return self.__mul__(scalar)

    def __abs__(self):
        """向量长度"""
        return (self.x ** 2 + self.y ** 2) ** 0.5

    def __bool__(self):
        """布尔值"""
        return self.x != 0 or self.y != 0

    def __len__(self):
        """维度"""
        return 2

    def __getitem__(self, index):
        """索引访问"""
        if index == 0:
            return self.x
        elif index == 1:
            return self.y
        raise IndexError("Vector index out of range")

    def __iter__(self):
        """迭代支持"""
        yield self.x
        yield self.y


v1 = Vector(3, 4)
v2 = Vector(1, 2)

print(f"v1 = {v1}")
print(f"repr(v1) = {repr(v1)}")
print(f"v1 + v2 = {v1 + v2}")
print(f"v1 - v2 = {v1 - v2}")
print(f"v1 * 2 = {v1 * 2}")
print(f"3 * v1 = {3 * v1}")
print(f"abs(v1) = {abs(v1)}")
print(f"len(v1) = {len(v1)}")
print(f"v1[0] = {v1[0]}, v1[1] = {v1[1]}")
print(f"list(v1) = {list(v1)}")
print(f"v1 == Vector(3, 4): {v1 == Vector(3, 4)}")

# =============================================================================
# 9. 数据类（dataclass）
# =============================================================================

print("\n=== 数据类 (dataclass) ===")

from dataclasses import dataclass, field


@dataclass
class Point:
    """数据类示例"""
    x: float
    y: float
    label: str = "origin"


@dataclass
class Rectangle:
    """带计算属性的数据类"""
    width: float
    height: float
    tags: list = field(default_factory=list)  # 可变默认值

    @property
    def area(self):
        return self.width * self.height


p1 = Point(3, 4)
p2 = Point(3, 4)

print(f"p1 = {p1}")
print(f"p1 == p2: {p1 == p2}")  # 自动生成 __eq__

rect = Rectangle(10, 20)
print(f"rect = {rect}")
print(f"rect.area = {rect.area}")

# =============================================================================
# 10. 抽象基类
# =============================================================================

print("\n=== 抽象基类 ===")

from abc import ABC, abstractmethod


class Shape(ABC):
    """形状抽象基类"""

    @abstractmethod
    def area(self):
        """计算面积 - 必须实现"""
        pass

    @abstractmethod
    def perimeter(self):
        """计算周长 - 必须实现"""
        pass

    def describe(self):
        """非抽象方法 - 可选重写"""
        return f"A shape with area {self.area():.2f}"


class Rectangle2(Shape):
    def __init__(self, width, height):
        self.width = width
        self.height = height

    def area(self):
        return self.width * self.height

    def perimeter(self):
        return 2 * (self.width + self.height)


class Circle2(Shape):
    def __init__(self, radius):
        self.radius = radius

    def area(self):
        import math
        return math.pi * self.radius ** 2

    def perimeter(self):
        import math
        return 2 * math.pi * self.radius


# 不能实例化抽象类
try:
    s = Shape()
except TypeError as e:
    print(f"不能实例化抽象类: {e}")

rect = Rectangle2(10, 5)
circle = Circle2(7)

print(f"\nrect.area() = {rect.area()}")
print(f"rect.perimeter() = {rect.perimeter()}")
print(f"rect.describe() = {rect.describe()}")

print(f"\ncircle.area() = {circle.area():.2f}")
print(f"circle.perimeter() = {circle.perimeter():.2f}")

# =============================================================================
# 11. __slots__ 优化
# =============================================================================

print("\n=== __slots__ 优化 ===")


class PointWithSlots:
    """使用 __slots__ 减少内存占用"""
    __slots__ = ['x', 'y']

    def __init__(self, x, y):
        self.x = x
        self.y = y


class PointWithoutSlots:
    """普通类"""

    def __init__(self, x, y):
        self.x = x
        self.y = y


p_slots = PointWithSlots(1, 2)
p_normal = PointWithoutSlots(1, 2)

print(f"PointWithSlots: x={p_slots.x}, y={p_slots.y}")

# 不能添加新属性
try:
    p_slots.z = 3
except AttributeError as e:
    print(f"slots 限制: {e}")

# 普通类可以
p_normal.z = 3
print(f"PointWithoutSlots 可以添加: z={p_normal.z}")

# =============================================================================
# 12. 描述符
# =============================================================================

print("\n=== 描述符 ===")


class Validator:
    """描述符 - 数据验证"""

    def __init__(self, min_value=None, max_value=None):
        self.min_value = min_value
        self.max_value = max_value

    def __set_name__(self, owner, name):
        self.name = name

    def __get__(self, obj, objtype=None):
        if obj is None:
            return self
        return obj.__dict__.get(self.name)

    def __set__(self, obj, value):
        if self.min_value is not None and value < self.min_value:
            raise ValueError(f"{self.name} must be >= {self.min_value}")
        if self.max_value is not None and value > self.max_value:
            raise ValueError(f"{self.name} must be <= {self.max_value}")
        obj.__dict__[self.name] = value


class Product:
    price = Validator(min_value=0)
    quantity = Validator(min_value=0, max_value=1000)

    def __init__(self, name, price, quantity):
        self.name = name
        self.price = price
        self.quantity = quantity


product = Product("Widget", 9.99, 100)
print(f"product: {product.name}, ${product.price}, qty: {product.quantity}")

try:
    product.price = -10
except ValueError as e:
    print(f"验证失败: {e}")


if __name__ == "__main__":
    print("\n" + "=" * 50)
    print("08_classes.py 运行完成！")
    print("=" * 50)
