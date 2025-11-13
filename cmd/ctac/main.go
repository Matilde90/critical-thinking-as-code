package main

import (
	"bufio"
	"ctac/pkg/ctac"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	version = "dev" // set via -ldflags
)

func usage() {
	fmt.Println(`ctac -- Critical Thinking as Code
	
	Usage:
		ctac analyse	[flags]		Analyse an argument file
		ctac ignore		[subcmd]	Manage ignore file
		ctac create		[subcmd]	Create argument file
		ctac version				Version
	
	Examples:
		ctac analyse -inputFile file.yaml -outputFile results.md -pretty
		ctac analyse -inputFile file.yaml -parallel -workers 2 -outputFile results.md -pretty
		ctac ignore print-template
		ctac create -filePath myargument.yaml
		ctac version

	Run "ctac <command> -h" for more information about a command.`)

}

func writeConfidenceLevel(file *os.File, scanner *bufio.Scanner) {
	fmt.Print("Please provide the confidence level:\n1.high\n2.medium\n3.low\ndefault: medium\n> ")
	if !scanner.Scan() {
		log.Fatalf("Could not read confidence level: %s", scanner.Err())
	}

	switch strings.ToLower(strings.TrimSpace(scanner.Text())) {
	case "1", "high", "h", "H":
		fmt.Fprintf(file, "    confidence: high\n")
	case "2", "medium", "m", "M":
		fmt.Fprintf(file, "    confidence: medium\n")
	case "3", "low", "l", "L":
		fmt.Fprintf(file, "    confidence: low\n")
	default:
		fmt.Fprintf(file, "    confidence: medium\n")
	}
}
func writePremise(file *os.File, id int, scanner *bufio.Scanner) {
	fmt.Printf("Please provide the premise text (single line):\n> ")
	if !scanner.Scan() {
		log.Fatalf("Could not read premise: %v", scanner.Err())
	}

	premise := strings.TrimSpace(scanner.Text())
	fmt.Fprintf(file, "-   id: P%d\n    text: %q\n", id, premise)

	writeConfidenceLevel(file, scanner)

	fmt.Printf("Do you want to add another premise\n [y/n] > ")
	if !scanner.Scan() {
		log.Fatalf("Could not read answer %v", scanner.Err())
	}

	switch strings.ToLower(strings.TrimSpace(scanner.Text())) {
	case "yes", "y":
		writePremise(file, id+1, scanner)
	case "no", "n":
		return
	default:
		fmt.Print("I do not recognise this. Assuming no")
		return
	}
}

func writeModality(file *os.File, scanner *bufio.Scanner) {
	fmt.Print("Please provide the modality:\n1.must\n2.should\n3.could\ndefault: should\n> ")
	if !scanner.Scan() {
		log.Fatalf("Could not read modality: %s", scanner.Err())
	}

	switch strings.ToLower(strings.TrimSpace(scanner.Text())) {
	case "1", "must", "m":
		fmt.Fprintf(file, "    modality: must\n")
	case "2", "should", "s":
		fmt.Fprintf(file, "    modality: should\n")
	case "3", "could", "c":
		fmt.Fprintf(file, "    modality: could\n")
	default:
		fmt.Println("I do not recognise this. Assuming modality: should")
		fmt.Fprintf(file, "    modality: should\n")
	}
}

func writeConclusion(file *os.File, scanner *bufio.Scanner) {
	fmt.Printf("Please provide the conclusion text (single line):\n> ")
	if !scanner.Scan() {
		log.Fatalf("Could not read the conclusion. Encountered error: %v", scanner.Err())
	}
	conclusion := strings.TrimSpace(scanner.Text())
	fmt.Fprintf(file, "conclusion:\n    text: %q\n", conclusion)
	writeConfidenceLevel(file, scanner)
	writeModality(file, scanner)
}

func createCmd(args []string) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	flagSet := flag.NewFlagSet("create", flag.ContinueOnError)
	flagSet.SetOutput(os.Stderr)
	filePath := flagSet.String("filePath", "", "Path to input argument yaml file. Creates or truncates the name file with the argument provided in the interactive steps")

	if err := flagSet.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return
		}
		os.Exit(2)
	}

	log.SetFlags(0)

	if *filePath == "" {
		log.Fatalf("error: Path to input argument yaml file is required via -filePath")
	}

	file, err := os.Create(*filePath)
	if err != nil {
		log.Fatalf("Failed to create input file: encountered error: %v", err)
	}
	fmt.Printf("Please provide a title\n>  ")
	if !scanner.Scan() {
		log.Fatalf("Could not read title. Encountered error: %v", scanner.Err())
	}
	fmt.Fprintf(file, "title: %q\npremises:\n", scanner.Text())
	id := 1
	writePremise(file, id, scanner)
	writeConclusion(file, scanner)
	fmt.Printf("yaml file created at %s\n", file.Name())

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("warning: closing file failed: %v", err)
		}
	}()
}

func analyseCmd(args []string) {
	flagSet := flag.NewFlagSet("analyse", flag.ContinueOnError)
	flagSet.SetOutput(os.Stderr)

	inputFile := flagSet.String("inputFile", "", "Path to input argument yaml file")
	parallel := flagSet.Bool("parallel", false, "Run rules in parallel (default: false)")
	workers := flagSet.Int("workers", 3, "Max concurrent workers (only used with parallel flag set as true)")
	outputFile := flagSet.String("outputFile", "", "Path to results JSON file")
	pretty := flagSet.Bool("pretty", false, "Pretty-print JSON")
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
	for _, rule := range ignoreSpec.Rules {
		ignored[rule] = true
	}
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

func ignoreCmd(args []string) {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Println(`Usage:
  ctac ignore print-template   # print a template to stdout`)
		return
	}
	switch args[0] {
	case "print-template":
		fmt.Println(`# ctac.ignore.yaml
		rules:
			- CTAC002_VAGUENESS_DETECTED
		reason:
			- "Describe why this ignore exists"`)
	default:
		fmt.Fprintf(os.Stderr, "unknown ignore subcommand %q\n", args[0])
		os.Exit(2)
	}
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "create", "-c":
		createCmd(os.Args[2:])
	case "analyze", "analyse", "-a":
		analyseCmd(os.Args[2:])
	case "ignore", "-i":
		ignoreCmd(os.Args[2:])
	case "help", "-h", "--help", "man":
		usage()
	case "version", "-v":
		fmt.Printf("ctac %s\n", version)
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", os.Args[1])
	}
}
