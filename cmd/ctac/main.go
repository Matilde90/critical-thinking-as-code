package main

import (
	"ctac/pkg/ctac"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

func usage() {
	fmt.Println(`ctac -- Critical Thinking as Code
	
	Usage:
		ctac analyse	[flags]		Analyse an argument file
		ctac ignore		[subcmd]	Manage ignore file
		ctac create		[subcmd]	Create argument file
	
	Examples:
		ctac analyse -inputFile file.yaml -outputFile results.md -pretty
		ctac analyse -inputFile file.yaml parallel -workers 2 -outputFile results.md -pretty
		ctac ignore print-template

	Run "ctac <command> -h" for more information about a command.`)

}

func analyseCmd(args []string) {
	flagSet := flag.NewFlagSet("analsze", flag.ContinueOnError)
	flagSet.SetOutput(os.Stderr)

	inputFile := flagSet.String("inputFile", "", "Path to input argument yaml file")
	parallel := flagSet.Bool("parallel", false, "Run rules in parallel - false by default")
	workers := flagSet.Int("workers", 3, "Max concurrent workers (only used with parallel flag set as true)")
	outputFile := flagSet.String("outputFile", "", "Path to results JSON file")
	pretty := flagSet.Bool("pretty", false, "Pretty print JSON")
	silent := flagSet.Bool("silent", false, "Quiet mode to silence output written to standard out")
	ignoreFile := flagSet.String("ignoreFile", "", "Path to ignore file")

	if err := flagSet.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return
		}
		os.Exit(2)
	}

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
		if !*silent {
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

	ignored := make(map[string]bool, len(ignoreSpec.Rules))
	fmt.Printf("%v", ignored)
	for _, rule := range ignoreSpec.Rules {
		ignored[rule] = true
	}
	fmt.Printf("%v", ignored)
	for _, issue := range issues {
		fmt.Printf("Processing issue %v with rule id %s", issue, issue.RuleID)
		if _, skip := ignored[issue.RuleID]; !skip {
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
		if err := os.WriteFile(*outputFile, b, 0o644); err != nil {
			log.Fatalf("Write outputfile: %v", err)
		}
	}
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
	fmt.Println(os.Args)
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "analyze", "analyse":
		analyseCmd(os.Args[2:])
	case "ignore":
		fmt.Print("ignore")
	case "create":
		fmt.Print("create")
	case "help", "-h", "--help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", os.Args[1])
	}
}
