package tbl_plf_tbl_user

import "time"

//用户表注释
type PlfTblUser struct{
	Id        int       `gorm:"column:id" `         //是否可空:NO 
	Address   string    `gorm:"column:address" `    //是否可空:YES 地址
	Age       int       `gorm:"column:age" `        //是否可空:YES 年龄（非负整数,默认3岁）
	CreatedAt time.Time `gorm:"column:created_at" ` //是否可空:NO 创建记录的时间，默认自动创建
	Intro     string    `gorm:"column:intro" `      //是否可空:YES 简单介绍
	Nickname  string    `gorm:"column:nickname" `   //是否可空:YES 昵称
	Password  string    `gorm:"column:password" `   //是否可空:YES 密码
	Phone     string    `gorm:"column:phone" `      //是否可空:YES 手机号(可以为空,长度11)
	UpdatedAt time.Time `gorm:"column:updated_at" ` //是否可空:NO 上一次更新时间，默认自动更新
	Username  string    `gorm:"column:username" `   //是否可空:NO 用户名
}

func (*PlfTblUser) TableName() string {
	 return "plf_tbl_user"
}
