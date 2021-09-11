package all_tbl_model

import "time"

//my_user表说明 用于测试
type MyUser struct{
	Id         int       `gorm:"column:id" json:"id" `                   //是否可空:NO 唯一ID
	UserName   string    `gorm:"column:user_name" json:"user_name" `     //是否可空:NO 用户名
	UserPass   string    `gorm:"column:user_pass" json:"user_pass" `     //是否可空:YES 密码
	UserStatus string    `gorm:"column:user_status" json:"user_status" ` //是否可空:YES 用户状态|active 活跃 这是一个很长的用户状态
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at" `   //是否可空:YES 
}

func (*MyUser) TableName() string {
	 return "my_user"
}
