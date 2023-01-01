package main

import (
	"fmt"
	"video_cutter/configs"
)

func main() {
	// 从`configs/conf.yml`加载配置信息
	if err := configs.Init(); err != nil {
		fmt.Printf("load config failed, err:%v\n", err)
		return
	}
}
