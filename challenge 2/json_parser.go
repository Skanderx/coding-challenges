package parser

import (
	"errors"
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
func findValue(json string, start int) (int, string, error) {
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

func Parse(json string) (map[string]interface{}, error) {
	json = strings.TrimSpace(json)
	if len(json) == 0 {
		return nil, errors.New("empty file")
	}
	// neither an object nor an array
	if (json[0] != '{' || json[len(json)-1] != '}') && (json[0] != '[' || json[len(json)-1] != ']') {
		return nil, errors.New("invalid json")
	}

	result := make(map[string]interface{})
	jsonLength := len(json)
	i := 0
	for i < jsonLength {

		// open object
		if json[i] == '{' {
			i++
		}
		var err error
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
		fieldValue := ""
		i, fieldValue, err = findValue(json, i) // TODO: test field name isn't empty
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
				return nil, err
			}
		}
		i++
	}

	return result, nil
}
