package ctac

import (
	"testing"
)

func TestMissingPremiseRule(t *testing.T){

	arg1:= Argument{
		Title: "Test",
		Premises:[]Premise{
			{Id: "P1", Text: "Some people like programming", Confidence: Medium},
		},
		Conclusion: Conclusion{Text:"Programming is fun"},
	}

	arg2:= Argument{
		Title: "Test",
		Premises:[]Premise{
		},
		Conclusion: Conclusion{Text:"Programming is fun"},
	}

	issues1:=MissingPremiseRule{}.Check(arg1)
	if len(issues1) != 0 {
		t.Fatalf("esting argument %+v: got %d issue%s", arg1, len(issues1), plural(len(issues1)))
	}

		issues2:=MissingPremiseRule{}.Check(arg2)
	if len(issues2) == 0 {
		t.Fatalf("Testing argument %+v: got %d issue%s", arg2, len(issues2), plural(len(issues2)))
	}

}