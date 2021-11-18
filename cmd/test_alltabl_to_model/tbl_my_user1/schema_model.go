package tbl_my_user1

import "time"

//my_user表说明 用于测试
type MyUser1 struct{
	Id         int       `gorm:"column:id,omitempty" `          //是否可空:NO 唯一ID
	UserName   string    `gorm:"column:user_name,omitempty" `   //是否可空:NO 用户名
	UserPass   string    `gorm:"column:user_pass,omitempty" `   //是否可空:YES 密码
	UserStatus string    `gorm:"column:user_status,omitempty" ` //是否可空:YES 用户状态|active 活跃 这是一个很长的用户状态
	CreatedAt  time.Time `gorm:"column:created_at,omitempty" `  //是否可空:YES 
}

func (*MyUser1) TableName() string {
	 return "my_user1"
}
