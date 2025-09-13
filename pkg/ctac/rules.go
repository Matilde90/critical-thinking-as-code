package ctac

import (
	"fmt"
	"regexp"
)

type Rule interface {
	ID() string
	Check(a Argument) []Issue
}

type Issue struct {
	RuleID   string
	Severity Severity
	Message  string
	Hint     string
}

type Severity string

const (
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

type MissingPremiseRule struct{}
type VaguenessDetector struct{}
type MissingConclusionRule struct{}
type SinglePremiseRule struct{}
type ModalityMismatchRule struct{}
type QuantificationRequiredRule struct{}

func (r MissingPremiseRule) ID() string {
	return "CTAC001_MISSING_PREMISES"
}

func (rule VaguenessDetector) ID() string {
	return "CTAC002_VAGUENESS_DETECTOR"
}

func (rule MissingConclusionRule) ID() string {
	return "CTAC003_MISSING_CONCLUSION_RULE"
}

func (rule SinglePremiseRule) ID() string {
	return "CTAC004_SINGLE_PREMISE_RULE"
}

func (rule ModalityMismatchRule) ID() string {
	return "CTAC005_MODALITY_MISMATCH_RULE"
}

func (rule QuantificationRequiredRule) ID() string {
	return "CTAC006_QUANTIFICATION_REQUIRED"
}

type vaguePhrase struct {
	phrase string
	reg    *regexp.Regexp
}

var vaguePhrases = []vaguePhrase{
	{
		phrase: "someone",
		reg:    regexp.MustCompile(`(?i)\bsomeone\b`),
	},
	{
		phrase: "some",
		reg:    regexp.MustCompile(`(?i)\bsome\b`),
	},
	{
		phrase: "everyone thinks",
		reg:    regexp.MustCompile(`(?i)\beveryone thinks\b`),
	},
	{
		phrase: "maybe",
		reg:    regexp.MustCompile(`(?i)\bmaybe\b`),
	},
	{
		phrase: "everyone knows",
		reg:    regexp.MustCompile(`(?i)\beveryone knows\b`),
	},
}

var regexDigit = regexp.MustCompile("[0-9]+")
var regexQuantificationPhrase = regexp.MustCompile(`(?i)(\bsignificant|\bdecrease|\bmost\b|\bincrease|\bdecline\b|\bpercent(age?)\b|%|\bmore\b|\bless\b|\brate\b|\btrend\b)`)

func (rule VaguenessDetector) Check(argument Argument) []Issue {
	var issues []Issue

	premises := argument.Premises

	for i, p := range premises {

		for _, vaguePhrase := range vaguePhrases {
			if vaguePhrase.reg.MatchString(p.Text) {
				issues = append(issues, Issue{
					RuleID:   rule.ID(),
					Severity: SeverityWarning,
					Message:  fmt.Sprintf("Premise %d '%q' contains vague words '%s'", i+1, p.Text, vaguePhrase.phrase),
					Hint:     "",
				})

			}
		}
	}
	return issues

}

func (r MissingPremiseRule) Check(argument Argument) []Issue {

	if len(argument.Premises) == 0 {

		return []Issue{{
			RuleID:   r.ID(),
			Severity: SeverityError,
			Message:  "This argument has no premises",
			Hint:     "Add a premise",
		}}
	}
	return nil
}

func (r MissingConclusionRule) Check(argument Argument) []Issue {

	if argument.Conclusion.Text == "" {
		return []Issue{{
			RuleID:   r.ID(),
			Severity: SeverityError,
			Message:  "This argument has no conclusion",
			Hint:     "Add the conclusion",
		}}
	}
	return nil
}

func (r SinglePremiseRule) Check(argument Argument) []Issue {

	if len(argument.Premises) == 1 {

		return []Issue{{
			RuleID:   r.ID(),
			Severity: SeverityWarning,
			Message:  "Single-premise arguments are often weak",
			Hint:     "Add another premise",
		}}
	}
	return nil
}

func (r ModalityMismatchRule) Check(argument Argument) []Issue {

	if argument.Conclusion.Modality == ModalityMust {

		count := 0
		for _, p := range argument.Premises {
			if p.Confidence == High {
				count++
			}
		}
		if count == 0 {
			return []Issue{{
				RuleID:   r.ID(),
				Severity: SeverityError,
				Message:  "Strong conclusion modality (‘must’) with weak/insufficient support.",
				Hint:     "Add at least one high-confidence premise or lower the modality (‘must’ → ‘should’)",
			}}
		}
	}
	return nil
}

func (rule QuantificationRequiredRule) Check(argument Argument) []Issue {

	var issues []Issue

	premises := argument.Premises

	for i, p := range premises {

		if regexQuantificationPhrase.MatchString(p.Text) && !regexDigit.MatchString(p.Text) {

			//TODO: raise one issue per premise and list all hits
			issues = append(issues, Issue{
				RuleID:   rule.ID(),
				Severity: SeverityError,
				Message:  fmt.Sprintf("premise %d '%q' uses quantification but omits reference to actual numbers", i+1, p.Text),
				Hint:     "Provide a number (e.g., ‘18%’) or sample size supporting significant/most/increase'",
			})

		}
	}

	return issues
}

func RunAllRules(a Argument) []Issue {
	rules := []Rule{
		MissingPremiseRule{},
		VaguenessDetector{},
		MissingConclusionRule{},
		SinglePremiseRule{},
		ModalityMismatchRule{},
		QuantificationRequiredRule{},
	}
	var issues []Issue

	for _, r := range rules {
		issues = append(issues, r.Check(a)...)
	}
	return issues
}
