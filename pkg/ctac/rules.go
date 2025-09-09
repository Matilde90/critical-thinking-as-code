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

func (r MissingPremiseRule) ID() string {
	return "CTAC001_MISSING_PREMISES"
}

func (rule VaguenessDetector) ID() string {
	return "CTAC002_VAGUENESS_DETECTOR"
}

func (rule VaguenessDetector) Check(argument Argument) []Issue {
	var issues []Issue

	premises := argument.Premises

	for i, p := range premises {
		type vaguePhrase struct{
			phrase string
			reg *regexp.Regexp
		}
		vaguePhrases := []vaguePhrase{
			{
			phrase: "someone",
			reg: regexp.MustCompile(`(?i)\bsomeone\b`),
		}, {
			phrase: "likely",
			reg: regexp.MustCompile(`(?i)\blikely\b`),
		},
		{
			phrase: "everyone thinks",
			reg: regexp.MustCompile(`(?i)\beveryone thinks\b`),
		},
				{
			phrase: "bprobably",
			reg: regexp.MustCompile(`(?i)\bprobably\b`),
		},
				{
			phrase: "emaybe",
			reg: regexp.MustCompile(`(?i)\bmaybe\b`),
		},
				{
			phrase: "everyone knowns",
			reg: regexp.MustCompile(`(?i)\beveryone knows\b`),
		},
	}

		for _, vaguePhrase := range vaguePhrases {
			if vaguePhrase.reg.MatchString(p.Text){
				issues = append(issues, Issue{
					RuleID:   rule.ID(),
					Severity: "Error",
					Message:  fmt.Sprintf("Premise %d '%v' contains vague words '%s'", i+1, p.Text, vaguePhrase.phrase),
				})

				break
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
