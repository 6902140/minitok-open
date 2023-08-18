package main

import (
	"minitok/config"
	"minitok/log"
	"minitok/routes"
	"minitok/storage"
	"minitok/usal"

	"github.com/gin-gonic/gin"
)

func init() { //初始化项目配置
	config.LoadConfig() //加载配置信息
	log.InitLog()       //初始化日志系统
	usal.InitDatabase() //初始化数据库
	storage.InitMinio() //初始化minIO对象存储系统
	usal.RedisInit()    //初始化redis缓存
}

func main() {
	defer log.Sync()           //Sync方法用于将缓冲的日志条目刷新到底层的io.Writer接口。
	defer usal.CloseDataBase() //defer 数据库的关闭
	defer usal.CloseRedis()    //redis数据库的关闭
	rou := gin.Default()       //获得gin的Engine
	rou = routes.SetRoute(rou)
	rou.Run()
}
