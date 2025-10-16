package main

import (
	"ctac/pkg/ctac"
	"flag"
	"fmt"
	"os"
)

func main() {

	file := flag.String("file", "examples/decision.yaml", "Path to argument yaml file")
	parallel := flag.Bool("parallel", false, "Run rules in parallel - false by default")
	flag.Parse()

	fmt.Println("Welcome to ctac, critical thinking as code\n")

	argument, err := ctac.Loader(*file)
	if err != nil {
		fmt.Printf("cannot unmarshal data: %v", err)
		os.Exit(1)
	}

	fmt.Println(ctac.SummariseArgument(*argument))

	var issues []ctac.Issue
	if *parallel {
		issues = ctac.RunAllRulesParallel(*argument, 3)
	} else {
		fmt.Println("Running all rules sequentially")
		issues = ctac.RunAllRulesSequential(*argument)
	}

	fmt.Println(ctac.FormatIssueMessage(issues))

	// TODO: emit JSON to stdout
}
