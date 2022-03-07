package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrorIsNotStruct = errors.New("object is not struct")
	ErrorFIsNotVal   = errors.New("field is not validate")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var s string
	for _, q := range v {
		s += fmt.Sprint(q)
	}
	return s
	// return fmt.Sprint(v) // так не срабьатывает
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func isValidField(fieldVal reflect.Value, oneTagVal string) bool {
	tagValSl := strings.Split(oneTagVal, ":") // получение ключа (tagValSl[0])/значения (tagValSl[1])
	var res bool
	if fieldVal.Kind() == reflect.Slice {
		for i := 0; i < fieldVal.Len(); i++ {
			res = switchTag(fieldVal.Index(i), tagValSl[0], tagValSl[1])
			if !res {
				return false
			}
		}
	} else {
		res = switchTag(fieldVal, tagValSl[0], tagValSl[1])
	}
	return res
}

func switchTag(fieldVal reflect.Value, tagKey string, tagVal string) bool {
	switch tagKey {
	case "len":
		fieldLen, err := strconv.Atoi(tagVal)
		if err != nil {
			return false
		}
		if len(fieldVal.String()) != fieldLen {
			return false
		}
	case "in":
		slStr := strings.Split(tagVal, ",")
		var fieldValStr string
		if fieldVal.Kind() == reflect.Int {
			fieldValStr = strconv.Itoa(int(fieldVal.Int()))
		} else {
			fieldValStr = fieldVal.String()
		}

		isContain := Contains(slStr, fieldValStr)
		if !isContain {
			return false
		}
	case "max":
		tagValInt, err := strconv.Atoi(tagVal)
		if err != nil {
			return false
		}
		if fieldVal.Kind() != reflect.Int {
			return false
		}
		if int(fieldVal.Int()) > tagValInt {
			return false
		}
	case "min":
		tagValInt, err := strconv.Atoi(tagVal)
		if err != nil {
			return false
		}
		if fieldVal.Kind() != reflect.Int {
			return false
		}
		if int(fieldVal.Int()) < tagValInt {
			return false
		}
	case "regexp":
		myPattern := strings.ReplaceAll(tagVal, "\\\\", "\\")
		matched, err := regexp.MatchString(myPattern, fieldVal.String())
		if err != nil || !matched {
			return false
		}
	}
	return true
}

func Validate(v interface{}) error {
	iv := reflect.ValueOf(v)
	if iv.Kind() != reflect.Struct {
		return ErrorIsNotStruct
	}
	vOf := reflect.ValueOf(v)
	tOf := reflect.TypeOf(v)
	fCount := vOf.NumField()
	myValidationErrors := make(ValidationErrors, 0)
	for i := 0; i < fCount; i++ {
		fieldName := tOf.Field(i).Name
		fieldVal := vOf.Field(i)
		tagVal, ok := tOf.Field(i).Tag.Lookup("validate") // ищем тег  validate и полчаем его значение
		if ok {
			tagValSlAll := strings.Split(tagVal, "|") // если тег групповой
			for _, oneTagVal := range tagValSlAll {
				if !isValidField(fieldVal, oneTagVal) { // если поле не валидно
					myValidationErrors = append(myValidationErrors, ValidationError{fieldName, ErrorFIsNotVal})
				}
			}
		}
	}
	if len(myValidationErrors) > 0 {
		return myValidationErrors
	}
	return nil
}
