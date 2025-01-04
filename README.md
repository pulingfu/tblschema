# tblschema
通过yaml配置文件生成对应数据库字段
```
Table:
  type: entity
  table: company_test  # 表名
  sharding_tables: company_test_01,company_test_02  # 分表指定表明，如果指定此字段则会按照schema创建分表
  options:
    charset: utf8
    collate: utf8_general_ci
    comment: 测试公司表
  indexes:  ## 索引
    idx_name:
      columns:
        - name
    idx_creator_uuid:
      columns:
        - creator_uuid
    idx_created_at:
      columns:
        - created_at
    idx_search_helper:
      columns:
        - search_helper
  unique_indexes:  ## 唯一索引
      uniq_uuid:
          columns: 
            - uuid 
            - uuid2
  fulltext_indexes:  ## 全文索引
      fulltext_intro:
          with_parser: simple ## 分词器
          columns: 
            - intro 
            - uuid2
  id:  ## 主键 至少要有一个字段
    id:  #第一个主键字段
      type: integer unsigned
      nullable: false
      generator: AUTO_INCREMENT
    uuid:  #第二个主键字段
      type: integer unsigned
      nullable: false
    primary_key2:  #第三个主键字段
      type: integer unsigned
      nullable: false
  fields:  ## 字段
    uuid2:
      type: varchar
      nullable: false
      comment:  用户 uuid
    name:
      type: varchar
      nullable: false
      comment:  公司名称
    creator_uuid:
      type: varchar
      nullable: false
      comment:  创建者
    cnt_project:
      type: integer
      nullable: false
      comment:  项目数量
    search_helper:
      type: varchar
      nullable: false
      comment: search_helper
    intro:
      type: text
      nullable: false
      comment: 简介
    created_at:
      type: datetime
      nullable: false
```
```
func main() {
	yts := tblschema.NewYamlToSqlHandler().SetYamlPath("./etc2/").
		SetDsn("root:tiger@(127.0.0.1:3306)/pulingfu?charset=utf8mb4&parseTime=True&loc=Local")
	yts.ExecuteSchemaSafeCheck()
}
```
将会生成数据库语句
```
CREATE TABLE company_test_01(
        id integer unsigned AUTO_INCREMENT NOT NULL,
        primary_key2 integer unsigned  NOT NULL,
        uuid integer unsigned  NOT NULL,
        cnt_project integer NOT NULL  COMMENT '项目数量' ,
        created_at datetime NOT NULL  COMMENT '' ,
        creator_uuid varchar(255) NOT NULL  COMMENT '创建者' ,
        intro text NOT NULL COMMENT '简介' ,
        name varchar(255) NOT NULL  COMMENT '公司名称' ,
        search_helper varchar(255) NOT NULL  COMMENT 'search_helper' ,
        uuid2 varchar(255) NOT NULL  COMMENT '用户 uuid' ,
        INDEX idx_created_at (created_at),
        INDEX idx_creator_uuid (creator_uuid),
        INDEX idx_name (name),
        INDEX idx_search_helper (search_helper),
        UNIQUE INDEX uniq_uuid (uuid,uuid2),
        FULLTEXT INDEX fulltext_intro (intro,uuid2),
        PRIMARY KEY( id , primary_key2 , uuid )
)
DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci ENGINE = InnoDB  COMMENT = '测试公司表' ;
CREATE TABLE company_test_02(
        id integer unsigned AUTO_INCREMENT NOT NULL,
        primary_key2 integer unsigned  NOT NULL,
        uuid integer unsigned  NOT NULL,
        cnt_project integer NOT NULL  COMMENT '项目数量' ,
        created_at datetime NOT NULL  COMMENT '' ,
        creator_uuid varchar(255) NOT NULL  COMMENT '创建者' ,
        intro text NOT NULL COMMENT '简介' ,
        name varchar(255) NOT NULL  COMMENT '公司名称' ,
        search_helper varchar(255) NOT NULL  COMMENT 'search_helper' ,
        uuid2 varchar(255) NOT NULL  COMMENT '用户 uuid' ,
        INDEX idx_created_at (created_at),
        INDEX idx_creator_uuid (creator_uuid),
        INDEX idx_name (name),
        INDEX idx_search_helper (search_helper),
        UNIQUE INDEX uniq_uuid (uuid,uuid2),
        FULLTEXT INDEX fulltext_intro (intro,uuid2),
        PRIMARY KEY( id , primary_key2 , uuid )
)
DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci ENGINE = InnoDB  COMMENT = '测试公司表' ;
```

更多案例请查看 /cmd/test...

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