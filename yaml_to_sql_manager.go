package tblschema

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"strings"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/pulingfu/tblschema/information_schema"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type YamlToSqlHandler struct {
	IsOutputBuildSchema      bool   // 是否输出编译后的结构
	IsEncryOutputBuildSchema bool   // 是否加密 编译后的结构
	EncryKey                 string // 加密key
	BuildSchemaDest          string //

	dsn string   //数据库连接dsn,列：用户:密码@(127.0.0.1:3306)/数据库?charset=utf8mb4&parseTime=True&loc=Local
	db  *gorm.DB //数据库连接

	YamlPath          string //yaml文件路径
	yamlFileFullPaths []string
	tables            []string

	sql []string
}

func NewYamlToSqlHandler() *YamlToSqlHandler {
	return &YamlToSqlHandler{
		IsOutputBuildSchema:      false,
		IsEncryOutputBuildSchema: false,
		BuildSchemaDest:          "./tblschema.value",
	}
}

func (ts *YamlToSqlHandler) SetDsn(dsn string) *YamlToSqlHandler {
	ts.dsn = dsn
	return ts
}
func (ts *YamlToSqlHandler) SetDB(db *gorm.DB) *YamlToSqlHandler {
	ts.db = db
	return ts
}

func (ts *YamlToSqlHandler) connectSql() {
	if ts.db == nil {
		if ts.dsn == "" {
			panic("数据库连接不能为空")
		}
		var configs = &gorm.Config{}
		db, err := gorm.Open(mysql.Open(ts.dsn), configs)
		if err != nil {
			panic(err)
		}
		ts.db = db
	}
}

func (ts *YamlToSqlHandler) SetYamlPath(yamlPath string) *YamlToSqlHandler {
	ts.YamlPath = yamlPath
	return ts
}

func (ts *YamlToSqlHandler) SetIsOutputBuildSchema(value bool, encry bool, key string) *YamlToSqlHandler {
	ts.IsOutputBuildSchema = value
	ts.IsEncryOutputBuildSchema = encry
	ts.EncryKey = key
	return ts
}
func (ts *YamlToSqlHandler) SetBuildSchemaDest(dest string) *YamlToSqlHandler {
	ts.BuildSchemaDest = dest
	return ts
}
func (ts *YamlToSqlHandler) getyamlFileFullPaths() *YamlToSqlHandler {

	files, err := os.ReadDir(ts.YamlPath)
	if err != nil {
		fmt.Println(err)
	}
	for _, f := range files {
		// fmt.Println(f.Name())
		filename := string(f.Name())
		if strings.Contains(filename, ".yml") {
			filename = fmt.Sprintf("%s%s", ts.YamlPath, filename)
			ts.yamlFileFullPaths = append(ts.yamlFileFullPaths, filename)
			// fmt.Println(ts.yamlFileFullPaths)
		}
	}

	return ts
}

func (ts *YamlToSqlHandler) getYamlDatas() *YamlToSqlHandler {
	var buildmapping = map[string]interface{}{}
	for _, v := range ts.yamlFileFullPaths {

		yamlFile, err := os.ReadFile(v)
		if err != nil {
			fmt.Println(err.Error())
		}
		table := map[string]interface{}{}
		err = yaml.Unmarshal(yamlFile, &table)
		if err != nil {
			fmt.Println(err.Error())
		}
		jb, err := json.Marshal(&table)

		tvalue := string(jb)
		tname := gjson.Parse(tvalue).Get("Table.table").String()
		if _, ok := buildmapping[tname]; ok {
			fmt.Printf("\x1b[%dm配置文件: %s 序列化失败，重复定义的表 \x1b[0m\n", 31, v)
			panic(fmt.Sprintf("\x1b[%dm配置文件: %s 序列化失败\x1b[0m\n", 31, v))
		}

		sharding_tables := gjson.Parse(tvalue).Get("Table.sharding_tables").String()
		if sharding_tables != "" {
			for _, sharding_name := range strings.Split(sharding_tables, ",") {
				if sharding_name == "" {
					continue
				}
				sharding_tblv, _ := sjson.Set(tvalue, "Table.table", sharding_name)
				ts.tables = append(ts.tables, sharding_tblv)
			}
		} else {
			ts.tables = append(ts.tables, tvalue)
		}
		buildmapping[tname] = table
		// t, err := json.Marshal(&table)
		// fmt.Println("json:", string(jb))
		if err != nil {
			fmt.Printf("\x1b[%dm配置文件: %s 序列化失败\x1b[0m\n", 31, v)
			panic(fmt.Sprintf("\x1b[%dm配置文件: %s 序列化失败\x1b[0m\n", 31, v))
		}

	}

	if ts.IsOutputBuildSchema {
		jb, err := json.Marshal(&buildmapping)
		if err != nil {
			fmt.Printf("\x1b[%dm 序列化编译产物失败 \x1b[0m\n", 31)
			panic(fmt.Sprintf("\x1b[%dm 序列化编译产物失败 \x1b[0m\n", 31))
		}

		bvalue := string(jb)
		if ts.IsEncryOutputBuildSchema {
			bvalue, err = EncryptString(bvalue, []byte(ts.EncryKey))
			if err != nil {
				fmt.Printf("\x1b[%dm 序列化编译产物加密失败: %s \x1b[0m\n", 31, err.Error())
				panic(fmt.Sprintf("\x1b[%dm 序列化编译产物加密失败: %s \x1b[0m\n", 31, err.Error()))
			}
		}

		os.MkdirAll(path.Dir(ts.BuildSchemaDest), os.ModePerm)
		file, err := os.Create(ts.BuildSchemaDest)
		if err != nil {
			fmt.Printf("\x1b[%dm 序列化产物写入失败 \x1b[0m\n", 31)
			panic(fmt.Sprintf("\x1b[%dm 序列化产物写入失败 %s \x1b[0m\n", 31, err.Error()))
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		_, err = fmt.Fprint(writer, bvalue)
		if err != nil {
			fmt.Printf("\x1b[%dm 序列化编译产物写入失败 %s\x1b[0m\n", 31, err.Error())
			panic(fmt.Sprintf("\x1b[%dm 序列化编译产物写入失败 %s \x1b[0m\n", 31, err.Error()))
		}
		writer.Flush()

	}

	return ts
}

func (ts *YamlToSqlHandler) loadFromBuildSchema() *YamlToSqlHandler {

	bvalue, err := os.ReadFile(ts.BuildSchemaDest)
	if err != nil {
		fmt.Printf("\x1b[%dm 序列化编译产物读取失败 \x1b[0m\n", 31)
		panic(fmt.Sprintf("\x1b[%dm 序列化编译产物读取失败 \x1b[0m\n", 31))
	}

	bvaluestr := string(bvalue)
	if ts.IsEncryOutputBuildSchema {
		bvaluestr, err = DecryptString(bvaluestr, []byte(ts.EncryKey))
		if err != nil {
			fmt.Printf("\x1b[%dm 序列化编译产物解密失败 \x1b[0m\n", 31)
			panic(fmt.Sprintf("\x1b[%dm 序列化编译产物解密失败 \x1b[0m\n", 31))
		}
	}
	gjson.Parse(bvaluestr).ForEach(func(key, value gjson.Result) bool {
		sharding_tables := value.Get("Table.sharding_tables").String()
		_value := value.String()
		if sharding_tables != "" {
			for _, sharding_name := range strings.Split(sharding_tables, ",") {
				if sharding_name == "" {
					continue
				}
				sharding_tblv, _ := sjson.Set(_value, "Table.table", sharding_name)
				ts.tables = append(ts.tables, sharding_tblv)
			}
		} else {
			ts.tables = append(ts.tables, value.String())
		}
		return true
	})

	return ts
}

func (ts *YamlToSqlHandler) doSchema() *YamlToSqlHandler {
	// fmt.Println(ts.tables)

	for i, tbl := range ts.tables {
		// sql := ""
		tbJson := gjson.Get(tbl, "Table")
		if tbJson.Exists() {
			tname := tbJson.Get("table")
			// fmt.Println(tname)
			if tname.Exists() {
				var sqlTbl information_schema.SqlTable
				ts.db.Table("INFORMATION_SCHEMA.TABLES").
					Select("*").
					Where("TABLE_SCHEMA=database()").
					Where("TABLE_NAME=?", tname.String()).Find(&sqlTbl)
				// fmt.Println(sqlTbl)
				//数据库里没有这张表
				if sqlTbl.TableName == "" {
					create := ts.getCreateTableSql(tbJson)
					// sql = fmt.Sprintf("%s;\n%s", sql, create)
					ts.sql = append(ts.sql, create)
					// fmt.Println(create)
				} else {
					change := ts.getGetChangeTableSql(tbJson, sqlTbl)
					ts.sql = append(ts.sql, change)
					// fmt.Println(change)
				}
			} else {
				fmt.Printf("\x1b[%dm 文件: %s 缺少表名\x1b[0m\n", 31, ts.yamlFileFullPaths[i])
				panic("缺少表名")
			}

		} else {
			fmt.Printf("\x1b[%dm 文件: %s 不正确\x1b[0m\n", 31, ts.yamlFileFullPaths[i])
			panic("配置文件不正确")
		}

	}

	return ts
}

func (ts *YamlToSqlHandler) getCreateTableSql(tbl gjson.Result) string {

	tname := tbl.Get("table")

	// sql := fmt.Sprintf("CREATE TABLE %s()", tname.String())
	sql := ""

	createPrefix := fmt.Sprintf("CREATE TABLE %s(\n", tname.String())
	if tbl.Get("options.charset").String() == "" {
		fmt.Printf("\x1b[%dm 表: %s charset 不正确\x1b[0m\n", 31, tname)
		panic("配置文件不正确")
	}
	if tbl.Get("options.collate").String() == "" {
		fmt.Printf("\x1b[%dm 表: %s collate 不正确\x1b[0m\n", 31, tname)
		panic("配置文件不正确")
	}
	createSuffix := fmt.Sprintf(")\nDEFAULT CHARACTER SET %s COLLATE %s ENGINE = InnoDB ",
		tbl.Get("options.charset").String(),
		tbl.Get("options.collate").String(),
	)
	if tbl.Get("options.comment").String() != "" {
		createSuffix = fmt.Sprintf("%s COMMENT = '%s' ;",
			createSuffix,
			tbl.Get("options.comment").String(),
		)
	}

	// var Ids []information_schema.TalbeIdInfo
	// gjson.ForEachLine(tbl.Get("id").String(), func(line gjson.Result) bool {
	// 	fmt.Println(line)
	// 	return true
	// })
	// fmt.Println(tbl.Get("id.id.type").IsObject())
	// fmt.Println(tbl.Get("id.id").IsObject())
	// fmt.Println(tbl.Get("id").IsObject())

	var noId bool
	// var noUniqueIndex bool
	// var noIndex bool

	columns := ""
	primary_key := "PRIMARY KEY("
	if tbl.Get("id").IsObject() {
		tbl.Get("id").ForEach(func(key, value gjson.Result) bool {
			if value.IsObject() {
				columnType := value.Get("type").String()
				if columnType == "varchar" {
					columnType = "varchar(255)"
				}

				columns = fmt.Sprintf("%s\t%s %s %s NOT NULL,\n",
					columns,
					key.String(),
					columnType,
					value.Get("generator").String(),
				)
				primary_key = fmt.Sprintf("%s %s ,", primary_key, key.String())
			} else {
				noId = true
			}
			// fmt.Println("key=", key, "====value=", value)
			return true
		})
		primary_key = fmt.Sprintf("%s)", primary_key[:len(primary_key)-1])
	} else {
		noId = true
	}

	if tbl.Get("fields").IsObject() {
		tbl.Get("fields").ForEach(func(key, value gjson.Result) bool {
			if value.IsObject() {

				def := value.Get("default").String()
				generator := value.Get("generator").String()
				comment := value.Get("comment")
				columnType := value.Get("type").String()
				if columnType == "varchar" {
					columnType = "varchar(255)"
				}
				if value.Get("nullable").Bool() {
					//不能有默认值的或者不写默认值
					if !value.Get("default").Exists() || isNoDefaultType(columnType) {
						if isNoDefaultType(columnType) {
							columns = fmt.Sprintf("%s\t%s %s COMMENT '%s' ,\n",
								columns,
								key.String(),
								columnType,
								comment,
							)
						} else {
							columns = fmt.Sprintf("%s\t%s %s DEFAULT NULL %s COMMENT '%s' ,\n",
								columns,
								key.String(),
								columnType,
								generator,
								comment,
							)
						}

					} else {
						columns = fmt.Sprintf("%s\t%s %s DEFAULT '%s' %s COMMENT '%s' ,\n",
							columns,
							key.String(),
							columnType,
							def,
							generator,
							comment,
						)
					}

				} else {
					if !value.Get("default").Exists() || isNoDefaultType(columnType) {
						if isNoDefaultType(columnType) {
							columns = fmt.Sprintf("%s\t%s %s NOT NULL COMMENT '%s' ,\n",
								columns,
								key.String(),
								columnType,
								comment,
							)
						} else {
							columns = fmt.Sprintf("%s\t%s %s NOT NULL %s COMMENT '%s' ,\n",
								columns,
								key.String(),
								columnType,
								generator,
								comment,
							)
						}

					} else {
						columns = fmt.Sprintf("%s\t%s %s DEFAULT '%s'  NOT NULL %s COMMENT '%s' ,\n",
							columns,
							key.String(),
							columnType,
							def,
							generator,
							comment,
						)
					}

				}

			} else {
				fmt.Printf("\x1b[%dm 表: %s 的 %s 不正确\x1b[0m\n", 31, tname, key.String())
				panic("配置文件不正确")
			}

			return true
		})

	}

	if tbl.Get("indexes").IsObject() {
		tbl.Get("indexes").ForEach(func(key, value gjson.Result) bool {
			if value.Get("columns").IsArray() {
				indexColumns := value.Get("columns").Array()
				indexKeys := ""
				for _, ic := range indexColumns {
					indexKeys = fmt.Sprintf("%s%s,", indexKeys, ic.String())
				}
				indexKeys = indexKeys[:len(indexKeys)-1]
				columns = fmt.Sprintf("%s\tINDEX %s (%s),\n",
					columns,
					key.String(),
					indexKeys,
				)
			}

			return true
		})
	}

	if tbl.Get("unique_indexes").IsObject() {
		tbl.Get("unique_indexes").ForEach(func(key, value gjson.Result) bool {
			if value.Get("columns").IsArray() {
				indexColumns := value.Get("columns").Array()
				indexKeys := ""
				for _, ic := range indexColumns {
					indexKeys = fmt.Sprintf("%s%s,", indexKeys, ic.String())
				}
				indexKeys = indexKeys[:len(indexKeys)-1]
				columns = fmt.Sprintf("%s\tUNIQUE INDEX %s (%s),\n",
					columns,
					key.String(),
					indexKeys,
				)
			}

			return true
		})
	}
	if tbl.Get("fulltext_indexes").IsObject() {
		tbl.Get("fulltext_indexes").ForEach(func(key, value gjson.Result) bool {
			if value.Get("columns").IsArray() {
				indexColumns := value.Get("columns").Array()
				indexKeys := ""
				for _, ic := range indexColumns {
					indexKeys = fmt.Sprintf("%s%s,", indexKeys, ic.String())
				}
				indexKeys = indexKeys[:len(indexKeys)-1]
				columns = fmt.Sprintf("%s\tFULLTEXT INDEX %s (%s),\n",
					columns,
					key.String(),
					indexKeys,
				)
			}

			return true
		})
	}

	if !noId {
		columns = fmt.Sprintf("%s\t%s\n",
			columns,
			primary_key,
		)
	} else {
		columns = columns[:len(columns)-2]
	}
	sql = fmt.Sprintf("%s%s%s\n", createPrefix, columns, createSuffix)
	return sql
}

func isNoDefaultType(dataType string) bool {
	switch dataType {
	case "tinytext", "mediumtext", "text", "longtext", "blob", "tinyblob",
		"mediumblob", "longblob":
		return true
	}
	return false
}

func verifyDataType(dataType string) int {
	switch dataType {
	case "int", "integer", "tinyint", "smallint", "mediumint", "bigint",
		"int unsigned", "integer unsigned", "tinyint unsigned", "smallint unsigned",
		"mediumint unsigned", "bigint unsigned", "bit":
		return 1
	case "float", "double", "decimal":
		return 2
	case "bool":
		return 3
	case "enum", "set", "varchar", "char", "tinytext", "mediumtext", "text", "longtext", "blob", "tinyblob",
		"mediumblob", "longblob", "binary", "varbinary":
		return 4
	case "date", "datetime", "timestamp", "time":
		return 5

		//不能有默认值的类型
	}

	return 0
}

func (ts *YamlToSqlHandler) doSqlSafe() *YamlToSqlHandler {
	// fmt.Println("您将要执行的结构操作为：")
	fmt.Printf("\x1b[%dm您将要执行的结构操作为： \x1b[0m\n", 34)
	for k, v := range ts.sql {
		vv := strings.ReplaceAll(v, "\n", "")
		vv = strings.ReplaceAll(vv, " ", "")
		if vv == "" {
			continue
		}
		fmt.Println(">>>>>>>>>>>>>", gjson.Parse(ts.tables[k]).Get("Table.table").String(), ">>>>>>>>>>>>>")
		fmt.Printf("\x1b[%dm%s \x1b[0m\n", 33, v)
		fmt.Println("<<<<<<<<<<<<<", gjson.Parse(ts.tables[k]).Get("Table.table").String(), "<<<<<<<<<<<<<")
	}
	fmt.Printf("\x1b[%dm>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>： \x1b[0m\n", 34)
	fmt.Printf("\x1b[%dm确认执行请输入[ Y ]： \x1b[0m\n", 34)
	fmt.Printf("\x1b[%dm>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>： \x1b[0m\n", 34)
	commend := ""
	fmt.Scanln(&commend)
	if commend == "Y" || commend == "y" {

		// ts.db.Exec("sql")
		// tx := ts.db.Begin()
		errsql := ""
		err := ts.db.Transaction(func(tx *gorm.DB) error {
			for _, tsql := range ts.sql {
				// fmt.Println(">>>>>>>>>>>>>", ts.yamlFileFullPaths[k], ">>>>>>>>>>>>>")
				// fmt.Printf("\x1b[%dm正在执行sql:\n%s \x1b[0m\n", 34, v)
				vv := strings.ReplaceAll(tsql, "\n", "")
				vv = strings.ReplaceAll(vv, " ", "")
				if vv == "" {
					continue
				}

				subsqls := strings.Split(tsql, ";")
				for _, subsql := range subsqls {

					ss := strings.ReplaceAll(subsql, ";", "")
					ss = strings.ReplaceAll(ss, " ", "")
					ss = strings.ReplaceAll(ss, "\n", "")
					if ss == "" {
						continue
					}
					fmt.Printf("\x1b[%dm>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>： \x1b[0m\n", 34)
					fmt.Printf("\x1b[%dm正在执行sql:\n%s \x1b[0m\n", 34,
						subsql+";")
					fmt.Printf("\x1b[%dm>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>： \x1b[0m\n", 34)
					subsql += ";"
					err := tx.Exec(subsql).Error
					if err != nil {
						tx.Rollback()
						errsql = subsql
						return err
					}
				}
				// fmt.Println("<<<<<<<<<<<<<", ts.yamlFileFullPaths[k], "<<<<<<<<<<<<<")
			}
			// tx.Commit()
			return nil
		})
		if err != nil {
			fmt.Printf("\x1b[%dm执行sql:\n%s\n时出现错误 \x1b[0m\n", 31, errsql)
			panic(err)
		}
		fmt.Printf("\x1b[%dm>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>： \x1b[0m\n", 34)
		fmt.Printf("\x1b[%dmSQL更新完毕： \x1b[0m\n", 36)
		fmt.Printf("\x1b[%dm<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<： \x1b[0m\n", 34)
	}

	return ts
}

func (ts *YamlToSqlHandler) doSql() *YamlToSqlHandler {
	// fmt.Println("您将要执行的结构操作为：")
	fmt.Printf("\x1b[%dm您将要执行的结构操作为： \x1b[0m\n", 34)
	for k, v := range ts.sql {
		vv := strings.ReplaceAll(v, "\n", "")
		vv = strings.ReplaceAll(vv, " ", "")
		if vv == "" {
			continue
		}
		fmt.Println(">>>>>>>>>>>>>", ts.yamlFileFullPaths[k], ">>>>>>>>>>>>>")
		fmt.Printf("\x1b[%dm%s \x1b[0m\n", 33, v)
		fmt.Println("<<<<<<<<<<<<<", ts.yamlFileFullPaths[k], "<<<<<<<<<<<<<")
	}
	fmt.Printf("\x1b[%dm>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>： \x1b[0m\n", 34)
	fmt.Printf("\x1b[%dm确认执行请输入[ Y ]： \x1b[0m\n", 34)
	fmt.Printf("\x1b[%dm>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>： \x1b[0m\n", 34)

	errsql := ""
	err := ts.db.Transaction(func(tx *gorm.DB) error {
		for _, tsql := range ts.sql {

			vv := strings.ReplaceAll(tsql, "\n", "")
			vv = strings.ReplaceAll(vv, " ", "")
			if vv == "" {
				continue
			}

			subsqls := strings.Split(tsql, ";")
			for _, subsql := range subsqls {

				ss := strings.ReplaceAll(subsql, ";", "")
				ss = strings.ReplaceAll(ss, " ", "")
				ss = strings.ReplaceAll(ss, "\n", "")
				if ss == "" {
					continue
				}
				fmt.Printf("\x1b[%dm>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>： \x1b[0m\n", 34)
				fmt.Printf("\x1b[%dm正在执行sql:\n%s \x1b[0m\n", 34,
					subsql+";")
				fmt.Printf("\x1b[%dm>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>： \x1b[0m\n", 34)
				subsql += ";"
				err := tx.Exec(subsql).Error
				if err != nil {
					tx.Rollback()
					errsql = subsql
					return err
				}
			}
			// fmt.Println("<<<<<<<<<<<<<", ts.yamlFileFullPaths[k], "<<<<<<<<<<<<<")
		}
		// tx.Commit()
		return nil
	})
	if err != nil {
		fmt.Printf("\x1b[%dm执行sql:\n%s\n时出现错误 \x1b[0m\n", 31, errsql)
		panic(err)
	}

	fmt.Printf("\x1b[%dm>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>： \x1b[0m\n", 34)
	fmt.Printf("\x1b[%dmSQL更新完毕： \x1b[0m\n", 36)
	fmt.Printf("\x1b[%dm<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<： \x1b[0m\n", 34)

	return ts
}

func (ts *YamlToSqlHandler) getGetChangeTableSql(tbl gjson.Result, sqlTbl information_schema.SqlTable) string {
	tname := tbl.Get("table").String()
	sql := "\n"
	if sqlTbl.TableComment != tbl.Get("options.comment").String() {
		sql = fmt.Sprintf("%sALTER TABLE %s comment '%s';\n", sql, tname, tbl.Get("options.comment").String())
	}
	//行
	//计算sql行
	var sqlColumns []information_schema.SqlTableColumns
	ts.db.Table("`INFORMATION_SCHEMA`.`COLUMNS`").
		Where("TABLE_SCHEMA=database()").
		Where("TABLE_NAME=?", tname).
		Find(&sqlColumns)
	var sqlColumnsSerialize = information_schema.SqlColumnsSerialize{}
	// sqlColumnsSerialize = map[string]map[string]string{}
	for _, sc := range sqlColumns {
		if strings.ToLower(sc.ColumnName) == "id" {
			continue
		}
		if sqlColumnsSerialize[sc.ColumnName] == nil {
			sqlColumnsSerialize[sc.ColumnName] = map[string]string{}
		}
		if !strings.Contains(sc.ColumnType, "varchar") && strings.Contains(sc.ColumnType, "(") {
			b := strings.Index(sc.ColumnType, "(")
			e := strings.Index(sc.ColumnType, ")")
			sc.ColumnType = sc.ColumnType[:b] + sc.ColumnType[e+1:]
		}
		sqlColumnsSerialize[sc.ColumnName]["type"] = sc.ColumnType

		if strings.ToLower(sc.IsNullable) == "yes" {
			sqlColumnsSerialize[sc.ColumnName]["nullable"] = "true"
		} else {
			sqlColumnsSerialize[sc.ColumnName]["nullable"] = "false"
		}

		sqlColumnsSerialize[sc.ColumnName]["comment"] = sc.ColumnComment
		if sc.ColumnDefault == nil {
			// sqlColumnsSerialize[sc.ColumnName]["default"] = ""
		} else {
			sqlColumnsSerialize[sc.ColumnName]["default"] = *sc.ColumnDefault
		}
		sqlColumnsSerialize[sc.ColumnName]["generator"] = sc.Extra
	}
	sqlColumnsJ, err := json.Marshal(&sqlColumnsSerialize)
	if err != nil {
		fmt.Printf("\x1b[%dm 表: %s 序列化失败 \x1b[0m\n", 31, tname)
		panic("配置文件不正确")
	}
	//计算删除和修改
	// var dropCloumns map[string]bool = map[string]bool{}
	// dropCloumns
	dropColumnsSql := ""
	sqlColumnsgj := gjson.Parse(string(sqlColumnsJ))
	sqlColumnsgj.ForEach(func(key, value gjson.Result) bool {
		if tbl.Get("fields." + key.String()).Exists() {
			var refresh bool
			var yy string
			//类型
			if tbl.Get("fields." + key.String() + ".type").Exists() {
				ymlt := tbl.Get("fields." + key.String() + ".type").String()
				ymlt = strings.ReplaceAll(ymlt, "integer", "int")
				if ymlt == "varchar" {
					ymlt = "varchar(255)"
				}
				sqlt := value.Get("type").String()
				sqlt = strings.ReplaceAll(sqlt, "integer", "int")

				if strings.Compare(strings.ToLower(sqlt),
					strings.ToLower(ymlt)) != 0 {
					refresh = true
					yy += "类型/"
				}
			} else {
				fmt.Printf("\x1b[%dm 表: %s 配置文件不正确 \x1b[0m\n", 31, tname)
				panic("配置文件不正确")
			}

			//nullable
			if tbl.Get("fields." + key.String() + ".nullable").Exists() {
				if strings.Compare(strings.ToLower(value.Get("nullable").String()),
					strings.ToLower(tbl.Get("fields."+key.String()+".nullable").String())) != 0 {
					refresh = true
					yy += "空不空/"
				}
			} else {
				if strings.ToLower(value.Get("nullable").String()) != "false" {
					refresh = true
					yy += "空不空/"
				}
			}

			//comment
			if tbl.Get("fields." + key.String() + ".comment").Exists() {
				if strings.Compare(value.Get("comment").String(),
					tbl.Get("fields."+key.String()+".comment").String()) != 0 {
					refresh = true
					yy += "备注/"
				}
			} else {
				if value.Get("comment").String() != "" {
					refresh = true
					yy += "备注/"
				}
			}

			//default
			//判定是否为自动插入数据
			//有些数据库类型不允许有默认值
			if !isNoDefaultType(value.Get("type").String()) {
				if strings.ToLower(value.Get("default").String()) == "current_timestamp" &&
					strings.ToLower(value.Get("generator").String()) == "default_generated" {
					if tbl.Get("fields." + key.String() + ".generator").Exists() {
						if strings.ToLower(tbl.Get("fields."+key.String()+".generator").String()) !=
							"default current_timestamp" {
							refresh = true
							yy += "默认/"
						}
					} else {
						// fmt.Println(tbl)
						// fmt.Println(tbl.Get("fields." + key.String() + ".generator"))
						refresh = true
						yy += "默认/"
					}
				} else {
					if tbl.Get("fields." + key.String() + ".default").Exists() {
						if value.Get("default").Exists() {
							if strings.Compare(value.Get("default").String(),
								tbl.Get("fields."+key.String()+".default").String()) != 0 {
								refresh = true
								yy += "默认/"
							}
						} else {
							refresh = true
							yy += "默认/"
						}
					} else {
						if value.Get("default").String() != "" {
							refresh = true
							yy += "默认/"
						}
					}
				}
			}

			//generator
			sqlg := strings.ToLower(value.Get("generator").String())
			if tbl.Get("fields." + key.String() + ".generator").Exists() {
				// if !strings.Contains(sqlg, "default_generated") {
				// 	if strings.Compare(sqlg,
				// 		strings.ToLower(tbl.Get("fields."+key.String()+".generator").String())) != 0 {
				// 		refresh = true
				// 		yy += "自动/"
				// 	}
				// }

			} else {
				if value.Get("generator").String() != "" && !strings.Contains(sqlg, "default_generated") {
					refresh = true
					yy += "自动/"
				}
			}
			if refresh {
				ymlt := tbl.Get("fields." + key.String() + ".type").String()
				if ymlt == "varchar" {
					ymlt = "varchar(255)"
				}
				if isNoDefaultType(ymlt) {
					if tbl.Get("fields."+key.String()+".nullable").Exists() &&
						tbl.Get("fields."+key.String()+".nullable").String() == "true" {
						sql = fmt.Sprintf("%sALTER TABLE %s MODIFY COLUMN %s %s %s COMMENT '%s';\n",
							sql,
							tname,
							key.String(),
							ymlt,
							tbl.Get("fields."+key.String()+".generator").String(),
							tbl.Get("fields."+key.String()+".comment").String(),
						)
					} else {

						sql = fmt.Sprintf("%sALTER TABLE %s MODIFY COLUMN %s %s NOT NULL %s COMMENT '%s';\n",
							sql,
							tname,
							key.String(),
							ymlt,
							tbl.Get("fields."+key.String()+".generator").String(),
							tbl.Get("fields."+key.String()+".comment").String(),
						)
					}

				} else {

					if tbl.Get("fields."+key.String()+".nullable").Exists() &&
						tbl.Get("fields."+key.String()+".nullable").String() == "true" {
						if tbl.Get("fields." + key.String() + ".default").Exists() {
							sql = fmt.Sprintf("%sALTER TABLE %s MODIFY COLUMN %s %s %s DEFAULT '%s' COMMENT '%s';\n",
								sql,
								tname,
								key.String(),
								ymlt,
								tbl.Get("fields."+key.String()+".generator").String(),
								tbl.Get("fields."+key.String()+".default").String(),
								tbl.Get("fields."+key.String()+".comment").String(),
							)
						} else {
							sql = fmt.Sprintf("%sALTER TABLE %s MODIFY COLUMN %s %s %s DEFAULT NULL COMMENT '%s';\n",
								sql,
								tname,
								key.String(),
								ymlt,
								tbl.Get("fields."+key.String()+".generator").String(),
								// tbl.Get("fields."+key.String()+".default").String(),
								tbl.Get("fields."+key.String()+".comment").String(),
							)
						}

					} else {

						if tbl.Get("fields." + key.String() + ".default").Exists() {
							sql = fmt.Sprintf("%sALTER TABLE %s MODIFY COLUMN %s %s NOT NULL %s DEFAULT '%s' COMMENT '%s';\n",
								sql,
								tname,
								key.String(),
								ymlt,
								tbl.Get("fields."+key.String()+".generator").String(),
								tbl.Get("fields."+key.String()+".default").String(),
								tbl.Get("fields."+key.String()+".comment").String(),
							)
						} else {
							sql = fmt.Sprintf("%sALTER TABLE %s MODIFY COLUMN %s %s NOT NULL %s COMMENT '%s';\n",
								sql,
								tname,
								key.String(),
								ymlt,
								tbl.Get("fields."+key.String()+".generator").String(),
								// tbl.Get("fields."+key.String()+".default").String(),
								tbl.Get("fields."+key.String()+".comment").String(),
							)
						}

					}

				}

				// fmt.Println(">>>>>>更新行==")
				// fmt.Println("库= ", key.String())
				// fmt.Println("yml= ", tbl.Get("fields."+key.String()).String())
				// fmt.Println("sql= ", value.String())
				// fmt.Println("原因：", yy)
				// fmt.Println("<<<<<===")
			} else {

				// fmt.Println(">>>>>>保留行==")
				// fmt.Println("库= ", key.String())
				// fmt.Println("yml= ", tbl.Get("fields."+key.String()).String())
				// fmt.Println("sql= ", value.String())
				// fmt.Println("原因：", yy)
				// fmt.Println("<<<<<===")
			}

		} else if tbl.Get("id." + key.String()).Exists() {
			//主键区，暂时用不上
		} else {
			dropColumnsSql = fmt.Sprintf("%sALTER  TABLE %s DROP %s;\n",
				dropColumnsSql,
				tname,
				key.String(),
			)
		}

		return true
	})
	//计算新增
	tbl.Get("fields").ForEach(func(key, value gjson.Result) bool {

		if !sqlColumnsgj.Get(key.String()).Exists() {

			ymlt := value.Get("type").String()
			if ymlt == "varchar" {
				ymlt = "varchar(255)"
			}
			if isNoDefaultType(ymlt) {
				if value.Get("nullable").Exists() &&
					value.Get("nullable").String() != "true" {
					sql = fmt.Sprintf("%sALTER TABLE %s ADD COLUMN %s %s NOT NULL %s COMMENT '%s';\n",
						sql,
						tname,
						key.String(),
						ymlt,
						tbl.Get("fields."+key.String()+".generator").String(),
						tbl.Get("fields."+key.String()+".comment").String(),
					)
				} else {
					sql = fmt.Sprintf("%sALTER TABLE %s ADD COLUMN %s %s %s COMMENT '%s';\n",
						sql,
						tname,
						key.String(),
						ymlt,
						tbl.Get("fields."+key.String()+".generator").String(),
						tbl.Get("fields."+key.String()+".comment").String(),
					)
				}

			} else {

				if value.Get("nullable").Exists() &&
					value.Get("nullable").String() != "true" {
					if value.Get("default").Exists() {
						sql = fmt.Sprintf("%sALTER TABLE %s ADD COLUMN %s %s NOT NULL %s DEFAULT '%s' COMMENT '%s';\n",
							sql,
							tname,
							key.String(),
							ymlt,
							tbl.Get("fields."+key.String()+".generator").String(),
							tbl.Get("fields."+key.String()+".default").String(),
							tbl.Get("fields."+key.String()+".comment").String(),
						)
					} else {
						sql = fmt.Sprintf("%sALTER TABLE %s ADD COLUMN %s %s NOT NULL %s COMMENT '%s';\n",
							sql,
							tname,
							key.String(),
							ymlt,
							tbl.Get("fields."+key.String()+".generator").String(),
							// tbl.Get("fields."+key.String()+".default").String(),
							tbl.Get("fields."+key.String()+".comment").String(),
						)
					}

				} else {
					if value.Get("default").Exists() {
						sql = fmt.Sprintf("%sALTER TABLE %s ADD COLUMN %s %s %s DEFAULT '%s' COMMENT '%s';\n",
							sql,
							tname,
							key.String(),
							ymlt,
							tbl.Get("fields."+key.String()+".generator").String(),
							tbl.Get("fields."+key.String()+".default").String(),
							tbl.Get("fields."+key.String()+".comment").String(),
						)
					} else {
						sql = fmt.Sprintf("%sALTER TABLE %s ADD COLUMN %s %s %s DEFAULT NULL COMMENT '%s';\n",
							sql,
							tname,
							key.String(),
							ymlt,
							tbl.Get("fields."+key.String()+".generator").String(),
							// tbl.Get("fields."+key.String()+".default").String(),
							tbl.Get("fields."+key.String()+".comment").String(),
						)
					}

				}

			}

		}

		return true
	})
	// fmt.Println(string(sqlColumnsJ))
	///
	//索引
	var sqlIndexes []information_schema.SqlIndexes
	ts.db.Raw(fmt.Sprintf("show indexes from %s", tname)).Scan(&sqlIndexes)

	var sqlIndexesSerialize information_schema.SqlIndexesSerialize
	sqlIndexesSerialize.UnqIndexes = map[string]map[string][]string{}
	sqlIndexesSerialize.Indexes = map[string]map[string][]string{}
	sqlIndexesSerialize.FulltextIndexes = map[string]map[string][]string{}

	//计算sql
	for _, sqlIndex := range sqlIndexes {
		if strings.ToLower(sqlIndex.Key_name) == "primary" {
			continue
		}
		if strings.ToLower(sqlIndex.IndexType) == "fulltext" {
			if sqlIndexesSerialize.FulltextIndexes[sqlIndex.Key_name] == nil {
				sqlIndexesSerialize.FulltextIndexes[sqlIndex.Key_name] = map[string][]string{}
			}
			sqlIndexesSerialize.FulltextIndexes[sqlIndex.Key_name]["columns"] = append(sqlIndexesSerialize.FulltextIndexes[sqlIndex.Key_name]["columns"], sqlIndex.Column_name)
		} else if strings.ToLower(sqlIndex.IndexType) == "btree" {
			if sqlIndex.Non_unique == 0 {
				// fmt.Println(sqlIndex.Key_name)
				if sqlIndexesSerialize.UnqIndexes[sqlIndex.Key_name] == nil {
					sqlIndexesSerialize.UnqIndexes[sqlIndex.Key_name] = map[string][]string{}
				}
				sqlIndexesSerialize.UnqIndexes[sqlIndex.Key_name]["columns"] = append(sqlIndexesSerialize.UnqIndexes[sqlIndex.Key_name]["columns"], sqlIndex.Column_name)
			}
			if sqlIndex.Non_unique == 1 {
				// fmt.Println(sqlIndex.Key_name)
				if sqlIndexesSerialize.Indexes[sqlIndex.Key_name] == nil {
					sqlIndexesSerialize.Indexes[sqlIndex.Key_name] = map[string][]string{}
				}
				sqlIndexesSerialize.Indexes[sqlIndex.Key_name]["columns"] = append(sqlIndexesSerialize.Indexes[sqlIndex.Key_name]["columns"], sqlIndex.Column_name)
			}
		}
	}

	sqlIndexesJ, err := json.Marshal(&sqlIndexesSerialize)
	if err != nil {
		fmt.Printf("\x1b[%dm 表: %s 序列化失败 \x1b[0m\n", 31, tname)
		panic("配置文件不正确")
	}
	// fmt.Println(string(sqlIndexesJ))
	// var keepThisKey bool
	//计算删除+修改
	// var dropIndexes map[string]bool = map[string]bool{}
	dropIndexesSql := ""
	if gjson.Get(string(sqlIndexesJ), "unique_indexes").Exists() {
		gjson.Get(string(sqlIndexesJ), "unique_indexes").ForEach(func(key, value gjson.Result) bool {
			if tbl.Get("unique_indexes." + key.String()).Exists() {
				if strings.Compare(value.String(), tbl.Get("unique_indexes."+key.String()).String()) == 0 {
					// keepThisKey = true
					// fmt.Println(">>>>>>保留==")
					// fmt.Println("unqkey= ", key.String())
					// fmt.Println("yml= ", tbl.Get("unique_indexes."+key.String()).String())
					// fmt.Println("sql= ", value.String())
					// fmt.Println("<<<<<===")
				} else {
					sql = fmt.Sprintf("%sDROP INDEX %s ON %s;\n",
						sql,
						key.String(),
						tname,
					)
					sqlcol := ""
					for _, v := range tbl.Get("unique_indexes." + key.String() + ".columns").Array() {
						sqlcol = fmt.Sprintf("%s%s,", sqlcol, v)
					}
					sqlcol = sqlcol[:len(sqlcol)-1]
					sql = fmt.Sprintf("%sCREATE UNIQUE INDEX %s ON %s(%s);\n",
						sql,
						key.String(),
						tname,
						sqlcol,
					)
					// keepThisKey = false
					// needReaddKey[key.String()] = true
				}
			} else {
				dropIndexesSql = fmt.Sprintf("%sDROP INDEX %s ON %s;\n",
					dropIndexesSql,
					key.String(),
					tname,
				)
				// if value.Get("columns")
				// keepThisKey = false
			}
			return true
		})
	}
	if gjson.Get(string(sqlIndexesJ), "fulltext_indexes").Exists() {
		gjson.Get(string(sqlIndexesJ), "fulltext_indexes").ForEach(func(key, value gjson.Result) bool {
			if tbl.Get("fulltext_indexes." + key.String()).Exists() {
				// fmt.Println(tbl.Get("fulltext_indexes." + key.String() + ".columns").String())
				// fmt.Println(tbl.Get("fulltext_indexes." + key.String() + ".with_parser").String())
				if strings.Compare(value.Get("columns").String(), tbl.Get("fulltext_indexes."+key.String()+".columns").String()) == 0 {
					// keepThisKey = true
					// fmt.Println(">>>>>>保留==")
					// fmt.Println("unqkey= ", key.String())
					// fmt.Println("yml= ", tbl.Get("fulltext_indexes."+key.String()).String())
					// fmt.Println("sql= ", value.String())
					// fmt.Println("<<<<<===")
				} else {
					sql = fmt.Sprintf("%sDROP INDEX %s ON %s;\n",
						sql,
						key.String(),
						tname,
					)
					sqlcol := ""
					for _, v := range tbl.Get("fulltext_indexes." + key.String() + ".columns").Array() {
						sqlcol = fmt.Sprintf("%s%s,", sqlcol, v)
					}
					sqlcol = sqlcol[:len(sqlcol)-1]
					with_parser := tbl.Get("fulltext_indexes." + key.String() + ".with_parser").String()
					if with_parser == "" {
						with_parser = "ngram"
					}
					sql = fmt.Sprintf("%sCREATE FULLTEXT INDEX %s ON %s(%s) WITH PARSER %s;\n",
						sql,
						key.String(),
						tname,
						sqlcol,
						with_parser,
					)
					// keepThisKey = false
					// needReaddKey[key.String()] = true
				}
			} else {
				dropIndexesSql = fmt.Sprintf("%sDROP INDEX %s ON %s;\n",
					dropIndexesSql,
					key.String(),
					tname,
				)
				// if value.Get("columns")
				// keepThisKey = false
			}
			return true
		})
	}
	if gjson.Get(string(sqlIndexesJ), "indexes").Exists() {
		gjson.Get(string(sqlIndexesJ), "indexes").ForEach(func(key, value gjson.Result) bool {
			if tbl.Get("indexes." + key.String()).Exists() {
				if strings.Compare(value.String(), tbl.Get("indexes."+key.String()).String()) == 0 {

					// fmt.Println(">>>>>>保留==")
					// fmt.Println("unqkey= ", key.String())
					// fmt.Println("yml= ", tbl.Get("indexes."+key.String()).String())
					// fmt.Println("sql= ", value.String())
					// fmt.Println("<<<<<===")
				} else {
					sql = fmt.Sprintf("%sDROP INDEX %s ON %s;\n",
						sql,
						key.String(),
						tname,
					)
					sqlcol := ""
					for _, v := range tbl.Get("indexes." + key.String() + ".columns").Array() {
						sqlcol = fmt.Sprintf("%s%s,", sqlcol, v)
					}
					sqlcol = sqlcol[:len(sqlcol)-1]
					sql = fmt.Sprintf("%sCREATE INDEX %s ON %s(%s);\n",
						sql,
						key.String(),
						tname,
						sqlcol,
					)

					// fmt.Println(">>>>>>修改==")
					// fmt.Println("unqkey= ", key.String())
					// fmt.Println("yml= ", tbl.Get("indexes."+key.String()).String())
					// fmt.Println("sql= ", value.String())
					// fmt.Println("<<<<<===")
				}
			} else {
				dropIndexesSql = fmt.Sprintf("%sDROP INDEX %s ON %s;\n",
					dropIndexesSql,
					key.String(),
					tname,
				)
				// fmt.Println(">>>>>>丢弃==")
				// fmt.Println("unqkey= ", key.String())
				// fmt.Println("yml= ", tbl.Get("indexes."+key.String()).String())
				// fmt.Println("sql= ", value.String())
				// fmt.Println("<<<<<===")
			}
			return true
		})
	}
	//计算新增索引
	tbl.Get("unique_indexes").ForEach(func(key, value gjson.Result) bool {
		if !gjson.Get(string(sqlIndexesJ), "unique_indexes."+key.String()).Exists() {
			sqlcol := ""
			for _, v := range value.Get("columns").Array() {
				sqlcol = fmt.Sprintf("%s%s,", sqlcol, v)
			}
			sqlcol = sqlcol[:len(sqlcol)-1]
			sql = fmt.Sprintf("%sCREATE UNIQUE INDEX %s ON %s(%s);\n",
				sql,
				key.String(),
				tname,
				sqlcol,
			)
		}

		return true
	})
	tbl.Get("fulltext_indexes").ForEach(func(key, value gjson.Result) bool {
		if !gjson.Get(string(sqlIndexesJ), "fulltext_indexes."+key.String()).Exists() {
			sqlcol := ""
			for _, v := range value.Get("columns").Array() {
				sqlcol = fmt.Sprintf("%s%s,", sqlcol, v)
			}
			sqlcol = sqlcol[:len(sqlcol)-1]

			with_parser := tbl.Get("fulltext_indexes." + key.String() + ".with_parser").String()
			if with_parser == "" {
				with_parser = "ngram"
			}

			sql = fmt.Sprintf("%sCREATE FULLTEXT INDEX %s ON %s(%s) WITH PARSER %s;\n",
				sql,
				key.String(),
				tname,
				sqlcol,
				with_parser,
			)
		}
		return true
	})
	tbl.Get("indexes").ForEach(func(key, value gjson.Result) bool {

		if !gjson.Get(string(sqlIndexesJ), "indexes."+key.String()).Exists() {
			sqlcol := ""
			for _, v := range value.Get("columns").Array() {
				sqlcol = fmt.Sprintf("%s%s,", sqlcol, v)
			}
			sqlcol = sqlcol[:len(sqlcol)-1]
			sql = fmt.Sprintf("%sCREATE INDEX %s ON %s(%s);\n",
				sql,
				key.String(),
				tname,
				sqlcol,
			)
		}
		return true
	})
	sql = fmt.Sprintf("%s%s%s", sql, dropIndexesSql, dropColumnsSql)

	return sql
}

// 校验yml的合法行
func (ts *YamlToSqlHandler) verifyYmlFile() *YamlToSqlHandler {
	for k, table := range ts.tables {
		tbJson := gjson.Get(table, "Table")
		// fieldsMap := map[string]string{}
		if !tbJson.Get("id").Exists() {
			fmt.Printf("\x1b[%dm 配置文件不正确:'%s' \x1b[0m\n", 31, ts.yamlFileFullPaths[k])
			fmt.Printf("\x1b[%dm 缺少主键id \x1b[0m\n", 31)
			panic("配置文件不正确")
		}
		// if !tbJson.Get("id.id").Exists() {
		// 	fmt.Printf("\x1b[%dm 配置文件不正确:'%s' \x1b[0m\n", 31, ts.yamlFileFullPaths[k])
		// 	fmt.Printf("\x1b[%dm 缺少主键id \x1b[0m\n", 31)
		// 	panic("配置文件不正确")
		// }
		tbJson.Get("indexes").ForEach(func(key, value gjson.Result) bool {
			if !value.Get("columns").IsArray() {
				fmt.Printf("\x1b[%dm 配置文件不正确:'%s' \x1b[0m\n", 31, ts.yamlFileFullPaths[k])
				fmt.Printf("\x1b[%dm indexes:'%s' is not array\x1b[0m\n", 31, key.String())
				panic("配置文件不正确")
			}
			for _, v := range value.Get("columns").Array() {
				if !tbJson.Get("fields."+v.String()).Exists() && !tbJson.Get("id."+v.String()).Exists() {
					fmt.Printf("\x1b[%dm 配置文件不正确:'%s' \x1b[0m\n", 31, ts.yamlFileFullPaths[k])
					fmt.Printf("\x1b[%dm indexes columns:'%s' is not find\x1b[0m\n", 31, v.String())
					panic("配置文件不正确")
				}
			}
			return true
		})
		tbJson.Get("unique_indexes").ForEach(func(key, value gjson.Result) bool {
			if !value.Get("columns").IsArray() {
				fmt.Printf("\x1b[%dm 配置文件不正确:'%s' \x1b[0m\n", 31, ts.yamlFileFullPaths[k])
				fmt.Printf("\x1b[%dm indexes:'%s' is not array\x1b[0m\n", 31, key.String())
				panic("配置文件不正确")
			}
			for _, v := range value.Get("columns").Array() {
				if !tbJson.Get("fields."+v.String()).Exists() && !tbJson.Get("id."+v.String()).Exists() {
					fmt.Printf("\x1b[%dm 配置文件不正确:'%s' \x1b[0m\n", 31, ts.yamlFileFullPaths[k])
					fmt.Printf("\x1b[%dm indexes columns:'%s' is not find\x1b[0m\n", 31, v.String())
					panic("配置文件不正确")
				}
			}
			return true
		})
		tbJson.Get("fulltext_indexes").ForEach(func(key, value gjson.Result) bool {
			if !value.Get("columns").IsArray() {
				fmt.Printf("\x1b[%dm 配置文件不正确:'%s' \x1b[0m\n", 31, ts.yamlFileFullPaths[k])
				fmt.Printf("\x1b[%dm indexes:'%s' is not array\x1b[0m\n", 31, key.String())
				panic("配置文件不正确")
			}
			for _, v := range value.Get("columns").Array() {
				if !tbJson.Get("fields."+v.String()).Exists() && !tbJson.Get("id."+v.String()).Exists() {
					fmt.Printf("\x1b[%dm 配置文件不正确:'%s' \x1b[0m\n", 31, ts.yamlFileFullPaths[k])
					fmt.Printf("\x1b[%dm indexes columns:'%s' is not find\x1b[0m\n", 31, v.String())
					panic("配置文件不正确")
				}
			}
			return true
		})
	}

	return ts
}

func (ts *YamlToSqlHandler) ExecuteSchemaSafeCheck() *YamlToSqlHandler {
	ts.connectSql()
	ts.getyamlFileFullPaths().
		getYamlDatas().verifyYmlFile().doSchema().doSqlSafe()

	return ts
}

func (ts *YamlToSqlHandler) ExecuteSchema() *YamlToSqlHandler {
	ts.connectSql()
	ts.getyamlFileFullPaths().
		getYamlDatas().verifyYmlFile().doSchema().doSql()

	return ts
}

func (ts *YamlToSqlHandler) LoadSchema() *YamlToSqlHandler {
	ts.connectSql()
	ts.loadFromBuildSchema().verifyYmlFile().doSchema()
	return ts
}

func (ts *YamlToSqlHandler) trimSql() *YamlToSqlHandler {

	var newsql []string
	for _, v := range ts.sql {
		vv := strings.ReplaceAll(v, "\n", "")
		vv = strings.ReplaceAll(vv, " ", "")
		if vv == "" {
			continue
		}
		newsql = append(newsql, v)
	}
	ts.sql = newsql
	return ts
}

func (ts *YamlToSqlHandler) VerifyIsCleanSchema() bool {
	ts.trimSql()
	return len(ts.sql) < 1
}

func (ts *YamlToSqlHandler) GetSql() []string {
	return ts.sql
}

func (ts *YamlToSqlHandler) DoSql() *YamlToSqlHandler {
	ts.doSqlSafe()
	return ts
}
