package values

import (
	"os"

	"dario.cat/mergo"
	"sigs.k8s.io/yaml"
)

// MustOverlayValues performs an OverlayValues call but ignores any errors that occur while reading the values file.
func MustMergeValues(values ...string) string {
	finalValues, _ := Merge(append([]string{BuildBaseValues()}, values...)...)
	return finalValues
}

func Merge(layers ...string) (string, error) {
	mergedLayers := map[string]interface{}{}

	for _, layer := range layers {
		if layer == "" {
			continue
		}

		var rawMapData map[string]interface{}
		err := yaml.Unmarshal([]byte(layer), &rawMapData)
		if err != nil {
			return "", err
		}

		err = mergo.Merge(&mergedLayers, rawMapData, mergo.WithOverride)
		if err != nil {
			return "", err
		}
	}

	data, err := yaml.Marshal(mergedLayers)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// MustLoadValuesFile attempts to load a values file from the provided filePath and if fails returns an empty string
func MustLoadValuesFile(filePath string) string {
	fileBytes, err := os.ReadFile(filePath) // nolint:gosec
	if err != nil {
		return ""
	}
	return string(fileBytes)
}
