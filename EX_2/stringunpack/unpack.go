package stringunpack

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

func UnpackString(input string) (string, error) {
	var result strings.Builder
	runes := []rune(input)
	escaped := false

	for i := 0; i < len(runes); i++ {
		switch {
		case escaped:
			result.WriteRune(runes[i])
			escaped = false
		case runes[i] == '\\':
			escaped = true
		case unicode.IsDigit(runes[i]):
			if i == 0 {
				return "", errors.New("invalid string format")
			}
			repeatCount, _ := strconv.Atoi(string(runes[i]))
			if repeatCount == 0 {
				continue
			}
			lastChar := result.String()[result.Len()-1]
			result.WriteString(strings.Repeat(string(lastChar), repeatCount-1))
		default:
			result.WriteRune(runes[i])
		}
	}

	if escaped {
		return "", errors.New("invalid string format")
	}

	return result.String(), nil
}
