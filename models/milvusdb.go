package models

import (
	"github.com/githubchry/gomdb/drivers"
	"errors"
	"github.com/milvus-io/milvus-sdk-go/milvus"
	"strconv"
)

//获取配置信息
func GetConfig() (string, error) {
	//GetConfig
	configInfo, status, _ := drivers.MilvusDbConn.GetConfig("*")
	if !status.Ok() {
		println("Get config failed: " + status.GetMessage())
		return configInfo, errors.New(status.GetMessage())
	}
	println("config: ")
	println(configInfo)
	return configInfo, nil
}

/*
创建集合
collectionName 集合名称
dimension 维度
index_file_size 自动创建索引的数据文件大小, 单位为M
metricType 距离度量方式
*/
func CreateCollection(collectionName string, dimension, indexFileSize, metricType int64) {
	//1.准备创建集合所需参数：
	collectionParam := milvus.CollectionParam{collectionName, dimension, indexFileSize, metricType}

	var hasCollection bool
	hasCollection, status, err := drivers.MilvusDbConn.HasCollection(collectionName)
	if err != nil {
		println("HasCollection rpc failed: " + err.Error())
	}
	if hasCollection == false {
		status, err = drivers.MilvusDbConn.CreateCollection(collectionParam)
		if err != nil {
			println("CreateCollection rpc failed: " + err.Error())
			return
		}
		if !status.Ok() {
			println("Create collection failed: " + status.GetMessage())
			return
		}
		println("Create collection " + collectionName + " success")
	}

	hasCollection, status, err = drivers.MilvusDbConn.HasCollection(collectionName)
	if err != nil {
		println("HasCollection rpc failed: " + err.Error())
		return
	}
	if hasCollection == false {
		println("Create collection failed: " + status.GetMessage())
		return
	}
	println("Collection: " + collectionName + " exist")
}

// 删除集合
func DropCollection(collectionName string) error {
	//Drop collection
	status, _ := drivers.MilvusDbConn.DropCollection(collectionName)
	hasCollection, status1, _ := drivers.MilvusDbConn.HasCollection(collectionName)
	if !status.Ok() || !status1.Ok() || hasCollection == true {
		println("Drop collection failed: " + status.GetMessage())
		return errors.New(status.GetMessage())
	}
	println("Drop collection " + collectionName + " success!")

	return nil
}

//预加载集合  搜索前加载可提高速度
func LoadCollection(collectionName string) error {
	//Preload collection
	status, err := drivers.MilvusDbConn.LoadCollection(collectionName)
	if err != nil {
		println("PreloadCollection rpc failed: " + err.Error())
		return err
	}
	if !status.Ok() {
		println(status.GetMessage())
		return errors.New(status.GetMessage())
	}
	println("Preload collection success")
	return nil
}

// 返回所有集合collection
func ListCollections() ([]string, error) {

	collections, status, err := drivers.MilvusDbConn.ListCollections()
	if err != nil {
		println("ShowCollections rpc failed: " + err.Error())
		return nil, err
	}

	if !status.Ok() {
		println("Show collections failed: " + status.GetMessage())
		return nil, errors.New(status.GetMessage())
	}
	println("ShowCollections: ")
	for i := 0; i < len(collections); i++ {
		println(" - " + collections[i])
	}
	return collections, nil
}

//获取集合信息
func GetCollectionInfo(collectionName string) (int64, int64, error){
	//test describe collection
	collectionParam, status, err := drivers.MilvusDbConn.GetCollectionInfo(collectionName)
	if err != nil {
		println("DescribeCollection rpc failed: " + err.Error())
		return 0, 0, err
	}
	if !status.Ok() {
		println("Create index failed: " + status.GetMessage())
		return 0, 0, errors.New(status.GetMessage())
	}

	println("CollectionName:" + collectionParam.CollectionName +
		"----Dimension:" + strconv.Itoa(int(collectionParam.Dimension)) +
		"----IndexFileSize:" + strconv.Itoa(int(collectionParam.IndexFileSize)))

	return collectionParam.Dimension, collectionParam.IndexFileSize, nil
}

/*
创建分区
通过标签将集合分割为若干个分区，从而提高搜索效率。每个分区实际上也是一个集合。

*/
func CreatePartition(collectionName, PartitionTag string) error {

	hasCollection, status, err := drivers.MilvusDbConn.HasCollection(collectionName)
	if err != nil {
		println("HasCollection rpc failed: " + err.Error())
		return err
	}

	if hasCollection == false {
		println(collectionName + " 不存在!")
		return errors.New(status.GetMessage())
	}

	partitionParam := milvus.PartitionParam{collectionName, PartitionTag}
	status, err = drivers.MilvusDbConn.CreatePartition(partitionParam)
	if err != nil {
		println("CreateCollection rpc failed: " + err.Error())
		return err
	}
	if !status.Ok() {
		println("Create collection failed: " + status.GetMessage())
		return errors.New(status.GetMessage())
	}

	return nil
}

// 删除分区
func DropPartition(collectionName, PartitionTag string) error {
	partitionParam := milvus.PartitionParam{collectionName, PartitionTag}
	//Drop Partition
	status, err := drivers.MilvusDbConn.DropPartition(partitionParam)
	if err != nil {
		println("DropPartition rpc failed: " + err.Error())
		return err
	}
	if !status.Ok() {
		println("Create Partition failed: " + status.GetMessage())
		return errors.New(status.GetMessage())
	}

	return nil
}

// 插入多个
func InsertMany(collectionName, partitionTag string, records []milvus.Entity) ([]int64, error) {

	insertParam := milvus.InsertParam{collectionName, partitionTag, records, nil}
	id_array, status, err := drivers.MilvusDbConn.Insert(&insertParam)
	if err != nil {
		println("Insert rpc failed: " + err.Error())
		return nil, err
	}
	if !status.Ok() {
		println("Insert vector failed: " + status.GetMessage())
		return nil, errors.New(status.GetMessage())
	}
	if len(id_array) != len(records) {
		println("ERROR: return id array is null")
	}
	println("Insert vectors success!")
	return id_array, nil
}

// 根据ID删除特征向量
func DeleteEntity(collectionName string, id_array []int64) error {
	status, err := drivers.MilvusDbConn.DeleteEntityByID(collectionName, id_array)
	if err != nil {
		println("DeleteByID failed: " + err.Error())
		return err
	}
	if !status.Ok() {
		println("DeleteByID status check error: " + status.GetMessage())
		return errors.New(status.GetMessage())
	}
	return nil
}

/*
数据落盘
Milvus 也会执行自动落盘。自动落盘会在固定的时间周期（1 秒）将所有现存集合的数据进行落盘操作。
在调用 delete 接口后，用户可以选择再调用 flush，保证新增的数据可见，被删除的数据不会再被搜到。

为什么数据插入后不能马上被搜索到？ 因为数据还没有落盘。要确保数据插入后立刻能搜索到，可以手动调用 flush 接口。
但是频繁调用 flush 接口可能会产生大量小数据文件，从而导致查询变慢。
*/
func Flush(collectionNames []string) error {
	status, err := drivers.MilvusDbConn.Flush(collectionNames)
	if err != nil {
		println("Flush error: " + err.Error())
		return err
	}
	if !status.Ok() {
		println("Flush status check error: " + status.GetMessage())
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
		println("Compact error: " + err.Error())
		return err
	}
	if !status.Ok() {
		println("Compact status check error: " + status.GetMessage())
		return errors.New(status.GetMessage())
	}
	return nil
}


/*
搜索
collectionName 在哪一个集合中搜索
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
*/
func Search(collectionName string, queryRecords []milvus.Entity, topk int64) (milvus.TopkQueryResult, error) {

	extraParams := "{\"nprobe\" : 32}"	//查询取的单元数	topk:查询返回的单元数
	searchParam := milvus.SearchParam{
		collectionName,
		queryRecords,
		topk,
		nil,
		extraParams}

	topkQueryResult, _, err := drivers.MilvusDbConn.Search(searchParam)
	if err != nil {
		println("Search rpc failed: " + err.Error())
		return topkQueryResult, err
	}

	return topkQueryResult, nil
}

// 查询集合总数
func Count(collectionName string) int64 {
	collectionCount, status, err := drivers.MilvusDbConn.CountEntities(collectionName)
	if err != nil {
		println("CountCollection rpc failed: " + err.Error())
		return -1
	}
	if !status.Ok() {
		println("Get collection count failed: " + status.GetMessage())
		return -2
	}
	println("Collection count:" + strconv.Itoa(int(collectionCount)))
	return collectionCount
}

// 创建索引
func CreateIndex(collectionName string, indexType milvus.IndexType) error {
	println("Start create index...", indexType)
	extraParams := "{\"nlist\" : 16384}"
	indexParam := milvus.IndexParam{collectionName, indexType, extraParams}
	status, err := drivers.MilvusDbConn.CreateIndex(&indexParam)
	if err != nil {
		println("CreateIndex rpc failed: " + err.Error())
		return err
	}
	if !status.Ok() {
		println("Create index failed: " + status.GetMessage())
		return errors.New(status.GetMessage())
	}

	println("Create index success!")
	return nil
}

//查看索引信息
func GetIndexInfo(collectionName string) (milvus.IndexParam, error) {
	//Describe index
	indexParam, status, err := drivers.MilvusDbConn.GetIndexInfo(collectionName)
	if err != nil {
		println("DescribeIndex rpc failed: " + err.Error())
		return indexParam, err
	}
	if !status.Ok() {
		println("Describe index failed: " + status.GetMessage())
	}
	println(indexParam.CollectionName + "----index type:" + strconv.Itoa(int(indexParam.IndexType)))
	return indexParam, nil
}

//删除索引
func DropIndex(collectionName string) error {

	status, err := drivers.MilvusDbConn.DropIndex(collectionName)
	if err != nil {
		println("DropIndex rpc failed: " + err.Error())
		return err
	}
	if !status.Ok() {
		println("Drop index failed: " + status.GetMessage())
		return errors.New(status.GetMessage())
	}
	return nil
}










