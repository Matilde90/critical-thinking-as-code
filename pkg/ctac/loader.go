package ctac

import (
	"os"

	"gopkg.in/yaml.v3"
)

func Loader(filePath string) (*Argument, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	argument := Argument{}
	err = yaml.Unmarshal(data, &argument)
	if err != nil {
		return nil, err
	}

	return &argument, err
}
