package multiple

import (
	context "context"
	"fmt"
	"log"
	"time"

	"github.com/BinChenn/interview_topn/pkg/kv"
	"google.golang.org/grpc"
)

var (
	serverAddr string
)

// GetMulMachineTopN user interface
func GetMulMachineTopN(minkey, maxkey int64, topn int) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := NewGetTopNClient(conn)
	fmt.Println("user client inialized in: ", serverAddr)

	req := &TopNReq{
		MinKey: minkey,
		MaxKey: maxkey,
		Topn:   int64(topn),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	rsp, err := client.GetAllTopN(ctx, req)

	var topnList kv.KVList
	for _, item := range rsp.KvList {
		newkv := kv.KV{
			Key:   item.Key,
			Value: item.Value,
		}
		topnList = append(topnList, &newkv)
	}
	fmt.Println("Get top ", len(topnList), " kv")
	kv.PrintKeys(topnList)
}

// InitUser init service addr
func InitUser(masterAddr string) {
	serverAddr = masterAddr
}
