package main

import (
	"compressor/compressor"
	"fmt"
)

func main() {
	_, err := compressor.Compress()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println()
}
