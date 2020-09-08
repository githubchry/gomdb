package main

import (
	"github.com/githubchry/gomdb/drivers"
)

func main() {
	// 初始化连接到MongoDB数据库
	drivers.Init()

	// 断开连接
	drivers.Close()
}
