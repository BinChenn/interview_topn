package single

import (
	"container/heap"
	"fmt"
	"runtime"
	"sort"

	"github.com/BinChenn/interview_topn/pkg/kv"
)

// TopNFunc interface
type TopNFunc func(kvlist kv.KVList, topn int) kv.KVList

// SplitFunc interface
type SplitFunc func(kvlist kv.KVList, routinePerCPU int) []kv.KVList

// GetBaseLineTopN get real topn key value
func GetBaseLineTopN(kvlist kv.KVList, topn int) kv.KVList {
	if len(kvlist) < topn {
		return kvlist
	}
	sort.Sort(kvlist)
	return kvlist[len(kvlist)-topn:]
}

// GetSingleTopN get topn with
func GetSingleTopN(kvlist kv.KVList, topn int) kv.KVList {
	if len(kvlist) < topn {
		return kvlist
	}
	newkvlist := kvlist[:topn]
	heap.Init(&newkvlist)

	for _, item := range kvlist[topn:] {
		if item.Key > newkvlist[0].Key {
			item := item
			newkvlist[0] = item
			heap.Fix(&newkvlist, 0)
		}
	}
	return newkvlist
}

// GetSingleTopNbyRange get topn by range
func GetSingleTopNbyRange(kvlist kv.KVList, topn int, minkey, maxkey int64) kv.KVList {
	var newlist kv.KVList

	for _, item := range kvlist {
		if item.Key > minkey && item.Key < maxkey {
			newlist = append(newlist, item)
		}
	}
	return GetSingleTopN(newlist, topn)
}

// DataSplitbySize split data by data size to multi go routine
func DataSplitbySize(kvlist kv.KVList, routinePerCPU int) []kv.KVList {
	// setup max cpus can use
	runtime.GOMAXPROCS(runtime.NumCPU())
	sliceNum := runtime.NumCPU() * routinePerCPU
	fmt.Println("cpu num: ", runtime.NumCPU(), " slice: ", sliceNum)
	sliceLen := len(kvlist) / sliceNum

	segments := make([]kv.KVList, sliceNum)
	for i := 0; i <= sliceNum-1; i++ {
		if i < sliceNum {
			segments[i] = kvlist[sliceLen*i : sliceLen*(i+1)]
		} else {
			segments[i] = kvlist[sliceLen*i:]
		}
	}
	return segments
}

// DataSplitbyHash split data by hash to multi go routine
func DataSplitbyHash(kvlist kv.KVList, routinePerCPU int) []kv.KVList {
	sliceNum := runtime.NumCPU() * routinePerCPU
	fmt.Println("cpu num: ", runtime.NumCPU(), " slice: ", sliceNum)
	segments := make([]kv.KVList, sliceNum)
	// hash func: data%sliceNum
	for key, value := range kvlist {
		segments[key%sliceNum] = append(segments[key%sliceNum], value)
	}
	return segments
}

// GetMultiCoreTopN get topn goroutine version
func GetMultiCoreTopN(kvlist kv.KVList, topn int, getTopN TopNFunc, split SplitFunc) kv.KVList {
	if len(kvlist) < topn {
		return kvlist
	}
	segments := split(kvlist, 1)

	channel := make(chan kv.KVList, len(segments))

	for _, kvlist := range segments {
		kvlist := kvlist
		go func() {
			channel <- getTopN(kvlist, topn)
		}()
	}
	var mergeList kv.KVList

	for i := 0; i < len(segments); i++ {
		if topnList, ok := <-channel; ok {
			mergeList = append(mergeList, topnList...)
		} else {
			fmt.Println("multi core channel error")
		}
	}

	result := getTopN(mergeList, topn)
	return result
}

// GetMultiCoreTopNbyRange get topn goroutine version
func GetMultiCoreTopNbyRange(kvlist kv.KVList, topn int, minkey, maxkey int64, getTopN TopNFunc, split SplitFunc) kv.KVList {
	var newlist kv.KVList

	for _, item := range kvlist {
		if item.Key > minkey && item.Key < maxkey {
			newlist = append(newlist, item)
		}
	}
	return GetMultiCoreTopN(newlist, topn, getTopN, split)

}
