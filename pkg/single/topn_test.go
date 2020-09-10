package single

import (
	"fmt"
	"testing"

	"github.com/BinChenn/interview_topn/pkg/kv"
)

func TestGetBaseLineTopN(t *testing.T) {
	fmt.Println("=====TestGetBaseLineTopN=====")
	l := kv.GenRandKV(100, 1, 1000, false)
	// fmt.Print("Original Keys: ")
	// kv.PrintKeys(l)

	topn := 5
	topnList := GetBaseLineTopN(l, topn)
	fmt.Print("Top ", topn, ": ")
	kv.PrintKeys(topnList)
}

func TestGetSingleTopN(t *testing.T) {
	fmt.Println("=====TestGetSingleTopN=====")

	l := kv.GenRandKV(100, 1, 1000, false)
	// fmt.Print("Original Keys: ")
	// kv.PrintKeys(l)

	topn := 5
	topnList := GetSingleTopN(l, topn)
	fmt.Print("Top ", topn, ": ")
	kv.PrintKeys(topnList)
}

func TestGetMultiCoreTopN(t *testing.T) {
	fmt.Println("=====TestGetMultiCoreTopN=====")
	l := kv.GenRandKV(100, 1, 1000, false)
	topn := 5

	topnList := GetMultiCoreTopN(l, topn, GetSingleTopN, DataSplitbyHash)
	// topnList := GetMultiCoreTopN(l, topn, GetSingleTopN, DataSplitbySize)
	fmt.Print("Top ", topn, ": ")
	kv.PrintKeys(topnList)
}
