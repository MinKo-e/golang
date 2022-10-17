package main

import (
	"github.com/sirupsen/logrus"
	"k8sManager/config"
	"k8sManager/route"
	"k8sManager/service"
)

func main() {

	service.K8s.InitClient()
	route.InitRouter()
	err := route.R.Run(config.ListenAddress)
	if err != nil {
		logrus.Error("启动Gin服务失败!" + err.Error())
	}

	//r := gin.Default()
	//r.POST("/api", TestHandler)
	//err := r.Run(":80")
	//if err != nil {
	//	logrus.Error(err.Error())
	//}

}

//func TestHandler(c *gin.Context) {
//	//podfied := &PodFied{}
//	//
//	//podfied.CName = c.PostForm("cname")
//	//podfied.PName = c.PostForm("pname")
//	//podfied.Image = c.PostForm("image")
//	//podfied.Namespace = c.PostForm("namespace")
//	//fmt.Printf("%#v", podfied)
//	podfied2 := PodFied{}
//	if err := c.Bind(&podfied2); err != nil {
//		logrus.Error(err)
//		c.JSON(500, gin.H{
//			"code": 500,
//			"msg":  err,
//		})
//		return
//	}
//	fmt.Printf("%#v", podfied2)
//	c.JSON(200, gin.H{
//		"msg":  "SuccessFully",
//		"code": 200,
//		"data": podfied2,
//	})
//}

//type PodFied struct {
//	PName              string            `form:"pname" binding:"required" json:"pname"`
//	Namespace          string            `form:"namespace" binding:"required" json:"namespace"`
//	Labels             map[string]string `form:"labels"`
//	Image              string            `form:"image" binding:"required" json:"image"`
//	CName              string            `form:"cname" binding:"required" json:"cname"`
//	NodeSelector       map[string]string `form:"node_selector"`
//	ServiceAccountName string            `form:"service_account_name"`
//	NodeName           string            `form:"node_name"`
//}
