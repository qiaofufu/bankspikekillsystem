package Random

import (
	"math/rand"
	"time"
)

type Random struct{}

const base string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (receiver Random) RandomString(n int) (str []byte) {
	rand.Seed(time.Now().Unix())
	for i := 0; i < n; i++ {
		idx := rand.Intn(26 * 2)
		str = append(str, base[idx])
	}
	return
}
