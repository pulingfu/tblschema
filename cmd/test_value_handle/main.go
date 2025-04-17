package main

import (
	"fmt"

	"github.com/king-kkong/dataschema/dvap"
)

type AA struct {
	Name string
	Age  int
}

func main() {

	var as = []AA{
		{
			Name: "a",
			Age:  1,
		},
		{
			Name: "b",
			Age:  2,
		},
	}

	if ok, err := dvap.HasValueInSlice(as, 1,
		func(a AA, b string) bool { return a.Name == b }); ok {
		fmt.Println("ok")
	} else {
		fmt.Println(err.Error())
	}

}
