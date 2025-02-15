package cut

import "strings"

func Fields(dataIn []byte, index int, delimiter rune) string {
	out := ""
	for _, line := range strings.Split(string(dataIn), "\n") {
		out += strings.Split(line, string(delimiter))[index] + "\n"
	}
	return out
}
