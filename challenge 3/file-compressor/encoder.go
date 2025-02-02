package compressor

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

const (
	STOPSEQUENCE = "%%%%"
)

// fileHeader returns the header of the compressed file
// headerformat: char1 length1 prefix1 char2 length2 prefix2 ... charn lengthn prefixn STOPSEQUENCE
func fileHeader(prefixMap PrefixMap) []byte {

	var buffer bytes.Buffer
	for k, v := range prefixMap {
		nbrParts := v.Length / 8
		if v.Length%8 != 0 {
			nbrParts++
		}
		parts := make([]byte, nbrParts)

		for i := range parts {
			// given the uneven bit sequence : XX_XXXX_YYYY_YYYY_ZZZZ_ZZZZ ; length = 22
			// we cut it into  8 bits chunks like so
			// Part i = 0 : 00XX_XXXX = byte(v.Code >> 16 )
			// Part i = 1 : YYYY_YYYY = byte(v.Code >> 8 )
			// Part i = 2 : ZZZZ_ZZZZ = byte((v.Code >> 0 ))
			// shift right length = max( 8 * (nbrParts - i - 1) , 0)
			parts[i] = byte((v.Code >> max(8*(int(nbrParts)-1-i), 0)))
		}

		buffer.WriteByte(k)
		buffer.WriteByte(v.Length)
		buffer.WriteString(string(parts))
	}
	buffer.WriteString(STOPSEQUENCE)

	return buffer.Bytes()
}

func Compress(data []byte, f io.Writer) error {

	// Count frequency of characters
	frequencies := countFrequency(data)
	// Make Huffman coding tree
	root := makeTree(frequencies)

	prefixMap := make(PrefixMap, len(frequencies))
	generatePrefixCodes(prefixMap, root, Code{0, 0})

	_, err := f.Write(fileHeader(prefixMap))
	if err != nil {
		return fmt.Errorf("error writing file header: %w", err)
	}

	w := bufio.NewWriter(f)
	// Writing 8 bits at a time; but having space for at least two
	var buffer struct {
		maxBitIndex int
		buffer      uint
	}

	for _, r := range data {
		// Adding prefix_Code in queue to be written
		// Solution without google:
		value := prefixMap[r]

		// example
		// 	buffer	= 00000aaa | ( xxxxxxxx_yyyyyyyy << 3 )
		// => buffer	= 00000aaa | xxxxx_yyyyyyyy000
		// => buffer	= xxxxx_yyyyyyyyaaa
		buffer.buffer = uint(value.Code) | buffer.buffer<<value.Length
		buffer.maxBitIndex += int(value.Length)

		// buffer	= xxxxxxxx_1yyyyyyy
		// if buffer reached this value, we at least have one byte we can write
		for buffer.maxBitIndex > 8 {

			// Write 8 bit slices at a time
			// buffer	= 000000yy_yyyyyyzz ;  maxBitIndex = 10
			// slice	= 000000yy_yyyyyyzz >> ( maxBitIndex - 7 + 1)
			// => slice	= 000000yy_yyyyyyzz >> 2
			// => slice	= 00000000Myyyyyyy : 8 bits
			slice := buffer.buffer >> (buffer.maxBitIndex - 8)
			err = w.WriteByte(byte(slice))
			if err != nil {
				return fmt.Errorf("error writing to file: %w", err)
			}
			// what should stay in buffer = 000000_000000zz = buffer % math.Pow(2, float64(maxBitIndex-7))
			buffer.buffer = buffer.buffer % (1 << (buffer.maxBitIndex - 7 + 1))
			buffer.maxBitIndex -= 8
		}
	}
	if buffer.maxBitIndex > 0 {
		// byte rest = 0000abcd // maxBitIndex = 4
		// so we shift left by (8-index) 4 bits to get abcd0000
		slice := byte(buffer.buffer) << (8 - buffer.maxBitIndex)
		err = w.WriteByte(slice)
		if err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}

	return w.Flush()
}
