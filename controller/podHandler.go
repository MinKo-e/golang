package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"k8sManager/service"
)

func GetPodHandler(c *gin.Context) {
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
	data, err := service.Pod.GetPods(params.FilterName, params.Namespace, params.Limit, params.Page)
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

func GetPodLogHandler(c *gin.Context) {
	params := new(struct {
		Name          string `form:"name" binding:"required"`
		ContainerName string `form:"container_name" binding:"required"`
		Namespace     string `form:"namespace" binding:"required"`
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
	data, err := service.Pod.GetPodLog(params.Name, params.Namespace, params.ContainerName)
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

func GetPodDetailsHandler(c *gin.Context) {
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
	data, err := service.Pod.GetPodDetails(params.Name, params.Namespace)
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

func CreatePodHandler(c *gin.Context) {
	params := &service.PodFied{}
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
	err := service.Pod.CreatePod(params)
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
		"name": params.PName,
	})
}

func GetNsPodNumHandler(c *gin.Context) {
	data, err := service.Pod.GetPodNum()
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

func DeletePodHandler(c *gin.Context) {
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
	err := service.Pod.DeletePod(params.Name, params.Namespace)
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
func UpdatePodHandler(c *gin.Context) {
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
	err := service.Pod.UpdatePod(params.Namespace, params.Content)
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

func GetContainerNameHandler(c *gin.Context) {
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
	data, err := service.Pod.GetContainerName(params.Name, params.Namespace)
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
