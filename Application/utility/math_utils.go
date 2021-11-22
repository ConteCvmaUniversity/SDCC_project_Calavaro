package utility

import (
	"math/rand"
	"time"
)

func MaxOf(vars ...uint64) uint64 {
	max := vars[0]

	for _, i := range vars {
		if max < i {
			max = i
		}
	}
	return max
}

func GetRandInt(max int) int {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(max)
	return n
}
