package compressor

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func decodeHeader(data []byte) (int, map[byte]byte, error) {
	prefixMap := make(map[byte]byte, 64)
	// Header sequence = [char1 prefix1 char2 prefix2 STOPSEQUENCE]
	startIndex := 0
	for i := 0; i < len(data); i += 2 {
		// Search for STOPSEQUENCE
		if i > 0 && strings.HasPrefix(string(data[i:]), STOPSEQUENCE) {
			startIndex = i + 4
			break
		}

		prefixMap[data[i]] = data[i+1]
	}
	if startIndex == 0 {
		return -1, nil, errors.New("file header without a stop sequence")
	}
	return startIndex, prefixMap, nil
}

var Pow2 = map[int]byte{
	0: 1,
	1: 2,
	2: 4,
	3: 8,
	4: 16,
	5: 32,
	6: 64,
	7: 128,
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
		for j := 0; j < 8; j++ {
			// previousSeq = o******** : o bit will overflow
			if byteSeq > 128 {
				return errors.New("file content does not have prefix codes")
			}
			mask = Pow2[7-j]
			// adding b bits to byteSeq from left to right until we found a match
			// b =  aaaaaaSa and mask = 00000010 and j = 6
			// mask = ( mask & b ) >> ( 7 - j )
			// mask = 000000S0 >> ( 7 - j )
			// mask = 0000000S
			mask &= b
			mask >>= (7 - j)
			// previousSeq = 00000abc
			// nextSeq = ( previousSeq << 1 ) | mask
			// => nextSeq = 00000abc0 | 0000000S
			// => nextSeq = 00000abcS
			byteSeq <<= 1
			byteSeq |= mask

			if r, ok := codePrefixMap[byteSeq]; ok {
				// sequence complete starting a new one
				w.WriteByte(r)
				byteSeq = 0
			}
		}
	}

	return w.Flush()
}
