package compressor

import (
	"cmp"
	"encoding/json"
	"fmt"
	"slices"
)

func countFrequency(data []byte) map[byte]int {
	freq := make(map[byte]int)
	for _, r := range data {
		freq[r]++
	}
	return freq
}

type TreeNode struct {
	Freq  int       `json:"freq"`
	Char  byte      `json:"char"`
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

func makeTree(freq map[byte]int) *TreeNode {
	result := make([]*TreeNode, len(freq))
	i := 0
	for k, v := range freq {
		result[i] = &TreeNode{Freq: v, Char: k}
		i++
	}
	slices.SortFunc(result, func(a, b *TreeNode) int {
		if a.Freq == b.Freq {
			return cmp.Compare(a.Char, b.Char)
		}
		return a.Freq - b.Freq
	})
	for {
		if len(result) == 1 {
			break
		}
		min1 := result[0] // min1.Freq <= min2.Freq
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

type Code struct {
	Code   uint
	Length byte
}
type PrefixMap map[byte]Code

func generatePrefixCodes(prefixMap PrefixMap, node *TreeNode, edge Code) {
	if node == nil {
		return
	}

	if node.Char != 0 {
		prefixMap[node.Char] = edge
	}

	generatePrefixCodes(prefixMap, node.Left, Code{edge.Code * 2, edge.Length + 1})

	generatePrefixCodes(prefixMap, node.Right, Code{edge.Code*2 + 1, edge.Length + 1})

}

type CodePrefixMap = map[uint]map[byte]byte

func generateCodePrefixMap(prefixMap PrefixMap) CodePrefixMap {
	codePrefixMap := make(CodePrefixMap, len(prefixMap))
	for r, b := range prefixMap {
		if _, ok := codePrefixMap[b.Code]; !ok {
			codePrefixMap[b.Code] = make(map[byte]byte)
		}
		codePrefixMap[b.Code][b.Length] = r
	}
	return codePrefixMap
}
