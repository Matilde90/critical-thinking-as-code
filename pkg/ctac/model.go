package ctac

struct Argument {
	Title string
	Premises []Premise
	Conclusion string
}

struct Premise {
	Id string `yaml: "id"`
	Text string `yaml: "text"`
	Confidence Confidence `yaml: "confidence"`
}

struct Conclusion {
	Text string `yaml: text`
	Modality Modality `yaml: "modality"`
	Confidence Confidence `yaml: "confidence"`
}

type Modality string

const (
	Must Modality = "must"
	Should Modality = "should"
)

type Confidence string

const (
	Low Confidence = "low"
	Medium Confidence = "medium"
	High Confidence = "high"
)