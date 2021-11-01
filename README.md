# tblschema
根据mysql表结构生成对应的go struct model

找到./cmd/..下的main.go  执行命令: go  run main.go 文件开启自动化之旅


<<<<<<<<<<<<<<<<<<<<<<<<
详情请查看演示demo

生成当前数据库中全部表对应的 go struct model
演示demo：./cmd/test_alltabl_to_model  

生成当前数据库中全部表对应的 go struct model
演示demo: ./cmd/test_alltbl_by_youself

生成指定表 的go struct model
演示demo: ./cmd/test_tbl_to_model

简易版生成指定表的 go struct model
演示demo: ./cmd/test_tbl_to_model

使用配置文件维护表结构（
    写struct tag是一件很麻烦的事情，
    写纯 sql文也是一件很麻烦的事情，
    用工具建表一个字段一个字段建立，更慢了
    用同一个配置文件自动维护多个数据库是一个很不错的选择哦
)
演示demo: ./cmd/test_yaml_to_sql

>>>>>>>>>>>>>>>>>>>>>>>>




<<<<<<<<<<<<<<<<<<<
mysql数据类型：
"int",
"integer", 
"tinyint",
"smallint",
"mediumint", 
"bigint",
"int unsigned", 
"integer unsigned", 
"tinyint unsigned", 
"smallint unsigned",
"mediumint unsigned", 
"bigint unsigned", 
"bit"

浮点型
"float", 
"double",
"decimal"

bool型
"bool"

字符
"enum", 
"set", 
"varchar", 
"char", 
"tinytext",
"mediumtext", 
"text", 
"longtext", 
"blob", 
"tinyblob",
"mediumblob", 
"longblob", 
"binary", 
"varbinary"

时间
"date", 
"datetime", 
"timestamp", 
"time"
>>>>>>>>>>>>>>>>>>>


yml配置文件格式参考test里面的yml文件哦

