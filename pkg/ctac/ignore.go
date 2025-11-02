package ctac

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type IgnoreSpec struct {
	Rules  []string `yaml:"rules" json:"rules"`
	Reason []string `yaml:"reason" json:"reason"`
}

func resolveIgnorePath(userPath string) string {
	if userPath != "" {
		return userPath
	}

	for _, defaultIgnoreFilePath := range []string{
		"ctac.ignore.yaml",
		"ctacignore.yaml",
		"ctacIgnore.yaml",
		"ctac.ignore.yml",
		"ctacignore.yml",
		"ctacIgnore.yml",
	} {
		if _, err := os.Stat(defaultIgnoreFilePath); err == nil {
			return defaultIgnoreFilePath
		}
	}
	return ""
}

func LoadIgnore(filePath string) (*IgnoreSpec, error) {
	ignoreFilePath := resolveIgnorePath(filePath)
	ignoreSpec := IgnoreSpec{}
	if ignoreFilePath == "" {
		return &ignoreSpec, nil
	}
	if _, err := os.Stat(ignoreFilePath); err != nil {
		return nil, fmt.Errorf("no such file or directory for ignore file at %s. Please provide a valid path to your ignore file", ignoreFilePath)
	}
	data, err := os.ReadFile(ignoreFilePath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &ignoreSpec)
	if err != nil {
		return nil, err
	}
	return &ignoreSpec, err
}
