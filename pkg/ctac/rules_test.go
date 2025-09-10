package ctac

import (
	"testing"
)

type TestCases []struct {
	name       string
	argument   Argument
	wantIssues int
}

func TestMissingPremiseRule(t *testing.T) {

	rule := MissingPremiseRule{}
	cases := TestCases{
		{
			name: "No premise, we want it to raise one issue",
			argument: Argument{
				Title:      "Test",
				Premises:   []Premise{},
				Conclusion: Conclusion{Text: "Programming is fun"},
			},
			wantIssues: 1,
		},
		{
			name: "One premise, we want it to not raise any issue",
			argument: Argument{
				Title: "Test",
				Premises: []Premise{
					{Id: "P1", Text: "Some people like programming", Confidence: Medium},
				},
				Conclusion: Conclusion{Text: "Programming is fun"},
			},
			wantIssues: 0,
		},
	}

	for _, tc := range cases {
		tc := tc // making a copy per iteration so parallel subtests got correct case data

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			issues := rule.Check(tc.argument)
			if got := len(issues); got != tc.wantIssues {
				t.Fatalf("Testing argument %q: got %d issue%s but we wanted %d", tc.argument.Title, got, plural(got), tc.wantIssues)
			}
		})
	}
}

func TestVaguenessDetector(t *testing.T) {

	rule := VaguenessDetector{}

	cases := TestCases{{
		name: "One vague word included, one issue",
		argument: Argument{
			Title: "One vague word",
			Premises: []Premise{
				{Text: "Everyone knows that people slack off when working from home"},
			},
			Conclusion: Conclusion{
				Text: "Working from home should be banned",
			},
		},
		wantIssues: 1,
	},
		{
			name: "Two vague words included, two issues",
			argument: Argument{
				Title: "Two vague words",
				Premises: []Premise{
					{Text: "Everyone knows that it is likely that people slack off when working from home"},
				},
				Conclusion: Conclusion{
					Text: "Working from home should be banned",
				},
			},
			wantIssues: 2,
		}}

	for _, tc := range cases {

		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			issues := rule.Check(tc.argument)

			if got := len(issues); got != tc.wantIssues {
				t.Fatalf("Testing argument %q: got %d issue%s but we wanted %d", tc.argument.Title, got, plural(got), tc.wantIssues)
			}
		})

	}
}

func TestMissingConclusionRule(t *testing.T) {

	rule := MissingConclusionRule{}

	cases := TestCases{{
		name: "No conclusion included, one issue",
		argument: Argument{
			Title: "Slacking at work",
			Premises: []Premise{
				{Text: "Everyone knows that people slack off when working from home"},
			},
			Conclusion: Conclusion{
				Text: "",
			},
		},
		wantIssues: 1,
	}}
	for _, tc := range cases {

		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			issues := rule.Check(tc.argument)
			if got := len(issues); got != tc.wantIssues {
				t.Fatalf("Testing argument %q: got %d issue%s but we wanted %d", tc.argument.Title, got, plural(got), tc.wantIssues)
			}
		})

	}
}
