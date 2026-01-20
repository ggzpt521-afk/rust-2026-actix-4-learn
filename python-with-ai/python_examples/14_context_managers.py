#!/usr/bin/env python3
"""
14_context_managers.py - 上下文管理器

Python 上下文管理器：with 语句、__enter__/__exit__、contextlib 等
"""

import contextlib
import time
import tempfile
import os
from pathlib import Path

# =============================================================================
# 1. 基本使用
# =============================================================================

print("=== 上下文管理器基本使用 ===")

# 文件操作是最常见的上下文管理器
with open(__file__, "r") as f:
    first_line = f.readline()
    print(f"第一行: {first_line.strip()}")

# 等价于：
# f = open(__file__, "r")
# try:
#     first_line = f.readline()
# finally:
#     f.close()

# 多个上下文管理器
with open(__file__, "r") as f1, tempfile.NamedTemporaryFile("w", delete=False) as f2:
    content = f1.read()[:100]
    f2.write(content)
    temp_path = f2.name

print(f"临时文件: {temp_path}")
os.unlink(temp_path)  # 清理

# =============================================================================
# 2. 自定义上下文管理器（类）
# =============================================================================

print("\n=== 自定义上下文管理器（类）===")


class Timer:
    """计时器上下文管理器"""

    def __init__(self, name="Timer"):
        self.name = name
        self.elapsed = 0

    def __enter__(self):
        """进入 with 块时调用"""
        print(f"  [{self.name}] 开始计时")
        self.start = time.perf_counter()
        return self  # 返回值赋给 as 后的变量

    def __exit__(self, exc_type, exc_val, exc_tb):
        """退出 with 块时调用"""
        self.elapsed = time.perf_counter() - self.start
        print(f"  [{self.name}] 结束，耗时: {self.elapsed:.4f}秒")
        # 返回 False（或 None）: 异常继续传播
        # 返回 True: 异常被抑制
        return False


with Timer("计算") as t:
    total = sum(range(1000000))
print(f"  结果: {total}, 计时器对象: {t.elapsed:.4f}秒")


class ManagedResource:
    """管理资源的上下文管理器"""

    def __init__(self, name):
        self.name = name
        print(f"  创建资源: {name}")

    def __enter__(self):
        print(f"  获取资源: {self.name}")
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        print(f"  释放资源: {self.name}")
        if exc_type is not None:
            print(f"  处理异常: {exc_type.__name__}: {exc_val}")
        return False  # 不抑制异常

    def do_work(self):
        print(f"  使用资源: {self.name}")


print("\n正常使用:")
with ManagedResource("DB连接") as resource:
    resource.do_work()

print("\n发生异常:")
try:
    with ManagedResource("文件") as resource:
        raise ValueError("模拟错误")
except ValueError:
    print("  外部捕获异常")

# =============================================================================
# 3. 使用 contextlib（推荐）
# =============================================================================

print("\n=== contextlib 模块 ===")


# @contextmanager 装饰器
@contextlib.contextmanager
def timer_context(name="Timer"):
    """使用生成器创建上下文管理器"""
    print(f"  [{name}] 开始")
    start = time.perf_counter()
    try:
        yield {"name": name}  # yield 的值是 as 后的变量
    except Exception as e:
        print(f"  [{name}] 发生异常: {e}")
        raise
    finally:
        elapsed = time.perf_counter() - start
        print(f"  [{name}] 结束，耗时: {elapsed:.4f}秒")


with timer_context("测试") as ctx:
    print(f"  上下文: {ctx}")
    time.sleep(0.1)


# 更多 contextlib 工具
print("\n--- contextlib.suppress ---")
# 抑制指定异常
with contextlib.suppress(FileNotFoundError):
    os.remove("nonexistent_file.txt")  # 不会抛出异常
    print("  文件不存在，但异常被抑制")

print("  继续执行")


print("\n--- contextlib.redirect_stdout ---")
import io

# 重定向标准输出
buffer = io.StringIO()
with contextlib.redirect_stdout(buffer):
    print("这会被重定向到 buffer")
    print("而不是终端")

print(f"捕获的输出: {buffer.getvalue()!r}")


print("\n--- contextlib.closing ---")


# 为没有上下文管理器支持的对象添加 close() 调用
class LegacyResource:
    def __init__(self):
        print("  LegacyResource 打开")

    def close(self):
        print("  LegacyResource 关闭")


with contextlib.closing(LegacyResource()) as resource:
    print("  使用 LegacyResource")


print("\n--- contextlib.nullcontext ---")


# 空上下文管理器（占位符）
def process_data(data, lock=None):
    """可选的锁"""
    ctx = lock if lock else contextlib.nullcontext()
    with ctx:
        return data * 2


print(f"  无锁: {process_data(5)}")

# =============================================================================
# 4. 异步上下文管理器
# =============================================================================

print("\n=== 异步上下文管理器 ===")


class AsyncResource:
    """异步上下文管理器"""

    async def __aenter__(self):
        print("  异步获取资源")
        return self

    async def __aexit__(self, exc_type, exc_val, exc_tb):
        print("  异步释放资源")
        return False


# 使用 @asynccontextmanager
@contextlib.asynccontextmanager
async def async_timer(name="AsyncTimer"):
    """异步计时器"""
    print(f"  [{name}] 异步开始")
    start = time.perf_counter()
    try:
        yield
    finally:
        elapsed = time.perf_counter() - start
        print(f"  [{name}] 异步结束，耗时: {elapsed:.4f}秒")


# 运行异步代码
import asyncio


async def async_demo():
    async with AsyncResource():
        print("  使用异步资源")

    async with async_timer("异步测试"):
        await asyncio.sleep(0.1)


asyncio.run(async_demo())

# =============================================================================
# 5. ExitStack
# =============================================================================

print("\n=== ExitStack ===")


# 动态管理多个上下文管理器
@contextlib.contextmanager
def numbered_resource(n):
    print(f"  获取资源 {n}")
    try:
        yield n
    finally:
        print(f"  释放资源 {n}")


with contextlib.ExitStack() as stack:
    # 动态添加上下文管理器
    resources = [stack.enter_context(numbered_resource(i)) for i in range(3)]
    print(f"  资源列表: {resources}")


# 条件性上下文管理器
print("\n条件性使用:")


def conditional_context(use_timer):
    with contextlib.ExitStack() as stack:
        if use_timer:
            stack.enter_context(timer_context("条件计时器"))
        print("  执行任务")


conditional_context(True)
conditional_context(False)

# =============================================================================
# 6. 实用示例
# =============================================================================

print("\n=== 实用示例 ===")


# 数据库事务模拟
@contextlib.contextmanager
def transaction(db_name):
    """模拟数据库事务"""
    print(f"  开始事务: {db_name}")
    try:
        yield
        print(f"  提交事务: {db_name}")
    except Exception as e:
        print(f"  回滚事务: {db_name}, 原因: {e}")
        raise


print("事务管理器:")
try:
    with transaction("mydb"):
        print("  执行数据库操作")
        # raise ValueError("模拟错误")  # 取消注释测试回滚
except ValueError:
    pass


# 临时改变工作目录
@contextlib.contextmanager
def working_directory(path):
    """临时改变工作目录"""
    old_dir = os.getcwd()
    try:
        os.chdir(path)
        yield
    finally:
        os.chdir(old_dir)


print("\n临时目录切换:")
print(f"  当前目录: {os.getcwd()}")
with working_directory("/tmp"):
    print(f"  临时目录: {os.getcwd()}")
print(f"  恢复目录: {os.getcwd()}")


# 设置/恢复环境变量
@contextlib.contextmanager
def env_var(key, value):
    """临时设置环境变量"""
    old_value = os.environ.get(key)
    os.environ[key] = value
    try:
        yield
    finally:
        if old_value is None:
            del os.environ[key]
        else:
            os.environ[key] = old_value


print("\n临时环境变量:")
print(f"  DEBUG = {os.environ.get('DEBUG', 'None')}")
with env_var("DEBUG", "true"):
    print(f"  DEBUG = {os.environ.get('DEBUG')}")
print(f"  DEBUG = {os.environ.get('DEBUG', 'None')}")


# 锁管理器（线程安全）
import threading


@contextlib.contextmanager
def locked(lock):
    """带超时的锁获取"""
    acquired = lock.acquire(timeout=1)
    if not acquired:
        raise TimeoutError("无法获取锁")
    try:
        yield
    finally:
        lock.release()


print("\n锁管理器:")
my_lock = threading.Lock()
with locked(my_lock):
    print("  持有锁")


# 原子文件写入
@contextlib.contextmanager
def atomic_write(filepath):
    """原子文件写入（先写临时文件，成功后替换）"""
    temp_path = filepath + ".tmp"
    try:
        with open(temp_path, "w") as f:
            yield f
        os.replace(temp_path, filepath)  # 原子替换
    except:
        if os.path.exists(temp_path):
            os.unlink(temp_path)
        raise


print("\n原子文件写入:")
test_file = "/tmp/atomic_test.txt"
with atomic_write(test_file) as f:
    f.write("原子写入的内容")
print(f"  文件内容: {Path(test_file).read_text()}")
os.unlink(test_file)

# =============================================================================
# 7. 重入上下文管理器
# =============================================================================

print("\n=== 重入上下文管理器 ===")


class ReentrantLock:
    """可重入的锁"""

    def __init__(self):
        self._lock = threading.RLock()

    def __enter__(self):
        self._lock.acquire()
        return self

    def __exit__(self, *args):
        self._lock.release()


lock = ReentrantLock()
with lock:
    print("  第一层")
    with lock:  # 可以重入
        print("  第二层")
    print("  返回第一层")


if __name__ == "__main__":
    print("\n" + "=" * 50)
    print("14_context_managers.py 运行完成！")
    print("=" * 50)
