package main

import (
	"github.com/BinChenn/interview_topn/pkg/multiple"
)

var (
	masterAddr string = "127.0.0.1:10000"
	// workerAddrs []string = []string{"127.0.0.1:10001", "127.0.0.1:10002", "127.0.0.1:10003"}
	workerAddrs []string = []string{"127.0.0.1:10001"}
	minkey      int64    = 1
	maxkey      int64    = 10000000
	topn                 = 5
)

// RunMultiple run
func RunMultiple() {
	multiple.InitUser(masterAddr)
	go func() {
		multiple.StartMaster(masterAddr, workerAddrs)
	}()

	go func() {
		multiple.StartWorker(workerAddrs[0])
	}()

	multiple.GetMulMachineTopN(minkey, maxkey, topn)
}

func main() {
	RunMultiple()
}
