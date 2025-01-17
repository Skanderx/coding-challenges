package compressor

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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

func Compress() (string, error) {

	// Read file
	args := os.Args
	fileName := filepath.Clean(args[len(args)-1])
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

	for _, r := range data {
		err = w.WriteByte(prefixMap[rune(r)])
		if err != nil {
			return "", fmt.Errorf("error writing to file: %w", err)
		}
	}

	w.Flush()
	return "", nil
}
