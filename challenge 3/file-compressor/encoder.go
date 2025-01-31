package compressor

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

const (
	STOPSEQUENCE = "%%%%"
)

// fileHeader returns the header of the compressed file
// headerformat: char1 prefix1 char2 prefix2 ... charn prefixn STOPSEQUENCE
func fileHeader(prefixMap map[byte]byte) []byte {

	var buffer bytes.Buffer
	for k, v := range prefixMap {
		buffer.WriteByte(k)
		buffer.WriteByte(v)
	}
	buffer.WriteString(STOPSEQUENCE)

	return buffer.Bytes()
}

var Pow2_16 = map[int]uint16{
	0:  1,
	1:  2,
	2:  4,
	3:  8,
	4:  16,
	5:  32,
	6:  64,
	7:  128,
	8:  256,
	9:  512,
	10: 1024,
	11: 2048,
	12: 4096,
	13: 8102,
	14: 16204,
	15: 32408,
	16: 64816,
}

func Compress(fileName string) error {

	data, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)

	}
	// Count frequency of characters
	frequencies := countFrequency(data)
	// Make Huffman coding tree
	root := makeTree(frequencies)

	prefixMap := make(map[byte]byte, len(frequencies))
	generatePrefixCodes(prefixMap, root, 0)

	f, err := os.Create("./" + fileName + "_compressed.txt")
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer f.Close()
	_, err = f.Write(fileHeader(prefixMap))
	if err != nil {
		return fmt.Errorf("error writing file header: %w", err)
	}

	w := bufio.NewWriter(f)
	// Writing 8 bits at a time; but having space for at least two
	var buffer struct {
		maxBitIndex int
		buffer      uint16
	}

	for _, r := range data {
		// Adding prefix_Code in queue to be written
		// Solution without google:
		value := prefixMap[r]
		switch {
		// example
		// 	buffer	= 00000aaa | ( xxxxxxxx_yyyyyyyy << 3 )
		// => buffer	= 00000aaa | xxxxx_yyyyyyyy000
		// => buffer	= xxxxx_yyyyyyyyaaa
		case value >= 128: // 8 bit value
			buffer.buffer = uint16(value) | buffer.buffer<<8
			buffer.maxBitIndex += 8
		case value >= 64: // 7 bit value
			buffer.buffer = uint16(value) | buffer.buffer<<7
			buffer.maxBitIndex += 7
		case value >= 32: // 6 bit value
			buffer.buffer = uint16(value) | buffer.buffer<<6
			buffer.maxBitIndex += 6
		case value >= 16:
			buffer.buffer = uint16(value) | buffer.buffer<<5
			buffer.maxBitIndex += 5
		case value >= 8:
			buffer.buffer = uint16(value) | buffer.buffer<<4
			buffer.maxBitIndex += 4
		case value >= 4:
			buffer.buffer = uint16(value) | buffer.buffer<<3
			buffer.maxBitIndex += 3
		case value >= 2:
			buffer.buffer = uint16(value) | buffer.buffer<<2
			buffer.maxBitIndex += 2
		default:
			buffer.buffer = uint16(value) | buffer.buffer<<1
			buffer.maxBitIndex += 1
		}

		// buffer	= xxxxxxxx_1yyyyyyy
		// if buffer reached this value, we at least have one byte we can write
		if buffer.maxBitIndex > 8 {

			// Write 8 bit slices at a time
			// buffer	= 000000yy_yyyyyyzz ;  maxBitIndex = 10
			// slice	= 000000yy_yyyyyyzz >> ( maxBitIndex - 7 + 1)
			// => slice	= 000000yy_yyyyyyzz >> 2
			// => slice	= 00000000Myyyyyyy : 8 bits
			slice := buffer.buffer >> (buffer.maxBitIndex - 7 + 1)
			err = w.WriteByte(byte(slice))
			if err != nil {
				return fmt.Errorf("error writing to file: %w", err)
			}
			// what should stay in buffer = 000000_000000zz = buffer % math.Pow(2, float64(maxBitIndex-7))
			buffer.buffer = buffer.buffer % Pow2_16[(buffer.maxBitIndex-7+1)]
			buffer.maxBitIndex -= 8
		}
	}
	if buffer.maxBitIndex > 0 {
		slice := byte(buffer.buffer >> (buffer.maxBitIndex - 7 + 1))
		err = w.WriteByte(slice)
		if err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}

	return w.Flush()
}
