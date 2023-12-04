package model

type TblProjectDetail struct {
	Id        int `gorm:"column:id" json:"id"`
	ProjectId int `gorm:"column:project_id" json:"project_id"`
	CreatorId int `gorm:"column:creator_id" json:"creator_id"`
}

func (*TblProjectDetail) TableName() string {
	return "tbl_project_detail"
}
