![brain image](./image/ctac-image.png)

# Critical Thinking as Code

Critical thinking as code is a command-line tool that help you keep track, record and strentghen your arguments using codified critical-thinking rules. It will help engineers to make reasoning explicit, clear and easy to version control.
It will let you codify arguments as structured data and it will give instant and transparent feedback on potential weak reasoning patterns. 
Ctac will guide you to write your argument as YAML and will detect weaknesses such as:

- Missing premises
- Circular reasoning
- Vagueness/hedging
- Overgeneralization
- Modality-confidence mismatch

## 🚀 Quickstart

### 1. Install

```bash
go install github.com/maliffi/ctac@latest
```
Or clone and build locally

```
git clone https://github.com/Matilde90/ctac.git
cd ctac
make build
```

## ✍🏼 Create a new argument

You can create a YAML file for the new argument using the ctac tool:

```sh
ctac create -filePath argument.yaml
```

You will be guided through prompts to add an argument title and provide premises and conclusion. For each premise you will be asked to provide a confidence level and for every conclusion you will be asked to provide confidence level and modality.

A YAML file like [decision.yaml](./examples/decision.yaml) will be created

---

## 🔎 Analyse the argument

```sh
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

## 🧠 Implemented Reasoning Rules

| RuleID| Description | Severity |
|---|--|---|
| CTAC001_MISSING_PREMISES  |   | error |
| CTAC002_VAGUENESS_DETECTED  |   | warning|
| CTAC003_MISSING_CONCLUSION_RULE  |  | error |
| CTAC004_SINGLE_PREMISE_RULE  |   | warning |
| CTAC005_MODALITY_MISMATCH_RULE  |   | error |
| CTAC006_QUANTIFICATION_REQUIRED  |   | error |
| CTAC007_EMOTIONAL_LANGUAGE_DETECTED  |  | error|
