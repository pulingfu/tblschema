package model

type TblProject struct {
	Id        int    `gorm:"column:id" json:"id"`
	Name      string `gorm:"column:name" json:"name"`
	CompanyId int    `gorm:"column:company_id" json:"company_id"`
	CreatorId int    `gorm:"column:creator_id" json:"creator_id"`
}

func (*TblProject) TableName() string {
	return "tbl_project"
}
