package testdata

import (
	"github.com/githubchry/gomdb/models"
	"log"
	"math/rand"
	"time"
)

/*
模拟描述摄像头识别到一个行人的信息

*/
type DescFace struct {
	IsWearMask  	bool	`json:"isWearMask"`
	IsWearGlasses 	bool	`json:"isWearGlasses"`
	Gender			int8	`json:"gender"`
	AgeRange		int8	`json:"ageRange"`
}

type DescBody struct {
	IsWearHat 		bool		`json:"isWearHat"`
	IsWearBackpack 	bool		`json:"isWearBackpack"`
	DressColor		[]string	`json:"dressColor"`
}

type Describe struct {
	Body DescBody `json:"body"`
	Face DescFace `json:"face"`
}

type Pedestrian struct {
	EventId	  	uint64   	`json:"eventId"`
	Location  	[2]float64 	`json:"location"`
	Timestamp 	uint64     	`json:"timestamp"`
	ImgUrl    	string     	`json:"imgUrl"`
	Desc      	Describe   	`json:"desc"`
}

var gIntCNT int64

//生成随机数
func genRandomNum(min, max int) int {
	//设置随机数种子
	rand.Seed(time.Now().UnixNano()+gIntCNT)
	return rand.Intn(max-min) + min
}

//生成随机字符串
func  genRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()+gIntCNT))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

/**
从数组中随机取一个数据出来
*/
var color_arr = []string{"red", "black", "white", "green", "blue", "purple", "yellow"}
var arr [2]string
func getRandomColorArr() []string{
	var len = len(color_arr)
	arr[0] = color_arr[genRandomNum(0,len - 1)]
	arr[1] = color_arr[genRandomNum(1,len - 1)]
	return arr[:];
}

var gbool bool
func genGBool() bool {
	gbool = !gbool
	return gbool
}

//生成随机字符串
func SetRandomPedestrian(p *Pedestrian){

	gIntCNT++;
	/*
	中国四极范围:
	135.230000
	 73.400000
	  3.520000
	 53.330000
	范围内随机生成 gps坐标
	*/
	p.Location[0] = float64(genRandomNum(73400000,135230000))/1000000
	p.Location[1] = float64(genRandomNum(3520000,53330000))/1000000
	/*
	2001-01-01 00:00:00 => 978278400
	2021-01-01 00:00:00 => 1609430400
	*/
	p.Timestamp = uint64(genRandomNum(978278400,1609430400))
	p.ImgUrl = genRandomString(10)
	p.EventId = uint64(gIntCNT);

	p.Desc.Face.IsWearMask = genGBool()
	p.Desc.Face.IsWearGlasses = genGBool()
	p.Desc.Face.AgeRange = int8(gIntCNT % 100)+1
	p.Desc.Face.Gender = int8(gIntCNT % 2)+1

	p.Desc.Body.IsWearHat = genGBool()
	p.Desc.Body.IsWearBackpack = genGBool()
	p.Desc.Body.DressColor = getRandomColorArr()
}


// 插入随机数据每次插入batch条,共插入count次 共count*batch条, 注意一次插入不要超过48M
func InsertRandomPedestrian(mgo *models.Mdb, batch, count int){
	var p Pedestrian
	pedestrians := [1000]interface{}{}
	// 2000 * 1000
	timeStart := time.Now()
	for i := 0; i < count; i++ {
		timeStart := time.Now()
		for j := 0; j < batch; j++ {
			//每次写入1000条数据
			SetRandomPedestrian(&p)
			pedestrians[j] = p
		}
		log.Printf("1000 SetRandomPedestrian need: %v\n", time.Since(timeStart))
		timeStart = time.Now()
		mgo.InsertMany(pedestrians[:])
		log.Printf("Inserted %d documents need: %v\n", len(pedestrians), time.Since(timeStart))
	}

	log.Println("Inserted 2000000 documents need:", time.Since(timeStart))
}


// 插入随机数据每次插入batch条,共插入count次 共count*batch条, 注意一次插入不要超过48M
func DeletePedestrianCollection(mgo *models.Mdb){
	err := mgo.DeleteCollection()
	if err != nil {
		log.Println("DeleteCollection:", err)
	}
}


/*
测试结果:
插入200W条, 每次1000条, 插入2000次:耗时202秒

eventid不是索引的情况下
> show dbs
chrydb  0.149GB
查询eventid字段, 查询一个不存在的实体, 首次查询1200ms左右, 后面耗时960ms左右 -- 缓存的作用
查询eventid字段, 查询findOne出来的实体(第一个), 3ms左右 -- 关乎插入时间, 默认的排序
查询eventid字段, 查询db.pedestrian.find().limit(1).sort({_id:-1})出来的实体(最后一个), 960ms左右 -- 遍历耗时

接下来对eventid字段建立索引  耗时4s左右
> show dbs
chrydb  0.174GB
查询eventid字段, 查询一个不存在的实体, 耗时3ms左右
查询eventid字段, 查询findOne出来的实体(第一个), 耗时2ms左右
查询eventid字段, 查询db.pedestrian.find().limit(1).sort({_id:-1})出来的实体(最后一个), 3ms左右

*/
