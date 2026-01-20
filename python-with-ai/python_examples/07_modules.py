#!/usr/bin/env python3
"""
07_modules.py - 模块与包管理

Python 模块导入、包结构、标准库使用、虚拟环境概念
"""

# =============================================================================
# 1. 模块导入方式
# =============================================================================

print("=== 模块导入方式 ===")

# 导入整个模块
import math

print(f"import math: math.pi = {math.pi}")
print(f"import math: math.sqrt(16) = {math.sqrt(16)}")

# 导入特定对象
from os import getcwd, listdir

print(f"\nfrom os import getcwd: {getcwd()}")

# 导入并重命名
import datetime as dt

now = dt.datetime.now()
print(f"import datetime as dt: {now}")

# 从模块导入并重命名
from collections import Counter as Cnt

print(f"from collections import Counter as Cnt: {Cnt('hello')}")

# 导入模块中的所有公开对象（不推荐）
# from math import *  # 会污染命名空间

# =============================================================================
# 2. 模块搜索路径
# =============================================================================

print("\n=== 模块搜索路径 ===")

import sys

print("sys.path 搜索顺序:")
for i, path in enumerate(sys.path[:5]):
    print(f"  {i}: {path}")
print("  ...")

# 动态添加搜索路径
# sys.path.append('/custom/module/path')

# =============================================================================
# 3. 模块属性
# =============================================================================

print("\n=== 模块属性 ===")

print(f"当前模块名: {__name__}")
print(f"当前文件: {__file__}")

# 模块信息
print(f"\nmath 模块信息:")
print(f"  __name__: {math.__name__}")
print(f"  __doc__: {math.__doc__[:50]}...")

# dir() 查看模块内容
print(f"\nmath 模块包含 {len(dir(math))} 个对象")
print(f"部分内容: {[x for x in dir(math) if not x.startswith('_')][:10]}")

# =============================================================================
# 4. __name__ 和 __main__
# =============================================================================

print("\n=== __name__ 和 __main__ ===")

# 当前模块作为主程序运行时，__name__ == '__main__'
# 当前模块被导入时，__name__ == 模块名

print(f"__name__ = {__name__}")

# 常见用法（在文件末尾）
# if __name__ == "__main__":
#     # 只在直接运行时执行，被导入时不执行
#     main()

# =============================================================================
# 5. 包（Package）
# =============================================================================

print("\n=== 包结构说明 ===")

package_structure = """
典型的包结构：

mypackage/
    __init__.py          # 包初始化文件（可以为空）
    module1.py           # 模块1
    module2.py           # 模块2
    subpackage/          # 子包
        __init__.py
        module3.py

导入方式：
    import mypackage
    import mypackage.module1
    from mypackage import module1
    from mypackage.module1 import some_function
    from mypackage.subpackage import module3
"""
print(package_structure)

# =============================================================================
# 6. __init__.py 示例
# =============================================================================

print("=== __init__.py 示例 ===")

init_example = """
# mypackage/__init__.py 示例内容：

# 定义包级别的变量
__version__ = '1.0.0'
__author__ = 'Your Name'

# 控制 from package import * 的行为
__all__ = ['module1', 'module2', 'some_function']

# 在包导入时执行的代码
print(f"Initializing mypackage {__version__}")

# 提升子模块的对象到包级别
from .module1 import some_function
from .module2 import SomeClass

# 现在可以直接: from mypackage import some_function
"""
print(init_example)

# =============================================================================
# 7. 相对导入与绝对导入
# =============================================================================

print("=== 相对导入与绝对导入 ===")

import_examples = """
在包内部的模块中：

# 绝对导入（推荐）
from mypackage.subpackage import module3
from mypackage.module1 import func1

# 相对导入
from . import module1           # 同级目录
from .module1 import func1      # 同级目录的模块
from .. import module2          # 上级目录
from ..subpackage import module3  # 上级的其他子包

注意：相对导入只能在包内使用，不能在主程序中使用
"""
print(import_examples)

# =============================================================================
# 8. 常用标准库模块
# =============================================================================

print("\n=== 常用标准库模块 ===")

# os 模块
import os

print("os 模块:")
print(f"  当前目录: {os.getcwd()}")
print(f"  环境变量 HOME: {os.environ.get('HOME', 'N/A')}")
print(f"  目录分隔符: {os.sep}")
print(f"  路径分隔符: {os.pathsep}")

# sys 模块
print("\nsys 模块:")
print(f"  Python 版本: {sys.version}")
print(f"  平台: {sys.platform}")
print(f"  命令行参数: {sys.argv}")

# json 模块
import json

print("\njson 模块:")
data = {"name": "Alice", "age": 25, "skills": ["Python", "SQL"]}
json_str = json.dumps(data, indent=2, ensure_ascii=False)
print(f"  序列化:\n{json_str}")

parsed = json.loads(json_str)
print(f"  反序列化: {parsed}")

# pathlib 模块（推荐用于路径操作）
from pathlib import Path

print("\npathlib 模块:")
current = Path.cwd()
print(f"  当前目录: {current}")
print(f"  home 目录: {Path.home()}")

example_path = Path("/usr/local/bin/python")
print(f"  路径名: {example_path.name}")
print(f"  父目录: {example_path.parent}")
print(f"  后缀: {example_path.suffix}")
print(f"  各部分: {example_path.parts}")

# 路径操作
new_path = current / "subdir" / "file.txt"
print(f"  路径拼接: {new_path}")

# =============================================================================
# 9. 虚拟环境
# =============================================================================

print("\n=== 虚拟环境 ===")

venv_info = """
创建和使用虚拟环境：

# 创建虚拟环境
python -m venv myenv

# 激活虚拟环境
# Windows:
myenv\\Scripts\\activate
# Unix/macOS:
source myenv/bin/activate

# 退出虚拟环境
deactivate

# 安装包
pip install package_name

# 导出依赖
pip freeze > requirements.txt

# 安装依赖
pip install -r requirements.txt
"""
print(venv_info)

# =============================================================================
# 10. 包管理工具
# =============================================================================

print("=== 包管理工具 ===")

pip_info = """
pip 常用命令：

pip install package              # 安装包
pip install package==1.0.0       # 安装特定版本
pip install package>=1.0.0       # 安装最低版本
pip install -U package           # 升级包
pip uninstall package            # 卸载包
pip list                         # 列出已安装包
pip show package                 # 显示包信息
pip search keyword               # 搜索包（已废弃）
pip freeze                       # 输出已安装包列表

其他工具：
- pipenv: pip + virtualenv 的组合
- poetry: 现代化的依赖管理工具
- conda: Anaconda 的包管理器
"""
print(pip_info)

# =============================================================================
# 11. 模块的重新加载
# =============================================================================

print("=== 模块重新加载 ===")

from importlib import reload

# 重新加载模块（开发时调试用）
# reload(some_module)

print("使用 importlib.reload(module) 重新加载已导入的模块")
print("注意：这主要用于开发调试，生产环境中通常不需要")

# =============================================================================
# 12. 检查模块是否存在
# =============================================================================

print("\n=== 检查模块是否存在 ===")

from importlib.util import find_spec


def module_exists(module_name):
    """检查模块是否可导入"""
    return find_spec(module_name) is not None


modules_to_check = ["os", "numpy", "nonexistent_module"]
for mod in modules_to_check:
    exists = "✓" if module_exists(mod) else "✗"
    print(f"  {mod}: {exists}")

# =============================================================================
# 13. 延迟导入（Lazy Import）
# =============================================================================

print("\n=== 延迟导入 ===")


def process_data(data):
    """只在需要时才导入模块"""
    # 延迟导入，只在函数被调用时才加载
    import json

    return json.dumps(data)


print("延迟导入示例：模块只在函数调用时才加载")
print(f"结果: {process_data({'key': 'value'})}")

# =============================================================================
# 14. 创建可执行包
# =============================================================================

print("\n=== 可执行包 ===")

executable_package = """
创建可执行包（python -m mypackage）：

mypackage/
    __init__.py
    __main__.py    # 关键文件！

# __main__.py 示例内容：
def main():
    print("Package is running!")

if __name__ == "__main__":
    main()

# 运行方式：
python -m mypackage
"""
print(executable_package)


if __name__ == "__main__":
    print("\n" + "=" * 50)
    print("07_modules.py 运行完成！")
    print("=" * 50)
