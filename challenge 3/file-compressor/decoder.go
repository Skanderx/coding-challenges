package compressor

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
)

func decodeHeader(data []byte) (int, map[rune]byte, error) {
	prefixMap := make(map[rune]byte, 64)
	// Header sequence = [char1 prefix1 char2 prefix2 STOPSEQUENCE]
	startIndex := 0
	for i := 0; i < len(data); i += 2 {
		// Search for STOPSEQUENCE
		if i > 0 && strings.HasPrefix(string(data[i:]), STOPSEQUENCE) {
			startIndex = i + 4
			break
		}

		prefixMap[rune(data[i])] = data[i+1]
	}
	if startIndex == 0 {
		return -1, nil, errors.New("file header without a stop sequence")
	}
	return startIndex, prefixMap, nil
}

func Decompress(fileName string) error {

	data, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)

	}
	startIndex, prefixMap, err := decodeHeader(data)
	if err != nil {
		return fmt.Errorf("file cannot be decompressed: %w", err)
	}
	codePrefixMap := generateCodePrefixMap(prefixMap)

	f, err := os.Create("./" + fileName + "decompressed.txt")
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	var byteSeq byte
	for _, b := range data[startIndex:] {
		var mask byte
		for j := range 7 {
			// previousSeq = o******** : o bit will overflow
			if byteSeq/128 == 1 {
				return errors.New("file content does not have prefix codes")
			}
			mask = byte(math.Pow(2, float64(7-j)))
			// previousSeq = *****abc
			// b =  aaaaaaSa and mask = 00000010
			// nextSeq = *****abc << 1 | 000000S0 >> j (=1)
			// => nextSeq = ****abc0 | 0000000S
			// => nextSeq = ****abcS
			byteSeq <<= 1
			mask &= b >> j
			byteSeq = byteSeq | mask

			if r, ok := codePrefixMap[byteSeq]; ok {
				// sequence complete starting a new one
				w.WriteRune(r)
				byteSeq = 0
			}
		}
	}

	return w.Flush()
}
