package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"k8sManager/service"
)

func GetNsHandler(c *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
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
	data, err := service.Ns.GetNss(params.FilterName, params.Limit, params.Page)
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

func GetNsDetailsHandler(c *gin.Context) {
	params := new(struct {
		Name string `form:"name" binding:"required"`
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
	data, err := service.Ns.GetNsDetails(params.Name)
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

func CreateNsHandler(c *gin.Context) {
	params := &service.NsFied{}
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
	err := service.Ns.CreateNs(params)
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

func GetNsNumHandler(c *gin.Context) {
	data, err := service.Ns.GetNsNum()
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

func DeleteNsHandler(c *gin.Context) {
	params := new(struct {
		Name string `form:"name" binding:"required"`
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
	err := service.Ns.DeleteNs(params.Name)
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
func UpdateNsHandler(c *gin.Context) {
	params := new(struct {
		Content string `form:"content" binding:"required"`
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
	err := service.Ns.UpdateNs(params.Content)
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
