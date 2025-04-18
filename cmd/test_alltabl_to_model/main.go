package main

import dataschema "github.com/k-kkong/dataschema"

// 调用一键生成表结构API自动生成表结构
// 输入命令go run main.go
func main() {
	th := dataschema.NewTblToStructHandler()
	//简易
	th.SetDsn("root:tiger@(127.0.0.1:3306)/pulingfu?charset=utf8mb4&parseTime=True&loc=Local").
		GenerateAllTblStruct()

	//复杂
	// th.SetDsn("root:tiger@(127.0.0.1:3306)/pulingfu?charset=utf8mb4&parseTime=True&loc=Local").
	// 	SetStructOrmTag("gorm").
	// 	SeTblStructColumnNameInfo(
	// 		dataschema.CAMEL_CASE, dataschema.
	// 			FIELD_ORDER_FIELD_NAME, "column_prefix_", "_column_suffix",
	// 	).SetTblStructNameInfo(dataschema.CAMEL_CASE, "tbl_prefix_", "_tbl_suffix").
	// 	GenerateAllTblStruct()
}
