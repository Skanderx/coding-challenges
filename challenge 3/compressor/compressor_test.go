package compressor

import (
	"fmt"
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
		freq := map[rune]int{
			'C': 32,
			'D': 42,
			'E': 120,
			'K': 7,
			'L': 42,
			'M': 24,
			'U': 37,
			'Z': 2,
		}
		tree := makeTree(freq)
		reference := &TreeNode{
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
		if tree.String() != reference.String() {
			t.Errorf("expected tree to be %v, got %v\n", reference, tree)
			fmt.Println("tree on example correct")
		} else {
			fmt.Println("tree on example correct")
		}
	})
}
