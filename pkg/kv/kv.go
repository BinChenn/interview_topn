package kv

import (
	"math/rand"
)

// KV pair
type KV struct {
	Key   int64
	Value []byte
}

// KVList list
type KVList []*KV

// Length of all kv
func (all_kv KVList) Len() int { return len(all_kv) }

// Lesser i.key is lesser than j.key
func (all_kv KVList) Less(i, j int) bool { return all_kv[i].Key < all_kv[j].Key }

// Swap i, j
func (all_kv KVList) Swap(i, j int) { all_kv[i], all_kv[j] = all_kv[j], all_kv[i] }

// Push interface
func (all_kv *KVList) Push(x interface{}) {
	*all_kv = append(*all_kv, x.(*KV))
}

// Pop interface
func (all_kv *KVList) Pop() interface{} {
	old := *all_kv
	n := len(old)
	x := old[n-1]
	*all_kv = old[:n-1]
	return x
}

// GenRandKV generate random KV
func GenRandKV(n int, minKey, maxKey int64) KVList {
	if maxKey-minKey < int64(n) {
		return nil
	}
	rand.Seed(1)
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	var randkvs KVList
	exist := make(map[int64]bool, 1000)
	count := 0

	for {
		newkv := KV{
			Key:   minKey + rand.Int63()%(maxKey-minKey),
			Value: make([]byte, rand.Int()%1024),
		}
		for j := range newkv.Value {
			newkv.Value[j] = letters[rand.Intn(len(letters))]
		}
		if _, ok := exist[newkv.Key]; ok {
			continue
		}
		exist[newkv.Key] = true
		randkvs = append(randkvs, &newkv)
		count++
		if count >= n {
			break
		}

	}
	return randkvs
}
