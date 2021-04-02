package tests

import (
	"math/rand"
	"strings"
)

func makeID(prefix string, count int) string {
	var (
		chars = "012345678abcdefghijkmnpqrstuvwxy"
		u     = rand.Int63()
		buf   = new(strings.Builder)
	)

	buf.WriteString(prefix)
	buf.WriteByte('-')
	for i, j := 0, 0; i < count; i++ {
		if j+5 > 64 {
			j, u = 0, rand.Int63()
		}
		buf.WriteByte(chars[u&0x1F])
		j, u = j+5, u>>5
	}
	return buf.String()
}
