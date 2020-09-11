# Topn Problem
## 1. 题目
`
题目;
内存中的行列结构的数据集，存在主键 k，求 TopN 算法
上述题目在多核环境下的优化
数据集大小为 1TB，分布规律未知。存储在某存储服务上，以 get(min_k, max_k) 接口获取数据，求多台服务器的计算方案
`
## 2. 问题分解
上述问题可分解为以下三个部分：
- 求内存中主键key最大的n条数据（topn）
- 利用多核加速求topn算法
- 分布式环境下，多台机器协作求全局topn

## 3. 方案设计与实现
### 3.1 数据结构设计
存在主键key的行列结构数据可抽象为(Key, Value)数据对，本项目中对(Key, Value)数据结构定义如下：
```go
type KV struct {
	Key   int64
	Value []byte
}
```

其中，Key表示数据主键，Value表述数据值，用字符数组表示。在实际系统中value值可以为序列化的json或protoc结构数据，可以统一抽象为字符数组。

### 3.2 内存topn问题
基于上述数据结构定义，使用小顶堆算法求(Key, Value)数组中的topn，具体算法原理为：首先使用(Key, Value)数组中前n个记录的key建立小顶堆，然后依次将第n+1个到最后一个key与堆顶比较。如果当前key大于堆顶，则替换堆顶并调整堆。如果小于堆顶则跳过。最终，堆中的n个(Key, Value)即为topn。

为了验证小顶堆topn算法的正确性，本项目实现了基于直接排序的算法作为baseline。即直接对所有(Key, Value)根据key排序，并取其中key最大的n个(Key, Value)作为topn。

本部分代码实现在`interview_topn/pkg/single`。

**复杂度分析**

假设内存中有m条(Key, Value)数据，欲求top n。
- 直接排序算法
排序算法时间复杂度为O(m logm)
- 小顶堆算法
建堆时间复杂度为O(n),之后m-n个key，每次调整的时间复杂度为O(log n)， 因此算法整体时间复杂度为：O((m-n)log n + n), 若n与m无关，则可视为常数，简化后的时间复杂度为：O(m log n)。

### 3.3 多核topn加速
现代cpu普遍采用多核架构，因此可以利用多核CPU并行加速topn计算过程。
本项目使用不同的切分方法，将(Key, Value)数据切分到不同cpu核上，核之间并行计算各自数据的topn，最终汇总计算得到全局数据的topn。

本部分代码实现在`interview_topn/pkg/single`。

**数据切分**

本项目实现了两种数据切分方法，分别基于数据块大小(size based)和Key值哈希(Hash based)。

**并行计算**

本项目利用Go语言中的goroutine机制，在完成数据切分以后，并行启动多个协程，分别计算各自数据块的topn，最后汇总。

需要注意的是，go语言中启动的多个goroutine默认单核运行，不能充分利用多核优势。因此需要设置允许go使用的最大CPU核数来充分利用多核。方法如下：
```go
runtime.GOMAXPROCS(runtime.NumCPU())
```

**并行计算正确性**

可以通过反证法证明：全局topn集合，必然是切分后多个局部数据块求得的topn集合交集的子集。即若有一个key出现在全局topn中，则它必然存在于某个局部数据块的topn集合中。

假设某个主键k存在于全局topn集合，但不存在于任意一个局部数据块的topn集合。假设k被切分到chunk A中，如果k不在chunk A的topn集合中，那就意味着在A中有n个比k大的主键，将数据全局降序排序，主键k必然在这n个比k大的主键的后面，因此主键k不在全局topn集合中，与假设矛盾，结论得证。


### 3.4 分布式方案

包括存储，计算两部分，本部分代码实现在`interview_topn/pkg/multiple`。

**存储**

借鉴GFS的设计方案，将(Key, Value)数据保存为64MB大小的chunk，chunk内key分布随机，chunk分布在不同服务器上，暂未实现chunk replicas机制。

**计算**

使用master，worker结构处理topn请求。具体请求处理流程如下：
1. user向master节点发出topn请求
2. 全局唯一的master节点接收用户请求，并调用各个worker的topn接口，等待结果返回。
3. 分布在不同服务器上的worker节点收到从master发出的请求，从本地存储chunk中读出数据并调用并行topn算法计算topn，最终将结果返回master。
4. master汇总worker返回的结果，最终将全局topn返回给用户。

**Locality**

Locality在分布式系统中尤为重要，即计算节点应该尽量利用本地数据计算，而非使用远程数据，因为将大量数据从远程节点传输至本地需要耗费大量网络带宽，会严重拖慢计算过程。因此，本项目中worker基于本地数据计算topn，避免了网络传输过程中的瓶颈。

## 4. 测试结果

单机topn算法测试结果如下：

```bash
(base) chenbin@netlab-Z390-GAMING-X:~/go/interview_topn$ cd ./pkg/single/ && go test
=====TestGetBaseLineTopN=====
Top 5: 958 959 968 973 975  
=====TestGetSingleTopN=====
Top 5: 958 959 975 968 973  
=====TestGetMultiCoreTopN=====
cpu num:  8  slice:  8
Top 5: 958 959 975 973 968  
PASS
ok      github.com/BinChenn/interview_topn/pkg/single   0.004s
```

分布式topn算法使用3台linux服务器作为worker，服务器处于同一局域网内，分别在三台服务器上以不同的范围生成本地数据，使用linux台式机作为master节点处理请求。测试中可以正确返回topn结果。为便于单机运行，代码中保留了在一台机器上同时启动master节点和worker节点。测试结果如下：
```bash
(base) chenbin@netlab-Z390-GAMING-X:~/go/interview_topn$ go run main.go 
user client inialized in:  127.0.0.1:10000
master client inialized in:  127.0.0.1:10001
Get top  5  kv
2645411 6276998 7191835 9486845 9531697 
```

## 5. 后续优化
1. chunk replicas

本项目考虑locality，将worker的计算范围限制在本地数据。优点是避免了通过网络传输数据的通信瓶颈，缺点则在于一旦某个worker故障，则无法返回正确结果。因此后续可以考虑实现chunk replicas机制，避免单个worker故障导致服务不可用。

2. 数据切分

基于固定size和hash的方法实现简单，但一旦有新节点接入就要重新划分数据，因此可以考虑使用一致性哈希的方式在worker间切分数据。

3. Key索引

目前的topn算法需要遍历读取本地所有chunk求得topn，而题目中给定了求取topn的范围值(min_k, max_k),因此可以考虑对key-chunk index的对应关系建立索引，避免读取key值范围不在(min_k, max_k)内的chunk，减少文件io时间。







