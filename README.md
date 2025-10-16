# Critical Thinking as Code

Critical thinking as code is a linter for reasoning. It will help engineers to make reasoning explicit, and encourage robust reasoning.
It will let you codify arguments as structured data and it will give instant and transparent feedback on potential weak reasoning patterns:

- Missing premises
- Circular reasoning
- Vagueness/hedging
- Overgeneralization
- Modality-confidence mismatch


Todo:

concurrency: run rules in parallel with a worker pool or errgroup + context. This has bounded concurrency. Added panic safety to ensure a biggy rule won't crash the whole run and a deterministic order by collecting results by rule index

generics: add small generic helper for dedup/collect;

--json flag (CLI) to emit json issues
exit code policy when errors found when -- fail-on=error

ignore issues

tiny cmd/ctac-server that accepts YAML and returns JSON issues (Gin or net/http). Shows HTTP and JSON.

Rule docs: link each rule ID in README to a short markdown explaining heuristics and limitations.

Readme updated with explanation