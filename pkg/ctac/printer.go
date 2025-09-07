package ctac

import (
	"fmt"
)

func SummariseArgument (argument Argument) string {

	summaryArgument := fmt.Sprintf("Title: %s\nPremises: %d\n", argument.Title, len(argument.Premises) )
	for i, p := range argument.Premises {
		summaryArgument += fmt.Sprintf("P%d. %s | Confidence: %s\n", i+1, p.Text, p.Confidence)
	}
	return summaryArgument
}