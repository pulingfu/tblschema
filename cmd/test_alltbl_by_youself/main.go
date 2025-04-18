package main

import (
	"fmt"

	tblschema "github.com/k-kkong/dataschema"
)

// 自定义生成想要的表结构
func main() {
	th := tblschema.NewTblToStructHandler().
		SetDsn("root:tiger@(127.0.0.1:3306)/pulingfu?charset=utf8mb4&parseTime=True&loc=Local").

		// 添加其他标签？比如json
		SetOtherTags("json").
		//设置包名
		SetPackageInfo("all_tbl_model", "", "")

	for _, tname := range th.GetAllTableNames() {
		th.
			//设置
			SetSavePath(fmt.Sprintf("./all_tbl_model/%s.go", tname)).
			SetTableName(tname).
			GenerateTblStruct()
	}

}
