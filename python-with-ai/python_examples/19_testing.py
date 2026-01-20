#!/usr/bin/env python3
"""
19_testing.py - 单元测试

Python 单元测试：unittest、pytest、mock、测试覆盖率等
"""

import unittest
from unittest import mock
from unittest.mock import Mock, MagicMock, patch
import doctest

# =============================================================================
# 1. 被测试的代码
# =============================================================================

class Calculator:
    """计算器类 - 用于测试示例"""

    def add(self, a, b):
        """加法

        >>> calc = Calculator()
        >>> calc.add(2, 3)
        5
        """
        return a + b

    def subtract(self, a, b):
        """减法"""
        return a - b

    def multiply(self, a, b):
        """乘法"""
        return a * b

    def divide(self, a, b):
        """除法"""
        if b == 0:
            raise ValueError("Cannot divide by zero")
        return a / b


def fetch_data(url):
    """模拟从URL获取数据"""
    import urllib.request
    with urllib.request.urlopen(url) as response:
        return response.read().decode('utf-8')


class UserService:
    """用户服务类 - 用于演示 mock"""

    def __init__(self, database):
        self.database = database

    def get_user(self, user_id):
        """获取用户"""
        return self.database.find_user(user_id)

    def create_user(self, name, email):
        """创建用户"""
        user = {"name": name, "email": email}
        return self.database.save_user(user)


# =============================================================================
# 2. unittest 基础
# =============================================================================

class TestCalculator(unittest.TestCase):
    """Calculator 类的单元测试"""

    @classmethod
    def setUpClass(cls):
        """测试类初始化（只执行一次）"""
        print("\n[setUpClass] 初始化测试类")
        cls.shared_resource = "共享资源"

    @classmethod
    def tearDownClass(cls):
        """测试类清理（只执行一次）"""
        print("[tearDownClass] 清理测试类")

    def setUp(self):
        """每个测试方法前执行"""
        self.calc = Calculator()

    def tearDown(self):
        """每个测试方法后执行"""
        pass

    # 基本断言
    def test_add(self):
        """测试加法"""
        self.assertEqual(self.calc.add(2, 3), 5)
        self.assertEqual(self.calc.add(-1, 1), 0)
        self.assertEqual(self.calc.add(0, 0), 0)

    def test_subtract(self):
        """测试减法"""
        self.assertEqual(self.calc.subtract(5, 3), 2)
        self.assertEqual(self.calc.subtract(3, 5), -2)

    def test_multiply(self):
        """测试乘法"""
        self.assertEqual(self.calc.multiply(3, 4), 12)
        self.assertEqual(self.calc.multiply(-2, 3), -6)

    def test_divide(self):
        """测试除法"""
        self.assertEqual(self.calc.divide(10, 2), 5)
        self.assertAlmostEqual(self.calc.divide(1, 3), 0.333, places=3)

    def test_divide_by_zero(self):
        """测试除以零的异常"""
        with self.assertRaises(ValueError) as context:
            self.calc.divide(10, 0)
        self.assertEqual(str(context.exception), "Cannot divide by zero")

    # 更多断言方法
    def test_assertions_demo(self):
        """演示各种断言方法"""
        # 相等性
        self.assertEqual(1 + 1, 2)
        self.assertNotEqual(1 + 1, 3)

        # 布尔
        self.assertTrue(1 < 2)
        self.assertFalse(1 > 2)

        # None
        self.assertIsNone(None)
        self.assertIsNotNone("value")

        # 身份
        a = [1, 2, 3]
        b = a
        self.assertIs(a, b)
        self.assertIsNot(a, [1, 2, 3])

        # 成员
        self.assertIn(2, [1, 2, 3])
        self.assertNotIn(4, [1, 2, 3])

        # 类型
        self.assertIsInstance([], list)

        # 近似
        self.assertAlmostEqual(0.1 + 0.2, 0.3, places=1)

        # 比较
        self.assertGreater(2, 1)
        self.assertGreaterEqual(2, 2)
        self.assertLess(1, 2)
        self.assertLessEqual(2, 2)

        # 正则
        self.assertRegex("hello world", r"world")

    # 跳过测试
    @unittest.skip("演示跳过测试")
    def test_skipped(self):
        """这个测试会被跳过"""
        pass

    @unittest.skipIf(True, "条件为真时跳过")
    def test_skip_if(self):
        """条件跳过"""
        pass

    @unittest.expectedFailure
    def test_expected_failure(self):
        """预期会失败的测试"""
        self.assertEqual(1, 2)


# =============================================================================
# 3. Mock 测试
# =============================================================================

class TestUserService(unittest.TestCase):
    """演示 mock 测试"""

    def test_get_user_with_mock(self):
        """使用 Mock 对象"""
        # 创建 mock 数据库
        mock_db = Mock()
        mock_db.find_user.return_value = {"id": 1, "name": "Alice"}

        # 注入 mock
        service = UserService(mock_db)
        user = service.get_user(1)

        # 验证返回值
        self.assertEqual(user["name"], "Alice")

        # 验证 mock 被正确调用
        mock_db.find_user.assert_called_once_with(1)

    def test_create_user_with_mock(self):
        """测试创建用户"""
        mock_db = Mock()
        mock_db.save_user.return_value = {"id": 1, "name": "Bob"}

        service = UserService(mock_db)
        result = service.create_user("Bob", "bob@example.com")

        # 验证调用参数
        mock_db.save_user.assert_called_once()
        call_args = mock_db.save_user.call_args[0][0]
        self.assertEqual(call_args["name"], "Bob")
        self.assertEqual(call_args["email"], "bob@example.com")

    def test_with_side_effect(self):
        """演示 side_effect"""
        mock_db = Mock()

        # side_effect 可以是异常
        mock_db.find_user.side_effect = Exception("Database error")

        service = UserService(mock_db)
        with self.assertRaises(Exception):
            service.get_user(1)

        # side_effect 也可以是函数
        mock_db.find_user.side_effect = lambda uid: {"id": uid, "name": f"User{uid}"}
        user = service.get_user(42)
        self.assertEqual(user["id"], 42)

    def test_with_patch(self):
        """使用 patch 装饰器"""
        with patch('urllib.request.urlopen') as mock_urlopen:
            # 配置 mock
            mock_response = Mock()
            mock_response.read.return_value = b'{"data": "test"}'
            mock_response.__enter__ = Mock(return_value=mock_response)
            mock_response.__exit__ = Mock(return_value=False)
            mock_urlopen.return_value = mock_response

            # 调用被测函数
            result = fetch_data("http://example.com")

            # 验证
            self.assertEqual(result, '{"data": "test"}')
            mock_urlopen.assert_called_once_with("http://example.com")


# =============================================================================
# 4. 参数化测试
# =============================================================================

class TestParameterized(unittest.TestCase):
    """参数化测试示例"""

    def test_add_parametrized(self):
        """手动参数化"""
        calc = Calculator()
        test_cases = [
            (1, 2, 3),
            (0, 0, 0),
            (-1, 1, 0),
            (100, 200, 300),
        ]
        for a, b, expected in test_cases:
            with self.subTest(a=a, b=b):
                self.assertEqual(calc.add(a, b), expected)


# =============================================================================
# 5. doctest
# =============================================================================

def factorial(n):
    """
    计算阶乘

    >>> factorial(0)
    1
    >>> factorial(1)
    1
    >>> factorial(5)
    120
    >>> factorial(-1)
    Traceback (most recent call last):
        ...
    ValueError: n must be >= 0
    """
    if n < 0:
        raise ValueError("n must be >= 0")
    if n <= 1:
        return 1
    return n * factorial(n - 1)


# =============================================================================
# 6. pytest 风格（可与 pytest 运行）
# =============================================================================

# pytest 不需要类，可以直接写函数
def test_simple_add():
    """pytest 风格的简单测试"""
    calc = Calculator()
    assert calc.add(2, 3) == 5


def test_simple_divide():
    """pytest 风格的异常测试"""
    calc = Calculator()
    import pytest
    # 如果安装了 pytest，可以使用：
    # with pytest.raises(ValueError):
    #     calc.divide(10, 0)


# pytest 参数化示例（需要 pytest）
# import pytest
# @pytest.mark.parametrize("a,b,expected", [
#     (1, 2, 3),
#     (0, 0, 0),
#     (-1, 1, 0),
# ])
# def test_add_parametrized(a, b, expected):
#     calc = Calculator()
#     assert calc.add(a, b) == expected


# =============================================================================
# 7. 测试覆盖率说明
# =============================================================================

coverage_info = """
=== 测试覆盖率 ===

安装 coverage:
    pip install coverage

运行测试并收集覆盖率:
    coverage run -m pytest
    # 或
    coverage run -m unittest discover

生成报告:
    coverage report              # 终端报告
    coverage html                # HTML 报告

配置文件 (.coveragerc):
    [run]
    source = src
    omit = */tests/*

    [report]
    exclude_lines =
        pragma: no cover
        def __repr__
        raise NotImplementedError
"""

# =============================================================================
# 8. pytest 高级特性说明
# =============================================================================

pytest_info = """
=== pytest 高级特性 ===

1. Fixtures:
   @pytest.fixture
   def calculator():
       return Calculator()

   def test_add(calculator):
       assert calculator.add(2, 3) == 5

2. 参数化:
   @pytest.mark.parametrize("input,expected", [
       (1, 2),
       (2, 4),
   ])
   def test_double(input, expected):
       assert input * 2 == expected

3. 标记:
   @pytest.mark.slow
   def test_slow_operation():
       pass

   # 运行: pytest -m slow

4. 配置 (pytest.ini):
   [pytest]
   testpaths = tests
   python_files = test_*.py
   python_functions = test_*
   addopts = -v --tb=short

5. 插件:
   - pytest-cov: 覆盖率
   - pytest-xdist: 并行测试
   - pytest-mock: mock 增强
   - pytest-asyncio: 异步测试
"""


# =============================================================================
# 9. 运行测试
# =============================================================================

if __name__ == "__main__":
    print("=" * 60)
    print("Python 单元测试示例")
    print("=" * 60)

    # 运行 doctest
    print("\n--- 运行 doctest ---")
    doctest.testmod(verbose=True)

    # 运行 unittest
    print("\n--- 运行 unittest ---")

    # 创建测试套件
    loader = unittest.TestLoader()
    suite = unittest.TestSuite()

    # 添加测试类
    suite.addTests(loader.loadTestsFromTestCase(TestCalculator))
    suite.addTests(loader.loadTestsFromTestCase(TestUserService))
    suite.addTests(loader.loadTestsFromTestCase(TestParameterized))

    # 运行测试
    runner = unittest.TextTestRunner(verbosity=2)
    result = runner.run(suite)

    # 输出说明
    print(coverage_info)
    print(pytest_info)

    print("\n" + "=" * 60)
    print("测试完成！")
    print("=" * 60)
    print("\n使用方法:")
    print("  python 19_testing.py          # 运行所有测试")
    print("  python -m pytest 19_testing.py  # 使用 pytest 运行")
    print("  python -m doctest 19_testing.py # 只运行 doctest")
