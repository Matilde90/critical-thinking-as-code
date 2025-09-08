package ctac

import (
	"testing"
)

func TestMissingPremiseRule(t *testing.T) {

	rule := MissingPremiseRule{}
	cases := []struct {
		name       string
		argument   Argument
		wantIssues int
	}{
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
