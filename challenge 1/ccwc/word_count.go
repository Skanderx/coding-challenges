package ccwc

import (
	"bufio"
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

func WordCount() string {

	args := os.Args

	f := flag.NewFlagSet("wordcount", flag.ContinueOnError)
	nbrOfBytes := f.Bool("c", false, "number of bytes")
	nbrOfLines := f.Bool("l", false, "number of lines")
	nbrOfWords := f.Bool("w", false, "number of words")
	nbrOfMultiBytes := f.Bool("m", false, "number of bytes excluding multi byte characters")

	err := f.Parse(args[1:])
	if err != nil {
		return err.Error()
	}
	// default to all=
	if !(*nbrOfBytes || *nbrOfLines || *nbrOfWords || *nbrOfMultiBytes) {
		*nbrOfBytes = true
		*nbrOfLines = true
		*nbrOfWords = true
	}

	dat := make([]byte, 0)
	fileName := ""
	if (len(args) > 1 && args[len(args)-1][0] == '-') || len(args) == 1 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			dat = append(dat, scanner.Bytes()...)
			dat = append(dat, '\n')
		}
		dat = dat[:len(dat)-1]
	} else {
		fileName = filepath.Clean(args[len(args)-1])
		var err error
		dat, err = os.ReadFile(fileName)
		if err != nil {
			panic(err)
		}
	}

	result := ""
	// number of lines
	if *nbrOfLines {
		result += "\t" + strconv.Itoa(strings.Count(string(dat), "\n"))
	}
	// number of words
	if *nbrOfWords {
		result += "\t" + strconv.Itoa(len(strings.Fields(string(dat))))
	}
	// number of carachterBytes
	if *nbrOfBytes {
		result += "\t" + strconv.Itoa(len(dat))
	}
	// number of bytes
	if *nbrOfMultiBytes {
		result += "\t" + strconv.Itoa(utf8.RuneCountInString(string(dat)))
	}

	if result == "" {
		return "unknown error"
	}

	return result + " " + fileName
}
