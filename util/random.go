package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "qwertyuiopasdfghjklmnbvcxz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// random Int generate a random integer between min and max
func RandomInt(min, max int64) int64 {
	if min > max {
		min, max = max, min // đảm bảo min <= max
	}
	return min + rand.Int63n(max-min+1) //ham and.Int63n(n int64) Trả về một số ngẫu nhiên trong khoảng [0, n)
}

// random string of length
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}
func RandomOwner() string {
	return RandomString(6)
}
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	return currencies[rand.Intn(len(currencies))]
}
