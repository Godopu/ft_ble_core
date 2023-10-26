package alg

import (
	"math/rand"
	"time"
)

func Detect(list []int) uint64 {
	slTime := uint64(320) + uint64(rand.Intn(80))
	time.Sleep(time.Millisecond * time.Duration(slTime))

	return slTime
}
