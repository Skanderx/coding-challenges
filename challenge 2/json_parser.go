package parser

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// give the index of the first non-space character, from the start index
func skipSpaces(json string, start int) (int, error) {
	i := start
	// if !unicode.IsSpace(rune(json[i])) {
	// 	return i, nil
	// }
	for {
		if i >= len(json) {
			return -1, errors.New("invalid json")
		}
		if !unicode.IsSpace(rune(json[i])) {
			return i, nil
		}
		i++
	}
}
func findString(json string, start int) (int, string, error) {
	i := start
	if json[i] != '"' {
		return -1, "", errors.New("invalid json expected \" but got " + string(json[i]))
	}
	i++
	val := ""
	for {
		if i >= len(json) {
			return -1, "", errors.New("invalid json")
		}
		if json[i] == '"' {
			break
		}
		val += string(json[i])
		i++
	}
	return i + 1, val, nil
}
func findValue(json string, start int) (int, any, error) {
	i := start
	// Object
	if json[i] == '{' {
		closingBracketIndex := len(json) - 1
		for {
			closingBracketIndex = strings.Index(json[i:closingBracketIndex+1], "}") + i
			if closingBracketIndex == -1 {
				return -1, "", errors.New("non closing object")
			}
			if !strings.Contains(json[i+1:closingBracketIndex], "{") {
				// closing bracket found
				break
			}
		}
		objVal, err := Parse(json[i : closingBracketIndex+1])
		if err != nil {
			return -1, "", err
		}
		return closingBracketIndex + 1, objVal, nil
	}
	// Array
	if json[i] == '[' {
		if json[i+1] == ']' {
			return i + 2, make([]any, 0), nil
		}
		closingBracketIndex := len(json) - 1
		for {
			closingBracketIndex = strings.Index(json[i:closingBracketIndex+1], "]") + i
			if closingBracketIndex == -1 {
				return -1, "", errors.New("non closing array")
			}
			if !strings.Contains(json[i+1:closingBracketIndex], "[") {
				// closing bracket found
				break
			}
		}
		parts := strings.Split(json[i+1:closingBracketIndex], ",")
		arrayVal := make([]any, len(parts))
		var err error
		for index, element := range parts {
			_, arrayVal[index], err = findValue(element, 0)
			if err != nil {
				return -1, "", err
			}
		}
		return closingBracketIndex + 1, arrayVal, nil
	}
	// string value
	if json[i] == '"' {
		i++
		val := ""
		for {
			if i >= len(json) {
				return -1, "", errors.New("invalid json")
			}
			if json[i] == '"' {
				break
			}
			val += string(json[i])
			i++
		}
		return i + 1, val, nil
	}
	// boolean value
	if json[i] == 't' {
		if i+4 >= len(json) {
			return -1, "", errors.New("invalid json")
		}
		if json[i:i+4] == "true" {
			return i + 4, true, nil
		}
	}
	if json[i] == 'f' {
		if i+5 >= len(json) {
			return -1, "", errors.New("invalid json")
		}
		if json[i:i+5] == "false" {
			return i + 5, false, nil
		}
	}
	// null value
	if json[i] == 'n' {
		if i+3 >= len(json) {
			return -1, "", errors.New("invalid json")
		}
		if json[i:i+4] == "null" {
			return i + 4, nil, nil
		}
	}
	// number value
	val := ""
	for {
		if i >= len(json) {
			return -1, "", errors.New("invalid json")
		}
		if unicode.IsSpace(rune(json[i])) || json[i] == ',' || json[i] == '}' || json[i] == ']' {
			break
		}
		val += string(json[i])
		i++
	}
	if strings.Contains(val, ".") {
		fval, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return -1, "", errors.New("invalid value " + val + " : " + err.Error())
		}
		return i, fval, nil
	}
	Ival, err := strconv.Atoi(val)
	if err != nil {
		return -1, "", errors.New("invalid value " + val + " : " + err.Error())
	}
	return i, Ival, nil
}

func Parse(json string) (any, error) {
	json = strings.TrimSpace(json)
	if len(json) == 0 {
		return nil, errors.New("empty file")
	}
	if json == "{}" {
		return make(map[string]any, 0), nil
	}
	// neither an object nor an array
	if (json[0] != '{' || json[len(json)-1] != '}') && (json[0] != '[' || json[len(json)-1] != ']') {
		return nil, errors.New("invalid json")
	}

	if json[0] == '[' {
		parts := strings.Split(string(json)[1:len(json)-2], ",")
		result := make([]any, len(parts))

		var err error
		for index, element := range parts {
			_, result[index], err = findValue(element, 0)
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	}

	result := make(map[string]interface{})
	jsonLength := len(json)
	i := 0
	for i < jsonLength {

		var err error
		i, err = skipSpaces(json, i)
		if err != nil {
			return nil, err
		}
		// open object
		if json[i] == '{' {
			i++
		}
		i, err = skipSpaces(json, i)
		if err != nil {
			return nil, err
		}
		fieldName := ""
		i, fieldName, err = findString(json, i)
		if err != nil {
			return nil, err
		}
		// find :
		i, err = skipSpaces(json, i)
		if err != nil {
			return nil, err
		}
		if json[i] != ':' {
			return nil, errors.New("invalid json expected : but got " + string(json[i]))
		}
		i++
		// field value
		i, err = skipSpaces(json, i)
		if err != nil {
			return nil, err
		}
		var fieldValue any
		i, fieldValue, err = findValue(json, i)
		if err != nil {
			return nil, err
		}
		result[fieldName] = fieldValue
		// either , or }
		i, err = skipSpaces(json, i)
		if err != nil {
			return nil, err
		}
		// either , or }
		if json[i] != ',' {
			if json[i] != '}' {
				return nil, errors.New("missing }")
			}
		}
		i++
	}

	return result, nil
}
