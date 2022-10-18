package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"k8sManager/service"
)

func GetPvHandler(c *gin.Context) {
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
	data, err := service.Pv.GetPv(params.FilterName, params.Limit, params.Page)
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

func GetPVDetailsHandler(c *gin.Context) {
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
	data, err := service.Pv.GetPvDetails(params.Name)
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

func CreatePvHandler(c *gin.Context) {
	params := &service.PvFied{}
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
	err := service.Pv.CreatePv(params)
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

func GetPvNumHandler(c *gin.Context) {
	data, err := service.Pv.GetPvNum()
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

func DeletePvHandler(c *gin.Context) {
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
	err := service.Pv.DeletePv(params.Name)
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
func UpdatePvHandler(c *gin.Context) {
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
	err := service.Pv.UpdatePv(params.Content)
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

func GetPvBindHandler(c *gin.Context) {
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
	data, err := service.Pv.GetPvBind(params.Name)
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
