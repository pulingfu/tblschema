// 测试将mysql单表生成 go struct model
package main

import (
	dataschema "github.com/k-kkong/dataschema"
)

// //输入命令go run main.go
func main() {
	//简单用法
	simple := dataschema.NewTblToStructHandler()
	simple.
		SetDsn("root:tiger@(127.0.0.1:3306)/pulingfu?charset=utf8mb4&parseTime=True&loc=Local").
		SetTableName("my_user1").

		//默认路径为当前目录
		SetSavePath("./simple/model.go").

		// SetPackageInfo("plf_test_package", "", "").
		GenerateTblStruct()

	///复杂用法
	th := dataschema.NewTblToStructHandler()
	th = th.
		//设置数据库dsn连接地址
		SetDsn("root:tiger@(127.0.0.1:3306)/pulingfu?charset=utf8mb4&parseTime=True&loc=Local").
		//设置生成的model包名
		SetPackageInfo("tbl_my_user", "prefix_", "_suffix").
		//设置要生成哪张数据库表的结构
		SetTableName("my_user").

		//设置结构体字段名 骆驼写法/首字母大写  不设置则默认骆驼写法如 user_name => UserName
		//写法、前缀、后缀
		SetTblStructNameInfo(dataschema.CAMEL_CASE, "a我是_table_name前缀_", "_b我是_table_name后缀").

		//设置行信息  写法、排序方式、前缀、后缀
		SeTblStructColumnNameInfo(dataschema.CAMEL_CASE, "", "a我是_column_name前缀_", "_b我是_column_name后缀").

		//生成结构体标记的orm类型 默认为gorm
		// orm => `orm:"column_name"`
		// gorm => `gorm:"column:column_name"`
		SetStructOrmTag("gorm").

		//设置其他结构体标签，如json，request
		// SetOtherTag("json", "request").
		// SetOtherTag("json").

		//设置生成模型的保存位置
		SetSavePath("./tbl_my_user/model.go").
		GenerateTblStruct()

	///************/////
	//继承上面的配置生成另外 的结构，如time类型对应成string
	th.
		//依次为包名、包名前缀、包名后缀
		SetPackageInfo("tbl_my_user_timestring", "prefix_", "_suffix").
		SetOtherTags("").
		SetSavePath("./tbl_my_user_time_string/model.go").

		//设置数据库time类型变为结构体string  默认为time.Time
		SetTimeType(dataschema.TIMETYPE_STRING).

		//设置结构体名 骆驼写法/首字母大写  不设置则默CameCase写法如 tbl_user => TblUser
		SetTblStructNameInfo(dataschema.CAMEL_CASE, "", "").
		//设置字段排序方式FIELD_ORDER_ORDINAL_POSITION 按数据库字段建立顺序 默认字典顺序
		SeTblStructColumnNameInfo(dataschema.CAMEL_CASE, dataschema.FIELD_ORDER_FIELD_NAME, "", "").
		GenerateTblStruct()

}
