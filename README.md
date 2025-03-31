# Tblschema
## 
使用yaml配置文件，快速简便的管理数据库的表结构，提升开发效率：

[参考文档：ExecuteSchemaSafeCheck](https://pkg.go.dev/github.com/pulingfu/tblschema#example-YamlToSqlHandler.ExecuteSchemaSafeCheck)


## 
可以将数据库的表结构，一键翻译成go语言的结构体，避免繁琐的手写字段，提升开发效率

[参考文档：GenerateAllTblStruct](https://pkg.go.dev/github.com/pulingfu/tblschema#example-TblToStructHandler.GenerateAllTblStruct)

## 
使用dvaplugin 插件，可快捷加载嵌套的数据结构，实现动态的结构加载，可省去定义不同的结构体的内外键，解藕结构体防止互相引用，提升开发效率
[参考文档：NewRelationLoader](https://pkg.go.dev/github.com/pulingfu/tblschema/dvaplugin#NewRelationLoader)

##
任意切片数据处理，Map,Reduce,Find查找符合,take查找一个,Divide分割,Page翻页,Pop,remove,sort,判断....等许多操作，详细参考
[参考文档：NewSlicer](https://pkg.go.dev/github.com/pulingfu/tblschema/dvaplugin#NewSlicer)

## 
提供了一些专注于处理数据结构的func，可以提升处理数据的开放效率

## 其他使用案例参考各个test文件
#### dataer用法
```go
package main

import (
	"fmt"
	"github.com/pulingfu/tblschema/dvaplugin"
)

func main() {
	dataer := dvaplugin.NewDataer()
	// Do something
}
```






