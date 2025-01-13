package compressor

import (
	"fmt"
	"os"
	"path/filepath"
)

func countFrequency(data []byte) map[rune]int {
	freq := make(map[rune]int)
	for _, r := range data {
		freq[rune(r)]++
	}
	return freq
}

func Compress() (string, error) {

	// Read file
	args := os.Args
	fileName := filepath.Clean(args[len(args)-1])
	data, err := os.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)

	}
	// Count frequency of characters
	countFrequency(data)

	return "", nil
}
