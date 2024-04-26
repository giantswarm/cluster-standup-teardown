package values

import "os"

// OverlayValues takes in a default values string and an optional path to a file containing values
// If the values file if found it will be loaded and used for the values, otherwise the default is used.
func OverlayValues(defaultValues string, valuesFile string) (string, error) {
	finalValues := defaultValues
	if valuesFile != "" {
		fileBytes, err := os.ReadFile(valuesFile)
		if err != nil && !os.IsNotExist(err) {
			return finalValues, err
		}
		// TODO: Override / merge values together
		finalValues = string(fileBytes)
	}

	return finalValues, nil
}

// MustOverlayValues performs an OverlayValues call but ignores any errors that occur while reading the values file.
func MustOverlayValues(defaultValues string, valuesFile string) string {
	finalValues, _ := OverlayValues(defaultValues, valuesFile)
	return finalValues
}
