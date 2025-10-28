package main

import (
	"ctac/pkg/ctac"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {

	inputFile := flag.String("inputFile", "", "Path to input argument yaml file")
	parallel := flag.Bool("parallel", false, "Run rules in parallel - false by default")
	workers := flag.Int("workers", 3, "Max concurrent workers (only used with parallel flag set as true)")
	outputFile := flag.String("outputFile", "", "Path to results JSON file")
	pretty := flag.Bool("pretty", false, "Pretty print JSON")
	silent := flag.Bool("silent", false, "Quite mode to silence output written to standard out")

	flag.Parse()

	log.SetFlags(0)

	if *inputFile == "" {
		log.Fatalf("error: -inputFile is required")
	}

	argument, err := ctac.Loader(*inputFile)
	if err != nil {
		log.Fatalf("load input error: %v", err)
	}

	if !*silent {
		fmt.Println("Welcome to ctac, critical thinking as code")
		fmt.Println(ctac.SummariseArgument(*argument))
	}


	var issues []ctac.Issue
	if *parallel {
		fmt.Println("Running all rules in parallel")
		issues = ctac.RunAllRulesParallel(*argument, *workers)
	} else {
		issues = ctac.RunAllRulesSequential(*argument)
	}

	if !*silent {
	fmt.Println(ctac.FormatIssueMessage(issues))
	}
	if *outputFile != "" {
		var b []byte
		if *pretty {
			b, err = json.MarshalIndent(issues, "", " ")
		} else {
		b, err = json.Marshal(issues)
		}
		if err != nil {
			log.Fatalf("error encoding JSON: %v", err)
		}
		os.WriteFile(*outputFile, b, 0o644)
	}
}
