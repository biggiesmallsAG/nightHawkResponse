
package common

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)


const Layout = "2006-01-02T15:04:05Z"
const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"


func NewSessionDir(szdir int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]byte, szdir)
	for i := range b {
		b[i] = chars[rand.Int63()%int64(len(chars))]
	}
	return string(b)
}

func GenerateCaseName() string {
	part_a := strings.ToUpper(NewSessionDir(5))
	part_b := strings.ToUpper(NewSessionDir(3))
	casename := fmt.Sprintf("%s-%s", part_a, part_b)
	return casename
}