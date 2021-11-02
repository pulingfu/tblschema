package tblschema

import (
	"fmt"
	"os"
	"path/filepath"
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
	"float":              "float64",
	"double":             "float64",
	"decimal":            "float64",
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
	"binary":             "string",
	"varbinary":          "string",
	"json":               "string",

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
	CAMEL_CASE  = "CamelCase"
	FIRST_UPPER = "First_upper"

	ORM  = "orm"
	GORM = "gorm"

	TIMETYPE_STRING = "string"
	TIMETYPE_TIME   = "time"

	//字典顺序
	FIELD_ORDER_FIELD_NAME = "COLUMN_NAME"
	//FIELD_ORDER_ORDINAL_POSITION数据库字段建立顺序
	FIELD_ORDER_ORDINAL_POSITION = "ORDINAL_POSITION"
)

type TblToStructHandler struct {
	dsn string   //数据库连接dsn,列：用户:密码@(127.0.0.1:3306)/数据库?charset=utf8mb4&parseTime=True&loc=Local
	db  *gorm.DB //数据库连接

	tableName string //要生成model的数据库表名
	savePath  string //保存model文件的位置
	timeType  string //时间类型对应go类型

	packageInfo         packageInfo         //模型文件包名配置
	tblStructNameInfo   tblStructNameInfo   //结构体模型名配置
	tblStructColumnInfo tblStructColumnInfo //结构体内容配置
}

type packageInfo struct {
	PackageName   string //包名
	PackagePrefix string //包名前缀
	PackageSuffix string //包名后缀
}
type tblStructNameInfo struct {
	TblStructName     string
	TblStructPrefix   string //包名前缀
	TblStructSuffix   string //包名后缀
	TblStructNameType string //结构体名类型 CAMEL_CASE骆驼命名/FIRST_UPPER首字母大写
}

type tblStructColumnInfo struct {
	ColumnOrder      string //排序方式
	ColumnNameType   string //字段名类型 CAMEL_CASE骆驼命名/FIRST_UPPER首字母大写
	ColumnNameSuffix string //生成模型行后缀
	ColumnNamePrefix string //生成模型行前缀

	ModelOrmTagType string   //生成orm结构题标签类型
	OtherTag        []string //其他标签

	//数据库对应行信息
	Columns         []column
	MaxLenFieldType int
	MaxLenFieldTag  int
	MaxLenFieldName int
}

func NewTblToStructHandler() *TblToStructHandler {
	return &TblToStructHandler{
		packageInfo: packageInfo{
			PackageName:   "model",
			PackageSuffix: "",
			PackagePrefix: "",
		},
		tblStructNameInfo: tblStructNameInfo{
			TblStructPrefix:   "",
			TblStructSuffix:   "",
			TblStructNameType: CAMEL_CASE,
		},
		tblStructColumnInfo: tblStructColumnInfo{
			ModelOrmTagType:  "gorm",
			ColumnNameType:   CAMEL_CASE,
			ColumnOrder:      FIELD_ORDER_ORDINAL_POSITION,
			ColumnNamePrefix: "",
			ColumnNameSuffix: "",
		},

		savePath: "model.go",
	}
}

func (ts *TblToStructHandler) SetDsn(dsn string) *TblToStructHandler {
	ts.dsn = dsn
	return ts
}
func (ts *TblToStructHandler) SetDB(db *gorm.DB) *TblToStructHandler {
	ts.db = db
	return ts
}
func (ts *TblToStructHandler) SetSavePath(savePath string) *TblToStructHandler {
	ts.savePath = savePath
	return ts
}
func (ts *TblToStructHandler) SetTableName(tableName string) *TblToStructHandler {
	ts.tableName = tableName
	return ts
}

//设置所生成对应的orm 标记类型  默认为 `gorm:"column:xxx"`
func (ts *TblToStructHandler) SetStructOrmTag(modelOrmTagType string) *TblToStructHandler {
	ts.tblStructColumnInfo.ModelOrmTagType = modelOrmTagType
	return ts
}

//添加其他的标签 如json ==> `json:"xxx"`
func (ts *TblToStructHandler) SetOtherTags(otherTag ...string) *TblToStructHandler {
	ts.tblStructColumnInfo.OtherTag = otherTag
	return ts
}

//设置行信息
func (ts *TblToStructHandler) SeTblStructColumnNameInfo(columnNameType, columnOrder, prefix, suffix string) *TblToStructHandler {
	ts.tblStructColumnInfo.ColumnNameType = columnNameType
	ts.tblStructColumnInfo.ColumnNameSuffix = suffix
	ts.tblStructColumnInfo.ColumnNamePrefix = prefix
	ts.tblStructColumnInfo.ColumnOrder = columnOrder
	return ts
}

//默认生成的结构体名类型为CamelCase写法，无前后后缀
func (ts *TblToStructHandler) SetTblStructNameInfo(structNameType, prifix, suffix string) *TblToStructHandler {
	ts.tblStructNameInfo.TblStructNameType = structNameType
	ts.tblStructNameInfo.TblStructPrefix = prifix
	ts.tblStructNameInfo.TblStructSuffix = suffix
	return ts
}

//设置包信息
func (ts *TblToStructHandler) SetPackageInfo(packageName, prifix, suffix string) *TblToStructHandler {
	ts.packageInfo.PackageName = packageName
	ts.packageInfo.PackagePrefix = prifix
	ts.packageInfo.PackageSuffix = suffix
	return ts
}

//设置数据库中的时间类型对应 modelstruct中的什么 time.Time/string
func (ts *TblToStructHandler) SetTimeType(timeType string) *TblToStructHandler {
	ts.timeType = timeType
	return ts
}

func (ts *TblToStructHandler) GenerateTblStruct() *TblToStructHandler {
	ts.connectSql()
	ts.getColumns()
	packageName := fmt.Sprintf("package %s%s%s\n",
		ts.packageInfo.PackagePrefix,
		ts.packageInfo.PackageName,
		ts.packageInfo.PackageSuffix,
	)
	var timetypeCount int64
	ts.db.Table("INFORMATION_SCHEMA.COLUMNS").
		Where("TABLE_SCHEMA=database()").
		Where("TABLE_NAME", ts.tableName).
		Where("DATA_TYPE in ?", []string{
			"date", "datetime", "timestamp", "time",
			"DATE", "DATETIME", "TIMESTAMP", "TIME",
		}).
		Count(&timetypeCount)

	packageimport := ""
	if ts.timeType != TIMETYPE_STRING && timetypeCount > 0 {
		packageimport = "import \"time\"\n"
	}

	var tableComment string
	// select * from INFORMATION_SCHEMA.TABLES where TABLE_SCHEMA=database()
	ts.db.Table("INFORMATION_SCHEMA.TABLES").
		Select("TABLE_COMMENT").
		Where("TABLE_SCHEMA=database()").
		Where("TABLE_NAME", ts.tableName).Find(&tableComment)
	tableComment = fmt.Sprintf("//%s\n", tableComment)
	ts.tblStructNameInfo.TblStructName = fmt.Sprintf("%s%s%s",
		ts.tblStructNameInfo.TblStructPrefix,
		ts.tableName,
		ts.tblStructNameInfo.TblStructSuffix,
	)
	structName := ts.generateChangeChara(ts.tblStructNameInfo.TblStructName, ts.tblStructNameInfo.TblStructNameType)
	structContent := fmt.Sprintf("type %s struct{\n", structName)

	for _, v := range ts.tblStructColumnInfo.Columns {
		match := fmt.Sprint("\t%-",
			ts.tblStructColumnInfo.MaxLenFieldName, "s %-",
			ts.tblStructColumnInfo.MaxLenFieldType, "s %-",
			ts.tblStructColumnInfo.MaxLenFieldTag, "s %s\n")
		structContent += fmt.Sprintf(match, v.FieldContent.Name, v.FieldContent.Type, v.FieldContent.Tag, v.FieldContent.Comment)
	}
	structContent += "}\n"

	functableName := fmt.Sprintf("func (*%s) TableName() string {\n", structName) +
		fmt.Sprintf("\t return \"%s\"\n", ts.tableName) +
		"}\n"

	fileContent := fmt.Sprintf("%s\n%s\n%s%s\n%s", packageName, packageimport, tableComment, structContent, functableName)

	// fmt.Println(fileContent)
	filePath := fmt.Sprint(ts.savePath)

	paths, _ := filepath.Split(ts.savePath)
	// fmt.Println(paths)
	os.MkdirAll(paths, os.ModePerm)
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("\x1b[%dm->table: %s 生成失败\x1b[0m\n", 31, ts.tableName)
		return ts
	}
	defer f.Close()

	f.WriteString(fileContent)
	fmt.Printf("\x1b[%dm->table: %s 生成成功\x1b[0m\n", 32, ts.tableName)
	// fmt.Printf("", )
	return ts
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

func (ts *TblToStructHandler) getColumns() {
	if ts.tableName == "" {
		panic("请先调用SetTableName设置要生成结构的数据库表哦")
	}
	db := ts.db
	var cols []column
	qr := db.Table("information_schema.COLUMNS").
		Select("COLUMN_NAME,DATA_TYPE,IS_NULLABLE,TABLE_NAME,COLUMN_COMMENT").
		Where("table_schema = DATABASE()").
		Where("TABLE_NAME", ts.tableName)
	switch ts.tblStructColumnInfo.ColumnOrder {
	case FIELD_ORDER_FIELD_NAME:
		qr.Order("COLUMN_NAME").
			Find(&cols)
	case FIELD_ORDER_ORDINAL_POSITION:
		qr.Order("ORDINAL_POSITION").
			Find(&cols)
	case "":
		qr.Order("COLUMN_NAME").
			Find(&cols)
	default:
		qr.Order(ts.tblStructColumnInfo.ColumnOrder).
			Find(&cols)
	}
	if len(cols) < 1 {
		panic("此表不存在或者数据库连接 不正确请检查哦")
	}
	ts.tblStructColumnInfo.MaxLenFieldName = 0
	ts.tblStructColumnInfo.MaxLenFieldTag = 0
	ts.tblStructColumnInfo.MaxLenFieldType = 0
	var tscolunm []column
	for _, col := range cols {
		switch ts.timeType {
		case TIMETYPE_STRING:
			switch col.Type {
			case "date", "datetime", "timestamp", "time":
				col.Type = fmt.Sprintf("%s_string", col.Type)
			}
		}
		var tag string
		switch ts.tblStructColumnInfo.ModelOrmTagType {
		case ORM:
			tag = fmt.Sprintf("`orm:\"%s\" ", col.ColumnName)
		case GORM:
			tag = fmt.Sprintf("`gorm:\"column:%s\" ", col.ColumnName)
		default:
			tag = fmt.Sprintf("`gorm:\"column:%s\" ", col.ColumnName)
		}
		for _, v := range ts.tblStructColumnInfo.OtherTag {
			if v != "" {
				tag += fmt.Sprintf("%s:\"%s\" ", v, col.ColumnName)
			}
		}
		tag += "`"
		fieldName := fmt.Sprintf("%s%s%s",
			ts.tblStructColumnInfo.ColumnNamePrefix,
			col.ColumnName,
			ts.tblStructColumnInfo.ColumnNameSuffix,
		)
		fieldName = ts.generateChangeChara(fieldName, ts.tblStructColumnInfo.ColumnNameType)

		if len(fieldName) > ts.tblStructColumnInfo.MaxLenFieldName {
			ts.tblStructColumnInfo.MaxLenFieldName = len(fieldName)
		}
		if len(sqlTypeToGoType[col.Type]) > ts.tblStructColumnInfo.MaxLenFieldType {
			ts.tblStructColumnInfo.MaxLenFieldType = len(sqlTypeToGoType[col.Type])
		}
		if len(tag) > ts.tblStructColumnInfo.MaxLenFieldTag {
			ts.tblStructColumnInfo.MaxLenFieldTag = len(tag)
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

	ts.tblStructColumnInfo.Columns = tscolunm

}

func (ts *TblToStructHandler) generateChangeChara(str string, Type string) string {

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

func (ts *TblToStructHandler) connectSql() {
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

func (ts *TblToStructHandler) GetAllTableNames() []string {
	ts.connectSql()
	var allTname []string
	ts.db.Table("INFORMATION_SCHEMA.TABLES").
		Select("TABLE_NAME").
		Where("TABLE_SCHEMA=database()").
		Find(&allTname)
	return allTname
}

func (ts *TblToStructHandler) GenerateAllTblStruct() {
	ts.connectSql()
	allTname := ts.GetAllTableNames()
	for _, tname := range allTname {
		ts.
			SetPackageInfo(tname, "tbl_", "").
			SetSavePath(fmt.Sprintf("./tbl_%s/schema_model.go", tname)).
			SetTableName(tname).
			GenerateTblStruct()
	}
}
