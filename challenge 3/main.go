package main

import (
	compressor "compressor/file-compressor"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Read file
	args := os.Args
	fileName := filepath.Clean(args[len(args)-1])
	_, err := compressor.Compress(fileName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println()
}
