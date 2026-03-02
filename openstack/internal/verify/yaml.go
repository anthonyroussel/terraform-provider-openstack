package verify

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

func NormalizeYAMLString(templateString any) (string, error) {
	if v, ok := templateString.(string); ok {
		templateString = strings.ReplaceAll(v, "\r\n", "\n")
	}

	return checkYAMLString(templateString)
}

func ValidStringIsYAML(v any, k string) (ws []string, errors []error) {
	if _, err := checkYAMLString(v); err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid YAML: %w", k, err))
	}

	return
}

// Takes a value containing YAML string and passes it through
// the YAML parser. Returns either a parsing
// error or original YAML string.
func checkYAMLString(yamlString any) (string, error) {
	var y any

	if yamlString == nil || yamlString.(string) == "" {
		return "", nil
	}

	s := yamlString.(string)

	err := yaml.Unmarshal([]byte(s), &y)

	return s, err
}
