package multiple

import (
	context "context"
	"fmt"
	"log"
	"net"
	"sort"

	"github.com/BinChenn/interview_topn/pkg/kv"
	"github.com/BinChenn/interview_topn/pkg/single"
	grpc "google.golang.org/grpc"
)

type workerServer struct {
	workerAddr string
}

func (ms *workerServer) GetWorkerTopN(ctx context.Context, req *TopNReq) (*TopNRsp, error) {
	minkey := req.MinKey
	maxkey := req.MaxKey
	topn := req.Topn

	chunks := kv.GetChunks()
	var mergeList kv.KVList

	for _, chunk := range chunks {
		kvList := kv.ReadByChunk(chunk)
		topnList := single.GetMultiCoreTopNbyRange(kvList, int(topn), minkey, maxkey, single.GetSingleTopN, single.DataSplitbySize)
		mergeList = append(mergeList, topnList...)
	}
	topnList := single.GetSingleTopNbyRange(mergeList, int(topn), minkey, maxkey)
	sort.Sort(topnList)
	result := make([]*KV, len(topnList))
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

// StartWorker start worker instance only once
func StartWorker(workerAddr string) {
	lis, err := net.Listen("tcp", workerAddr)
	if err != nil {
		log.Fatalf("Worker failed to listen: %s", workerAddr)
	}
	grpcServer := grpc.NewServer()
	ws := workerServer{
		workerAddr: workerAddr,
	}
	RegisterGetWorkerTopNServer(grpcServer, &ws)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Worker failed to serve: %v", err)
	}
	fmt.Println("start worker success")

}
