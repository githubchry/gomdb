package models

import (
	"context"
	"github.com/githubchry/gomdb/drivers"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//https://studygolang.com/articles/16846?fr=sidebar

// mongodb model
type Mdb struct {
	collection *mongo.Collection
}

//  new collection conn
func NewMdb(collection string) *Mdb {
	mdb := new(Mdb)
	mdb.collection = drivers.MongoDbConn.Database(drivers.MongoDbName).Collection(collection)
	return mdb
}


// 插入单个
func (m *Mdb) InsertOne(document interface{}) (insertResult *mongo.InsertOneResult) {
	insertResult, err := m.collection.InsertOne(context.TODO(), document)
	if err != nil {
		log.Fatal("??",err)
	}
	return
}

// 插入多个
func (m *Mdb) InsertMany(documents []interface{}) (insertManyResult *mongo.InsertManyResult) {
	insertManyResult, err := m.collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// 删除单个
func (m *Mdb) DeleteOne(key string, value interface{}) int64 {
	filter := bson.D{{key, value}}
	count, err := m.collection.DeleteOne(context.TODO(), filter, nil)
	if err != nil {
		log.Fatal(err)
	}
	return count.DeletedCount

}

// 删除多个
func (m *Mdb) DeleteMany(key string, value interface{}) int64 {
	filter := bson.D{{key, value}}
	count, err := m.collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	return count.DeletedCount
}

// 更新一个
func (m *Mdb) UpdateOne(key string, value interface{}, update interface{}) (updateResult *mongo.UpdateResult) {
	filter := bson.D{{key, value}}
	updateResult, err := m.collection.UpdateOne(context.TODO(), filter, update)
	//singleResult := m.collection.FindOneAndUpdate(context.TODO(), filter, update, nil)
	if err != nil {
		log.Fatal(err)
	}
	return updateResult
}

// 更新多个
func (m *Mdb) UpdateMany(key string, value interface{}, update interface{}) (updateResult *mongo.UpdateResult) {
	filter := bson.D{{key, value}}
	updateResult, err := m.collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	return updateResult
}

// 查询单个
func (m *Mdb) FindOne(key string, value interface{}) *mongo.SingleResult {
	filter := bson.D{{key, value}}
	singleResult := m.collection.FindOne(context.TODO(), filter)
	if singleResult != nil {
		//log.Println(singleResult)
	}
	return singleResult
}

// 按选项查询集合
// skip 	跳过
// limit 	读取数量
// _sort  	排序   1 倒序;-1 正序
func (m *Mdb) FindAll(skip, limit int64, _sort int) *mongo.Cursor {
	sort := bson.D{{"_id", _sort}}
	filter := bson.D{{}}

	// where
	findOptions := options.Find()
	findOptions.SetSort(sort)
	findOptions.SetLimit(limit)
	findOptions.SetSkip(skip)

	cur, err := m.collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	return cur
}

// 查询集合文档总数
func (m *Mdb) Count() int64 {
	size, _ := m.collection.EstimatedDocumentCount(context.TODO())
	return size
}

// 获取集合创建时间和编号 没什么卵用
func (m *Mdb) ParsingId(result string) (time.Time, uint64) {
	temp1 := result[:8]
	timestamp, _ := strconv.ParseInt(temp1, 16, 64)
	dateTime := time.Unix(timestamp, 0) // 这是截获情报时间 时间格式 2019-04-24 09:23:39 +0800 CST
	temp2 := result[18:]
	count, _ := strconv.ParseUint(temp2, 16, 64) // 截获情报的编号
	return dateTime, count
}

// 创建索引
func (m *Mdb) CreateIndex(keys ...string) (string, error) {
	opts := options.CreateIndexes().SetMaxTime(1000 * time.Second)
	indexView := m.collection.Indexes()
	keysDoc := bsonx.Doc{}
	// 复合索引
	for _, key := range keys {
		if strings.HasPrefix(key, "-") {
			keysDoc = keysDoc.Append(strings.TrimLeft(key, "-"), bsonx.Int32(-1))
		} else {
			keysDoc = keysDoc.Append(key, bsonx.Int32(1))
		}
	}
	// 创建索引
	result, err := indexView.CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    keysDoc,
			Options: options.Index().SetUnique(true),
		},
		opts,
	)
	if result == "" || err != nil {
		log.Print("EnsureIndex error", err)
	}
	return result, err
}


// 获取集合创建时间和编号 没什么卵用
func (m *Mdb) DeleteCollection() error {
	return m.collection.Drop(context.TODO())
}
