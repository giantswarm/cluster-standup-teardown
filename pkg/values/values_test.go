package values

import (
	"strings"
	"testing"
)

func TestMerge(t *testing.T) {
	tests := []struct {
		name        string
		inputValues []string
		expected    string
	}{
		{
			name: "Single input",
			inputValues: []string{
				`foo: bar`,
			},
			expected: `foo: bar`,
		},
		{
			name: "Single input, nested",
			inputValues: []string{
				`foo:
  bar: baz`,
			},
			expected: `foo:
  bar: baz`,
		},
		{
			name: "Override top lever",
			inputValues: []string{
				`foo: bar`,
				`foo: baz`,
			},
			expected: `foo: baz`,
		},
		{
			name: "Override nested",
			inputValues: []string{
				`foo:
  bar: baz`,
				`foo:
  bar: 123`,
			},
			expected: `foo:
  bar: 123`,
		},
		{
			name: "Add nested",
			inputValues: []string{
				`foo:
  bar: baz`,
				`foo:
  extra: property`,
			},
			expected: `foo:
  bar: baz
  extra: property`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actualResult, err := Merge(tc.inputValues...)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if strings.TrimSpace(actualResult) != strings.TrimSpace(tc.expected) {
				t.Fatalf("Actual value didn't match expected value\n\nexpected: %q\n\nactual: %q", strings.TrimSpace(tc.expected), strings.TrimSpace(actualResult))
			}
		})
	}
}
