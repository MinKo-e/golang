package main

import (
	"github.com/sirupsen/logrus"
	"k8sManager/config"
	db "k8sManager/db/mysql"
	"k8sManager/middler"
	"k8sManager/route"
	"k8sManager/service"
)

func main() {
	db.InitMysql()
	service.K8s.InitClient()
	route.R.Use(middler.Core())
	route.InitRouter()
	err := route.R.Run(config.ListenAddress)
	if err != nil {
		logrus.Error("启动Gin服务失败!" + err.Error())
	}
}
