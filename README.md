![brain image](./image/ctac-image.png)

# Critical Thinking as Code

Critical thinking as code (CTAC) is a command-line tool that helps engineers to keep track, record and strentghen their arguments using codified critical-thinking rules. It helps engineers to make reasoning explicit, clear and easy to version control. It lets codify arguments as structured data and it gives instant and transparent feedback on potential weak reasoning patterns.
CTAC guides engineers to write their argument as YAML and will detect weaknesses such as:

- Missing premises
- Circular reasoning
- Vagueness/hedging
- Overgeneralization
- Modality-confidence mismatch

*Please note: this is work in progress. More rules and functionalities will be added in the coming weeks*

## 🚀 Quickstart

### 1. Install

```bash
go install github.com/Matilde90/ctac@latest
```

Or clone and build locally

```bash
git clone https://github.com/Matilde90/ctac.git
cd ctac
make build
```

## ✍🏼 Create a new argument

Use CTAC to create a YAML file for a new argument

```bash
ctac create -filePath argument.yaml
```

You will be guided through prompts to:
- add an argument title
- provide premises and their confidence level
- add a conclusion with a confidence level and modality.

CTAC will generate a structured YAML file like [decision.yaml](./examples/decision.yaml).

---

## 🔎 Analyse the argument

```bash
ctac analyse -inputFile argument.yaml -pretty
```

## 🤖 Available Commands
|Command | Description | Example |
|--|--|--|
| ctac create | Interactive wizard to create a YAML argument file| ctac create -filePath argument.yaml
| ctac analyse | Analyse argument against built-in rules | ctac analyse -inputFile argument.yaml|
| ctac ignore | Prints a sample ignore file | ctac ignore print-template|
| ctac version| Prints version (set via -ldflags) | ctac version |
| ctac help | Displays usage help | ctac help

## 🤖 Options

### Create

`ctac create`
  -filePath string
        Path to input argument yaml file. Creates or truncates the name file with the argument provided in the interactive steps

### Analyse

`ctac analyse`
  -ignoreFile string
        Path to ignore file
  -inputFile string
        Path to input argument yaml file
  -outputFile string
        Path to results JSON file
  -parallel
        Run rules in parallel (default: false)
  -pretty
        Pretty-print JSON
  -silent
        Quiet mode to silence output written to standard out
  -workers int
        Max concurrent workers (only used with parallel flag set as true) (default 3)

### Ignore

`ctac ignore`
  ctac ignore print-template   # print a template to stdout

## 🧠 Implemented Reasoning Rules

| RuleID| Description | Severity |
|---|--|---|
| CTAC001_MISSING_PREMISES  | Flags arguments with no premise  | error |
| CTAC002_VAGUENESS_DETECTED  |  Flags arguments whose premises use vague words | warning|
| CTAC003_MISSING_CONCLUSION_RULE  | Flags arguments with no conclusion | error |
| CTAC004_SINGLE_PREMISE_RULE  | Flags arguments that have only one premise as these are often weak | warning |
| CTAC005_MODALITY_MISMATCH_RULE  | It flags arguments with a strong conclusion (modality must) with weak/insufficient support.  | error |
| CTAC006_QUANTIFICATION_REQUIRED  | It flags argument with premises using quantification terms that omits reference to actual numbers  | error |
| CTAC007_EMOTIONAL_LANGUAGE_DETECTED  | The argument uses emotional language as it can involve appeal to emotions bias | error|


## 🤝 Contributing

CTAC is in active development.
Feedback, suggestions, and pull requests are welcome!