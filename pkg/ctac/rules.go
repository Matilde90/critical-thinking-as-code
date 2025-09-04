package ctac

type Rule interface {
	ID() string
	Check(a Argument) []Issue
}

type Issue struct {
	RuleID string
	Severity Severity
	Message string
} 

type Severity string

const (
	Info Severity = "low"
	Warning Severity = "medium"
	Error Severity = "high"
)

type MissingPremiseRule struct{}

func (r MissingPremiseRule) ID() string {
	return "CTAC001_MISSING_PREMISES"
}

func  (r MissingPremiseRule) Check (argument Argument) []Issue {

	if len(argument.Premises) == 0 {

		return []Issue{{
			RuleID: "Missing Premises",
			Severity: "Error",
			Message: "This argument has no premises",
		}}
	}
	return  nil	
}
