package db

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"k8sManager/config"
	"k8sManager/model"
)

var DB *gorm.DB
var IsInit bool

func InitMysql() {
	if IsInit {
		return
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBDatabase, config.DBCharset)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		logrus.Errorln("数据库连接失败" + err.Error())
		return
	}
	err = db.AutoMigrate(&model.Workflow{})
	if err != nil {
		logrus.Error("建表失败!")
		panic(err)
	}
	DB = db

	IsInit = true

	logrus.Info("数据库初始化成功!")

}
