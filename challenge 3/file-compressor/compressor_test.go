package compressor

import (
	"fmt"
	"maps"
	"os"
	"testing"
)

const TESTFILEPATH = "../../challenge 3/test.txt"

func TestCountFrequency(t *testing.T) {
	data, err := os.ReadFile(TESTFILEPATH)
	if err != nil {
		t.Errorf("error reading file: %v\n", err)
		return
	}
	freq := countFrequency(data)
	t.Run("X occurences on test file", func(t *testing.T) {
		if freq['X'] != 333 {
			t.Errorf("expected 333 occurence of X, got %d\n", freq['X'])
		} else {
			fmt.Println("X occurences on test file correct")
		}
	})
	t.Run("Y occurences on test file", func(t *testing.T) {
		if freq['t'] != 223000 {
			t.Errorf("expected 223000 occurence of t, got %d\n", freq['Y'])
		} else {
			fmt.Println("t occurences on test file correct")
		}
	})
}

var freqExample = map[rune]int{
	'C': 32,
	'D': 42,
	'E': 120,
	'K': 7,
	'L': 42,
	'M': 24,
	'U': 37,
	'Z': 2,
}
var referenceTree = &TreeNode{
	Freq: 306, Rune: 0,
	Left: &TreeNode{
		Freq: 120, Rune: 'E',
		Left: nil, Right: nil,
	},
	Right: &TreeNode{
		Freq: 186, Rune: 0,
		Left: &TreeNode{
			Freq: 79, Rune: 0,
			Left: &TreeNode{
				Freq: 37, Rune: 'U',
				Left: nil, Right: nil,
			},
			Right: &TreeNode{
				Freq: 42, Rune: 'D',
				Left: nil, Right: nil,
			},
		},
		Right: &TreeNode{
			Freq: 107, Rune: 0,
			Left: &TreeNode{
				Freq: 42, Rune: 'L',
				Left: nil, Right: nil,
			},
			Right: &TreeNode{
				Freq: 65, Rune: 0,
				Left: &TreeNode{
					Freq: 32, Rune: 'C',
					Left: nil, Right: nil,
				},
				Right: &TreeNode{
					Freq: 33, Rune: 0,
					Left: &TreeNode{
						Freq: 9, Rune: 0,
						Left: &TreeNode{
							Freq: 2, Rune: 'Z',
							Left: nil, Right: nil,
						},
						Right: &TreeNode{
							Freq: 7, Rune: 'K',
							Left: nil, Right: nil,
						},
					},
					Right: &TreeNode{
						Freq: 24, Rune: 'M',
						Left: nil, Right: nil,
					},
				},
			},
		},
	},
}
var referencePrefixMap = map[rune]byte{
	'C': 0b1110,
	'D': 0b101,
	'E': 0b0,
	'K': 0b111101,
	'L': 0b110,
	'M': 0b11111,
	'U': 0b100,
	'Z': 0b111100,
}

func TestMakeTree(t *testing.T) {
	data, err := os.ReadFile(TESTFILEPATH)
	if err != nil {
		t.Errorf("error reading file: %v\n", err)
		return
	}
	freq := countFrequency(data)
	root := makeTree(freq)
	t.Run("root node has correct frequency", func(t *testing.T) {
		freqSum := 0
		for _, v := range freq {
			freqSum += v
		}
		if root.Freq != freqSum {
			t.Errorf("expected sum of frequencies to be %d, got %d\n", freqSum, root.Freq)
		} else {
			fmt.Println("root node has correct frequency")
		}
	})
	t.Run("tree on example", func(t *testing.T) {

		tree := makeTree(freqExample)
		if tree.String() != referenceTree.String() {
			t.Errorf("expected tree to be %v, got %v\n", referenceTree, tree)
			fmt.Println("tree on example correct")
		} else {
			fmt.Println("tree on example correct")
		}
	})
}

func TestPrefixMap(t *testing.T) {
	t.Run("prefix map on example", func(t *testing.T) {

		prefixMap := make(map[rune]byte, len(freqExample))

		generatePrefixCodes(&prefixMap, referenceTree, 0)

		if !maps.Equal(prefixMap, referencePrefixMap) {
			t.Errorf("PrefixMap() returned %v, want %v", prefixMap, referencePrefixMap)
		} else {
			fmt.Println("PrefixMap() correct")
		}
	})
}

func TestFileHeader(t *testing.T) {
	t.Run("file header on example", func(t *testing.T) {
		prefixMap := make(map[rune]byte, len(freqExample))
		generatePrefixCodes(&prefixMap, referenceTree, 0)
		header := fileHeader(&prefixMap)
		fmt.Println(string(header))
		// TODO create and decode header
	})
}
