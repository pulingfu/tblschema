package model

type TblUser struct {
	Id   int    `gorm:"column:id" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

func (*TblUser) TableName() string {
	return "tbl_user"
}
