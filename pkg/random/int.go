package random

import (
	"crypto/rand"
	"math/big"
)

// Int 产生随机数
// [min, max] random number
func Int(min, max int64) int64 {
	r, _ := rand.Int(rand.Reader, big.NewInt(max-min))
	return r.Int64() + min
}
