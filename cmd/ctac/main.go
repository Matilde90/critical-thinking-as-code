package main

import (
	"ctac/pkg/ctac"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
)

func main() {

	inputFile := flag.String("inputFile", "", "Path to input argument yaml file")
	parallel := flag.Bool("parallel", false, "Run rules in parallel - false by default")
	workers := flag.Int("workers", 3, "Max concurrent workers (only used with parallel flag set as true)")
	outputFile := flag.String("outputFile", "", "Path to results JSON file")
	pretty := flag.Bool("pretty", false, "Pretty print JSON")
	silent := flag.Bool("silent", false, "Quiet mode to silence output written to standard out")
	ignoreFile := flag.String("ignoreFile", "", "Path to ignore file")

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
		if !*silent{
			fmt.Println("Running all rules in parallel")
		}
		issues = ctac.RunAllRulesParallel(*argument, *workers)
	} else {
		issues = ctac.RunAllRulesSequential(*argument)
	}

	var filteredIssues []ctac.Issue
	ignoreSpec, err := ctac.LoadIgnore(*ignoreFile)
	if err != nil {
		log.Fatalf("Load ignore file error: %v", err)
	}
	for _, issue := range issues {
		if !slices.Contains(ignoreSpec.Rules, issue.RuleID) {
			filteredIssues = append(filteredIssues, issue)
		}
	}
	if !*silent {
		fmt.Println(ctac.FormatIssueMessage(filteredIssues))
	}
	if *outputFile != "" {
		var b []byte
		if *pretty {
			b, err = json.MarshalIndent(filteredIssues, "", "  ")
		} else {
			b, err = json.Marshal(filteredIssues)
		}
		if err != nil {
			log.Fatalf("error encoding JSON: %v", err)
		}
		if err := os.WriteFile(*outputFile, b, 0o644) ; err != nil {
			log.Fatalf("Write outputfile: %v", err)
		}
	}
}
