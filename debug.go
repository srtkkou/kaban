package kaban

import (
	"fmt"
	"strings"
)

func xdump(blob []byte) string {
	b := new(strings.Builder)
	fmt.Fprint(b,
		"    |00----02----04----06----08----0A----0C----0E---")
	for i, v := range blob {
		if i%16 == 0 {
			fmt.Fprintln(b)
			fmt.Fprintf(b, "%04X|", (i / 16))
		}
		fmt.Fprintf(b, "%02X ", v)
	}
	if len(blob)%16 != 0 {
		fmt.Fprintln(b)
	}
	return b.String()
}
