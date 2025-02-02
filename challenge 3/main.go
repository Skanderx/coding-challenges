package main

import (
	compressor "compressor/file-compressor"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var flags = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
var decompress bool

func init() {
	flags.BoolVar(&decompress, "d", false, "decompress file")
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments")
		return
	}

	flags.Parse(os.Args[1:])

	args := os.Args[len(os.Args)-2:]
	fileNameIn := filepath.Clean(args[0])
	fileNameOut := filepath.Clean(args[1])

	dataIn, err := os.ReadFile(fileNameIn)
	if err != nil {
		fmt.Println(err)
		return
	}

	fileOut, err := os.Create(fileNameOut)
	if err != nil {
		fmt.Println(fmt.Errorf("error creating file: %w", err))
	}
	defer fileOut.Close()

	if decompress {
		err = compressor.Decompress(dataIn, fileOut)
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	} else {
		err = compressor.Compress(dataIn, fileOut)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}
