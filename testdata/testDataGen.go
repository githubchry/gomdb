package testdata

import (
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
	Location  [2]float64 `json:"location"`
	timestamp uint64     `json:"timestamp"`
	ImgUrl    string     `json:"imgUrl"`
	Desc      Describe   `json:"desc"`
}

//生成随机数
func genRandomNum(min, max int) int {
	//设置随机数种子
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

//生成布尔值
func genRandomBool() bool {
	//设置随机数种子
	rand.Seed(time.Now().UnixNano())
	if rand.Int() % 2 != 0 {
		return true
	} else {
		return false
	}
}


//生成随机字符串
func  genRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
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
	arr[1] = color_arr[genRandomNum(0,len - 2)]
	return arr[:];
}

//生成随机字符串
func SetRandomPedestrian(p *Pedestrian){
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
	p.timestamp = uint64(genRandomNum(978278400,1609430400))
	p.ImgUrl = genRandomString(10)

	p.Desc.Face.IsWearMask = genRandomBool()
	p.Desc.Face.IsWearGlasses = genRandomBool()
	p.Desc.Face.AgeRange = int8(genRandomNum(1,100))
	p.Desc.Face.Gender = int8(genRandomNum(1,2))

	p.Desc.Body.IsWearHat = genRandomBool()
	p.Desc.Body.IsWearBackpack = genRandomBool()
	p.Desc.Body.DressColor = getRandomColorArr()

}