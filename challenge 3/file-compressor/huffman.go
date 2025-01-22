package compressor

import (
	"encoding/json"
	"fmt"
	"slices"
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

func generatePrefixCodes(prefixMap map[rune]byte, node *TreeNode, edge byte) {
	if node == nil {
		return
	}

	if node.Rune != 0 {
		prefixMap[node.Rune] = edge
	}

	generatePrefixCodes(prefixMap, node.Left, edge*2)

	generatePrefixCodes(prefixMap, node.Right, edge*2+1)

}

func generateCodePrefixMap(prefixMap map[rune]byte) map[byte]rune {
	codePrefixMap := make(map[byte]rune, len(prefixMap))
	for r, b := range prefixMap {
		codePrefixMap[b] = r
	}
	return codePrefixMap
}
