package cut

import "strings"

func Fields(dataIn []byte, index int) string {
	out := ""
	for _, line := range strings.Split(string(dataIn), "\n") {
		out += strings.Split(line, "\t")[index] + "\n"
	}
	return out
}
