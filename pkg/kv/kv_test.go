package kv

import (
	"fmt"
	"testing"
)

func TestGenRandKV(t *testing.T) {
	l := GenRandKV(10, 1, 100)
	for _, kv := range l {
		fmt.Println(kv.Key, " ", len(kv.Value))
	}
}
