package lang

import (
	"errors"
	"strings"
)

const LATIN1 byte = 0
const UTF16 byte = 1

const SINGLE_TOKEN byte = '?'
const MULTIPLE_TOKEN byte = '*'

var STRING_VALUE_OFFSET int64
var STRING_CODER_OFFSET int64

func init() {

}

//var NEW_STRING MethodHandle

func GetBytes(s string) ([]byte, error) {
	if s == "" {
		return nil, errors.New("the input string must not be empty")
	}
	return []byte(s), nil
}

func PadStart(str string, minLength int, padRune byte) string {
	length := len(str)
	if length >= minLength {
		return str
	}
	bytes, err := GetBytes(str)
	if err != nil {
		println(err)
	}
	dest := make([]byte, minLength)
	padLength := minLength - length
	for i := 0; i < padLength; i++ {
		dest[i] = padRune
	}
	for i := padLength; i < minLength; i++ {
		dest[i] = bytes[i-padLength]
	}
	return string(dest)
}

func IsBlank(str string) bool {
	return len(str) == 0
}

func TokenizeToStringArray(str, delimiter string) []string {
	tokens := strings.Split(str, delimiter)

	var result []string
	for _, token := range tokens {
		trimmed := strings.TrimSpace(token)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
