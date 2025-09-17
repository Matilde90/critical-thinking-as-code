package ctac

import (
	"fmt"
	"regexp"
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
type EmotionalLanguageDetector struct{}

func (r MissingPremiseRule) ID() string {
	return "CTAC001_MISSING_PREMISES"
}

func (rule VaguenessDetector) ID() string {
	return "CTAC002_VAGUENESS_DETECTED"
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

func (rule EmotionalLanguageDetector) ID() string {
	return "CTAC007_EMOTIONAL_LANGUAGE_DETECTED"
}

type vaguePhrase struct {
	Phrase string
	Reg    *regexp.Regexp
}

var vaguePhrases = []vaguePhrase{
	{
		Phrase: "someone",
		Reg:    regexp.MustCompile(`(?i)\bsomeone\b`),
	},
	{
		Phrase: "some",
		Reg:    regexp.MustCompile(`(?i)\bsome\b`),
	},
	{
		Phrase: "everyone thinks",
		Reg:    regexp.MustCompile(`(?i)\beveryone thinks\b`),
	},
	{
		Phrase: "maybe",
		Reg:    regexp.MustCompile(`(?i)\bmaybe\b`),
	},
	{
		Phrase: "everyone knows",
		Reg:    regexp.MustCompile(`(?i)\beveryone knows\b`),
	},
}

type emotionalLanguagePhrase struct {
	Phrase string
	Reg    *regexp.Regexp
}

var negativeEmotionWords = []string{"terrible", "horrible", "disastrous", "catastrophic", "evil", "awful", "tragic", "shocking", "outrageous"}
var positiveEmotionWords = []string{"amazing", "brilliant", "fantastic", "heroic", "wonderful", "incredible", "terrific"}
var persuasiveIntensifiers = []string{"obviously", "clearly", "undeniably", "absolutely", "definitively"}

func buildPhrases(words []string) []emotionalLanguagePhrase {
	phrases := make([]emotionalLanguagePhrase, 0, len(words))
	for _, w := range words {
		pattern := fmt.Sprintf(`(?i)\b%s`, w)
		phrases = append(phrases, emotionalLanguagePhrase{
			Phrase: w,
			Reg:    regexp.MustCompile(pattern),
		})
	}
	return phrases
}

var negativePhrases = buildPhrases(negativeEmotionWords)
var positivePhrases = buildPhrases(positiveEmotionWords)
var intensifierPhrases = buildPhrases(persuasiveIntensifiers)

var regexDigit = regexp.MustCompile("[0-9]+")
var regexQuantificationPhrase = regexp.MustCompile(`(?i)(\bsignificant|\bdecrease|\bmost\b|\bincrease|\bdecline\b|\bpercent(age?)\b|%|\bmore\b|\bless\b|\brate\b|\btrend\b)`)

func (rule VaguenessDetector) Check(argument Argument) []Issue {
	var issues []Issue

	premises := argument.Premises
	var spottedVagueWords string

	for _, p := range premises {

		for _, vaguePhrase := range vaguePhrases {
			if vaguePhrase.Reg.MatchString(p.Text) {

				spottedVagueWords = spottedVagueWords + ", " + vaguePhrase.Phrase
			}
		}
		if len(spottedVagueWords) > 0 {

			issues = append(issues, Issue{
				RuleID:   rule.ID(),
				Severity: SeverityWarning,
				Message:  fmt.Sprintf("Premise %s %q contains vague words '%s'", p.Id, p.Text, strings.TrimLeft(spottedVagueWords, " ,")),
				Hint:     "Remove use of vague words by using more precise language",
			})
		}

		spottedVagueWords = ""

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

	for _, p := range premises {

		if regexQuantificationPhrase.MatchString(p.Text) && !regexDigit.MatchString(p.Text) {

			issues = append(issues, Issue{
				RuleID:   rule.ID(),
				Severity: SeverityError,
				Message:  fmt.Sprintf("Premise %s '%q' uses quantification but omits reference to actual numbers", p.Id, p.Text),
				Hint:     "Provide a number (e.g., ‘18%’) or sample size supporting significant/most/increase'",
			})
		}
	}

	return issues
}

func (rule EmotionalLanguageDetector) Check(argument Argument) []Issue {

	var issues []Issue

	premises := argument.Premises
	var spottedEmotionalWords string

	for _, p := range premises {

		for _, np := range negativePhrases {
			if np.Reg.MatchString(p.Text) {
				spottedEmotionalWords = spottedEmotionalWords + ", " + np.Phrase
			}
		}

		for _, pp := range positivePhrases {
			if pp.Reg.MatchString(p.Text) {
				spottedEmotionalWords = spottedEmotionalWords + ", " + pp.Phrase
			}
		}

		for _, ip := range intensifierPhrases {
			if ip.Reg.MatchString(p.Text) {
				spottedEmotionalWords = spottedEmotionalWords + ", " + ip.Phrase
			}
		}

		if len(spottedEmotionalWords) > 0 {
			issues = append(issues, Issue{
				RuleID:   rule.ID(),
				Severity: SeverityError,
				Message:  fmt.Sprintf("Premise %s '%q' uses emotional language %s", p.Id, p.Text, strings.TrimLeft(spottedEmotionalWords, " ,")),
				Hint:     "Please rewrite the premises without using unnecessary emotional language'",
			})
		}

		spottedEmotionalWords = ""

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
		EmotionalLanguageDetector{},
	}
	var issues []Issue

	for _, r := range rules {
		issues = append(issues, r.Check(a)...)
	}
	return issues
}
