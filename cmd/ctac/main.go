package main

import (
	"ctac/pkg/ctac"
	"fmt"
	"flag"
)


func main (){
	fmt.Println("Welcome to ctac, critical thinking as code")
	fmt.Println()
	file := flag.String("file", "examples/decision.yaml", "path to argument yaml file")
	flag.Parse()

	argument, err := ctac.Loader(*file)
	if err != nil {
		fmt.Printf("cannot unmarshal data: %v", err)
	}
	fmt.Println(ctac.SummariseArgument(*argument))
	result:= ctac.RunAllRules(*argument)
	fmt.Println(result)
}
