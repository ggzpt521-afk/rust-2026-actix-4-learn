// 08_collections.rs - Rust集合类型详解

// Rust标准库提供了多种集合类型，主要分为三类：
// 1. 有序序列：Vec<T>、String、VecDeque<T>、LinkedList<T>
// 2. 映射：HashMap<K, V>、BTreeMap<K, V>
// 3. 集合：HashSet<T>、BTreeSet<T>

use std::collections::{HashMap, HashSet, VecDeque, LinkedList, BTreeMap, BTreeSet};
use std::hash::Hash;

pub fn run_example() {
    println!("=== Rust学习示例 ===\n");
    // 1. Vec<T> - 动态数组
    println!("=== Vec<T> - 动态数组 ===");
    
    // 创建Vec
    let mut v1: Vec<i32> = Vec::new();
    let mut v2 = vec![1, 2, 3, 4, 5];
    let v3 = vec![0; 10]; // 创建包含10个0的Vec
    
    println!("v2: {:?}", v2);
    println!("v3: {:?}", v3);
    
    // 添加元素
    v1.push(6);
    v1.push(7);
    v1.extend(&[8, 9, 10]);
    
    println!("添加元素后v1: {:?}", v1);
    
    // 访问元素
    println!("v2的第一个元素: {}", v2[0]); // 索引访问（越界会panic）
    println!("v2的第一个元素: {:?}", v2.get(0)); // get方法（返回Option<T>）
    println!("v2的第10个元素: {:?}", v2.get(10)); // 安全访问不存在的元素
    
    // 修改元素
    v2[1] = 20;
    
    // 遍历元素
    println!("遍历v2的元素:");
    for i in &v2 {
        println!("{}", i);
    }
    
    // 遍历并修改元素
    println!("遍历并修改v2的元素:");
    for i in &mut v2 {
        *i *= 2;
        println!("{}", i);
    }
    
    // 删除元素
    let removed = v2.pop(); // 删除并返回最后一个元素
    println!("删除的元素: {:?}", removed);
    println!("删除后v2: {:?}", v2);
    
    v2.remove(0); // 删除指定索引的元素
    println!("删除索引0的元素后v2: {:?}", v2);
    
    // Vec的其他常用方法
    println!("v2的长度: {}", v2.len());
    println!("v2是否为空: {}", v2.is_empty());
    v2.clear(); // 清空Vec
    println!("清空后v2: {:?}", v2);
    
    // 2. String - UTF-8字符串
    println!("\n=== String - UTF-8字符串 ===");
    
    // 创建String
    let s1 = String::new();
    let s2 = String::from("Hello");
    let s3 = " world".to_string();
    
    // 字符串拼接
    let mut s4 = s2 + &s3; // s2被移动，不再可用
    s4.push('!'); // 添加单个字符
    s4.push_str(" Rust"); // 添加字符串切片
    
    println!("s4: {}", s4);
    
    // 格式化字符串
    let name = "Alice";
    let age = 30;
    let s5 = format!("My name is {} and I'm {} years old.", name, age);
    println!("s5: {}", s5);
    
    // 访问字符串
    println!("s4的长度: {}", s4.len());
    println!("s4的字节数: {}", s4.as_bytes().len());
    
    // 注意：Rust的String不支持索引访问，因为UTF-8字符可能占用多个字节
    // println!("s4[0]: {}", s4[0]); // 这会报错
    
    // 遍历字符串的字符
    println!("遍历s4的字符:");
    for c in s4.chars() {
        println!("{}", c);
    }
    
    // 遍历字符串的字节
    println!("遍历s4的字节:");
    for b in s4.bytes() {
        println!("{}", b);
    }
    
    // 字符串切片（注意：必须在字符边界上切片）
    let hello = "你好，世界";
    let s6 = &hello[0..6]; // "你好"，因为每个中文字符占用3个字节
    println!("s6: {}", s6);
    
    // 3. HashMap<K, V> - 哈希映射
    println!("\n=== HashMap<K, V> - 哈希映射 ===");
    
    // 创建HashMap
    let mut scores: HashMap<String, i32> = HashMap::new();
    
    // 添加键值对
    scores.insert(String::from("Alice"), 100);
    scores.insert(String::from("Bob"), 85);
    scores.insert(String::from("Charlie"), 90);
    
    // 从元组向量创建HashMap
    let teams = vec![String::from("Blue"), String::from("Yellow")];
    let initial_scores = vec![10, 50];
    let mut team_scores: HashMap<_, _> = teams.into_iter().zip(initial_scores.into_iter()).collect();
    
    println!("scores: {:?}", scores);
    println!("team_scores: {:?}", team_scores);
    
    // 访问HashMap的值
    let alice_score = scores.get(&String::from("Alice"));
    println!("Alice的分数: {:?}", alice_score);
    
    // 遍历HashMap
    println!("遍历scores:");
    for (key, value) in &scores {
        println!("{}: {}", key, value);
    }
    
    // 修改HashMap
    scores.insert(String::from("Alice"), 105); // 覆盖现有值
    println!("修改后Alice的分数: {:?}", scores.get(&String::from("Alice")));
    
    // 只有在键不存在时插入
    scores.entry(String::from("Alice")).or_insert(110); // 不会插入，因为键已存在
    scores.entry(String::from("David")).or_insert(95); // 会插入，因为键不存在
    
    println!("使用entry后scores: {:?}", scores);
    
    // 基于旧值更新新值
    let text = "hello world hello rust hello world";
    let mut word_count = HashMap::new();
    
    for word in text.split_whitespace() {
        let count = word_count.entry(word).or_insert(0);
        *count += 1;
    }
    
    println!("单词计数: {:?}", word_count);
    
    // 4. HashSet<T> - 哈希集合
    println!("\n=== HashSet<T> - 哈希集合 ===");
    
    // 创建HashSet
    let mut set: HashSet<i32> = HashSet::new();
    
    // 添加元素
    set.insert(1);
    set.insert(2);
    set.insert(3);
    set.insert(3); // 重复元素，不会被添加
    
    println!("set: {:?}", set);
    
    // 检查元素是否存在
    println!("set是否包含2: {}", set.contains(&2));
    println!("set是否包含4: {}", set.contains(&4));
    
    // 删除元素
    set.remove(&2);
    println!("删除2后set: {:?}", set);
    
    // 遍历HashSet
    println!("遍历set:");
    for num in &set {
        println!("{}", num);
    }
    
    // 集合操作
    let set1: HashSet<_> = [1, 2, 3, 4, 5].iter().cloned().collect();
    let set2: HashSet<_> = [4, 5, 6, 7, 8].iter().cloned().collect();
    
    println!("set1: {:?}", set1);
    println!("set2: {:?}", set2);
    
    // 交集
    let intersection: HashSet<_> = set1.intersection(&set2).cloned().collect();
    println!("交集: {:?}", intersection);
    
    // 并集
    let union: HashSet<_> = set1.union(&set2).cloned().collect();
    println!("并集: {:?}", union);
    
    // 差集
    let difference: HashSet<_> = set1.difference(&set2).cloned().collect();
    println!("差集 (set1 - set2): {:?}", difference);
    
    // 对称差集
    let symmetric_difference: HashSet<_> = set1.symmetric_difference(&set2).cloned().collect();
    println!("对称差集: {:?}", symmetric_difference);
    
    // 5. VecDeque<T> - 双端队列
    println!("\n=== VecDeque<T> - 双端队列 ===");
    
    // 创建VecDeque
    let mut deque: VecDeque<i32> = VecDeque::new();
    
    // 添加元素
    deque.push_back(1);
    deque.push_back(2);
    deque.push_back(3);
    deque.push_front(0);
    
    println!("deque: {:?}", deque);
    
    // 访问元素
    println!("deque的第一个元素: {:?}", deque.front());
    println!("deque的最后一个元素: {:?}", deque.back());
    
    // 修改元素
    if let Some(first) = deque.front_mut() {
        *first = -1;
    }
    
    println!("修改第一个元素后deque: {:?}", deque);
    
    // 删除元素
    deque.pop_front();
    deque.pop_back();
    
    println!("删除首尾元素后deque: {:?}", deque);
    
    // 6. LinkedList<T> - 双向链表
    println!("\n=== LinkedList<T> - 双向链表 ===");
    
    // 创建LinkedList
    let mut list: LinkedList<i32> = LinkedList::new();
    
    // 添加元素
    list.push_back(1);
    list.push_back(2);
    list.push_front(0);
    
    println!("list: {:?}", list);
    
    // 访问元素
    println!("list的第一个元素: {:?}", list.front());
    println!("list的最后一个元素: {:?}", list.back());
    
    // 删除元素
    list.pop_front();
    list.pop_back();
    
    println!("删除首尾元素后list: {:?}", list);
    
    // 7. BTreeMap<K, V> - 基于B树的有序映射
    println!("\n=== BTreeMap<K, V> - 有序映射 ===");
    
    // 创建BTreeMap
    let mut btree_map: BTreeMap<i32, String> = BTreeMap::new();
    
    // 添加键值对
    btree_map.insert(3, String::from("three"));
    btree_map.insert(1, String::from("one"));
    btree_map.insert(4, String::from("four"));
    btree_map.insert(2, String::from("two"));
    
    // BTreeMap会自动按键排序
    println!("btree_map: {:?}", btree_map);
    
    // 遍历BTreeMap（按键的顺序）
    println!("遍历btree_map:");
    for (key, value) in &btree_map {
        println!("{}: {}", key, value);
    }
    
    // 8. BTreeSet<T> - 基于B树的有序集合
    println!("\n=== BTreeSet<T> - 有序集合 ===");
    
    // 创建BTreeSet
    let mut btree_set: BTreeSet<i32> = BTreeSet::new();
    
    // 添加元素
    btree_set.insert(3);
    btree_set.insert(1);
    btree_set.insert(4);
    btree_set.insert(1); // 重复元素，不会被添加
    btree_set.insert(2);
    
    // BTreeSet会自动按元素排序
    println!("btree_set: {:?}", btree_set);
    
    // 遍历BTreeSet（按元素的顺序）
    println!("遍历btree_set:");
    for num in &btree_set {
        println!("{}", num);
    }
    
    // 9. 集合类型的选择指南
    println!("\n=== 集合类型的选择指南 ===");
    println!("- 动态数组: Vec<T> - 用于需要快速随机访问和追加元素的场景");
    println!("- 字符串: String - 用于处理UTF-8文本");
    println!("- 双端队列: VecDeque<T> - 用于需要在两端快速添加/删除元素的场景");
    println!("- 链表: LinkedList<T> - 用于需要在任意位置快速插入/删除元素的场景");
    println!("- 无序映射: HashMap<K, V> - 用于需要快速查找键值对的场景");
    println!("- 有序映射: BTreeMap<K, V> - 用于需要按键排序的映射场景");
    println!("- 无序集合: HashSet<T> - 用于需要快速判断元素是否存在的场景");
    println!("- 有序集合: BTreeSet<T> - 用于需要按元素排序的集合场景");
}

// 10. 为自定义类型实现Hash和Eq trait以用于HashMap和HashSet
#[derive(Debug, PartialEq, Eq, Hash)]
struct Person {
    name: String,
    age: u32,
}

impl Person {
    fn new(name: &str, age: u32) -> Self {
        Self {
            name: name.to_string(),
            age,
        }
    }
}

fn custom_type_in_hashmap() {
    println!("\n=== 自定义类型在HashMap中的使用 ===");
    
    let mut person_map = HashMap::new();
    
    let alice = Person::new("Alice", 30);
    let bob = Person::new("Bob", 25);
    
    person_map.insert(alice, "Engineer");
    person_map.insert(bob, "Designer");
    
    // 注意：由于alice和bob已被移动到HashMap中，无法再直接使用
    
    // 可以使用相同字段创建新实例来查找
    let alice_search = Person::new("Alice", 30);
    println!("Alice的职业: {:?}", person_map.get(&alice_search));
}

// 11. 集合的性能比较
// Vec<T>:
// - 随机访问: O(1)
// - 末尾添加/删除: O(1) 均摊
// - 中间插入/删除: O(n)

// String:
// - 末尾添加/删除: O(1) 均摊
// - 中间插入/删除: O(n)

// HashMap<K, V>:
// - 插入: O(1) 平均，O(n) 最坏
// - 查找: O(1) 平均，O(n) 最坏
// - 删除: O(1) 平均，O(n) 最坏

// BTreeMap<K, V>:
// - 插入: O(log n)
// - 查找: O(log n)
// - 删除: O(log n)

// HashSet<T> 和 BTreeSet<T> 的性能与对应的映射类型类似

// 12. 集合的所有权
// 集合类型会获取它们所包含元素的所有权
// 如果需要在集合中存储引用，必须使用生命周期参数

fn ownership_example() {
    let s1 = String::from("hello");
    let s2 = String::from("world");
    
    let mut vec = Vec::new();
    vec.push(s1); // s1的所有权被转移到vec
    vec.push(s2); // s2的所有权被转移到vec
    
    // println!("s1: {}", s1); // 这会报错，因为s1的所有权已被转移
    
    // 使用引用
    let s3 = String::from("rust");
    let s4 = String::from("programming");
    
    let mut vec_ref = Vec::new();
    vec_ref.push(&s3);
    vec_ref.push(&s4);
    
    println!("vec_ref: {:?}", vec_ref);
    println!("s3: {}", s3); // 可以正常访问，因为只是借用
}
// 用于单独运行本文件的main函数
fn main() {
    run_example();
}
