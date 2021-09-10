//测试将mysql单表生成 go struct model
package main

import (
	"github.com/pulingfu/tblschema"
)

func main() {
	th := tblschema.NewTblSchemaHandler()
	th = th.
		//设置数据库dsn连接地址
		SetDsn("root:tiger@(127.0.0.1:3306)/pulingfu?charset=utf8mb4&parseTime=True&loc=Local").
		//设置生成的model包名
		SetPackageName("tbl_my_user").
		//设置要生成哪张数据库表的结构
		SetTableName("my_user").

		//设置结构体名 骆驼写法/首字母大写  不设置则默认骆驼写法如 tbl_user => TblUser
		SetStructNameType(tblschema.CAMEL_CASE).

		//设置结构体字段名 骆驼写法/首字母大写  不设置则默认骆驼写法如 user_name => UserName
		SefieldNameType(tblschema.CAMEL_CASE).

		//生成结构体标记的orm类型 默认为gorm
		// orm => `orm:"column_name"`
		// gorm => `gorm:"column:column_name"`
		SetModelOrmTagType("gorm").

		//设置其他结构体标签，如json，request
		// SetOtherTag("json", "request").
		SetOtherTag("json").

		//设置生成模型的保存位置
		SetSavePath("./tbl_my_user/model.go").
		Run()

	///************/////
	//继承上面的配置生成另外 的结构，如time类型对应成string
	th.
		SetPackageName("tbl_my_user_timestring").
		SetOtherTag("").
		SetSavePath("./tbl_my_user_time_string/model.go").

		//设置数据库time类型变为结构体string  默认为time.Time
		SetTimeType(tblschema.TIMETYPE_STRING).

		//设置字段排序方式FIELD_ORDER_ORDINAL_POSITION 按数据库字段建立顺序 默认字典顺序
		SetFieldOrder(tblschema.FIELD_ORDER_ORDINAL_POSITION).
		Run()

}
