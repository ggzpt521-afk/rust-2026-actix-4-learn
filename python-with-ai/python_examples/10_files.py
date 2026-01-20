#!/usr/bin/env python3
"""
10_files.py - 文件与操作系统交互

Python 文件读写、目录操作、pathlib、os 模块等
"""

import os
import shutil
from pathlib import Path
import tempfile
import json

# =============================================================================
# 1. 基本文件操作
# =============================================================================

print("=== 基本文件操作 ===")

# 创建临时目录用于演示
temp_dir = Path(tempfile.mkdtemp(prefix="python_demo_"))
print(f"临时目录: {temp_dir}")

# 写入文件
file_path = temp_dir / "sample.txt"
with open(file_path, "w", encoding="utf-8") as f:
    f.write("Hello, Python!\n")
    f.write("这是第二行\n")
    f.write("This is line 3\n")

print(f"文件已写入: {file_path}")

# 读取整个文件
with open(file_path, "r", encoding="utf-8") as f:
    content = f.read()
print(f"\n读取整个文件:\n{content}")

# 读取所有行
with open(file_path, "r", encoding="utf-8") as f:
    lines = f.readlines()
print(f"readlines(): {lines}")

# 逐行读取（内存友好）
print("\n逐行读取:")
with open(file_path, "r", encoding="utf-8") as f:
    for line_num, line in enumerate(f, 1):
        print(f"  {line_num}: {line.strip()}")

# =============================================================================
# 2. 文件打开模式
# =============================================================================

print("\n=== 文件打开模式 ===")

modes_info = """
常用文件打开模式：
'r'  - 读取（默认）
'w'  - 写入（覆盖）
'x'  - 排他创建（文件存在则失败）
'a'  - 追加
'b'  - 二进制模式
't'  - 文本模式（默认）
'+'  - 读写模式

组合示例：
'rb' - 二进制读取
'wb' - 二进制写入
'r+' - 读写（文件必须存在）
'w+' - 读写（覆盖）
'a+' - 读写（追加）
"""
print(modes_info)

# 追加模式
with open(file_path, "a", encoding="utf-8") as f:
    f.write("追加的内容\n")

# 二进制模式
binary_file = temp_dir / "binary.bin"
with open(binary_file, "wb") as f:
    f.write(b"\x00\x01\x02\x03\x04")

with open(binary_file, "rb") as f:
    data = f.read()
print(f"二进制数据: {data}")

# =============================================================================
# 3. 文件对象方法
# =============================================================================

print("\n=== 文件对象方法 ===")

with open(file_path, "r+", encoding="utf-8") as f:
    # tell() - 获取当前位置
    print(f"初始位置: {f.tell()}")

    # read(n) - 读取 n 个字符
    chunk = f.read(10)
    print(f"读取 10 字符: {repr(chunk)}")
    print(f"当前位置: {f.tell()}")

    # seek(offset, whence) - 移动位置
    f.seek(0)  # 回到开头
    print(f"seek(0) 后位置: {f.tell()}")

    # readline() - 读取一行
    line = f.readline()
    print(f"readline(): {repr(line)}")

    # 获取文件信息
    print(f"文件名: {f.name}")
    print(f"模式: {f.mode}")
    print(f"是否关闭: {f.closed}")

# =============================================================================
# 4. pathlib 模块（推荐）
# =============================================================================

print("\n=== pathlib 模块 ===")

# 创建 Path 对象
current = Path.cwd()
home = Path.home()
file_p = Path("/usr/local/bin/python")

print(f"当前目录: {current}")
print(f"home 目录: {home}")

# 路径操作
path = Path("/Users/name/documents/file.txt")
print(f"\n路径分析: {path}")
print(f"  name: {path.name}")          # file.txt
print(f"  stem: {path.stem}")          # file
print(f"  suffix: {path.suffix}")      # .txt
print(f"  parent: {path.parent}")      # /Users/name/documents
print(f"  parts: {path.parts}")        # 各部分

# 路径拼接
new_path = temp_dir / "subdir" / "file.txt"
print(f"\n路径拼接: {new_path}")

# 路径修改
path = Path("document.txt")
print(f"with_suffix('.pdf'): {path.with_suffix('.pdf')}")
print(f"with_stem('report'): {path.with_stem('report')}")

# 创建目录
subdir = temp_dir / "subdir"
subdir.mkdir(exist_ok=True)
print(f"\n创建目录: {subdir}")

# 文件操作
test_file = subdir / "test.txt"
test_file.write_text("Hello from pathlib!", encoding="utf-8")
print(f"write_text: {test_file}")

content = test_file.read_text(encoding="utf-8")
print(f"read_text: {content}")

# 文件信息
print(f"\n文件信息:")
print(f"  exists: {test_file.exists()}")
print(f"  is_file: {test_file.is_file()}")
print(f"  is_dir: {test_file.is_dir()}")
print(f"  stat: {test_file.stat().st_size} bytes")

# 遍历目录
print(f"\n遍历目录 {temp_dir}:")
for item in temp_dir.iterdir():
    item_type = "目录" if item.is_dir() else "文件"
    print(f"  [{item_type}] {item.name}")

# glob 模式匹配
print(f"\nglob('*.txt'):")
for txt_file in temp_dir.glob("**/*.txt"):
    print(f"  {txt_file}")

# =============================================================================
# 5. os 模块
# =============================================================================

print("\n=== os 模块 ===")

# 环境变量
print(f"HOME: {os.environ.get('HOME', 'N/A')}")
print(f"PATH 条目数: {len(os.environ.get('PATH', '').split(os.pathsep))}")

# 系统信息
print(f"\n系统信息:")
print(f"  os.name: {os.name}")
print(f"  os.sep: {repr(os.sep)}")
print(f"  os.pathsep: {repr(os.pathsep)}")
print(f"  os.linesep: {repr(os.linesep)}")

# os.path（传统方式）
import os.path as osp

path = "/usr/local/bin/python"
print(f"\nos.path 操作:")
print(f"  basename: {osp.basename(path)}")
print(f"  dirname: {osp.dirname(path)}")
print(f"  split: {osp.split(path)}")
print(f"  splitext: {osp.splitext('file.txt')}")
print(f"  join: {osp.join('a', 'b', 'c')}")

# =============================================================================
# 6. 文件和目录操作
# =============================================================================

print("\n=== 文件和目录操作 ===")

# 复制文件
src = test_file
dst = temp_dir / "copied.txt"
shutil.copy2(src, dst)  # copy2 保留元数据
print(f"复制文件: {src.name} -> {dst.name}")

# 移动/重命名
new_name = temp_dir / "renamed.txt"
dst.rename(new_name)
print(f"重命名: copied.txt -> renamed.txt")

# 删除文件
new_name.unlink()
print(f"删除文件: renamed.txt")

# 创建多级目录
deep_dir = temp_dir / "a" / "b" / "c"
deep_dir.mkdir(parents=True, exist_ok=True)
print(f"创建多级目录: {deep_dir}")

# 删除目录树
shutil.rmtree(temp_dir / "a")
print(f"删除目录树: a/b/c")

# =============================================================================
# 7. 临时文件
# =============================================================================

print("\n=== 临时文件 ===")

# 临时文件（自动删除）
with tempfile.NamedTemporaryFile(mode="w", suffix=".txt", delete=False) as tf:
    tf.write("临时内容")
    temp_path = tf.name
    print(f"临时文件: {temp_path}")

# 临时目录
with tempfile.TemporaryDirectory() as td:
    print(f"临时目录: {td}")
    temp_file = Path(td) / "temp.txt"
    temp_file.write_text("test")
    print(f"  创建文件: {temp_file.name}")
# 退出 with 块后目录自动删除

# =============================================================================
# 8. JSON 文件
# =============================================================================

print("\n=== JSON 文件 ===")

data = {
    "name": "Alice",
    "age": 25,
    "skills": ["Python", "SQL", "Machine Learning"],
    "active": True,
    "score": 95.5
}

json_file = temp_dir / "data.json"

# 写入 JSON
with open(json_file, "w", encoding="utf-8") as f:
    json.dump(data, f, indent=2, ensure_ascii=False)
print(f"JSON 写入: {json_file}")

# 读取 JSON
with open(json_file, "r", encoding="utf-8") as f:
    loaded = json.load(f)
print(f"JSON 读取: {loaded}")

# 直接使用 pathlib（Python 3.9+）
# json_file.write_text(json.dumps(data))
# loaded = json.loads(json_file.read_text())

# =============================================================================
# 9. CSV 文件
# =============================================================================

print("\n=== CSV 文件 ===")

import csv

csv_file = temp_dir / "data.csv"

# 写入 CSV
rows = [
    ["Name", "Age", "City"],
    ["Alice", 25, "Beijing"],
    ["Bob", 30, "Shanghai"],
    ["Charlie", 35, "Shenzhen"]
]

with open(csv_file, "w", newline="", encoding="utf-8") as f:
    writer = csv.writer(f)
    writer.writerows(rows)
print(f"CSV 写入: {csv_file}")

# 读取 CSV
with open(csv_file, "r", encoding="utf-8") as f:
    reader = csv.reader(f)
    for row in reader:
        print(f"  {row}")

# 使用 DictReader/DictWriter
print("\nDictReader:")
with open(csv_file, "r", encoding="utf-8") as f:
    reader = csv.DictReader(f)
    for row in reader:
        print(f"  {row}")

# =============================================================================
# 10. 文件锁定（跨平台）
# =============================================================================

print("\n=== 文件锁定 ===")

# 简单的文件锁定示例
import fcntl


def write_with_lock(filepath, content):
    """带文件锁的写入"""
    with open(filepath, "w") as f:
        try:
            fcntl.flock(f.fileno(), fcntl.LOCK_EX)  # 排他锁
            f.write(content)
            print(f"  带锁写入完成")
        finally:
            fcntl.flock(f.fileno(), fcntl.LOCK_UN)  # 释放锁


lock_file = temp_dir / "locked.txt"
write_with_lock(lock_file, "Locked content")

# =============================================================================
# 11. 监控文件变化
# =============================================================================

print("\n=== 监控文件变化 ===")

monitor_info = """
使用 watchdog 库监控文件变化（需要安装）：

from watchdog.observers import Observer
from watchdog.events import FileSystemEventHandler

class MyHandler(FileSystemEventHandler):
    def on_modified(self, event):
        print(f"Modified: {event.src_path}")

    def on_created(self, event):
        print(f"Created: {event.src_path}")

observer = Observer()
observer.schedule(MyHandler(), path=".", recursive=True)
observer.start()
"""
print(monitor_info)

# =============================================================================
# 12. 清理临时文件
# =============================================================================

print("\n=== 清理临时文件 ===")

# 删除整个临时目录
shutil.rmtree(temp_dir)
print(f"已清理临时目录: {temp_dir}")

# 删除之前创建的临时文件
if os.path.exists(temp_path):
    os.remove(temp_path)
    print(f"已删除临时文件: {temp_path}")


if __name__ == "__main__":
    print("\n" + "=" * 50)
    print("10_files.py 运行完成！")
    print("=" * 50)
