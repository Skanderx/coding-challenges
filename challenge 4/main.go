package main

import (
	"flag"
	"fmt"
	"main/cut"

	"os"
	"path/filepath"
)

var flags = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
var fields int // []int
var delimiter string

func init() {
	flags.IntVar(&fields, "f", 1, "Fields to cut")
	flags.StringVar(&delimiter, "d", "\t", "Delimiter")
}

func main() {
	flags.Parse(os.Args[1:])

	fileName := filepath.Clean(os.Args[len(os.Args)-1])

	dataIn, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	if fields < 1 {
		fmt.Println("Fields must be greater than 0")
		return
	}

	fmt.Println(cut.Fields(dataIn, fields-1, []rune(delimiter)[0]))
}
