package ssp

import (
	"math/rand"
	"time"
)

var pool = "abcdefghijklmnopqrstuvwxqz0123456789"

func init() {
	rand.Seed(time.Now().Unix())
}

func RandomID(n int) string {
	res := make([]byte, n)
	for i := range res {
		res[i] = pool[rand.Intn(len(pool))]
	}
	return string(res)
}
