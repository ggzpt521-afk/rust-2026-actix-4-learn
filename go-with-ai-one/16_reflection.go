// ============================================================================
// 16_reflection.go - 反射原理与实践
// ============================================================================
// 运行: go run 16_reflection.go
//
// 【本文件学习目标】
// 1. 理解反射的核心概念：运行时类型检查和操作
// 2. 掌握 reflect.Type 和 reflect.Value 的使用
// 3. 深入理解反射三大定律
// 4. 了解反射的底层原理（interface{} 内存结构）
// 5. 学会结构体反射：字段遍历、Tag 解析、动态修改
// 6. 掌握方法反射：获取方法、动态调用
// 7. 理解反射的性能代价和优化技巧
//
// 【反射的核心概念】
// - 反射是程序在运行时检查和操作自身结构的能力
// - Go 的反射基于 interface{} 的类型信息
// - reflect.Type 描述类型，reflect.Value 描述值
//
// 【反射 vs 普通代码】
// | 特性       | 普通代码           | 反射代码              |
// |------------|--------------------|-----------------------|
// | 类型检查   | 编译时             | 运行时                |
// | 性能       | 快                 | 慢 100-400 倍         |
// | 类型安全   | 编译器保证         | 需要自己处理 panic    |
// | 适用场景   | 已知类型           | 动态/未知类型         |
// ============================================================================

package main

import (
	"fmt"
	"reflect"
	"strings"
)

// ============================================================================
// 【基础类型定义】
// ============================================================================

// ReflectInt 是自定义类型，用于演示 Type vs Kind
type ReflectInt int

// ReflectUser 结构体，用于演示结构体反射
type ReflectUser struct {
	Name    string `json:"name" validate:"required"`
	Age     int    `json:"age" validate:"min=0,max=150"`
	Email   string `json:"email,omitempty"`
	private string // 私有字段，反射无法访问其值
}

// ReflectCalc 用于演示方法反射
type ReflectCalc struct {
	Name string
}

func (c ReflectCalc) Add(a, b int) int  { return a + b }
func (c ReflectCalc) Sub(a, b int) int  { return a - b }
func (c *ReflectCalc) Mul(a, b int) int { return a * b } // 指针接收者

// ============================================================================
// 【第一部分：反射基础 - TypeOf 和 ValueOf】
// ============================================================================

func demoBasicReflection() {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("【第一部分：反射基础 - TypeOf 和 ValueOf】")
	fmt.Println(strings.Repeat("=", 70))

	// -------------------------------------------------------------------------
	// 1. 基本类型的反射
	// -------------------------------------------------------------------------
	fmt.Println("\n--- 1. 基本类型的反射 ---")

	var x float64 = 3.14

	// reflect.TypeOf() 获取类型信息
	t := reflect.TypeOf(x)
	fmt.Printf("TypeOf(x)  = %v\n", t)        // float64
	fmt.Printf("Type.Name  = %v\n", t.Name()) // float64
	fmt.Printf("Type.Kind  = %v\n", t.Kind()) // float64

	// reflect.ValueOf() 获取值信息
	v := reflect.ValueOf(x)
	fmt.Printf("ValueOf(x) = %v\n", v)        // 3.14
	fmt.Printf("Value.Type = %v\n", v.Type()) // float64
	fmt.Printf("Value.Kind = %v\n", v.Kind()) // float64
	fmt.Printf("Value.Float= %v\n", v.Float()) // 3.14

	// -------------------------------------------------------------------------
	// 2. Type vs Kind 的区别（重要！）
	// -------------------------------------------------------------------------
	fmt.Println("\n--- 2. Type vs Kind 的区别 ---")

	var myInt ReflectInt = 42

	t2 := reflect.TypeOf(myInt)
	fmt.Printf("类型名 (Name): %s\n", t2.Name()) // MyInt（自定义类型名）
	fmt.Printf("种类 (Kind):  %s\n", t2.Kind())  // int（底层种类）

	// Kind 是有限的枚举，Type 可以是任意自定义类型
	fmt.Println("\n常见 Kind 值:")
	fmt.Println("  Bool, Int, Int8...Int64, Uint...Uint64")
	fmt.Println("  Float32, Float64, Complex64, Complex128")
	fmt.Println("  Array, Chan, Func, Interface, Map, Pointer, Slice, String, Struct")

	// -------------------------------------------------------------------------
	// 3. 各种类型的反射
	// -------------------------------------------------------------------------
	fmt.Println("\n--- 3. 各种类型的反射 ---")

	values := []interface{}{
		42,
		3.14,
		"hello",
		true,
		[]int{1, 2, 3},
		map[string]int{"a": 1},
		ReflectUser{Name: "test"},
		&ReflectUser{Name: "ptr"},
		func(x int) int { return x * 2 },
		make(chan int),
	}

	for _, val := range values {
		t := reflect.TypeOf(val)
		v := reflect.ValueOf(val)
		fmt.Printf("  值: %-20v | Type: %-20s | Kind: %s\n",
			truncate(fmt.Sprintf("%v", val), 18), t, t.Kind())
		_ = v // 避免未使用警告
	}
}

// ============================================================================
// 【第二部分：反射三大定律】
// ============================================================================

func demoThreeLaws() {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("【第二部分：反射三大定律】")
	fmt.Println(strings.Repeat("=", 70))

	// -------------------------------------------------------------------------
	// 定律一：从接口值到反射对象
	// Reflection goes from interface value to reflection object.
	// -------------------------------------------------------------------------
	fmt.Println("\n--- 定律一：接口 → 反射对象 ---")

	var x int = 42

	// x 被隐式转为 interface{}，然后反射提取信息
	v := reflect.ValueOf(x)
	t := reflect.TypeOf(x)

	fmt.Printf("原值: %d (类型: int)\n", x)
	fmt.Printf("  → reflect.ValueOf → %v\n", v)
	fmt.Printf("  → reflect.TypeOf  → %v\n", t)

	// -------------------------------------------------------------------------
	// 定律二：从反射对象到接口值
	// Reflection goes from reflection object to interface value.
	// -------------------------------------------------------------------------
	fmt.Println("\n--- 定律二：反射对象 → 接口 ---")

	// 用 Interface() 方法转回 interface{}
	i := v.Interface()
	fmt.Printf("v.Interface() = %v (类型: %T)\n", i, i)

	// 类型断言恢复原始类型
	y := i.(int)
	fmt.Printf("类型断言后: %d (类型: %T)\n", y, y)

	// -------------------------------------------------------------------------
	// 定律三：要修改反射对象，值必须可设置
	// To modify a reflection object, the value must be settable.
	// -------------------------------------------------------------------------
	fmt.Println("\n--- 定律三：修改需要可设置性 ---")

	// 错误示例：传值，不可修改
	var a float64 = 3.14
	va := reflect.ValueOf(a)
	fmt.Printf("传值: CanSet() = %v (不能修改副本)\n", va.CanSet())

	// 正确示例：传指针，可以修改
	vp := reflect.ValueOf(&a)
	fmt.Printf("传指针: CanSet() = %v (指针本身不能 Set)\n", vp.CanSet())

	// 用 Elem() 获取指针指向的元素
	ve := vp.Elem()
	fmt.Printf("Elem(): CanSet() = %v (可以修改!)\n", ve.CanSet())

	// 修改值
	ve.SetFloat(7.28)
	fmt.Printf("修改后 a = %v\n", a)
}

// ============================================================================
// 【第三部分：结构体反射】
// ============================================================================

func demoStructReflection() {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("【第三部分：结构体反射】")
	fmt.Println(strings.Repeat("=", 70))

	// -------------------------------------------------------------------------
	// 1. 遍历结构体字段
	// -------------------------------------------------------------------------
	fmt.Println("\n--- 1. 遍历结构体字段 ---")

	user := ReflectUser{
		Name:    "张三",
		Age:     25,
		Email:   "zhang@example.com",
		private: "secret",
	}

	t := reflect.TypeOf(user)
	v := reflect.ValueOf(user)

	fmt.Printf("类型: %s, 字段数: %d\n\n", t.Name(), t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)  // StructField 类型信息
		value := v.Field(i)  // Value 值信息

		// 检查是否可导出（PkgPath 为空表示可导出）
		exported := field.PkgPath == ""

		fmt.Printf("字段 %d: %s\n", i, field.Name)
		fmt.Printf("  类型:    %s\n", field.Type)
		fmt.Printf("  Tag:     %s\n", field.Tag)
		fmt.Printf("  可导出:  %v\n", exported)

		if exported {
			fmt.Printf("  值:      %v\n", value.Interface())
		} else {
			fmt.Printf("  值:      (不可访问 - 私有字段)\n")
		}
		fmt.Println()
	}

	// -------------------------------------------------------------------------
	// 2. Tag 解析
	// -------------------------------------------------------------------------
	fmt.Println("--- 2. Tag 解析 ---")

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue // 跳过私有字段
		}

		jsonTag := field.Tag.Get("json")
		validateTag := field.Tag.Get("validate")

		fmt.Printf("%s:\n", field.Name)
		fmt.Printf("  json:     %q\n", jsonTag)
		fmt.Printf("  validate: %q\n", validateTag)
	}

	// -------------------------------------------------------------------------
	// 3. 按名称获取字段
	// -------------------------------------------------------------------------
	fmt.Println("\n--- 3. 按名称获取字段 ---")

	// FieldByName 返回 (StructField, bool)
	if field, ok := t.FieldByName("Name"); ok {
		fmt.Printf("找到字段 'Name': 类型=%s, Tag=%s\n", field.Type, field.Tag)
	}

	if field, ok := t.FieldByName("NotExist"); !ok {
		fmt.Printf("字段 'NotExist' 不存在: %v\n", field.Name == "")
	}

	// -------------------------------------------------------------------------
	// 4. 动态修改结构体字段
	// -------------------------------------------------------------------------
	fmt.Println("\n--- 4. 动态修改结构体字段 ---")

	// 必须传指针才能修改
	userPtr := &ReflectUser{Name: "原始名", Age: 20}
	fmt.Printf("修改前: %+v\n", userPtr)

	vp := reflect.ValueOf(userPtr).Elem()

	// 通过字段名修改
	nameField := vp.FieldByName("Name")
	if nameField.CanSet() {
		nameField.SetString("新名字")
	}

	ageField := vp.FieldByName("Age")
	if ageField.CanSet() {
		ageField.SetInt(30)
	}

	fmt.Printf("修改后: %+v\n", userPtr)

	// -------------------------------------------------------------------------
	// 5. 通用的字段设置函数
	// -------------------------------------------------------------------------
	fmt.Println("\n--- 5. 通用字段设置函数 ---")

	user2 := &ReflectUser{Name: "test"}
	fmt.Printf("设置前: %+v\n", user2)

	err := setField(user2, "Name", "动态设置")
	if err != nil {
		fmt.Println("错误:", err)
	}

	err = setField(user2, "Age", 100)
	if err != nil {
		fmt.Println("错误:", err)
	}

	fmt.Printf("设置后: %+v\n", user2)
}

// setField 通用字段设置函数
func setField(obj interface{}, name string, value interface{}) error {
	v := reflect.ValueOf(obj)

	// 检查是否是指针
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("必须传入指针，实际: %s", v.Kind())
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("必须是结构体指针，实际: %s", v.Kind())
	}

	field := v.FieldByName(name)
	if !field.IsValid() {
		return fmt.Errorf("字段 '%s' 不存在", name)
	}
	if !field.CanSet() {
		return fmt.Errorf("字段 '%s' 不可设置", name)
	}

	val := reflect.ValueOf(value)
	if field.Type() != val.Type() {
		return fmt.Errorf("类型不匹配: 字段是 %s, 值是 %s", field.Type(), val.Type())
	}

	field.Set(val)
	return nil
}

// ============================================================================
// 【第四部分：方法反射】
// ============================================================================

func demoMethodReflection() {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("【第四部分：方法反射】")
	fmt.Println(strings.Repeat("=", 70))

	// -------------------------------------------------------------------------
	// 1. 获取类型的方法列表
	// -------------------------------------------------------------------------
	fmt.Println("\n--- 1. 获取方法列表 ---")

	calc := ReflectCalc{Name: "计算器"}

	// 值类型只能看到值接收者的方法
	fmt.Println("值类型 Calculator 的方法:")
	listMethods(calc)

	// 指针类型能看到所有方法（值接收者 + 指针接收者）
	fmt.Println("\n指针类型 *Calculator 的方法:")
	listMethods(&calc)

	// -------------------------------------------------------------------------
	// 2. 动态调用方法
	// -------------------------------------------------------------------------
	fmt.Println("\n--- 2. 动态调用方法 ---")

	// 调用 Add 方法
	results, err := callMethod(calc, "Add", 10, 20)
	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Printf("calc.Add(10, 20) = %v\n", results[0])
	}

	// 调用 Sub 方法
	results, err = callMethod(calc, "Sub", 50, 30)
	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Printf("calc.Sub(50, 30) = %v\n", results[0])
	}

	// 调用指针接收者的方法
	results, err = callMethod(&calc, "Mul", 6, 7)
	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Printf("calc.Mul(6, 7) = %v\n", results[0])
	}

	// -------------------------------------------------------------------------
	// 3. 方法的详细信息
	// -------------------------------------------------------------------------
	fmt.Println("\n--- 3. 方法的详细信息 ---")

	t := reflect.TypeOf(calc)
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		fmt.Printf("方法 %d: %s\n", i, method.Name)
		fmt.Printf("  类型: %s\n", method.Type)
		fmt.Printf("  参数数量: %d (包含接收者)\n", method.Type.NumIn())
		fmt.Printf("  返回值数量: %d\n", method.Type.NumOut())

		// 打印参数类型
		fmt.Print("  参数类型: ")
		for j := 0; j < method.Type.NumIn(); j++ {
			if j > 0 {
				fmt.Print(", ")
			}
			fmt.Print(method.Type.In(j))
		}
		fmt.Println()

		// 打印返回值类型
		fmt.Print("  返回类型: ")
		for j := 0; j < method.Type.NumOut(); j++ {
			if j > 0 {
				fmt.Print(", ")
			}
			fmt.Print(method.Type.Out(j))
		}
		fmt.Println()
	}
}

// listMethods 列出类型的所有方法
func listMethods(x interface{}) {
	t := reflect.TypeOf(x)
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		fmt.Printf("  %s: %s\n", method.Name, method.Type)
	}
}

// callMethod 动态调用方法
func callMethod(obj interface{}, methodName string, args ...interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(obj)
	method := v.MethodByName(methodName)

	if !method.IsValid() {
		return nil, fmt.Errorf("方法 '%s' 不存在", methodName)
	}

	// 准备参数
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	// 调用方法
	results := method.Call(in)

	// 转换结果
	out := make([]interface{}, len(results))
	for i, r := range results {
		out[i] = r.Interface()
	}

	return out, nil
}

// ============================================================================
// 【第五部分：实用工具函数】
// ============================================================================

func demoUtilityFunctions() {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("【第五部分：实用工具函数】")
	fmt.Println(strings.Repeat("=", 70))

	// -------------------------------------------------------------------------
	// 1. 深度比较 (DeepEqual)
	// -------------------------------------------------------------------------
	fmt.Println("\n--- 1. DeepEqual 深度比较 ---")

	a := []int{1, 2, 3}
	b := []int{1, 2, 3}
	c := []int{1, 2, 4}

	fmt.Printf("a=%v, b=%v: DeepEqual=%v\n", a, b, reflect.DeepEqual(a, b)) // true
	fmt.Printf("a=%v, c=%v: DeepEqual=%v\n", a, c, reflect.DeepEqual(a, c)) // false

	// map 比较
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 2, "a": 1}
	m3 := map[string]int{"a": 1, "b": 3}

	fmt.Printf("m1=%v, m2=%v: DeepEqual=%v\n", m1, m2, reflect.DeepEqual(m1, m2)) // true
	fmt.Printf("m1=%v, m3=%v: DeepEqual=%v\n", m1, m3, reflect.DeepEqual(m1, m3)) // false

	// -------------------------------------------------------------------------
	// 2. 零值检查 (IsZero)
	// -------------------------------------------------------------------------
	fmt.Println("\n--- 2. IsZero 零值检查 ---")

	values := []interface{}{
		0,
		"",
		false,
		[]int(nil),
		map[string]int(nil),
		(*ReflectUser)(nil),
		ReflectUser{},
		1,
		"hello",
		[]int{1},
	}

	for _, val := range values {
		v := reflect.ValueOf(val)
		fmt.Printf("  %T(%v): IsZero=%v\n", val, truncate(fmt.Sprintf("%v", val), 10), v.IsZero())
	}

	// -------------------------------------------------------------------------
	// 3. 创建新值 (New, MakeSlice, MakeMap, MakeChan)
	// -------------------------------------------------------------------------
	fmt.Println("\n--- 3. 动态创建值 ---")

	// 创建 *int
	intType := reflect.TypeOf(0)
	ptrVal := reflect.New(intType)
	ptrVal.Elem().SetInt(42)
	fmt.Printf("reflect.New(int): %v (值=%d)\n", ptrVal.Type(), ptrVal.Elem().Int())

	// 创建 slice
	sliceType := reflect.SliceOf(intType)
	sliceVal := reflect.MakeSlice(sliceType, 3, 5)
	sliceVal.Index(0).SetInt(10)
	sliceVal.Index(1).SetInt(20)
	sliceVal.Index(2).SetInt(30)
	fmt.Printf("reflect.MakeSlice: %v = %v\n", sliceVal.Type(), sliceVal.Interface())

	// 创建 map
	mapType := reflect.MapOf(reflect.TypeOf(""), intType)
	mapVal := reflect.MakeMap(mapType)
	mapVal.SetMapIndex(reflect.ValueOf("key1"), reflect.ValueOf(100))
	mapVal.SetMapIndex(reflect.ValueOf("key2"), reflect.ValueOf(200))
	fmt.Printf("reflect.MakeMap: %v = %v\n", mapVal.Type(), mapVal.Interface())

	// -------------------------------------------------------------------------
	// 4. 类型断言检查
	// -------------------------------------------------------------------------
	fmt.Println("\n--- 4. 类型检查工具 ---")

	checkTypes := []interface{}{
		42,
		"hello",
		[]int{1, 2},
		map[string]int{},
		ReflectUser{},
		&ReflectUser{},
		func() {},
		make(chan int),
	}

	for _, val := range checkTypes {
		v := reflect.ValueOf(val)
		t := v.Type()

		isPtr := t.Kind() == reflect.Ptr
		isSlice := t.Kind() == reflect.Slice
		isMap := t.Kind() == reflect.Map
		isStruct := t.Kind() == reflect.Struct
		isFunc := t.Kind() == reflect.Func
		isChan := t.Kind() == reflect.Chan

		fmt.Printf("  %T: Ptr=%v Slice=%v Map=%v Struct=%v Func=%v Chan=%v\n",
			val, isPtr, isSlice, isMap, isStruct, isFunc, isChan)
	}
}

// ============================================================================
// 【第六部分：实战示例 - 简易 JSON 序列化】
// ============================================================================

func demoJSONSerializer() {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("【第六部分：实战示例 - 简易 JSON 序列化】")
	fmt.Println(strings.Repeat("=", 70))

	type Address struct {
		City   string `json:"city"`
		Street string `json:"street"`
	}

	type Person struct {
		Name    string   `json:"name"`
		Age     int      `json:"age"`
		Address Address  `json:"address"`
		Hobbies []string `json:"hobbies"`
	}

	p := Person{
		Name: "张三",
		Age:  25,
		Address: Address{
			City:   "北京",
			Street: "长安街",
		},
		Hobbies: []string{"编程", "阅读", "游泳"},
	}

	fmt.Println("\n原始结构体:")
	fmt.Printf("  %+v\n", p)

	fmt.Println("\n序列化为 JSON:")
	jsonStr := toJSON(reflect.ValueOf(p))
	fmt.Printf("  %s\n", jsonStr)
}

// toJSON 简易 JSON 序列化（使用反射）
func toJSON(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return fmt.Sprintf(`"%s"`, v.String())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", v.Uint())

	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%g", v.Float())

	case reflect.Bool:
		return fmt.Sprintf("%t", v.Bool())

	case reflect.Slice, reflect.Array:
		var items []string
		for i := 0; i < v.Len(); i++ {
			items = append(items, toJSON(v.Index(i)))
		}
		return "[" + strings.Join(items, ", ") + "]"

	case reflect.Map:
		var pairs []string
		iter := v.MapRange()
		for iter.Next() {
			key := toJSON(iter.Key())
			val := toJSON(iter.Value())
			pairs = append(pairs, fmt.Sprintf("%s: %s", key, val))
		}
		return "{" + strings.Join(pairs, ", ") + "}"

	case reflect.Struct:
		var fields []string
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath != "" { // 跳过私有字段
				continue
			}

			// 获取 json tag，如果没有则使用字段名
			key := field.Tag.Get("json")
			if key == "" {
				key = field.Name
			}
			// 处理 omitempty 等选项（简化：只取第一部分）
			if idx := strings.Index(key, ","); idx != -1 {
				key = key[:idx]
			}

			val := toJSON(v.Field(i))
			fields = append(fields, fmt.Sprintf(`"%s": %s`, key, val))
		}
		return "{" + strings.Join(fields, ", ") + "}"

	case reflect.Ptr:
		if v.IsNil() {
			return "null"
		}
		return toJSON(v.Elem())

	case reflect.Interface:
		if v.IsNil() {
			return "null"
		}
		return toJSON(v.Elem())

	default:
		return "null"
	}
}

// ============================================================================
// 【第七部分：性能注意事项】
// ============================================================================

func demoPerformanceNotes() {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("【第七部分：性能注意事项】")
	fmt.Println(strings.Repeat("=", 70))

	fmt.Println(`
【反射性能对比】
┌─────────────────────┬──────────────┬────────────────┐
│ 操作                │ 耗时         │ 相对性能       │
├─────────────────────┼──────────────┼────────────────┤
│ 直接字段访问        │ ~0.5 ns/op   │ 基准           │
│ 反射（每次查找）    │ ~200 ns/op   │ 慢 400 倍      │
│ 反射（缓存索引）    │ ~50 ns/op    │ 慢 100 倍      │
└─────────────────────┴──────────────┴────────────────┘

【性能优化建议】

1. 缓存类型信息
   ❌ 错误：每次循环都调用 TypeOf
   for _, item := range items {
       t := reflect.TypeOf(item)  // 重复获取
       ...
   }

   ✅ 正确：提前缓存
   t := reflect.TypeOf(items[0])
   for _, item := range items {
       v := reflect.ValueOf(item)
       ...
   }

2. 缓存字段索引
   ❌ 错误：每次都按名称查找
   v.FieldByName("Name")  // 需要遍历字段

   ✅ 正确：使用索引
   idx := t.FieldByName("Name").Index
   v.FieldByIndex(idx)  // 直接定位

3. 避免热点路径使用反射
   - Web 框架的路由匹配：预先编译
   - 高频调用的序列化：考虑代码生成
   - 性能关键的业务逻辑：使用泛型或手写代码

4. 考虑代码生成替代反射
   - go generate + text/template
   - 第三方工具如 easyjson, go-swagger
`)

	// 演示缓存优化
	fmt.Println("--- 缓存优化示例 ---")

	type Data struct {
		ID   int
		Name string
	}

	data := Data{ID: 1, Name: "test"}

	// 缓存类型信息
	t := reflect.TypeOf(data)
	idField, _ := t.FieldByName("ID")
	nameField, _ := t.FieldByName("Name")

	fmt.Printf("缓存的字段索引: ID=%v, Name=%v\n", idField.Index, nameField.Index)

	// 使用缓存的索引
	v := reflect.ValueOf(data)
	idVal := v.FieldByIndex(idField.Index).Int()
	nameVal := v.FieldByIndex(nameField.Index).String()

	fmt.Printf("使用缓存索引获取值: ID=%d, Name=%s\n", idVal, nameVal)
}

// ============================================================================
// 【辅助函数】
// ============================================================================

// truncate 截断字符串
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// ============================================================================
// 【主函数】
// ============================================================================

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║              Go 语言反射 (Reflection) 完全指南                       ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")

	// 第一部分：基础
	demoBasicReflection()

	// 第二部分：三大定律
	demoThreeLaws()

	// 第三部分：结构体反射
	demoStructReflection()

	// 第四部分：方法反射
	demoMethodReflection()

	// 第五部分：实用工具
	demoUtilityFunctions()

	// 第六部分：实战示例
	demoJSONSerializer()

	// 第七部分：性能注意事项
	demoPerformanceNotes()

	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("【反射使用总结】")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println(`
✅ 适合使用反射的场景：
   • 处理 interface{} 类型的动态数据
   • 实现序列化/反序列化（JSON, XML, etc.）
   • ORM 框架、依赖注入
   • 基于 struct tag 的功能

❌ 不适合使用反射的场景：
   • 已知类型的普通业务逻辑
   • 性能关键路径
   • 可以用泛型解决的问题（Go 1.18+）

💡 反射口诀：
   1. 接口到反射，TypeOf 和 ValueOf
   2. 反射到接口，Interface() 拿回来
   3. 要改值，必传指针，Elem() 取元素
`)
}
