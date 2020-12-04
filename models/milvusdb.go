package models

import (
	"errors"
	"github.com/githubchry/gomdb/drivers"
	"github.com/milvus-io/milvus-sdk-go/milvus"
	"log"
)

//获取配置信息
func GetConfig() (string, error) {
	//GetConfig
	configInfo, status, _ := drivers.MilvusDbConn.GetConfig("*")
	if !status.Ok() {
		log.Println("Get config failed: " + status.GetMessage())
		return configInfo, errors.New(status.GetMessage())
	}
	return configInfo, nil
}

/*
创建集合
collectionName 集合名称
dimension 维度
index_file_size 自动创建索引的数据文件大小, 单位为M
metricType 距离度量方式

一个集合collection可以有多个分区partition, 一个分区可以有多个数据段segment
*/
func CreateCollection(collectionName string, dimension, indexFileSize, metricType int64) {
	//1.准备创建集合所需参数：
	collectionParam := milvus.CollectionParam{collectionName, dimension, indexFileSize, metricType}

	var hasCollection bool
	hasCollection, status, err := drivers.MilvusDbConn.HasCollection(collectionName)
	if err != nil {
		log.Println("HasCollection rpc failed: " + err.Error())
	}
	if hasCollection == false {
		status, err = drivers.MilvusDbConn.CreateCollection(collectionParam)
		if err != nil {
			log.Println("CreateCollection rpc failed: " + err.Error())
			return
		}
		if !status.Ok() {
			log.Println("Create collection failed: " + status.GetMessage())
			return
		}
		log.Println("Create collection " + collectionName + " success")
	}

	hasCollection, status, err = drivers.MilvusDbConn.HasCollection(collectionName)
	if err != nil {
		log.Println("HasCollection rpc failed: " + err.Error())
		return
	}
	if hasCollection == false {
		log.Println("Create collection failed: " + status.GetMessage())
		return
	}
	log.Println("Collection: " + collectionName + " exist")
}

// 删除集合
func DropCollection(collectionName string) error {
	//Drop collection
	status, _ := drivers.MilvusDbConn.DropCollection(collectionName)
	hasCollection, status1, _ := drivers.MilvusDbConn.HasCollection(collectionName)
	if !status.Ok() || !status1.Ok() || hasCollection == true {
		log.Println("Drop collection failed: " + status.GetMessage())
		return errors.New(status.GetMessage())
	}
	log.Println("Drop collection " + collectionName + " success!")

	return nil
}

//预加载集合  搜索前加载可提高速度
func LoadCollection(collectionName string) error {
	//Preload collection
	status, err := drivers.MilvusDbConn.LoadCollection(collectionName)
	if err != nil {
		log.Println("PreloadCollection rpc failed: " + err.Error())
		return err
	}
	if !status.Ok() {
		log.Println(status.GetMessage())
		return errors.New(status.GetMessage())
	}
	log.Printf("Preload collection[%s] success", collectionName)
	return nil
}

// 返回所有集合collection
func ListCollections() ([]string, error) {
	collections, status, err := drivers.MilvusDbConn.ListCollections()
	if err != nil {
		log.Println("ShowCollections rpc failed: " + err.Error())
		return nil, err
	}

	if !status.Ok() {
		log.Println("Show collections failed: " + status.GetMessage())
		return nil, errors.New(status.GetMessage())
	}

	return collections, nil
}

//获取集合信息
func GetCollectionInfo(collectionName string) (int64, int64, error){
	//test describe collection
	collectionParam, status, err := drivers.MilvusDbConn.GetCollectionInfo(collectionName)
	if err != nil {
		log.Println("DescribeCollection rpc failed: " + err.Error())
		return 0, 0, err
	}
	if !status.Ok() {
		log.Println("Create index failed: " + status.GetMessage())
		return 0, 0, errors.New(status.GetMessage())
	}

	return collectionParam.Dimension, collectionParam.IndexFileSize, nil
}

/*
创建分区
通过标签将集合分割为若干个分区，从而提高搜索效率。每个分区实际上也是一个集合。
一个集合最多创建4096个分区
每个集合都有一个 _default 分区。插入数据时如果没有指定分区，Milvus 会将数据插入该分区中。
*/
func CreatePartition(collectionName, PartitionTag string) error {

	hasCollection, status, err := drivers.MilvusDbConn.HasCollection(collectionName)
	if err != nil {
		log.Println("HasCollection rpc failed: " + err.Error())
		return err
	}

	if hasCollection == false {
		log.Println(collectionName + " 不存在!")
		return errors.New(status.GetMessage())
	}

	partitionParam := milvus.PartitionParam{collectionName, PartitionTag}
	status, err = drivers.MilvusDbConn.CreatePartition(partitionParam)
	if err != nil {
		log.Println("CreateCollection rpc failed: " + err.Error())
		return err
	}
	if !status.Ok() {
		log.Println("Create collection failed: " + status.GetMessage())
		return errors.New(status.GetMessage())
	}

	return nil
}

// 删除分区
func DropPartition(collectionName, partitionTag string) error {
	partitionParam := milvus.PartitionParam{collectionName, partitionTag}
	//Drop Partition
	status, err := drivers.MilvusDbConn.DropPartition(partitionParam)
	if err != nil {
		log.Println("DropPartition rpc failed: " + err.Error())
		return err
	}
	if !status.Ok() {
		log.Println("Create Partition failed: " + status.GetMessage())
		return errors.New(status.GetMessage())
	}

	return nil
}

// 返回所有集合collection
func ListPartitions(collectionName string) ([]string, error) {
	partitionParams, status, err := drivers.MilvusDbConn.ListPartitions(collectionName)
	if err != nil {
		log.Println("ListPartitions rpc failed: " + err.Error())
		return nil, err
	}

	if !status.Ok() {
		log.Println("Show partitions failed: " + status.GetMessage())
		return nil, errors.New(status.GetMessage())
	}

	var partitionNames []string
	for _, param := range partitionParams {
		partitionNames = append(partitionNames, param.PartitionTag)
	}
	return partitionNames, nil
}

/*
数据批量插入
单次插入的数据量不能大于 256 MB。插入数据的流程如下：
	1.服务端接收到插入请求后，将数据写入预写日志（WAL）。
	2.当预写日志成功记录后，返回插入操作。
	3.将数据写入可写缓冲区（mutable buffer）。
每个集合都有独立的可写缓冲区。每个可写缓冲区的容量上限是 128 MB。
所有集合的可写缓冲区总容量上限由系统参数 insert_buffer_size 决定，默认是 1 GB。
*/
func Insert(collectionName, partitionTag string, records []milvus.Entity) ([]int64, error) {
	insertParam := milvus.InsertParam{collectionName, partitionTag, records, nil}
	id_array, status, err := drivers.MilvusDbConn.Insert(&insertParam)
	if err != nil {
		log.Println("Insert rpc failed: " + err.Error())
		return nil, err
	}
	if !status.Ok() {
		log.Println("Insert vector failed: " + status.GetMessage())
		return nil, errors.New(status.GetMessage())
	}
	if len(id_array) != len(records) {
		log.Println("ERROR: return id array is null")
	}
	return id_array, nil
}

// 根据ID批量删除特征向量
func DeleteEntity(collectionName string, id_array []int64) error {
	status, err := drivers.MilvusDbConn.DeleteEntityByID(collectionName, id_array)
	if err != nil {
		log.Println("DeleteByID failed: " + err.Error())
		return err
	}
	if !status.Ok() {
		log.Println("DeleteByID status check error: " + status.GetMessage())
		return errors.New(status.GetMessage())
	}
	return nil
}

/*
数据落盘
落盘操作的流程如下：
	1.系统开辟一块新的可写缓冲区，用于容纳后续插入的数据。
	2.系统将之前的可写缓冲区设为只读（immutable buffer）。
	3.系统把只读缓冲区的数据写入磁盘，并将新数据段的描述信息写入元数据后端服务。
完成以上流程后，系统就成功创建了一个数据段（segment）。

自动触发:
	1.定时触发, 定时间隔由系统参数 auto_flush_interval 决定，默认是 1 秒。
	2.缓冲区达到上限触发, 累积数据达到可写缓冲区的上限（128MB）会触发落盘操作。

在调用 delete 接口后，用户可以选择再调用 flush，保证新增的数据可见，被删除的数据不会再被搜到。

为什么数据插入后不能马上被搜索到？ 因为数据还没有落盘。要确保数据插入后立刻能搜索到，可以手动调用 flush 接口。
但是频繁调用 flush 接口可能会产生大量小数据文件，从而导致查询变慢。
*/
func Flush(collectionNames []string) error {
	status, err := drivers.MilvusDbConn.Flush(collectionNames)
	if err != nil {
		log.Println("Flush error: " + err.Error())
		return err
	}
	if !status.Ok() {
		log.Println("Flush status check error: " + status.GetMessage())
		return errors.New(status.GetMessage())
	}
	return nil
}

/*
数据段整理
数据段是 Milvus 自动将插入的向量数据合并所获得的数据文件。一个集合可包含多个数据段。
如果一个数据段中的向量数据被删除，被删除的向量数据占据的空间并不会自动释放。
你可以对集合中的数据段进行 compact 操作以释放多余空间。
*/
func Compact(collectionName string) error {
	status, err := drivers.MilvusDbConn.Compact(collectionName)
	if err != nil {
		log.Println("Compact error: " + err.Error())
		return err
	}
	if !status.Ok() {
		log.Println("Compact status check error: " + status.GetMessage())
		return errors.New(status.GetMessage())
	}
	return nil
}


/*
搜索
collectionName 在哪一个集合中搜索
partitionTags 在哪集合下的哪些分区中搜索, 全部则可置为nil
queryRecords	查询向量
topk 搜索与查询向量相似度最高的前 topk 个结果
查询向量是数组, 返回的结构体也是数组, 数组内各自包含了最多 topk 个结果

返回结构体:
type TopkQueryResult struct {
	QueryResultList []QueryResult	//该数组长度对应queryRecords的长度
}

type QueryResult struct {
	// Ids id array
	Ids []int64				//该数组长度<=topk
	// Distances distance array
	Distances []float32		//该数组长度<=topk
}
查询向量数据时，你可以根据标签来指定在某个分区的数据中进行查询。Milvus 既支持对分区标签的精确匹配，也支持正则表达式匹配。
*/
func Search(collectionName string, partitionTags []string, queryRecords []milvus.Entity, topk int64) (milvus.TopkQueryResult, error) {

	extraParams := "{\"nprobe\" : 32}"	//查询取的单元数	topk:查询返回的单元数
	searchParam := milvus.SearchParam{
		collectionName,
		queryRecords,
		topk,
		partitionTags,
		extraParams}

	topkQueryResult, _, err := drivers.MilvusDbConn.Search(searchParam)
	if err != nil {
		log.Println("Search rpc failed: " + err.Error())
		return topkQueryResult, err
	}

	return topkQueryResult, nil
}

// 查询集合总数 1秒耗时
func Count(collectionName string) int64 {
	collectionCount, status, err := drivers.MilvusDbConn.CountEntities(collectionName)
	if err != nil {
		log.Println("CountCollection rpc failed: " + err.Error())
		return -1
	}
	if !status.Ok() {
		log.Println("Get collection count failed: " + status.GetMessage())
		return -2
	}
	return collectionCount
}

/*
创建索引
同时只能有一种索引, 创建时会自动把旧的索引文件删掉, 所以无需手动删除索引
索引参数看这里: https://www.milvus.io/cn/docs/v0.10.4/index.md
当插入的数据段少于 4096 行时，Milvus 不会为其建立索引。
FLAT		N/A				查询数据规模小，对查询速度要求不高。需要 100% 的召回率。
IVF_FLAT	基于量化的索引		高速查询，要求尽可能高的召回率。
IVF_SQ8		基于量化的索引		高速查询，磁盘和内存资源有限，仅有 CPU 资源。
IVF_SQ8H	基于量化的索引		高速查询，磁盘、内存、显存有限。
IVF_PQ		基于量化的索引
RNSG		基于图的索引
HNSW		基于图的索引
ANNOY		基于树的索引
*/
func CreateIndex(indexParam milvus.IndexParam) error {
	log.Println("Start create index...", indexParam)
	status, err := drivers.MilvusDbConn.CreateIndex(&indexParam)
	if err != nil {
		log.Println("CreateIndex rpc failed: " + err.Error())
		return err
	}
	if !status.Ok() {
		log.Println("Create index failed: " + status.GetMessage())
		return errors.New(status.GetMessage())
	}

	log.Println("Create index success!")
	return nil
}

//删除索引 会恢复成默认的FLAT索引
func DropIndex(collectionName string) error {
	status, err := drivers.MilvusDbConn.DropIndex(collectionName)
	if err != nil {
		log.Println("DropIndex rpc failed: " + err.Error())
		return err
	}
	if !status.Ok() {
		log.Println("Drop index failed: " + status.GetMessage())
		return errors.New(status.GetMessage())
	}
	return nil
}

//查看索引信息
func GetIndexInfo(collectionName string) (milvus.IndexParam, error) {
	//Describe index
	indexParam, status, err := drivers.MilvusDbConn.GetIndexInfo(collectionName)
	if err != nil {
		log.Println("DescribeIndex rpc failed: " + err.Error())
		return indexParam, err
	}
	if !status.Ok() {
		log.Println("Describe index failed: " + status.GetMessage())
	}
	return indexParam, nil
}










