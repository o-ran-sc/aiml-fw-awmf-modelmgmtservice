package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncrementArtifactVersion(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expected       string
		expected_error string
	}{
		{"InitialVersion", "0.0.0", "1.0.0", ""},
		{"MinorIncrement", "1.0.0", "1.1.0", ""},
		{"AnotherIncrement", "1.5.0", "1.6.0", ""},
		{"InvalidFormat", "a.b.c", "", "failed to parse artifactVersion numbers: strconv.Atoi: parsing \"a\": invalid syntax, strconv.Atoi: parsing \"b\": invalid syntax, strconv.Atoi: parsing \"c\": invalid syntax"},
		{"InvalidFormat", "1.0", "", "invalid artifactVersion format: 1.0"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := IncrementArtifactVersion(tc.input)
			var err_str string
			if err == nil {
				err_str = ""
			} else {
				err_str = err.Error()
			}

			assert.Equal(t, tc.expected_error, err_str)
			assert.Equal(t, tc.expected, got)
		})
	}
}
