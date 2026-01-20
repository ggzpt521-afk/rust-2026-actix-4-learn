#!/usr/bin/env python3
"""
16_concurrency.py - 并发与异步编程

Python 并发：threading、multiprocessing、asyncio、concurrent.futures
"""

import threading
import multiprocessing
import asyncio
import time
from concurrent.futures import ThreadPoolExecutor, ProcessPoolExecutor, as_completed
import queue

# =============================================================================
# 1. threading 基础
# =============================================================================

print("=== threading 基础 ===")


def worker(name, delay):
    """工作线程函数"""
    print(f"  线程 {name} 开始")
    time.sleep(delay)
    print(f"  线程 {name} 结束")
    return f"{name} 完成"


# 创建并启动线程
t1 = threading.Thread(target=worker, args=("A", 0.2))
t2 = threading.Thread(target=worker, args=("B", 0.1))

t1.start()
t2.start()

# 等待线程完成
t1.join()
t2.join()
print("所有线程完成")

# 线程属性
print(f"\n当前线程: {threading.current_thread().name}")
print(f"活动线程数: {threading.active_count()}")

# =============================================================================
# 2. 线程同步
# =============================================================================

print("\n=== 线程同步 ===")

# 共享变量和锁
counter = 0
lock = threading.Lock()


def increment_with_lock():
    global counter
    for _ in range(10000):
        with lock:  # 使用上下文管理器
            counter += 1


def increment_without_lock():
    global counter
    for _ in range(10000):
        counter += 1  # 不安全！


# 使用锁
counter = 0
threads = [threading.Thread(target=increment_with_lock) for _ in range(5)]
for t in threads:
    t.start()
for t in threads:
    t.join()
print(f"有锁计数器: {counter}")

# 不使用锁（可能不正确）
counter = 0
threads = [threading.Thread(target=increment_without_lock) for _ in range(5)]
for t in threads:
    t.start()
for t in threads:
    t.join()
print(f"无锁计数器: {counter} (可能不是 50000)")


# RLock - 可重入锁
rlock = threading.RLock()


def recursive_function(n):
    with rlock:
        if n > 0:
            recursive_function(n - 1)


# Semaphore - 信号量
semaphore = threading.Semaphore(3)  # 最多3个线程同时访问


def limited_access(name):
    with semaphore:
        print(f"  {name} 获得访问权")
        time.sleep(0.1)


print("\n信号量示例:")
threads = [threading.Thread(target=limited_access, args=(f"Thread-{i}",)) for i in range(5)]
for t in threads:
    t.start()
for t in threads:
    t.join()


# Event - 事件
event = threading.Event()


def waiter(name):
    print(f"  {name} 等待事件...")
    event.wait()
    print(f"  {name} 收到事件!")


print("\nEvent 示例:")
t = threading.Thread(target=waiter, args=("Waiter",))
t.start()
time.sleep(0.1)
print("  触发事件")
event.set()
t.join()


# Condition - 条件变量
condition = threading.Condition()
items = []


def producer():
    with condition:
        items.append("item")
        print("  生产了一个 item")
        condition.notify()


def consumer():
    with condition:
        while not items:
            condition.wait()
        item = items.pop()
        print(f"  消费了: {item}")


print("\nCondition 示例:")
c = threading.Thread(target=consumer)
p = threading.Thread(target=producer)
c.start()
time.sleep(0.1)
p.start()
c.join()
p.join()

# =============================================================================
# 3. 线程安全队列
# =============================================================================

print("\n=== 线程安全队列 ===")

task_queue = queue.Queue()
results = []


def worker_queue():
    while True:
        item = task_queue.get()
        if item is None:
            break
        results.append(item * 2)
        task_queue.task_done()


# 创建工作线程
threads = [threading.Thread(target=worker_queue) for _ in range(3)]
for t in threads:
    t.start()

# 添加任务
for i in range(10):
    task_queue.put(i)

# 等待所有任务完成
task_queue.join()

# 停止工作线程
for _ in threads:
    task_queue.put(None)
for t in threads:
    t.join()

print(f"队列处理结果: {sorted(results)}")

# =============================================================================
# 4. concurrent.futures
# =============================================================================

print("\n=== concurrent.futures ===")


def compute(n):
    """模拟计算密集型任务"""
    time.sleep(0.1)
    return n * n


# ThreadPoolExecutor
print("ThreadPoolExecutor:")
with ThreadPoolExecutor(max_workers=4) as executor:
    # 提交单个任务
    future = executor.submit(compute, 5)
    print(f"  单个结果: {future.result()}")

    # 批量提交
    futures = [executor.submit(compute, i) for i in range(5)]
    for future in as_completed(futures):
        print(f"  完成: {future.result()}")

    # 使用 map
    results = list(executor.map(compute, range(5)))
    print(f"  map 结果: {results}")


# ProcessPoolExecutor（多进程）
def cpu_intensive(n):
    """CPU 密集型任务"""
    return sum(i * i for i in range(n))


print("\nProcessPoolExecutor:")
if __name__ == "__main__":  # 多进程需要这个保护
    with ProcessPoolExecutor(max_workers=2) as executor:
        results = list(executor.map(cpu_intensive, [10000, 20000, 30000]))
        print(f"  结果: {results}")

# =============================================================================
# 5. asyncio 基础
# =============================================================================

print("\n=== asyncio 基础 ===")


async def async_worker(name, delay):
    """异步工作函数"""
    print(f"  {name} 开始")
    await asyncio.sleep(delay)
    print(f"  {name} 结束")
    return f"{name} 完成"


async def main_async():
    """主异步函数"""
    # 顺序执行
    print("顺序执行:")
    result1 = await async_worker("Task1", 0.1)
    result2 = await async_worker("Task2", 0.1)

    # 并发执行
    print("\n并发执行:")
    results = await asyncio.gather(
        async_worker("Task-A", 0.2),
        async_worker("Task-B", 0.1),
        async_worker("Task-C", 0.15),
    )
    print(f"  结果: {results}")


asyncio.run(main_async())

# =============================================================================
# 6. asyncio 高级特性
# =============================================================================

print("\n=== asyncio 高级特性 ===")


# create_task - 创建任务
async def example_tasks():
    print("create_task 示例:")

    async def background_task():
        await asyncio.sleep(0.1)
        return "后台任务完成"

    task = asyncio.create_task(background_task())
    print("  任务已创建，继续执行其他代码")
    result = await task
    print(f"  {result}")


asyncio.run(example_tasks())


# 超时处理
async def example_timeout():
    print("\n超时处理示例:")

    async def slow_operation():
        await asyncio.sleep(1)
        return "完成"

    try:
        result = await asyncio.wait_for(slow_operation(), timeout=0.1)
    except asyncio.TimeoutError:
        print("  操作超时!")


asyncio.run(example_timeout())


# asyncio.Queue
async def example_queue():
    print("\nasyncio.Queue 示例:")
    q = asyncio.Queue()

    async def producer(q):
        for i in range(3):
            await q.put(i)
            print(f"  生产: {i}")

    async def consumer(q):
        while True:
            item = await q.get()
            print(f"  消费: {item}")
            q.task_done()
            if item == 2:
                break

    await asyncio.gather(producer(q), consumer(q))


asyncio.run(example_queue())


# Semaphore 限制并发
async def example_semaphore():
    print("\n异步信号量示例:")
    semaphore = asyncio.Semaphore(2)

    async def limited_task(name):
        async with semaphore:
            print(f"  {name} 开始")
            await asyncio.sleep(0.1)
            print(f"  {name} 结束")

    await asyncio.gather(*[limited_task(f"Task-{i}") for i in range(4)])


asyncio.run(example_semaphore())

# =============================================================================
# 7. 异步上下文管理器和迭代器
# =============================================================================

print("\n=== 异步上下文管理器和迭代器 ===")


class AsyncResource:
    """异步上下文管理器"""

    async def __aenter__(self):
        print("  获取异步资源")
        await asyncio.sleep(0.05)
        return self

    async def __aexit__(self, *args):
        print("  释放异步资源")
        await asyncio.sleep(0.05)


class AsyncRange:
    """异步迭代器"""

    def __init__(self, n):
        self.n = n
        self.i = 0

    def __aiter__(self):
        return self

    async def __anext__(self):
        if self.i >= self.n:
            raise StopAsyncIteration
        await asyncio.sleep(0.05)
        result = self.i
        self.i += 1
        return result


async def async_context_demo():
    # 异步上下文管理器
    async with AsyncResource():
        print("  使用资源")

    # 异步迭代器
    print("\n异步迭代:")
    async for i in AsyncRange(3):
        print(f"  i = {i}")


asyncio.run(async_context_demo())


# 异步生成器
async def async_generator(n):
    """异步生成器"""
    for i in range(n):
        await asyncio.sleep(0.05)
        yield i


async def async_gen_demo():
    print("\n异步生成器:")
    async for value in async_generator(3):
        print(f"  value = {value}")


asyncio.run(async_gen_demo())

# =============================================================================
# 8. 在同步代码中运行异步函数
# =============================================================================

print("\n=== 同步/异步互操作 ===")


async def async_operation():
    await asyncio.sleep(0.1)
    return "异步结果"


# 从同步代码调用异步函数
print("从同步代码调用异步函数:")
result = asyncio.run(async_operation())
print(f"  结果: {result}")


# 在异步中运行同步函数
def blocking_io():
    time.sleep(0.1)
    return "IO 完成"


async def run_blocking():
    print("\n在异步中运行同步阻塞函数:")
    loop = asyncio.get_event_loop()
    result = await loop.run_in_executor(None, blocking_io)
    print(f"  结果: {result}")


asyncio.run(run_blocking())

# =============================================================================
# 9. 多进程基础
# =============================================================================

print("\n=== 多进程基础 ===")


def process_worker(name):
    """进程工作函数"""
    print(f"  进程 {name} (PID: {multiprocessing.current_process().pid})")
    return name


# 只在主进程中执行
if __name__ == "__main__":
    # 创建进程
    processes = []
    for i in range(3):
        p = multiprocessing.Process(target=process_worker, args=(f"P{i}",))
        processes.append(p)
        p.start()

    for p in processes:
        p.join()

    print(f"CPU 核心数: {multiprocessing.cpu_count()}")

# =============================================================================
# 10. 选择合适的并发方式
# =============================================================================

print("\n=== 并发方式选择指南 ===")

guide = """
1. threading（多线程）
   - 适用于 I/O 密集型任务
   - 受 GIL 限制，不适合 CPU 密集型
   - 共享内存，需要锁同步
   - 例如：文件操作、网络请求

2. multiprocessing（多进程）
   - 适用于 CPU 密集型任务
   - 绕过 GIL
   - 进程间通信开销较大
   - 例如：数据处理、科学计算

3. asyncio（异步）
   - 适用于高并发 I/O 操作
   - 单线程，无锁开销
   - 需要异步库支持
   - 例如：Web 服务器、爬虫

4. concurrent.futures
   - 统一的高级接口
   - ThreadPoolExecutor: 线程池
   - ProcessPoolExecutor: 进程池
   - 简化任务提交和结果获取

选择建议：
- I/O 密集 + 高并发 → asyncio
- I/O 密集 + 简单场景 → threading
- CPU 密集 → multiprocessing
- 需要简单接口 → concurrent.futures
"""
print(guide)


if __name__ == "__main__":
    print("\n" + "=" * 50)
    print("16_concurrency.py 运行完成！")
    print("=" * 50)
