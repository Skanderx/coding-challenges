package cut

import "strings"

func Fields(dataIn []byte, index []int, delimiter rune) string {
	out := ""
	for _, line := range strings.Split(string(dataIn), "\n") {
		columns := strings.Split(line, string(delimiter))
		for _, index := range index {
			if index >= len(columns) {
				continue
			}
			out += columns[index] + string(delimiter)
		}
		out += "\n"
	}
	return out
}
