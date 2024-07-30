package stringunpack

import (
	"testing"
)

func TestUnpackString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"a4bc2d5e", "aaaabccddddde", false},
		{"abcd", "abcd", false},
		{"45", "", true},
		{"", "", false},
		{"qwe\\4\\5", "qwe45", false},
		{"qwe\\45", "qwe44444", false},
		{"qwe\\\\5", "qwe\\\\\\\\\\", false},
		{"abc\\", "", true},
	}

	for _, test := range tests {
		result, err := UnpackString(test.input)
		if test.hasError {
			if err == nil {
				t.Errorf("expected an error for input %s", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("did not expect an error for input %s, but got %v", test.input, err)
			}
			if result != test.expected {
				t.Errorf("expected %s for input %s, but got %s", test.expected, test.input, result)
			}
		}
	}
}