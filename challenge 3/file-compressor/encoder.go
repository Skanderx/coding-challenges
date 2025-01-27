package compressor

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
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
	var buffer uint16

	for _, r := range data {

		// buffer	= xxxxxxxx_1yyyyyyy
		// if buffer reached this value, we at least have one byte we can write
		if buffer > 128 {
			maxBitIndex := 8
			for {
				if maxBitIndex < 15 && buffer > uint16(math.Pow(2, float64(maxBitIndex)+1)) {
					maxBitIndex = maxBitIndex + 1
				} else {
					break
				}
			}
			// Write 8 bit slices at a time
			// buffer	= 000000yy_yyyyyyzz ;  maxBitIndex = 9
			// slice	= ( 000000yy_yyyyyyzz >> ( maxBitIndex - 7 )
			// => slice	= (xxxxxxMy_yyyyyyzz >> 2
			// => slice	= 00000000Myyyyyyy : 8 bits
			slice := buffer
			slice >>= (maxBitIndex - 7)
			// what should stay of buffer = 000000_000000zz = buffer % math.Pow(2, float64(maxBitIndex-7))
			buffer = buffer % uint16(math.Pow(2, float64(maxBitIndex-7)))
			err = w.WriteByte(byte(slice))
			if err != nil {
				return fmt.Errorf("error writing to file: %w", err)
			}
		}
		// Adding prefix_Code in queue to be written
		// Solution without google:
		value := prefixMap[r]
		switch {
		// example
		// 	buffer	= 00000aaa | ( xxxxxxxx_yyyyyyyy << 3 )
		// => buffer	= 00000aaa | xxxxx_yyyyyyyy000
		// => buffer	= xxxxx_yyyyyyyyaaa
		case value >= 128: // 8 bit value
			buffer = uint16(prefixMap[r]) | buffer<<8
		case value >= 64: // 7 bit value
			buffer = uint16(prefixMap[r]) | buffer<<7
		case value >= 32: // 6 bit value
			buffer = uint16(prefixMap[r]) | buffer<<6
		case value >= 16:
			buffer = uint16(prefixMap[r]) | buffer<<5
		case value >= 8:
			buffer = uint16(prefixMap[r]) | buffer<<4
		case value >= 4:
			buffer = uint16(prefixMap[r]) | buffer<<3
		case value >= 2:
			buffer = uint16(prefixMap[r]) | buffer<<2
		default:
			buffer = uint16(prefixMap[r]) | buffer<<1
		}
	}
	if buffer > 0 {
		slice := byte(buffer >> 8)
		err = w.WriteByte(slice)
		if err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}

	return w.Flush()
}
