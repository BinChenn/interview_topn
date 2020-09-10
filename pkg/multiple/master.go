package multiple

import (
	context "context"
	"fmt"
	"log"
	"net"
	"sort"
	"time"

	"github.com/BinChenn/interview_topn/pkg/kv"
	"github.com/BinChenn/interview_topn/pkg/single"
	grpc "google.golang.org/grpc"
)

// masterServer master server
type masterServer struct {
	masterAddr  string
	workerAddrs []string
}

// GetAllTopN master get all topn
func (ms *masterServer) GetAllTopN(ctx context.Context, req *TopNReq) (*TopNRsp, error) {
	minkey := req.MinKey
	maxkey := req.MaxKey
	topn := req.Topn

	// 使用二层hash构建key索引
	// get worker ip:port
	channel := make(chan kv.KVList, len(ms.workerAddrs))

	for i, server := range ms.workerAddrs {
		go func(i int, server string) {
			conn, err := grpc.Dial(server, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("fail to dial: %v", err)
			}
			defer conn.Close()
			client := NewGetWorkerTopNClient(conn)
			fmt.Println("master client inialized in: ", server)

			req := &TopNReq{
				MinKey: minkey,
				MaxKey: maxkey,
				Topn:   int64(topn),
			}

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			rsp, err := client.GetWorkerTopN(ctx, req)
			var topnList kv.KVList
			for _, item := range rsp.KvList {
				newkv := kv.KV{
					Key:   item.Key,
					Value: item.Value,
				}
				topnList = append(topnList, &newkv)
			}
			channel <- topnList
		}(i, server)
	}

	var mergeList kv.KVList
	// merge topnlist from workers
	for i := 0; i < len(ms.workerAddrs); i++ {
		if topnList, ok := <-channel; ok {
			mergeList = append(mergeList, topnList...)
		} else {
			fmt.Println("master channel error")
		}
	}
	topnList := single.GetSingleTopNbyRange(mergeList, int(topn), minkey, maxkey)
	sort.Sort(topnList)
	var result []*KV
	for _, item := range topnList {
		item := KV{
			Key:   item.Key,
			Value: item.Value,
		}
		result = append(result, &item)
	}
	return &TopNRsp{
		KvList: result,
	}, nil

}

// StartMaster start master instance only once
func StartMaster(masterAddr string, workerAddrs []string) {
	lis, err := net.Listen("tcp", masterAddr)
	if err != nil {
		log.Fatalf("Master failed to listen: %s", masterAddr)
	}
	grpcServer := grpc.NewServer()
	ms := masterServer{
		masterAddr:  masterAddr,
		workerAddrs: workerAddrs,
	}
	RegisterGetTopNServer(grpcServer, &ms)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Master failed to serve: %v", err)
	}
	fmt.Println("start master success")
}
