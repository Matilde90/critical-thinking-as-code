package main

import (
	"ctac/pkg/ctac"
	"fmt"
	"flag"
)


func main (){
	fmt.Println("Welcome to ctac, critical thinking as code")
	file := flag.String("file", "examples/decision.yaml", "path to argument yaml file")
	flag.Parse()

	data, err := ctac.Loader(*file)
	if err != nil {
		fmt.Printf("cannot unmarshal data: %v", err)
	}
	fmt.Println(*data)
	result:= ctac.RunAllRules(*data)
	fmt.Println(result)
}
