package main

import (
	tblschema "github.com/king-kkong/dataschema"
)

func main() {
	// files, err := ioutil.ReadDir(".")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// for _, f := range files {
	// 	s := string(f.Name()) + "------"
	// 	fmt.Println(s)
	// }
	// InitConfig()
	yts := tblschema.NewYamlToSqlHandler().SetYamlPath("./etc2/").
		SetDsn("root:tiger@(127.0.0.1:3306)/pulingfu?charset=utf8mb4&parseTime=True&loc=Local")

	yts.ExecuteSchemaSafeCheck()
}

// func InitConfig() {
// 	files, err := ioutil.ReadDir(".")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	for _, f := range files {
// 		fmt.Println(f.Name())
// 	}
// 	yamlFile, err := ioutil.ReadFile("./etc/Entity.BiCountPointToSkin.dcm.yml")
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	table := map[string]interface{}{}
// 	err = yaml.Unmarshal(yamlFile, &table)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	jb, _ := json.Marshal(&table)
// 	fmt.Println(string(jb))
// 	// v, ok := table["Table"].(map[string]interface{})
// 	// if ok {
// 	// 	for _, field := range v {
// 	// 		fmt.Println("===")
// 	// 		fmt.Println(field)
// 	// 	}

// 	// }

// 	// fmt.Printf("config.app: %#v\n", table)
// 	// fmt.Printf("config.log: %#v\n", table)

// }
