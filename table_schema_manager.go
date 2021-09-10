package tblschema

import (
	"fmt"
	"os"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//mysql数据类型<=>go数据类型
var sqlTypeToGoType = map[string]string{
	"int":                "int",
	"integer":            "int",
	"tinyint":            "int",
	"smallint":           "int",
	"mediumint":          "int",
	"bigint":             "int",
	"int unsigned":       "int",
	"integer unsigned":   "int",
	"tinyint unsigned":   "int",
	"smallint unsigned":  "int",
	"mediumint unsigned": "int",
	"bigint unsigned":    "int",
	"bit":                "int",
	"bool":               "bool",
	"enum":               "string",
	"set":                "string",
	"varchar":            "string",
	"char":               "string",
	"tinytext":           "string",
	"mediumtext":         "string",
	"text":               "string",
	"longtext":           "string",
	"blob":               "string",
	"tinyblob":           "string",
	"mediumblob":         "string",
	"longblob":           "string",
	"float":              "float64",
	"double":             "float64",
	"decimal":            "float64",
	"binary":             "string",
	"varbinary":          "string",

	"date":      "time.Time", //
	"datetime":  "time.Time", //
	"timestamp": "time.Time", //
	"time":      "time.Time", //

	"date_string":      "string", // time.Time
	"datetime_string":  "string", // time.Time
	"timestamp_string": "string", // time.Time
	"time_string":      "string", // time.Time
}

const (
	CAMEL_CASE  = "camel_case"
	FIRST_UPPER = "first_upper"

	ORM  = "orm"
	GORM = "gorm"
)

type TblSchemaHandler struct {
	dsn string   //数据库连接dsn,列：用户:密码@(127.0.0.1:3306)/数据库?charset=utf8mb4&parseTime=True&loc=Local
	db  *gorm.DB //数据库连接

	tableName       string   //要生成model的数据库表名
	modelOrmTagType string   //生成orm结构题标签类型
	otherTag        []string //其他标签
	fieldNameType   string   //字段名类型 CAMEL_CASE骆驼命名/FIRST_UPPER首字母大写
	structNameType  string   //结构体名类型 CAMEL_CASE骆驼命名/FIRST_UPPER首字母大写
	savePath        string   //保存model文件的位置
	packageName     string   //生成model包名

	columns         []column
	maxLenFieldType int
	maxLenFieldTag  int
	maxLenFieldName int
}

func NewTblSchemaHandler() *TblSchemaHandler {
	return &TblSchemaHandler{
		packageName:     "model",
		savePath:        "model.go",
		modelOrmTagType: "gorm",
		fieldNameType:   CAMEL_CASE,
		structNameType:  CAMEL_CASE,
	}
}

func (ts *TblSchemaHandler) SetDsn(dsn string) *TblSchemaHandler {
	ts.dsn = dsn
	return ts
}
func (ts *TblSchemaHandler) SetDB(db *gorm.DB) *TblSchemaHandler {
	ts.db = db
	return ts
}
func (ts *TblSchemaHandler) SetSavePath(savePath string) *TblSchemaHandler {
	ts.savePath = savePath
	return ts
}
func (ts *TblSchemaHandler) SetTableName(tableName string) *TblSchemaHandler {
	ts.tableName = tableName
	return ts
}

func (ts *TblSchemaHandler) SetModelOrmTagType(modelOrmTagType string) *TblSchemaHandler {
	ts.modelOrmTagType = modelOrmTagType
	return ts
}
func (ts *TblSchemaHandler) SetOtherTag(otherTag ...string) *TblSchemaHandler {
	ts.otherTag = otherTag
	return ts
}
func (ts *TblSchemaHandler) SefieldType(fieldNameType string) *TblSchemaHandler {
	ts.fieldNameType = fieldNameType
	return ts
}

func (ts *TblSchemaHandler) SetTableNameType(structNameType string) *TblSchemaHandler {
	ts.structNameType = structNameType
	return ts
}
func (ts *TblSchemaHandler) SetPackageName(packageName string) *TblSchemaHandler {
	ts.packageName = packageName
	return ts
}

func (ts *TblSchemaHandler) Run() {
	ts.connectSql()
	ts.getColumns()
	packageName := "package model"
	if ts.packageName != "" {
		packageName = fmt.Sprintf("package %s\n", ts.packageName)
	}
	structName := ts.camelCase(ts.tableName, ts.structNameType)
	structContent := fmt.Sprintf("type %s struct{\n", structName)

	for _, v := range ts.columns {
		match := fmt.Sprint("\t%-", ts.maxLenFieldName, "s %-", ts.maxLenFieldType, "s %-", ts.maxLenFieldTag, "s %s\n")
		structContent += fmt.Sprintf(match, v.FieldContent.Name, v.FieldContent.Type, v.FieldContent.Tag, v.FieldContent.Comment)
	}
	structContent += "}\n"

	functableName := fmt.Sprintf("func (*%s) TableName() string {\n", structName) +
		fmt.Sprintf("\t return \"%s\"\n", ts.tableName) +
		"}\n"

	fileContent := fmt.Sprintf("%s\n%s\n%s", packageName, structContent, functableName)

	fmt.Println(fileContent)
	filePath := fmt.Sprint(ts.savePath)
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println("打开文件失败")
		return
	}
	defer f.Close()

	f.WriteString(fileContent)

}

type column struct {
	ColumnName    string `gorm:"column:COLUMN_NAME"`
	Type          string `gorm:"column:DATA_TYPE"`
	Nullable      string `gorm:"column:IS_NULLABLE"`
	TableName     string `gorm:"column:TABLE_NAME"`
	ColumnComment string `gorm:"column:COLUMN_COMMENT"`

	FieldContent Field `gorm:"-"`
}

type Field struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}

func (ts *TblSchemaHandler) getColumns() {
	db := ts.db
	var cols []column
	db.Table("information_schema.COLUMNS").
		Select("COLUMN_NAME,DATA_TYPE,IS_NULLABLE,TABLE_NAME,COLUMN_COMMENT").
		Where("table_schema = DATABASE()").
		Where("TABLE_NAME", ts.tableName).
		Order("COLUMN_NAME").
		Find(&cols)
	var tscolunm []column
	for _, col := range cols {
		var tag string
		switch ts.modelOrmTagType {
		case ORM:
			tag = fmt.Sprintf("`orm:\"%s\" ", col.ColumnName)
		case GORM:
			tag = fmt.Sprintf("`gorm:\"column:%s\" ", col.ColumnName)
		default:
			tag = fmt.Sprintf("`gorm:\"column:%s\" ", col.ColumnName)
		}
		for _, v := range ts.otherTag {
			tag += fmt.Sprintf("%s:\"%s\" ", v, col.ColumnName)
		}
		tag += "`"
		fieldName := ts.camelCase(col.ColumnName, ts.fieldNameType)

		if len(fieldName) > ts.maxLenFieldName {
			ts.maxLenFieldName = len(fieldName)
		}
		if len(sqlTypeToGoType[col.Type]) > ts.maxLenFieldType {
			ts.maxLenFieldType = len(sqlTypeToGoType[col.Type])
		}
		if len(tag) > ts.maxLenFieldTag {
			ts.maxLenFieldTag = len(tag)
		}

		col.FieldContent = Field{
			Name:    fieldName,
			Type:    sqlTypeToGoType[col.Type],
			Tag:     tag,
			Comment: fmt.Sprintf("//是否可空:%s %s", col.Nullable, col.ColumnComment),
		}

		// col.ColunmContent = fmt.Sprintf("%s %s %s//是否可空：%s %s\n",
		// 	fieldName,
		// 	sqlTypeToGoType[col.Type],
		// 	tag,
		// 	col.Nullable,
		// 	col.ColumnComment)
		tscolunm = append(tscolunm, col)

		// ts.columns = append(ts.columns, col)
	}

	ts.columns = tscolunm

}

func (ts *TblSchemaHandler) camelCase(str string, Type string) string {

	var text string
	//不开启字段转为骆驼写法则仅仅将首字母大写
	switch Type {
	case CAMEL_CASE:
		for _, p := range strings.Split(str, "_") {
			text += strings.ToUpper(p[0:1]) + p[1:]
		}
	case FIRST_UPPER:
		text += strings.ToUpper(str[0:1]) + strings.ToLower(str[1:])
	default:
		for _, p := range strings.Split(str, "_") {
			text += strings.ToUpper(p[0:1]) + p[1:]
		}
	}

	return text
}

func (ts *TblSchemaHandler) connectSql() {
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
