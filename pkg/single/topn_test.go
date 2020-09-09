package single

import (
	"fmt"
	"testing"

	"github.com/BinChenn/interview_topn/pkg/kv"
)

func printKeys(l kv.KVList) {
	for _, kv := range l {
		fmt.Print(kv.Key, " ")
	}
	fmt.Println(" ")
}

func TestGetBaseLineTopN(t *testing.T) {
	fmt.Println("=====TestGetBaseLineTopN=====")
	l := kv.GenRandKV(100, 1, 1000)
	// fmt.Print("Original Keys: ")
	// printKeys(l)

	topn := 5
	topnList := GetBaseLineTopN(l, topn)
	fmt.Print("Top ", topn, ": ")
	printKeys(topnList)
}

func TestGetSingleTopN(t *testing.T) {
	fmt.Println("=====TestGetSingleTopN=====")

	l := kv.GenRandKV(100, 1, 1000)
	// fmt.Print("Original Keys: ")
	// printKeys(l)

	topn := 5
	topnList := GetSingleTopN(l, topn)
	fmt.Print("Top ", topn, ": ")
	printKeys(topnList)
}

func TestGetMultiCoreTopN(t *testing.T) {
	fmt.Println("=====TestGetMultiCoreTopN=====")
	l := kv.GenRandKV(100, 1, 1000)
	topn := 5

	topnList := GetMultiCoreTopN(l, topn, GetSingleTopN, DataSplitbyHash)
	// topnList := GetMultiCoreTopN(l, topn, GetSingleTopN, DataSplitbySize)
	fmt.Print("Top ", topn, ": ")
	printKeys(topnList)

}
