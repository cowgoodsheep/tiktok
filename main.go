package main

import (
	"fmt"
	"tiktok/config"
	"tiktok/dao"
	"tiktok/models"
	"tiktok/routers"
)

func main() {
	//数据库初始化
	dao.InitMySQL()
	//数据库迁移
	dao.DB.AutoMigrate(&models.User{}, &models.Video{}, &models.Favorite{}, &models.Comment{}, &models.Follow{}, &models.Reply{}, &models.History{})
	//关闭数据库
	defer dao.DB.Close()
	//开启路由
	r := routers.SetupRouter()
	if err := r.Run(fmt.Sprintf(":%d", config.Conf.Port)); err != nil {
		fmt.Printf("server start failed, error:%v\n", err)
	}
}
