// 用于初始化数据库
package usal

import (
	"fmt"
	"minitok/config"
	"minitok/log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DataBase *gorm.DB

func InitDatabase() {
	//初始化数据库链接
	var err error
	conf := config.GetConfig()
	host := conf.Mysql.Host
	port := conf.Mysql.Port
	database := conf.Mysql.Database
	username := conf.Mysql.Username
	password := conf.Mysql.Password
	//拼接dsn
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", username, password, host, port, database)

	log.Info("try to open DSN: " + args)
	//使用orm框架对数据库进行链接
	DataBase, err = gorm.Open("mysql", args)
	if err != nil {
		panic("failed to connect database ,err:" + err.Error())
	}

	log.Infof("DATABASE CONNECTED SUCCESSFULLY,USERNAME:%s,DATABASE:%s", username, database)
}

func CloseDataBase() {
	err := DataBase.Close() //调用gorm的接口对数据库的连接进行关闭
	if err != nil {
		return
	}
}

func GetDB() *gorm.DB {
	return DataBase
}
