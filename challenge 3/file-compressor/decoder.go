package compressor

import (
	"errors"
	"strings"
)

func decodeHeader(data []byte) (*map[rune]byte, error) {
	prefixMap := make(map[rune]byte, 64)

	// Header sequence = [char1 CHARASSIGN prefix1 CHARSEP char2 CHARASSIGN prefix2 CHARSEP CHARSEP]
	foundStop := false
	for i := 0; i < len(data); i += 2 {
		// Search for STOPSEQUENCE
		if i > 0 && strings.HasPrefix(string(data[i:]), STOPSEQUENCE) {
			foundStop = true
			break
		}

		prefixMap[rune(data[i])] = data[i+1]
	}
	if !foundStop {
		return nil, errors.New("file header without a stop sequence")
	}
	return &prefixMap, nil
}
