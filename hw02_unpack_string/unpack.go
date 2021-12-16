package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var newStr string
	var lastR rune
	var isDigit bool
	var lastIsDigit bool
	//	fmt.Println(len(str))
	//	fmt.Println(utf8.RuneCountInString(str))
	for _, r := range str {
		isDigit = unicode.IsDigit(r)
		if lastR == 0 && isDigit {
			return "", ErrInvalidString
		}
		if isDigit && lastIsDigit {
			return "", ErrInvalidString
		}
		if isDigit {
			i, err := strconv.Atoi(string(r))
			if err != nil {
				return "", err
			}
			newStr += strings.Repeat(string(lastR), i)
		} else if !lastIsDigit && lastR > 0 {
			newStr += string(lastR)
		}
		lastR = r
		lastIsDigit = isDigit
	}
	if !lastIsDigit && lastR > 0 {
		newStr += string(lastR)
	}

	return newStr, nil
}
