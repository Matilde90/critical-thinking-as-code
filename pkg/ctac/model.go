package ctac

type Modality string

const (
	ModalityMust Modality = "must"
	ModalityShould Modality = "should"
)

type Confidence string

const (
	Low Confidence = "low"
	Medium Confidence = "medium"
	High Confidence = "high"
)

type Argument struct {
	Title string `yaml:"title"`
	Premises []Premise `yaml:"premises"`
	Conclusion Conclusion `yaml:"conclusion"`
}

type Premise struct {
	Id string `yaml:"id"`
	Text string `yaml:"text"`
	Confidence Confidence `yaml:"confidence"`
}

type Conclusion struct {
	Text string `yaml:"text"`
	Modality Modality `yaml:"modality"`
	Confidence Confidence `yaml:"confidence"`
}
