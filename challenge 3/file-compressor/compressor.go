package compressor

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/exp/slices"
)

func countFrequency(data []byte) map[rune]int {
	freq := make(map[rune]int)
	for _, r := range data {
		freq[rune(r)]++
	}
	return freq
}

type TreeNode struct {
	Freq  int       `json:"freq"`
	Rune  rune      `json:"char"`
	Left  *TreeNode `json:"left"`
	Right *TreeNode `json:"right"`
}

func (t *TreeNode) String() string {
	result, err := json.Marshal(t)
	if err != nil {
		fmt.Println("cannot convert tree", err)
		return ""
	}
	return string(result)
}
func makeTree(freq map[rune]int) *TreeNode {
	result := make([]*TreeNode, len(freq))
	i := 0
	for k, v := range freq {
		result[i] = &TreeNode{Freq: v, Rune: k}
		i++
	}
	slices.SortFunc(result, func(a, b *TreeNode) int {
		return a.Freq - b.Freq
	})
	for {
		if len(result) == 1 {
			break
		}
		min1 := result[0]
		min2 := result[1]
		result = result[2:]
		joinedNode := &TreeNode{
			Freq:  min1.Freq + min2.Freq,
			Left:  min1,
			Right: min2,
		}
		insertionIndex := slices.IndexFunc(result, func(a *TreeNode) bool {
			return a.Freq > joinedNode.Freq
		})
		if insertionIndex == -1 {
			insertionIndex = len(result)
		}
		result = slices.Insert(result, insertionIndex, joinedNode)
	}

	return result[0]
}

func generatePrefixCodes(prefixMap *map[rune]byte, node *TreeNode, edge byte) {
	if node == nil {
		return
	}

	if node.Rune != 0 {
		(*prefixMap)[node.Rune] = edge
	}

	generatePrefixCodes(prefixMap, node.Left, edge*2)

	generatePrefixCodes(prefixMap, node.Right, edge*2+1)

}

// fileHeader returns the header of the compressed file
// headerformat: char1:prefix1,char2:prefix2,...,charn:prefixn,,
func fileHeader(prefixMap *map[rune]byte) []byte {

	var buffer bytes.Buffer
	for k, v := range *prefixMap {
		buffer.WriteRune(k)
		buffer.WriteRune(':')
		buffer.WriteByte(v)
		buffer.WriteRune(',')
	}
	buffer.WriteRune(',')

	return buffer.Bytes()
}

func Compress(fileName string) (string, error) {

	data, err := os.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)

	}
	// Count frequency of characters
	ferequencies := countFrequency(data)
	// Make Huffman coding tree
	root := makeTree(ferequencies)

	prefixMap := make(map[rune]byte, len(ferequencies))
	generatePrefixCodes(&prefixMap, root, 0)

	f, err := os.Create("./" + fileName + "_compressed.txt")
	if err != nil {
		return "", fmt.Errorf("error creating file: %w", err)
	}
	defer f.Close()
	_, err = f.Write(fileHeader(&prefixMap))
	if err != nil {
		return "", fmt.Errorf("error writing file header: %w", err)
	}

	w := bufio.NewWriter(f)
	// Writing 8 bits at a time;
	var buffer uint16

	for _, r := range data {
		if buffer >= 255 {
			// buffer	= 0bxxxxxxxx|yyyyyyyy
			// slice	= 0byyyyyyyy => Write 8 bits at a time
			// buffer	= 0bxxxxxxxx => leave the rest until we have 8 other bits set
			slice := byte(buffer - buffer%255)
			buffer = buffer % 255
			err = w.WriteByte(slice)
			if err != nil {
				return "", fmt.Errorf("error writing to file: %w", err)
			}
		}
		// Adding prefix_Code in queue to be written
		// Solution without google:
		value := prefixMap[rune(r)]
		switch {
		case value >= 128:
			buffer = uint16(prefixMap[rune(r)]) | buffer<<8
		case value >= 64:
			buffer = uint16(prefixMap[rune(r)]) | buffer<<7
		case value >= 32:
			buffer = uint16(prefixMap[rune(r)]) | buffer<<6
		case value >= 16:
			buffer = uint16(prefixMap[rune(r)]) | buffer<<5
		case value >= 8:
			buffer = uint16(prefixMap[rune(r)]) | buffer<<4
		case value >= 4:
			buffer = uint16(prefixMap[rune(r)]) | buffer<<3
		case value >= 2:
			buffer = uint16(prefixMap[rune(r)]) | buffer<<2
		default:
			buffer = uint16(prefixMap[rune(r)]) | buffer<<1
		}
	}
	if buffer > 0 {
		slice := byte(buffer >> 8)
		err = w.WriteByte(slice)
		if err != nil {
			return "", fmt.Errorf("error writing to file: %w", err)
		}
	}

	w.Flush()
	return "", nil
}
