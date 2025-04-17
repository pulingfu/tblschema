package dvaplugin

import (
	"fmt"
	"time"

	"github.com/king-kkong/dataschema/cmd/test_alltabl_to_model/tbl_my_user"
	"github.com/tidwall/gjson"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ExampleGetUpdateMapping() {

	var configs = &gorm.Config{}
	db, err := gorm.Open(mysql.Open("root:tiger@(127.0.0.1:3306)/pulingfu?charset=utf8mb4&parseTime=True&loc=Local"), configs)
	if err != nil {
		panic(err)
	}

	reqjson := `
	{
		"id": 1,
		"user_name": "name",
		"user_pass": "pass",
		"user_status": "status"
	}`

	reqdata := gjson.Parse(reqjson)

	var metadata tbl_my_user.MyUser
	uuid := reqdata.Get("uuid").String()
	db.Model(&metadata).Where("uuid", uuid).Take(&metadata)
	if metadata.Id < 1 {
		fmt.Print("数据不存在")
		return
	}

	updates, err := GetUpdateMapping(metadata, reqdata)
	if err != nil {
		fmt.Printf("%v", err.Error())
		return
	}

	updates["updated_at"] = time.Now()
	db.Model(&metadata).Where("uuid", uuid).UpdateColumns(updates)

}

func ExampleGetCreateMapping() {

	var configs = &gorm.Config{}
	db, err := gorm.Open(mysql.Open("root:tiger@(127.0.0.1:3306)/pulingfu?charset=utf8mb4&parseTime=True&loc=Local"), configs)
	if err != nil {
		panic(err)
	}

	reqjson := `
	{
		"user_name": "name",
		"user_pass": "pass",
		"user_status": "status"
	}`

	reqdata := gjson.Parse(reqjson)

	var metadata tbl_my_user.MyUser
	creates, err := GetUpdateMapping(metadata, reqdata)
	if err != nil {
		fmt.Printf("%v", err.Error())
		return
	}
	creates["updated_at"] = time.Now()

	tx := db.Begin()

	metadata.Uuid = "abcd"
	if err := tx.Create(&metadata).Error; err != nil {
		tx.Rollback()
		fmt.Printf("%v", err.Error())
		return
	}

	if err := tx.Model(&metadata).Where("uuid", metadata.Uuid).UpdateColumns(creates).Error; err != nil {
		tx.Rollback()
		fmt.Printf("%v", err.Error())
		return
	}

	tx.Commit()
}
