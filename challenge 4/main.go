package main

import (
	"flag"
	"fmt"
	"main/cut"
	"strconv"
	"strings"

	"os"
	"path/filepath"
)

var flags = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
var fields string
var delimiter string

func init() {
	flags.StringVar(&fields, "f", "1", "Fields to cut")
	flags.StringVar(&delimiter, "d", "\t", "Delimiter")
}
func splitFields(fields string) []string {
	if unquotedFields := strings.Trim(fields, "\""); unquotedFields != fields {
		return strings.Split(unquotedFields, " ")
	}
	return strings.Split(fields, ",")
}

func main() {
	flags.Parse(os.Args[1:])

	fileName := filepath.Clean(os.Args[len(os.Args)-1])

	dataIn, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(fields) == 0 {
		fmt.Println("Fields must be greater than 0")
		return
	}
	fieldsString := splitFields(fields)
	fieldIndexes := make([]int, len(fieldsString))
	for i, field := range fieldsString {
		fieldIndex, err := strconv.Atoi(field)
		if err != nil {
			fmt.Printf("field %s is not a number\n : %s", field, err)
			return
		}
		fieldIndexes[i] = fieldIndex - 1
	}

	fmt.Println(cut.Fields(dataIn, fieldIndexes, []rune(delimiter)[0]))
}
