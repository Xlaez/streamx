package utils

import (
	"math/rand"
	"strings"
	"time"
)

const num = "1234567890"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomIntegers(n int) string {
	var sb strings.Builder

	k := len(num)

	for i := 0; i < n; i++ {
		c := num[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}
