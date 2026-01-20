#!/usr/bin/env python3
"""
09_exceptions.py - 异常处理

Python 异常处理机制：try/except/else/finally、自定义异常、异常链等
"""

# =============================================================================
# 1. 基本异常处理
# =============================================================================

print("=== 基本异常处理 ===")

# try-except 基本结构
try:
    result = 10 / 0
except ZeroDivisionError:
    print("捕获到除零错误！")

# 捕获异常对象
try:
    result = 10 / 0
except ZeroDivisionError as e:
    print(f"异常信息: {e}")
    print(f"异常类型: {type(e).__name__}")

# =============================================================================
# 2. 多个异常处理
# =============================================================================

print("\n=== 多个异常处理 ===")


def risky_operation(value):
    """可能抛出多种异常的函数"""
    try:
        number = int(value)
        result = 10 / number
        my_list = [1, 2, 3]
        return my_list[number]
    except ValueError:
        print("值错误：无法转换为整数")
    except ZeroDivisionError:
        print("除零错误：不能除以零")
    except IndexError:
        print("索引错误：列表索引越界")
    except (TypeError, AttributeError) as e:
        # 同时捕获多种异常
        print(f"类型或属性错误: {e}")
    return None


risky_operation("abc")  # ValueError
risky_operation("0")    # ZeroDivisionError
risky_operation("10")   # IndexError

# =============================================================================
# 3. try-except-else-finally
# =============================================================================

print("\n=== try-except-else-finally ===")


def divide(a, b):
    """完整的异常处理结构"""
    try:
        result = a / b
    except ZeroDivisionError:
        print("  except: 除零错误")
        return None
    except TypeError:
        print("  except: 类型错误")
        return None
    else:
        # 没有异常时执行
        print(f"  else: 计算成功，结果是 {result}")
        return result
    finally:
        # 无论是否有异常都会执行
        print("  finally: 清理工作")


print("divide(10, 2):")
divide(10, 2)

print("\ndivide(10, 0):")
divide(10, 0)

# =============================================================================
# 4. 捕获所有异常
# =============================================================================

print("\n=== 捕获所有异常 ===")

# 捕获所有异常（不推荐，但有时需要）
try:
    # 某些危险操作
    x = 1 / 0
except Exception as e:
    # 捕获所有 Exception 子类
    print(f"捕获到异常: {type(e).__name__}: {e}")

# BaseException 是所有异常的基类
# 包括 SystemExit, KeyboardInterrupt, GeneratorExit
# 通常不应该捕获这些

# =============================================================================
# 5. 抛出异常
# =============================================================================

print("\n=== 抛出异常 ===")


def validate_age(age):
    """验证年龄"""
    if not isinstance(age, int):
        raise TypeError("Age must be an integer")
    if age < 0:
        raise ValueError("Age cannot be negative")
    if age > 150:
        raise ValueError("Age seems unrealistic")
    return True


# 测试
for test_age in [25, -5, "abc"]:
    try:
        validate_age(test_age)
        print(f"  {test_age}: 有效")
    except (TypeError, ValueError) as e:
        print(f"  {test_age}: {e}")

# =============================================================================
# 6. 自定义异常
# =============================================================================

print("\n=== 自定义异常 ===")


class ValidationError(Exception):
    """验证错误基类"""
    pass


class EmailValidationError(ValidationError):
    """邮箱验证错误"""

    def __init__(self, email, message="Invalid email format"):
        self.email = email
        self.message = message
        super().__init__(f"{message}: {email}")


class PasswordValidationError(ValidationError):
    """密码验证错误"""

    def __init__(self, message, requirements=None):
        self.message = message
        self.requirements = requirements or []
        super().__init__(message)


def validate_email(email):
    if "@" not in email:
        raise EmailValidationError(email, "Missing @ symbol")
    if not email.endswith((".com", ".org", ".net")):
        raise EmailValidationError(email, "Invalid domain")
    return True


def validate_password(password):
    errors = []
    if len(password) < 8:
        errors.append("At least 8 characters")
    if not any(c.isupper() for c in password):
        errors.append("At least one uppercase letter")
    if not any(c.isdigit() for c in password):
        errors.append("At least one digit")

    if errors:
        raise PasswordValidationError(
            "Password doesn't meet requirements",
            requirements=errors
        )
    return True


# 测试自定义异常
test_data = [
    ("email", "invalid-email"),
    ("email", "test@example.com"),
    ("password", "weak"),
    ("password", "StrongPass123"),
]

for field, value in test_data:
    try:
        if field == "email":
            validate_email(value)
        else:
            validate_password(value)
        print(f"  {field} '{value}': 有效")
    except EmailValidationError as e:
        print(f"  Email '{e.email}': {e.message}")
    except PasswordValidationError as e:
        print(f"  Password: {e.message}")
        for req in e.requirements:
            print(f"    - {req}")

# =============================================================================
# 7. 异常链
# =============================================================================

print("\n=== 异常链 ===")


class DatabaseError(Exception):
    """数据库错误"""
    pass


def fetch_user(user_id):
    """模拟数据库操作"""
    try:
        # 模拟数据库操作失败
        data = {"users": {}}
        return data["users"][user_id]
    except KeyError as e:
        # raise from 保留原始异常信息
        raise DatabaseError(f"User {user_id} not found") from e


try:
    fetch_user(123)
except DatabaseError as e:
    print(f"DatabaseError: {e}")
    print(f"原始异常: {e.__cause__}")

# 禁止异常链
# raise DatabaseError("...") from None

# =============================================================================
# 8. 异常信息获取
# =============================================================================

print("\n=== 异常信息获取 ===")

import traceback
import sys


def problematic_function():
    return 1 / 0


try:
    problematic_function()
except ZeroDivisionError:
    # 获取异常信息
    exc_type, exc_value, exc_tb = sys.exc_info()
    print(f"类型: {exc_type.__name__}")
    print(f"值: {exc_value}")

    # 格式化堆栈跟踪
    print("\n堆栈跟踪:")
    tb_lines = traceback.format_exception(exc_type, exc_value, exc_tb)
    for line in tb_lines:
        print(line, end="")

# =============================================================================
# 9. 断言（assert）
# =============================================================================

print("\n=== 断言 ===")


def calculate_average(numbers):
    """使用断言进行前置条件检查"""
    assert len(numbers) > 0, "List cannot be empty"
    assert all(isinstance(n, (int, float)) for n in numbers), "All elements must be numbers"
    return sum(numbers) / len(numbers)


# 正常情况
result = calculate_average([1, 2, 3, 4, 5])
print(f"平均值: {result}")

# 断言失败
try:
    calculate_average([])
except AssertionError as e:
    print(f"断言失败: {e}")

# 注意：断言可以被禁用（python -O）
# 不应该用于数据验证，只用于调试

# =============================================================================
# 10. warnings 模块
# =============================================================================

print("\n=== 警告 ===")

import warnings


def deprecated_function():
    """已弃用的函数"""
    warnings.warn(
        "deprecated_function is deprecated, use new_function instead",
        DeprecationWarning,
        stacklevel=2
    )
    return "old result"


# 显示警告
warnings.filterwarnings("always", category=DeprecationWarning)
result = deprecated_function()

# 将警告转为异常
# warnings.filterwarnings("error", category=DeprecationWarning)

# =============================================================================
# 11. 上下文管理器中的异常
# =============================================================================

print("\n=== 上下文管理器中的异常 ===")


class ManagedResource:
    """演示上下文管理器中的异常处理"""

    def __enter__(self):
        print("  获取资源")
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        print("  释放资源")
        if exc_type is not None:
            print(f"  处理异常: {exc_type.__name__}: {exc_val}")
            # 返回 True 表示异常已处理，不再向上传播
            # 返回 False 或 None 表示异常继续传播
        return False

    def do_something(self):
        raise ValueError("Something went wrong")


try:
    with ManagedResource() as resource:
        resource.do_something()
except ValueError as e:
    print(f"  外部捕获: {e}")

# =============================================================================
# 12. 常见内置异常
# =============================================================================

print("\n=== 常见内置异常 ===")

exceptions_info = """
BaseException
├── SystemExit          # sys.exit() 抛出
├── KeyboardInterrupt   # Ctrl+C
├── GeneratorExit       # 生成器关闭
└── Exception
    ├── StopIteration       # 迭代结束
    ├── ArithmeticError
    │   ├── ZeroDivisionError
    │   ├── FloatingPointError
    │   └── OverflowError
    ├── AssertionError      # assert 失败
    ├── AttributeError      # 属性不存在
    ├── BufferError
    ├── EOFError           # input() 遇到 EOF
    ├── ImportError
    │   └── ModuleNotFoundError
    ├── LookupError
    │   ├── IndexError      # 索引越界
    │   └── KeyError        # 键不存在
    ├── MemoryError
    ├── NameError           # 名称未定义
    │   └── UnboundLocalError
    ├── OSError
    │   ├── FileNotFoundError
    │   ├── PermissionError
    │   └── TimeoutError
    ├── RuntimeError
    │   ├── NotImplementedError
    │   └── RecursionError
    ├── SyntaxError
    │   └── IndentationError
    ├── TypeError           # 类型错误
    ├── ValueError          # 值错误
    └── Warning
        ├── DeprecationWarning
        ├── FutureWarning
        └── UserWarning
"""
print(exceptions_info)

# =============================================================================
# 13. 最佳实践
# =============================================================================

print("=== 异常处理最佳实践 ===")

best_practices = """
1. 具体异常：捕获具体的异常类型，而不是使用裸的 except

   # 不好
   try:
       ...
   except:
       pass

   # 好
   try:
       ...
   except ValueError as e:
       handle_error(e)

2. EAFP vs LBYL
   EAFP (Easier to Ask Forgiveness than Permission) - Pythonic 风格
   LBYL (Look Before You Leap) - 传统风格

   # LBYL
   if key in dictionary:
       return dictionary[key]

   # EAFP (推荐)
   try:
       return dictionary[key]
   except KeyError:
       return default

3. 不要吞掉异常：至少记录日志

   try:
       ...
   except Exception as e:
       logging.error(f"Error: {e}")
       raise  # 重新抛出

4. 使用 finally 或上下文管理器确保资源释放

5. 自定义异常应继承自 Exception，而不是 BaseException

6. 异常信息应该有意义，便于调试
"""
print(best_practices)


if __name__ == "__main__":
    print("\n" + "=" * 50)
    print("09_exceptions.py 运行完成！")
    print("=" * 50)
