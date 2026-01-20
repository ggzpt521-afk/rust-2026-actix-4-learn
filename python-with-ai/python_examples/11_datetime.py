#!/usr/bin/env python3
"""
11_datetime.py - 日期时间处理

Python datetime 模块、时区处理、日期计算、格式化等
"""

import datetime
from datetime import date, time, datetime as dt, timedelta
from zoneinfo import ZoneInfo  # Python 3.9+
import calendar

# =============================================================================
# 1. date 对象
# =============================================================================

print("=== date 对象 ===")

# 创建日期
today = date.today()
specific_date = date(2024, 6, 15)
from_string = date.fromisoformat("2024-12-25")

print(f"今天: {today}")
print(f"指定日期: {specific_date}")
print(f"从字符串: {from_string}")

# 日期属性
print(f"\n日期属性:")
print(f"  year: {today.year}")
print(f"  month: {today.month}")
print(f"  day: {today.day}")
print(f"  weekday(): {today.weekday()}")  # 0=周一
print(f"  isoweekday(): {today.isoweekday()}")  # 1=周一
print(f"  isocalendar(): {today.isocalendar()}")

# 日期格式化
print(f"\n日期格式化:")
print(f"  isoformat(): {today.isoformat()}")
print(f"  strftime: {today.strftime('%Y年%m月%d日')}")
print(f"  strftime: {today.strftime('%A, %B %d, %Y')}")

# =============================================================================
# 2. time 对象
# =============================================================================

print("\n=== time 对象 ===")

# 创建时间
t1 = time(14, 30, 45)
t2 = time(14, 30, 45, 123456)  # 带微秒
t3 = time.fromisoformat("15:30:00")

print(f"t1: {t1}")
print(f"t2 (带微秒): {t2}")
print(f"t3: {t3}")

# 时间属性
print(f"\n时间属性:")
print(f"  hour: {t1.hour}")
print(f"  minute: {t1.minute}")
print(f"  second: {t1.second}")
print(f"  microsecond: {t2.microsecond}")

# 时间格式化
print(f"\n时间格式化:")
print(f"  isoformat(): {t1.isoformat()}")
print(f"  strftime: {t1.strftime('%H:%M:%S')}")
print(f"  strftime: {t1.strftime('%I:%M %p')}")

# =============================================================================
# 3. datetime 对象
# =============================================================================

print("\n=== datetime 对象 ===")

# 创建 datetime
now = dt.now()
utc_now = dt.utcnow()  # 已弃用，建议使用 datetime.now(timezone.utc)
specific = dt(2024, 6, 15, 14, 30, 0)
combined = dt.combine(date.today(), time(10, 30))

print(f"现在: {now}")
print(f"UTC: {utc_now}")
print(f"指定时间: {specific}")
print(f"组合: {combined}")

# 从时间戳创建
timestamp = 1700000000
from_ts = dt.fromtimestamp(timestamp)
print(f"从时间戳 {timestamp}: {from_ts}")

# 转换为时间戳
ts = now.timestamp()
print(f"当前时间戳: {ts}")

# 解析字符串
parsed = dt.strptime("2024-06-15 14:30:00", "%Y-%m-%d %H:%M:%S")
print(f"解析字符串: {parsed}")

# 格式化
print(f"\ndatetime 格式化:")
print(f"  ISO: {now.isoformat()}")
print(f"  自定义: {now.strftime('%Y-%m-%d %H:%M:%S')}")
print(f"  中文: {now.strftime('%Y年%m月%d日 %H时%M分%S秒')}")

# =============================================================================
# 4. timedelta - 时间差
# =============================================================================

print("\n=== timedelta 时间差 ===")

# 创建 timedelta
delta1 = timedelta(days=7)
delta2 = timedelta(hours=2, minutes=30)
delta3 = timedelta(weeks=2, days=3, hours=5)

print(f"delta1: {delta1}")
print(f"delta2: {delta2}")
print(f"delta3: {delta3}")

# timedelta 属性
print(f"\ntimedelta 属性:")
print(f"  days: {delta3.days}")
print(f"  seconds: {delta3.seconds}")
print(f"  total_seconds(): {delta3.total_seconds()}")

# 日期计算
today = date.today()
print(f"\n日期计算:")
print(f"  今天: {today}")
print(f"  一周后: {today + timedelta(days=7)}")
print(f"  30天前: {today - timedelta(days=30)}")

# datetime 计算
now = dt.now()
print(f"\ndatetime 计算:")
print(f"  现在: {now}")
print(f"  2小时后: {now + timedelta(hours=2)}")
print(f"  昨天这个时间: {now - timedelta(days=1)}")

# 计算时间差
d1 = dt(2024, 1, 1)
d2 = dt(2024, 12, 31)
diff = d2 - d1
print(f"\n{d1} 到 {d2} 相差:")
print(f"  {diff.days} 天")
print(f"  {diff.total_seconds()} 秒")

# =============================================================================
# 5. 时区处理
# =============================================================================

print("\n=== 时区处理 ===")

# 使用 zoneinfo（Python 3.9+）
utc = ZoneInfo("UTC")
beijing = ZoneInfo("Asia/Shanghai")
new_york = ZoneInfo("America/New_York")
tokyo = ZoneInfo("Asia/Tokyo")

# 创建带时区的 datetime
now_utc = dt.now(utc)
now_beijing = dt.now(beijing)
now_ny = dt.now(new_york)

print(f"UTC: {now_utc}")
print(f"北京: {now_beijing}")
print(f"纽约: {now_ny}")

# 时区转换
utc_time = dt(2024, 6, 15, 12, 0, 0, tzinfo=utc)
beijing_time = utc_time.astimezone(beijing)
ny_time = utc_time.astimezone(new_york)
tokyo_time = utc_time.astimezone(tokyo)

print(f"\n时区转换 (UTC 12:00):")
print(f"  北京: {beijing_time.strftime('%Y-%m-%d %H:%M %Z')}")
print(f"  纽约: {ny_time.strftime('%Y-%m-%d %H:%M %Z')}")
print(f"  东京: {tokyo_time.strftime('%Y-%m-%d %H:%M %Z')}")

# naive datetime 转 aware datetime
naive = dt(2024, 6, 15, 12, 0, 0)
aware = naive.replace(tzinfo=beijing)
print(f"\nnaive: {naive}")
print(f"aware: {aware}")

# =============================================================================
# 6. calendar 模块
# =============================================================================

print("\n=== calendar 模块 ===")

# 月历
print("2024年6月:")
print(calendar.month(2024, 6))

# 判断闰年
print(f"2024 是闰年: {calendar.isleap(2024)}")
print(f"2023 是闰年: {calendar.isleap(2023)}")

# 某月的天数范围
print(f"2024年2月: {calendar.monthrange(2024, 2)}")  # (周几开始, 天数)

# 某年某月的所有周
print("2024年6月的周:")
cal = calendar.Calendar()
for week in cal.monthdayscalendar(2024, 6):
    print(f"  {week}")

# =============================================================================
# 7. 常用日期操作
# =============================================================================

print("\n=== 常用日期操作 ===")


def get_month_start_end(year, month):
    """获取某月的第一天和最后一天"""
    first_day = date(year, month, 1)
    _, last_day_num = calendar.monthrange(year, month)
    last_day = date(year, month, last_day_num)
    return first_day, last_day


def get_week_range(d):
    """获取某天所在周的起止日期（周一到周日）"""
    start = d - timedelta(days=d.weekday())
    end = start + timedelta(days=6)
    return start, end


def get_age(birth_date):
    """计算年龄"""
    today = date.today()
    age = today.year - birth_date.year
    if (today.month, today.day) < (birth_date.month, birth_date.day):
        age -= 1
    return age


# 测试
start, end = get_month_start_end(2024, 2)
print(f"2024年2月: {start} ~ {end}")

week_start, week_end = get_week_range(date.today())
print(f"本周: {week_start} ~ {week_end}")

birth = date(1990, 6, 15)
print(f"出生日期 {birth}，年龄: {get_age(birth)}")

# =============================================================================
# 8. 日期格式化代码
# =============================================================================

print("\n=== 日期格式化代码 ===")

format_codes = """
常用格式化代码：
%Y - 4位年份 (2024)
%y - 2位年份 (24)
%m - 月份 (01-12)
%d - 日期 (01-31)
%H - 小时24制 (00-23)
%I - 小时12制 (01-12)
%M - 分钟 (00-59)
%S - 秒 (00-59)
%f - 微秒 (000000-999999)
%p - AM/PM
%A - 星期全名 (Monday)
%a - 星期缩写 (Mon)
%B - 月份全名 (January)
%b - 月份缩写 (Jan)
%j - 一年中的第几天 (001-366)
%U - 一年中的第几周（周日开始）
%W - 一年中的第几周（周一开始）
%w - 星期几 (0=周日)
%Z - 时区名
%z - UTC 偏移 (+0800)
"""
print(format_codes)

# 示例
now = dt.now()
formats = [
    "%Y-%m-%d",
    "%Y/%m/%d %H:%M:%S",
    "%Y年%m月%d日",
    "%A, %B %d, %Y",
    "%I:%M %p",
    "%Y-%m-%d %H:%M:%S.%f",
]

print("格式化示例:")
for fmt in formats:
    print(f"  '{fmt}': {now.strftime(fmt)}")

# =============================================================================
# 9. ISO 8601 格式
# =============================================================================

print("\n=== ISO 8601 格式 ===")

now = dt.now(ZoneInfo("Asia/Shanghai"))

print(f"isoformat(): {now.isoformat()}")
print(f"isoformat(sep=' '): {now.isoformat(sep=' ')}")
print(f"isoformat(timespec='seconds'): {now.isoformat(timespec='seconds')}")
print(f"isoformat(timespec='minutes'): {now.isoformat(timespec='minutes')}")

# 解析 ISO 格式
iso_string = "2024-06-15T14:30:00+08:00"
parsed = dt.fromisoformat(iso_string)
print(f"\n解析 '{iso_string}':")
print(f"  结果: {parsed}")
print(f"  时区: {parsed.tzinfo}")

# =============================================================================
# 10. 性能提示
# =============================================================================

print("\n=== 性能提示 ===")

performance_tips = """
1. 避免在循环中重复调用 datetime.now()
   # 不好
   for item in items:
       item.timestamp = datetime.now()

   # 好
   now = datetime.now()
   for item in items:
       item.timestamp = now

2. 字符串解析使用 strptime 比较慢
   考虑使用第三方库如 dateutil.parser

3. 时区转换可能较慢
   如果需要大量转换，考虑使用批量处理

4. 使用 timestamp() 进行日期比较比直接比较 datetime 对象更快
"""
print(performance_tips)


if __name__ == "__main__":
    print("\n" + "=" * 50)
    print("11_datetime.py 运行完成！")
    print("=" * 50)
