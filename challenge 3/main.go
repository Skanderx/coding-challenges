package main

import (
	compressor "compressor/file-compressor"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Read file
	args := os.Args
	fileName := filepath.Clean(args[len(args)-1])

	f := flag.NewFlagSet("wordcount", flag.ContinueOnError)
	decompress := f.Bool("d", false, "decompress file")
	err := f.Parse(args[1:])
	if err != nil {
		fmt.Println(err)
		return
	}
	if decompress != nil && *decompress {
		err := compressor.Decompress(fileName)
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}

	err = compressor.Compress(fileName)
	if err != nil {
		fmt.Println(err)
	}

}
