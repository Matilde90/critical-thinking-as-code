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
}

type Severity string

const (
	Info    Severity = "low"
	Warning Severity = "medium"
	Error   Severity = "high"
)

type MissingPremiseRule struct{}
type VaguenessDetector struct{}
type MissingConclusionRule struct{}
type SinglePremiseRule struct{}
type ModalityMismatchRule struct{}

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
	return "CTA004_SINGLE_PREMISE_RULE"
}

func (rule ModalityMismatchRule) ID() string {
	return "CTA005_MODALITY_MISMATCH_RULE"
}

func (rule VaguenessDetector) Check(argument Argument) []Issue {
	var issues []Issue

	premises := argument.Premises

	for i, p := range premises {
		type vaguePhrase struct {
			phrase string
			reg    *regexp.Regexp
		}
		vaguePhrases := []vaguePhrase{
			{
				phrase: "someone",
				reg:    regexp.MustCompile(`(?i)\bsomeone\b`),
			}, {
				phrase: "likely",
				reg:    regexp.MustCompile(`(?i)\blikely\b`),
			},
			{
				phrase: "everyone thinks",
				reg:    regexp.MustCompile(`(?i)\beveryone thinks\b`),
			},
			{
				phrase: "probably",
				reg:    regexp.MustCompile(`(?i)\bprobably\b`),
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

		for _, vaguePhrase := range vaguePhrases {
			if vaguePhrase.reg.MatchString(p.Text) {
				issues = append(issues, Issue{
					RuleID:   rule.ID(),
					Severity: "Error",
					Message:  fmt.Sprintf("Premise %d '%v' contains vague words '%s'", i+1, p.Text, vaguePhrase.phrase),
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
			Severity: "Error",
			Message:  "This argument has no premises",
		}}
	}
	return nil
}

func (r MissingConclusionRule) Check(argument Argument) []Issue {

	if argument.Conclusion.Text == "" {
		return []Issue{{
			RuleID:   r.ID(),
			Severity: "Error",
			Message:  "This argument has no conclusion",
		}}
	}
	return nil
}

func (r SinglePremiseRule) Check(argument Argument) []Issue {

	if len(argument.Premises) == 1 {

		return []Issue{{
			RuleID:   r.ID(),
			Severity: "Warning",
			Message:  "Single-premise arguments are often weak",
		}}
	}
	return nil
}

func (r ModalityMismatchRule) Check(argument Argument) []Issue {

	if argument.Conclusion.Confidence == "high" && argument.Conclusion.Modality == "must" {

		count := 0
		for _, p := range argument.Premises {
			if p.Confidence == "high" {
				count++
			}
		}
		if count == 0 {
			return []Issue{{
				RuleID:   r.ID(),
				Severity: "Error",
				Message:  "Modality mismatch",
			}}
		}
	}
	return nil
}

func RunAllRules(a Argument) []Issue {
	rules := []Rule{
		MissingPremiseRule{},
		VaguenessDetector{},
	}
	var issues []Issue

	for _, r := range rules {
		issues = append(issues, r.Check(a)...)
	}
	return issues
}
