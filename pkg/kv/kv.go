package kv

// KV pair
type KV struct {
	Key   int64
	Value string
}

// KVList list
type KVList []*KV

// Length of all kv
func (all_kv KVList) Len() int { return len(all_kv) }

// Lesser j.key is lesser than i.key
func (all_kv KVList) Less(i, j int) bool { return all_kv[i].Key > all_kv[j].Key }

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
