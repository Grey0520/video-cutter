package main

import (
	"fmt"
	"video_cutter/configs"
	"video_cutter/dao"
	"video_cutter/logger"
	"video_cutter/routers"
)

func main() {
	// 从`configs/conf.yml`加载配置信息
	if err := configs.Init(); err != nil {
		fmt.Printf("load config failed, err:%v\n", err)
		return
	}
	// 自定义logger
	if err := logger.Init(configs.Conf.LogConfig, configs.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	// mysql的连接
	if err := dao.Init(configs.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer dao.Close() // 程序退出关闭数据库连接

	// 初始化路由
	r := routers.SetupRouter(configs.Conf.Mode)
	err := r.Run(fmt.Sprintf(":%d", configs.Conf.Port))
	if err != nil {
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}
}
