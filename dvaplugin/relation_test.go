package dvaplugin

// 用户
type TblUser struct {
	Id   int    `json:"id" gorm:"column:id"`
	Name string `json:"name" gorm:"column:name"`
}

func (u TblUser) TableName() string {
	return "user"
}

// 学生
type TblStudent struct {
	Id      int    `json:"id" gorm:"column:id"`
	ClassId string `json:"class_id" gorm:"column:class_id"` //班级
	UserId  int    `json:"user_id" gorm:"column:user_id"`   // 用户ID
}

func (s TblStudent) TableName() string {
	return "student"
}

// 老师
type TblTeacher struct {
	Id      int    `json:"id" gorm:"column:id"`
	ClassId string `json:"class_id" gorm:"column:class_id"` //班级
	UserId  int    `json:"user_id" gorm:"column:user_id"`   // 用户ID
	Subject string `json:"subject" gorm:"column:subject"`   // 科目
}

func (t TblTeacher) TableName() string {
	return "teacher"
}

// 班级
type TblClass struct {
	Id   int    `json:"id" gorm:"column:id"`
	Name string `json:"name" gorm:"column:name"`
}

func (c TblClass) TableName() string {
	return "class"
}

// 加载多个表或者多个结构时间的嵌套关系
func ExampleRelationLoader_LoadResult() {

}
