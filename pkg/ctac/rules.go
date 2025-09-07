package ctac

import (
	"fmt"
	"strings"
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

func (r VaguenessDetector) ID() string {
	return "CTAC002_VAGUENESS_DETECTOR"
}

func (r VaguenessDetector) Check(argument Argument) []Issue {
	var issues []Issue

	premises := argument.Premises

	for i, p := range premises {
		vagueWord := []string{"someone", "probably", "likely"}
		for _, v := range vagueWord {
			if strings.Contains(p.Text, v) {
				issues = append(issues, Issue{
					RuleID:   r.ID(),
					Severity: "Error",
					Message:  fmt.Sprintf("Premise %d '%v' contains vague words '%s'", i+1, p.Text, v),
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
