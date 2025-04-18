package dataschema

import (
	"fmt"
	"strings"
)

func ExampleTblToStructHandler_GenerateAllTblStruct() {

	//案例1 简易
	{
		th := NewTblToStructHandler()
		th.SetDsn("root:tiger@(127.0.0.1:3306)/pulingfu?charset=utf8mb4&parseTime=True&loc=Local").
			GenerateAllTblStruct()
	}

	//案例2 复杂
	{
		th := NewTblToStructHandler()
		th.SetDsn("root:tiger@(127.0.0.1:3306)/pulingfu?charset=utf8mb4&parseTime=True&loc=Local").
			SetStructOrmTag("gorm").     //设置所生成对应的orm 标记类型
			SetOtherTags("json", "msg"). //添加其他的标签 如json ==> `json:"xxx"` msg ==> `msg:"xxx"`
			SeTblStructColumnNameInfo(
				CAMEL_CASE,                         //设置字段名写法类型为骆驼写法
				FIELD_ORDER_FIELD_NAME,             // 设置字段名排序方式为字段名称排序
				"column_prefix_", "_column_suffix", // 设置字段名前缀和后缀 可以为空字符串
			).
			SetTblStructNameInfo(CAMEL_CASE, "tbl_prefix_", "_tbl_suffix"). //设置生成的结构体名类型为CamelCase写法，以及前后缀
			GenerateAllTblStruct()
	}

	// 案例3 高度自定义
	{
		th := NewTblToStructHandler()
		savePrefix := "./pkg/models/tbl_sql_auto_model/" //设置保存路径前缀
		th.SetDsn("root:tiger@(127.0.0.1:3306)/pulingfu?charset=utf8mb4&parseTime=True&loc=Local")
		fmt.Printf("\x1b[%dm>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>： \x1b[0m\n", 34)
		fmt.Printf("\x1b[%dm您即将自动生成go struct model： \x1b[0m\n", 34)
		fmt.Printf("\x1b[%dm多表可以使用逗号隔开： \x1b[0m\n", 34)
		fmt.Printf("\x1b[%dm请输入要生成的sql表名： \x1b[0m\n", 34)
		fmt.Printf("\x1b[%dm>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>： \x1b[0m\n", 34)
		commend := ""
		fmt.Scanln(&commend)
		commend = strings.ReplaceAll(commend, "，", ",")
		commend = strings.ReplaceAll(commend, "-", ",")
		commend = strings.ReplaceAll(commend, ".", ",")
		commend = strings.ReplaceAll(commend, "。", ",")
		commends := strings.Split(commend, ",")
		for _, commend = range commends {
			if commend == "" {
				continue
			}
			th.SetTableName(commend).
				SetSavePath(fmt.Sprintf("%s/tbl_%s/tbl_%s.go", savePrefix, commend, commend)).
				SetPackageInfo(commend, "tbl_", "").
				GenerateTblStruct()
		}
	}

}

func ExampleYamlToSqlHandler_ExecuteSchemaSafeCheck() {

	// 配置文件请参考 目录./cmd/test_yaml_to_sql/etc2/ 下的案例
	{
		yts := NewYamlToSqlHandler().SetYamlPath("./cmd/test_yaml_to_sql/etc2/").
			SetDsn("root:tiger@(127.0.0.1:3306)/pulingfu?charset=utf8mb4&parseTime=True&loc=Local")
		yts.ExecuteSchemaSafeCheck()
	}

	// 配置文件请参考 目录./cmd/test_yaml_to_sql/etc/ 下的案例
	{
		yts := NewYamlToSqlHandler().SetYamlPath("./cmd/test_yaml_to_sql/etc/").
			SetDsn("root:tiger@(127.0.0.1:3306)/pulingfu?charset=utf8mb4&parseTime=True&loc=Local")
		yts.ExecuteSchemaSafeCheck()
	}

}
