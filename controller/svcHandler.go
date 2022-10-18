package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"k8sManager/service"
)

func GetSvcHandler(c *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace" binding:"required"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	if err := c.Bind(params); err != nil {
		logrus.Error("Bind绑定form参数失败" + err.Error())
		c.JSON(500, gin.H{
			"code": "500",
			"data": nil,
			"msg":  err.Error(),
		})
		return
	}
	data, err := service.NetworkSvc.GetNetworkSvcs(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(500, gin.H{
			"code": "500",
			"data": nil,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"msg":  "Successfully",
		"code": 200,
		"data": data,
	})
}

func GetSvcDetailsHandler(c *gin.Context) {
	params := new(struct {
		Name      string `form:"name" binding:"required"`
		Namespace string `form:"namespace" binding:"required"`
	})
	if err := c.Bind(params); err != nil {
		logrus.Error("Bind绑定form参数失败" + err.Error())
		c.JSON(500, gin.H{
			"code": "500",
			"data": nil,
			"msg":  err.Error(),
		})
		return
	}
	data, err := service.NetworkSvc.GetNetworkSvcDetails(params.Name, params.Namespace)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(500, gin.H{
			"code": "500",
			"data": nil,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"msg":  "Successfully",
		"code": 200,
		"data": data,
	})
}

func CreateSvcHandler(c *gin.Context) {
	params := &service.NetworkSvcFied{}
	if err := c.Bind(params); err != nil {
		logrus.Error("Bind绑定form参数失败" + err.Error())
		c.JSON(500, gin.H{
			"code": "500",
			"data": nil,
			"msg":  err.Error(),
		})
		return
	}
	fmt.Printf("%#v\n", params)
	err := service.NetworkSvc.CreateNetworkSvc(params)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(500, gin.H{
			"code": "500",
			"data": nil,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"msg":  "Successfully",
		"code": 200,
		"name": params.Name,
	})
}

func GetNsSvcNumHandler(c *gin.Context) {
	data, err := service.NetworkSvc.GetNetworkSvcNum()
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(500, gin.H{
			"code": "500",
			"data": nil,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"msg":  "Successfully",
		"code": 200,
		"data": data,
	})
}

func DeleteSvcHandler(c *gin.Context) {
	params := new(struct {
		Name      string `form:"name" binding:"required"`
		Namespace string `form:"namespace" binding:"required"`
	})
	if err := c.Bind(params); err != nil {
		logrus.Error("Bind绑定form参数失败" + err.Error())
		c.JSON(500, gin.H{
			"code": "500",
			"data": nil,
			"msg":  err.Error(),
		})
		return
	}
	err := service.NetworkSvc.DeleteNetworkSvc(params.Name, params.Namespace)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(500, gin.H{
			"code": "500",
			"data": nil,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"msg":  "Successfully",
		"code": 200,
	})
}
func UpdateSvcHandler(c *gin.Context) {
	params := new(struct {
		Namespace string `form:"namespace" binding:"required"`
		Content   string `form:"content" binding:"required"`
	})
	if err := c.Bind(params); err != nil {
		logrus.Error("Bind绑定form参数失败" + err.Error())
		c.JSON(500, gin.H{
			"code": "500",
			"data": nil,
			"msg":  err.Error(),
		})
		return
	}
	err := service.NetworkSvc.UpdateNetworkSvc(params.Namespace, params.Content)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(500, gin.H{
			"code": "500",
			"data": nil,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"msg":  "Successfully",
		"code": 200,
	})
}
