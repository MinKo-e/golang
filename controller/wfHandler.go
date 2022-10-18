package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"k8sManager/service"
)

func GetWFHandler(c *gin.Context) {
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
	data, err := service.Workflow.GetList(params.FilterName, params.Limit, params.Page)
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

func CreateWFHandler(c *gin.Context) {
	params := &service.WorkflowCreate{}
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
	err := service.Workflow.Create(params)
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

func DeleteWFHandler(c *gin.Context) {
	params := new(struct {
		ID int `form:"id" binding:"required"`
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
	err := service.Workflow.DeleteWorkFlow(params.ID)
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
