package parser

import "errors"

func Parse(json string) (map[string]interface{}, error) {
	if len(json) == 0 {
		return nil, errors.New("empty file")
	}
	// neither an object nor an array
	if (json[0] != '{' || json[len(json)-1] != '}') && (json[0] != '[' || json[len(json)-1] != ']') {
		return nil, errors.New("invalid json")
	}

	return nil, nil
}
