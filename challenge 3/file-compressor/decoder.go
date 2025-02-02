package compressor

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

func decodeHeader(data []byte) (int, PrefixMap, error) {
	prefixMap := make(PrefixMap, 64)
	// Header sequence = [char1 length1 prefix1 char2 length2 prefix2 STOPSEQUENCE]
	startIndex := 0
	i := 0
	for i < len(data) {
		// Search for STOPSEQUENCE
		if i > 0 && strings.HasPrefix(string(data[i:]), STOPSEQUENCE) {
			startIndex = i + len(STOPSEQUENCE)
			break
		}
		char := data[i]
		codeLength := byte(data[i+1])
		i += 2 // read Char and Length

		var codeSequence uint

		nbrParts := codeLength / 8
		if codeLength%8 != 0 {
			nbrParts++
		}

		if int(nbrParts)+i > len(data) {
			return -1, nil, errors.New("file header is wrong")
		}

		for partIndex := byte(0); partIndex < nbrParts; partIndex++ {
			codeSequence <<= 8
			part := data[i]
			i++ // read one byte of the sequence
			codeSequence |= uint(part)
		}
		prefixMap[char] = Code{codeSequence, codeLength}
	}
	if startIndex == 0 {
		return -1, nil, errors.New("file header without a stop sequence")
	}
	return startIndex, prefixMap, nil
}

func Decompress(data []byte, f io.Writer) error {

	startIndex, prefixMap, err := decodeHeader(data)
	if err != nil {
		return fmt.Errorf("file cannot be decompressed: %w", err)
	}
	codePrefixMap := generateCodePrefixMap(prefixMap)

	w := bufio.NewWriter(f)

	var byteSeq struct {
		sequence uint
		length   byte
	}
	for _, b := range data[startIndex:] {
		// reading from left to right, bit by bit
		for j := 7; j >= 0; j-- {
			mask := byte(1) << j

			// adding b bits to byteSeq from left to right until we found a match
			// b =  aaaaaaSa and mask = 00000010 and j = 1
			// mask = ( mask & b ) >> j
			// mask = 000000S0 >> j
			// mask = 0000000S
			mask &= b
			mask >>= j
			// previousSeq = 00000abc
			// nextSeq = ( previousSeq << 1 ) | mask
			// => nextSeq = 00000abc0 | 0000000S
			// => nextSeq = 00000abcS
			byteSeq.sequence <<= 1
			byteSeq.sequence |= uint(mask)
			byteSeq.length++

			if r, ok := codePrefixMap[byteSeq.sequence][byteSeq.length]; ok {
				// sequence complete starting a new one
				w.WriteByte(r)
				byteSeq.sequence = 0
				byteSeq.length = 0
			} else if byteSeq.length == 32 {
				return errors.New("file content does not have prefix codes")
			}
		}
	}

	return w.Flush()
}
