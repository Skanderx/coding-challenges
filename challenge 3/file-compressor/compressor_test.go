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
		t.Skipf("error reading file: %v\n", err)
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

func TestEncoderDecoder(t *testing.T) {
	tests := []struct {
		freq               map[byte]int
		referenceTree      *TreeNode
		referencePrefixMap map[byte]byte
	}{{
		freq: map[byte]int{
			'C': 32,
			'D': 42,
			'E': 120,
			'K': 7,
			'L': 42,
			'M': 24,
			'U': 37,
			'Z': 2,
		},
		referenceTree: &TreeNode{
			Freq: 306, Char: 0,
			Left: &TreeNode{
				Freq: 120, Char: 'E',
				Left: nil, Right: nil,
			},
			Right: &TreeNode{
				Freq: 186, Char: 0,
				Left: &TreeNode{
					Freq: 79, Char: 0,
					Left: &TreeNode{
						Freq: 37, Char: 'U',
						Left: nil, Right: nil,
					},
					Right: &TreeNode{
						Freq: 42, Char: 'D',
						Left: nil, Right: nil,
					},
				},
				Right: &TreeNode{
					Freq: 107, Char: 0,
					Left: &TreeNode{
						Freq: 42, Char: 'L',
						Left: nil, Right: nil,
					},
					Right: &TreeNode{
						Freq: 65, Char: 0,
						Left: &TreeNode{
							Freq: 32, Char: 'C',
							Left: nil, Right: nil,
						},
						Right: &TreeNode{
							Freq: 33, Char: 0,
							Left: &TreeNode{
								Freq: 9, Char: 0,
								Left: &TreeNode{
									Freq: 2, Char: 'Z',
									Left: nil, Right: nil,
								},
								Right: &TreeNode{
									Freq: 7, Char: 'K',
									Left: nil, Right: nil,
								},
							},
							Right: &TreeNode{
								Freq: 24, Char: 'M',
								Left: nil, Right: nil,
							},
						},
					},
				},
			},
		},
		referencePrefixMap: map[byte]byte{
			'C': 0b1110,
			'D': 0b101,
			'E': 0b0, // Starting with 0 is a problem
			'K': 0b111101,
			'L': 0b110,
			'M': 0b11111,
			'U': 0b100,
			'Z': 0b111100,
		},
	}}
	for _, tt := range tests {
		t.Run("Make tree", func(t *testing.T) {
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

				tree := makeTree(tt.freq)
				if tree.String() != tt.referenceTree.String() {
					t.Errorf("expected tree to be %v, got %v\n", tt.referenceTree, tree)
					fmt.Println("tree on example correct")
				} else {
					fmt.Println("tree on example correct")
				}
			})
		})

		t.Run("Make prefix map", func(t *testing.T) {
			t.Run("prefix map on example", func(t *testing.T) {

				prefixMap := make(map[byte]byte, len(tt.freq))

				generatePrefixCodes(prefixMap, tt.referenceTree, 0)

				if !maps.Equal(prefixMap, tt.referencePrefixMap) {
					t.Errorf("PrefixMap() returned %v, want %v", prefixMap, tt.referencePrefixMap)
				} else {
					fmt.Println("PrefixMap() correct")
				}
			})
		})

		t.Run("File Header", func(t *testing.T) {
			t.Run("file header on example", func(t *testing.T) {
				prefixMap := make(map[byte]byte, len(tt.freq))
				generatePrefixCodes(prefixMap, tt.referenceTree, 0)
				header := fileHeader(prefixMap)

				_, prefixMapDecoded, err := decodeHeader(header)
				if err != nil {
					t.Errorf("error decoding header: %v\n", err)
					return
				}
				if !maps.Equal(prefixMap, prefixMapDecoded) {
					t.Errorf("PrefixMap() returned %v, want %v", prefixMap, prefixMapDecoded)
				} else {
					fmt.Println("PrefixMap() correct")
				}
			})
		})
	}
}
