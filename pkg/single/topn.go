package single

import (
	"container/heap"
	"sort"

	"github.com/BinChenn/interview_topn/pkg/kv"
)

// GetRealTopN : get real topn key value
func GetBaseLineTopN(kvlist kv.KVList, topn int) kv.KVList {
	if len(kvlist) < topn {
		return kvlist
	}
	sort.Sort(kvlist)
	return kvlist[:topn]
}

// GetSingleTopN get topn with
func GetSingleTopN(kvlist kv.KVList, topn int) kv.KVList {
	if len(kvlist) < topn {
		return kvlist
	}
	newkvlist := kvlist[:topn]
	heap.Init(&newkvlist)

	for _, item := range kvlist[topn:] {
		if item.Key < newkvlist[0].Key {
			newkvlist[0] = item
			heap.Fix(&newkvlist, 0)
		}
	}
	return newkvlist
}
