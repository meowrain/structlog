# 前言

今天在写项目的时候测试函数打印结构体，发现用 golang 的 fmt 输出整个结构体的内容实在是太费劲了，而且也不美观（美观需要花费大量精力写）
突然就想到，可以利用 golang 的反射获取 tag 里面的 Key value 作为结构体字段的说明，结构体字段作为 value

就写了 testlog 这个库

# 说明

本文将介绍一个名为 testlogs 的 Go 语言工具包，它可以帮助开发者轻松地解析和输出结构体字段信息，并支持通过 testlog 标签自定义输出字段名。

## testlog 包含的函数

testlogs 包提供了两个主要函数：

`LogStructFields(v interface{}) map[string]interface{}`：解析结构体并返回一个包含字段信息的 map，其中键为字段名（或 testlog 标签的值），值为字段的值。

`LogStruct(v interface{}) string`：将结构体内容格式化为字符串，使用 testlog 标签作为输出标识。

## 源码说明

---

### 1. 包声明和导入

```go
package testlogs

import (
	"fmt"
	"reflect"
	"strings"
)
```

- **`package testlogs`**：声明当前文件属于 `testlogs` 包。
- **`import`**：导入依赖的 Go 标准库包：
  - `fmt`：用于格式化输出。
  - `reflect`：用于反射操作，动态获取结构体的字段信息。
  - `strings`：用于字符串操作，例如拼接字符串。

---

### 2. `LogStructFields` 函数

```go
func LogStructFields(v interface{}) map[string]interface{} {
```

- **`LogStructFields`**：函数名，用于解析结构体并返回一个 `map`，其中键为字段名（或 `testlog` 标签的值），值为字段的值。
- **`v interface{}`**：参数是一个空接口类型，可以接收任意类型的值。

---

#### 2.1 反射获取值的类型

```go
val := reflect.ValueOf(v)
```

- **`reflect.ValueOf(v)`**：通过反射获取参数 `v` 的 `reflect.Value` 对象，用于后续操作。

---

#### 2.2 处理指针类型

```go
if val.Kind() == reflect.Ptr {
    val = val.Elem()
}
```

- **`val.Kind()`**：获取 `val` 的类型种类（`Kind`）。
- **`reflect.Ptr`**：判断 `val` 是否为指针类型。
- **`val.Elem()`**：如果 `val` 是指针类型，解引用指针，获取指针指向的实际值。

---

#### 2.3 检查是否为结构体

```go
if val.Kind() != reflect.Struct {
    return nil
}
```

- **`reflect.Struct`**：判断 `val` 是否为结构体类型。
- 如果不是结构体类型，返回 `nil`。

---

#### 2.4 初始化结果 `map`

```go
result := make(map[string]interface{})
```

- **`make(map[string]interface{})`**：初始化一个 `map`，用于存储解析后的字段信息。

---

#### 2.5 获取结构体类型信息

```go
typ := val.Type()
```

- **`val.Type()`**：获取 `val` 的类型信息（`reflect.Type` 对象），用于后续获取字段信息。

---

#### 2.6 遍历结构体字段

```go
for i := 0; i < val.NumField(); i++ {
```

- **`val.NumField()`**：获取结构体的字段数量。
- 使用 `for` 循环遍历结构体的每一个字段。

---

#### 2.7 获取字段信息

```go
field := typ.Field(i)
fieldValue := val.Field(i)
```

- **`typ.Field(i)`**：获取第 `i` 个字段的类型信息（`reflect.StructField` 对象）。
- **`val.Field(i)`**：获取第 `i` 个字段的值（`reflect.Value` 对象）。

---

#### 2.8 获取 `testlog` 标签

```go
tag := field.Tag.Get("testlog")
if tag == "" {
    tag = field.Name
}
```

- **`field.Tag.Get("testlog")`**：获取字段的 `testlog` 标签值。
- 如果 `testlog` 标签为空，则使用字段名作为 `tag`。

---

#### 2.9 检查字段是否可导出

```go
if !fieldValue.CanInterface() {
    continue
}
```

- **`fieldValue.CanInterface()`**：检查字段是否可以被 `Interface()` 方法访问（即是否为导出字段）。
- 如果字段不可导出，跳过该字段。

---

#### 2.10 处理嵌套结构体

```go
if fieldValue.Kind() == reflect.Struct {
    nested := LogStructFields(fieldValue.Interface())
    for k, v := range nested {
        result[tag+"."+k] = v
    }
    continue
}
```

- **`fieldValue.Kind() == reflect.Struct`**：判断字段是否为结构体类型。
- **`LogStructFields(fieldValue.Interface())`**：递归调用 `LogStructFields`，解析嵌套结构体。
- **`result[tag+"."+k] = v`**：将嵌套结构体的字段信息合并到当前结果 `map` 中，键为 `父字段.子字段` 的形式。

---

#### 2.11 处理指针类型字段

```go
if fieldValue.Kind() == reflect.Ptr {
    if fieldValue.IsNil() {
        result[tag] = nil
        continue
    }
    fieldValue = fieldValue.Elem()
}
```

- **`fieldValue.Kind() == reflect.Ptr`**：判断字段是否为指针类型。
- **`fieldValue.IsNil()`**：检查指针是否为 `nil`。
  - 如果为 `nil`，将字段值设置为 `nil`，并跳过该字段。
- **`fieldValue.Elem()`**：解引用指针，获取指针指向的实际值。

---

#### 2.12 存储字段值

```go
result[tag] = fieldValue.Interface()
```

- **`fieldValue.Interface()`**：将字段值转换为 `interface{}` 类型。
- **`result[tag]`**：将字段值存储到结果 `map` 中，键为 `tag`。

---

#### 2.13 返回结果

```go
return result
```

- 返回解析后的字段信息 `map`。

---

### 3. `LogStruct` 函数

```go
func LogStruct(v interface{}) string {
```

- **`LogStruct`**：函数名，用于将结构体内容格式化为字符串。
- **`v interface{}`**：参数是一个空接口类型，可以接收任意类型的值。

---

#### 3.1 获取字段信息

```go
fields := LogStructFields(v)
```

- 调用 `LogStructFields` 函数，获取结构体的字段信息。

---

#### 3.2 初始化字符串构建器

```go
var sb strings.Builder
```

- **`strings.Builder`**：用于高效拼接字符串。

---

#### 3.3 遍历字段信息并构建字符串

```go
for k, v := range fields {
    sb.WriteString(k)
    sb.WriteString(": ")
    sb.WriteString(fmt.Sprintf("%v", v))
    sb.WriteString("\n")
}
```

- **`sb.WriteString(k)`**：将字段名写入字符串构建器。
- **`sb.WriteString(": ")`**：写入分隔符 `: `。
- **`fmt.Sprintf("%v", v)`**：将字段值格式化为字符串。
- **`sb.WriteString("\n")`**：写入换行符。

---

#### 3.4 返回结果字符串

```go
return sb.String()
```

- 返回构建好的字符串。

---

### 4. 总结

`testlogs` 包的核心功能是通过反射解析结构体，并根据 `testlog` 标签输出字段信息。以下是关键点：

1. **反射**：使用 `reflect` 包动态获取结构体的字段信息。
2. **递归处理嵌套结构体**：支持解析嵌套的结构体字段。
3. **指针处理**：支持解析指针类型的字段。
4. **标签支持**：通过 `testlog` 标签自定义输出字段名。
5. **字符串构建**：使用 `strings.Builder` 高效拼接字符串。

## 功能详解

3.1 LogStructFields 函数
LogStructFields 函数的主要功能是解析结构体并返回一个 map[string]interface{}，其中键为字段名（或 testlog 标签的值），值为字段的值。

3.1.1 处理指针和嵌套结构体
指针类型：如果字段是指针类型，函数会检查指针是否为 nil。如果为 nil，则将该字段的值设置为 nil；否则，函数会解引用指针并继续处理。

嵌套结构体：如果字段是结构体类型，函数会递归调用 LogStructFields 来处理嵌套结构体，并将结果合并到当前的结果 map 中。

3.1.2 处理未导出字段
Go 语言中的未导出字段（即首字母小写的字段）无法通过 Interface() 方法访问。因此，LogStructFields 会跳过这些字段，避免运行时错误。

3.1.3 testlog 标签
如果字段有 testlog 标签，则使用标签的值作为 map 中的键。

如果没有 testlog 标签，则使用字段名作为键。

3.2 LogStruct 函数
LogStruct 函数的主要功能是将结构体内容格式化为字符串。它调用了 LogStructFields 函数来获取字段信息，并将其格式化为易于阅读的字符串。

3.2.1 格式化输出
每个字段的键值对以 key: value 的形式输出。

嵌套结构体的字段会以 parent.child 的形式显示，其中 parent 是父结构体的字段名或 testlog 标签，child 是子结构体的字段名或 testlog 标签。

4. 适用场景
   testlogs 包适用于以下场景：

调试：在调试复杂结构体时，快速查看结构体的内容。

日志记录：在日志中记录结构体的状态，便于后续分析。

测试：在单元测试中验证结构体的字段值是否符合预期。

# 使用效果

![](https://blog.meowrain.cn/api/i/2025/01/19/gIQCz91737295989649076069.avif)

测试代码：

```go
package testlogs

import (
	"fmt"
	"testing"
)

type Address struct {
	City    string `testlog:"城市"`
	Country string `testlog:"国家"`
}
type User struct {
	Name    string `testlog:"用户名"`
	Age     int    `testlog:"年龄"`
	Email   string `testlog:"邮箱"`
	Address `testlog:"地址"`
}

func TestTestLogs(t *testing.T) {
	user := User{
		Name:  "Alice",
		Age:   30,
		Email: "alice@example.com",
		Address: Address{
			City:    "上海",
			Country: "中国",
		},
	}
	fmt.Println("==================================")
	fmt.Println("不使用testlogs: ", user)
	fmt.Println("==================================")

	fmt.Println()
	fmt.Println()
	fmt.Println("==================================")
	fmt.Println("使用testlogs: ")
	fmt.Println(LogStruct(user))
	fmt.Println("==================================")
}

```

# 源代码

```go
package testlogs

import (
	"fmt"
	"reflect"
	"strings"
)

// LogStructFields 解析结构体并输出带有testlog标签的字段信息
func LogStructFields(v interface{}) map[string]interface{} {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	result := make(map[string]interface{})
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// 获取testlog标签
		tag := field.Tag.Get("testlog")
		if tag == "" {
			tag = field.Name // 如果没有testlog标签，使用字段名
		}

		// 检查字段是否可以被 Interface() 方法访问
		if !fieldValue.CanInterface() {
			continue // 跳过未导出的字段
		}

		// 处理嵌套结构体
		if fieldValue.Kind() == reflect.Struct {
			nested := LogStructFields(fieldValue.Interface())
			for k, v := range nested {
				result[tag+"."+k] = v
			}
			continue
		}

		// 处理指针类型
		if fieldValue.Kind() == reflect.Ptr {
			if fieldValue.IsNil() {
				result[tag] = nil
				continue
			}
			fieldValue = fieldValue.Elem()
		}

		result[tag] = fieldValue.Interface()
	}

	return result
}

// LogStruct 打印结构体内容，使用testlog标签作为输出标识
func LogStruct(v interface{}) string {
	fields := LogStructFields(v)
	var sb strings.Builder

	for k, v := range fields {
		sb.WriteString(k)
		sb.WriteString(": ")
		sb.WriteString(fmt.Sprintf("%v", v))
		sb.WriteString("\n")
	}

	return sb.String()
}

```
