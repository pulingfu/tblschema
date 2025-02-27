package dvaplugin

import (
	"fmt"

	"github.com/tidwall/gjson"
)

var test_json_parent = `
[
	{
		"id": 1,
		"name": "a",
		"age": 1,
		"body2":
		{
			"bid":11,
			"height":11,
			"desc":"11"
		},
		"body":[
			{
				"bid":1,
				"height":1,
				"desc":"1"
			},
			{
				"bid":11,
				"height":11,
				"desc":"11"
			}
		]
	},
	{
		"id": 2,
		"name": "b",
		"age": 2,
		"body2":
		{
			"bid":22,
			"height":22,
			"desc":"22"
		},
		"body":[
			{
				"bid":2,
				"height":2,
				"desc":"2"
			},
			{
				"bid":22,
				"height":22,
				"desc":"22"
			}
		]
	},
	{
		"id": 3,
		"name": "c",
		"age": 3,
		"bid":333,
		"body2":
		{
			"bid":33,
			"height":33,
			"desc":"33"
		},
		"body":[
			{
				"bid":33,
				"height":33,
				"desc":"33"
			},
			{
				"bid":333,
				"height":33,
				"desc":"33"
			},
			{
				"bid":333,
				"height":33,
				"desc":"33"
			}
		]
	}
]
`

var test_json_sub = `
[
	{
		"bid":"11",
		"type":"foo"
	},
	{
		"bid":"11",
		"type":"bar"
	},
	{
		"bid":"2",
		"type":"foo2"
	},
	{
		"bid":"333",
		"type":"bar333"
	},
	{
		"bid":"333",
		"type":"bar3333"
	}
]

`

func ExampleDataer_GetKeys() {
	dataer := NewDataer()
	dataer.GetKeys(gjson.Parse(test_json_parent), "body2|bid")
	fmt.Printf("\nkeys:=%v", dataer.Keys)
}

// 将subdata arry 中符合条件的单个元素，加入到parent 指定位置中
func ExampleDataer_HasOne() {

	dataer := NewDataer().
		SetMeta(test_json_parent). //测试json在本example所在文件中，可到github查看
		SetCompareFunc(func(p, s gjson.Result) bool {
			return p.Get("bid").Int() == s.Get("bid").Int()
		}).
		SetSubGroup(gjson.Parse(test_json_sub))

	dataer.HasOne(gjson.Parse(test_json_parent), "", "body2|newbd")

	fmt.Printf("\nresult:=%v", VtoJsonString(dataer.GetResult()))
}

// 将subdata arry 中符合条件的多个元素，加入到parent 指定位置中
func ExampleDataer_HasMany() {

	dataer := NewDataer().
		SetMeta(test_json_parent). //测试json在本example所在文件中，可到github查看
		SetCompareFunc(func(p, s gjson.Result) bool {
			return p.Get("bid").Int() == s.Get("bid").Int()
		}).
		SetSubGroup(gjson.Parse(test_json_sub))

	dataer.HasMany(gjson.Parse(test_json_parent), "", "body2|newbd")

	fmt.Printf("\nresult:=%v", VtoJsonString(dataer.GetResult()))
}
