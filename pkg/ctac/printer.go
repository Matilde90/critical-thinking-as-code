package ctac

import (
	"fmt"
)

func SummariseArgument(argument Argument) string {

	summaryArgument := fmt.Sprintf("Title: %s\nPremises: %d\n", argument.Title, len(argument.Premises))
	for i, p := range argument.Premises {
		summaryArgument += fmt.Sprintf("P%d. %s | Confidence: %s\n", i+1, p.Text, p.Confidence)
	}

	summaryArgument += fmt.Sprintf("--------------\nConclusion: %s | Confidence: %s\n", argument.Conclusion.Text, argument.Conclusion.Confidence)
	return summaryArgument
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

func FormatIssueMessage(issues []Issue) string {

	if len(issues) == 0 {
		return "✅ No issues found.\n"
	}
	var formattedIssues string

	formattedIssues += fmt.Sprintf("Found %d issue%s:\n\n", len(issues), plural((len(issues))))
	for _, issue := range issues {
		formattedIssues += fmt.Sprintf("- %s |  %s | %s |\n", issue.RuleID, issue.Severity, issue.Message)
	}
	return formattedIssues
}

