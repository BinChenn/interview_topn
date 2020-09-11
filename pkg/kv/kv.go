package kv

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
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

// PrintKeys print all keys
func PrintKeys(l KVList) {
	for _, kv := range l {
		fmt.Print(kv.Key, " ")
	}
	fmt.Println(" ")
}

// GenRandKV generate random KV
func GenRandKV(n int, minKey, maxKey int64, writeFlag bool) KVList {
	if maxKey-minKey < int64(n) {
		return nil
	}
	rand.Seed(1)
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	var randkvs KVList
	exist := make(map[int64]bool, 1000)
	count := 0

	currDataSize := 0
	chunkSize := 67108864 // chunk size: 64MB
	chunkIndex := 0

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
		// kv size  = 8(size of key) + size of value
		currDataSize += (8 + len(newkv.Value))
		if writeFlag && currDataSize > chunkSize {
			WriteChunk(randkvs, int64(chunkIndex))
			chunkIndex++
			currDataSize = 0
			randkvs = KVList{}
		}
		count++
		if count >= n {
			break
		}
	}
	return randkvs
}

// WriteChunk write generate kv to a chunk file, chunkSize(MB)
func WriteChunk(kvlist KVList, chunkIndex int64) {
	dataFName := fmt.Sprintf(".data/chunk-%d.data", chunkIndex)
	// record the size of each value in kv
	lenFname := fmt.Sprintf(".data/len-%d.data", chunkIndex)

	dataFOut, err := os.Create(dataFName)
	if err != nil {
		panic(err)
	}

	lenFOut, err := os.Create(lenFname)
	if err != nil {
		panic(err)
	}

	for _, item := range kvlist {
		keyByted := make([]byte, 8)
		binary.LittleEndian.PutUint64(keyByted, uint64(item.Key))
		if _, err := dataFOut.Write(keyByted); err != nil {
			panic(err)
		}
		if _, err := dataFOut.Write(item.Value); err != nil {
			panic(err)
		}
		lenRecord := fmt.Sprintf("%d\n", len(item.Value))
		if _, err := lenFOut.WriteString(lenRecord); err != nil {
			panic(err)
		}
	}
	if err := dataFOut.Close(); err != nil {
		panic(err)
	}
	if err := lenFOut.Close(); err != nil {
		panic(err)
	}

}

// ReadByChunk read kv list by chunk index
func ReadByChunk(chunkIndex int64) KVList {
	dataFName := fmt.Sprintf("./pkg/kv/data/chunk-%d.data", chunkIndex)
	// record the size of each value in kv
	lenFname := fmt.Sprintf("./pkg/kv/data/len-%d.data", chunkIndex)
	dataByted, err := ioutil.ReadFile(dataFName)
	if err != nil {
		panic(err)
	}
	lenByted, err := ioutil.ReadFile(lenFname)
	if err != nil {
		panic(err)
	}

	// get length of each kv(offset)
	lenStr := string(lenByted)
	lenStrs := strings.Split(lenStr, "\n")
	lenStrs = lenStrs[:len(lenStrs)-1]
	lenRecords := make([]int, len(lenStrs))
	for i, str := range lenStrs {
		lenRecords[i], err = strconv.Atoi(str)
		if err != nil {
			panic(err)
		}
	}
	// fmt.Println("Read file: ", dataFName, " len: ", len(lenRecords))
	kvlist := make([]*KV, len(lenRecords))
	offset := 0
	for i, valueLen := range lenRecords {
		keyByted := dataByted[offset : offset+8]
		newkv := KV{
			Key:   int64(binary.LittleEndian.Uint64(keyByted)),
			Value: dataByted[offset+8 : offset+8+valueLen],
		}
		kvlist[i] = &newkv
		offset += (8 + valueLen)
	}

	return kvlist
}

// GetWorkers return chunk indexs within range[minkey, maxkey]
func GetWorkers() []string {
	// TODO: 索引构建
	indexs := make([]string, 10)
	return indexs
}

// GetChunks return chunks id store in worker i
func GetChunks() []int64 {
	files, err := ioutil.ReadDir("./pkg/kv/data")

	if err != nil {
		log.Fatal(err)
	}

	chunkNum := len(files) / 2
	chunks := make([]int64, chunkNum)
	for i := 0; i < chunkNum; i++ {
		chunks[i] = int64(i)
	}
	return chunks
}
