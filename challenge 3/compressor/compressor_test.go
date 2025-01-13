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
