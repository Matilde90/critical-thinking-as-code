# Critical Thinking as Code

Critical thinking as code is a linter for reasoning. It will help engineers to make reasoning explicit, and encourage robust reasoning.
It will let you codify arguments as structured data and it will give instant and transparent feedback on potential weak reasoning patterns:

- Missing premises
- Circular reasoning
- Vagueness/hedging
- Overgeneralization
- Modality-confidence mismatch

## Build

`make build`

## Installation

go install ./cmd/ctac

Put the line below in your ~/.bashrc or ~/.zshrc to persist it
export PATH="$(go env GOPATH)/bin:$PATH"