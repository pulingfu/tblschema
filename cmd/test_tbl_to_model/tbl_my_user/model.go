package prefix_tbl_my_user_suffix

import "time"

//my_user表说明 用于测试
type A我是TableName前缀MyUserB我是TableName后缀 struct{
	A我是ColumnName前缀CreatedAtB我是ColumnName后缀                  time.Time `gorm:"column:created_at" `  //是否可空:YES 
	A我是ColumnName前缀IdB我是ColumnName后缀                         int       `gorm:"column:id" `          //是否可空:NO 唯一ID
	A我是ColumnName前缀UserNameB我是ColumnName后缀                   string    `gorm:"column:user_name" `   //是否可空:NO 用户名
	A我是ColumnName前缀UserPassB我是ColumnName后缀                   string    `gorm:"column:user_pass" `   //是否可空:YES 密码
	A我是ColumnName前缀UserStatusB我是ColumnName后缀                 string    `gorm:"column:user_status" ` //是否可空:YES 用户状态|active 活跃 这是一个很长的用户状态
}

func (*A我是TableName前缀MyUserB我是TableName后缀) TableName() string {
	 return "my_user"
}
