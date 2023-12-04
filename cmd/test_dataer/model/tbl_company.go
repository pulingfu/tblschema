package model

type TblCompany struct {
	Id        int    `gorm:"column:id" json:"id"`
	Name      string `gorm:"column:name" json:"name"`
	CreatorId int    `gorm:"column:creator_id" json:"creator_id"`
}

func (*TblCompany) TableName() string {
	return "tbl_company"
}
